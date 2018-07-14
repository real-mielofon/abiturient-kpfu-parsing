package app

import (
	"reflect"
	"testing"
)

func TestListAbiturents_Get(t *testing.T) {
	tests := []struct {
		name    string
		l       ListAbiturents
		want    []Abiturient
		wantErr bool
	}{{name: "test1", l: ListAbiturents{}, want: []Abiturient{}, wantErr: false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAbiturents.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListAbiturents.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
