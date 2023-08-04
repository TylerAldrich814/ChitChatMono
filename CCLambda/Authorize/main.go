package main

import (
	"context"
	"encoding/json"
	"fmt"
	"utils"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

const ChitChatUserId = utils.ChitChatUserID
type AWSGatewayRequest = utils.AWSGatewayRequest
type AWSGatewayResponse = utils.AWSGatewayResponse
type Response = utils.Response
type Tokens = utils.Tokens

func Authorize(
  ctx     context.Context,
  request AWSGatewayRequest,
)( AWSGatewayResponse, error ){
  var user utils.User
  var resp = utils.Response{}

  err := json.Unmarshal([]byte(request.Body), &user)
  if err != nil {
    resp.Message = fmt.Sprintf("Error unmarshaling Request Body: %v\n", err)

    return resp.RespondWith(400)
  }
  sess, err := session.NewSession(&aws.Config{
    Region: aws.String("us-east-2"),
  })
  if err != nil {
    resp.Message = "Failed to create AWS Session"
    resp.RespondWith(400)
  }

  cognitoSvc := cognitoidentityprovider.New(sess)

  // svc := secretsmanager.New(sess)
  // clientIdSec, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
  //   SecretId: aws.String(utils.ChitChatUserID),
  // })
  // if err != nil {
  //   resp.Message = fmt.Sprintf("Error occured while trying to receive ClientID\n%v\n", err)
  //   return resp.RespondWith(400)
  // }
  // clientId := *clientIdSec.SecretString

	params := &cognitoidentityprovider.InitiateAuthInput{
    AuthFlow: aws.String("USER_PASSWORD_AUTH"),
    ClientId: aws.String("NEED TO SET THIS UP TYLER"),
    AuthParameters: map[string]*string{
      "USERNAME": aws.String(user.Username),
      "PASSWORD": aws.String(user.HashedPswd),
    },
  }

  authed, err := cognitoSvc.InitiateAuth(params)
  if err != nil {
    resp.Message = "User Failed Authorization."
    return resp.RespondWith(401)
  }

  // If authorization was successfull, generate Tokens, store RefreshToken,
  // and return the users tokens
  tokens := Tokens{
    AccessToken:  *authed.AuthenticationResult.AccessToken,
    IDToken:      *authed.AuthenticationResult.IdToken,
  }

  utils.HandleRefreshToken(
    sess,
    user.Uid.String(),
    authed.AuthenticationResult.RefreshToken,
  )

  marshaledTokens, err := json.Marshal(tokens)
  if err != nil {
    resp.Message = fmt.Sprintf("Failed to Marshal Tokens: %v\n", err)
    return resp.RespondWith(400)
  }

  resp = Response{Message: "User was successfully Authorized!"}
  resp.Body = marshaledTokens

  return resp.RespondWith(200)
}


func main(){
  lambda.Start(Authorize)
}

