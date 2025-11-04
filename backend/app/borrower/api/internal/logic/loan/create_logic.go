package loan

import (
	"context"

	"github.com/jinzhu/copier"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateLogic create loan
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.LoanCreateRequest) (resp *types.LoanCreateResponse, err error) {
	rpcReq := &borrower.LoanCreateRequest{
		Data: &borrower.LoanCreateRequestData{},
	}
	resp = &types.LoanCreateResponse{}
	rpcResp := &borrower.LoanCreateResponse{}

	err = copier.CopyWithOption(&rpcReq.Data, req, copier.Option{
		DeepCopy:      true,
		IgnoreEmpty:   true,
		CaseSensitive: true,
	})
	if err != nil {
		return nil, errorx.GetFailedTransReqErr(err)
	}
	// todo handle applicants
	// applicants := make([]*borrower.BorrowerApplicantEditEntry, 0, len(rpcReq.Data.Applicants))
	// for _, applicant := range rpcReq.Data.Applicants {
	// 	applicant := applicant
	// 	if len(applicant.FirstName) > 0 && len(applicant.LastName) > 0 {
	// 		applicants = append(applicants, applicant)
	// 	}
	// }
	// if len(applicants) > 0 {
	// 	primaryCount := slice.CountBy(applicants, func(index int, item *borrower.BorrowerApplicantEditEntry) bool {
	// 		return item.IsPrimary
	// 	})
	// 	if primaryCount == 0 {
	// 		applicants[0].IsPrimary = true
	// 	}
	// 	if primaryCount > 1 {
	// 		return nil, errors.WithStack(errors.New("Multiple Primary Applicant"))
	// 	}
	// }
	// rpcReq.Data.Applicants = applicants

	// get borrower id
	rpcReq.Data.BorrowerId, err = utils.GetAccountIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, err = l.svcCtx.LoanRpcClient.CreateBorrowerLoan(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(resp, &rpcResp.Data)
	if err != nil {
		return nil, errorx.GetFailedTransRespErr(err)
	}
	return resp, err
}
