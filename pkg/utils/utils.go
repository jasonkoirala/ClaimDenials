package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func ParseRequestBody(r *http.Request, X interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, X); err != nil {
		return err
	}
	return nil
}

func PrepareSecrets() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		// handle error
	}
	svc := secretsmanager.New(sess)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String("ClaimDenialsUserPoolID"),
	}

	output, err := svc.GetSecretValue(input)
	if err != nil {
		// handle error
	}
	if output != nil {
		if output.SecretString != nil {
			fmt.Println(*output.SecretString)
		} else if output.SecretBinary != nil {
			fmt.Println(string(output.SecretBinary))
		} else {
			fmt.Println("Secret value is empty.")
		}
	} else {
		fmt.Println("Output is nil.")
	}
}
