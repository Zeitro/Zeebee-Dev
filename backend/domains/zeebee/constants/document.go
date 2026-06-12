package constants

type DocumentType string

const (
	DocumentTypeRoot   DocumentType = "Root"
	DocumentTypeBranch DocumentType = "Branch"
	DocumentTypeLeaf   DocumentType = "Leaf"
)
