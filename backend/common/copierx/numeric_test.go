package copierx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// numeric/pointer shapes mirroring model.RangeInt -> proto.RangeInt32 and the
// nested-pointer metadata mapping in discover_helper.go.

type srcRangeInt struct {
	Min int `json:"min"`
	Max int `json:"max"`
}
type srcRangeFloat struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}
type dstRangeInt32 struct {
	Min int32 `json:"min"`
	Max int32 `json:"max"`
}
type dstRangeFloat struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type srcNested struct {
	Amount *srcRangeInt
	Rate   *srcRangeFloat
}
type dstNested struct {
	Amount *dstRangeInt32
	Rate   *dstRangeFloat
}

func TestTrans_IntToInt32PtrNested(t *testing.T) {
	// int -> int32, value
	d, err := Trans[srcRangeInt, dstRangeInt32](&srcRangeInt{Min: 100, Max: 200})
	assert.NoError(t, err)
	assert.Equal(t, int32(100), d.Min)
	assert.Equal(t, int32(200), d.Max)

	// nested pointers: *srcRangeInt -> *dstRangeInt32 (int->int32) + float passthrough
	n := &srcNested{Amount: &srcRangeInt{Min: 1, Max: 5}, Rate: &srcRangeFloat{Min: 0.5, Max: 9.9}}
	nd, err := Trans[srcNested, dstNested](n)
	assert.NoError(t, err)
	if assert.NotNil(t, nd.Amount) {
		assert.Equal(t, int32(1), nd.Amount.Min)
		assert.Equal(t, int32(5), nd.Amount.Max)
	}
	if assert.NotNil(t, nd.Rate) {
		assert.Equal(t, 0.5, nd.Rate.Min)
		assert.Equal(t, 9.9, nd.Rate.Max)
	}

	// nil nested pointer stays nil
	nd2, err := Trans[srcNested, dstNested](&srcNested{})
	assert.NoError(t, err)
	assert.Nil(t, nd2.Amount)
}
