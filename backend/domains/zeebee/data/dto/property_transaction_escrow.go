package dto

type PropertyTransactionEscrowChecklist struct {
	HasWireInstruction bool `json:"has_wire_instruction"`
	HasEmdReceipt      bool `json:"has_emd_receipt"`
	HasPrelimReport    bool `json:"has_prelim_report"`
	HasNhd             bool `json:"has_nhd"`
	HasColoredMap      bool `json:"has_colored_map"`
}
