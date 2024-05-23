package config

type Config struct {
	TLS     bool
	CA      string
	Cert    string
	Key     string
	AssetUI string
	Listen  string
	Domain  string
	UIPath  string
}
