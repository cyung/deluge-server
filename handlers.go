package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "io"
  "io/ioutil"
)

var magnets map[string]bool
var key string

func init() {
  magnets = make(map[string]bool)
  key = LoadConfig()
}

func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Welcome!")
}

// returns all magnets
func GetMagnets(w http.ResponseWriter, r *http.Request) {
  if client_key != key {
    w.WriteHeader(401)
    return
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  err := json.NewEncoder(w).Encode(magnets)

  if err != nil {
    panic(err)
  }
}

func AddMagnet(w http.ResponseWriter, r *http.Request) {
  if client_key != key {
    w.WriteHeader(401)
    return
  }

  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 500000))

  if err != nil {
    panic(err)
  }

  err = r.Body.Close()
  if err != nil {
    panic(err)
  }

  var magnet Magnet
  err = json.Unmarshal(body, &magnet)

  // data is not formatted correctly
  if err != nil {
    w.WriteHeader(422) // cannot be processed
    return
  }

  // there was no field "magnet"
  if magnet.Magnet == "" {
    w.WriteHeader(400)
    return
  }

  magnets[magnet.Magnet] = true

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)
  err = json.NewEncoder(w).Encode(magnets)

  if err != nil {
    panic(err)
  }
}

func DeleteMagnet(w http.ResponseWriter, r *http.Request) {
  if client_key != key {
    w.WriteHeader(401)
    return
  }

  magnet := r.URL.Query().Get("magnet")
  fmt.Printf("deleting %v\n", magnet)

  _, ok := magnets[magnet]

  if !ok {
    w.WriteHeader(404)
    return
  }

  delete(torrents, magnet)

  w.WriteHeader(200)
}
