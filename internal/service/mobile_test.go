package service

import "testing"

func TestMd5V(t *testing.T) {
	mobile := "17744582020"
	mobileMd5 := Md5V(mobile)
	t.Logf(mobileMd5)
}
