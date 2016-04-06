package goviewnew

import (
      "github.com/gorilla/mux"
      "github.com/gorilla/securecookie"
      "html/template"
      "io/ioutil"
      "net/http"
      "log"
      "dbconfig"
)

var (
  rootPage []byte
  signupPage []byte
)

//Fetch all templates
var templates, templates_err = template.ParseGlob("templates/*")

type Page struct {
    Title string
    Body  []byte
}
  
func init() {
  dbconfig.InitDB("root:mosalis119988@tcp(127.0.0.1:3306)/Users")
  r := mux.NewRouter()
  //Manage templates
    r.HandleFunc("/", root)
    r.HandleFunc("/signup", signup)
    r.HandleFunc("/login", loginHandler).Methods("POST")
    r.HandleFunc("/signupauth", signupauthHandler).Methods("POST")
    r.HandleFunc("/logout", logoutHandler).Methods("POST")
    r.HandleFunc("/dash", dash)
    http.Handle("/", r)

  signup_content, err := ioutil.ReadFile("templates/signup.html")
  
  if err != nil {
    panic(err)
  }
    
  if err != nil {
    panic(err)
  }

  if err != nil {
    panic(err)
  }
  
  signupPage = signup_content
}


func root(w http.ResponseWriter, r *http.Request) { 
  login_content, err := ioutil.ReadFile("templates/login.html")
  if err != nil {
    panic(err)
  }
    
  if err != nil {
    panic(err)
  }

  if err != nil {
    panic(err)
  }
  
  loginPage := login_content
   w.Write(loginPage) 
}

func signup(w http.ResponseWriter, r *http.Request) { 
  signup_content, err := ioutil.ReadFile("templates/signup.html")
  if err != nil {
    panic(err)
  }
    
  if err != nil {
    panic(err)
  }

  if err != nil {
    panic(err)
  }
  
  signupPage := signup_content
   w.Write(signupPage) 
}

func dash(w http.ResponseWriter, request *http.Request) { 
  root_data := make(map[string]string)
  root_data["username"]   = getUserName(request)
  root_data["body"] = ""

  templates.ExecuteTemplate(w, "dashboard.html", root_data)
}


var cookieHandler = securecookie.New(
  securecookie.GenerateRandomKey(64),
  securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
  if cookie, err := request.Cookie("session"); err == nil {
    cookieValue := make(map[string]string)
    if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
      userName = cookieValue["name"]
    }
  }
  return userName
}

func setSession(userName string, response http.ResponseWriter) {
  value := map[string]string{
    "name": userName,
  }
  if encoded, err := cookieHandler.Encode("session", value); err == nil {
    cookie := &http.Cookie{
      Name:  "session",
      Value: encoded,
      Path:  "/",
    }
    http.SetCookie(response, cookie)
  }
}

func clearSession(response http.ResponseWriter) {
  cookie := &http.Cookie{
    Name:   "session",
    Value:  "",
    Path:   "/",
    MaxAge: -1,
  }
  http.SetCookie(response, cookie)
}

// login handler

func loginHandler(response http.ResponseWriter, request *http.Request) {
  usermail := request.FormValue("email")
  pass := request.FormValue("password")
  redirectTarget := "/"
  var uid, username, email, password, created string ;
  if usermail != "" && pass != "" {
    err := dbconfig.DB.QueryRow("SELECT * FROM Users.userinfo where email=?", usermail).Scan(&uid, &username, &email, &password, &created)
      if err != nil {
        log.Panic("User cant be found")
      } else if password == pass {
         // .. check credentials ..
        setSession(username, response)
        redirectTarget = "/dash"
      }
  }
  http.Redirect(response, request, redirectTarget, 302)
}

// Sign up handler
func signupauthHandler(response http.ResponseWriter, request *http.Request) {
  name := request.FormValue("name")
  email := request.FormValue("email")
  pass := request.FormValue("password")
  redirectTarget := "/"
  stmt, err := dbconfig.DB.Prepare("INSERT userinfo SET username=?,email=?,password=?,created=?")
  if err != nil {
    log.Panic(err)
  }
  stmt.Exec(name, email, pass, "1:00 Am")
  if name != "" && pass != "" {
    setSession(name, response)
    redirectTarget = "/dash"
  }
  http.Redirect(response, request, redirectTarget, 302)
}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
  clearSession(response)
  http.Redirect(response, request, "/", 302)
}
