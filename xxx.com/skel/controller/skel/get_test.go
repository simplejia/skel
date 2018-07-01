package skel

import (
	"testing"

	"xxx.com/skel/test"
)

// 测试/skel/get
func TestGet(t *testing.T) {
	test.Setup()

	skel := GetSkel()
	if skel == nil {
		t.Fatal("test fail")
	}
}
