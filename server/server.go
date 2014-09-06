package server

import (
  "log"
  "net/http"
  "LittleDictionary/config"
  "LittleDictionary/server/middlewares"
  "LittleDictionary/server/handlers"
  "github.com/icub3d/httpauth"
)



func Server() {
  cfg := config.Get("config.gcfg")

  http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("server/public/"))))

  http.HandleFunc("/", handlers.IndexHandler)
  http.HandleFunc("/add", handlers.AddHandler)
  http.HandleFunc("/find", handlers.FindHandler)

  http.Handle("/word/add",
    middlewares.Mongo(
      http.HandlerFunc(handlers.AddWord), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName,
      cfg.Database.CollectionName ))

  http.Handle("/words/filter",
    middlewares.Mongo(
      http.HandlerFunc(handlers.FilterWords), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName,
      cfg.Database.CollectionName ))

  http.Handle("/words",
    middlewares.Mongo(
      http.HandlerFunc(handlers.FindWords), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName,
      cfg.Database.CollectionName ))

  http.Handle("/words/random",
    middlewares.Mongo(
      http.HandlerFunc(handlers.RandomWord), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName,
      cfg.Database.CollectionName ))



  log.Println("Http Server Listening on "+cfg.Http.Host+":"+cfg.Http.Port)
  log.Fatal(
    http.ListenAndServe(
      cfg.Http.Host+":"+cfg.Http.Port, 
      httpauth.Basic("private area", http.DefaultServeMux,
      func(user, pass string) bool {
          if user == cfg.Auth.Username && pass == cfg.Auth.Password {
              return true
          }
          return false
      })))
  
}
