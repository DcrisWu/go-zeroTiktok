package logic

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"go-zeroTiktok/model"
	"go-zeroTiktok/user-service/internal/config"
	"go-zeroTiktok/user-service/internal/svc"
	"go-zeroTiktok/user-service/pb/user"
	"go-zeroTiktok/utils"
	"golang.org/x/crypto/argon2"
	"google.golang.org/grpc/status"

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
	// todo 注册成功直接登陆
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
