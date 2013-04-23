package gogoogleoauth2

import (
  "net/http"
  "net/url"
  "fmt"

  "io/ioutil"

  "appengine"
  "appengine/urlfetch"

  "code.google.com/p/goauth2/oauth"
)

const (
  BASE_SITE = "https://accounts.google.com/o"
  AUTH_PATH = "/oauth2/auth"
  TOKEN_PATH = "/oauth2/token"
  PROFILE_URL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

  CLIENT_ID = "670315590273-04lcsdb09rom5d3uejvnet15fti0affi"
  CLIENT_SECRET = "DwdAYjN92XJL3HpnGkfFd7JE"
  REDIRECT_URI = "http://ts-go-oauth2.appspot.com/auth/callback"
)

var oauthCfg = &oauth.Config {
  ClientId: CLIENT_ID,
  ClientSecret: CLIENT_SECRET,
  AuthURL: BASE_SITE + AUTH_PATH,
  TokenURL: BASE_SITE + TOKEN_PATH,
  RedirectURL: REDIRECT_URI,
  Scope: "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
}

var templates = template.Must(template.ParseFiles())

func handleRoot(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "root");
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
  url := oauthCfg.AuthCodeURL("")

  http.Redirect(w, r, url, http.StatusFound)
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
  /**
   * @TODO
   * get token from response
   * handle denied access
   */
  code := r.FormValue("code")
  error := r.FormValue("error")

  //access denied
  if error != "" {
    fmt.Fprintln(w, "access denied")
    return
  }

  c := appengine.NewContext(r)
  client := urlfetch.Client(c)

  data := url.Values{}
  data.Add("code", code)
  data.Add("client_id", CLIENT_ID)
  data.Add("client_secret", CLIENT_SECRET)
  data.Add("redirect_uri", REDIRECT_URI)
  data.Add("grant_type", "authorization_code")

  resp, _ := client.PostForm(BASE_SITE + TOKEN_PATH, data)

  if resp.StatusCode != 200 {
    http.Error(w, "error", http.StatusInternalServerError)
    return
  }
  defer resp.Body.Close()

  bb, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    http.Error(w, "error", http.StatusInternalServerError)
    return
  }
  bs := string(bb)

  fmt.Fprintln(w, bs)
}

func init() {
  http.HandleFunc("/", handleRoot)
  http.HandleFunc("/auth", handleAuth);
  http.HandleFunc("/auth/callback", handleAuthCallback)
  http.HandleFunc("/auth/success", handleAuthSuccess)
  http.HandleFunc("/auth/failure", handleAuthFailure)
}

