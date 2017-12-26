package mysql_test

import (
	"testing"
	"time"

	"openauth/store/token"
)

func TestToken(t *testing.T) {
	s := newTestStore()
	defer s.Close()

	tk := token.Token{
		UserID:       "unit-test-user-01",
		AccessToken:  "aabbcc",
		ClientID:     "client-id",
		GrantType:    token.PASSWORD,
		RefreshToken: "bbccdd",
		TokenType:    "bearer",
		CreatedAt:    time.Now().Unix(),
		ExpiresIn:    3600,
	}

	if _, err := s.SaveToken(&tk); err != nil {
		t.Fatal(err)
	}

	tGet, err := s.GetToken(tk.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if tGet.AccessToken != tk.AccessToken {
		t.Fatal("get access token not equal")
	}

	if err := s.DeleteToken(tk.AccessToken); err != nil {
		t.Fatal(err)
	}

}
