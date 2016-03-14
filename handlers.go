package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "io"
  "io/ioutil"
)

var torrents map[string]bool

func init() {
  torrents = make(map[string]bool)
}

func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Welcome!")
}

// returns all torrents
func GetTorrents(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  err := json.NewEncoder(w).Encode(torrents)

  if err != nil {
    panic(err)
  }
}

func AddTorrent(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 500000))

  if err != nil {
    panic(err)
  }

  err = r.Body.Close()
  if err != nil {
    panic(err)
  }

  var torrent Torrent
  err = json.Unmarshal(body, &torrent)

  // data is not formatted correctly
  if err != nil {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(422) // cannot be processed

    err = json.NewEncoder(w).Encode(err)

    if err != nil {
      panic(err)
    }

    return
  }

  // there was no field "magnet"
  if torrent.Magnet == "" {
    w.WriteHeader(400)
    return
  }

  torrents[torrent.Magnet] = true

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)
  err = json.NewEncoder(w).Encode(torrents)

  if err != nil {
    panic(err)
  }
}

func DeleteTorrent(w http.ResponseWriter, r *http.Request) {
  magnet := r.URL.Query().Get("magnet")
  fmt.Printf("deleting %v\n", magnet)

  _, ok := torrents[magnet]

  if !ok {
    w.WriteHeader(404)
    return
  }

  delete(torrents, magnet)

  w.WriteHeader(200)
}

