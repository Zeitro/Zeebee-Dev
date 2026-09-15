package copierx

import (
	"reflect"

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
	prepareDestinationCollections(src, dist)
	opt := copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters:    opts,
	}
	err = copier.CopyWithOption(dist, src, opt)
	if err != nil {
		return nil, errorx.GetFailedTransErr(err)
	}
	return dist, err
}

// Assign assigns the values from src to dist. It performs a deep copy of the source value to the destination with case sensitivity.
func Assign[T, U any](src *T, dist *U, opts ...copier.TypeConverter) (err error) {
	if src == nil {
		return nil
	}

	if dist == nil {
		dist = new(U)
	}
	prepareDestinationCollections(src, dist)
	opt := copier.Option{
		DeepCopy:      true,
		CaseSensitive: true,
		Converters:    opts,
	}
	err = copier.CopyWithOption(dist, src, opt)
	if err != nil {
		return errorx.GetFailedTransErr(err)
	}
	return err
}

// prepareDestinationCollections clears non-nil slice targets before copier copies them.
// copier overwrites existing slices without removing stale entries.
func prepareDestinationCollections(src, dst any) {
	prepareCollections(reflect.ValueOf(src), reflect.ValueOf(dst))
}

func prepareCollections(src, dst reflect.Value) {
	for src.IsValid() && src.Kind() == reflect.Pointer {
		if src.IsNil() {
			return
		}
		src = src.Elem()
	}
	for dst.IsValid() && dst.Kind() == reflect.Pointer {
		if dst.IsNil() {
			return
		}
		dst = dst.Elem()
	}
	if !src.IsValid() || !dst.IsValid() {
		return
	}
	if src.Kind() == reflect.Slice && dst.Kind() == reflect.Slice {
		if !src.IsNil() && dst.CanSet() {
			dst.Set(reflect.Zero(dst.Type()))
		}
		return
	}
	if src.Kind() != reflect.Struct || dst.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		dstField := dst.FieldByName(src.Type().Field(i).Name)
		if !dstField.IsValid() {
			continue
		}
		prepareCollections(srcField, dstField)
	}
}
