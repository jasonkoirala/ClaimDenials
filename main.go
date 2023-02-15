package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	acr "github.com/jasonkoirala/ClaimDenials/pkg/routes/aws_cognito_routes"
)

func main() {
	r := chi.NewRouter()
	acr.RegisterAWSCognitoRoutes(r)
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
