package drivelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncTransactionStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewGetSyncTransactionStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncTransactionStatusLogic {
	return &GetSyncTransactionStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *GetSyncTransactionStatusLogic) GetSyncTransactionStatus(in *borrower.DriveSyncTransactionRequest) (resp *borrower.DriveSyncTransactionResponse, err error) {
	defer func() {
		if err != nil {
			l.svcCtx.Logger.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
		}
	}()
	claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if in.Data == nil {
		return nil, errorx.New("request data is empty")
	}

	job, err := l.baseLogic.DriveLogic.GetSyncJobStatus(l.ctx, in.Data.TransactionId, claim.OrganizationId, constants.UserTypeBorrower)
	if err != nil {
		return nil, err
	}
	// no job yet: empty status means "never synced".
	if job == nil {
		return &borrower.DriveSyncTransactionResponse{Data: &borrower.DriveSyncTransactionResponseData{}}, nil
	}

	data, err := copierx.Trans[model.DriveSyncJob, borrower.DriveSyncTransactionResponseData](job)
	if err != nil {
		return nil, err
	}
	// JobId 与源字段名（job.Id）不一致，copier 不会自动映射，单独赋值。
	data.JobId = job.Id

	return &borrower.DriveSyncTransactionResponse{Data: data}, nil
}
