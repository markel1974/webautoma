package email

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

//https://learn.microsoft.com/en-us/exchange/client-developer/legacy-protocols/how-to-authenticate-an-imap-pop-smtp-application-by-using-oauth

const OAuth2Bearer = "XOAUTH2"

var ErrUnexpectedClientResponse = errors.New("sasl: unexpected client response")

type OAuth2BearerError struct {
	Status  string `json:"status"`
	Schemes string `json:"schemes"`
	Scope   string `json:"scope"`
}

type OAuth2BearerOptions struct {
	Username string
	Token    string
	Host     string
	Port     int
}

// Implements error
func (err *OAuth2BearerError) Error() string {
	return fmt.Sprintf("OAUTHBEARER authentication error (%v)", err.Status)
}

type Oauth2BearerClient struct {
	OAuth2BearerOptions
}

// NewOAuth2BearerClient An implementation of the OAUTHBEARER authentication mechanism, as
// described in RFC 7628.
func NewOAuth2BearerClient(opt *OAuth2BearerOptions) *Oauth2BearerClient {
	return &Oauth2BearerClient{*opt}
}

func (a *Oauth2BearerClient) Start() (string, []byte, error) {
	const sep = "\x01"
	//var auth string
	//if a.Username != "" {
	auth := "user=" + a.Username
	str := auth
	//}
	//str := "n," + auth + ","
	if a.Host != "" {
		str += sep + "host=" + a.Host
	}

	if a.Port != 0 {
		str += sep + "port=" + strconv.Itoa(a.Port)
	}
	str += sep + "auth=Bearer " + a.Token + sep + sep
	//ir := []byte(str)
	ir := b64.StdEncoding.EncodeToString([]byte(str))
	//fmt.Println(sEnc)
	return OAuth2Bearer, []byte(ir), nil
}

func (a *Oauth2BearerClient) Next(challenge []byte) ([]byte, error) {
	authBearerErr := &OAuth2BearerError{}
	if err := json.Unmarshal(challenge, authBearerErr); err != nil {
		return nil, err
	} else {
		return nil, authBearerErr
	}
}

/*

type OAuthBearerAuthenticator func(opts OAuth2BearerOptions) *OAuth2BearerError

type oauthBearerServer struct {
	done         bool
	failErr      error
	authenticate OAuthBearerAuthenticator
}

func (a *oauthBearerServer) fail(desc string) ([]byte, bool, error) {
	blob, err := json.Marshal(OAuth2BearerError{
		Status:  "invalid_request",
		Schemes: "bearer",
	})
	if err != nil {
		panic(err) // wtf
	}
	a.failErr = errors.New("sasl: client error: " + desc)
	return blob, false, nil
}

func (a *oauthBearerServer) Next(response []byte) (challenge []byte, done bool, err error) {
	// Per RFC, we cannot just send an error, we need to return JSON-structured
	// value as a challenge and then after getting dummy response from the
	// client stop the exchange.
	if a.failErr != nil {
		// Server libraries (go-smtp, go-imap) will not call Next on
		// protocol-specific SASL cancel response ('*'). However, GS2 (and
		// indirectly OAUTHBEARER) defines a protocol-independent way to do so
		// using 0x01.
		if len(response) != 1 && response[0] != 0x01 {
			return nil, true, errors.New("sasl: invalid response")
		}
		return nil, true, a.failErr
	}

	if a.done {
		err = ErrUnexpectedClientResponse
		return
	}

	// Generate empty challenge.
	if response == nil {
		return []byte{}, false, nil
	}

	a.done = true

	// Cut n,a=username,\x01host=...\x01auth=...
	// into
	//   n
	//   a=username
	//   \x01host=...\x01auth=...\x01\x01
	parts := bytes.SplitN(response, []byte{','}, 3)
	if len(parts) != 3 {
		return a.fail("Invalid response")
	}
	flag := parts[0]
	authId := parts[1]
	if !bytes.Equal(flag, []byte{'n'}) {
		return a.fail("Invalid response, missing 'n' in gs2-cb-flag")
	}
	opts := OAuth2BearerOptions{}
	if len(authId) > 0 {
		if !bytes.HasPrefix(authId, []byte("a=")) {
			return a.fail("Invalid response, missing 'a=' in gs2-authzid")
		}
		opts.Username = string(bytes.TrimPrefix(authId, []byte("a=")))
	}

	// Cut \x01host=...\x01auth=...\x01\x01
	// into
	//   *empty*
	//   host=...
	//   auth=...
	//   *empty*
	//
	// Note that this code does not do a lot of checks to make sure the input
	// follows the exact format specified by RFC.
	params := bytes.Split(parts[2], []byte{0x01})
	for _, p := range params {
		// Skip empty fields (one at start and end).
		if len(p) == 0 {
			continue
		}

		pParts := bytes.SplitN(p, []byte{'='}, 2)
		if len(pParts) != 2 {
			return a.fail("Invalid response, missing '='")
		}

		switch string(pParts[0]) {
		case "host":
			opts.Host = string(pParts[1])
		case "port":
			port, err := strconv.ParseUint(string(pParts[1]), 10, 16)
			if err != nil {
				return a.fail("Invalid response, malformed 'port' value")
			}
			opts.Port = int(port)
		case "auth":
			const prefix = "bearer "
			strValue := string(pParts[1])
			// Token type is case-insensitive.
			if !strings.HasPrefix(strings.ToLower(strValue), prefix) {
				return a.fail("Unsupported token type")
			}
			opts.Token = strValue[len(prefix):]
		default:
			return a.fail("Invalid response, unknown parameter: " + string(pParts[0]))
		}
	}

	authErr := a.authenticate(opts)
	if authErr != nil {
		blob, err := json.Marshal(authErr)
		if err != nil {
			panic(err) // wtf
		}
		a.failErr = authErr
		return blob, false, nil
	}

	return nil, true, nil
}

*/
