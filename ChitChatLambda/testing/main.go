package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)
type AWSGatewayResponse = events.APIGatewayProxyResponse
type AWSGatewayRequest = events.APIGatewayProxyRequest

func Test1(
  ctx context.Context,
  request AWSGatewayRequest,
)( AWSGatewayResponse, error ){
  // req := request.Path
  // msg := fmt.Sprintf("PATH: %v", req)

  return AWSGatewayResponse{
    Body: "FUCK",
    StatusCode: 200,
  }, nil
}

func main(){
  lambda.Start(func(
    ctx context.Context,
    request AWSGatewayRequest,
  )( AWSGatewayResponse, error ){
    return Test1(ctx, request)
    // switch request.Path {
    // case "/Test1":
    //   return Test1(ctx, request)
    // // case "/Test2":
    // // case "/Test3":
    // default:
    //   return AWSGatewayResponse{
    //     Body: "Provided Path is missing",
    //     StatusCode: 400,
    //   }, nil
    // }
  })
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
//
