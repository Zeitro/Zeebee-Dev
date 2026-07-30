package file

import "strings"

// IsDocx reports whether name has a .docx extension (OOXML), case-insensitive.
func IsDocx(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".docx")
}

// IsDoc reports whether name has a .doc extension (97-2003 OLE binary),
// case-insensitive. Excludes .docx/.docm (those don't end in ".doc").
func IsDoc(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), ".doc")
}

// IsWord reports whether name is a Word doc/docx. Excludes .docm/.dotm/.rtf
// (out of scope: only .doc/.docx are parsed).
func IsWord(name string) bool {
	return IsDocx(name) || IsDoc(name)
}

// WordExt returns the soffice import filter ext ("docx"/"doc") for a Word file,
// "" otherwise.
func WordExt(name string) string {
	switch {
	case IsDocx(name):
		return "docx"
	case IsDoc(name):
		return "doc"
	default:
		return ""
	}
}

// PdfDisplayName swaps a Word filename's extension to .pdf. Provider Files APIs
// classify by filename extension, not declared content type — a .docx name
// carrying pdf bytes is misclassified, so the converted stream is renamed.
// Non-Word names pass through unchanged.
func PdfDisplayName(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.HasSuffix(lower, ".docx"):
		return name[:len(name)-len(".docx")] + ".pdf"
	case strings.HasSuffix(lower, ".doc"):
		return name[:len(name)-len(".doc")] + ".pdf"
	default:
		return name
	}
}
