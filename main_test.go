package main

import (
	"fmt"
	"strconv"
	"testing"
)

type testStruct struct {
	name    string
	want    []Abiturient
	wantErr bool
}

func Test_getListAbiturient(t *testing.T) {
	tests := []testStruct{{"Test1", nil, false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getListAbiturient()
			for _, a := range got {
				fmt.Printf("%4d %40s %3d %s\n", a.Num, a.Fio, a.Points[4], strconv.FormatBool(a.Original))
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("getListAbiturient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//			if !reflect.DeepEqual(got, tt.want) {
			//				t.Errorf("getListAbiturient() = %v, want %v", got, tt.want)
			//			}
		})
	}
}
