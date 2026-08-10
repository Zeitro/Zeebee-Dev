package logic

import (
	"context"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"

	"zeebee-backend/app/tenant/config"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/email"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/types"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	logicpkg "zeebee-backend/domains/zeebee/logic"
	"zeebee-backend/domains/zeebee/model"
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
	// tenant runs ONE unfiltered Casbin enforcer (central authority) and is the
	// sole policy executor, so it subscribes to every platform's reload channel.
	// Each platform still publishes on its OWN channel — a borrower change never
	// wakes lender/admin subscribers. Reloads are scoped per msg.Type, so cross-
	// platform noise is bounded even though one enforcer holds all policies.
	eventTopic := svcCtx.Config.CasbinConf.EventTopic
	authTopicMap := map[constants.UserType]string{
		constants.UserTypeBorrower: eventTopic.BorrowerTopic,
		constants.UserTypeLender:   eventTopic.LenderTopic,
		constants.UserTypeAdmin:    eventTopic.AdminTopic,
	}

	l := &BaseLogic{
		svcCtx: svcCtx,
		AllLogics: logicpkg.NewAllLogics(logicpkg.AllLogicConfig{
			Repos:                  svcCtx.AllRepos,
			RedisCli:               svcCtx.RedisCli,
			JwtConfig:              toJwtConfig(svcCtx.Config.JwtConf),
			CurrentEnforcer:        svcCtx.Enforcer,
			CurrentAuthEventTopics: []string{eventTopic.BorrowerTopic, eventTopic.LenderTopic, eventTopic.AdminTopic},
			AuthEventTopicMap:      authTopicMap,
			Logger:                 svcCtx.Logger,
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

// RecordAudit writes a security audit event to the Operation table (NFR-003).
// Async + best-effort: runs on a detached context so the request returning (and
// canceling its ctx) cannot abort the write, and a store failure is logged but
// never returned. Callers stay synchronous — no goroutine needed at call sites.
func (l *BaseLogic) RecordAudit(ctx context.Context, input dto.OperationCreateDto) {
	auditCtx := context.WithoutCancel(ctx)
	threading.GoSafe(func() {
		op, err := copierx.Trans[dto.OperationCreateDto, model.Operation](&input)
		if err != nil {
			l.svcCtx.Logger.WithContext(auditCtx).
				WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
				Error("fail to map audit record")
			return
		}
		if _, err := l.svcCtx.OperationRepo.Create(auditCtx, op); err != nil {
			l.svcCtx.Logger.WithContext(auditCtx).
				WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
				WithFields(logx.Field("action", input.Operation)).
				Error("fail to record audit")
		}
	})
}

// RecordAuditFromClaim records an audit event attributed to the caller in ctx
// (operator id/type resolved from the JWT claim). Best-effort: a missing claim
// or store failure is logged, never returned, and never blocks the caller. Use
// RecordAudit directly for pre-auth events (login) where there is no claim yet.
func (l *BaseLogic) RecordAuditFromClaim(ctx context.Context, entityType constants.OperationEntityType, entityId int64, operation constants.OperationAction) {
	claim, err := l.UserLogic.GetUserClaimFromCtx(ctx)
	if err != nil {
		l.svcCtx.Logger.WithContext(ctx).
			WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
			WithFields(logx.Field("operation", operation)).
			Error("fail to resolve audit operator from claim")
		return
	}
	l.RecordAudit(ctx, dto.OperationCreateDto{
		EntityType:   entityType,
		EntityId:     entityId,
		OperatorId:   claim.UserId,
		OperatorType: claim.UserType,
		Operation:    operation,
	})
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
