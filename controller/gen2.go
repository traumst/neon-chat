package controller

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"go.chat/utils"
)

func Gen2(w http.ResponseWriter, r *http.Request) {
	log.Printf("--%s-> Gen2", utils.GetReqId(r))
	user, err := utils.GetCurrentUser(r)
	if err != nil {
		log.Printf("--%s-> Gen2 WARN user, %s\n", utils.GetReqId(r), err)
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}

	var buf bytes.Buffer
	tmpl := template.Must(template.ParseFiles("html/generated.html"))
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		log.Printf("--%s-> Gen2 TRACE, user[%s] gets content\n", utils.GetReqId(r), user)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Error parsing template"))
	} else {
		log.Printf("--%s-> Gen2 TRACE, user[%s] gets content\n", utils.GetReqId(r), user)
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}
