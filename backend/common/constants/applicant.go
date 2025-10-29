package constants

import (
	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/common/enums"
)

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

func (c ApplicantType) IsValid() bool {
	var ApplicantTypeValues = []ApplicantType{
		ApplicantTypeIndividual,
		ApplicantTypeCorporation,
		ApplicantTypeLLC,
		ApplicantTypePartnership,
		ApplicantTypeTrust,
	}
	_, isValid := slice.FindBy(ApplicantTypeValues, func(index int, value ApplicantType) bool {
		return value == c
	})
	return isValid
}

type ApplicantInvitationStatus string

const ApplicantInvitationStatusName = "applicant_invitation_status"

const (
	ApplicantInvitationStatusPending  ApplicantInvitationStatus = "Pending"
	ApplicantInvitationStatusAccepted ApplicantInvitationStatus = "Accepted"
	ApplicantInvitationStatusRejected ApplicantInvitationStatus = "Rejected"
)

type MaritalStatus string

const MaritalStatusName = "marital_status"
const (
	MaritalStatusMarried   MaritalStatus = "Married"
	MaritalStatusUnmarried MaritalStatus = "Unmarried"
	MaritalStatusSeparated MaritalStatus = "Separated"
)

type TINType string

const TINTypeName = "tin_type"
const (
	TINTypeSSN  TINType = "SSN"
	TINTypeEIN  TINType = "EIN"
	TINTypeITIN TINType = "ITIN"
)

type ApplicantSex string

const ApplicantSexName = "applicant_sex"
const (
	ApplicantSexMale   ApplicantSex = "Male"
	ApplicantSexFemale ApplicantSex = "Female"
	ApplicantSexHidden ApplicantSex = "Hidden"
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
	// register marital status
	enums.Register(MaritalStatusName, []enums.EnumOpt{
		{
			Label: string(MaritalStatusMarried),
			Value: string(MaritalStatusMarried),
		},
		{
			Label: string(MaritalStatusUnmarried),
			Value: string(MaritalStatusUnmarried),
		},
		{
			Label: string(MaritalStatusSeparated),
			Value: string(MaritalStatusSeparated),
		},
	})
	// register tin type
	enums.Register(TINTypeName, []enums.EnumOpt{
		{
			Label: string(TINTypeSSN),
			Value: string(TINTypeSSN),
		},
		{
			Label: string(TINTypeEIN),
			Value: string(TINTypeEIN),
		},
		{
			Label: string(TINTypeITIN),
			Value: string(TINTypeITIN),
		},
	})
	// register applicant sex
	enums.Register(ApplicantSexName, []enums.EnumOpt{
		{
			Label: string(ApplicantSexMale),
			Value: string(ApplicantSexMale),
		},
		{
			Label: string(ApplicantSexFemale),
			Value: string(ApplicantSexFemale),
		},
		{
			Label: "I do not wish to provide this information",
			Value: string(ApplicantSexHidden),
		},
	})
}
