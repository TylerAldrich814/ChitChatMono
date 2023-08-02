package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"

	// "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	// "golang.org/x/crypto/bcrypt"
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

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-2"))

// queryUser :: For querying a User object from DynamoDB.
func queryUser(
  request AWSGatewayRequest,
)( *User, error ){
  var user User
  err := json.Unmarshal([]byte(request.Body), &user)
  if err != nil {
    return  nil, err
  }
  input := &dynamodb.QueryInput{
    TableName: aws.String(CCUserTable),
    KeyConditionExpression: aws.String("username = :e AND email = :e AND password = :p"),
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
      ":u": {S: aws.String(user.Username)},
      ":e": {S: aws.String(user.Email)},
      ":p": {S: aws.String(user.HashedPswd)},
    },
  }
  result, err := db.Query(input)
  if err != nil {
    return nil, err
  }
  if *result.Count == 0 {
    return nil, nil
  }
  item := result.Items[0]
  user.Email = aws.StringValue(item["email"].S)
  user.HashedPswd = aws.StringValue(item["username"].S)

  return &user, nil
}

func signup(ctx context.Context, request AWSGatewayRequest)( AWSGatewayResponse, error ){
  user, err := queryUser(request)
  if err != nil {
    resp := Response{Message: fmt.Sprintf("Error while Querying for User: %v", err)}
    return resp.RespondWithMessage(500)
  }
  if user == nil {
    resp := Response{
      Message: "Username and/or Email already exist",
    }
    return resp.RespondWithMessage(401)
  }

  err = json.Unmarshal([]byte(request.Body), user)
  if err != nil {
    resp := Response{Message: "Failed to extract Data from Response Body"}
    return resp.RespondWithMessage(502)
  }

  av, err := dynamodbattribute.MarshalMap(user)
  if err != nil {
    resp := Response{Message: "Failed to extract Data from Response Body"}
    return resp.RespondWithMessage(502)
  }
  input := &dynamodb.PutItemInput{
    Item:      av,
    TableName: aws.String(CCUserTable),
  }
  _, err = db.PutItem(input)
  if err != nil {
    resp := Response{Message: "Failed to place user in Database"}
    return resp.RespondWithMessage(502)
  }
  resp := Response{Message: "User has Successsfully Signed up!"}
  return resp.RespondWithMessage(200)
}


func signin(ctx context.Context, request AWSGatewayRequest)( AWSGatewayResponse, error) {
  user, err := queryUser(request)
  if err != nil {
    resp := Response{Message: fmt.Sprintf("Error while Querying for User: %v", err)}
    return resp.RespondWithMessage(500)
  }
  if user == nil {
    resp := Response{
      Message: "User does not exist!",
    }
    return resp.RespondWithMessage(401)
  }
  err = json.Unmarshal([]byte(request.Body), &user)
  if err != nil {
    resp := Response{Message: "Failed to extract Data from Response Body"}
    return resp.RespondWithMessage(400)
  }

  av, err := dynamodbattribute.MarshalMap(user)
  if err != nil {
    resp := Response{Message: "Failed to extract Data from Response Body"}
    return resp.RespondWithMessage(400)
  }
  input := &dynamodb.PutItemInput{
    Item:      av,
    TableName: aws.String(CCUserTable),
  }
  _, err = db.PutItem(input)
  if err != nil {
    resp := Response{Message: "Failed to place user in Database"}
    return resp.RespondWithMessage(502)
  }
  resp := Response{Message: "User has Successsfully Signed up!"}
  return resp.RespondWithMessage(200)

}

func joinChatroom(ctx context.Context, request AWSGatewayRequest)( AWSGatewayResponse, error ){

  return events.APIGatewayProxyResponse{}, nil
}

func sendMessage(ctx context.Context, request AWSGatewayRequest)( AWSGatewayResponse, error ){

  return events.APIGatewayProxyResponse{}, nil
}


func receiveMessage(ctx context.Context, request AWSGatewayRequest)( AWSGatewayResponse, error ){

  return events.APIGatewayProxyResponse{}, nil
}

func handler(ctx context.Context, request AWSGatewayRequest)( AWSGatewayResponse, error ){
  switch request.HTTPMethod {
  case "POST":
    switch request.Resource {
    case "/signup":
      return signup(ctx, request)
    case "/signin":
      return signin(ctx, request)
    case "/joiSn":
      return joinChatroom(ctx, request)
    case "/send":
      return sendMessage(ctx, request)
    }
  case "GET":
    switch request.Resource {
    case "/receive":
      return receiveMessage(ctx, request)
    }
  }

  return events.APIGatewayProxyResponse{}, nil
}

// func handler(request AWSGatewayRequest) (AWSGatewayResponse, error) {
// 	var greeting string
// 	sourceIP := request.RequestContext.Identity.SourceIP
//
// 	if sourceIP == "" {
// 		greeting = "Hello, world!\n"
// 	} else {
// 		greeting = fmt.Sprintf("Hello MOTHERFUCKER!, %s!\n", sourceIP)
// 	}
//
// 	return events.APIGatewayProxyResponse{
// 		Body:       greeting,
// 		StatusCode: 200,
// 	}, nil
// }

func main() {
	lambda.Start(handler)
}
