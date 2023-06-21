package service

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/RakhimovAns/Alif/pkg"
	"github.com/RakhimovAns/Alif/types"
	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

var Pool *pgxpool.Pool

func Auth(channel chan *string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			auth := request.Header.Get("X-UserId")
			XDigest := request.Header.Get("X-Digest")
			if auth == "" && XDigest == "" {
				requestBody := Parse(request)
				if requestBody.CustomerData["login"] != "" && requestBody.CustomerData["password"] != "" {
					var id string
					var hash string
					password := requestBody.CustomerData["password"]
					login := requestBody.CustomerData["login"]
					err := Pool.QueryRow(request.Context(), `
    select id, password from customers where login=$1
`, login).Scan(&id, &hash)
					if err == pgx.ErrNoRows {
						log.Println(err)
						return
					}
					hashed, err := hex.DecodeString(hash)
					if err != nil {
						log.Println(err)
						return
					}
					err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
					if err != nil {
						log.Println(err)
						return
					}

					channel <- &id
					next.ServeHTTP(writer, request)
					return
				} else {
					return
				}
			}
			bearerToken := strings.Split(auth, " ")
			if bearerToken[0] != "Bearer" {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if len(bearerToken) != 2 {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			token, id, err := ParseToken(bearerToken[1])
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if token.Valid {
				requestBody := Parse(request)
				Digest, err := pkg.CalculateHMACSHA1(requestBody.AuthData, "Ansar")
				if err != nil {
					http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
				if XDigest != Digest {
					http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
				channel <- &id
				next.ServeHTTP(writer, request)
				return
			}
		})
	}
}

func ParseToken(accessToken string) (*jwt.Token, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &types.TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing methon: %v ", token.Header["alg"])
		}
		return []byte("My Key"), nil
	})
	if err != nil {
		return nil, "", err
	}
	claims, ok := token.Claims.(*types.TokenClaim)
	if !ok {
		return nil, "", err
	}
	return token, claims.UserID, err
}

func Parse(request *http.Request) *types.RequestBody {
	var requestBody *types.RequestBody
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		fmt.Println("Failed to parse request body:", err)
		return nil
	}
	return requestBody
}
