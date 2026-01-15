package copierx

import (
	"database/sql"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"zeebee-backend/common/errorx"
)

// TimeToRFC3339Str time.Time -> string (2006-01-02T15:04:05Z07:00)
func TimeToRFC3339Str() copier.TypeConverter {
	return TimeToStr(time.RFC3339)
}

// TimeToStr time.Time -> string
func TimeToStr(layout string) copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: time.Time{},
		DstType: copier.String,
		Fn: func(src interface{}) (dst interface{}, err error) {
			val, ok := src.(time.Time)
			if !ok {
				return "", errorx.Errorf("Invalid Src Type:%T", val)
			}
			if val.IsZero() {
				return "", nil
			}

			return val.Format(layout), nil
		},
	}
}

func TimePtrToRFC3339Str() copier.TypeConverter {
	return TimePtrToStr(time.RFC3339)
}
func TimePtrToStr(layout string) copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: &time.Time{},
		DstType: copier.String,
		Fn: func(src interface{}) (dst interface{}, err error) {
			val, ok := src.(*time.Time)
			if !ok {
				return "", errorx.Errorf("Invalid Src Type:%T", val)
			}
			if val == nil || val.IsZero() {
				return "", nil
			}
			return val.Format(layout), nil
		},
	}
}

// SqlTimePtrToRFC3339Str *sql.NullTime -> string (2006-01-02T15:04:05Z07:00)
func SqlTimePtrToRFC3339Str() copier.TypeConverter {
	return SqlTimePtrToStr(time.RFC3339)
}

// SqlTimePtrToStr *sql.NullTime -> string
func SqlTimePtrToStr(layout string) copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: &sql.NullTime{},
		DstType: copier.String,
		Fn: func(src interface{}) (dst interface{}, err error) {
			val, ok := src.(*sql.NullTime)
			if !ok {
				return "", errorx.Errorf("Invalid Src Type:%T", val)
			}
			if val == nil || !val.Valid || val.Time.IsZero() {
				return "", nil
			}

			return val.Time.Format(layout), nil
		},
	}
}

// StrToSqlTimePtr string -> *sql.NullTime
func StrToSqlTimePtr(layout string) copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: copier.String,
		DstType: &sql.NullTime{},
		Fn: func(src interface{}) (dst interface{}, err error) {
			val, ok := src.(string)
			if !ok {
				return "", errors.Errorf("Invalid Src Type:%T", val)
			}
			dst = &sql.NullTime{}
			if len(val) == 0 {
				return dst, nil
			}
			t, err := time.Parse(layout, val)
			if err != nil {
				return dst, errorx.WithStack(err)
			}

			dst = &sql.NullTime{
				Time:  t,
				Valid: true,
			}
			return dst, nil
		},
	}
}

// UsStrToTime string(mm/dd/yyyy) -> time
func UsStrToTime() copier.TypeConverter {
	return StrToTime("01/02/2006")
}

// StrToTime string -> time
func StrToTime(format string) copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: copier.String,
		DstType: time.Time{},
		Fn: func(src interface{}) (dst interface{}, err error) {
			val, ok := src.(string)
			if !ok {
				return "", errorx.Errorf("Invalid Src Type:%T", val)
			}
			dst, err = time.Parse(format, val)

			return dst, err
		},
	}
}
