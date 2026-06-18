package uuid

import "testing"

func TestGetSonyflakeID(t *testing.T) {
	id, err := GetSonyflakeID()
	if err != nil {
		t.Fatalf("fail to get uuid error:%v", err)
	}
	t.Logf("sonyflake id: %d", id)
}

func TestGetSonyflakeIDString(t *testing.T) {
	id, err := GetSonyflakeIDString()
	if err != nil {
		t.Fatalf("fail to get uuid error:%v", err)
	}
	t.Logf("sonyflake id: %s", id)
}
