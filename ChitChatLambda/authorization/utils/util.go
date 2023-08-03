package utils

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type AWSGatewayResponse = events.APIGatewayProxyResponse
type AWSGatewayRequest = events.APIGatewayProxyRequest

const CCUserTable     string = "CCUserTable"
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
