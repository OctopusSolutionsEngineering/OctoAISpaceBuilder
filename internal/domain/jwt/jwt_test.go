package jwt

import (
	"crypto/rsa"
	"testing"
)

// An expired token to test with
const testJwt = "eyJhbGciOiJQUzI1NiIsImtpZCI6IjU2YjI4ZWM4N2E1MzQ3ZWQ4ZjdjZGNhZDJkYzU1ZWE0IiwidHlwIjoiSldUIn0.eyJhdWQiOiJodHRwOi8vbG9jYWxob3N0OjgwODAiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwODAiLCJleHAiOjE3NDQxNzYzNDcsImlhdCI6MTc0NDE3Mjc0NywibmJmIjoxNzQ0MTcyNzQ3LCJqdGkiOiJkZDA0N2U0YWVlMDI0ZmI3YmFiNmUzODhkNzAwMzFlOSIsInN1YiI6ImY2NDY1M2QxLWRhNjctNDc4Zi05OTY5LTA4OGFjN2YyZmU1MiJ9.PY9mZjz6YsGCG-gTYU89KQPrem7ecx5mxqOboKj4RG80ILb03-NLPwejliW_qTwWYthbalwtM0R2rGont0k2YRNirIT0g37ivz1qPOPWpjZQXyF0OS6ulYgIKJXINjWu2EY1uqqpyskZ5kVTIQ6p9CQAssTjr2qeBgSZfdBNSEhs2TU0Qvc4GMkjloU8OxT-UeDXSMqewoIQkrlw-dCupmF-tzEBrPk8p6xdt9Bc6CyfcL6Bj5JINAu2GsazOWKf7twiSEDEonwEDxcd3ssNvQZmxVTV49gEym184JsNXI9bPJGD38bvGgNY7YXYd9RPpg8CQw-IoI4ZoQRTRrnVJw"

func TestGetAud(t *testing.T) {
	if aud, err := GetJwtAud(testJwt); err != nil || aud != "http://localhost:8080" {
		t.Fatal("failed to get JWT audience, was " + aud)
	}
}

func TestValidateJWTWithJWKS(t *testing.T) {
	jwksURL := "http://localhost:8080/.well-known/jwks"

	err := ValidateJWTWithJWKS(testJwt, jwksURL, func(_, kid string) (*rsa.PublicKey, error) {
		return GetKeyFromJWKSByteArray([]byte(`{
		  "keys": [
			{
			  "e": "AQAB",
			  "kid": "56b28ec87a5347ed8f7cdcad2dc55ea4",
			  "kty": "RSA",
			  "n": "yZj-ASROgl5xbl80snLLc1djqJmt1HlIIgomy9rjXiCNmpJZwKsGWoUxMZBanvr2qgMWW0e73mHoaZSOASvkvCnc3l_4h0MP-C98Ogt2q1z3oQBtLDcIr_zLxK8n3pJIJV9x-OM9dDW7ESLOKMW8syAJ6V5y24IwWzTUT5kEeLE6q1JAV1vDnSGu-crYHYrDUZpPh2vXL_SMEL19CKfjdYKE_Pg83S8FBA7JRZ_s4oWDBuCw1kIZvZRlX0S88yivazslKgWQG1Sbczph2kNvuo6ypzCX3rhY6YosyIYlCdPpnHJ-Ej6EfsL-7HAMay7G9_7TaZ-BRb3TaTnqtppFIQ",
			  "use": "sig"
			}
		  ]
		}`), kid)
	}, true)

	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}
}
