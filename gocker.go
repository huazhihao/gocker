package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const BtrfsPath = "/var/gocker/"
const ImgPrefix = "img_"
const ContainerPrefix = "ps_"
const CpuShares = 512
const Memory = 512 * 100000

func help() {
	fmt.Println("help")
	//TODO
}

func panicRun(args ...string) []byte {
	if os.Getenv("DEBUG") != "" {
		fmt.Println(args)
	}
	out, err := exec.Command(args[0], args[1:]...).Output()
	if os.Getenv("DEBUG") != "" {
		log.Println(string(out))
	}
	if err != nil {
		log.Fatal(err)
	}
	return out
}
