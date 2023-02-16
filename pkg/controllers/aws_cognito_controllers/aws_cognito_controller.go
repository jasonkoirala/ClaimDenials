package aws_cognito_controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	model "github.com/jasonkoirala/ClaimDenials/pkg/models/CognitoUserModels"
	"github.com/jasonkoirala/ClaimDenials/pkg/utils"
)

// User Registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser := &model.UserRequestBody{}

	if err := utils.ParseRequestBody(r, newUser); err != nil {
		fmt.Println("Inside parse request body error.")
		errorBody := &model.GlobalResponse{}
		errorBody.Error = true
		errorBody.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	resp, err := model.RegisterUserInCognitoUserPool(newUser)
	if err != nil {

		errorBody := &model.GlobalResponse{}

		errorBody.Error = true
		errorBody.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorBody)
		return
	}
	response := model.UserRegistrationResponse{
		CodeDeliveryDetails: struct {
			AttributeName  string `json:"attributeName"`
			DeliveryMedium string `json:"deliveryMedium"`
			Destination    string `json:"destination"`
		}{
			AttributeName:  *resp.CodeDeliveryDetails.AttributeName,
			DeliveryMedium: *resp.CodeDeliveryDetails.DeliveryMedium,
			Destination:    *resp.CodeDeliveryDetails.Destination,
		},
		UserConfirmed: *resp.UserConfirmed,
		UserSub:       *resp.UserSub,
		Message:       "Successfully registered the user",
		Error:         false,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// User verfication/ confirmation
func ConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newConfirmation := &model.UserConfirmationRequestBody{}
	if err := utils.ParseRequestBody(r, newConfirmation); err != nil {
		errorBody := &model.GlobalResponse{}
		errorBody.Error = true
		errorBody.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	if newConfirmation.Username == "" || newConfirmation.Confirmation == "" {
		errorBody := &model.GlobalResponse{}
		errorBody.Error = true
		errorBody.Message = "Username & Confirmation code are required fields."
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	_, err := model.ConfirmUserInCognitoUserPool(newConfirmation)
	if err != nil {
		errorBody := &model.GlobalResponse{}
		errorBody.Error = true
		errorBody.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorBody)
		return
	}
	response := model.GlobalResponse{}
	response.Error = false
	response.Message = "Successfully Verified."
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// User Login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loginRequest := &model.UserLoginRequestBody{}
	if err := utils.ParseRequestBody(r, loginRequest); err != nil {
		errorBody := &model.GlobalResponse{}
		errorBody.Error = true
		errorBody.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	if loginRequest.Username == "" || loginRequest.Password == "" {
		errorBody := &model.GlobalResponse{}
		errorBody.Error = true
		errorBody.Message = "Username & Password are required fields."
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	resp, err := model.UserAuthenticationInCognitoUserPool(loginRequest)
	if err != nil {
		errorBody := &model.GlobalResponse{}
		errorBody.Error = true
		errorBody.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorBody)
		return
	}

	role, _ := model.GetUserRole(*resp.AuthenticationResult.AccessToken)
	response := model.CognitoLoginResponse{
		AuthenticationResult: struct {
			AccessToken  string `json:"AccessToken"`
			ExpiresIn    int    `json:"ExpiresIn"`
			TokenType    string `json:"TokenType"`
			RefreshToken string `json:"RefreshToken"`
			IdToken      string `json:"IdToken"`
		}{
			AccessToken:  *resp.AuthenticationResult.AccessToken,
			TokenType:    *resp.AuthenticationResult.TokenType,
			RefreshToken: *resp.AuthenticationResult.RefreshToken,
			IdToken:      *resp.AuthenticationResult.IdToken,
		},
		ChallengeName: *&resp.ChallengeName,
		Session:       *&resp.Session,
		Role:          role,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
