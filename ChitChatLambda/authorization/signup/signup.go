package signup

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/TylerAldrich814/chitchat/authorization/utils"
)

type AWSGatewayResponse = utils.AWSGatewayResponse
type AWSGatewayRequest = utils.AWSGatewayRequest
type User = utils.User
type Response = utils.Response

func Signup(
  ctx     context.Context,
  request AWSGatewayRequest,
  sess    *session.Session,
)( AWSGatewayResponse, error ){
  cognitoSvc := cognitoidentityprovider.New(sess)
  dynamoSvc  := dynamodb.New(sess)

  var resp = Response{}
  var user User
  err := json.Unmarshal([]byte(request.Body), &user)
  if err != nil {
    resp.Message = "Failed to unmarshal Request Body."
    resp.RespondWith(400)
  }

  // Useing AWS Cognito To Sign up new user. Failed if user exists already,
  _, err = cognitoSvc.SignUp(&cognitoidentityprovider.SignUpInput{
    ClientId: aws.String(utils.ChitChatUserID),
    Username: aws.String(user.Username),
    Password: aws.String(user.HashedPswd),

    UserAttributes: []*cognitoidentityprovider.AttributeType{
      {
        Name:  aws.String("email"),
        Value: aws.String(user.Email),
      },
    },
  })
  if err != nil {
    awsErr, ok := err.(awserr.Error)
    if ok &&( awsErr.Code() == cognitoidentityprovider.ErrCodeUsernameExistsException ){
      resp.Message = fmt.Sprintf("The Username '%v' already exists!", user.Username)
    } else {
      resp.Message = fmt.Sprintf("Could not sign user up with Cognito: %v", err)
    }
    resp.RespondWith(400)
  }

  // After Authorization && Signup is complete. We'll store the rest of the user's
  // information in our Dynamodb instance.
  av, err := dynamodbattribute.MarshalMap(user)
  if err != nil {
    resp.Message = "Error occurred while Marshaling User"
    return resp.RespondWith(400)
  }

  _, err = dynamoSvc.PutItem(&dynamodb.PutItemInput{
    TableName: aws.String(utils.CCUserTable),
    Item: av,
  })
  if err != nil {
    resp.Message = "Failed to save User Data in our DynamoDB User's Table.."
    return resp.RespondWith(400)
  }

  resp.Message = "Successfully Signed up user! User is now saved in DynamoDB"
  return resp.RespondWith(200)
}
