package gogoogleoauth2

import (
  "net/http"
  "fmt"

  "code.google.com/p/goauth2/oauth"
)

const (
  BASE_SITE = "https://accounts.google.com/o"
  AUTH_PATH = "/oauth2/auth"
  TOKEN_PATH = "/oauth2/token"
)

var oauthCfg = &oauth.Config {
  ClientId: "670315590273-04lcsdb09rom5d3uejvnet15fti0affi",
  ClientSecret: "DwdAYjN92XJL3HpnGkfFd7JE",
  AuthURL: BASE_SITE + AUTH_PATH,
  TokenURL: BASE_SITE + TOKEN_PATH,
  RedirectURL: "http://ts-go-oauth2.appspot.com/auth/callback",
  Scope: "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "root");
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
  url := oauthCfg.AuthCodeURL("")

  http.Redirect(w, r, url, http.StatusFound)
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
  // TODO implement token
  fmt.Fprintln(w, "callback")
}

func init() {
  http.HandleFunc("/", handleRoot)
  http.HandleFunc("/auth", handleAuth);
  http.HandleFunc("/auth/callback", handleAuthCallback)
}

