package main

import (
	"context"
	"flag"
	"os"

	"libdb.so/hserve"
	"libdb.so/spottyproxy/server/api"
)

var (
	addr = ":8080"
)

func main() {
	flag.StringVar(&addr, "addr", addr, "address to listen on")
	flag.Parse()

	handler := api.Mount(nil, api.Opts{
		LoginSecret: os.Getenv("APP_LOGIN_SECRET"),
	})

	hserve.MustListenAndServe(context.Background(), addr, handler)
}
