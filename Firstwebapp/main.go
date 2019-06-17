
package main

import (
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"

	//"github.com/stretchr/gomniauth/providers/facebook"
	//"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
	"time"
	"trace"
)
type message struct {
	Name string
	Message string
	When time.Time
}
type templateHandler struct {
	once sync.Once
	filename string
	templ    *template.Template
}
// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data:=map[string]interface{}{
		"Host":r.Host,
	}
	if authCookies, err:=r.Cookie("auth");err==nil{

		data["UserData"]=objx.MustFromBase64(authCookies.Value)

	}
	t.templ.Execute(w,data)
}

func main() {
	var adr string=os.Args[1]
	gomniauth.SetSecurityKey("98dfbg7iu2nb4uywevihjw4tuiyub34noilk")
	gomniauth.WithProviders(
		google.New("530631165756-t2r83rl435hnpg7o2ruh2lrfbs3tsloh.apps.googleusercontent.com", "ThF8q1JBi60XoZAxg3TSDtE6", "http://localhost:8080/auth/callback/google"),

	)

	// root
	r := NewRoom()
	r.tracer=trace.New(os.Stdout)
	//http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/",loginHandler)
	http.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	if err := http.ListenAndServe(adr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}