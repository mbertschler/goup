package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
)

type UpTarget struct {
	Host   string
	Port   int
	User   string
	Target string
	Os     string
	Arch   string
}

func main() {
	var targets = map[string]UpTarget{}
	file, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), "goup.toml"))
	if err != nil {
		log.Fatal(err)
	}
	err = toml.Unmarshal(file, &targets)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) < 2 {
		fmt.Println("need at least 3 arguments (goup [target] [package])")
		return
	}

	conf, ok := targets[os.Args[1]]
	if !ok {
		fmt.Println("can't find target", os.Args[1], "in config file")
	}
	err = build(conf, os.Args[2])
	if err != nil {
		fmt.Println("go install error:", err)
		return
	}
	err = upload(conf, os.Args[2])
	if err != nil {
		fmt.Println("scp error:", err)
	}
}

func build(conf UpTarget, pkg string) error {
	cmd := exec.Command("go", "install", pkg)
	cmd.Env = append(os.Environ(), "GOARCH="+conf.Arch, "GOOS="+conf.Os)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func upload(conf UpTarget, pkg string) error {
	parts := strings.Split(pkg, "/")
	binary := parts[len(parts)-1]
	cmd := exec.Command("scp", "-P "+fmt.Sprint(conf.Port), findBin(conf, binary),
		conf.User+"@"+conf.Host+":"+filepath.Join(conf.Target, binary))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func findBin(conf UpTarget, name string) string {
	gopaths := filepath.SplitList(os.Getenv("GOPATH"))
	for _, p := range gopaths {
		var path string
		if runtime.GOOS == conf.Os && runtime.GOARCH == conf.Arch {
			path = filepath.Join(p, "bin", name)
		} else {
			path = filepath.Join(p, "bin", conf.Os+"_"+conf.Arch, name)
		}

		f, err := os.Open(path)
		if err == nil {
			f.Close()
			return path
		}
		if !os.IsNotExist(err) {
			fmt.Println("unexpected open error:", err)
		}
	}
	return ""
}
