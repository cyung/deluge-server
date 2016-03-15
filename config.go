package main

import (
  "io/ioutil"
  "log"
  "encoding/json"
)

type Configuration struct {
  ChrisKey string `json:"CHRIS_KEY"`
}

func LoadConfig() string {
  file, err := ioutil.ReadFile("./config.json")
  if err != nil {
    log.Fatal(err)
  }

  var config Configuration
  err = json.Unmarshal(file, &config)
  if err != nil {
    log.Fatal(err)
  }

  return config.ChrisKey
}