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

// NewOAuth2BearerClient An implementation of the OAUTHBEARER authentication mechanism, as described in RFC 7628.
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
	ir := b64.StdEncoding.EncodeToString([]byte(str))
	return OAuth2Bearer, []byte(ir), nil
}

func (a *Oauth2BearerClient) Next(challenge []byte) ([]byte, error) {
	authBearerErr := &OAuth2BearerError{}
	if err := json.Unmarshal(challenge, authBearerErr); err != nil {
		return nil, err
	}
	return nil, authBearerErr
}
