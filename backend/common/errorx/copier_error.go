package errorx

func GetFailedTransReqErr(err error) error {
	return Wrap(err, "Trans Request Error")
}

func GetFailedTransRespErr(err error) error {
	return Wrap(err, "Trans Response Error")
}

func GetFailedTransErr(err error) error {
	return Wrap(err, "Trans Data Error")
}
