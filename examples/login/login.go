// Trustmaster login example
//
// If you don't have an idea what trustmaster is, this is most probably not for you, sorry ;)
//
// In order to try this example you need to have an approved trustmaster account and need to apply for api credentials.
//
//
package main // udico.de/trustmaster/examples/login

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

func main() {
	log.SetLevel(log.InfoLevel) // you could set to DebugLevel to see raw http responses in the logs for debugging

  	tm, err := trustmaster.NewClient(TrustmasterClientId, TrustmasterClientSecret, TrustmasterRedirectUri)
	if err != nil {
		log.WithError(err).Fatal("cannot get trustmaster instance")
	}

	// First, create the url for the user to authenticate
	// chose any value for the state parameter which allows you to identify the user within your callback function.
	// Select the required scopes as well. Technically you need none of those for just the login to work. However,
	// you should check the trust of the user after the login.
	log.Info(tm.AuthURL(
		"test",
		trustmaster.GetAgentName,
		trustmaster.GetEmailAddress,
		trustmaster.GetGoogleName,
		trustmaster.GetTelegram,
	))

	// Create a callback handler for oauth requests.
	// The default handler forwards logins to the given channel and redirects the user to a website of your choice
	handler := trustmaster.CallbackHandler{
		Events:   make(chan trustmaster.AuthEvent),
		Redirect: "https://google.com/",
	}

	// Next, spin up the event listener …
	go func() {
		for ev := range handler.Events {
			// Exchange the code for a user access token
			usertoken, err := tm.UserAccessToken(ev.Code)
			if err != nil {
				log.WithError(err).Error("cannot not exchange code for access token")
				continue
			}

			// Make sure to store the token!
			// Having the user to authenticate each time is a bad practice.
			//
			// Note: Access tokens refresh automatically, once they're used after expiration.
			//       You should always save them.
			// You can also refresh a token manually:
			//     if usertoken.Expired() {
			//       tm.RefreshToken(usertoken)
			//       // check for errors and save
			//     }

			// Read the users Profile
			prof, err := tm.Profile(usertoken)
			if err != nil {
				log.WithError(err).Error("cannot read profile")
				continue
			}
			log.Infof("%+v", prof)

		}
	}()

	// Finally, spin up an http endpoint
	hSrv := &http.Server{Addr: ":80", Handler: nil}
	http.Handle("/oauth/callback", handler)
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
