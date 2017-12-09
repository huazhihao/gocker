package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Image struct {
	image string
	tag   string
}

func (im *Image) name() string {
	return im.image + ":" + im.tag
}

func (im *Image) path() string {
	return BtrfsPath + ImgPrefix + im.name()
}

func (im *Image) tmp() string {
	return "/tmp/" + im.name() + "/"
}

func (im *Image) layer() string {
	type Manifest struct {
		Config   string
		RepoTags []string
		Layers   []string
	}
	bs, _ := ioutil.ReadFile(im.tmp() + "manifest.json")
	var m []Manifest
	_ = json.Unmarshal(bs, &m)
	return im.tmp() + m[0].Layers[0]
}

func (im *Image) exists() bool {
	if im.image == "" || im.tag == "" {
		log.Fatal("image name or tag is not specified")
	}
	_, err := os.Stat(im.path())
	return err == nil
}

func Init(image, tag string) *Image {
	if tag == "" {
		tag = "latest"
	}
	return &Image{image, tag}
}

func Pull(image, tag string) {
	im := *Init(image, tag)
	if !im.exists() {
		panicRun("btrfs", "subvolume", "create", im.path())
	}
	panicRun("download-frozen-image-v2", im.tmp(), im.name())
	panicRun("tar", "-xf", im.layer(), "-C", im.path())
	_ = os.Remove(im.tmp())
}

func Commit(image, tag string) {
	//im := *Init(image, tag)
	//todo
}
func Rmi(image, tag string) {
	im := *Init(image, tag)
	if im.exists() {
		panicRun("btrfs", "subvolume", "delete", im.path())
	}
}
func Images() {
	files, err := ioutil.ReadDir(BtrfsPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), ImgPrefix) {
			fmt.Println(strings.TrimLeft(file.Name(), ImgPrefix))
		}
	}
}
