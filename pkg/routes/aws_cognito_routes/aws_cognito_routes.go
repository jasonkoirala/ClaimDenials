package aws_cognito_routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jasonkoirala/ClaimDenials/auth"
)

var RegisterAWSCognitoRoutes = func(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Claim Denials...!!!"))
	})
	r.Post("/register", auth.RegisterHandler)
	r.Post("/verify", auth.VerifyEmailConfirmation)
	r.Post("/login", auth.LoginHandler)
}
