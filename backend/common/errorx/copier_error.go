package errorx

import "github.com/pkg/errors"

func GetFailedTransReqErr(err error) error {
	return errors.Wrap(err, "trans request error")
}

func GetFailedTransRespErr(err error) error {
	return errors.Wrap(err, "trans response error")
}

func GetFailedTransErr(err error) error {
	return errors.Wrap(err, "trans data error")
}
