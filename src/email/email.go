package email

import (
	"fmt"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
	"io"
	"strings"
	"time"
)

//https://github.com/emersion/go-imap/blob/v2/imapclient/example_test.go

const defaultFolder = "INBOX"

type Mode int

const (
	ModeInsecure Mode = 0
	ModeStartTLS Mode = 1
	ModeTLS      Mode = 2
)

func (m Mode) String() string {
	switch m {
	case ModeInsecure:
		return "insecure"
	case ModeStartTLS:
		return "startTLS"
	case ModeTLS:
		return "TLS"
	}
	return ""
}

type Client struct {
	mode      Mode
	server    string
	user      string
	password  string
	folder    string
	options   *imapclient.Options
	useOAuth2 bool
}

func NewClientFromTarget(target string) (*Client, error) {
	const protoSep = "://"
	const serverSep = "@"
	const userSep = ":"
	const oauth2Mode = "[oauth2]"
	pos := strings.LastIndex(target, serverSep)
	if pos <= 0 {
		return nil, fmt.Errorf("invalid target, missing server section")
	}
	p1 := target[:pos]
	server := target[pos+len(serverSep):]
	pos = strings.Index(target, protoSep)
	if pos <= 0 {
		return nil, fmt.Errorf("invalid target, missing protocol section")
	}
	k := p1[:pos]
	p2 := p1[pos+len(protoSep):]
	pos = strings.Index(p2, userSep)
	if pos <= 0 {
		return nil, fmt.Errorf("invalid target, missing user section")
	}
	user := p2[:pos]
	password := p2[pos+len(userSep):]
	mode := ModeTLS
	useOAuth2 := false
	if strings.Contains(k, oauth2Mode) {
		k = strings.Replace(k, oauth2Mode, "", -1)
		useOAuth2 = true
	}
	switch strings.ToLower(k) {
	case "insecure":
		mode = ModeInsecure
	case "starttls":
		mode = ModeStartTLS
	case "tls":
		mode = ModeTLS
	default:
		return nil, fmt.Errorf("invalid target, unknown mode %s", k)
	}
	return NewClient(server, user, password, mode, useOAuth2), nil
}

func NewClient(server string, user string, password string, mode Mode, useOAuth2 bool) *Client {
	return &Client{
		server:    server,
		user:      user,
		password:  password,
		mode:      mode,
		folder:    defaultFolder,
		options:   nil,
		useOAuth2: useOAuth2,
	}
}

func (cl *Client) SetFolder(folder string) {
	cl.folder = folder
}

func (cl *Client) Retrieve(opt *Options) (string, error) {
	var c *imapclient.Client
	var err error
	switch cl.mode {
	case ModeInsecure:
		c, err = imapclient.DialInsecure(cl.server, cl.options)
	case ModeStartTLS:
		c, err = imapclient.DialStartTLS(cl.server, cl.options)
	case ModeTLS:
		c, err = imapclient.DialTLS(cl.server, cl.options)
	}
	if err != nil {
		return "", err
	}
	if c == nil {
		return "", fmt.Errorf("unsupported mode %v", cl.mode.String())
	}
	defer c.Close()
	if cl.useOAuth2 {
		saslClient := NewOAuthBearerClient(&OAuth2BearerOptions{Username: cl.user, Token: cl.password})
		if err = c.Authenticate(saslClient); err != nil {
			return "", err
		}
	} else {
		if err = c.Login(cl.user, cl.password).Wait(); err != nil {
			return "", err
		}
	}
	mailboxes, err := c.List("", "%", nil).Collect()
	if err != nil {
		return "", err
	}
	for _, mbox := range mailboxes {
		fmt.Printf(" - %v\n", mbox.Mailbox)
	}
	selectedMbox, err := c.Select(cl.folder, nil).Wait()
	if err != nil {
		return "", err
	}
	seqNum := selectedMbox.NumMessages
	if seqNum == 0 {
		return "", fmt.Errorf("empty mailbox")
	}
	fmt.Printf("%s contains %v messages\n", cl.folder, selectedMbox.NumMessages)
	var data string
	if selectedMbox.NumMessages > 0 {
		fetchOptions := &imap.FetchOptions{
			Envelope: true,
		}
		minTime := time.Now().Add(-time.Duration(opt.validity) * time.Minute)
		for {
			seqSet := imap.SeqSetNum(seqNum)
			var messages []*imapclient.FetchMessageBuffer
			messages, err = c.Fetch(seqSet, fetchOptions).Collect()
			if err != nil {
				err = fmt.Errorf("failed to fetch message: %v", err)
				break
			}
			if len(messages) > 0 {
				msg := messages[0]
				if g1 := msg.Envelope.Date.Before(minTime); g1 {
					err = fmt.Errorf("invalid date")
					break
				}
				if opt.subjectRgx.Match([]byte(msg.Envelope.Subject)) {
					var body []byte
					if body, err = cl.fetchBody(c, messages[0].SeqNum); err == nil {
						tmp := strings.Replace(string(body), "\n", " ", -1)
						tmp = strings.Replace(tmp, "\r", " ", -1)
						if v := opt.bodyRgx.FindStringSubmatch(tmp); len(v) > 0 {
							err = nil
							data = v[1]
							break
						}
					} else {
						//err =
						//log.Fatalf("failed to fetch first message body in INBOX: %v", err)
					}
				}
			}
			fmt.Printf("not found %d\n", seqNum)
			seqNum--
			if seqNum <= 0 {
				err = fmt.Errorf("")
				break
			}
		}
	}
	_ = c.Logout().Wait()
	if err == nil && len(data) == 0 {
		err = fmt.Errorf("not found")
	}
	return data, err
}

func (cl *Client) fetchBody(c *imapclient.Client, seqNum uint32) ([]byte, error) {
	seqSet := imap.SeqSetNum(seqNum)
	fetchOptions := &imap.FetchOptions{BodySection: []*imap.FetchItemBodySection{{}}}
	fetchCmd := c.Fetch(seqSet, fetchOptions)
	defer fetchCmd.Close()

	msg := fetchCmd.Next()
	if msg == nil {
		return nil, fmt.Errorf("FETCH command did not return any message")
	}

	var bs imapclient.FetchItemDataBodySection
	ok := false
	for {
		item := msg.Next()
		if item == nil {
			break
		}
		bs, ok = item.(imapclient.FetchItemDataBodySection)
		if ok {
			break
		}
	}
	if !ok {
		return nil, fmt.Errorf("FETCH command did not return body section")
	}
	mr, err := mail.CreateReader(bs.Literal)
	if mr == nil {
		if err != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf("nil reader")
		}
	}
	//if err != nil {
	//	return nil, err
	//}

	/*
		// Print a few header fields
		h := mr.Header
		if date, err := h.Date(); err != nil {
			log.Printf("failed to parse Date header field: %v", err)
		} else {
			log.Printf("Date: %v", date)
		}
		if to, err := h.AddressList("To"); err != nil {
			log.Printf("failed to parse To header field: %v", err)
		} else {
			log.Printf("To: %v", to)
		}
		if subject, err := h.Text("Subject"); err != nil {
			log.Printf("failed to parse Subject header field: %v", err)
		} else {
			log.Printf("Subject: %v", subject)
		}
	*/

	var parts []byte
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		switch p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, err := io.ReadAll(p.Body)
			if err != nil {
				return nil, err
			}
			parts = append(parts, b...)
			return b, nil
			//case *mail.AttachmentHeader:
			//	filename, _ := h.Filename()
			//	log.Printf("Attachment: %v", filename)
			//}
		}
	}
	_ = fetchCmd.Close()
	return parts, err
}
