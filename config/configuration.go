package config

import (
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	SearchUrl         string
	ItemNameNode      string
	ItemPriceNode     string
	ItemUrlNode       string
	NumberOfPagesNode string
}

var Conf Configuration

func init() {
	err := gonfig.GetConf("./config.json", &Conf)
	if err != nil {
		panic(err)
	}
}
