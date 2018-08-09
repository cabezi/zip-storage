package main

import (
	"io/ioutil"
	"os"

	"github.com/cabezi/zip-storage/config"
	"github.com/cabezi/zip-storage/controller"
	"github.com/zipper-project/zipper/common/log"
	yaml "gopkg.in/yaml.v2"
)

func main() {

	log.Println("zip-storge start...")

	yamlFile, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Errorf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, config.Cfg)
	if err != nil {
		log.Errorf("Unmarshal: %v", err)
	}

	log.Println("------>", config.Cfg)

	c := controller.NewControler(config.Cfg)

	c.ReceivePushHeight()

}
