// Package libreoffice converts doc/docx to pdf via a headless LibreOffice
// (`soffice`) subprocess. Concurrency is gated by a profile pool: each slot
// owns a pre-warmed UserInstallation dir so concurrent conversions don't trip
// LibreOffice's single-instance-per-profile lock.
package libreoffice

import (
	"os"
	"path/filepath"
)

// safeProfileXcu hardens a profile: macro security set to VeryHigh so embedded
// macros never auto-run during conversion.
//
// External-link / remote-image fetching is blocked at the network layer
// (pdf_generator container has no outbound egress except S3) — that is the
// primary, reliable control; this profile adds macro hardening on top.
const safeProfileXcu = `<?xml version="1.0" encoding="UTF-8"?>
<oor:items xmlns:oor="http://openoffice.org/2001/registry" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
 <item oor:path="/org.openoffice.Office.Common/Security/Scripting"><prop oor:name="MacroSecurityLevel" oor:op="fuse"><value>3</value></prop></item>
</oor:items>
`

// ensureSafeProfile writes the hardened registrymodifications.xcu into the slot's
// UserInstallation dir. Idempotent — skipped once present so a reused slot pays
// the write only at pool warm-up.
func ensureSafeProfile(profileDir string) error {
	userDir := filepath.Join(profileDir, "user")
	if err := os.MkdirAll(userDir, 0o750); err != nil {
		return err
	}
	xcu := filepath.Join(userDir, "registrymodifications.xcu")
	if _, err := os.Stat(xcu); err == nil {
		return nil
	}
	return os.WriteFile(xcu, []byte(safeProfileXcu), 0o600)
}
