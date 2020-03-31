package core

import "testing"

func TestTokenManager_CreateToken(t *testing.T) {
	tm := TokenManager{}
	tm.CreateToken("a", "b")
	tm.CreateToken("a", "d")
	tm.CreateToken("a", "c")

}
