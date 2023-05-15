package utility

import (
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type Claims interface {
	jwt.Claims
	SetExpiration(exp int64)
	BeforeCreate()
}

func GetTokenWithExpiry(claims Claims, expiry int64, key *[]byte, token_details map[string]interface{}) (string, error) {
	claims.SetExpiration(expiry)
	return getToken(claims, token_details, key)
}

func getToken(claims Claims, token_details map[string]interface{}, key *[]byte) (string, error) {
	jsonbody, err := json.Marshal(token_details)
	if err != nil {
		return "", fmt.Errorf("token generation with err: %v", err.Error())
	}
	if err := json.Unmarshal(jsonbody, claims); err != nil {
		return "", fmt.Errorf("token generation with err: %v", err.Error())
	}
	claims.BeforeCreate()
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(*key)
	if err != nil {
		return "", fmt.Errorf("error parsing RSA private key: %v", err)
	}
	signedToken, err := t.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("error in token creation with error: %v", err)
	}
	return signedToken, nil
}
