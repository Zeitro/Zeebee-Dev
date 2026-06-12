package common

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/data/dto"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// sync resource
func NewSyncResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncResourceLogic {
	return &SyncResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncResourceLogic) SyncResource(req *types.EmptyRequest) (resp *types.EmptyResponse, err error) {
	var resourceList []dto.ResourceSyncItemDto
	err = json.Unmarshal(resourceBytes, &resourceList)
	if err != nil {
		return nil, errorx.Wrap(err, "failed to unmarshal resources")
	}
	if len(resourceList) == 0 {
		return nil, errorx.New("no resources found")
	}
	resourceContent, err := json.Marshal(resourceList)
	if err != nil {
		return nil, errorx.Wrap(err, "failed to marshal resource content")
	}

	_, err = l.svcCtx.ResourceRpcClient.SyncResource(l.ctx, &borrower.SyncResourceRequest{
		Data: &borrower.SyncResourceRequestData{
			ResourceJson: string(resourceContent),
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

var resourceBytes []byte = []byte(`[
	{
		"id": 141,
		"parent_id": 0,
		"type": "Borrower",
		"resource_type": "Menu",
		"name": "My transactions",
		"path": "/myTransactions",
		"icon": "transaction",
		"tag": "Transactions",
		"is_show": true,
		"status": "Active",
		"sort_order": 3,
		"enable_data_scope": false,
		"description": "",
		"is_public": false,
		"created_at": "2026-05-26T02:51:35Z",
		"updated_at": "2026-05-29T05:55:58Z",
		"apis": [
			{
				"id": 402,
				"parent_id": 394,
				"name": "page transactions",
				"path": "/api/v1/borrower_property_transaction/page",
				"method": "POST",
				"status": "Active",
				"description": "page transactions",
				"type": "Borrower",
				"sort_order": 8,
				"auth_required": true,
				"permission_required": true,
				"created_at": "2026-05-28T03:18:25Z",
				"updated_at": "2026-05-28T03:18:25Z"
			}
		],
		"children": [
			{
				"id": 143,
				"parent_id": 141,
				"type": "Borrower",
				"resource_type": "Button",
				"name": "Delete",
				"path": "",
				"icon": "transaction",
				"tag": "Transactions:Delete",
				"is_show": false,
				"status": "Active",
				"sort_order": 2,
				"enable_data_scope": false,
				"description": "",
				"is_public": false,
				"created_at": "2026-05-29T05:31:42Z",
				"updated_at": "2026-05-29T05:31:42Z",
				"apis": [
					{
						"id": 398,
						"parent_id": 394,
						"name": "delete transaction",
						"path": "/api/v1/borrower_property_transaction/:id",
						"method": "DELETE",
						"status": "Active",
						"description": "delete transaction",
						"type": "Borrower",
						"sort_order": 4,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					}
				]
			},
			{
				"id": 147,
				"parent_id": 141,
				"type": "Borrower",
				"resource_type": "Button",
				"name": "Add",
				"path": "",
				"icon": "",
				"tag": "Transactions:Add",
				"is_show": false,
				"status": "Active",
				"sort_order": 2,
				"enable_data_scope": false,
				"description": "",
				"is_public": false,
				"created_at": "2026-05-29T05:33:58Z",
				"updated_at": "2026-05-29T05:33:58Z",
				"apis": [
					{
						"id": 305,
						"parent_id": 303,
						"name": "all",
						"path": "/api/v1/property/all",
						"method": "GET",
						"status": "Active",
						"description": "all",
						"type": "Borrower",
						"sort_order": 2,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-03-11T07:44:46Z",
						"updated_at": "2026-03-11T07:44:46Z"
					},
					{
						"id": 395,
						"parent_id": 394,
						"name": "create transaction",
						"path": "/api/v1/borrower_property_transaction",
						"method": "POST",
						"status": "Active",
						"description": "create transaction",
						"type": "Borrower",
						"sort_order": 1,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					}
				]
			},
			{
				"id": 148,
				"parent_id": 141,
				"type": "Borrower",
				"resource_type": "Button",
				"name": "Edit",
				"path": "",
				"icon": "",
				"tag": "Transactions:Edit",
				"is_show": false,
				"status": "Active",
				"sort_order": 3,
				"enable_data_scope": false,
				"description": "",
				"is_public": false,
				"created_at": "2026-05-29T05:36:32Z",
				"updated_at": "2026-05-29T05:36:32Z",
				"apis": [
					{
						"id": 305,
						"parent_id": 303,
						"name": "all",
						"path": "/api/v1/property/all",
						"method": "GET",
						"status": "Active",
						"description": "all",
						"type": "Borrower",
						"sort_order": 2,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-03-11T07:44:46Z",
						"updated_at": "2026-03-11T07:44:46Z"
					},
					{
						"id": 396,
						"parent_id": 394,
						"name": "transaction detail",
						"path": "/api/v1/borrower_property_transaction/:id",
						"method": "GET",
						"status": "Active",
						"description": "transaction detail",
						"type": "Borrower",
						"sort_order": 2,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					},
					{
						"id": 397,
						"parent_id": 394,
						"name": "update transaction",
						"path": "/api/v1/borrower_property_transaction/:id",
						"method": "PUT",
						"status": "Active",
						"description": "update transaction",
						"type": "Borrower",
						"sort_order": 3,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					}
				]
			},
			{
				"id": 142,
				"parent_id": 141,
				"type": "Borrower",
				"resource_type": "Menu",
				"name": "Details",
				"path": "/myTransactions/:id",
				"icon": "",
				"tag": "Transactions:Details",
				"is_show": false,
				"status": "Active",
				"sort_order": 4,
				"enable_data_scope": false,
				"description": "",
				"is_public": false,
				"created_at": "2026-05-26T03:19:52Z",
				"updated_at": "2026-05-29T05:53:16Z",
				"apis": [
					{
						"id": 291,
						"parent_id": 286,
						"name": "update pending status",
						"path": "/api/v1/loan/:id/pending_status",
						"method": "PATCH",
						"status": "Active",
						"description": "update pending status",
						"type": "Borrower",
						"sort_order": 5,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-03-11T07:44:46Z",
						"updated_at": "2026-03-11T07:44:46Z"
					},
					{
						"id": 292,
						"parent_id": 286,
						"name": "delete",
						"path": "/api/v1/loan/:id",
						"method": "DELETE",
						"status": "Active",
						"description": "delete",
						"type": "Borrower",
						"sort_order": 6,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-03-11T07:44:46Z",
						"updated_at": "2026-03-11T07:44:46Z"
					},
					{
						"id": 294,
						"parent_id": 286,
						"name": "detail",
						"path": "/api/v1/loan/:id",
						"method": "GET",
						"status": "Active",
						"description": "detail",
						"type": "Borrower",
						"sort_order": 8,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-03-11T07:44:46Z",
						"updated_at": "2026-03-11T07:44:46Z"
					},
					{
						"id": 295,
						"parent_id": 286,
						"name": "page",
						"path": "/api/v1/loan/page",
						"method": "POST",
						"status": "Active",
						"description": "page",
						"type": "Borrower",
						"sort_order": 9,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-03-11T07:44:46Z",
						"updated_at": "2026-03-11T07:44:46Z"
					},
					{
						"id": 332,
						"parent_id": 331,
						"name": "create aws s3 presign url",
						"path": "/api/v1/storage/presign_url",
						"method": "POST",
						"status": "Active",
						"description": "create aws s3 presign url",
						"type": "Borrower",
						"sort_order": 1,
						"auth_required": true,
						"permission_required": false,
						"created_at": "2026-03-11T07:44:46Z",
						"updated_at": "2026-03-11T07:44:46Z"
					},
					{
						"id": 390,
						"parent_id": 286,
						"name": "update stage",
						"path": "/api/v1/loan/:id/stage",
						"method": "PATCH",
						"status": "Active",
						"description": "update third party loan's stage",
						"type": "Borrower",
						"sort_order": 7,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-22T11:17:55Z",
						"updated_at": "2026-05-22T11:17:55Z"
					},
					{
						"id": 392,
						"parent_id": 391,
						"name": "document files",
						"path": "/api/v1/document/:id/files",
						"method": "GET",
						"status": "Active",
						"description": "document files",
						"type": "Borrower",
						"sort_order": 1,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					},
					{
						"id": 396,
						"parent_id": 394,
						"name": "transaction detail",
						"path": "/api/v1/borrower_property_transaction/:id",
						"method": "GET",
						"status": "Active",
						"description": "transaction detail",
						"type": "Borrower",
						"sort_order": 2,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					},
					{
						"id": 399,
						"parent_id": 394,
						"name": "switch transaction status",
						"path": "/api/v1/borrower_property_transaction/:id/status",
						"method": "PATCH",
						"status": "Active",
						"description": "switch transaction status",
						"type": "Borrower",
						"sort_order": 5,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					},
					{
						"id": 403,
						"parent_id": 394,
						"name": "upload file to document",
						"path": "/api/v1/borrower_property_transaction/:transaction_id/document/:document_id/file",
						"method": "POST",
						"status": "Active",
						"description": "upload file to document",
						"type": "Borrower",
						"sort_order": 9,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					},
					{
						"id": 404,
						"parent_id": 394,
						"name": "delete document file",
						"path": "/api/v1/borrower_property_transaction/:transaction_id/file/:file_id",
						"method": "DELETE",
						"status": "Active",
						"description": "delete document file",
						"type": "Borrower",
						"sort_order": 10,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					},
					{
						"id": 405,
						"parent_id": 394,
						"name": "create leaf document",
						"path": "/api/v1/borrower_property_transaction/:transaction_id/document",
						"method": "POST",
						"status": "Active",
						"description": "create leaf document",
						"type": "Borrower",
						"sort_order": 11,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					},
					{
						"id": 406,
						"parent_id": 394,
						"name": "transaction operation stats",
						"path": "/api/v1/borrower_property_transaction/:transaction_id/operations/stats",
						"method": "GET",
						"status": "Active",
						"description": "transaction operation stats",
						"type": "Borrower",
						"sort_order": 12,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					},
					{
						"id": 407,
						"parent_id": 394,
						"name": "transaction operation list",
						"path": "/api/v1/borrower_property_transaction/:transaction_id/operations",
						"method": "GET",
						"status": "Active",
						"description": "transaction operation list",
						"type": "Borrower",
						"sort_order": 13,
						"auth_required": true,
						"permission_required": true,
						"created_at": "2026-05-28T03:18:25Z",
						"updated_at": "2026-05-28T03:18:25Z"
					}
				]
			}
		]
	}
]`)
