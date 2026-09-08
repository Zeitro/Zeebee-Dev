package operation

import (
	"testing"

	"zeebee-backend/app/borrower/api/internal/types"
	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/domains/zeebee/proto"
)

func TestOperationPageResponseConvertsMetadataToObject(t *testing.T) {
	data := &borrower.OperationPageResponseData{List: []*proto.OperationEntity{{Metadata: `{"task_id":7}`}}}

	resp, err := copierx.Trans[borrower.OperationPageResponseData, types.OperationPageResponse](data, copierx.StrToStrMap())
	if err != nil {
		t.Fatal(err)
	}
	if resp.List[0].Metadata["task_id"] != float64(7) {
		t.Fatalf("unexpected metadata: %#v", resp.List[0].Metadata)
	}
}
