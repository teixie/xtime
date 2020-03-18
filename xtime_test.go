package xtime

import (
	"testing"
)

func TestParse(t *testing.T) {
	for _, test := range []struct {
		expr interface{}
		want string
	}{
		{"2020-01-01 00:00:01", "2020-01-01 00:00:01"},
	} {
		got, err := Parse(test.expr)
		if err != nil {
			t.Errorf("%v", err)
		}
		if got.Format(LayoutYmdHis) != test.want {
			t.Errorf("Parse(%s) = %v, want %v", test.expr, got.Format(LayoutYmdHis), test.want)
		}
	}
}
