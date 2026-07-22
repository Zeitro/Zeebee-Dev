package borrowerpropertytransactionlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"

	"zeebee-backend/common/errorx"
)

type ParseFileExtractLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewParseFileExtractLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseFileExtractLogic {
	return &ParseFileExtractLogic{
		ctx:       ctx,
		svcCtx:    svcCtx,
		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

// ParseFileExtract triggers a manual parse (manual_tree / manual_select). IDOR
// + rule resolution live in parse_file.Dispatch; here we only resolve the
// caller's org from the claim. A duplicate trigger surfaces as gRPC
// AlreadyExists (mapped to HTTP 409 at the API layer).
func (l *ParseFileExtractLogic) ParseFileExtract(in *borrower.ParseFileExtractReq) (resp *borrower.ParseFileExtractResp, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	requestId, err := l.svcCtx.FileParser.Dispatch(l.ctx, parse_file.DispatchParam{
		TransactionId: in.TransactionId,
		DocumentId:    in.DocumentId,
		OrgId:         claim.OrganizationId,
		Parser:        in.Parser,
		Trigger:       parse_file.TriggerManual,
	})
	if err != nil {
		return nil, err
	}

	return &borrower.ParseFileExtractResp{RequestId: requestId}, nil
}
