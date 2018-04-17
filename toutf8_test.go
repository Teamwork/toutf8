package toutf8

import (
	"testing"
	"unicode/utf8"

	"github.com/teamwork/test"
)

func Test(t *testing.T) {
	tests := []struct {
		in      []byte
		charset string
		want    string
		wantErr string
	}{
		{[]byte("H€llo"), "what is this", "", "unknown character set: what is this"},
		{[]byte("€"[1:]), "utf-8", "��", ""},

		{[]byte(""), "us-ascii", "", ""},
		{[]byte("H€llo"), "utf-8", "H€llo", ""},
		{[]byte("H€llo"), "utf8", "H€llo", ""},
		{[]byte("hello"), "us-ascii", "hello", ""},
		{[]byte{0x48, 0xa4, 0x6c, 0x6c, 0x6f}, "iso-8859-15", "H€llo", ""},
		{[]byte{0x93, 0xfa, 0x96, 0x7b, 0x8c, 0xea}, "shift_jis", "日本語", ""},
		{[]byte{0xff, 0xfe, 0x48, 0x00, 0xac, 0x20, 0x6c, 0x00, 0x6c, 0x00, 0x6f, 0x00}, "utf-16", "H€llo", ""},
	}

	for _, tt := range tests {
		t.Run(tt.charset, func(t *testing.T) {
			out, err := Byte(tt.charset, tt.in)
			if !test.ErrorContains(err, tt.wantErr) {
				t.Fatalf("unexpected error\ngot:  %v\nwant: %v", err, tt.wantErr)
			}

			s := string(out)
			if s != tt.want {
				t.Errorf("wrong output\nout:  %#v\nwant: %#v", s, tt.want)
			}

			if !utf8.ValidString(s) {
				t.Errorf("not a valid string: %#v", s)
			}
		})
	}
}
