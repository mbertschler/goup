package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"git.exahome.net/tools/elephant"
)

type Config struct {
	OS         string
	ARCH       string
	Env        string // GOARM=7
	IP         string
	Port       string
	SSHUser    string
	RemotePath string // /root/bin/
}

// IP of Beaglebone
const IP = "10.0.0.10"

// Port of SSH on BeagleBone
const Port = "8022"

func main() {
	var c Config
	err := elephant.Load("goup.config", &c)
	if err != nil {
		fmt.Println(err)
		err := elephant.Store("goup.config", &c)
		if err != nil {
			fmt.Println(err)
			elephant.Delete("goup.config")
		}
		err = elephant.Store("goup.config", &c)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("please configure goup in elephant and then restart")
		os.Exit(1)
	}
	fmt.Printf("Config: %#v", c)
	flag.Parse()
	args := flag.Args()
	//fmt.Println("os.Args:", os.Args)
	if len(args) < 2 {
		fmt.Println("need at least 2 arguments (eg. goarm install package)")
		return
	}
	if args[0] == "install" {
		cmd := exec.Command("env", "GOOS="+c.OS, "GOARCH="+c.ARCH, c.Env, "go", "install", args[1])
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("go install error:", err)
			return
		}
		fmt.Print("built... ")
		argparts := strings.Split(args[1], "/")
		lastElement := argparts[len(argparts)-1]
		cmd = exec.Command("scp", "-P "+c.Port, os.ExpandEnv("$GOPATH/bin/"+c.OS+"_"+c.ARCH+"/")+lastElement,
			c.SSHUser+"@"+c.IP+":"+c.RemotePath+lastElement)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("scp error:", err)
			return
		}
		fmt.Println("uploaded to remote")
	} else {
		fmt.Println("Command", args[1], "not implemented")
	}
}
