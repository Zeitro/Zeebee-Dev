package types

// LoanExtendValues represents extended information related to a loan.
type LoanExtendValues struct {
	Pricing           string `json:"pricing"`
	LoanStatusUpdates string `json:"loan_status_updates"`
	MeetingNotes      string `json:"meeting_notes"`
	ActionList        string `json:"action_list"`
	BusinessLine      string `json:"business_line"`
}
