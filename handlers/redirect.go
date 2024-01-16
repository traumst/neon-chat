package handlers

import "net/http"

func SendBack(w http.ResponseWriter, r *http.Request) {
	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusBadRequest)
}
