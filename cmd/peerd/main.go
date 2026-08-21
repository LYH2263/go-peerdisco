package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"example.com/peerdisco"
	"example.com/peerdisco/internal/transport"
)

func main() {
	addr := flag.String("addr", ":8105", "listen address")
	web := flag.String("web", "web", "static web dir")
	data := flag.String("data", "data", "data dir")
	flag.Parse()
	_ = os.MkdirAll(*data, 0o755)

	tr := &transport.Mock{}
	c, err := peerdisco.New(peerdisco.Options{
		SelfID:      "peerd-self",
		SelfAddr:    "127.0.0.1:19001",
		Transport:   tr,
		PersistPath: filepath.Join(*data, "members.json"),
		AuditPath:   filepath.Join(*data, "audit.log"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	_ = c.Join(peerdisco.JoinSpec{ID: "peerd-self", Addr: "127.0.0.1:19001", Tags: []string{"self"}})

	var mu sync.Mutex
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(*web)))
	mux.HandleFunc("/api/members", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		writeJSON(w, c.Members())
	})
	mux.HandleFunc("/api/suspicion", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		writeJSON(w, c.SuspicionList())
	})
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		writeJSON(w, c.Stats())
	})
	mux.HandleFunc("/api/join", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", 405)
			return
		}
		var spec peerdisco.JoinSpec
		if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		mu.Lock()
		err := c.Join(spec)
		mu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, map[string]string{"ok": "1"})
	})
	mux.HandleFunc("/api/suspect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", 405)
			return
		}
		var body struct{ ID, From string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		mu.Lock()
		err := c.MarkSuspect(body.ID, body.From)
		mu.Unlock()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		writeJSON(w, map[string]string{"ok": "1"})
	})
	log.Printf("peerd listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
