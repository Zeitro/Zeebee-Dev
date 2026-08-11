package userlogic

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/utils"
)

type GetClaimLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClaimLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClaimLogic {
	return &GetClaimLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetClaimLogic) GetClaim(in *tenant.GetClaimRequest) (*tenant.GetClaimResponse, error) {
	if in == nil || in.Data == nil || in.Data.AccessToken == "" {
		return nil, status.Error(codes.Unauthenticated, "not login")
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return nil, err
	}
	claim, err := utils.ParseToken(in.Data.AccessToken, l.svcCtx.Config.JwtConf.AccessSecret)
	if err != nil || claim.UserType != userType {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	if claim.UserId == 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return &tenant.GetClaimResponse{Data: &tenant.GetClaimResponseData{
		UserId:             claim.UserId,
		OrganizationId:     claim.OrganizationId,
		OrganizationUserId: claim.OrganizationUserId,
		Username:           claim.Username,
		Role:               claim.Role,
		UserType:           string(claim.UserType),
	}}, nil
}
