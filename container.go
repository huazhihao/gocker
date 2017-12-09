package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Container struct {
	id int
}

func (c *Container) path() string {
	return fmt.Sprintf("%s%s%d/", BtrfsPath, ContainerPrefix, c.id)
}

func (c *Container) exists() bool {
	_, err := os.Stat(c.path())
	return err == nil
}

func (c *Container) namespace() string {
	return fmt.Sprintf("%d", c.id)
}

func (c *Container) cgroups() string {
	return "cpu,cpuacct,memory:/" + c.namespace()
}

func (c *Container) attr(s string) string {
	switch s {
	case "mac":
		return fmt.Sprintf("02:42:ac:11:0%d:%d", (c.id-42000)/100, c.id%100)
	case "ip":
		return fmt.Sprintf("10.0.0.%d/24", c.id-42000)
	}
	return fmt.Sprintf("%s_%d", s, c.id)
}

func Rm(sid string) {
	id, err := strconv.Atoi(sid)
	if err != nil {
		return
	}
	c := &Container{id}
	if c.exists() {
		panicRun("btrfs", "subvolume", "delete", c.path())
		panicRun("cgdelete", "-g", c.cgroups())
	}
}
func Ps() {
	files, err := ioutil.ReadDir(BtrfsPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("CONTAINER\tCOMMAND")
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), ContainerPrefix) {
			out, _ := ioutil.ReadFile(BtrfsPath + file.Name() + "/CMD")
			fmt.Println(strings.TrimLeft(file.Name(), ContainerPrefix) + "\t" + string(out))
		}
	}
}
func Run(image, tag string, cmds ...string) {
	cmd := strings.Join(cmds, " ")

	rand.Seed(time.Now().UnixNano())
	c := &Container{42002 + rand.Intn(252)} //42002-42254
	im := *Init(image, tag)
	if !im.exists() {
		Pull(image, tag)
	}
	panicRun("ip", "link", "add", "dev", c.attr("veth0"), "type", "veth", "peer", "name", c.attr("veth1"))
	panicRun("ip", "link", "set", "dev", c.attr("veth0"), "up")
	panicRun("ip", "link", "set", c.attr("veth0"), "master", "bridge0")
	panicRun("ip", "netns", "add", c.attr("netns"))
	panicRun("ip", "link", "set", c.attr("veth1"), "netns", c.attr("netns"))
	panicRun("ip", "netns", "exec", c.attr("netns"), "ip", "link", "set", "dev", "lo", "up")
	panicRun("ip", "netns", "exec", c.attr("netns"), "ip", "link", "set", c.attr("veth1"), "address", c.attr("mac"))
	panicRun("ip", "netns", "exec", c.attr("netns"), "ip", "addr", "add", c.attr("ip"), "dev", c.attr("veth1"))
	panicRun("ip", "netns", "exec", c.attr("netns"), "ip", "link", "set", "dev", c.attr("veth1"), "up")
	panicRun("ip", "netns", "exec", c.attr("netns"), "ip", "route", "add", "default", "via", "10.0.0.1")
	panicRun("btrfs", "subvolume", "snapshot", im.path(), c.path())

	ioutil.WriteFile(c.path()+"/etc/resolv.conf", []byte("nameserver 8.8.8.8"), 0644)
	ioutil.WriteFile(c.path()+"/CMD", []byte(cmd), 0644)

	panicRun("cgcreate", "-g", c.cgroups())
	panicRun("cgset", "-r", fmt.Sprintf("cpu.shares=%d", CpuShares), c.namespace())
	panicRun("cgset", "-r", fmt.Sprintf("memory.limit_in_bytes=%d", Memory), c.namespace())
	out := panicRun("cgexec", "-g", c.cgroups(),
		"ip", "netns", "exec", c.attr("netns"),
		"unshare", "-fmuip", "--mount-proc",
		"chroot", c.path(),
		"/bin/sh", "-c", "/bin/mount -t proc proc /proc && "+cmd)
	fmt.Println(string(out))
	ioutil.WriteFile(c.path()+"/LOG", out, 0644)
	panicRun("ip", "link", "del", "dev", c.attr("veth0"))
	panicRun("ip", "netns", "del", c.attr("netns"))
}
func Exec(sid string, cmds ...string) {
	id, err := strconv.Atoi(sid)
	if err != nil {
		return
	}
	_ = &Container{id}

	//TODO
}
func Logs(sid string) {
	id, err := strconv.Atoi(sid)
	if err != nil {
		return
	}
	c := &Container{id}
	if c.exists() {
		out, _ := ioutil.ReadFile(c.path() + "/LOG")
		fmt.Println(string(out))
	}
}
