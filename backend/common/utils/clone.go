package utils

import (
	"encoding/json"
	"fmt"
)

func Clone[T any](src T) (T, error) {
	var dst T

	data, err := json.Marshal(src)
	if err != nil {
		return dst, fmt.Errorf("Failed To Marshal Source Object: %w", err)
	}

	err = json.Unmarshal(data, &dst)
	if err != nil {
		return dst, fmt.Errorf("Failed To Unmarshal Into Destination Object: %w", err)
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
