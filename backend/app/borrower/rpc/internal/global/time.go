package global

import "time"

func initTime() (err error) {
	// set utc time zone
	time.Local = time.UTC

	return err
}
