package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/boltdb/bolt"
	"github.com/powerman/narada-go/narada/bootstrap"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/qarea/ctxtg"
	"github.com/qarea/jirams/api"
	"github.com/qarea/jirams/cfg"
	"github.com/qarea/jirams/jira"
	"github.com/qarea/jirams/store"
)

type appParams struct {
	BoltDB    *bolt.DB
	PublicKey []byte
}

func start(params appParams) {
	var (
		httpListener net.Listener
	)
	userStore := store.New(params.BoltDB)
	jiraClient := jira.NewClient(userStore)
	tokenParser, err := ctxtg.NewRSATokenParser(params.PublicKey)
	if err != nil {
		panic(err)
	}
	rpcInterface := api.NewRPCAPI(jiraClient, tokenParser)
	rpc.Register(rpcInterface)

	http.Handle(cfg.HTTP.BasePath+"/rpc", jsonrpc2.HTTPHandler(nil))
	httpListener, err = net.Listen("tcp", cfg.HTTP.Listen)
	if err != nil {
		panic(err)
	}
	defer httpListener.Close()

	if err := bootstrap.Unlock(); err != nil {
		log.Fatal(err)
	}

	l.NOTICE("Listening on %s", httpListener.Addr().String())
	l.Fatal(http.Serve(httpListener, nil))
}
