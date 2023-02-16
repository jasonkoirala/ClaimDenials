package CognitoUserModels

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type GlobalResponse struct {
	Message string `json:"message"`
	Error   bool   `json:"error"`
}

type UserRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserConfirmationRequestBody struct {
	Username     string `json:"username"`
	Confirmation string `json:"confirmation"`
}

type UserRegistrationResponse struct {
	CodeDeliveryDetails struct {
		AttributeName  string `json:"attributeName"`
		DeliveryMedium string `json:"deliveryMedium"`
		Destination    string `json:"destination"`
	} `json:"codeDeliveryDetails"`
	UserConfirmed bool   `json:"userConfirmed"`
	UserSub       string `json:"userSub"`
	Message       string `json:"message"`
	Error         bool   `json:"error"`
}

type UserLoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CognitoLoginResponse struct {
	AuthenticationResult struct {
		AccessToken  string `json:"AccessToken"`
		ExpiresIn    int    `json:"ExpiresIn"`
		TokenType    string `json:"TokenType"`
		RefreshToken string `json:"RefreshToken"`
		IdToken      string `json:"IdToken"`
	} `json:"AuthenticationResult"`
	ChallengeName *string `json:"ChallengeName"`
	Session       *string `json:"Session"`
	Role          string  `json:"role"`
}

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

func RegisterUserInCognitoUserPool(newUser *UserRequestBody) (*cognitoidentityprovider.SignUpOutput, error) {
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
		Username: aws.String(newUser.Username),
		Password: aws.String(newUser.Password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(newUser.Email),
			},
			{
				Name:  aws.String("custom:role"),
				Value: aws.String(newUser.Role),
			},
		},
	}

	response, err := cognitoClient.SignUpWithContext(context.Background(), input)
	return response, err
}

func ConfirmUserInCognitoUserPool(confirmation *UserConfirmationRequestBody) (*cognitoidentityprovider.ConfirmSignUpOutput, error) {
	response, err := cognitoClient.ConfirmSignUpWithContext(context.TODO(), &cognitoidentityprovider.ConfirmSignUpInput{
		Username:         aws.String(confirmation.Username),
		ConfirmationCode: aws.String(confirmation.Confirmation),
		ClientId:         aws.String(os.Getenv("CLIENT_ID")),
	})
	return response, err
}

func UserAuthenticationInCognitoUserPool(loginRequest *UserLoginRequestBody) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	input := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeUserPasswordAuth),
		ClientId: aws.String(os.Getenv("CLIENT_ID")),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(loginRequest.Username),
			"PASSWORD": aws.String(loginRequest.Password),
			"SCOPES":   aws.String("custom:role"),
		},
	}

	response, err := cognitoClient.InitiateAuthWithContext(context.Background(), input)
	return response, err
}

func GetUserRole(accessToken string) (string, string) {

	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(accessToken),
	}

	result, err := cognitoClient.GetUserWithContext(context.Background(), input)

	if err != nil {
		return "NA", "NA"
	}

	var role string
	for _, attr := range result.UserAttributes {
		if *attr.Name == "custom:role" {
			role = *attr.Value
			break
		}
	}

	return role, "NA"
}
