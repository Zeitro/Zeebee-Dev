package organizationsubscriptionlogic

import (
	"context"

	"zeebee-backend/app/tenant/rpc/internal/logic"
	"zeebee-backend/app/tenant/rpc/internal/svc"
	"zeebee-backend/app/tenant/rpc/pb/tenant"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
)

type PageLogic struct {
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	baseLogic *logic.BaseLogic
}

func NewPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageLogic {
	return &PageLogic{ctx: ctx, svcCtx: svcCtx, baseLogic: logic.GetBaseLogic(svcCtx)}
}

// Page returns a paginated list of organization subscriptions.
func (l *PageLogic) Page(in *tenant.OrganizationSubscriptionPageRequest) (resp *tenant.OrganizationSubscriptionPageResponse, err error) {
	if in.Data == nil {
		return nil, errorx.ErrorEmptyRequestData
	}
	userType, err := logic.ParseUserType(in.Data.UserType)
	if err != nil {
		return
	}

	// copierx maps page/limit/sorter/keyword and tolerates a request that omits
	// the filter or sorter (they map to zero values).
	pageReq, err := copierx.Trans[tenant.OrganizationSubscriptionPageRequestData, dto.OrganizationSubscriptionPageRequestDto](in.Data)
	if err != nil {
		return
	}
	if userType == constants.UserTypeAdmin {
		if in.Data.Filter != nil && in.Data.Filter.Type != "" {
			if _, err = logic.ParseUserType(in.Data.Filter.Type); err != nil {
				return
			}
		}
	} else {
		if in.Data.OrganizationId <= 0 {
			return nil, errorx.New("organization id is required")
		}
		pageReq.Filter.Type = userType
		pageReq.Filter.OrganizationId = in.Data.OrganizationId
	}
	// Status crosses a repeated-string -> []OrganizationSubscriptionStatus
	// boundary, so derive it explicitly (nil filter => no status filter).
	pageReq.Filter.Status = nil
	if in.Data.Filter != nil {
		for _, s := range in.Data.Filter.Status {
			pageReq.Filter.Status = append(pageReq.Filter.Status, constants.OrganizationSubscriptionStatus(s))
		}
	}

	page, err := l.baseLogic.OrganizationSubscriptionLogic.Page(l.ctx, *pageReq)
	if err != nil {
		return
	}

	data, err := copierx.Trans[gormx.PageResponse[dto.OrganizationSubscriptionPageEntity], tenant.OrganizationSubscriptionPageResponseData](
		&page, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToRFC3339Str(),
	)
	if err != nil {
		return
	}
	return &tenant.OrganizationSubscriptionPageResponse{Data: data}, nil
}
