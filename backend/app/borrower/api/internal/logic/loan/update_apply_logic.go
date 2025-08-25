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

type UpdateApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// update apply
func NewUpdateApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApplyLogic {
	return &UpdateApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateApplyLogic) UpdateApply(req *types.LoanApplyRequest) (resp *types.EmptyResponse, err error) {
	logx.Info("update apply")
	rpcReq := &borrower.LoanApplyRequest{
		Data: &borrower.LoanApplyRequestData{},
	}

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
			return nil, errors.WithStack(errors.New("multiple primary applicant"))
		}
	}
	rpcReq.Data.Applicants = applicants

	// get borrower id
	rpcReq.Data.BorrowerId, err = utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcReq.Data.LenderId = 1

	_, err = l.svcCtx.LoanRpcClient.UpdateApplyLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return
}
