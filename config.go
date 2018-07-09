package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// ConfigType its write to config
type Config struct {
	chats map[int64]StatusByName
}

// readLines reads a whole file into memory
// and returns a slice of its lines.

// ReadConfig read config
func (c *Config) ReadConfig() error {
	c.chats = make(map[int64]StatusByName)

	file, err := os.Open(fileConfig)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

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

	file, _ := os.OpenFile(fileConfig, os.O_WRONLY, 0666)
	//    if err != nil {
	//log.Fatal(err)
	//return err
	//}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err := encoder.Encode(&c)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	return nil
}

func (c *Config) Add(key int64, name string) {
	c.chats[key] = StatusByName{Name: name, Status: StatusAbiturienta{Num:0, NumWithOriginal:0}}
}
