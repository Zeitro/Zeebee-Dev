package projectmanagementprojectlogic

import (
	"database/sql"
	"testing"
	"time"

	"zeebee-backend/app/borrower/rpc/pb/borrower"
	"zeebee-backend/common/copierx"
	"zeebee-backend/common/timex"
	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/data/dto"
	"zeebee-backend/domains/zeebee/model"
	"zeebee-backend/domains/zeebee/proto"
)

func TestProjectScheduleRequestAndResponseConversion(t *testing.T) {
	input, err := copierx.Trans[borrower.ProjectManagementProjectCreateRequestData, dto.ProjectManagementProjectCreateDto](
		&borrower.ProjectManagementProjectCreateRequestData{ProjectType: "Planned", StartAt: "01/31/2026", PeriodValue: 2, PeriodUnit: "Month"},
		copierx.StrToSqlTimePtr(timex.LayoutUsStyle),
	)
	if err != nil || input.ProjectType != constants.ProjectManagementProjectTypePlanned || input.StartAt == nil || input.StartAt.Time.Format(time.DateOnly) != "2026-01-31" || input.PeriodValue != 2 || input.PeriodUnit != constants.ProjectManagementPeriodUnitMonth {
		t.Fatalf("schedule request conversion failed: %#v, %v", input, err)
	}
	endAt := &sql.NullTime{Time: time.Date(2026, time.February, 28, 0, 0, 0, 0, time.UTC), Valid: true}
	entity, err := copierx.Trans[model.ProjectManagementProject, proto.ProjectManagementProjectEntity](
		&model.ProjectManagementProject{ProjectType: constants.ProjectManagementProjectTypePlanned, EndAt: endAt},
		copierx.TimeToRFC3339Str(), copierx.SqlTimePtrToStr(timex.LayoutUsStyle),
	)
	if err != nil || entity.EndAt != "02/28/2026" {
		t.Fatalf("derived end date response conversion failed: %#v, %v", entity, err)
	}
}
