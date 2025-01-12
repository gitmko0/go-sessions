package main

import (
	"fmt"
	"net/http"

	"github.com/gitmko0/go-sessions"
	"github.com/gitmko0/go-sessions/sessiondb/redis"
	"github.com/gitmko0/go-sessions/sessiondb/redis/service"
)

func main() {
	// replace with your running redis' server settings:
	db := redis.New(service.Config{Network: service.DefaultRedisNetwork,
		Addr:        service.DefaultRedisAddr,
		Password:    "",
		Database:    "",
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: service.DefaultRedisIdleTimeout,
		Prefix:      "",
	}) // to use badger just use the sessiondb/badger#New func.

	defer db.Close()

	sess := sessions.New(sessions.Config{Cookie: "s"})

	//
	// IMPORTANT:
	//
	sess.UseDatabase(db)

	// the rest of the code stays the same.
	app := http.NewServeMux()

	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("You should navigate to the /set, /get, /delete, /clear,/destroy instead")))
	})
	app.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		s := sess.Start(w, r)
		//set session values
		s.Set("", "iris")

		//test if setted here
		w.Write([]byte(fmt.Sprintf("All ok session setted to: %s", s.GetString(""))))
	})

	app.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		// get a specific key, as string, if no found returns just an empty string
		name := sess.Start(w, r).GetString("")

		w.Write([]byte(fmt.Sprintf("The name on the /set was: %s", name)))
	})

	app.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		// delete a specific key
		sess.Start(w, r).Delete("")
	})

	app.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		// removes all entries
		sess.Start(w, r).Clear()
	})

	app.HandleFunc("/destroy", func(w http.ResponseWriter, r *http.Request) {
		//destroy, removes the entire session data and cookie
		sess.Destroy(w, r)
	})

	app.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		// updates expire date with a new date
		sess.ShiftExpiration(w, r)
	})

	http.ListenAndServe(":80", app)
}
