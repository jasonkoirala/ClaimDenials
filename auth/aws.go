package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var cognitoClient *cognitoidentityprovider.CognitoIdentityProvider

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("REGION")),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	cognitoClient = cognitoidentityprovider.New(sess)
}

type Response struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("-------------------------------")

	type RequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
		Username: aws.String(body.Username),
		Password: aws.String(body.Password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(body.Email),
			},
			{
				Name:  aws.String("custom:role"), // new role attribute
				Value: aws.String(body.Role),
			},
		},
	}

	response, err := cognitoClient.SignUpWithContext(context.Background(), input)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(&Response{Message: "User registered successfully"})
	json.NewEncoder(w).Encode(response)
}

func VerifyEmailConfirmation(w http.ResponseWriter, r *http.Request) {
	type RequestBody struct {
		Username     string `json:"username"`
		Confirmation string `json:"confirmation"`
	}

	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := body.Username
	confirmationCode := body.Confirmation

	if username == "" || confirmationCode == "" {
		http.Error(w, "Missing required parameters: username, confirmationCode", http.StatusBadRequest)
		return
	}

	response, err := cognitoClient.ConfirmSignUpWithContext(context.TODO(), &cognitoidentityprovider.ConfirmSignUpInput{
		Username:         aws.String(username),
		ConfirmationCode: aws.String(confirmationCode),
		ClientId:         aws.String(os.Getenv("CLIENT_ID")),
	})

	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	type RequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeUserPasswordAuth),
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(body.Username),
			"PASSWORD": aws.String(body.Password),
		},
	}

	response, err := cognitoClient.InitiateAuthWithContext(context.Background(), input)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
