package logic

import (
	"context"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/types"
	"zeebee-backend/common/utils"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	logicpkg "zeebee-backend/domains/zeebee/logic"
)

var baseLogic *BaseLogic
var onceBase sync.Once

// BaseLogic is the tenant logic entrypoint. It embeds ONE shared AllLogics (used
// for non-token operations) and owns the per-platform concerns tenant must handle
// itself because the shared layer is single-platform: jwtConfigs (issuance),
// claim resolution, and user_type parsing.
type BaseLogic struct {
	svcCtx *svc.ServiceContext
	*logicpkg.AllLogics
	jwtConfigs   map[constants.UserType]dto.JwtConfig
	platforms    map[constants.UserType]config.PlatformProfile
	emailClients map[constants.UserType]*email.Email
}

func GetBaseLogic(svcCtx *svc.ServiceContext) *BaseLogic {
	onceBase.Do(func() {
		baseLogic = NewBaseLogic(svcCtx)
	})
	return baseLogic
}

func NewBaseLogic(svcCtx *svc.ServiceContext) *BaseLogic {
	l := &BaseLogic{
		svcCtx: svcCtx,
		AllLogics: logicpkg.NewAllLogics(logicpkg.AllLogicConfig{
			Repos:           svcCtx.AllRepos,
			RedisCli:        svcCtx.RedisCli,
			CurrentEnforcer: svcCtx.Enforcer,
			Logger:          svcCtx.Logger,
		}),
		jwtConfigs:   buildJwtConfigs(svcCtx.Config.PlatformConf),
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

func buildJwtConfigs(pc config.PlatformConf) map[constants.UserType]dto.JwtConfig {
	return map[constants.UserType]dto.JwtConfig{
		constants.UserTypeBorrower: toJwtConfig(pc.Borrower.JwtAuth),
		constants.UserTypeLender:   toJwtConfig(pc.Lender.JwtAuth),
		constants.UserTypeAdmin:    toJwtConfig(pc.Admin.JwtAuth),
	}
}

func toJwtConfig(j types.JwtAuth) dto.JwtConfig {
	return dto.JwtConfig{
		AccessSecret:  j.AccessSecret,
		AccessExpire:  j.AccessExpire,
		RefreshSecret: j.RefreshSecret,
		RefreshExpire: j.RefreshExpire,
	}
}

// JwtConfigFor returns the platform JwtConfig; a missing entry is a config error.
func (l *BaseLogic) JwtConfigFor(userType constants.UserType) (dto.JwtConfig, error) {
	c, ok := l.jwtConfigs[userType]
	if !ok {
		return dto.JwtConfig{}, errorx.New("jwt config not configured for user type: " + string(userType))
	}
	return c, nil
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

// GetUserClaimFromCtx resolves the current user from the incoming gRPC token
// metadata. tenant serves all platforms, so it tries each platform's access
// secret and returns the claim that parses.
func (l *BaseLogic) GetUserClaimFromCtx(ctx context.Context, userType constants.UserType) (jwtUser utils.JwtUser, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return utils.JwtUser{}, errorx.WithStack(status.Error(codes.Unauthenticated, "not login"))
	}
	tokens := md[utils.JwtTokenCtxKey]
	if len(tokens) == 0 {
		return utils.JwtUser{}, errorx.WithStack(status.Error(codes.Unauthenticated, "not login"))
	}
	token := tokens[0]
	jwtConf, err := l.JwtConfigFor(userType)
	if err != nil {
		return utils.JwtUser{}, errorx.WithStack(status.Error(codes.InvalidArgument, err.Error()))
	}
	claim, err := utils.ParseToken(token, jwtConf.AccessSecret)
	if err != nil {
		err = errorx.WithStack(status.Error(codes.InvalidArgument, err.Error()))
	}
	jwtUser = claim.JwtUser
	return
}
