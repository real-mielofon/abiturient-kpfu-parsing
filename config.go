package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// ConfigType its write to config
type ConfigType struct {
	chats  map[int64]int
	status StatusAbiturienta
}

// readLines reads a whole file into memory
// and returns a slice of its lines.

// ReadConfig read config
func (c *ConfigType) ReadConfig() error {
	c.chats = make(map[int64]int)

	file, err := os.Open(fileConfig)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return scanner.Err()
	}
	n, err := strconv.ParseInt(scanner.Text(), 10, 32)
	if err != nil {
		return err
	}
	len := int(n)
	for i := 0; i < len; i++ {
		if !scanner.Scan() {
			return scanner.Err()
		}
		n, err := strconv.ParseInt(scanner.Text(), 10, 32)
		if err != nil {
			return err
		}
		c.chats[n] = 0
	}
	if !scanner.Scan() {
		return scanner.Err()
	}
	n, err = strconv.ParseInt(scanner.Text(), 10, 32)
	if err != nil {
		return err
	}
	c.status.Num = int(n)
	if !scanner.Scan() {
		return scanner.Err()
	}
	n, err = strconv.ParseInt(scanner.Text(), 10, 32)
	if err != nil {
		return err
	}
	c.status.NumWithOriginal = int(n)
	return nil

}

// WriteConfig write config
func (c *ConfigType) WriteConfig() error {

	file, _ := os.OpenFile(fileConfig, os.O_WRONLY, 0666)
	//    if err != nil {
	//log.Fatal(err)
	//return err
	//}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%d\n", len(c.chats))
	for key := range c.chats {
		fmt.Fprintf(w, "%d\n", key)
	}
	fmt.Fprintf(w, "%d\n", c.status.Num)
	fmt.Fprintf(w, "%d\n", c.status.NumWithOriginal)
	w.Flush()
	return nil
}

// Add - add key
func (c *ConfigType) Add(key int64) {
	c.chats[key] = 0
}
