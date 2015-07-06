package handlers

import (
	"log"
	"net/http"
  "encoding/json"
  "io"
  "fmt"
  "github.com/goland-amaurybrisou/LittleDictionary/server/middlewares"
  "gopkg.in/mgo.v2/bson"
  "gopkg.in/mgo.v2"
  "math/rand"
  "time"
  "html/template"
  "errors"
  "strings"
)

type Word struct {
  Id bson.ObjectId `_id,omitempty`
  Word string
  Pos string
  Definition string 
  Language string
  Example string `json:",omitempty"`
  Ethymology string `json:",omitempty"`
}

type UpdatedWord struct {
  Id bson.ObjectId `_id,omitempty`
  Word string `json:",omitempty"`
  Pos string `json:",omitempty"`
  Definition string `json:",omitempty"`
  Language string `json:",omitempty"`
  Example string `json:",omitempty"`
  Ethymology string `json:",omitempty"`
}


func AddWord(w http.ResponseWriter, r *http.Request){
  rObj, err := getBody(r.Body)
  
  if err != nil {
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    http.Error(w, jsonMsg, http.StatusOK)
    return
  }

  rObj.Id = bson.NewObjectId()

  c := middlewares.GetWords(r)

  err = c.Insert(&Word{
    rObj.Id,
    rObj.Word, 
    rObj.Pos, 
    rObj.Definition, 
    rObj.Language,
    rObj.Example, 
    rObj.Ethymology })
  if err != nil {
    log.Println(err)
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    http.Error(w, jsonMsg, http.StatusNoContent)
  } else {
    jsonMsg, _ := forgeResponse(Response{true, nil}) 
    log.Println("Word added")
    fmt.Fprintf(w, jsonMsg) 
  }
}

func DelWord(w http.ResponseWriter, r *http.Request){
  rObj, err := getUpdatedBody(r.Body)
  
  if err != nil {
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    http.Error(w, jsonMsg, http.StatusOK)
    return
  }

  c := middlewares.GetWords(r)


  err = c.Remove(bson.M{"_id": bson.ObjectId(rObj.Id)})

  if err != nil {
    log.Println("Error Removing  : ", err)
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
    http.Error(w, jsonMsg, http.StatusNoContent)
    return
  } else {
    var status = true
    
    jsonMsg, _ := forgeResponse(Response{status, nil})

    fmt.Fprintf(w, jsonMsg)
    return 
    
  }
}

func UpdateWord(w http.ResponseWriter, r *http.Request){
  rObj, err := getUpdatedBody(r.Body)
  
  if err != nil {
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    http.Error(w, jsonMsg, http.StatusOK)
    return
  }

  c := middlewares.GetWords(r)

  q := buildUpdateQuery(&rObj)

  change := mgo.Change{
    Update: bson.M{"$set": q},
    ReturnNew: false,
  }

  word := UpdatedWord{}
  _, err = c.Find(bson.M{"_id": bson.ObjectId(rObj.Id)}).Apply(change, &word)

  if err != nil {
    log.Println("Error fetching : ", err)
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
    http.Error(w, jsonMsg, http.StatusNoContent)
    return
  } else {
    var status = true
    
    jsonMsg, _ := forgeResponse(Response{status, nil})

    fmt.Fprintf(w, jsonMsg)
    return 
    
  }
}


func FilterWords(w http.ResponseWriter, r *http.Request){
  rObj, err := getBody(r.Body)
  
  if err != nil {
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    http.Error(w, jsonMsg, http.StatusOK)
    return
  }

  c := middlewares.GetWords(r)
  q := buildQuery(&rObj)

  words := []Word{}
  err = c.Find(bson.M{ "$or": q }).Sort("word").All(&words)

  if err != nil {
    log.Println("Error fetching : ", err)
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
    http.Error(w, jsonMsg, http.StatusNoContent)
    return
  } else {
    var status = true
    if len(words) == 0 {
      status = false
    }
    
    jsonMsg, _ := forgeResponse(Response{status, words})

    fmt.Fprintf(w, jsonMsg)
    return 
    
  }
}

func FindWords(w http.ResponseWriter, r *http.Request){
  c := middlewares.GetWords(r)

  words := []Word{} 

  err := c.Find(bson.M{}).Sort("word").All(&words)

  if err != nil {
    log.Println("Error fetching : ", err)
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
    http.Error(w, jsonMsg, http.StatusNoContent)
    return
  } else {

    var status = true
    if len(words) == 0 {
      status = false
    }
    
    jsonMsg, _ := forgeResponse(Response{status, words})

    fmt.Fprintf(w, jsonMsg)
    return 
  }

}

func contains (s []string)bool {
  for _, v := range s { if v == "text/json" || v == "application/json" { return true }}
  return false 
}
  

func RandomWord(w http.ResponseWriter, r *http.Request){
  Accept := strings.Split(r.Header.Get("Accept"), ",")
  returnJson := contains(Accept)

  c := middlewares.GetWords(r)

  words := []Word{} 

  err := c.Find(bson.M{}).All(&words)

  if err != nil {
    log.Println("Error fetching : ", err)
    jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
    http.Error(w, jsonMsg, http.StatusNoContent)
    return
  } else {
    if len(words) == 0 {
      return
    }
    rand.Seed(time.Now().UnixNano())
    index := rand.Intn(len(words))

    if returnJson == true {
      jsonMsg, _ := forgeResponse(Response{true, words[index:index+1]})
      fmt.Fprintf(w, jsonMsg)
    } else {
      tmpl, err := template.ParseFiles("server/views/word.html")
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
      }
      err = tmpl.Execute(w, words[index:index+1])
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
      }
    }

    return 
  }

  jsonMsg, _ := forgeResponse(Response{false, nil}) 
    
  http.Error(w, jsonMsg, http.StatusInternalServerError)
  return 
}

type Response struct {
  Result bool
  Words []Word `json:",omitempty"`
}

func forgeResponse(rep Response) (string, error){
  jbMsg, err := json.Marshal(rep)

  if err != nil {    
    return "", err
  }

  jsonMsg := string(jbMsg[:]) // converting byte array to string
  return jsonMsg, nil
}

func buildUpdateQuery(rObj *UpdatedWord) (q map[string]string){
  q = make(map[string]string)

  if rObj.Word != "" {
    q["word"] = rObj.Word
  }

  if rObj.Pos != "" {
    q["pos"] = rObj.Pos
  }

  if rObj.Definition != "" {
    q["definition"] = rObj.Definition
  }

  if rObj.Language != "" {
    q["language"] = rObj.Language
  }

  if rObj.Example != "" {
    q["example"] = rObj.Example
  }

  if rObj.Ethymology != "" {
    q["ethymology"] = rObj.Ethymology
  }
  return q
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

func getBody(body io.Reader) (t Word, e error) {
  decoder := json.NewDecoder(body)
  err := decoder.Decode(&t)
  if err != nil {
    log.Println("Error deconding request Body ")
  }

  emptyWord := Word{}

  if t == emptyWord {
    return t, errors.New("Empty Object")
  }
  return t, nil
}


func getUpdatedBody(body io.Reader) (t UpdatedWord, e error) {
  decoder := json.NewDecoder(body)
  err := decoder.Decode(&t)
  if err != nil {
    log.Println("Error deconding request Body ", err)
  }

  emptyWord := UpdatedWord{}

  if t == emptyWord {
    return t, errors.New("Empty Object")
  }


  return t, nil
  
}