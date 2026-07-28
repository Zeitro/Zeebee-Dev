package file

import "testing"

func TestIsImage(t *testing.T) {
	cases := map[string]bool{
		"scan.png": true, "photo.jpg": true, "IMG.JPEG": true,
		"wire.pdf": false, "notes.docx": false, "noext": false,
	}
	for name, want := range cases {
		if got := IsImage(name); got != want {
			t.Errorf("IsImage(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestImageContentType(t *testing.T) {
	if ct := ImageContentType("a.png"); ct != "image/png" {
		t.Errorf("png = %q, want image/png", ct)
	}
	if ct := ImageContentType("a.jpg"); ct != "image/jpeg" {
		t.Errorf("jpg = %q, want image/jpeg", ct)
	}
	if ct := ImageContentType("a.jpeg"); ct != "image/jpeg" {
		t.Errorf("jpeg = %q, want image/jpeg", ct)
	}
}
