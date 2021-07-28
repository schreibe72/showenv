package main

import (
	"container/list"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/logmatic/logmatic-go"
	log "github.com/sirupsen/logrus"
)

func showEnv(w http.ResponseWriter, r *http.Request) {
	kv := map[string]map[string]string{}
	kv["url"] = map[string]string{}
	kv["envvar"] = map[string]string{}
	kv["headers"] = map[string]string{}
	kv["url"]["Request-PATH"] = r.URL.Path
	kv["url"]["Request-HOST"] = r.Host
	for _, i := range os.Environ() {
		s := strings.SplitN(i, "=", 2)
		kv["envvar"][s[0]] = s[1]
	}
	for k := range r.Header {
		o := r.Header.Values(k)
		b, _ := json.Marshal(&o)
		kv["headers"][k] = string(b)
	}
	b, err := json.MarshalIndent(&kv, "", "  ")
	if err != nil {
		log.WithFields(log.Fields{"Type": "Error"}).Warn(err)
	}
	log.WithFields(log.Fields{"Type": "Access", "Path": r.URL.Path, "Host": r.Host}).Info()
	w.Write(b)
}

func die(w http.ResponseWriter, r *http.Request) {
	log.Fatalln("somebody wants me to die....")
}

func memory(w http.ResponseWriter, r *http.Request) {
	dstr := r.URL.Query().Get("sleep")
	d, err := strconv.ParseInt(dstr, 10, 32)
	if err != nil {
		d = 0
	}

	log.Printf("eating memory.... yummy - i like! (every %i msec)\n", d)
	s := "eating more memory.... yummy - i like!"
	go func() {
		l := list.New()
		for {
			time.Sleep(time.Duration(d) * time.Millisecond)
			l.PushBack(s)
		}
	}()
	w.Write([]byte(s))
}

func main() {
	log.SetFormatter(&logmatic.JSONFormatter{})
	http.HandleFunc("/", showEnv)
	http.HandleFunc("/die", die)
	http.HandleFunc("/memory", memory)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.WithFields(log.Fields{"Type": "Error"}).Panic(err)
	}
}
