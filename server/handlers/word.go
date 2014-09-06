package handlers

import (
	"log"
	"net/http"
  "encoding/json"
  "io"
  "fmt"
  "LittleDictionary/server/middlewares"
  "gopkg.in/mgo.v2/bson"
  "math/rand"
  "time"
)

type Word struct {
  Word string
  Pos string
  Definition string 
  Language string
  Example string `json:",omitempty"`
  Ethymology string `json:",omitempty"`
}


func getBody(body io.Reader) (t Word, ) {
  decoder := json.NewDecoder(body)
  err := decoder.Decode(&t)
  if err != nil {
    log.Println("Error deconding request Body ")
  }
  return t
}

type Response struct {
  Result bool
  Words []Word `json:",omitempty"`
}

func AddWord(w http.ResponseWriter, r *http.Request){
  rObj := getBody(r.Body)
  c := middlewares.GetWords(r)
  

  emptyWord := Word{}
  
  if rObj == emptyWord {
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
    http.Error(w, jsonMsg, http.StatusInternalServerError)
  }

  err := c.Insert(&Word{
    rObj.Word, 
    rObj.Pos, 
    rObj.Definition, 
    rObj.Language,
    rObj.Example, 
    rObj.Ethymology })
  if err != nil {
    log.Println(err)

    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
    http.Error(w, jsonMsg, http.StatusInternalServerError)
  } else {
    jsonMsg, _ := forgeResponse(Response{true, nil}) 
    log.Println("Word added")
    fmt.Fprintf(w, jsonMsg) 
  }

  
  
}

func forgeResponse(rep Response) (string, error){
  jbMsg, err := json.Marshal(rep)

  if err != nil {    
    return "", err
  }

  jsonMsg := string(jbMsg[:]) // converting byte array to string
  return jsonMsg, nil
}

func buildQuery(rObj *Word) (q []bson.M){
  if rObj.Word != "" {
    q = append(q, bson.M{ "word": rObj.Word })
  }

  if rObj.Pos != "" {
    q = append(q, bson.M{ "pos": rObj.Pos })
  }

  if rObj.Definition != "" {
    q = append(q, bson.M{ "definition": rObj.Definition })
  }

  if rObj.Language != "" {
    q = append(q, bson.M{ "language": rObj.Language })
  }

  if rObj.Example != "" {
    q = append(q, bson.M{ "example": rObj.Example })
  }

  if rObj.Ethymology != "" {
    q = append(q, bson.M{ "ethymology": rObj.Ethymology })
  }
  return q
}

func FilterWords(w http.ResponseWriter, r *http.Request){
  rObj := getBody(r.Body)
  c := middlewares.GetWords(r)

  q := buildQuery(&rObj)


  if len(q) == 0 {
    FindWords(w, r)
    return;
  }

  words := []Word{}
  err := c.Find(bson.M{ "$or": q }).All(&words)

  if err != nil {
    log.Println("Error fetching : ", err)
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
    http.Error(w, jsonMsg, http.StatusInternalServerError)
    return
  } else {

    if len(words) != 0 {
      jsonMsg, _ := forgeResponse(Response{true, words})

      fmt.Fprintf(w, jsonMsg)
      return 
    }
  }

  jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
  http.Error(w, jsonMsg, http.StatusInternalServerError)
  return
}

func FindWords(w http.ResponseWriter, r *http.Request){
  c := middlewares.GetWords(r)


  words := []Word{} 

  err := c.Find(bson.M{}).All(&words)

  if err != nil {
    log.Println("Error fetching : ", err)
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
    http.Error(w, jsonMsg, http.StatusInternalServerError)
    return
  } else {

    jsonMsg, _ := forgeResponse(Response{true, words})

    fmt.Fprintf(w, jsonMsg)
    return 
  }

  jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
  http.Error(w, jsonMsg, http.StatusInternalServerError)
  return 
}


func RandomWord(w http.ResponseWriter, r *http.Request){
  c := middlewares.GetWords(r)

  words := []Word{} 

  err := c.Find(bson.M{}).All(&words)

  if err != nil {
    log.Println("Error fetching : ", err)
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
    http.Error(w, jsonMsg, http.StatusInternalServerError)
    return
  } else {

    rand.Seed(time.Now().UnixNano())
    index := rand.Intn(len(words))

    log.Println(index)

    jsonMsg, _ := forgeResponse(Response{true, words[index:index+1]})

    fmt.Fprintf(w, jsonMsg)
    return 
  }

  jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
  http.Error(w, jsonMsg, http.StatusInternalServerError)
  return 
}