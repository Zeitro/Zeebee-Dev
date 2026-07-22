package constants

type DocumentTemplateType string

const (
	DocumentTemplateTypeBorrowerPropertyTransaction DocumentTemplateType = "BorrowerPropertyTransaction"
)

type DocumentTemplateSource string

const (
	DocumentTemplateFileParseSourceSingleFile  DocumentTemplateSource = "SingleFile"
	DocumentTemplateFileParseSourceFiles       DocumentTemplateSource = "Files"
	DocumentTemplateFileParseSourceDirectory   DocumentTemplateSource = "Directory"
	DocumentTemplateFileParseSourceTagFiltered DocumentTemplateSource = "TagFiltered"
)

type DocumentTemplateTrigger string

const (
	FileParseTriggeredTypeUpload DocumentTemplateTrigger = "Upload"
	FileParseTriggeredTypeManual DocumentTemplateTrigger = "Manual"
)

type DocumentTemplateParser string
const (
	DocumentTemplateFileParseParserRpa        DocumentTemplateParser = "Rpa"
	DocumentTemplateFileParseParserEscrowInfo DocumentTemplateParser = "EscrowInfo"
)
