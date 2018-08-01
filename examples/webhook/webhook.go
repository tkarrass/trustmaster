package main

import (
	"udico.de/trustmaster"
	log "github.com/sirupsen/logrus"
	"net/http"
	"udico.de/terminator"
	"context"
)

const TrustmasterClientId = "<your own trustmaster client id here>"
const TrustmasterClientSecret = "<your own trustmaster client secret here>"

// This has to exactly match the uri entered in your api client request, in order for the oauth login to work!
const TrustmasterRedirectUri = "https://example.com/oauth/callback"

const TrustmasterWebhookSecret = "<...>"

func main() {

	tm, err := trustmaster.NewClient(
		TrustmasterClientId,
		TrustmasterClientSecret,
		TrustmasterRedirectUri,
	)
	if err != nil {
		log.WithError(err).Fatal("cannot get trustmaster instance")
	}

	gtoken, err := tm.GenericAccessToken()
	if err != nil {
		log.WithError(err).Fatal("cannot get generic access token")
	}

	handler := trustmaster.WebhookHandler{
		Events: make(chan trustmaster.WebhookEvent),
		WebhookSecret: TrustmasterWebhookSecret,
	}

	go func() {
		for ev := range handler.Events {
			log.Infof("Trust status changed for user %v", ev.GoogleId)
			trust, err := tm.Trust(ev.GoogleId, gtoken)
			if err != nil {
				log.WithError(err).Error("cannot get trust information")
				continue
			}
			log.Infof("%+v", trust)
		}
	}()

	hSrv := &http.Server{Addr: ":80", Handler: nil}
	http.Handle("/hook", handler)
	closer := terminator.Terminator
	go func() {
		if err := hSrv.ListenAndServe(); err != nil {
			select {
			case _, opn := <-closer:
				if opn {
					log.WithError(err).Error("content in close channel")
				} else {
					log.Info("shutting down")
				}
			default:
				log.WithError(err).Fatal("unexpected failure")
			}
		}
	}()

	// wait for ctrl-c/sigint/sigterm/... and shutdown the server gracefully
	<-closer
	hSrv.Shutdown(context.Background())
	close(handler.Events)
}