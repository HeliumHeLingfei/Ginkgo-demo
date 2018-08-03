package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

/*
func main(){
	cli:=New("10.100.47.76","22","root","1234567890",5000)
	m:=cli.DoRun("1","echo hello!")
	m1:=cli.DoRun("2","echo hrello!")
	fmt.Print(m,m1)
}*/

func main() {
	clis, err := GetJsonFile("client.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	results := make([]chan map[string]string, len(clis.Clients))
	for i, cli := range clis.Clients {
		results[i] = make(chan map[string]string)
		go cli.Command(cli.CmdFile, results[i])
	}
	fp, err := os.OpenFile("result.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	m := make(map[string][]map[string]string)
	var out bytes.Buffer
	for i, result := range results {
		var t []map[string]string
		for {
			r, ok := <-result
			if !ok {
				break
			}
			t = append(t, r)
		}
		m[(clis.Clients[i].IP + ":" + clis.Clients[i].Port)] = t
	}
	data, _ := json.Marshal(m)
	json.Indent(&out, data, "", "\t")
	out.WriteTo(fp)
}
