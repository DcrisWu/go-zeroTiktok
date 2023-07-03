package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
)

const keyPayload = "payload"
const keyUid = "uid"
const keyExp = "exp"
const keyIat = "iat"

// GenerateToken 生成token
// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func GenerateToken(secretKey string, iat, seconds int64, payload map[string]interface{}) (string, error) {
	claims := make(jwt.MapClaims)
	claims[keyExp] = iat + seconds
	claims[keyIat] = iat
	marshal, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(marshal)
	claims[keyPayload] = str
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func GetUid(ctx context.Context) int64 {
	pl, _ := ctx.Value(keyPayload).(string)
	decodeString, err := base64.StdEncoding.DecodeString(pl)
	if err != nil {
		return PayLoadNotFound
	}
	plmap := make(map[string]interface{})
	err = json.Unmarshal(decodeString, &plmap)
	if err != nil {
		return PayLoadNotFound
	}
	if _, ok := plmap[keyUid]; !ok {
		return UidNotFound
	}
	uid, err := strconv.ParseInt(plmap[keyUid].(string), 10, 64)
	if err != nil {
		return UidNotFound
	}
	return uid
}
