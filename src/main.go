package main

import (
	"fmt"
	"golang-chat/src/models/auth"
	"golang-chat/src/models/chat"
	"golang-chat/src/models/util"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func main() {
	r := chat.NewRoom()
	err := auth.SetupProvider()
	if err != nil {
		fmt.Println(err.Error())
	}

	http.HandleFunc("/", util.Redirect)
	http.Handle("/chat", auth.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/room", r)
	http.HandleFunc("/auth/", auth.LoginHandler)
	// start chatroom
	go r.Run()
	// start web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("LintenAndServe: ", err)
	}

}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("html", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)

}
