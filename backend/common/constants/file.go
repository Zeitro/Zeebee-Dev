package constants

import (
	"github.com/duke-git/lancet/v2/slice"

	"zeebee-backend/common/enums"
)

type FileType string

const FileTypeName = "file_type"

const (
	FileTypeApplyForm        FileType = "ApplyForm"
	FileTypeIdentification   FileType = "Identification"
	FileTypeIncome           FileType = "Income"
	FileTypeProperty         FileType = "Property"
	FileTypeAssets           FileType = "Assets"
	FileTypeCredit           FileType = "Credit"
	FileTypeMiscellaneous    FileType = "Miscellaneous"
	FileTypeInternalDocument FileType = "InternalDocument"
)

func (c FileType) IsValid() bool {
	arr := []FileType{
		FileTypeApplyForm,
		FileTypeIdentification,
		FileTypeIncome,
		FileTypeProperty,
		FileTypeAssets,
		FileTypeCredit,
		FileTypeMiscellaneous,
		FileTypeInternalDocument,
	}
	_, ok := slice.FindBy(arr, func(index int, item FileType) bool {
		return c == item
	})
	return ok
}

type FileBizType string

const FileBizTypeName = "file_biz_type"

const (
	FileBizTypeLoan FileBizType = "Loan"
)

func init() {
	enums.Register(FileTypeName, []enums.EnumOpt{
		{
			Label: string(FileTypeApplyForm),
			Value: string(FileTypeApplyForm),
		},
		{
			Label: string(FileTypeIdentification),
			Value: string(FileTypeIdentification),
		},
		{
			Label: string(FileTypeIncome),
			Value: string(FileTypeIncome),
		},
		{
			Label: string(FileTypeProperty),
			Value: string(FileTypeProperty),
		},
		{
			Label: string(FileTypeAssets),
			Value: string(FileTypeAssets),
		},
		{
			Label: string(FileTypeCredit),
			Value: string(FileTypeCredit),
		},
		{
			Label: string(FileTypeMiscellaneous),
			Value: string(FileTypeMiscellaneous),
		},
		{
			Label: string(FileTypeInternalDocument),
			Value: string(FileTypeInternalDocument),
		},
	})
}
