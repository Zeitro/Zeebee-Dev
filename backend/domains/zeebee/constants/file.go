package constants

type FileEntityType string

const FileEntityTypeName = "file_entity_type"

const (
	FileEntityTypeLoanDocument          FileEntityType = "LoanDocument"
	FileEntityTypeDocument              FileEntityType = "Document"
	FileEntityTypeProjectManagementTask FileEntityType = "ProjectManagementTask"
)

func init() {

}
