package gatewayserver

import (
	"io"
	"net/http"
	"strings"

	"upspin.io/client"
	"upspin.io/errors"
	"upspin.io/upspin"
)

type Gateway struct {
	c   upspin.Client
	cfg upspin.Config
	mux *http.ServeMux
}

func NewGateway(cfg upspin.Config, options ...string) (*Gateway, error) {
	o := &Gateway{}
	const op errors.Op = "cmd/upspin-gateway/gatewayserver.New"
	if cfg == nil {
		return nil, errors.E(op, errors.Invalid, "nil config")
	}
	o.c = client.New(cfg)
	o.cfg = cfg
	o.route()
	return o, nil
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}
func (g *Gateway) route() {
	g.mux = http.NewServeMux()
	g.mux.Handle("/raw/", http.StripPrefix("/raw", http.HandlerFunc(g.handleRaw)))
}

func (g *Gateway) handleRaw(w http.ResponseWriter, r *http.Request) {
	names := g.expandUpspin([]string{strings.TrimPrefix(r.URL.Path, "/")}, r.URL.Query().Get("glob") == "true")
	if len(names) == 0 {
		http.Error(w, "File not Found", 404)
		return
	}
	file, err := g.c.Open(names[0])
	if err != nil {
		http.Error(w, "File not Found", 404)
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "File not Found", 404)
		return
	}
}
