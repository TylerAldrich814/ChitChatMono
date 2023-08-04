package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/google/uuid"
)

type AWSGatewayResponse = events.APIGatewayProxyResponse
type AWSGatewayRequest = events.APIGatewayProxyRequest

const CCUserTable     string = "CCUsers"
const SECRETTOKENAUTH string = "SECRET_TOKEN_AUTH"
const ChitChatUserID  string = "ChitChatClientId"

type User struct {
  Uid        uuid.UUID `json:"uid"`
  Username   string `json:"username"`
  HashedPswd string `json:"password"`
  Email      string `json:"email"`
}

type Response struct {
  Message string `json:"message"`
  Body    []byte
}

type Tokens struct {
  AccessToken  string `json:"accessToken"`
  IDToken      string `json:"idToken"`
}

func(resp *Response)RespondWith(statusCode int)( AWSGatewayResponse,error ){
  jsonResponse, err := json.Marshal(resp)
  if err != nil {
    return AWSGatewayResponse{
      StatusCode: 400,
      Body: "Failed to Marshal Response Error...",
    }, err
  }

  return AWSGatewayResponse{
    StatusCode: statusCode,
    Body: string(jsonResponse),
  }, nil
}

// This will check if a User already has a RefreshToken stored in our
// secretsmanager or not.
func HandleRefreshToken(
  sess         *session.Session,
  secretId     string,
  secretString *string,
) error {
  // Token-Helper Funcitons
  svc := secretsmanager.New(sess)

  // _, err := svc.DescribeSecret(&secretsmanager.DescribeSecretInput{
  //   SecretId: aws.String(secretId),
  // })
  // if err != nil {
  //   if aerr, ok := err.(awserr.Error); ok {
  //     switch aerr.Code() {
  //     case secretsmanager.ErrCodeResourceNotFoundException:
  //       return NewRefreshToken(svc, secretId, secretString)
  //     default:
  //       return errors.New(fmt.Sprintf("Error occured while checking if Secret exists: %v", err))
  //     }
  //   }
  // }
  return UpdateRefreshToken(svc, secretId, secretString)
}

func NewRefreshToken(
  svc          *secretsmanager.SecretsManager,
  secretid     string,
  secretstring *string,
) error {
  // input := &secretsmanager.CreateSecretInput{
  //   Name:         aws.String(secretid),
  //   SecretString: aws.String(*secretstring),
  // }
  //
  // _, err := svc.CreateSecret(input)
  return nil
}

func UpdateRefreshToken(
  svc          *secretsmanager.SecretsManager,
  secretid     string,
  secretstring *string,
) error {
  // input := &secretsmanager.UpdateSecretInput{
  //   SecretId: aws.String(secretid),
  //   SecretString: aws.String(*secretstring),
  // }
  //
  // _, err := svc.UpdateSecret(input)
  // return err

  return nil
}
