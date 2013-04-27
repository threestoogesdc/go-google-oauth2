package gogoogleoauth2

import (
  "net/http"
  "fmt"

  "html/template"

  "encoding/json"
  "io/ioutil"

  "code.google.com/p/goauth2/oauth"

  "appengine"
  "appengine/urlfetch"
)


type TokenResponse struct {
  AccessToken string
  RefreshToken string
  Expiry int64
}

// json package only accesses the
// the exported fields of struct types
// (those that begin with an uppercase letter)
type UserInfo struct {
  Id string
  Email string
  Verified_email bool
  Name string
  Given_name string
  Family_name string
  Link string
  Picture string
  Gender string
  Birthday string
}


const (
  VIEW_PATH = "app/views/"

  BASE_SITE = "https://accounts.google.com/o"
  AUTH_PATH = "/oauth2/auth"
  TOKEN_PATH = "/oauth2/token"
  PROFILE_URL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

  CLIENT_ID = "670315590273-04lcsdb09rom5d3uejvnet15fti0affi"
  CLIENT_SECRET = "DwdAYjN92XJL3HpnGkfFd7JE"
  REDIRECT_URI = "http://localhost:8080/auth/callback"

  REQUEST_API = "https://www.googleapis.com/oauth2/v1/userinfo?access_token="
)

var oauthCfg = &oauth.Config {
  ClientId: CLIENT_ID,
  ClientSecret: CLIENT_SECRET,
  AuthURL: BASE_SITE + AUTH_PATH,
  TokenURL: BASE_SITE + TOKEN_PATH,
  RedirectURL: REDIRECT_URI,
  Scope: "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
}

//var templates = template.Must(template.ParseFiles())

func handleRoot(w http.ResponseWriter, r *http.Request) {
  renderTemplate(w, "index")
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
  transport := &oauth.Transport{Config: oauthCfg, Transport: &urlfetch.Transport{Context: c}}


  token, err := transport.Exchange(code)
  if err != nil {
    fmt.Fprintf(w, "error exchange %#v", err)
  }

  url := REQUEST_API + token.AccessToken
  resp, err := transport.Client().Get(url)
  if err != nil {
    http.Error(w, "api error", http.StatusInternalServerError)
    return
  }
  defer resp.Body.Close()

  d, _ := ioutil.ReadAll(resp.Body)

  var ui UserInfo
  err = json.Unmarshal(d, &ui)

  renderUserTemplate(w, "success", &ui)
}

func handleSuccess(w http.ResponseWriter, r *http.Request) {
  var ui = UserInfo{
    "id",
    "test",
    true,
    "name",
    "given",
    "family",
    "link",
    "picture",
    "gender",
    "birthday",
  }

  renderUserTemplate(w, "success", &ui)
}

var templates = template.Must(template.ParseFiles(VIEW_PATH + "success.html", VIEW_PATH + "index.html"))

func renderUserTemplate(w http.ResponseWriter, tmpl string, ui *UserInfo) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", ui)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", nil)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

}

func init() {
  http.HandleFunc("/", handleRoot)
  http.HandleFunc("/auth", handleAuth);
  http.HandleFunc("/auth/callback", handleAuthCallback)
  http.HandleFunc("/success", handleSuccess)
}
