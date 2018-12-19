package middlewares

import (
	"log"
	"net/http"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
)

var (
	collectionName string
)

func Mongo(
	h http.Handler,
	url string,
	port string,
	dbName string,
	collection string) http.Handler {

	collectionName = collection

	session, err := mgo.Dial(url + ":" + port)
	if err != nil {
		log.Fatal("Error connecting to MongoDB Server : ", err)
	}

	f := func(w http.ResponseWriter, r *http.Request) {

		reqSession := session.Clone()
		defer reqSession.Close()

		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)

		db := reqSession.DB(dbName)

		collec := db.C(collection)

		index := mgo.Index{
			Key:    []string{"word", "pos", "definition"},
			Unique: true,
		}

		collec.EnsureIndex(index)

		set(r, collection, collec)

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

func GetWords(r *http.Request) (d *mgo.Collection) {
	mutex.Lock()
	defer mutex.Unlock()
	if data[r] != nil {
		return data[r][collectionName].(*mgo.Collection)
	}
	return d
}
