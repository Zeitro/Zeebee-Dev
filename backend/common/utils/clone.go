package utils

import (
	"encoding/json"
	"fmt"
)

func Clone[T any](src T) (T, error) {
	var dst T

	data, err := json.Marshal(src)
	if err != nil {
		return dst, fmt.Errorf("Failed to marshal source object: %w", err)
	}

	err = json.Unmarshal(data, &dst)
	if err != nil {
		return dst, fmt.Errorf("Failed to unmarshal into destination object: %w", err)
	}

	return dst, nil
}

func TransStruct[T any, V any](src T, dist V) error {
	contents, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contents, dist)
	return err
}
