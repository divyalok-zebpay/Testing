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
	Name string `json:"name,omitempty"`
}

func getBodyHash(body any) string {
	var bytesBody []byte
	bytesBody, ok := body.([]byte)
	if !ok {
		bytesBody, _ = json.Marshal(body)
	}
	hasher := sha256.New()
	hasher.Write(bytesBody)
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
	body := []map[string]any{
		{
			"transaction_id": "c6911095-ba83-4aa1-b0fb-15934568a65a",
			"destination":    1,
			"amount":         100,
			"from":           "c6911095-ba83-4aa1-b0fb-15934568a65a",
		},
		{
			"transaction_id": "c6911095-ba83-4aa1-b0fb-15934568a65a",
			"destination":    2,
			"amount":         100,
			"from":           "c6911095-ba83-4aa1-b0fb-15934568a65a",
		},
		{
			"transaction_id": "c6911095-ba83-4aa1-b0fb-15934568a65c",
			"destination":    1,
			"amount":         100.234234,
			"from":           "c6911095-ba83-4aa1-b0fb-15934568a65a",
		},
		{
			"transaction_id": "c6911095-ba83-4aa1-b0fb-15934568a65d",
			"destination":    2,
			"amount":         100,
			"from":           "c6911095-ba83-4aa1-b0fb-15934568a65a",
		},
		{
			"transaction_id": "c6911095-ba83-4aa1-b0fb-15934568a65e",
			"destination":    1,
			"amount":         100,
			"from":           "c6911095-ba83-4aa1-b0fb-15934568a65a",
		},
		{
			"transaction_id": "c6911095-ba83-4aa1-b0fb-15934568a65f",
			"destination":    1,
			"amount":         100,
			"from":           "c6911095-ba83-4aa1-b0fb-15934568a65a",
		},
	}
	uri := "/api/v1/transactions"
	url := fmt.Sprintf("http://localhost:8080%s", uri)
	bod, _ := json.Marshal(body)
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
	// return
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"x-api-key":     "90b2e8bb-ea3c-4849-8fb0-a8072825c2e4",
	}
	status, res, err := external.HTTPCall(&external.HTTPCallParams{
		Client:  httpclient.New(context.Background()),
		Method:  external.HttpMethodPost,
		URL:     url,
		Payload: bod,
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
