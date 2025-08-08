package constants

import "zeebee-backend/common/enums"

type ApplicantType string

const ApplicantTypeName = "applicant_type"

const (
	ApplicantTypeIndividual  ApplicantType = "Individual"
	ApplicantTypeCorporation ApplicantType = "Corporation"
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
