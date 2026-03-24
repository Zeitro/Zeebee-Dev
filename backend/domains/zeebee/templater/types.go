package templater

type Scope string
type FileType string
type Category string

type OrganizationType string

const (
	ScopeSystem       Scope = "System"
	ScopeOrganization Scope = "Organization"

	FileTypeHTML FileType = "Html"
	FileTypePDF  FileType = "Pdf"

	// OrganizationTypeLender   OrganizationType = "Lender"
	// OrganizationTypeBorrower OrganizationType = "Borrower"
)

type TemplateConf struct {
	Purpose          string
	Scope            Scope
	Category         Category
	OrganizationType string
	OrganizationId   int64
	FileType         FileType
}
