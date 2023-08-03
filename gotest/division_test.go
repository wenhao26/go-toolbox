package gotest

import (
	"testing"
)

func Test_Division1(t *testing.T) {
	if res, err := Division(6, 2); res !=3 || err != nil {
		t.Error("除法函数测试不通过")
	} else {
		t.Log("测试通过")
	}
}

func Test_Division2(t *testing.T) {
	if _, err := Division(6, 0); err == nil {
		t.Error("除法函数测试不通过")
	} else {
		t.Log("测试通过:", err)
	}
}

func Benchmark_Division1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Division(3, 8)
	}
}

func Benchmark_Division2(b *testing.B) {
	b.StopTimer()
	// todo ...
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Division(4, 7)
	}
}
