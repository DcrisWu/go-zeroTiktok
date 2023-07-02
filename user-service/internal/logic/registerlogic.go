package logic

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"go-zeroTiktok/model"
	"go-zeroTiktok/user-service/internal/config"
	"go-zeroTiktok/user-service/internal/svc"
	"go-zeroTiktok/user-service/pb/user"
	"go-zeroTiktok/utils"
	"golang.org/x/crypto/argon2"
	"google.golang.org/grpc/status"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {

	id, err := l.CreateUser(in, l.svcCtx.Config.Argon2ID)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		} else {
			return nil, status.Error(500, "注册失败")
		}
	}
	//注册成功直接登陆
	_, err = l.CheckUser(in, l.svcCtx.Config.Argon2ID)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		} else {
			return nil, status.Error(500, "注册失败")
		}
	}
	return &user.RegisterResp{
		Status: utils.SUCCESS,
		UserId: id,
		Token:  utils.GenerateToken(),
	}, nil
}

func (l *RegisterLogic) CreateUser(in *user.RegisterReq, argon2Params *config.Argon2Params) (int64, error) {
	exitUser, err := l.svcCtx.UserModel.FindOneByUserName(l.ctx, in.UserName)
	if err != nil {
		return 0, err
	}
	if exitUser != nil {
		return 0, status.Error(400, "用户名已存在")
	}
	password, err := generateFromPassword(in.Password, argon2Params)
	if err != nil {
		return 0, err
	}
	id, err := utils.NewBasicGenerator().GenerateId()
	if err != nil {
		return 0, err
	}
	_, err = l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Id:       id,
		UserName: in.UserName,
		Password: password,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

// generatePassword generate the hash from the password string with salt and iterations values.
// the encrypting algorithm is Argon2id.
func generateFromPassword(password string, argon2Params *config.Argon2Params) (string, error) {
	salt, err := generateRandomBytes(argon2Params.SaltLength)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, argon2Params.Iterations, argon2Params.Memory, argon2Params.Parallelism, argon2Params.KeyLength)

	// Base64 encode the salt and hashed password.
	base64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argon2Params.Memory, argon2Params.Iterations, argon2Params.Parallelism, base64Salt, base64Hash)

	return encodedHash, nil
}

func generateRandomBytes(saltLength uint32) ([]byte, error) {
	buf := make([]byte, saltLength)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (l *RegisterLogic) CheckUser(in *user.RegisterReq, argon2Params *config.Argon2Params) (int64, error) {
	u, err := l.svcCtx.UserModel.FindOneByUserName(l.ctx, in.UserName)
	if err != nil {
		return 0, err
	}
	if u == nil {
		return 0, status.Error(400, "用户不存在")
	}
	match, err := comparePasswordAndHash(in.Password, u.Password)
	if err != nil {
		return 0, err
	}
	if !match {
		return 0, status.Error(400, "密码错误")
	}
	return u.Id, nil
}

func comparePasswordAndHash(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	argon2Params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the input password using the same parameters.
	inputHash := argon2.IDKey([]byte(password), salt, argon2Params.Iterations, argon2Params.Memory, argon2Params.Parallelism, argon2Params.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, inputHash) == 1 {
		return true, nil
	}
	return false, nil
}

// decodeHash decode the hash of the password from the database.
//
// returns an error if the password is not valid.
func decodeHash(encodedHash string) (argon2Params *config.Argon2Params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, status.Error(500, "Invalid Hash")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, status.Error(500, "Incompatible Version")
	}

	argon2Params = &config.Argon2Params{}
	if _, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &argon2Params.Memory, &argon2Params.Iterations, &argon2Params.Parallelism); err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	argon2Params.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	argon2Params.KeyLength = uint32(len(hash))

	return argon2Params, salt, hash, nil
}
