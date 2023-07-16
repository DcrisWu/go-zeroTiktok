package utils

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
)

const keyPayload = "payload"
const keyUid = "uid"
const keyExp = "exp"
const keyIat = "iat"
const JwtPrefix = "jwt:"

// GenerateJwt 生成token
// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func GenerateJwt(secretKey string, iat, seconds int64, payload map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set some claims
	token.Claims = jwt.MapClaims{
		keyExp:     iat + seconds,
		keyIat:     iat,
		keyPayload: payload,
	}

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secretKey))

	return tokenString, err
}

func ParseJWT(tokenString string, secretKey string) (map[string]interface{}, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Convert the payload to map[string]interface{}
		payload, ok := claims[keyPayload].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid payload")
		}
		return payload, nil
	} else {
		return nil, err
	}
}

func GetUid(ctx context.Context) int64 {
	pl := ctx.Value(keyPayload).(map[string]interface{})
	s := pl[keyUid]
	uid, err := strconv.ParseInt(s.(string), 10, 64)
	if err != nil {
		return UidNotFound
	}
	return uid
}

func JwtToRedis(ctx context.Context, redis *redis.Redis, uid int64, exp int) error {
	key := JwtPrefix + strconv.FormatInt(uid, 10)
	_, err := redis.SetnxCtx(ctx, key, "1")
	if err != nil {
		return err
	}
	return redis.ExpireCtx(ctx, key, exp)
}

func IsJwtInRedis(ctx context.Context, redis *redis.Redis, uid int64) (bool, error) {
	key := JwtPrefix + strconv.FormatInt(uid, 10)
	return redis.ExistsCtx(ctx, key)
}

func DelJwtInRedis(ctx context.Context, redis *redis.Redis, uid int64) error {
	key := JwtPrefix + strconv.FormatInt(uid, 10)
	_, err := redis.DelCtx(ctx, key)
	return err
}
