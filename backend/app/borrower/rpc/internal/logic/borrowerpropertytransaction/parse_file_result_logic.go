package borrowerpropertytransactionlogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"

	"zeebee-backend/common/errorx"
)

type ParseFileResultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewParseFileResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseFileResultLogic {
	return &ParseFileResultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ParseFileResult resolves the staged progress either by requestId
// (mode a) or by transaction+document+parser params (mode b: rebuilds the
// ruleSig → reqmap → requestId). Returns the current stage + processed/total;
// there is no overall percent. Neither progress nor done present → not found.
func (l *ParseFileResultLogic) ParseFileResult(in *borrower.ParseFileResultRequest) (resp *borrower.ParseFileResultResponse, err error) {
	if in.Data == nil {
		return nil, errorx.New("data required")
	}

	requestId := in.Data.RequestId
	if in.Data.Type == "params" {
		if in.Data.Params == nil {
			return nil, errorx.New("params required")
		}
		sig := parse_file.RuleSig(in.Data.Params.DocumentId, constants.DocumentTemplateParser(in.Data.Params.Parser))
		requestId, err = l.svcCtx.Progress.GetReqMap(l.ctx, in.Data.Params.TransactionId, sig)
		if err != nil {
			return nil, errorx.Wrap(err, "lookup parse request")
		}
	}
	if requestId == "" {
		return nil, errorx.New("parse request not found")
	}
	return l.readProgress(requestId)
}

// readProgress prefers the live progress record; falls back to the terminal
// done record once the run key has expired.
func (l *ParseFileResultLogic) readProgress(requestId string) (*borrower.ParseFileResultResponse, error) {
	if prog, _ := l.svcCtx.Progress.GetProgress(l.ctx, requestId); prog != nil {
		return progressToResp(prog), nil
	}
	done, err := l.svcCtx.Progress.GetDone(l.ctx, requestId)
	if err != nil {
		return nil, errorx.Wrap(err, "load parse done")
	}
	if done == nil {
		return nil, errorx.New("parse request not found")
	}
	return &borrower.ParseFileResultResponse{
		Data: &borrower.ParseFileResultResponseData{
			Status: done.Status,
			Error:  done.Error,
		},
	}, nil
}

// progressToResp projects the live staged progress: current stage + its
// processed/total. If all stages are terminal, report Finished.
func progressToResp(p *progress.Progress) *borrower.ParseFileResultResponse {
	resp := &borrower.ParseFileResultResponse{
		Data: &borrower.ParseFileResultResponseData{
			Status: p.Status,
			Error:  p.Error,
			Stage:  p.Stage,
		},
	}
	for _, s := range p.Stages {
		if s.Name == p.Stage {
			resp.Data.StageProcessed = int32(s.Processed)
			resp.Data.StageTotal = int32(s.Total)
			break
		}
	}
	return resp
}
