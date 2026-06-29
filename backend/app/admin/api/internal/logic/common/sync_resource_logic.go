package common

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/admin/api/internal/svc"
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
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

	_, err = l.svcCtx.ResourceRpcClient.SyncResource(l.ctx, &admin.SyncResourceRequest{
		Data: &admin.SyncResourceRequestData{
			ResourceJson: string(resourceContent),
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

var resourceBytes []byte = []byte(`[{"id":152,"parent_id":0,"type":"Admin","resource_type":"Menu","name":"Pricing","path":"pricing","icon":"pricing","tag":"Pricing","is_show":true,"status":"Active","sort_order":6,"enable_data_scope":false,"description":"","is_public":false,"created_at":"2026-06-16T03:08:20Z","updated_at":"2026-06-17T07:41:44Z","children":[{"id":153,"parent_id":152,"type":"Admin","resource_type":"Menu","name":"Pricing plan","path":"/pricing/pricingPlan","icon":"","tag":"Pricing:PricingPlan","is_show":true,"status":"Active","sort_order":1,"enable_data_scope":false,"description":"","is_public":false,"created_at":"2026-06-16T03:09:37Z","updated_at":"2026-06-18T02:28:57Z","apis":[{"id":428,"parent_id":422,"name":"page","path":"/api/v1/pricing_plan/page","method":"POST","status":"Active","description":"page","type":"Admin","sort_order":6,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"}],"children":[{"id":155,"parent_id":153,"type":"Admin","resource_type":"Button","name":"Edit","path":"","icon":"","tag":"Pricing:PricingPlan:Edit","is_show":false,"status":"Active","sort_order":1,"enable_data_scope":false,"description":"","is_public":false,"created_at":"2026-06-17T08:08:13Z","updated_at":"2026-06-18T02:35:03Z","apis":[{"id":419,"parent_id":418,"name":"create","path":"/api/v1/pricing_plan_policy","method":"POST","status":"Active","description":"create","type":"Admin","sort_order":1,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"},{"id":420,"parent_id":418,"name":"update","path":"/api/v1/pricing_plan_policy/:id","method":"PUT","status":"Active","description":"update","type":"Admin","sort_order":2,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"},{"id":421,"parent_id":418,"name":"delete","path":"/api/v1/pricing_plan_policy/:id","method":"DELETE","status":"Active","description":"delete","type":"Admin","sort_order":3,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"},{"id":425,"parent_id":422,"name":"update","path":"/api/v1/pricing_plan/:id","method":"PUT","status":"Active","description":"update","type":"Admin","sort_order":3,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"},{"id":427,"parent_id":422,"name":"detail","path":"/api/v1/pricing_plan/:id","method":"GET","status":"Active","description":"detail","type":"Admin","sort_order":5,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"}]},{"id":156,"parent_id":153,"type":"Admin","resource_type":"Button","name":"Delete","path":"","icon":"","tag":"Pricing:PricingPlan:Delete","is_show":false,"status":"Active","sort_order":2,"enable_data_scope":false,"description":"","is_public":false,"created_at":"2026-06-17T08:09:05Z","updated_at":"2026-06-18T02:35:43Z","apis":[{"id":426,"parent_id":422,"name":"delete","path":"/api/v1/pricing_plan/:id","method":"DELETE","status":"Active","description":"delete","type":"Admin","sort_order":4,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"}]},{"id":157,"parent_id":153,"type":"Admin","resource_type":"Button","name":"Change status","path":"","icon":"","tag":"Pricing:PricingPlan:ChangeStatus","is_show":false,"status":"Active","sort_order":3,"enable_data_scope":false,"description":"","is_public":false,"created_at":"2026-06-17T08:09:53Z","updated_at":"2026-06-18T02:36:09Z","apis":[{"id":429,"parent_id":422,"name":"switch status","path":"/api/v1/pricing_plan/:id/status","method":"PATCH","status":"Active","description":"switch status","type":"Admin","sort_order":7,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"}]},{"id":158,"parent_id":153,"type":"Admin","resource_type":"Button","name":"Add","path":"","icon":"","tag":"Pricing:PricingPlan:Add","is_show":false,"status":"Active","sort_order":4,"enable_data_scope":false,"description":"","is_public":false,"created_at":"2026-06-17T08:10:39Z","updated_at":"2026-06-18T02:36:45Z","apis":[{"id":419,"parent_id":418,"name":"create","path":"/api/v1/pricing_plan_policy","method":"POST","status":"Active","description":"create","type":"Admin","sort_order":1,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"},{"id":420,"parent_id":418,"name":"update","path":"/api/v1/pricing_plan_policy/:id","method":"PUT","status":"Active","description":"update","type":"Admin","sort_order":2,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"},{"id":421,"parent_id":418,"name":"delete","path":"/api/v1/pricing_plan_policy/:id","method":"DELETE","status":"Active","description":"delete","type":"Admin","sort_order":3,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"},{"id":423,"parent_id":422,"name":"create","path":"/api/v1/pricing_plan","method":"POST","status":"Active","description":"create","type":"Admin","sort_order":1,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"}]},{"id":159,"parent_id":153,"type":"Admin","resource_type":"Button","name":"Menu permission","path":"","icon":"","tag":"Pricing:PricingPlan:MenuPermission","is_show":false,"status":"Active","sort_order":5,"enable_data_scope":false,"description":"","is_public":false,"created_at":"2026-06-17T08:12:18Z","updated_at":"2026-06-18T02:38:36Z","apis":[{"id":144,"parent_id":142,"name":"tree","path":"/api/v1/resource/tree","method":"GET","status":"Active","description":"tree","type":"Admin","sort_order":2,"auth_required":true,"permission_required":true,"created_at":"2026-03-11T07:44:46Z","updated_at":"2026-03-11T07:44:46Z"},{"id":427,"parent_id":422,"name":"detail","path":"/api/v1/pricing_plan/:id","method":"GET","status":"Active","description":"detail","type":"Admin","sort_order":5,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"},{"id":430,"parent_id":422,"name":"bind resources","path":"/api/v1/pricing_plan/:id/bind_resources","method":"POST","status":"Active","description":"bind resources","type":"Admin","sort_order":8,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"}]}]},{"id":154,"parent_id":152,"type":"Admin","resource_type":"Menu","name":"Subscription","path":"/pricing/subscription","icon":"","tag":"Pricing:Subscription","is_show":true,"status":"Active","sort_order":2,"enable_data_scope":false,"description":"","is_public":false,"created_at":"2026-06-16T03:16:01Z","updated_at":"2026-06-18T02:39:06Z","apis":[{"id":415,"parent_id":410,"name":"page","path":"/api/v1/organization_subscription/page","method":"POST","status":"Active","description":"page","type":"Admin","sort_order":5,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"}],"children":[{"id":160,"parent_id":154,"type":"Admin","resource_type":"Button","name":"Add","path":"","icon":"","tag":"Pricing:Subscription:Add","is_show":false,"status":"Active","sort_order":1,"enable_data_scope":false,"description":"","is_public":false,"created_at":"2026-06-17T08:13:31Z","updated_at":"2026-06-18T02:40:36Z","apis":[{"id":128,"parent_id":121,"name":"select option","path":"/api/v1/organization/select_options","method":"GET","status":"Active","description":"select option","type":"Admin","sort_order":7,"auth_required":true,"permission_required":true,"created_at":"2026-03-11T07:44:46Z","updated_at":"2026-03-11T07:44:46Z"},{"id":411,"parent_id":410,"name":"create","path":"/api/v1/organization_subscription","method":"POST","status":"Active","description":"create","type":"Admin","sort_order":1,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"},{"id":424,"parent_id":422,"name":"all","path":"/api/v1/pricing_plan/all","method":"GET","status":"Active","description":"all","type":"Admin","sort_order":2,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"}]},{"id":161,"parent_id":154,"type":"Admin","resource_type":"Button","name":"Change status","path":"","icon":"","tag":"Pricing:Subscription:ChangeStatus","is_show":false,"status":"Active","sort_order":2,"enable_data_scope":false,"description":"","is_public":false,"created_at":"2026-06-17T08:14:37Z","updated_at":"2026-06-18T02:40:51Z","apis":[{"id":414,"parent_id":410,"name":"switch status","path":"/api/v1/organization_subscription/:id/status","method":"PATCH","status":"Active","description":"switch status","type":"Admin","sort_order":4,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"}]},{"id":162,"parent_id":154,"type":"Admin","resource_type":"Button","name":"Change expire date","path":"","icon":"","tag":"Pricing:Subscription:ChangeExpireDate","is_show":false,"status":"Active","sort_order":3,"enable_data_scope":false,"description":"","is_public":false,"created_at":"2026-06-17T08:15:32Z","updated_at":"2026-06-18T02:41:08Z","apis":[{"id":413,"parent_id":410,"name":"update expire date","path":"/api/v1/organization_subscription/:id/expire_date","method":"PATCH","status":"Active","description":"update expire date","type":"Admin","sort_order":3,"auth_required":true,"permission_required":true,"created_at":"2026-06-18T02:21:35Z","updated_at":"2026-06-18T02:21:35Z"}]}]}]}]
`)
