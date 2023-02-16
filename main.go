package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	acr "github.com/jasonkoirala/ClaimDenials/pkg/routes/aws_cognito_routes"
	"github.com/rs/cors"
)

func main() {
	r := chi.NewRouter()
	acr.RegisterAWSCognitoRoutes(r)
	port := os.Getenv("PORT")
	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))
}
