package timex

import "time"

const (
	LayoutUsStyle = "01/02/2006"
)

func FormateStrToTime(str string, layout string) (t time.Time, err error) {
	return time.Parse(layout, str)
}

func FormateTimeToStrWithUsStyle(str string) (t time.Time, err error) {
	return FormateStrToTime(str, LayoutUsStyle)
}
