package server

import (
  "log"
  "net/http"
  "LittleDictionary/config"
  "LittleDictionary/server/middlewares"
  "LittleDictionary/server/handlers"
)



func Server() {
  cfg := config.Get("config.gcfg")

  http.Handle("/words/random",
    middlewares.Mongo(
      http.HandlerFunc(handlers.RandomWord), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName,
      cfg.Database.CollectionName ))
  
  auth := middlewares.NewBasicAuth(cfg.Auth.Username, cfg.Auth.Password)

  http.Handle("/public/", http.StripPrefix("/public/", auth.BasicAuthHandler(http.FileServer(http.Dir("server/public/")))))

  http.Handle("/", auth.BasicAuthHandler(http.HandlerFunc(handlers.IndexHandler)))
  http.Handle("/add", auth.BasicAuthHandler(http.HandlerFunc(handlers.AddHandler)))
  http.Handle("/find", auth.BasicAuthHandler(http.HandlerFunc(handlers.FindHandler)))

  http.Handle("/word/add",
    middlewares.Mongo(
      auth.BasicAuthHandler(http.HandlerFunc(handlers.AddWord)), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName,
      cfg.Database.CollectionName ))

  http.Handle("/word/delete",
    middlewares.Mongo(
      auth.BasicAuthHandler(http.HandlerFunc(handlers.DelWord)), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName,
      cfg.Database.CollectionName ))

  http.Handle("/word/update",
    middlewares.Mongo(
      auth.BasicAuthHandler(http.HandlerFunc(handlers.UpdateWord)), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName,
      cfg.Database.CollectionName ))

  http.Handle("/words/filter",
    middlewares.Mongo(
      auth.BasicAuthHandler(http.HandlerFunc(handlers.FilterWords)), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName,
      cfg.Database.CollectionName ))

  http.Handle("/words",
    middlewares.Mongo(
      auth.BasicAuthHandler(http.HandlerFunc(handlers.FindWords)), 
      cfg.Database.Url,
      cfg.Database.Port,
      cfg.Database.DbName,
      cfg.Database.CollectionName ))

  



  log.Println("Http Server Listening on "+cfg.Http.Host+":"+cfg.Http.Port)
  log.Fatal(http.ListenAndServe(cfg.Http.Host+":"+cfg.Http.Port, nil))
  
}
