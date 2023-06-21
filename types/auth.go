package types

import "github.com/golang-jwt/jwt"

type TokenClaim struct {
	jwt.StandardClaims
	UserID string `json:"userid"`
}
type RequestBody struct {
	AuthData     map[string]string `json:"authData"`
	CustomerData map[string]string `json:"customerData"`
	OtherData    map[string]string `json:"otherData"`
}
