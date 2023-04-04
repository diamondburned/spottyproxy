package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"libdb.so/spottyproxy/server/internal/crudkv"
)

// Opts contains options for the API.
type Opts struct {
	LoginSecret string
}

// Mount mounts the API on the given router.
func Mount(store crudkv.BasicStore, opts Opts) http.Handler {
	sessions := NewSessionStore(store)

	r := chi.NewMux()
	r.Mount("/login", newLoginHandler(sessions, opts.LoginSecret))

	return r
}
