package logic

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/oauth2"
	"gorm.io/datatypes"

	"zeebee-backend/common/aws_s3"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/drive/client"
	"zeebee-backend/domains/zeebee/model"
)

const (
	driveDefaultFolderName = "Zeebee"

	driveFileStatusCreated  = "created"
	driveFileStatusReplaced = "replaced"
	driveFileStatusFailed   = "failed"
	driveFileStatusRemoved  = "removed"

	driveDefaultContentType = "application/octet-stream"

	// async sync job statuses (mirror drive_sync_jobs.status comment)
	DriveSyncStatusRunning       = "Running"
	DriveSyncStatusSuccess       = "Success"
	DriveSyncStatusPartialFailed = "PartialFailed"
	DriveSyncStatusFailed        = "Failed"

	// upper bound for a single background sync run, independent of the HTTP request.
	driveSyncJobTimeout = 10 * time.Minute

	// 同步操作日志的结果（写入 operations.metadata.outcome）
	driveSyncOutcomeSuccess = "success"
	driveSyncOutcomeFailed  = "failed"
)

// DriveLogic provides Google Drive integration business logic (per organization).
type DriveLogic struct {
	*BaseLogic
}

// NewDriveLogic creates a new drive logic instance.
func NewDriveLogic(baseLogic *BaseLogic) *DriveLogic {
	return &DriveLogic{
		BaseLogic: baseLogic,
	}
}

// GetIntegration returns the connection state for an organization. A missing row
// is reported as "not connected" with sensible defaults.
func (l *DriveLogic) GetIntegration(ctx context.Context, organizationId int64, organizationType constants.UserType) (result *dto.DriveIntegrationDto, err error) {
	row, err := l.repos.DriveAuthorizationRepo.GetByOrganization(ctx, organizationId, organizationType)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return &dto.DriveIntegrationDto{
			IsAuthorized: false,
			FolderName:   driveDefaultFolderName,
			DriveName:    "My Drive",
		}, nil
	}
	return &dto.DriveIntegrationDto{
		IsAuthorized:       row.IsAuthorized,
		AuthEmail:          row.AuthEmail,
		FolderName:         row.FolderName,
		DriveId:            row.DriveId,
		DriveName:          row.DriveName,
		RootFolderUrl:      row.RootFolderUrl,
		AuthorizedByUserId: row.AuthorizedByUserId,
	}, nil
}

// GenerateAuthUrl ensures a row exists then builds a consent URL.
func (l *DriveLogic) GenerateAuthUrl(ctx context.Context, input dto.DriveAuthUrlRequestDto) (url string, err error) {
	if input.AuthCodeGenerator == nil {
		return "", errorx.New("Drive auth code url generator is nil")
	}
	row, err := l.repos.DriveAuthorizationRepo.GetByOrganization(ctx, input.OrganizationId, input.OrganizationType)
	if err != nil {
		return "", err
	}
	if row == nil {
		err = l.repos.DriveAuthorizationRepo.Upsert(ctx, &model.DriveAuthorization{
			OrganizationId:   input.OrganizationId,
			OrganizationType: input.OrganizationType,
		})
		if err != nil {
			return "", err
		}
	}
	// state is not security-critical: /drive/authorize is authenticated and re-derives
	// the organization from the JWT claim. A stable marker is enough.
	return input.AuthCodeGenerator.GenerateAuthCodeUrl("drive"), nil
}

// Authorize exchanges the OAuth code, verifies the account, and persists the token.
func (l *DriveLogic) Authorize(ctx context.Context, input dto.DriveAuthorizeRequestDto) (result *dto.DriveIntegrationDto, err error) {
	if input.Oauth2Config == nil {
		return nil, errorx.New("Drive oauth config is nil")
	}
	if input.Code == "" {
		return nil, errorx.New("Code is empty")
	}
	cli, err := client.NewDriveClient(input.Oauth2Config, input.Code)
	if err != nil {
		return nil, err
	}
	email, err := cli.GetEmail(ctx)
	if err != nil {
		return nil, err
	}
	token := cli.Token()
	if token == nil {
		return nil, errorx.New("Drive token is nil")
	}
	err = l.repos.DriveAuthorizationRepo.Upsert(ctx, &model.DriveAuthorization{
		OrganizationId:     input.OrganizationId,
		OrganizationType:   input.OrganizationType,
		AuthEmail:          email,
		IsAuthorized:       true,
		Token:              &model.DriveToken{Token: token},
		AuthorizedByUserId: input.UserId,
	})
	if err != nil {
		return nil, err
	}
	return l.GetIntegration(ctx, input.OrganizationId, input.OrganizationType)
}

// ListDrives lists the shared drives accessible to the connected account, with My
// Drive prepended as the default option.
func (l *DriveLogic) ListDrives(ctx context.Context, input dto.DriveListDrivesRequestDto) (result []client.DriveInfo, err error) {
	cli, _, err := l.buildClient(ctx, input.OrganizationId, input.OrganizationType, input.Oauth2Config)
	if err != nil {
		return nil, err
	}
	drives, err := cli.ListDrives(ctx)
	if err != nil {
		return nil, err
	}
	result = make([]client.DriveInfo, 0, len(drives)+1)
	result = append(result, client.DriveInfo{Id: "", Name: "My Drive"})
	result = append(result, drives...)
	return result, nil
}

// UpdateSelection switches the target drive.
func (l *DriveLogic) UpdateSelection(ctx context.Context, input dto.DriveUpdateSelectionRequestDto) (err error) {
	driveName := input.DriveName
	if input.DriveId == "" {
		driveName = "My Drive"
	}
	return l.repos.DriveAuthorizationRepo.UpdateDriveSelection(ctx, input.OrganizationId, input.OrganizationType, input.DriveId, driveName)
}

// UpdateFolder renames the root folder.
func (l *DriveLogic) UpdateFolder(ctx context.Context, input dto.DriveUpdateFolderRequestDto) (err error) {
	folderName := strings.TrimSpace(input.FolderName)
	if folderName == "" {
		return errorx.New("Folder name is empty")
	}
	return l.repos.DriveAuthorizationRepo.UpdateFolderName(ctx, input.OrganizationId, input.OrganizationType, folderName)
}

// Disconnect revokes the local connection (token/email cleared, settings kept).
func (l *DriveLogic) Disconnect(ctx context.Context, organizationId int64, organizationType constants.UserType) (err error) {
	return l.repos.DriveAuthorizationRepo.Disconnect(ctx, organizationId, organizationType)
}

// SyncTransaction mirrors one transaction's checklist documents into Drive following
// the directory layout: <FolderName>/<property address>/<transaction folder>/
// (Required|Optional)/<nested node.Name...>/<files>.
// StartSyncTransaction validates connectivity cheaply, records a Running job and
// launches the actual Drive sync in a detached goroutine, returning the job
// immediately. The sync itself may take minutes (many sequential Drive round
// trips) and must not be bound to the HTTP request, which the dev proxy aborts
// after ~30s. The frontend polls GetSyncJobStatus for progress/result.
func (l *DriveLogic) StartSyncTransaction(ctx context.Context, s3Client *aws_s3.S3Client, input dto.DriveSyncTransactionRequestDto) (job *model.DriveSyncJob, err error) {
	if s3Client == nil {
		return nil, errorx.New("S3 client is nil")
	}
	// fast, DB-only pre-check so "not connected" fails synchronously.
	row, err := l.repos.DriveAuthorizationRepo.GetByOrganization(ctx, input.OrganizationId, input.OrganizationType)
	if err != nil {
		return nil, err
	}
	if row == nil || !row.IsAuthorized || row.Token == nil || row.Token.Token == nil {
		return nil, errorx.New("Google Drive is not connected")
	}
	// ensure the transaction exists and belongs to the org before spawning work.
	if _, err = l.repos.BorrowerPropertyTransactionRepo.GetOrganizationTransaction(ctx, input.TransactionId, input.OrganizationId); err != nil {
		return nil, errorx.Wrap(err, "transaction not found")
	}

	// 防止重复同步：若已有一个尚未超时的 Running 任务，直接复用它（接口幂等），
	// 避免两个 goroutine 并发 EnsureFolder 竞态，在 Drive 中创建出重复文件夹。
	latest, err := l.repos.DriveSyncJobRepo.GetLatestByTransaction(ctx, input.TransactionId, input.OrganizationId, input.OrganizationType)
	if err != nil {
		return nil, err
	}
	if latest != nil && latest.Status == DriveSyncStatusRunning {
		if time.Since(latest.CreatedAt) < driveSyncJobTimeout {
			return latest, nil
		}
		// 超时仍停留在 Running（多为进程中断遗留的僵尸任务）：置为终态。
		// 既清理僵尸，也释放下方部分唯一索引占用的「每交易一个 Running」名额。
		if uerr := l.repos.DriveSyncJobRepo.UpdateResult(ctx, latest.Id, map[string]interface{}{
			"status": DriveSyncStatusFailed,
			"error":  truncateError("sync timed out (stale job cleaned up)"),
		}); uerr != nil {
			return nil, uerr
		}
	}

	job, err = l.repos.DriveSyncJobRepo.Create(ctx, &model.DriveSyncJob{
		OrganizationId:   input.OrganizationId,
		OrganizationType: input.OrganizationType,
		TransactionId:    input.TransactionId,
		Status:           DriveSyncStatusRunning,
	})
	if err != nil {
		// 两个请求几乎同时创建时，后一个会撞上「每交易至多一个 Running」的唯一索引冲突。
		// 此时复用已存在的 Running 任务而非报错，使接口保持幂等、避免重复同步。
		if gormx.IsUniqueViolation(err) {
			if existing, gerr := l.repos.DriveSyncJobRepo.GetLatestByTransaction(ctx, input.TransactionId, input.OrganizationId, input.OrganizationType); gerr == nil &&
				existing != nil && existing.Status == DriveSyncStatusRunning {
				return existing, nil
			}
		}
		return nil, err
	}

	jobId := job.Id
	go l.runSyncJob(jobId, s3Client, input)

	return job, nil
}

// runSyncJob executes the sync on a detached context and persists the outcome.
func (l *DriveLogic) runSyncJob(jobId int64, s3Client *aws_s3.S3Client, input dto.DriveSyncTransactionRequestDto) {
	ctx, cancel := context.WithTimeout(context.Background(), driveSyncJobTimeout)
	defer cancel()

	result, err := l.runSyncTransaction(ctx, s3Client, input)

	columns := map[string]interface{}{}
	var status string
	if err != nil {
		status = DriveSyncStatusFailed
		columns["status"] = status
		columns["error"] = truncateError(err.Error())
	} else {
		status = DriveSyncStatusSuccess
		if result.Failed > 0 {
			status = DriveSyncStatusPartialFailed
		}
		columns["status"] = status
		columns["total"] = result.Total
		columns["created"] = result.Created
		columns["replaced"] = result.Replaced
		columns["failed"] = result.Failed
		columns["removed"] = result.Removed
		columns["folder_web_link"] = result.FolderWebLink
		columns["files"] = toSyncFiles(result.Files)
	}

	if uerr := l.repos.DriveSyncJobRepo.UpdateResult(ctx, jobId, columns); uerr != nil {
		logx.WithContext(ctx).WithFields(logx.Field("job_id", jobId)).
			Errorf("drive sync job persist failed: %v", uerr)
	}

	// 记录交易操作日志：成功/失败（outcome）放 metadata，供 Log 页展示。
	l.recordSyncOperation(ctx, input, status, result, err)
}

// recordSyncOperation 在同步任务终态后写一条交易操作日志。
// 事件统一为 "Synced to Google Drive"，成功/失败及统计明细放在 metadata 里：
// 整体出错或存在文件失败（部分失败）→ outcome=failed；全部成功 → outcome=success。
func (l *DriveLogic) recordSyncOperation(ctx context.Context, input dto.DriveSyncTransactionRequestDto, status string, result *dto.DriveSyncTransactionResultDto, syncErr error) {
	meta := map[string]interface{}{"status": status}
	outcome := driveSyncOutcomeSuccess
	if syncErr != nil {
		outcome = driveSyncOutcomeFailed
		meta["error"] = truncateError(syncErr.Error())
	} else if result != nil {
		meta["total"] = result.Total
		meta["created"] = result.Created
		meta["replaced"] = result.Replaced
		meta["removed"] = result.Removed
		meta["failed"] = result.Failed
		// 部分失败也算失败：只要有文件没同步成功，outcome 即为 failed。
		if result.Failed > 0 {
			outcome = driveSyncOutcomeFailed
		}
	}
	meta["outcome"] = outcome

	metaJSON, merr := json.Marshal(meta)
	if merr != nil {
		metaJSON = []byte("{}")
	}

	op := &model.Operation{
		EntityId:       input.TransactionId,
		EntityType:     constants.OperationEntityTypeBorrowerPropertyTransaction,
		OperatorId:     input.OperatorId,
		OperatorName:   input.OperatorName,
		OperatorType:   input.OrganizationType,
		Operation:      constants.OperationActionBorrowerPropertyTransactionSyncToDrive,
		Metadata:       datatypes.JSON(metaJSON),
		OrganizationId: input.OrganizationId,
	}
	if _, oerr := l.repos.OperationRepo.Create(ctx, op); oerr != nil {
		logx.WithContext(ctx).WithFields(logx.Field("transaction_id", input.TransactionId)).
			Errorf("record drive sync operation failed: %v", oerr)
	}
}

// GetSyncJobStatus returns the latest sync job for a transaction (nil if none).
func (l *DriveLogic) GetSyncJobStatus(ctx context.Context, transactionId, organizationId int64, organizationType constants.UserType) (job *model.DriveSyncJob, err error) {
	return l.repos.DriveSyncJobRepo.GetLatestByTransaction(ctx, transactionId, organizationId, organizationType)
}

func toSyncFiles(files []dto.DriveSyncedFileDto) model.DriveSyncFiles {
	out := make(model.DriveSyncFiles, 0, len(files))
	for _, f := range files {
		out = append(out, model.DriveSyncFileItem{
			FileName: f.FileName,
			DriveUrl: f.DriveUrl,
			Status:   f.Status,
		})
	}
	return out
}

func truncateError(msg string) string {
	const max = 500
	if len(msg) > max {
		return msg[:max]
	}
	return msg
}

func (l *DriveLogic) runSyncTransaction(ctx context.Context, s3Client *aws_s3.S3Client, input dto.DriveSyncTransactionRequestDto) (result *dto.DriveSyncTransactionResultDto, err error) {
	if s3Client == nil {
		return nil, errorx.New("S3 client is nil")
	}
	cli, row, err := l.buildClient(ctx, input.OrganizationId, input.OrganizationType, input.Oauth2Config)
	if err != nil {
		return nil, err
	}

	transaction, err := l.repos.BorrowerPropertyTransactionRepo.GetOrganizationTransaction(ctx, input.TransactionId, input.OrganizationId)
	if err != nil {
		return nil, errorx.Wrap(err, "transaction not found")
	}

	// load files into checklist leaf nodes
	l.fillChecklistFiles(ctx, &transaction.Checklist)

	driveId := row.DriveId

	// 1. root folder
	folderName := row.FolderName
	if folderName == "" {
		folderName = driveDefaultFolderName
	}
	rootId, rootUrl, err := cli.EnsureFolder(ctx, sanitizeFolderName(folderName), "", driveId)
	if err != nil {
		return nil, err
	}
	if rootId != row.RootFolderId || rootUrl != row.RootFolderUrl {
		_ = l.repos.DriveAuthorizationRepo.UpdateRootFolder(ctx, input.OrganizationId, input.OrganizationType, rootId, rootUrl)
	}

	// 2. property address folder
	propId, _, err := cli.EnsureFolder(ctx, buildPropertyFolderName(transaction), rootId, driveId)
	if err != nil {
		return nil, err
	}

	// 3. transaction folder "<Type> Transaction <Code> (<Status>)"：同一交易固定复用
	//    同一个文件夹，状态变化时改名而非新建，避免遗留旧状态的孤儿文件夹。
	txId, txUrl, err := l.ensureTransactionFolder(ctx, cli, propId, driveId, transaction)
	if err != nil {
		return nil, err
	}

	result = &dto.DriveSyncTransactionResultDto{FolderWebLink: txUrl}

	// 4. 平铺 checklist 子文件夹：required/optional 文档的子节点直接挂在
	//    transaction 文件夹下，不再额外包一层 Required/Optional 目录。
	l.syncNodes(ctx, cli, s3Client, driveId, txId, transaction.Checklist.RequiredDocuments, result)
	l.syncNodes(ctx, cli, s3Client, driveId, txId, transaction.Checklist.OptionalDocuments, result)

	return result, nil
}

// buildClient loads the stored authorization, builds a client and refreshes the
// token when needed (persisting the refreshed token).
func (l *DriveLogic) buildClient(ctx context.Context, organizationId int64, organizationType constants.UserType, config *oauth2.Config) (cli *client.DriveClient, row *model.DriveAuthorization, err error) {
	if config == nil {
		return nil, nil, errorx.New("Drive oauth config is nil")
	}
	row, err = l.repos.DriveAuthorizationRepo.GetByOrganization(ctx, organizationId, organizationType)
	if err != nil {
		return nil, nil, err
	}
	if row == nil || !row.IsAuthorized || row.Token == nil || row.Token.Token == nil {
		return nil, nil, errorx.New("Google Drive is not connected")
	}
	cli, err = client.NewDriveClientFromToken(config, row.Token.Token)
	if err != nil {
		return nil, nil, err
	}
	token, err := cli.RefreshToken(ctx)
	if err != nil {
		return nil, nil, err
	}
	if token != nil && token.AccessToken != row.Token.Token.AccessToken {
		_ = l.repos.DriveAuthorizationRepo.UpdateToken(ctx, organizationId, organizationType, &model.DriveToken{Token: token})
	}
	return cli, row, nil
}

// syncNodes recursively mirrors the checklist tree: every node becomes a folder, and
// files attached to a leaf node are uploaded into that node's folder.
func (l *DriveLogic) syncNodes(ctx context.Context, cli *client.DriveClient, s3Client *aws_s3.S3Client, driveId, parentId string, nodes []*dto.ChecklistNode, result *dto.DriveSyncTransactionResultDto) {
	for i := range nodes {
		node := nodes[i]
		folderId, _, err := cli.EnsureFolder(ctx, sanitizeFolderName(node.Name), parentId, driveId)
		if err != nil {
			// folder creation failed: record every file under this node as failed
			l.markNodeFilesFailed(node, result)
			continue
		}
		if len(node.Children) > 0 {
			l.syncNodes(ctx, cli, s3Client, driveId, folderId, node.Children, result)
			continue
		}
		for _, f := range node.Files {
			l.syncFile(ctx, cli, s3Client, driveId, folderId, f, result)
		}
		// 把「我们之前同步上去、但系统侧当前已不存在」的文件从 Drive 移除，
		// 使系统侧的删除能反映到 Drive。仅清理带标记的文件，用户手动放进该文件夹的文件不受影响。
		l.reconcileFolderFiles(ctx, cli, driveId, folderId, node.Files, result)
	}
}

// reconcileFolderFiles trashes managed Drive files whose names are no longer present
// in the system's current file set for this folder.
func (l *DriveLogic) reconcileFolderFiles(ctx context.Context, cli *client.DriveClient, driveId, folderId string, currentFiles []*dto.FileEntityDto, result *dto.DriveSyncTransactionResultDto) {
	managed, err := cli.ListManagedFiles(ctx, folderId, driveId)
	if err != nil {
		logx.WithContext(ctx).WithFields(logx.Field("folder_id", folderId)).
			Errorf("drive reconcile list managed files failed: %v", err)
		return
	}
	keep := make(map[string]struct{}, len(currentFiles))
	for _, f := range currentFiles {
		keep[f.FileName] = struct{}{}
	}
	for _, df := range managed {
		if _, ok := keep[df.Name]; ok {
			continue
		}
		if terr := cli.TrashFile(ctx, df.Id); terr != nil {
			logx.WithContext(ctx).WithFields(logx.Field("drive_file_id", df.Id)).
				Errorf("drive reconcile trash file failed: %v", terr)
			continue
		}
		result.Removed++
		result.Files = append(result.Files, dto.DriveSyncedFileDto{FileName: df.Name, Status: driveFileStatusRemoved})
	}
}

// syncFile downloads a file from S3 and upserts it into the target Drive folder.
func (l *DriveLogic) syncFile(ctx context.Context, cli *client.DriveClient, s3Client *aws_s3.S3Client, driveId, parentId string, f *dto.FileEntityDto, result *dto.DriveSyncTransactionResultDto) {
	result.Total++
	synced := dto.DriveSyncedFileDto{FileName: f.FileName}

	key := s3Client.ParseKeyFromUrl(f.FileUrl)
	obj, err := s3Client.GetObject(ctx, key)
	if err != nil {
		synced.Status = driveFileStatusFailed
		result.Failed++
		result.Files = append(result.Files, synced)
		return
	}
	defer func() { _ = obj.Body.Close() }()

	_, webViewLink, replaced, err := cli.UpsertFile(ctx, f.FileName, contentTypeForName(f.FileName), obj.Body, parentId, driveId)
	if err != nil {
		synced.Status = driveFileStatusFailed
		result.Failed++
		result.Files = append(result.Files, synced)
		return
	}
	synced.DriveUrl = webViewLink
	if replaced {
		synced.Status = driveFileStatusReplaced
		result.Replaced++
	} else {
		synced.Status = driveFileStatusCreated
		result.Created++
	}
	// persist per-file Drive sync state so the frontend can show a "Google Drive"
	// tag for files currently mirrored into Drive.
	if f.Id > 0 {
		if uerr := l.repos.FileRepo.MarkDriveSynced(ctx, f.Id, webViewLink, time.Now()); uerr != nil {
			logx.WithContext(ctx).WithFields(logx.Field("file_id", f.Id)).
				Errorf("mark file drive synced failed: %v", uerr)
		}
	}
	result.Files = append(result.Files, synced)
}

func (l *DriveLogic) markNodeFilesFailed(node *dto.ChecklistNode, result *dto.DriveSyncTransactionResultDto) {
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			l.markNodeFilesFailed(child, result)
		}
		return
	}
	for _, f := range node.Files {
		result.Total++
		result.Failed++
		result.Files = append(result.Files, dto.DriveSyncedFileDto{FileName: f.FileName, Status: driveFileStatusFailed})
	}
}

// fillChecklistFiles loads files for every leaf node of the checklist tree.
func (l *DriveLogic) fillChecklistFiles(ctx context.Context, checklist *dto.ChecklistRoot) {
	l.fillNodeFiles(ctx, checklist.RequiredDocuments)
	l.fillNodeFiles(ctx, checklist.OptionalDocuments)
}

func (l *DriveLogic) fillNodeFiles(ctx context.Context, nodes []*dto.ChecklistNode) {
	for i := range nodes {
		if len(nodes[i].Children) > 0 {
			l.fillNodeFiles(ctx, nodes[i].Children)
			continue
		}
		files, err := l.repos.FileRepo.GetFiles(ctx, nodes[i].DocumentId, constants.FileEntityTypeDocument)
		if err != nil || len(files) == 0 {
			continue
		}
		dtoFiles := make([]*dto.FileEntityDto, 0, len(files))
		for _, f := range files {
			dtoFile, terr := copierx.Trans[model.File, dto.FileEntityDto](f, copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle))
			if terr != nil || dtoFile == nil {
				continue
			}
			dtoFiles = append(dtoFiles, dtoFile)
		}
		nodes[i].Files = dtoFiles
	}
}

// buildTransactionFolderName builds "<Type> Transaction (<Status>)", e.g.
// "Purchase Transaction (Pending)". The transaction code is intentionally
// omitted from the folder name.
func buildTransactionFolderName(transaction *model.BorrowerPropertyTransaction) string {
	name := transactionFolderNamePrefix(transaction) + strings.TrimSpace(string(transaction.Status)) + ")"
	return sanitizeFolderName(name)
}

// transactionFolderNamePrefix 返回交易文件夹名中不随状态变化的稳定前缀
// "<Type> Transaction <Code> ("，用于识别同一交易在状态变更前后的文件夹。
// 前缀含交易编号（Code），故能唯一对应一笔交易。
func transactionFolderNamePrefix(transaction *model.BorrowerPropertyTransaction) string {
	prefix := strings.TrimSpace(string(transaction.Type)) + " Transaction"
	if code := strings.TrimSpace(transaction.Code); code != "" {
		prefix += " " + code
	}
	return prefix + " ("
}

// ensureTransactionFolder resolves the single Drive folder for a transaction.
// 父目录下若已存在该交易的文件夹（按不含状态的前缀识别）：
//   - 名字已是当前状态：直接复用；
//   - 名字是旧状态：改名为当前状态后复用（folder id / 链接不变），不再新建，
//     从而把旧状态的文件夹「换」成当前状态，避免遗留孤儿文件夹。
//
// 找不到现成文件夹（或列举失败兜底）时再按名查找或新建。
func (l *DriveLogic) ensureTransactionFolder(ctx context.Context, cli *client.DriveClient, parentId, driveId string, transaction *model.BorrowerPropertyTransaction) (id, webViewLink string, err error) {
	targetName := buildTransactionFolderName(transaction)
	prefix := transactionFolderNamePrefix(transaction)

	folders, lerr := cli.ListFolders(ctx, parentId, driveId)
	if lerr == nil {
		if folder, needRename := pickTransactionFolder(folders, targetName, prefix); folder != nil {
			if needRename {
				if rerr := cli.RenameFile(ctx, folder.Id, targetName); rerr != nil {
					return "", "", rerr
				}
			}
			return folder.Id, folder.WebViewLink, nil
		}
	}
	return cli.EnsureFolder(ctx, targetName, parentId, driveId)
}

// pickTransactionFolder decides which existing sibling folder to use, by name only:
//   - 名字与 targetName 完全一致：直接复用（needRename=false）；
//   - 否则取第一个以 prefix 开头的旧状态文件夹：改名复用（needRename=true）；
//   - 都没有：返回 (nil, false)，交由调用方新建。
func pickTransactionFolder(folders []client.DriveFolderInfo, targetName, prefix string) (folder *client.DriveFolderInfo, needRename bool) {
	var stale *client.DriveFolderInfo
	for i := range folders {
		if folders[i].Name == targetName {
			return &folders[i], false
		}
		if stale == nil && strings.HasPrefix(folders[i].Name, prefix) {
			stale = &folders[i]
		}
	}
	return stale, stale != nil
}

// buildPropertyFolderName builds a readable folder name from the property address.
func buildPropertyFolderName(transaction *model.BorrowerPropertyTransaction) string {
	if transaction.BorrowerProperty != nil {
		if name := formatAddress(transaction.BorrowerProperty.Address); name != "" {
			return sanitizeFolderName(name)
		}
	}
	return sanitizeFolderName("Property " + transaction.Code)
}

func formatAddress(addr model.Address) string {
	parts := make([]string, 0, 4)
	for _, p := range []string{addr.Address1, addr.Address2, addr.City, addr.State, addr.ZipCode} {
		if s := strings.TrimSpace(p); s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, ", ")
}

// sanitizeFolderName removes characters that make Drive folder paths ambiguous.
func sanitizeFolderName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	if name == "" {
		return "Untitled"
	}
	return name
}

// contentTypeForName resolves a MIME type from the file extension.
func contentTypeForName(name string) string {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
	if ct, ok := aws_s3.ContentTypeMap[ext]; ok {
		return ct
	}
	return driveDefaultContentType
}
