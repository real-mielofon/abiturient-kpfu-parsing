package main

import (
	"reflect"
	"testing"
)

func TestConfigType_WriteConfig(t *testing.T) {
	tests := []struct {
		name    string
		c       *ConfigType
		wantErr bool
	}{
		{
			name:    "test1",
			c:       &ConfigType{chats: map[int64]int{0: 0, 1: 0}, status: StatusAbiturienta{30, 6}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.WriteConfig(); (err != nil) != tt.wantErr {
				t.Errorf("ConfigType.WriteConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigType_ReadConfig(t *testing.T) {
	tests := []struct {
		name    string
		c       *ConfigType
		want    *ConfigType
		wantErr bool
	}{
		{
			name:    "test1",
			c:       &ConfigType{chats: map[int64]int{}, status: StatusAbiturienta{0, 0}},
			want:    &ConfigType{chats: map[int64]int{0: 0, 1: 0}, status: StatusAbiturienta{30, 6}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if errWrite := tt.want.WriteConfig(); (errWrite != nil) != tt.wantErr {
				t.Errorf("ConfigType.ReadConfig() error = %v, wantErr %v", errWrite, tt.wantErr)
			}
			if err := tt.c.ReadConfig(); (err != nil) != tt.wantErr {
				t.Errorf("ConfigType.ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.want, tt.c) {
				t.Errorf("getListAbiturient() = %v, want %v", tt.c, tt.want)
			}
		})
	}
}
