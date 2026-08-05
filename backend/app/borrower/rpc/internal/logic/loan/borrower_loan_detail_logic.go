package loanlogic

import (
	"context"

	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

type BorrowerLoanDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewBorrowerLoanDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BorrowerLoanDetailLogic {
	return &BorrowerLoanDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *BorrowerLoanDetailLogic) BorrowerLoanDetail(in *borrower.LoanDetailRequest) (resp *borrower.LoanDetailResponse, err error) {

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}

	loan, err := l.svcCtx.LoanRepo.GetBorrowerLoan(l.ctx, in.Data.Id, claim.OrganizationId)
	if err != nil {
		return nil, err
	}
	resp = &borrower.LoanDetailResponse{
		Data: &proto.LoanEntity{},
	}
	resp.Data, err = copierx.Trans[model.Loan, proto.LoanEntity](loan,
		copierx.TimeToRFC3339Str(),
		copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil {
		return nil, err
	}

	// loan experiences
	var transExperiences []*proto.LoanExperienceEntity
	var realEstateExperiences []*proto.LoanExperienceEntity
	slice.ForEach(loan.LoanExperiences, func(index int, item *model.LoanExperience) {
		experience := &proto.LoanExperienceEntity{}
		experience, _ = copierx.Trans[model.LoanExperience, proto.LoanExperienceEntity](item,
			copierx.TimeToRFC3339Str(),
			copierx.TimePtrToRFC3339Str(),
			copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
		)
		if item.Type == constants.LoanExperienceTypeTransaction {
			transExperiences = append(transExperiences, experience)
		} else {
			realEstateExperiences = append(realEstateExperiences, experience)
		}
	})
	// requested borrower loan document count
	slice.ForEach(loan.LoanDocuments, func(index int, item *model.LoanDocument) {
		if item.Category != constants.LoanDocumentCategoryBorrower && item.Category != constants.LoanDocumentCategoryExtend {
			return
		}
		if item.Status == constants.LoanDocumentStatusRequested {
			resp.Data.RequestedLoanDocumentCount++
		}
	})

	resp.Data.TransactionExperiences = transExperiences
	resp.Data.RealEstateExperiences = realEstateExperiences

	return resp, err
}
