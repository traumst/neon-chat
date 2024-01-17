package utils

import "net/http"

func SendBack(w http.ResponseWriter, r *http.Request, status int) {
	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, status)
}
