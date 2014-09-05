package config


import (	
	"code.google.com/p/gcfg"
	"log"
)

type Config struct {
  Database struct {
    Url string
    Port string
    DbName string
  }
  Http struct {
  	Host string 
  	Port string
  }
}

func Get(file string)(Config){
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, file)
	if err != nil {
		log.Fatal("error ", err)
	}
	return cfg
}