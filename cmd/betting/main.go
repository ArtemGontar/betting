package main

import (
	"flag"
	"log"

	"github.com/ArtemGontar/betting/internal/app/apiserver"
	premierleague_reader "github.com/ArtemGontar/betting/internal/app/reader/premierleague"
	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	//premierleague_reader.Teams()
	premierleague_reader.Players()
}

func apiServer() {
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
