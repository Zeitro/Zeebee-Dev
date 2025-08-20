package constants

import "zeebee-backend/common/enums"

// Individual, Corporation,LLC,Partnership,Trust
type ApplicantType string

const ApplicantTypeName = "applicant_type"

const (
	ApplicantTypeIndividual  ApplicantType = "Individual"
	ApplicantTypeCorporation ApplicantType = "Corporation"
	ApplicantTypeLLC         ApplicantType = "LLC"
	ApplicantTypePartnership ApplicantType = "Partnership"
	ApplicantTypeTrust       ApplicantType = "Trust"
)

type ApplicantInvitationStatus string

const ApplicantInvitationStatusName = "applicant_invitation_status"

const (
	ApplicantInvitationStatusPending  ApplicantInvitationStatus = "Pending"
	ApplicantInvitationStatusAccepted ApplicantInvitationStatus = "Accepted"
	ApplicantInvitationStatusRejected ApplicantInvitationStatus = "Rejected"
)

func init() {
	// register applicant type
	enums.Register(ApplicantTypeName, []enums.EnumOpt{
		{
			Label: string(ApplicantTypeIndividual),
			Value: string(ApplicantTypeIndividual),
		},
		{
			Label: string(ApplicantTypeCorporation),
			Value: string(ApplicantTypeCorporation),
		},
		{
			Label: string(ApplicantTypePartnership),
			Value: string(ApplicantTypePartnership),
		},
		{
			Label: string(ApplicantTypeLLC),
			Value: string(ApplicantTypeLLC),
		},
		{
			Label: string(ApplicantTypeTrust),
			Value: string(ApplicantTypeTrust),
		},
	})
	// register applicant invitation status
	enums.Register(ApplicantInvitationStatusName, []enums.EnumOpt{
		{
			Label: string(ApplicantInvitationStatusPending),
			Value: string(ApplicantInvitationStatusPending),
		},
		{
			Label: string(ApplicantInvitationStatusAccepted),
			Value: string(ApplicantInvitationStatusAccepted),
		},
		{
			Label: string(ApplicantInvitationStatusRejected),
			Value: string(ApplicantInvitationStatusRejected),
		},
	})
}
