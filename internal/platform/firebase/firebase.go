package firebase

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

const firebasePublicKeyURL = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"

var ErrInvalidToken = errors.New("invalid firebase token")

type Claims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	jwt.RegisteredClaims
}

func VerifyIDToken(idToken, projectID string) (*Claims, error) {
	resp, err := http.Get(firebasePublicKeyURL)
	if err != nil {
		return nil, fmt.Errorf("fetch firebase public keys: %w", err)
	}
	defer resp.Body.Close()

	var certs map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&certs); err != nil {
		return nil, fmt.Errorf("decode firebase public keys: %w", err)
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(idToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("missing kid in token header")
		}

		certPEM, ok := certs[kid]
		if !ok {
			return nil, errors.New("unknown kid in firebase public keys")
		}

		block, _ := pem.Decode([]byte(certPEM))
		if block == nil {
			return nil, errors.New("failed to decode PEM block")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate: %w", err)
		}

		rsaKey, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("not an RSA public key")
		}

		return rsaKey, nil
	}, jwt.WithAudience(projectID))

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	expectedIss := "https://securetoken.google.com/" + projectID
	if iss, _ := claims.GetIssuer(); iss != expectedIss {
		return nil, fmt.Errorf("%w: unexpected issuer %s", ErrInvalidToken, iss)
	}

	if !claims.EmailVerified {
		return nil, fmt.Errorf("%w: email not verified", ErrInvalidToken)
	}

	return claims, nil
}
