package logic

import (
	"encoding/json"
	"reflect"
)

func projectManagementOperationMetadata(locator, before, after map[string]any) ([]byte, error) {
	metadata := make(map[string]any, len(locator)+2)
	for key, value := range locator {
		metadata[key] = value
	}
	if len(before) != 0 {
		metadata["before"] = before
	}
	if len(after) != 0 {
		metadata["after"] = after
	}
	return json.Marshal(metadata)
}

func operationValuesDiff(before, after map[string]any) (map[string]any, map[string]any) {
	changedBefore, changedAfter := make(map[string]any), make(map[string]any)
	// Operation value pairs are built by the same helper and therefore share keys.
	for key, beforeValue := range before {
		if afterValue := after[key]; !reflect.DeepEqual(beforeValue, afterValue) {
			changedBefore[key], changedAfter[key] = beforeValue, afterValue
		}
	}
	return changedBefore, changedAfter
}
