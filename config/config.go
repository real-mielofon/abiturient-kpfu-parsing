package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/real-mielofon/abiturient-kpfu-parsing/status"
)

// ConfigType its write to config
type Config struct {
	configFileName string
	Chats          map[int64]status.StatusByName
}

// readLines reads a whole file into memory
// and returns a slice of its lines.

// ReadConfig read config
func (c *Config) ReadConfig(fileConfig string) error {
	c.Chats = make(map[int64]status.StatusByName)

	file, err := os.Open(fileConfig)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	c.configFileName = fileConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	fmt.Println(c) // output: [UserA, UserB]
	return nil
}

// WriteConfig write config
func (c *Config) WriteConfig() error {

	file, err := os.OpenFile(c.configFileName, os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&c)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	return nil
}

func (c *Config) Add(key int64, name string) {
	c.Chats[key] = status.StatusByName{Name: name, Status: status.StatusAbiturienta{Num: 0, NumWithOriginal: 0}}
}
