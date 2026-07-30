package file

import "testing"

func TestIsWord(t *testing.T) {
	cases := map[string]bool{
		"wire.doc":              true,
		"wire.docx":             true,
		"WIRE.DOCX":             true,
		"Escrow.Doc":            true,
		"a.b.docx":              true,
		"macro.docm":            false, // out of scope
		"tmpl.dotm":             false,
		"old.rtf":               false,
		"contract.pdf":          false,
		"photo.png":             false,
		"doc":                   false, // no extension dot
		"notes.doc.txt":         false, // .txt tail
		"instructions.docx.pdf": false, // .pdf tail
	}
	for name, want := range cases {
		if got := IsWord(name); got != want {
			t.Errorf("IsWord(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestWordExt(t *testing.T) {
	cases := map[string]string{
		"a.docx":   "docx",
		"A.Doc":    "doc",
		"a.docm":   "",
		"a.pdf":    "",
		"noext":    "",
	}
	for name, want := range cases {
		if got := WordExt(name); got != want {
			t.Errorf("WordExt(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestPdfDisplayName(t *testing.T) {
	cases := map[string]string{
		"wire.doc":           "wire.pdf",
		"Wire.DOCX":          "Wire.pdf",
		"a.b.docx":           "a.b.pdf",
		"contract.pdf":       "contract.pdf", // non-Word unchanged
		"photo.png":          "photo.png",
		"macro.docm":         "macro.docm", // not in scope, unchanged
	}
	for name, want := range cases {
		if got := PdfDisplayName(name); got != want {
			t.Errorf("PdfDisplayName(%q) = %q, want %q", name, got, want)
		}
	}
}
