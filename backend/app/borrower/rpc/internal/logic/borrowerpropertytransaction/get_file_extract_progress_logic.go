package borrowerpropertytransactionlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file"
	"zeebee-backend/app/borrower/rpc/internal/biz_helper/parse_file/progress"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"

	"zeebee-backend/common/errorx"
)

type GetFileExtractProgressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFileExtractProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileExtractProgressLogic {
	return &GetFileExtractProgressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetFileExtractProgress resolves the staged progress either by requestId
// (mode a) or by transaction+document+parser params (mode b: rebuilds the
// ruleSig → reqmap → requestId). Returns the current stage + processed/total;
// there is no overall percent. Neither progress nor done present → not found.
func (l *GetFileExtractProgressLogic) GetFileExtractProgress(in *borrower.GetFileExtractProgressReq) (resp *borrower.GetFileExtractProgressResp, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(
				logx.Field("stacks", errorx.GetErrorStack(err)),
			).Error(err)
		}
	}()

	requestId := in.RequestId
	if in.Type == "params" {
		if in.Params == nil {
			return nil, errorx.New("params required for type=params")
		}
		sig := parse_file.RuleSig(in.Params.DocumentId, constants.DocumentTemplateParser(in.Params.Parser))
		requestId, err = l.svcCtx.Progress.GetReqMap(l.ctx, in.Params.TransactionId, sig)
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
func (l *GetFileExtractProgressLogic) readProgress(requestId string) (*borrower.GetFileExtractProgressResp, error) {
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
	return &borrower.GetFileExtractProgressResp{
		Status: done.Status,
		Error:  done.Error,
	}, nil
}

// progressToResp projects the live staged progress: current stage + its
// processed/total. If all stages are terminal, report Finished.
func progressToResp(p *progress.Progress) *borrower.GetFileExtractProgressResp {
	resp := &borrower.GetFileExtractProgressResp{Status: p.Status, Stage: p.Stage, Error: p.Error}
	for _, s := range p.Stages {
		if s.Name == p.Stage {
			resp.StageProcessed = int32(s.Processed)
			resp.StageTotal = int32(s.Total)
			break
		}
	}
	return resp
}
