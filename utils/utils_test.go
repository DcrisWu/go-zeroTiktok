package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	token, err := GenerateJwt("123", time.Now().Unix(), time.Now().Unix()+3600, map[string]interface{}{"uid": "1"})
	assert.NoError(t, err)
	jwt, err := ParseJWT(token, "123")
	assert.NoError(t, err)
	assert.Equal(t, jwt["uid"], "1")
}
