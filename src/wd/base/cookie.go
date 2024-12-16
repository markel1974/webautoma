package base

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

/*
// ChromeDriver returns the expiration date as a float. Handle both formats
// via a type switch.
type cookie struct {
	Name     string      `json:"name"`
	Value    string      `json:"value"`
	Path     string      `json:"path"`
	Domain   string      `json:"domain"`
	Secure   bool        `json:"secure"`
	Expiry   interface{} `json:"expiry"`
	HTTPOnly bool        `json:"httpOnly"`
	SameSite string      `json:"sameSite",omitempty`
}
*/

// Cookie represents an HTTP cookie.
type Cookie struct {
	Name     string      `json:"name"`
	Value    string      `json:"value"`
	Path     string      `json:"path"`
	Domain   string      `json:"domain"`
	Secure   bool        `json:"secure"`
	Expiry   interface{} `json:"expiry"`
	HTTPOnly bool        `json:"httpOnly"`
	SameSite string      `json:"sameSite,omitempty"`
}

const (
	SameSiteNone   = "None"
	SameSiteLax    = "Lax"
	SameSiteStrict = "Strict"
	SameSiteEmpty  = ""
)

func (c Cookie) Sanitize() Cookie {
	parseExpiry := func(e interface{}) uint {
		switch expiry := c.Expiry.(type) {
		case int:
			if expiry > 0 {
				return uint(expiry)
			}
		case float64:
			return uint(expiry)
		}
		return 0
	}

	parseSameSite := func(s string) string {
		if s == "" {
			return ""
		}
		for _, v := range []string{SameSiteNone, SameSiteLax, SameSiteStrict} {
			if strings.EqualFold(v, s) {
				return v
			}
		}
		return SameSiteLax
	}

	return Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		Secure:   c.Secure,
		Expiry:   parseExpiry(c.Expiry),
		HTTPOnly: c.HTTPOnly,
		SameSite: parseSameSite(c.SameSite),
	}
}

func (c Cookie) Http() *http.Cookie {
	expire, _ := c.Expiry.(uint)
	sameSite := http.SameSiteDefaultMode
	switch c.SameSite {
	case SameSiteNone:
		sameSite = http.SameSiteNoneMode
	case SameSiteLax:
		sameSite = http.SameSiteLaxMode
	case SameSiteStrict:
		sameSite = http.SameSiteStrictMode
	}
	out := &http.Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		Secure:   c.Secure,
		Expires:  time.Unix(int64(expire), 0),
		HttpOnly: c.HTTPOnly,
		SameSite: sameSite,
	}
	return out
}

func Cookies2Json(cookies []Cookie) (string, error) {
	src, err := json.Marshal(cookies)
	if err != nil {
		return "", err
	}
	return string(src), nil
}

func Cookies2Curl(cookies []Cookie) (string, error) {
	//string example.com - the domain name
	//boolean FALSE - include subdomains
	//string /foobar/ - path
	//boolean TRUE - send/receive over HTTPS only
	//number 1462299217 - expires at - seconds since Jan 1st 1970, or 0
	//string person - name of the cookie
	//string daniel - value of the cookie
	var out []byte
	for _, c := range cookies {
		httpOnly := "FALSE"
		if c.HTTPOnly {
			httpOnly = "TRUE"
		}
		//expire := strconv.FormatUint(uint64(c.Expire()), 10)
		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%v\t%s\t%s", c.Domain, "FALSE", c.Path, httpOnly, c.Expiry, c.Name, c.Value)
		out = append(out, []byte(line)...)
		out = append(out, '\n')
	}
	return string(out), nil
}
