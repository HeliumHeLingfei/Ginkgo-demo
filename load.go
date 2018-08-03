package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unsafe"
)


func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return *(*string)(unsafe.Pointer(&ln)), err
}

func Load(filename string) (*bufio.Reader, error) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(filename)
		return nil, err
	}
	buf := bufio.NewReader(f)
	return buf, nil
}

func GetFile(filename string) ([]string, error) {
	var (
		lines []string
	)
	file, err := Load(filename)
	if file == nil {
		return nil, err
	}
	line, readErr := Readln(file)
	for readErr == nil {
		lines = append(lines, line)
		line, readErr = Readln(file)
	}
	if readErr != io.EOF {
		return nil, readErr
	}
	return lines, nil
}

func GetJsonFile(filename string) (Clients, error) {
	result := Clients{}
	f, err := ioutil.ReadFile(filename)
	if  err != nil {
		fmt.Println("read file ", filename, err)
		return result, err
	}
	err = json.Unmarshal(f, &result)
	if err != nil {
		fmt.Println("read file ", filename, err)
		return result, err
	}
	return result, nil
}
