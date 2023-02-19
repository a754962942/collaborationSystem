package jwts

import "testing"

func TestParseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzczODE0MTMsInRva2VuIjoiMTAwNyJ9.hkKwbgsXbjoV2uae_kvVtHRpwoAi1_jTlfsP4msHJqs"
	ParseToken(tokenString, "manageSystem")
}
