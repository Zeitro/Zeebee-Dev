package constants

import "fmt"

type CommonStatus string

const (
	CommonStatusActive   CommonStatus = "Active"
	CommonStatusDisabled CommonStatus = "Disabled"
)

func (s CommonStatus) Validate() (err error) {
	if s != CommonStatusActive && s != CommonStatusDisabled {
		err = fmt.Errorf("invalid status:%v", s)
	}
	return
}
