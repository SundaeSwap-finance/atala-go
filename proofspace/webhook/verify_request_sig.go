package webhook

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const AuthHeaderName = "Authorization"

func VerifyRequestSignature(headers map[string][]string, publicKey string) error {

	if headers[AuthHeaderName] != nil {
		tokenFromHeader := headers[AuthHeaderName][0]
		tokenFromHeader = tokenFromHeader[strings.Index(tokenFromHeader, " ")+1:]
		token, signatureVerificationErr := jwt.Parse(tokenFromHeader, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				// SIGNATURE VERIFICATION FAILED DURING TOKEN PARSE
				return nil, errors.New("request signature verification failed: unable to parse JWT token")
			}
			pk, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
			if err != nil {
				// SIGNATURE VERIFICATION FAILED DURING PUBLIC KEY PARSE
				fmt.Printf("Error parsing RSA PublicKey from PEM%+v\n", err)
				return nil, err
			}
			return pk, nil

		})
		if signatureVerificationErr != nil {
			return signatureVerificationErr
		}
		if !token.Valid {
			return errors.New("request signature verification failed: invalid JWT token")
		}
		return nil
	}
	return errors.New("token not found in header")
}
