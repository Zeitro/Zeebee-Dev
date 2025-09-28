package loan

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"zeebee-backend/app/lender/api/internal/svc"
	"zeebee-backend/app/lender/api/internal/types"
	"zeebee-backend/app/lender/rpc/pb/lender"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.LoanUpdateRequest) (resp *types.LoanUpdateResponse, err error) {
	rpcReq := &lender.LoanUpdateRequest{
		Data: &lender.LoanUpdateRequestData{},
	}
	resp = &types.LoanUpdateResponse{}
	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		IgnoreEmpty:   false,
		DeepCopy:      true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	// get lender id
	rpcReq.Data.LenderCompanyId, err = utils.GetLenderCompanyIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// get employee id
	rpcReq.Data.LenderEmployeeId, err = utils.GetLenderAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	applicants := make([]*lender.ApplicantEntry, 0, len(rpcReq.Data.Applicants))
	for _, applicant := range rpcReq.Data.Applicants {
		applicant := applicant
		if len(applicant.FirstName) > 0 {
			applicants = append(applicants, applicant)
		}
	}
	if len(applicants) > 0 {
		primaryCount := slice.CountBy(applicants, func(index int, item *lender.ApplicantEntry) bool {
			return item.IsPrimary
		})
		if primaryCount == 0 {
			applicants[0].IsPrimary = true
		}
		if primaryCount > 1 {
			return nil, errors.WithStack(errors.New("multiple primary applicant"))
		}
	}
	rpcReq.Data.Applicants = applicants

	rpcResp, err := l.svcCtx.LoanRpcClient.UpdateLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}
