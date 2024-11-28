package email

import (
	"bytes"
	"fmt"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
	"io"
	"log"
	"regexp"
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
	mode       Mode
	server     string
	user       string
	password   string
	folder     string
	options    *imapclient.Options
	subjectRgx *regexp.Regexp
	bodyRgx    *regexp.Regexp
}

func NewClient(server string, user string, password string, mode Mode) *Client {
	subjectRgx := regexp.MustCompile("test .*")
	bodyRgx := regexp.MustCompile("<beta>([^>]+)<beta>")
	return &Client{
		server:     server,
		user:       user,
		password:   password,
		mode:       mode,
		folder:     defaultFolder,
		options:    nil,
		subjectRgx: subjectRgx,
		bodyRgx:    bodyRgx,
	}
}

func (cl *Client) SetFolder(folder string) {
	cl.folder = folder
}

func (cl *Client) Retrieve() (string, error) {
	//"mail.example.org:993"
	//"root", "asdf"
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

	if err = c.Login(cl.user, cl.password).Wait(); err != nil {
		return "", err
	}

	mailboxes, err := c.List("", "%", nil).Collect()
	if err != nil {
		return "", err
	}
	for _, mbox := range mailboxes {
		fmt.Printf(" - %v", mbox.Mailbox)
	}

	selectedMbox, err := c.Select(cl.folder, nil).Wait()
	if err != nil {
		return "", err
	}
	fmt.Printf("%s contains %v messages", cl.folder, selectedMbox.NumMessages)

	var data string

	if selectedMbox.NumMessages > 0 {
		seqSet := imap.SeqSetNum(1)
		fetchOptions := &imap.FetchOptions{Envelope: true}
		messages, err := c.Fetch(seqSet, fetchOptions).Collect()
		if err != nil {
			log.Fatalf("failed to fetch first message in INBOX: %v", err)
		}
		if cl.subjectRgx.Match([]byte(messages[0].Envelope.Subject)) {
			if body, err := cl.fetchBody(c, messages[0].SeqNum); err == nil {
				if v := cl.bodyRgx.FindSubmatch(body); len(v) > 0 {
					data = string(bytes.Join(v, []byte{' '}))
				}
			}
		}
	}
	if err = c.Logout().Wait(); err != nil {
		return "", err
	}
	return data, nil
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
	if err != nil {
		return nil, err
	}

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
