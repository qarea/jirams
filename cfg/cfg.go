package cfg

import (
	"strings"
	"time"

	"github.com/powerman/narada-go/narada"
)

var log = narada.NewLog("")

// Global configuration variables
var (
	Debug        bool
	LockTimeout  time.Duration
	RSAPublicKey []byte
	MySQL        struct {
		Host     string
		Port     int
		DB       string
		Login    string
		Password string
	}
	HTTP struct {
		Listen       string
		BasePath     string
		RealIPHeader string
	}
)

func init() {
	if err := load(); err != nil {
		log.Fatal(err)
	}
}

func load() error {
	Debug = narada.GetConfigLine("log/level") == "DEBUG"

	HTTP.Listen = narada.GetConfigLine("http/listen")
	if !strings.Contains(HTTP.Listen, ":") {
		log.Fatal("please setup config/http/listen")
	}

	HTTP.BasePath = narada.GetConfigLine("http/basepath")
	if HTTP.BasePath != "" && (HTTP.BasePath[0] != '/' || HTTP.BasePath[len(HTTP.BasePath)-1] == '/') {
		log.Fatal("config/http/basepath should begin with / and should not end with /")
	}

	HTTP.RealIPHeader = narada.GetConfigLine("http/real_ip_header")

	var err error
	RSAPublicKey, err = narada.GetConfig("rsa_public_key")
	if err != nil {
		return err
	}

	LockTimeout = narada.GetConfigDuration("lock_timeout")
	return nil
}
