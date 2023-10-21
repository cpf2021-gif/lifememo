package logic

import (
	"context"
	"lifememo/application/applet/api/internal/code"
	"lifememo/application/user/rpc/user"
	"lifememo/pkg/encrypt"
	"lifememo/pkg/jwt"
	"lifememo/pkg/util"
	"lifememo/pkg/xcode"
	"strings"

	"lifememo/application/applet/api/internal/svc"
	"lifememo/application/applet/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	// 校验字段
	req.Email = strings.TrimSpace(req.Email)
	if len(req.Email) == 0 {
		return nil, code.EmailEmpty
	}

	err := util.ParseEmail(req.Email)
	if err != nil {
		return nil, code.EmailFormatError
	}

	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	if len(req.VerificationCode) == 0 {
		return nil, code.VerificationCodeEmpty
	}

	err = checkVerificationCode(l.svcCtx.BizRedis, req.Email, req.VerificationCode)
	if err != nil {
		return nil, code.VerificationCodeError
	}
	userEmail, err := encrypt.EncEmail(req.Email)
	if err != nil {
		logx.Errorf("EncEmail error: %v", err)
		return nil, err
	}

	u, err := l.svcCtx.UserRpc.FindByEmail(l.ctx, &user.FindByEmailRequest{
		Email: userEmail,
	})

	if err != nil {
		logx.Errorf("FindByEmail error: %v", err)
		return nil, err
	}

	if u == nil || u.UserId == 0 {
		return nil, xcode.AccessDenied
	}

	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": u.UserId,
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		UserId: u.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil

}
