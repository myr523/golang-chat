package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"golang-chat/src/models/chat"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}


func main() {
	r := chat.NewRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// start chatroom
	go r.Run()
	// start web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("LintenAndServe: ", err)
	}

}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("chat", t.filename)))
	})
	t.templ.Execute(w, nil)
}
