package main_test

import (
	"fmt"
	. "gotest/ssh"
	"io/ioutil"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ssh", func() {
	var (
		clis     Clients
		c        Cli
		err      error
		filename string
		cmd      string
		rs       string
		time     float64
		serverid int
		cmdid    int
		try      int
		result   map[string]string
	)

	gettime := func(t string) (k float64) {
		if t[len(t)-2] == 'm' {
			k, _ = strconv.ParseFloat(t[:len(t)-2], 64)
			k = k / 1000
		} else {
			k, _ = strconv.ParseFloat(t[:len(t)-1], 64)
		}
		return
	}

	BeforeSuite(func() {
		clis, err = GetJsonFile("client.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		for i, _ := range clis.Clients {
			cmds, err := GetFile(clis.Clients[i].CmdFile)
			if err != nil {
				fmt.Println("error: wrong file name")
				Panic()
			}
			clis.Clients[i].Cmds = cmds
		}
		try = 0
	})

	AfterEach(func() {
		_, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Print(err)
		}
		//str := string(b)
		//fmt.Println(str)
	})

	BeforeEach(func() {
		serverid = 0
		cmdid = 0
		rs = "0"
		time = 1.7
	})

	JustBeforeEach(func() {
		c = clis.Clients[serverid]
		filename = strconv.Itoa(serverid) + "_" + strconv.Itoa(cmdid)
		cmd = c.Cmds[cmdid]
		result = c.DoRun(filename, cmd)
	})

	Describe("Connecting to remote server and get basic info", func() {
		Context("The first server", func() {
			It("Should get version", func() {
				Expect(result["command"]).To(Equal(cmd),"should be cephmgmtclient list-soft-version")

				Expect(result["return status"]).To(Equal(rs),"should return 0")

				Expect(gettime(result["time"])).To(BeNumerically("<", time),"should be less than 1.7s")
			})
		})

		Context("The second server", func() {
			BeforeEach(func() {
				serverid = 1
				cmdid = 0
				time = 0.1
				try = try + 1
			})

			It("Should be more than 0.1s", func() {
				Expect(gettime(result["time"])).To(BeNumerically(">", time))
			})
			It("Should echo twice", func() {
				Expect(try).To(Equal(2))
			})
		})
	})
})
