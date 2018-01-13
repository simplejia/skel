package skel

import (
	"testing"

	"xxx.com/skel/test"

	"github.com/simplejia/lc"
)

func init() {
	lc.Disabled = true
}

// 测试/skel/del
func TestDel(t *testing.T) {
	test.Setup()

	skel := GetSkel()
	id := skel.ID

	err := Del(id)
	if err != nil {
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
