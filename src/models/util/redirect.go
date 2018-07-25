package util

import (
	"net/http"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8080/chat", http.StatusTemporaryRedirect)
}
