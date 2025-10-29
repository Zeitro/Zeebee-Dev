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

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update loan
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.LoanUpdateRequest) (resp *types.EmptyResponse, err error) {
	rpcReq := &borrower.LoanUpdateRequest{
		Data: &borrower.LoanUpdateRequestData{},
	}
	resp = &types.EmptyResponse{}

	err = copier.Copy(&rpcReq.Data, req)
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}

	// propertyBaseId := req.LoanProperty.PropertyBaseId
	// _ = propertyBaseId
	// rpcReq.Data.LoanProperty.PropertyBase.Address.State = req.LoanProperty.PropertyBase.Address.State

	applicants := make([]*borrower.BorrowerEditApplicantEntry, 0, len(rpcReq.Data.Applicants))
	for _, applicant := range rpcReq.Data.Applicants {
		applicant := applicant
		if len(applicant.FirstName) > 0 && len(applicant.LastName) > 0 {
			applicants = append(applicants, applicant)
		}
	}
	if len(applicants) > 0 {
		primaryCount := slice.CountBy(applicants, func(index int, item *borrower.BorrowerEditApplicantEntry) bool {
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

	_, err = l.svcCtx.LoanRpcClient.UpdateBorrowerLoan(l.ctx, rpcReq)

	return nil, err
}
