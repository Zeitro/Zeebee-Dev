package uuid

import (
	"fmt"

	"github.com/sony/sonyflake/v2"
)

var defaultSonyflake *sonyflake.Sonyflake

func init() {
	st := sonyflake.Settings{}
	defaultSonyflake, _ = sonyflake.New(st)
}

func GetSonyflakeID() (id int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("GetSonyflakeID error: %v", r)
		}
	}()

	id, err = defaultSonyflake.NextID()
	return id, err
}

func GetSonyflakeIDString() (uuid string, err error) {
	id, err := GetSonyflakeID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), err
}
