package skel

import (
	"testing"

	"xxx.com/skel/test"
)

// 测试/skel/del
func TestDel(t *testing.T) {
	test.Setup()

	skel := GetSkel()
	id := skel.ID

	if err := Del(id); err != nil {
		t.Fatal(err)
	}

	skelNew, err := Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if skelNew != nil {
		t.Fatal("ret skel not valid")
	}
}
