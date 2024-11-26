package middleware

import (
	"log"
	"neon-chat/src/utils"
	h "neon-chat/src/utils/http"
	"net/http"
)

func MaintenanceMiddleware() Middleware {
	return Middleware{
		Name: "Maintenance",
		Func: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := utils.MaintenanceManager.IncrUserCount(); err != nil {
					log.Println("Server is under maintenance", err)
					h.SetRetryAfterHeader(&w, 10)
					http.Error(w, "Under Maintenance", http.StatusServiceUnavailable)
					return
				}
				defer utils.MaintenanceManager.DecrUserCount()
				next.ServeHTTP(w, r)
			})
		}}
}
