package oauth2_test

import (
	"testing"

	"openauth/pkg/oauth2"
	"openauth/store/token"
)

func TestValidate(t *testing.T) {
	t.Run("PasswordAuthOK", testValidatePasswordTokenOK)
	t.Run("PasswordAuthNotFound", testValidatePasswordTokenNotFound)
	t.Run("TokenExpired", testValidateTokenExpired)
}

func testValidatePasswordTokenOK(t *testing.T) {
	authsvr := NewOAuth2Ctroller(3600)

	_, err := authsvr.ValidateToken("validated-token-string")
	if err != nil {
		t.Fatal(err)
	}
}

func testValidatePasswordTokenNotFound(t *testing.T) {
	authsvr := NewOAuth2Ctroller(3600)

	_, err := authsvr.ValidateToken("validated-token-string-not-found")
	if err == nil {
		t.Fatal("want an not found error")
	}
}

func testValidateTokenExpired(t *testing.T) {
	authsvr := NewOAuth2Ctroller(0)

	_, err := authsvr.ValidateToken("validated-token-string")
	if err.Error() != "token has expired" {
		t.Fatal("want token expired")
	}
}

func TestIssueToken(t *testing.T) {
	t.Run("PasswordOK", testIssueByPassOK)
}

func testIssueByPassOK(t *testing.T) {
	authsvr := NewOAuth2Ctroller(3600)

	req := oauth2.TokenRequest{
		GrantType:    token.PASSWORD,
		ClientID:     "validated-client-id",
		ClientSecret: "validated-client-secret",
		DomainName:   "validated-domain",
		Username:     "validated-user",
		Password:     "validated-pass",
	}
	_, err := authsvr.IssueToken(&req)
	if err != nil {
		t.Fatal(err)
	}
}
