package main

import (
	"reflect"
	"testing"
)

func TestConfigType_WriteConfig(t *testing.T) {
	tests := []struct {
		name    string
		c       *Config
		wantErr bool
	}{
		name: "test1",
		c: &Config{chats: map[int64]StatusByName{
			0: StatusByName{Name: "Пономарев Степан Алексеевич", status: StatusAbiturienta{30, 6}},
			1: StatusByName{Name: "Иванов Иван Иванович", status: StatusAbiturienta{0, 1}},
		},
		},
		wantErr: false,
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
		c       *Config
		want    *Config
		wantErr bool
	}{
		{
			name: "test2",
			c:    &Config{chats: map[int64]int{}, status: StatusAbiturienta{0, 0}},
			want: &Config{
				chats: map[int64]StatusByName{
					0: {name: "Пономарёв Степан Алексеевич", status: StatusAbiturienta{30, 6}},
					1: {name: "Иванов Иван Иванович", status: StatusAbiturienta{10, 5}},
				},
			},
			wantErr: false,
		},
		{
			name: "test0",
			c:    &Config{chats: map[int64]int{}, status: StatusAbiturienta{0, 0}},
			want: &Config{
				chats: map[int64]StatusByName{},
			},
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
