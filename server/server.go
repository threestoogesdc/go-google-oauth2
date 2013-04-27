package gogoogleoauth2

import (
  "net/http"
  "fmt"

  "html/template"

  "encoding/json"
  "io/ioutil"
  "errors"
  "strings"

  "code.google.com/p/goauth2/oauth"

  "appengine"
  "appengine/urlfetch"
)

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

type Scope struct {
  Desc string
  Url string
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

var oauthCfg *oauth.Config

var templates = template.Must(template.ParseFiles(VIEW_PATH + "success.html", VIEW_PATH + "index.html", VIEW_PATH + "slice.html"))

var scopes = []Scope{
  Scope {
    "Manage your calendars",
    "https://www.googleapis.com/auth/calendar",
  },
  Scope {
    "View your calendars",
    "https://www.googleapis.com/auth/calendar.readonly",
  },
  Scope {
    "View and manage the files and documents in your Google",
    "https://www.googleapis.com/auth/drive",
  },
  Scope {
    "View your Google Drive apps",
    "https://www.googleapis.com/auth/drive.apps.readonly",
  },
  Scope {
    "View and manage Google Drive files that you have opened or created with this app",
    "https://www.googleapis.com/auth/drive.file",
  },
  Scope {
    "View metadata for files and documents in your Google Drive",
    "https://www.googleapis.com/auth/drive.metadata.readonly",
  },
  Scope {
    "View the files and documents in your Google Drive",
    "https://www.googleapis.com/auth/drive.readonly",
  },
  Scope {
    "Modify your Google Apps Script scripts' behavior",
    "https://www.googleapis.com/auth/drive.scripts",
  },
  Scope {
    "Manage your tasks",
    "https://www.googleapis.com/auth/tasks",
  },
  Scope {
    "View your tasks",
    "https://www.googleapis.com/auth/tasks.readonly",
  },
  Scope {
    "View and manage your Google Maps Coordinate jobs",
    "https://www.googleapis.com/auth/coordinate",
  },
  Scope {
    "View your Google Coordinate jobs",
    "https://www.googleapis.com/auth/coordinate.readonly",
  },
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
  renderScopeTemplate(w, "index", scopes)
}

func formValues(r *http.Request, key string) []string {
    if r.Form == nil {
         r.ParseMultipartForm(200)
    }
    if vs := r.Form[key]; len(vs) > 0 {
         return vs
    }
    return nil
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
    http.Error(w, "error parsing", http.StatusInternalServerError)
    return
  }

  v := r.Form
  s := v["scopes[]"]

  as := strings.Join(s, " ")

  if as != "" {
    oauthCfg = &oauth.Config {
      ClientId: CLIENT_ID,
      ClientSecret: CLIENT_SECRET,
      AuthURL: BASE_SITE + AUTH_PATH,
      TokenURL: BASE_SITE + TOKEN_PATH,
      RedirectURL: REDIRECT_URI,
      Scope: as,
    }
    url := oauthCfg.AuthCodeURL("")

    http.Redirect(w, r, url, http.StatusFound)
  }

}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
  code := r.FormValue("code")
  error := r.FormValue("error")

  if error != "" {
    http.Error(w, "access denied", http.StatusInternalServerError)
    return
  }

  c := appengine.NewContext(r)
  transport := &oauth.Transport{Config: oauthCfg, Transport: &urlfetch.Transport{Context: c}}


  token, err := transport.Exchange(code)
  if err != nil {
    fmt.Fprintf(w, "error exchange %#v", err)
  }

  var ui UserInfo
  ui, err = userData(w, transport, token.AccessToken, &ui)

  renderUserTemplate(w, "success", &ui)
}

// obtain user data from constant REQUEST_API
// UserInfo is used to Unmarshal the json response
// returns UserInfor strct for display
func userData(w http.ResponseWriter, transport *oauth.Transport, token string, ui *UserInfo) (UserInfo, error) {
  url := REQUEST_API + token
  resp, err := transport.Client().Get(url)
  if err != nil {
    http.Error(w, "api error", http.StatusInternalServerError)
    return *ui, errors.New("api error")
  }
  defer resp.Body.Close()

  d, _ := ioutil.ReadAll(resp.Body)

  err = json.Unmarshal(d, &ui)

  return *ui, nil
}

func renderUserTemplate(w http.ResponseWriter, tmpl string, ui *UserInfo) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", ui)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func renderScopeTemplate(w http.ResponseWriter, tmpl string, slice []Scope) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", slice)
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
}
