package config

import (
	"reflect"
	"testing"

	"github.com/real-mielofon/abiturient-kpfu-parsing/status"
)

const (
	testConfigFile = "./data/test_subscribe.txt"
)

func TestConfigType_WriteConfig(t *testing.T) {
	tests := []struct {
		name    string
		c       *Config
		wantErr bool
	}{
		{name: "test1",
			c: &Config{Chats: map[int64]status.StatusByName{
				0: status.StatusByName{Name: "Пономарев Степан Алексеевич", Status: status.StatusAbiturienta{30, 6}},
				1: status.StatusByName{Name: "Иванов Иван Иванович", Status: status.StatusAbiturienta{0, 1}},
			},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.configFileName = testConfigFile
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
			name: "test1",
			c: &Config{
				Chats: map[int64]status.StatusByName{
					0: {Name: "Пономарёв Степан Алексеевич", Status: status.StatusAbiturienta{30, 6}},
					1: {Name: "Иванов Иван Иванович", Status: status.StatusAbiturienta{10, 5}},
				},
			},
			want: &Config{
				Chats: map[int64]status.StatusByName{
					0: {Name: "Пономарёв Степан Алексеевич", Status: status.StatusAbiturienta{30, 6}},
					1: {Name: "Иванов Иван Иванович", Status: status.StatusAbiturienta{10, 5}},
				},
			},
			wantErr: false,
		},
		{
			name: "test_empty",
			c: &Config{
				Chats: map[int64]status.StatusByName{},
			},
			want: &Config{
				Chats: map[int64]status.StatusByName{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.want.configFileName = testConfigFile
			if errWrite := tt.want.WriteConfig(); (errWrite != nil) != tt.wantErr {
				t.Errorf("ConfigType.ReadConfig() error = %v, wantErr %v", errWrite, tt.wantErr)
			}
			if err := tt.c.ReadConfig(testConfigFile); (err != nil) != tt.wantErr {
				t.Errorf("ConfigType.ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.want, tt.c) {
				t.Errorf("getListAbiturient() = %v, want %v", tt.c, tt.want)
			}
		})
	}
}
