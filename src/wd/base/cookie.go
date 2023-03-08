package base

import "strings"

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
