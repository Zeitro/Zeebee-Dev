package types

type LoanExtendValues struct {
	Pricing          string `json:"pricing"`
	LoanStatusUpdate string `json:"loan_status_update"`
	MeetingNotes     string `json:"meeting_notes"`
	ActionList       string `json:"action_list"`
	BusinessLine     string `json:"business_line"`
}
