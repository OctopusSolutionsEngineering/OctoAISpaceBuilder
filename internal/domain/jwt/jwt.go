package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"io"
	"math/big"
	"net/http"
)

type JWKS struct {
	Keys []json.RawMessage `json:"keys"`
}

type GetKey func(jwksURL, kid string) (*rsa.PublicKey, error)

func GetKeyFromJWKSByteArray(body []byte, kid string) (*rsa.PublicKey, error) {
	var jwks JWKS
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}

	// Find the key with the matching kid
	for _, key := range jwks.Keys {
		var k map[string]interface{}
		if err := json.Unmarshal(key, &k); err != nil {
			continue
		}

		if k["kid"] == kid {
			// Extract modulus (n) and exponent (e)
			modulus, err := base64.RawURLEncoding.DecodeString(k["n"].(string))
			if err != nil {
				return nil, fmt.Errorf("failed to decode modulus: %w", err)
			}

			exponent, err := base64.RawURLEncoding.DecodeString(k["e"].(string))
			if err != nil {
				return nil, fmt.Errorf("failed to decode exponent: %w", err)
			}

			// Convert exponent to integer
			e := 0
			for _, b := range exponent {
				e = e<<8 + int(b)
			}

			// Construct RSA public key
			pubKey := &rsa.PublicKey{
				N: new(big.Int).SetBytes(modulus),
				E: e,
			}
			return pubKey, nil
		}
	}

	return nil, errors.New("key not found in JWKS")
}

func GetKeyFromJWKS(jwksURL, kid string) (*rsa.PublicKey, error) {
	// Fetch the JWKS
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			zap.L().Error("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS response: %w", err)
	}

	return GetKeyFromJWKSByteArray(body, kid)
}

func ValidateJWTWithJWKS(tokenString string, jwksURL string, getKey GetKey, skipValidation bool) error {
	options := []jwt.ParserOption{}

	if skipValidation {
		options = append(options, jwt.WithoutClaimsValidation())
	}

	// Parse the token to extract the kid
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSAPSS); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the kid from the header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("missing kid in token header")
		}

		// Fetch the public key from the JWKS
		return getKey(jwksURL, kid)
	}, options...)

	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

func GetJwtAud(tokenString string) (string, error) {
	// Parse the token without validation
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if aud, ok := claims["aud"]; ok {
			return aud.(string), nil
		}
	}

	return "", errors.New("aud claim not found")
}

func ValidateJWT(tokenString string, server string, skipValidation bool) error {
	return ValidateJWTWithJWKS(tokenString, server+"/.well-known/jwks", GetKeyFromJWKS, skipValidation)
}
