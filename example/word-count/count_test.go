package main

import (
	"testing"
)

func TestTotalWords(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "en0",
			args: args{s: "hello,web framework"},
			want: 3,
		},
		{
			name: "en1",
			args: args{s: " hello, web framework"},
			want: 3,
		},
		{
			name: "zh0",
			args: args{s: "你好，网站架构"},
			want: 7,
		},
		{
			name: "zh1",
			args: args{s: " 你好，网站 架构"},
			want: 7,
		},
		{
			name: "en_zh0",
			args: args{s: "这是架构：hello, web framework"},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TotalWords(tt.args.s); got != tt.want {
				t.Errorf("TotalWords() = %v, want %v", got, tt.want)
			}
		})
	}
}