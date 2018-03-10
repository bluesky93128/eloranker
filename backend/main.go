package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var pool *redis.Pool

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", os.Getenv("REDIS_ADDRESS"))
			if err != nil {
				panic(err)
			}
			return conn, err
		},
	}
	defer pool.Close()

	router := mux.NewRouter()
	router.HandleFunc("/ws", ServeWS)

	fs := http.FileServer(http.Dir("dist"))
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := os.Open(path.Join("dist", path.Clean(r.URL.Path)))
		if os.IsNotExist(err) {
			http.ServeFile(w, r, "dist/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	})

	var handler http.Handler = router
	handler = handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"}),
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedOrigins([]string{"*"}),
	)(router)
	// handler = handlers.ProxyHeaders(handler)
	handler = handlers.CombinedLoggingHandler(os.Stdout, handler)

	log.Println("Listening on :80")
	log.Fatal(http.ListenAndServe(":80", handler))
}
