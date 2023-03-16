package gatewayserver

import (
	"net/http"

	"upspin.io/config"
	"upspin.io/flags"
	"upspin.io/log"

	// TODO: Which of these are actually needed?

	// Load useful packers
	_ "upspin.io/pack/ee"
	_ "upspin.io/pack/eeintegrity"
	_ "upspin.io/pack/plain"

	// Load required transports
	_ "upspin.io/transports"
)

func Main() (ready chan<- struct{}) {
	flags.Parse(flags.Server)
	// Load configuration and keys for this server. It needs a real upspin username and keys.
	cfg, err := config.FromFile(flags.Config)
	if err != nil {
		log.Fatal(err)
	}

	httpDir, err := NewGateway(cfg)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", httpDir)

	return ready
}
