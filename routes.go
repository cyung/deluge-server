package main

import (
  "net/http"
  "github.com/gorilla/mux"
)

type Route struct {
  Name string
  Method string
  Pattern string
  HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
  router := mux.NewRouter().StrictSlash(true)

  router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))
  http.Handle("/", router)
  
  for _, route := range routes {
    router.
      Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(route.HandlerFunc)
  }


  return router
}

var routes = Routes {
  Route {
    "Magnets",
    "POST",
    "/magnets",
    AddMagnet,
  },
  Route {
    "Torrents",
    "GET",
    "/torrents",
    GetTorrents,
  },
  Route {
    "Torrents",
    "POST",
    "/torrents",
    AddTorrent,
  },
  Route {
    "Torrents",
    "DELETE",
    "/torrents",
    DeleteTorrent,
  },
}