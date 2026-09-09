package errorx

func GetFailedTransReqErr(err error) error {
	return Wrap(err, "trans request error")
}

func GetFailedTransRespErr(err error) error {
	return Wrap(err, "trans response error")
}

func GetFailedTransErr(err error) error {
	return Wrap(err, "trans data error")
}
