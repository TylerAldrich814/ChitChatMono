// package main
package chitchat

import (
	"github.com/TylerAldrich814/chitchat/authorization/auth"
	"github.com/TylerAldrich814/chitchat/authorization/signup"
	"github.com/TylerAldrich814/chitchat/authorization/tokens"
	"github.com/TylerAldrich814/chitchat/authorization/utils"
	"github.com/aws/aws-lambda-go/lambda"
)
type AWSGatewayRequest = utils.AWSGatewayRequest
type AWSGatewayResponse = utils.AWSGatewayResponse

const AWS_REGION = "us-east-2"

func handleAuth(){
  lambda.Start(auth.Authorize)
}

func handleSignup(){
  lambda.Start(signup.Signup)
}

func handleTokenRefresh(){
  lambda.Start(token.RefreshToken)
}

// func main(){
//   lambda.Start(func(
//     ctx context.Context,
//     request AWSGatewayRequest,
//   )( AWSGatewayResponse, error ){
//     sess, err := session.NewSession(&aws.Config{
//       Region: aws.String(AWS_REGION),
//     })
//     if err != nil {
//       return AWSGatewayResponse{
//         StatusCode: 404,
//         Body: "Error occured while setting up Session.",
//       }, err
//     }
//
//     switch request.Path {
//     case "/signup":
//       return auth.Authorize(ctx, request, sess)
//     case "/authorize":
//       return auth.Authorize(ctx, request, sess)
//     case "/refresh":
//       return token.RefreshToken(ctx, request, sess)
//     default:
//       return AWSGatewayResponse{
//         StatusCode: 404,
//         Body: "Unknown Request",
//       }, nil
//     }
//   })
// }
