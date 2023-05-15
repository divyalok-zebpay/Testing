package claims

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

type ZebBraveClaims struct {
	jwt.StandardClaims
	URI      string `json:"uri"`
	Nonce    string `json:"nonce"`
	IAT      int64  `json:"iat"`
	Sub      string `json:"sub"`
	BodyHash string `json:"bodyHash"`
}

func (b *ZebBraveClaims) Valid() error {
	if !b.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("token has expired")
	}
	return nil
}

func (b *ZebBraveClaims) SetExpiration(exp int64) {
	b.ExpiresAt = exp
}

func (b *ZebBraveClaims) BeforeCreate() {
	tokenID, _ := uuid.NewV4()
	b.Nonce = tokenID.String()
	b.IAT = time.Now().Unix()
}
