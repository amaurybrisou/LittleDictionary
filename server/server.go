package server

import (
  // "fmt"
  "log"
  "net/http"
  "LittleDictionary/config"
  "LittleDictionary/server/middlewares"
  "LittleDictionary/server/handlers"
)


func Server() {
  cfg := config.Get("config.default.gcfg")

  http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("server/public/"))))

  http.HandleFunc("/", handlers.IndexHandler)
  http.HandleFunc("/add", handlers.AddHandler)
  http.HandleFunc("/find", handlers.FindHandler)

  http.Handle("/word/add",
    middlewares.Mongo(
      http.HandlerFunc(handlers.AddWord), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName))

  http.Handle("/words/filter",
    middlewares.Mongo(
      http.HandlerFunc(handlers.FilterWords), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName))

  http.Handle("/words",
    middlewares.Mongo(
      http.HandlerFunc(handlers.FindWords), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName))

  log.Println("Http Server Listening on "+cfg.Http.Host+":"+cfg.Http.Port)
  log.Fatal(
    http.ListenAndServe(
      cfg.Http.Host+":"+cfg.Http.Port, nil))
  
}
