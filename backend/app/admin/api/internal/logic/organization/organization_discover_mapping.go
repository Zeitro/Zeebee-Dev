package organization

import (
	"zeebee-backend/app/admin/api/internal/types"
	"zeebee-backend/app/admin/rpc/pb/admin"
	"zeebee-backend/common/copierx"
)

func discoverProtoSettingsToTypes(p *admin.LenderDiscoverSettingsEntity) (types.LenderDiscoverSettings, error) {
	out, err := copierx.Trans[admin.LenderDiscoverSettingsEntity, types.LenderDiscoverSettings](p)
	if err != nil {
		return types.LenderDiscoverSettings{}, err
	}
	return *out, nil
}

func discoverProtoRunToTypes(p *admin.LenderDiscoverRunEntity) (types.LenderDiscoverRunEntity, error) {
	if p == nil {
		return types.LenderDiscoverRunEntity{Stats: types.LenderDiscoverRunStats{}}, nil
	}
	out, err := copierx.Trans[admin.LenderDiscoverRunEntity, types.LenderDiscoverRunEntity](p)
	if err != nil {
		return types.LenderDiscoverRunEntity{}, err
	}
	return *out, nil
}

func discoverProtoPreviewDataToTypes(p *admin.LenderDiscoverPreviewResponseData) (types.LenderDiscoverPreviewData, error) {
	if p == nil {
		return types.LenderDiscoverPreviewData{}, nil
	}
	out, err := copierx.Trans[admin.LenderDiscoverPreviewResponseData, types.LenderDiscoverPreviewData](p)
	if err != nil {
		return types.LenderDiscoverPreviewData{}, err
	}
	return *out, nil
}

func discoverTypesCommitRequestToProto(req *types.LenderDiscoverCommitRequest) (*admin.LenderDiscoverCommitRequest, error) {
	if req == nil {
		return &admin.LenderDiscoverCommitRequest{Data: &admin.LenderDiscoverCommitRequestData{}}, nil
	}
	list := make([]*admin.LenderDiscoverCommitLender, 0, len(req.Lenders))
	for i := range req.Lenders {
		ent, err := copierx.Trans[types.LenderDiscoverCommitLender, admin.LenderDiscoverCommitLender](&req.Lenders[i])
		if err != nil {
			return nil, err
		}
		list = append(list, ent)
	}
	return &admin.LenderDiscoverCommitRequest{
		Data: &admin.LenderDiscoverCommitRequestData{Lenders: list},
	}, nil
}
