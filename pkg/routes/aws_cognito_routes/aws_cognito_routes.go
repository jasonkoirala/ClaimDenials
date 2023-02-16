package aws_cognito_routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	awscontroller "github.com/jasonkoirala/ClaimDenials/pkg/controllers/aws_cognito_controllers"
)

var RegisterAWSCognitoRoutes = func(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Claim Denials...!!!"))
	})
	r.Post("/register", awscontroller.RegisterHandler)
	r.Post("/verify", awscontroller.ConfirmationHandler)
	r.Post("/login", awscontroller.LoginHandler)
}
