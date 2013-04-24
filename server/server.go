package gogoogleoauth2

import (
  "net/http"
  "net/url"
  "fmt"

  "encoding/json"

  "io/ioutil"

  "appengine"
  "appengine/urlfetch"
  "html/template"

  "code.google.com/p/goauth2/oauth"
)

/**
 * json package only accesses the
 * the exported fields of struct types
 * (those that begin with an uppercase letter)
 */
type TokenResponse struct {
  Access_token string
  Token_type string
  Expires_in int
  Id_token string
}

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
  var tr TokenResponse
  err = json.Unmarshal(bb, &tr)

  url := REQUEST_API + tr.Access_token

  cl := urlfetch.Client(c)
  respo, _ := cl.Get(url)
  if respo.StatusCode != 200 {
    http.Error(w, "api error", http.StatusInternalServerError)
    return
  }
  defer respo.Body.Close()

  d, _ := ioutil.ReadAll(respo.Body)

  var ui UserInfo
  err = json.Unmarshal(d, &ui)

  renderTemplate(w, "success", &ui)
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

  renderTemplate(w, "success", &ui)
}

var templates = template.Must(template.ParseFiles(VIEW_PATH + "success.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, ui *UserInfo) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", ui)
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
