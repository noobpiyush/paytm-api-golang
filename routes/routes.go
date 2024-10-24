package routes

import (
	"net/http"

	"github.com/noobpiyush/paytm-api/handlers"
)

func RegisteredRoutes() {
	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/signin", handlers.SigninHandler)
	http.HandleFunc("/", handleHealth)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	handlers.GetIP(r)
	w.Write([]byte("hiii there "))
}
