package copierx

import (
	"github.com/jinzhu/copier"

	"zeebee-backend/common/errorx"
)

// Trans converts a source value of type T to a destination value of type U.
// It performs a deep copy of the source value to the destination with case sensitivity.
// Type converters can be provided as optional arguments to handle custom type conversions.
//
// Parameters:
//   - src: A pointer to the source value of type T to be converted
//   - opts: Optional type converters for custom type conversion handling
//
// Returns:
//   - *U: A pointer to the newly created destination value of type U
//   - error: An error if the conversion fails, nil otherwise
//
// Notes:
//   - If src is nil, returns (nil, nil)
//   - Uses copier.CopyWithOption for the actual conversion
//   - Wraps any conversion error with errorx.GetFailedTransErr
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
