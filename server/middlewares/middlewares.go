package middlewares

import (	
	"log"
  "net/http"
  "sync"
  "time"
  "gopkg.in/mgo.v2"
  
  // "gopkg.in/mgo.v2/bson"
)


type (
	Context struct {
		MongoDB *mgo.Database
	}
)


func Mongo(h http.Handler, url string, port string, dbName string) http.Handler {
	session, err := mgo.Dial(url+":"+port)
  if err != nil {
    log.Fatal("Error connecting to MongoDB Server : ", err)
  }

  log.Println("Connected to MongoDB on "+dbName)

	f := func(w http.ResponseWriter, r *http.Request) {

    reqSession := session.Clone()
    defer reqSession.Close()

    // Optional. Switch the session to a monotonic behavior.
  	session.SetMode(mgo.Monotonic, true)

    db := reqSession.DB(dbName)

    c := db.C("Words")

    index := mgo.Index{
      Key: []string{"word", "pos"},
      Unique: true,
    }

    c.EnsureIndex(index)

    set(r, "db", db)

  	h.ServeHTTP(w, r)
  }

	return http.HandlerFunc(f)
}


var (
    mutex sync.Mutex
    data  = make(map[*http.Request]map[interface{}]interface{})
    datat = make(map[*http.Request]int64)
)

// Set stores a value for a given key in a given request.
func set(r *http.Request, key, val interface{}) {
  mutex.Lock()
  defer mutex.Unlock()
  if data[r] == nil {
    data[r] = make(map[interface{}]interface{})
    datat[r] = time.Now().Unix()
    // log.Println("Adding new Entry to Context")
  }
  data[r][key] = val
}

func GetDB(r *http.Request) (d *mgo.Database){
  mutex.Lock()
  defer mutex.Unlock()
  if data[r] != nil {
    return data[r]["db"].(*mgo.Database)
  }
  return d
}