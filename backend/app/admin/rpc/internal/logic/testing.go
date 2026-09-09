package logic

import "sync"

// ResetBaseLogicForTest clears the cached baseLogic singleton so each test
// wires its own ServiceContext. Test-only; production always uses the
// once-cached instance. Mirrors the lender / tenant base_logic helper.
func ResetBaseLogicForTest() {
	baseLogic = nil
	onceBase = sync.Once{}
}
