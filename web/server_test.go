package web

import (
	"net/url"
	"testing"
)

func Test_countFromQuery(t *testing.T) {
	type args struct {
		query url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"count exists", args{map[string][]string{"count": []string{"5"}}}, 5, false},
		{"non-number count", args{map[string][]string{"count": []string{"hello"}}}, 0, true},
		{"count not present, default used", args{map[string][]string{}}, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := countFromQuery(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("countFromQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("countFromQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
