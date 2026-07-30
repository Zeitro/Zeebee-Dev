package drivelogic

import (
	"context"

	"zeebee-backend/app/borrower/rpc/internal/global"
	"zeebee-backend/app/borrower/rpc/internal/logic"
	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/common/errorx"
)

type SyncTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext

	baseLogic *logic.BaseLogic
}

func NewSyncTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncTransactionLogic {
	return &SyncTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,

		baseLogic: logic.GetBaseLogic(svcCtx),
	}
}

func (l *SyncTransactionLogic) SyncTransaction(in *borrower.DriveSyncTransactionRequest) (resp *borrower.DriveSyncTransactionResponse, err error) {
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
		return nil, errorx.New("Request data is empty")
	}

	if global.DriveOAuthConfig == nil {
		return nil, errorx.New("Google Drive OAuth is not configured")
	}
	s3Client := global.GetS3Client()
	if s3Client == nil {
		return nil, errorx.New("S3 client is not configured")
	}

	input := dto.DriveSyncTransactionRequestDto{
		TransactionId:    in.Data.TransactionId,
		OrganizationId:   claim.OrganizationId,
		OrganizationType: constants.UserTypeBorrower,
		Oauth2Config:     global.DriveOAuthConfig,
		OperatorId:       claim.OrganizationUserId,
		OperatorName:     claim.Username,
	}

	// Kick off the sync asynchronously and return immediately: the actual Drive
	// sync may take minutes and must not be bound to the HTTP request.
	job, err := l.baseLogic.DriveLogic.StartSyncTransaction(l.ctx, s3Client, input)
	if err != nil {
		return nil, err
	}

	return &borrower.DriveSyncTransactionResponse{
		Data: &borrower.DriveSyncTransactionResponseData{
			Status: job.Status,
			JobId:  job.Id,
		},
	}, nil
}
