package commonlogic

import (
	"context"
	"encoding/json"

	"zeebee-backend/app/borrower/rpc/internal/svc"
	"zeebee-backend/common/errorx"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/proto"
)

type SyncDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncDataLogic {
	return &SyncDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncDataLogic) SyncData(in *proto.EmptyRequest) (resp *proto.EmptyResponse, err error) {
	if err = l.migrateDocumentTemplate(); err != nil {
		return nil, err
	}

	return &proto.EmptyResponse{}, nil
}

func (l *SyncDataLogic) migrateDocumentTemplate() error {

	tpl, err := l.svcCtx.AllRepos.DocumentTemplateRepo.GetByType(l.ctx, constants.DocumentTemplateTypeBorrowerPropertyTransaction)
	if err != nil {
		return err
	}
	if tpl == nil {
		return nil
	}
	changed := false
	for key, nodes := range tpl.Structure {
		if migrateStructureNodes(nodes) {
			changed = true
			tpl.Structure[key] = nodes
		}
	}
	if !changed {
		return nil
	}
	// Map-based Updates bypass the gorm json serializer, so pre-marshal and
	// hand the column a JSON string the jsonb column accepts verbatim.
	raw, err := json.Marshal(tpl.Structure)
	if err != nil {
		return errorx.Wrap(err, "marshal structure")
	}
	return l.svcCtx.AllRepos.DocumentTemplateRepo.UpdateColumnsById(l.ctx, tpl.Id, map[string]interface{}{
		"structure": string(raw),
	})
}

func migrateStructureNodes(nodes []*dto.DocumentTemplateStructure) bool {
	changed := false
	for i := range nodes {
		if migrateStructureNode(nodes[i]) {
			changed = true
		}
	}
	return changed
}

func migrateStructureNode(n *dto.DocumentTemplateStructure) bool {
	changed := false
	switch n.Name {
	case "Wire Instructions":
		if !hasParserRule(n.FileParseConfig, constants.DocumentTemplateFileParseParserEscrowInfo) {
			n.FileParseConfig = &dto.DocumentTemplateFileParseConfig{Rules: []dto.DocumentTemplateFileParseRule{{
				Parser:      constants.DocumentTemplateFileParseParserEscrowInfo,
				Source:      constants.DocumentTemplateFileParseSourceSingleFile,
				TriggeredBy: constants.FileParseTriggeredTypeUpload,
			}}}
			changed = true
		}
	case "Contract":
		// Promote rpa to the folder: parse the entire subtree on manual trigger.
		if !hasParserRule(n.FileParseConfig, constants.DocumentTemplateFileParseParserRPA) {
			if n.FileParseConfig == nil {
				n.FileParseConfig = &dto.DocumentTemplateFileParseConfig{}
			}
			n.FileParseConfig.Rules = append(n.FileParseConfig.Rules, dto.DocumentTemplateFileParseRule{
				Parser:      constants.DocumentTemplateFileParseParserRPA,
				Source:      constants.DocumentTemplateFileParseSourceDirectory,
				TriggeredBy: constants.FileParseTriggeredTypeManual,
				Label:       "Original RPA",
			})
			changed = true
		}
	case "Original RPA":
		// Parsed via the Contract folder now; drop the leaf-level rule.
		if n.FileParseConfig != nil {
			n.FileParseConfig = nil
			changed = true
		}
	}
	for i := range n.Children {
		if migrateStructureNode(n.Children[i]) {
			changed = true
		}
	}
	return changed
}

func hasParserRule(cfg *dto.DocumentTemplateFileParseConfig, parser constants.DocumentTemplateParser) bool {
	if cfg == nil {
		return false
	}
	for _, r := range cfg.Rules {
		if r.Parser == parser {
			return true
		}
	}
	return false
}
