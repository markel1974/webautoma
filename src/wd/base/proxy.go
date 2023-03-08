package base

// Proxy specifies configuration for proxies in the browser. Set the key
// "proxy" in Capabilities to an instance of this type.
type Proxy struct {
	// Type is the type of proxy to use. This is required to be populated.
	Type ProxyType `json:"proxyType"`

	// AutoconfigURL is the URL to be used for proxy auto configuration. This is
	// required if Type is set to PAC.
	AutoconfigURL string `json:"proxyAutoconfigUrl,omitempty"`

	// The following are used when Type is set to Manual.
	//
	// Note that in Firefox, connections to localhost are not proxied by default,
	// even if a proxy is set. This can be overridden via a preference setting.
	FTP           string   `json:"ftpProxy,omitempty"`
	HTTP          string   `json:"httpProxy,omitempty"`
	SSL           string   `json:"sslProxy,omitempty"`
	SOCKS         string   `json:"socksProxy,omitempty"`
	SOCKSVersion  int      `json:"socksVersion,omitempty"`
	SOCKSUsername string   `json:"socksUsername,omitempty"`
	SOCKSPassword string   `json:"socksPassword,omitempty"`
	NoProxy       []string `json:"noProxy,omitempty"`

	// The W3C draft spec includes port fields as well. According to the
	// specification, ports can also be included in the above addresses. However,
	// in the Geckodriver implementation, the ports must be specified by these
	// additional fields.
	HTTPPort  int `json:"httpProxyPort,omitempty"`
	SSLPort   int `json:"sslProxyPort,omitempty"`
	SocksPort int `json:"socksProxyPort,omitempty"`
}

// ProxyType is an enumeration of the types of proxies available.
type ProxyType string

const (
	// Direct connection - no proxy in use.
	Direct ProxyType = "direct"
	// Manual proxy settings configured, e.g. setting a proxy for HTTP, a proxy
	// for FTP, etc.
	Manual = "manual"
	// Autodetect proxy, probably with WPAD
	Autodetect = "autodetect"
	// System settings used.
	System = "system"
	// PAC - Proxy autoconfiguration from a URL.
	PAC = "pac"
)
