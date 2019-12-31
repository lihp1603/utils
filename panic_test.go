package utils

import (
	"testing"
)

//测试
func TestPanicHandler(t *testing.T) {
	defer t.Error("in main")
	testDemo()
	t.Error("process main")
}

func testDemo() {
	defer PanicHandler()
	panic("demo unknown err")
}
