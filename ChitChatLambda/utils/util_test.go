package utils

import (
	"reflect"
	"testing"
)

// Not sure how to test this. might try again later.
func TestRespond(t *testing.T){
  responses := []struct {
    Message string
    Body    []byte
    code    int
  }{
    {Message: "Unknown Error", Body: []byte("Error Has Occurred"), code: 400},
  }
  for _, tt := range responses{
    resp := Response{ Message: tt.Message, Body: tt.Body }
    got, _ := resp.RespondWith(tt.code)

    want := AWSGatewayResponse{
      StatusCode: tt.code,
      // Body: fmt.Sprintf("message: %v, body: %v"), tt.Message, string(tt.Body),
    }

    if !reflect.DeepEqual(want, got) {
      t.Errorf("Failed: Wanted '%v', got got: '%v'", want, got)
    }
  }

}
