package trustmaster

import (
	"net/http"
	"fmt"
	"time"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"crypto/hmac"
	"crypto/sha256"
	"strings"
	"encoding/hex"
	"encoding/json"
)

type AuthEvent struct {
	Code string
	Tag string
}

type CallbackHandler struct {
	Events chan AuthEvent
	Redirect string  // redirect to this target after reading the code and tag
	Error string  // redirect here, if any error occured. will redirect to Redirect if empty
}

// Creates a http callback handler which reads the code and additional tag
// from the request and forwards to a given redirect target.
func (handler CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get code and tag
	code := r.URL.Query().Get("code")
	tag := r.URL.Query().Get("state")
	target := handler.Redirect
	if code == "" || tag == "" {
		if handler.Error != "" {
			target = handler.Error
		}
	} else {
		handler.Events <- AuthEvent{code, tag}
	}

	// redirect
	w.Header().Set("Location", target)
	w.WriteHeader(307)
	fmt.Fprintf(w, "<html><body><h3>moved</h3></body></html>")
}


type WebhookEvent struct {
	GoogleId string
	Timestamp time.Time
	Tag string
}

type WebhookHandler struct {
	Events chan WebhookEvent
	WebhookSecret string
}

func (webhook WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	xsignature := r.Header.Get("X-Signature")
	if !strings.HasPrefix(xsignature, "sha256=") {
		log.Error("invalid hash algorithm")
		w.WriteHeader(500)
		fmt.Fprintf(w, "<html><body><h3>error</h3></body></html>")
		return
	}
	signature := strings.TrimPrefix(xsignature, "sha256=")
	sigbytes, err := hex.DecodeString(signature)
	rdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("cannot process request")
		w.WriteHeader(500)
		fmt.Fprintf(w, "<html><body><h3>error</h3></body></html>")
		return
	}
	log.Debug(string(rdata[:]))
	mac := hmac.New(sha256.New, []byte(webhook.WebhookSecret))
	mac.Write(rdata)
	expected := mac.Sum(nil)
	if !hmac.Equal(sigbytes, expected) {
		log.Error("signature mismatch")
		w.WriteHeader(500)
		fmt.Fprintf(w, "<html><body><h3>error</h3></body></html>")
		return
	}
	hook := &Webhook{}
	err = json.Unmarshal(rdata, hook)
	if err != nil {
		log.Error("error parsing webhook data")
		w.WriteHeader(500)
		fmt.Fprintf(w, "<html><body><h3>error</h3></body></html>")
		return
	}
	t, err := time.Parse("2006-01-02T15:04:05-0700", hook.Data.Timestamp)
	ev := WebhookEvent{
	Tag: hook.Tag,
	Timestamp: t,
	GoogleId: hook.Data.GoogleId,
	}
	webhook.Events <- ev
	w.WriteHeader(200)
	fmt.Fprintf(w, "<html><body><h3>ok</h3></body></html>")
}