package main

import (
	"brave/claims"
	external "brave/helper"
	"brave/helper/httpclient"
	"brave/utility"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Request struct {
	Name string
}

func getBodyHash(req interface{}) string {
	data, _ := json.Marshal(req)
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func getSecretKey() *[]byte {
	data, err := ioutil.ReadFile("./zebBrave_secret.key")
	if err != nil {
		fmt.Println("File reading error", err)
	}
	return &data
}

func main() {
	body := Request{}
	uri := "/api/v1/balance"
	url := fmt.Sprintf("http://localhost:8080%s", uri)
	bod, _ := json.Marshal(body)
	fmt.Printf("bod: %v\n\n", bod)
	token_details := map[string]interface{}{
		"bodyHash": getBodyHash(bod),
		"uri":      uri,
		"sub":      "90b2e8bb-ea3c-4849-8fb0-a8072825c2e4",
	}
	token, err := utility.GetTokenWithExpiry(&claims.ZebBraveClaims{}, time.Now().Add(25*time.Second).Unix(), getSecretKey(), token_details)
	if err != nil {
		fmt.Printf("token generation failed with error: %s\n", err.Error())
		return
	}
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"x-api-key":     "90b2e8bb-ea3c-4849-8fb0-a8072825c2e4",
	}
	pl, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("error1: %s\n", err.Error())
		return
	}
	status, res, err := external.HTTPCall(&external.HTTPCallParams{
		Client:  httpclient.New(context.Background()),
		Method:  external.HttpMethodGet,
		URL:     url,
		Payload: pl,
		Headers: headers,
	})
	if err != nil {
		fmt.Printf("error2: %s\n", err.Error())
		return
	}
	if status < 200 || status >= 300 {
		fmt.Printf("error3: %s status: %v\n", string(res), status)
		return
	}

	fmt.Printf("This is the response: %s\n", string(res))
}
