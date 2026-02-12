package errorx

func GetFailedTransReqErr(err error) error {
	return Wrap(err, "Trans request error")
}

func GetFailedTransRespErr(err error) error {
	return Wrap(err, "Trans response error")
}

func GetFailedTransErr(err error) error {
	return Wrap(err, "Trans data error")
}
