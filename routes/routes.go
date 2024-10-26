package routes

import (
	"net/http"

	"github.com/noobpiyush/paytm-api/handlers"
)

func RegisteredRoutes() {
	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/signin", handlers.SigninHandler)
}
