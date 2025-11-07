package errorx

import "github.com/pkg/errors"

func GetFailedTransReqErr(err error) error {
	return errors.Wrap(err, "Trans Request Error")
}

func GetFailedTransRespErr(err error) error {
	return errors.Wrap(err, "Trans Response Error")
}

func GetFailedTransErr(err error) error {
	return errors.Wrap(err, "Trans Data Error")
}

func New(text string) (err error) {
	return WithStack(errors.New(text))
}

func Errorf(format string, args ...interface{}) (err error) {
	return errors.Errorf(format, args...)
}

func Wrap(err error, text string) (newErr error) {
	return errors.Wrap(err, text)
}
func Wrapf(err error, format string, args ...interface{}) (newErr error) {
	return errors.Wrapf(err, format, args...)
}
func Unwrap(err error) (newErr error) {
	return errors.Unwrap(err)
}
func WithStack(err error) (newErr error) {
	return errors.WithStack(err)
}

func NotFound(record string) (err error) {
	if len(record) == 0 {
		return New("Record Not Found")
	}
	return Errorf("%s Not Found", record)
}
