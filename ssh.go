package main

import (
	"fmt"
	"io/ioutil"

	"regexp"

	"golang.org/x/crypto/ssh"

	"time"
)

type Clients struct {
	Clients []Cli
}

type Cli struct {
	IP         string //IP地址
	Port       string //端口号
	Username   string //用户名
	Password   string //密码
	Timeout    int
	LastResult chan string //最近一次Run的结果
	CmdFile    string
	Cmds       []string
}

func New(ip string, port string, username string, password string, timeout int) Cli {
	var cli Cli
	cli.IP = ip
	cli.Password = password
	cli.Port = port
	cli.Username = username
	cli.Timeout = timeout
	return cli
}

func (cli Cli) connect() (*ssh.Session, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", cli.IP, cli.Port), &ssh.ClientConfig{
		User:            cli.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(cli.Password)},
		Timeout:         3000 * time.Millisecond,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, err
}

func (cli Cli) start(shell string, sshchannel chan string, ch chan map[string]string) {
	starttime := time.Now()
	session, err := cli.connect()
	if err != nil {
		sshResult := fmt.Sprintf("<%s>", err.Error())
		sshchannel <- sshResult
		close(sshchannel)
		return
	}
	defer session.Close()
	buf, err := session.CombinedOutput("source /root/localrc;" + shell + ";echo return status: $?")
	gettime := time.Since(starttime).String()
	result := string(buf)
	re, _ := regexp.Compile("[0-9]+")
	rs := re.Find([]byte(result[len(result)-5:]))
	m := map[string]string{
		"command":       shell,
		"return status": string(rs),
		"time":          gettime,
	}
	ch <- m

	sshchannel <- fmt.Sprintf("%s\n%s\n%v%s%v\n%s\n", ("IP: " + cli.IP + ":" + cli.Port), ("command: " + shell + "\nresult: "), result, "using: ", gettime, "------------------------------------------------------------------------------------------")
	return
}

//执行shell
func (cli Cli) Run(shell string, ch chan map[string]string) {
	sshchannel := make(chan string)
	go cli.start(shell, sshchannel, ch)
	select {
	case <-time.After(time.Duration(cli.Timeout) * time.Millisecond):
		cli.LastResult <- fmt.Sprintf("%s\n%s%v%s\n%s\n", ("IP: " + cli.IP + ":" + cli.Port), "timeout: ", float64(cli.Timeout)/1000, " second.", "------------------------------------------------------------------------------------------")
	case result, ok := <-sshchannel:
		if !ok {
			cli.LastResult <- "error: " + result
		} else {
			cli.LastResult <- result
		}
	}
	return
}

func (cli Cli) DoRun(filename string, shell string) map[string]string {
	sshchannel := make(chan string)
	ch := make(chan map[string]string, 1)
	go cli.start(shell, sshchannel, ch)

	select {
	case <-time.After(time.Duration(cli.Timeout) * time.Millisecond):
		ioutil.WriteFile(filename, []byte(fmt.Sprintf("%s\n%s%v%s\n%s\n", ("IP: "+cli.IP+":"+cli.Port), "timeout: ", float64(cli.Timeout)/1000, " second.", "------------------------------------------------------------------------------------------")), 0644)
	case result, ok := <-sshchannel:
		if !ok {
			ioutil.WriteFile(filename, []byte("error: "+result), 0644)
		} else {
			ioutil.WriteFile(filename, []byte(result), 0644)
			return <-ch
		}
	}
	return nil

}

func (cli Cli) Command(filename string, result chan map[string]string) {
	defer close(result)
	cmds, err := GetFile(filename)
	if err != nil {
		result <- map[string]string{
			"error": "wrong file name",
		}
		return
	}
	cli.LastResult = make(chan string, len(cmds))
	for i, cmd := range cmds {
		fmt.Println(i, cmd)
		go cli.Run(cmd, result)
	}
	for i := 0; i < len(cmds); i++ {
		fmt.Print(<-cli.LastResult)
	}
}
