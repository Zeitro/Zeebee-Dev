package loan

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// apply
func NewApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyLogic {
	return &ApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyLogic) Apply(req *types.LoanApplyRequest) (resp *types.LoanApplyResponse, err error) {
	rpcReq := &borrower.LoanApplyRequest{
		Data: &borrower.LoanApplyRequestData{},
	}
	resp = &types.LoanApplyResponse{}
	rpcResp := &borrower.LoanApplyResponse{}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		IgnoreEmpty:   true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	applicants := make([]*borrower.ApplicantEntry, 0, len(rpcReq.Data.Applicants))
	for _, applicant := range rpcReq.Data.Applicants {
		applicant := applicant
		if len(applicant.FirstName) > 0 {
			applicants = append(applicants, applicant)
		}
	}
	if len(applicants) > 0 {
		primaryCount := slice.CountBy(applicants, func(index int, item *borrower.ApplicantEntry) bool {
			return item.IsPrimary
		})
		if primaryCount == 0 {
			applicants[0].IsPrimary = true
		}
		if primaryCount > 1 {
			return nil, errors.WithStack(errors.New("Multiple Primary Applicant"))
		}
	}
	rpcReq.Data.Applicants = applicants

	// get borrower id
	rpcReq.Data.BorrowerId, err = utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcReq.Data.LenderCompanyId = 1

	rpcResp, err = l.svcCtx.LoanRpcClient.ApplyLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}

	return
}
