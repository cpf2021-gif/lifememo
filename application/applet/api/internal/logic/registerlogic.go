package logic

import (
	"context"
	"strings"

	"lifememo/application/applet/api/internal/code"
	"lifememo/application/applet/api/internal/svc"
	"lifememo/application/applet/api/internal/types"
	"lifememo/application/user/rpc/user"
	"lifememo/pkg/encrypt"
	"lifememo/pkg/jwt"
	"lifememo/pkg/util"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (*types.RegisterResponse, error) {
	// 校验参数
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) == 0 {
		return nil, code.RegisterNameEmpty
	}
	if len(req.Name) > 20 {
		return nil, code.RegisterNameLimit
	}

	req.Email = strings.TrimSpace(req.Email)
	if len(req.Email) == 0 {
		return nil, code.EmailEmpty
	}
	err := util.ParseEmail(req.Email)
	if err != nil {
		return nil, code.EmailFormatError
	}

	// 加密邮箱
	email, err := encrypt.EncEmail(req.Email)

	// 邮箱是否已经注册
	u, err := l.svcCtx.UserRpc.FindByEmail(l.ctx, &user.FindByEmailRequest{
		Email: email,
	})

	if err != nil {
		logx.Errorf("FindByEmail error: %v", err)
		return nil, err
	}

	if u != nil && u.UserId > 0 {
		return nil, code.RegisterEmailExist
	}

	// 校验验证码
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	if len(req.VerificationCode) == 0 {
		return nil, code.VerificationCodeEmpty
	}

	err = checkVerificationCode(l.svcCtx.BizRedis, email, req.VerificationCode)
	if err != nil {
		logx.Errorf("checkVerificationCode error: %v", err)
		return nil, err
	}

	// 注册
	reqRet, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterRequest{
		Username: req.Name,
		Email:    email,
	})
	if err != nil {
		logx.Errorf("Register error: %v", err)
		return nil, err
	}

	// 生成token
	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": reqRet.UserId,
		},
	})
	if err != nil {
		logx.Errorf("BuildTokens error: %v", err)
		return nil, err
	}

	// 删除验证码缓存
	_ = delActivationCache(email, l.svcCtx.BizRedis)

	return &types.RegisterResponse{
		UserId: reqRet.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}

func checkVerificationCode(rds *redis.Redis, email, ActivationCode string) error {
	cacheCode, err := getActivationCache(email, rds)
	if err != nil {
		return err
	}
	if cacheCode == "" {
		return code.VerificationCodeExpired
	}
	if cacheCode != ActivationCode {
		return code.VerificationCodeError
	}
	return nil
}
