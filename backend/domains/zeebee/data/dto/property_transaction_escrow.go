package dto

import "strings"

type PropertyTransactionEscrowChecklist struct {
	HasWireInstruction bool `json:"has_wire_instruction"`
	HasEmdReceipt      bool `json:"has_emd_receipt"`
	HasPrelimReport    bool `json:"has_prelim_report"`
	HasNhd             bool `json:"has_nhd"`
	HasColoredMap      bool `json:"has_colored_map"`
}

// SetByName flips the checklist flag matching docName (case-insensitive) and
// reports whether name is a known escrow doc. Upload sets true; delete sets
// false once the doc has no files left.
func (c *PropertyTransactionEscrowChecklist) SetByName(docName string, v bool) bool {
	switch strings.ToLower(docName) {
	case "wire instructions":
		c.HasWireInstruction = v
	case "emd receipt":
		c.HasEmdReceipt = v
	case "prelim report":
		c.HasPrelimReport = v
	case "nhd report":
		c.HasNhd = v
	case "colored map":
		c.HasColoredMap = v
	default:
		return false
	}
	return true
}
