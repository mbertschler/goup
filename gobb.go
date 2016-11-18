package main

import (
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
	//fmt.Println("os.Args:", os.Args)
	if len(os.Args) < 3 {
		fmt.Println("need at least 2 arguments (eg. gobb install package)")
		return
	}
	if os.Args[1] == "install" {
		cmd := exec.Command("env", "GOOS=linux", "GOARCH=arm", "GOARM=7", "go", "install", os.Args[2])
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("go install error:", err)
			return
		}
		fmt.Print("built... ")
		argparts := strings.Split(os.Args[2], "/")
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
		fmt.Println("uploaded to BeagleBone")
	} else {
		fmt.Println("Command", os.Args[1])
	}
}
