package main

import (
	auth "github.com/TylerAldrich814/ChitChatMono/ChitChatLamda/auth/authorize"
	"github.com/aws/aws-lambda-go/lambda"
)

func main(){
  lambda.Start(auth.Authorize)
}
