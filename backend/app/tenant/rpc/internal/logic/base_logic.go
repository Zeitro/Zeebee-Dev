package logic

import (
	"context"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/types"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	logicpkg "zeebee-backend/domains/zeebee/logic"
)

var baseLogic *BaseLogic
var onceBase sync.Once

// BaseLogic is the tenant logic entrypoint. It embeds ONE shared AllLogics and
// owns the per-platform concerns tenant must handle itself: platform profiles
// (brand / URLs / expiry) and SMTP clients. JWT signing uses the single shared
// key (AllLogicConfig.JwtConfig); userType lives in the token claim.
type BaseLogic struct {
	svcCtx *svc.ServiceContext
	*logicpkg.AllLogics
	platforms    map[constants.UserType]config.PlatformProfile
	emailClients map[constants.UserType]*email.Email
}

func GetBaseLogic(svcCtx *svc.ServiceContext) *BaseLogic {
	onceBase.Do(func() {
		baseLogic = NewBaseLogic(svcCtx)
	})
	return baseLogic
}

// ResetBaseLogicForTest clears the cached singleton so each test wires its own
// ServiceContext. Test-only; production always uses the once-cached instance.
func ResetBaseLogicForTest() {
	baseLogic = nil
	onceBase = sync.Once{}
}

func NewBaseLogic(svcCtx *svc.ServiceContext) *BaseLogic {
	// tenant runs ONE unfiltered Casbin enforcer (central authority). All platforms
	// publish role-reload events to the same watcher channel so the enforcer reloads
	// regardless of which platform's data changed.
	topic := svcCtx.Config.CasbinConf.Watcher.Channel
	authTopicMap := map[constants.UserType]string{
		constants.UserTypeBorrower: topic,
		constants.UserTypeLender:   topic,
		constants.UserTypeAdmin:    topic,
	}

	l := &BaseLogic{
		svcCtx: svcCtx,
		AllLogics: logicpkg.NewAllLogics(logicpkg.AllLogicConfig{
			Repos:                 svcCtx.AllRepos,
			RedisCli:              svcCtx.RedisCli,
			JwtConfig:             toJwtConfig(svcCtx.Config.JwtConf),
			CurrentEnforcer:       svcCtx.Enforcer,
			CurrentAuthEventTopic: topic,
			AuthEventTopicMap:     authTopicMap,
			Logger:                svcCtx.Logger,
		}),
		platforms:    buildPlatforms(svcCtx.Config.PlatformConf),
		emailClients: buildEmailClients(svcCtx.Config.PlatformConf),
	}

	// Rebuild each platform's policy in the background (central authority loads all).
	for _, userType := range []constants.UserType{constants.UserTypeBorrower, constants.UserTypeLender, constants.UserTypeAdmin} {
		threading.GoSafe(func() {
			ctx := context.Background()
			if err := l.AuthLogic.RebuildPolicy(ctx, userType); err != nil {
				svcCtx.Logger.WithContext(ctx).WithFields(
					logx.Field("stacks", errorx.GetErrorStack(err)),
					logx.Field("user_type", userType),
				).Error("fail to rebuild tenant policy")
			}
		})
	}

	return l
}

func toJwtConfig(j types.JwtConf) dto.JwtConfig {
	return dto.JwtConfig{
		AccessSecret:  j.AccessSecret,
		AccessExpire:  j.AccessExpire,
		RefreshSecret: j.RefreshSecret,
		RefreshExpire: j.RefreshExpire,
	}
}

func buildPlatforms(pc config.PlatformConf) map[constants.UserType]config.PlatformProfile {
	return map[constants.UserType]config.PlatformProfile{
		constants.UserTypeBorrower: pc.Borrower,
		constants.UserTypeLender:   pc.Lender,
		constants.UserTypeAdmin:    pc.Admin,
	}
}

// buildEmailClients builds one SMTP client per platform from PlatformConf.
// Gmail OAuth sending is deferred (per scope decision: SMTP first).
func buildEmailClients(pc config.PlatformConf) map[constants.UserType]*email.Email {
	clients := map[constants.UserType]*email.Email{}
	for userType, profile := range buildPlatforms(pc) {
		ec := profile.EmailConf
		if ec.Host == "" || ec.Port == 0 {
			continue
		}
		clients[userType] = email.NewEmail(email.SMTPConfig{
			Host:     ec.Host,
			Port:     ec.Port,
			User:     ec.Username,
			Password: ec.Password,
		})
	}
	return clients
}

// PlatformFor returns the platform profile (brand, URLs, expiry); missing = config error.
func (l *BaseLogic) PlatformFor(userType constants.UserType) (config.PlatformProfile, error) {
	p, ok := l.platforms[userType]
	if !ok {
		return config.PlatformProfile{}, errorx.New("platform not configured for user type: " + string(userType))
	}
	return p, nil
}

// EmailClientFor returns the platform SMTP client; missing = config error.
func (l *BaseLogic) EmailClientFor(userType constants.UserType) (*email.Email, error) {
	c, ok := l.emailClients[userType]
	if !ok {
		return nil, errorx.New("email client not configured for user type: " + string(userType))
	}
	return c, nil
}

// ParseUserType maps the request user_type string (case-insensitive) to a constants.UserType.
func ParseUserType(s string) (constants.UserType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "borrower":
		return constants.UserTypeBorrower, nil
	case "lender":
		return constants.UserTypeLender, nil
	case "admin":
		return constants.UserTypeAdmin, nil
	default:
		return "", errorx.New("invalid user type")
	}
}
