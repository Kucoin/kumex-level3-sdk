package log

import "testing"

func TestInfo(t *testing.T) {
	Info("this is a info msg: %s", "hello world.")
}

func TestWarn(t *testing.T) {
	Warn("this is a warn msg: %s", "hello world.")
}

func TestError(t *testing.T) {
	Error("this is a error msg: %s", "hello world.")
}

func TestFatal(t *testing.T) {
	//Fatal("this is a error msg: %s", "hello world.")
}
