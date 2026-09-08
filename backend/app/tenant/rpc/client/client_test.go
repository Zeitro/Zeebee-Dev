package client

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/zrpc"
)

// TestNew_AllSubClientsPopulated asserts every sub-client field is wired. The
// constructors only store the zrpc.Client (no dial), so a nil client suffices.
// Iterating by reflection keeps this green as new sub-clients are added — a
// field that New forgets to populate fails here.
func TestNew_AllSubClientsPopulated(t *testing.T) {
	var cli zrpc.Client
	c := New(cli)
	require.NotNil(t, c)

	v := reflect.ValueOf(c).Elem()
	for i := range v.NumField() {
		require.NotNilf(t, v.Field(i).Interface(),
			"field %s must be populated by New", v.Type().Field(i).Name)
	}
}
