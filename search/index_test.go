package search

import (
	"reflect"
	"testing"
)

func Test_handleMultipleCommandsOnSameLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"one command", args{"grep(5) - This is a test"}, []string{"grep(5) - This is a test"}},
		{"two commands", args{"grep(5), egrep(6) - This is a test"}, []string{"grep(5) - This is a test", "egrep(6) - This is a test"}},
		{"comma in description", args{"grep(5), egrep(6) - This is, a test"}, []string{"grep(5) - This is, a test", "egrep(6) - This is, a test"}},
		{"extra spaces", args{"grep(5),   egrep(6)    -   This is, a test"}, []string{"grep(5) -   This is, a test", "egrep(6) -   This is, a test"}},
		{"graceful failure", args{"grep(5),egrep(6)"}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleMultipleCommandsOnSameLine(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMultipleCommandsOnSameLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getManpages(t *testing.T) {
	res := getManpages()
	if len(res) <= 0 {
		t.Fail()
	}
}

func Test_getManpage(t *testing.T) {
	t.Run("has manpage", func(t *testing.T) {
		res := getManpage("grep")
		if len(res) <= 0 {
			t.Fail()
		}
	})

	t.Run("no manpage", func(t *testing.T) {
		res := getManpage("fake_manpage")
		if len(res) > 0 {
			t.Fail()
		}
	})

	t.Run("lowercase needed", func(t *testing.T) {
		res := getManpage("Grep")
		if len(res) <= 0 {
			t.Fail()
		}
	})

	t.Run("alias to same manpage", func(t *testing.T) {
		first := getManpage("grep")
		second := getManpage("egrep")
		if first != second {
			t.Fail()
		}
	})
}
