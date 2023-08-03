package token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/TylerAldrich814/chitchat/authorization/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)
type AWSGatewayResponse = utils.AWSGatewayResponse
type AWSGatewayRequest = utils.AWSGatewayRequest
type Response = utils.Response
type User = utils.User

type Tokens struct {
  AccessToken  string `json:"accessToken"`
  IDToken      string `json:"idToken"`
}

func StoreSecretToken(
  userId       string,
  refreshToken string,
  sess         *session.Session,
) error {
  svc := secretsmanager.New(sess)

  input := &secretsmanager.CreateSecretInput{
    Name:         aws.String(userId),
    SecretString: aws.String(refreshToken),
  }

  _, err := svc.CreateSecret(input)
  if err != nil {
    return errors.New(fmt.Sprintf("Error while trying to store RefreshToken\n%v\n", err.Error()))
  }
  return nil
}

// Used for when a Users Refresh Token expires.
// This Lambda will Obtain the expired Refreshtoken from Cognito.
// Then, we'll refresh all of the users Tokens, store the new RefreshToken
// Then return a Token{ AccessToken, IDToken }
// func refreshToken(userId string, clientId string)( *Tokens, error ){
func RefreshToken(
  ctx     context.Context,
  request AWSGatewayRequest,
  sess    *session.Session,
)( AWSGatewayResponse, error ){
  var resp = Response{}
  svc := secretsmanager.New(sess)

  var user User

  err := json.Unmarshal([]byte(request.Body), &user)
  if err != nil {
    resp.Message = fmt.Sprintf("Error unmarshaling Request Body: %v\n", err)

    return resp.RespondWith(400)
  }


  secret, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
    SecretId: aws.String(user.Uid.String()),
  })
  if err != nil {
    resp.Message = fmt.Sprintf("Error while trying retrieve stored RefreshToken\n%v\n", err.Error())
    return resp.RespondWith(400)
  }

  refreshToken := *secret.SecretString

  cognitoSvc := cognitoidentityprovider.New(sess)
  input := &cognitoidentityprovider.InitiateAuthInput{
    AuthFlow: aws.String(utils.SECRETTOKENAUTH),
    ClientId: aws.String(user.Uid.String()),
    AuthParameters: map[string]*string{
      "REFRESH_TOKEN": aws.String(refreshToken),
    },
  }
  result, err := cognitoSvc.InitiateAuth(input)
  if err != nil {
    resp.Message = fmt.Sprintf("Error while trying refresh Tokens\n%v\n", err.Error())
    return resp.RespondWith(404)
  }
  toke := Tokens{
    AccessToken: *result.AuthenticationResult.AccessToken,
    IDToken: *result.AuthenticationResult.IdToken,
  }
  tokeJson, err := json.Marshal(toke)
  if err != nil {
    resp.Message = "An error occured while marshaling your Tokens"
    resp.RespondWith(400)
  }
  resp.Body = tokeJson

  return resp.RespondWith(200)
}
