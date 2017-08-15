package config

// Config is the structure for serving the acme-reverseproxy
type Config struct {
	CA      CA
	Mapping map[string]string
}

// CA provides configuration for passing to the certificate authority
type CA struct {
	LetsEncryptURL string
	CacheDir       string // Path to stash certificates
	Email          string // email to register with acme provider (let's encrypt)
}
