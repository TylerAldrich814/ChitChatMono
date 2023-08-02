package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/google/uuid"
)
type AWSGatewayResponse = events.APIGatewayProxyResponse
type AWSGatewayRequest = events.APIGatewayProxyRequest

const CCUserTable = "CCUserTable"

type User struct {
  Uid        uuid.UUID `json:"uid"`
  Username   string `json:"username"`
  HashedPswd string `json:"password"`
  Email      string `json:"email"`
}

type Tokens struct {
  AccessToken  string `json:"accessToken"`
  IDToken      string `json:"idToken"`
  RefreshToken string `json:"refreshToken"`
}

type Response struct {
  Message string `json:"message"`
  Body    []byte
}
func(resp *Response)RespondWithMessage(statusCode int)( AWSGatewayResponse,error ){
  jsonResponse, err := json.Marshal(resp)
  if err != nil {
    return AWSGatewayResponse{}, err
  }

  return AWSGatewayResponse{
    StatusCode: statusCode,
    Body: string(jsonResponse),
  }, nil
}

func authorize(
  ctx context.Context,
  request AWSGatewayRequest,
)( AWSGatewayResponse, error ){
  var user User
  var resp = Response{}

  err := json.Unmarshal([]byte(request.Body), &user)
  if err != nil {
    resp.Message = fmt.Sprintf("Error unmarshaling Request Body: %v\n", err)

    return resp.RespondWithMessage(400)
  }

  sess := session.Must(session.NewSession())
  cognitoSvc := cognitoidentityprovider.New(sess)

  userPoolId := "your-user-pool-id"
  clientId := "your-client-id"

	params := &cognitoidentityprovider.InitiateAuthInput{
    AuthFlow: aws.String("USER_PASSWORD_AUTH"),
    AuthParameters: map[string]*string{
      "USERNAME": aws.String(user.Username),
      "PASSWORD": aws.String(user.HashedPswd),
    },
    ClientId: aws.String(clientId),
  }
  authed, err := cognitoSvc.InitiateAuth(params)
  if err != nil {
    resp.Message = "User Failed Authorization."
    return resp.RespondWithMessage(401)
  }
  // If authorization is successfull,
  // 'authed' will contain the following tokens
  //    - Access Token
  //    - ID Token
  //    - Refresh Token
  tokens := Tokens{
    AccessToken:  *authed.AuthenticationResult.AccessToken,
    IDToken:      *authed.AuthenticationResult.IdToken,
    RefreshToken: *authed.AuthenticationResult.RefreshToken,
  }
  marshaledTokens, err := json.Marshal(tokens)
  if err != nil {
    resp.Message = fmt.Sprintf("Failed to Marshal Tokens: %v\n", err)
    return resp.RespondWithMessage(400)
  }

  resp = Response{Message: "User was successfully Authorized!"}
  resp.Body = marshaledTokens

  return resp.RespondWithMessage(200)
}

func main(){
  lambda.Start(authorize)
}
