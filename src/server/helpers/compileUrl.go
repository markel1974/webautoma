package helpers

import "strings"

func CompileUrl(listen string, domain string, tls bool) string {
	proto := "https"
	host := ""
	port := ""
	if !tls {
		proto = "http"
	}
	if hp := strings.Split(listen, ":"); len(hp) > 0 {
		host = hp[0]
		if len(hp) > 1 {
			port = ":" + hp[1]
		}
	}
	if len(domain) > 0 {
		host = domain
	}
	if len(host) == 0 {
		host = "*"
	}
	return proto + "://" + host + port
}
