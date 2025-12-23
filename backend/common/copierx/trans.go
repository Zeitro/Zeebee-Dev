package copierx

import (
	"github.com/jinzhu/copier"

	"zeebee-backend/common/errorx"
)

func Trans[T, U any](src *T, opts ...copier.TypeConverter) (dist *U, err error) {
	if src == nil {
		return nil, nil
	}
	dist = new(U)
	err = copier.CopyWithOption(dist, src, copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters:    opts,
	})
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	return dist, err
}
