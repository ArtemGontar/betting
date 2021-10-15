package main

import (
	"flag"
	"log"

	"github.com/ArtemGontar/betting/internal/app/apiserver"
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
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}

	// расчет мат. ожидания
	//(Вероятность выигрыша) х (сумму потенциального выигрыша по текущему пари) – (вероятность проигрыша) х (сумму потенциального проигрыша по текущему пари).
}
