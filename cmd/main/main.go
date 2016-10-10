// Package main provides main service.
package main

import (
	"github.com/powerman/narada-go/narada/bootstrap"
	"gitlab.qarea.org/tgms/ctxtg"

	"net/http"

	"github.com/powerman/narada-go/narada"
	"github.com/prometheus/client_golang/prometheus"

	"../../api/rpcsvc"
	"../../cfg"
	"../../tracker"
)

var log = narada.NewLog("")

func main() {
	p, err := ctxtg.NewRSATokenParser(cfg.RSAPublicKey)
	if err != nil {
		log.Fatal(err)
	}

	tr := tracker.NewClient()
	rpcsvc.Init(tr, p)

	if err := bootstrap.Unlock(); err != nil {
		log.Fatal(err)
	}
	http.Handle(cfg.HTTP.BasePath+"/metrics", prometheus.Handler())
	log.NOTICE("Listening on %s", cfg.HTTP.Listen+cfg.HTTP.BasePath)
	log.Fatal(http.ListenAndServe(cfg.HTTP.Listen, nil))
}
