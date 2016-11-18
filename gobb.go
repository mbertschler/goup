package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// IP of Beaglebone
const IP = "10.0.0.10"

// Port of SSH on BeagleBone
const Port = "8022"

func main() {
	flag.Parse()
	args := flag.Args()
	//fmt.Println("os.Args:", os.Args)
	if len(args) < 2 {
		fmt.Println("need at least 2 arguments (eg. goarm install package)")
		return
	}
	if args[0] == "install" {
		cmd := exec.Command("env", "GOOS=linux", "GOARCH=arm", "GOARM=7", "go", "install", args[1])
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
		cmd = exec.Command("scp", "-P "+Port, os.ExpandEnv("$GOPATH/bin/linux_arm/")+lastElement,
			"root@"+IP+":/root/bin/"+lastElement)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println("scp error:", err)
			return
		}
		fmt.Println("uploaded to ARM Linux")
	} else {
		fmt.Println("Command", args[1], "not implemented")
	}
}
