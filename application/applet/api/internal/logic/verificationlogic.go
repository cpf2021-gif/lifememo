package logic

import (
	"context"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"lifememo/application/applet/api/internal/code"
	"lifememo/application/applet/api/internal/svc"
	"lifememo/application/applet/api/internal/types"
	"lifememo/pkg/encrypt"
	"lifememo/pkg/util"

	"github.com/jordan-wright/email"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	prefixVerificationCount = "biz#verification#count#%s"
	prefixActivation        = "biz#activation#%s"
	verificationLimitPerDay = 10
	expireActivation        = 60 * 30
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequest) (*types.VerificationResponse, error) {
	// 验证邮箱格式
	req.Email = strings.TrimSpace(req.Email)
	err := util.ParseEmail(req.Email)
	if err != nil {
		return nil, code.EmailFormatError
	}

	// 加密邮箱
	encEmail, err := encrypt.EncEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// 查询验证码发送次数
	count, err := l.getVerificationCount(encEmail)
	if err != nil {
		logx.Errorf("getVerificationCount error: %v", err)
	}
	if count >= verificationLimitPerDay {
		return nil, code.VerificationCodeLimit
	}

	// 获取验证码(30分钟内有效)
	ActivationCode, err := getActivationCache(encEmail, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("getActivationCache error: %v", err)
	}

	if len(ActivationCode) == 0 {
		ActivationCode = util.RandomNumeric(6)
		// 保存验证码
		err = saveActivationCache(encEmail, ActivationCode, l.svcCtx.BizRedis)
		if err != nil {
			logx.Errorf("saveActivationCache error: %v", err)
			return nil, err
		}
	}

	// 发送验证码
	err = sendVerificationCode(req.Email, ActivationCode, l.svcCtx.Config.Email.Key)
	if err != nil {
		logx.Errorf("sendVerificationCode error: %v", err)
		return nil, code.SendEmailError
	}

	// 增加验证码发送次数
	err = l.incrVerificationCount(encEmail)
	if err != nil {
		logx.Errorf("incrVerificationCount error: %v", err)
	}

	return &types.VerificationResponse{}, nil
}

// 查询验证码发送次数
func (l *VerificationLogic) getVerificationCount(email string) (int, error) {
	key := fmt.Sprintf(prefixVerificationCount, email)
	val, err := l.svcCtx.BizRedis.Get(key)
	if err != nil {
		return 0, err
	}
	if len(val) == 0 {
		return 0, nil
	}

	return strconv.Atoi(val)
}

// 增加验证码发送次数
func (l *VerificationLogic) incrVerificationCount(email string) error {
	key := fmt.Sprintf(prefixVerificationCount, email)
	_, err := l.svcCtx.BizRedis.Incr(key)
	if err != nil {
		return err
	}

	// 设置过期时间
	return l.svcCtx.BizRedis.Expireat(key, util.EndOfDay(time.Now()).Unix())
}

func getActivationCache(email string, rds *redis.Redis) (string, error) {
	key := fmt.Sprintf(prefixActivation, email)
	return rds.Get(key)
}

func saveActivationCache(email, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, email)
	return rds.Setex(key, code, expireActivation)
}

func delActivationCache(email string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, email)
	_, err := rds.Del(key)
	return err
}

func sendVerificationCode(userEmail, code, key string) error {
	// create new email
	e := email.NewEmail()
	// set sender email
	e.From = "lifeMemo <2992247892@qq.com>"
	// set receiver email
	e.To = []string{userEmail}
	// set subject
	e.Subject = "lifeMemo验证码"
	// set content
	e.Text = []byte(code)
	// send email
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "2992247892@qq.com", key, "smtp.qq.com"))
	if err != nil {
		return err
	}

	return nil
}
