package controllers

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
)

type OAuth struct {
	ProviderConfigs map[string]*oauth2.Config
}

func (oa OAuth) Connect(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	provider = strings.ToLower(provider)

	config, ok := oa.ProviderConfigs[provider]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	state := csrf.Token(r)
	setCookie(w, "oauth_state", state)

	url := config.AuthCodeURL(
		state,
		// todo: read redirect_uri from env
		oauth2.SetAuthURLParam("redirect_uri", "http://localhost:3000/oauth/dropbox/callback"),
		oauth2.AccessTypeOffline,
	)

	http.Redirect(w, r, url, http.StatusFound)
}
