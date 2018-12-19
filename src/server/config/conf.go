package config


import (	
	"gopkg.in/gcfg.v1"
	"log"
)

type Config struct {
  Database struct {
    Url string
    Port string
    DbName string
    CollectionName string
  }
  Http struct {
  	Host string 
  	Port string
  }
  Auth struct {
    Username string 
    Password string
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