package loan

import (
	"context"

	"golang.org/x/sync/errgroup"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/constants"

	"github.com/zeromicro/go-zero/core/logx"
)

type DashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// dashboard
func NewDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DashboardLogic {
	return &DashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DashboardLogic) Dashboard(req *types.LoanDashboardRequest) (resp *types.LoanDashboardResponse, err error) {
	var pendingSummary *borrower.LoanSummaryResponseData
	var activeSummary *borrower.LoanSummaryResponseData
	var archivedPage *borrower.LoanPageResponseData
	eg := errgroup.Group{}
	eg.Go(func() error {
		// pending  loan summary
		var egErr error
		rpcReq := &borrower.LoanSummaryRequest{
			Data: &borrower.LoanSummaryRequestData{
				StartDate: req.StartDate,
				EndDate:   req.EndDate,
				Source:    req.Source,
				DueInDays: req.DueInDays,
				Sorter: &borrower.LoanPageRequestSorter{
					Column: "closing_date",
					Sort:   "asc",
				},
				Stage:    string(constants.LoanStagePending),
				DateType: "closing_date",
			},
		}

		rpcResp, egErr := l.svcCtx.LoanRpcClient.Summary(l.ctx, rpcReq)
		if egErr != nil {
			return egErr
		}
		pendingSummary = rpcResp.Data
		return egErr
	})
	eg.Go(func() error {
		// active loan summary
		var egErr error
		rpcReq := &borrower.LoanSummaryRequest{
			Data: &borrower.LoanSummaryRequestData{
				StartDate: req.StartDate,
				EndDate:   req.EndDate,
				Source:    req.Source,
				DueInDays: req.DueInDays,
				Sorter: &borrower.LoanPageRequestSorter{
					Column: "updated_at",
					Sort:   "asc",
				},
				Stage:    string(constants.LoanStageActive),
				DateType: "update",
			},
		}

		rpcResp, egErr := l.svcCtx.LoanRpcClient.Summary(l.ctx, rpcReq)
		if egErr != nil {
			return egErr
		}
		activeSummary = rpcResp.Data
		return egErr
	})
	eg.Go(func() error {
		// archived loan page
		var egErr error
		rpcReq := &borrower.LoanPageRequest{
			Data: &borrower.LoanPageRequestData{
				Page:  1,
				Limit: 1,
				Filter: &borrower.LoanPageRequestFilter{
					Stage:     string(constants.LoanStageArchived),
					Source:    req.Source,
					DateType:  "update",
					DueInDays: req.DueInDays,
					StartDate: req.StartDate,
					EndDate:   req.EndDate,
				},
				Sorter: &borrower.LoanPageRequestSorter{
					Column: "updated_at",
					Sort:   "asc",
				},
			},
		}
		rpcResp, egErr := l.svcCtx.LoanRpcClient.BorrowerLoanPage(l.ctx, rpcReq)
		if egErr != nil {
			return egErr
		}
		archivedPage = rpcResp.Data
		return egErr
	})

	err = eg.Wait()
	if err != nil {
		return
	}
	resp = &types.LoanDashboardResponse{
		TotalLoanCount:     pendingSummary.TotalLoanCount + activeSummary.TotalLoanCount + archivedPage.Total,
		PendingLoanCount:   pendingSummary.TotalLoanCount,
		ActiveLoanCount:    activeSummary.TotalLoanCount,
		ArchivedLoanCount:  archivedPage.Total,
		PendingLoanSummary: &types.LoanDashboardLoanSummary{},
		ActiveLoanSummary:  &types.LoanDashboardLoanSummary{},
	}
	resp.PendingLoanSummary, err = copierx.Trans[borrower.LoanSummaryResponseData, types.LoanDashboardLoanSummary](pendingSummary)
	if err != nil {
		return
	}
	resp.ActiveLoanSummary, err = copierx.Trans[borrower.LoanSummaryResponseData, types.LoanDashboardLoanSummary](activeSummary)
	if err != nil {
		return
	}

	return
}
