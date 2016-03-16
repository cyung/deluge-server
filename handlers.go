package main

import (
  "net/http"
  "math/rand"
  "time"
  "io"
  "io/ioutil"
  "os"
  "archive/zip"
  "fmt"
  "encoding/json"
)

// return a zip of all files stored on the server
func GetTorrents(w http.ResponseWriter, r *http.Request) {
  if !Validate(r.Header.Get("Authorization")) {
    w.WriteHeader(401)
    return
  }

  // get list of all torrents in directory
  files, err := ioutil.ReadDir("./torrents")
  if err != nil {
    w.WriteHeader(500)
    panic(err)
  }

  // return if no files
  if len(files) == 0 {
    w.WriteHeader(404)
    return
  }
  
  if len(files) == 1 && files[0].Name() == ".DS_Store" {
    w.WriteHeader(404)
    return
  }

  // create zip file
  zip_filename := "./tmp/" + RandomFilename() + ".zip"
  zip_file, err := os.Create(zip_filename)
  if err != nil {
    w.WriteHeader(500)
    panic(err)
  }

  zip_writer := zip.NewWriter(zip_file)

  // add files to zip file
  for _, file := range files {
    data, err := ioutil.ReadFile("./torrents/" + file.Name())
    if err != nil {
      w.WriteHeader(500)
      panic(err)
    }

    f, err := zip_writer.Create(file.Name())
    if err != nil {
      w.WriteHeader(500)
      panic(err)
    }

    _, err = f.Write(data)
    if err != nil {
      w.WriteHeader(500)
      panic(err)
    }
  }

  err = zip_writer.Close()
  if err != nil {
    w.WriteHeader(500)
    panic(err)
  }

  zip_file.Close()

  // serve zip file to client
  w.Header().Set("Content-Disposition", "attachment; filename=" + RandomFilename() + ".zip")
  w.Header().Set("Content-Type", "application/zip")
  http.ServeFile(w, r, zip_filename)

  // remove file after serving
  if err := os.Remove(zip_filename); err != nil {
    panic(err)
  }
}

// upload a torrent file to FS
// file checking occurs on front-end
func AddTorrent(w http.ResponseWriter, r *http.Request) {
  if !Validate(r.Header.Get("Authorization")) {
    w.WriteHeader(401)
    return
  }

  file, _, err := r.FormFile("torrent")
  if err != nil {
    w.WriteHeader(500)
    panic(err)
  }
  defer file.Close()

  out, err := os.Create("./torrents/" + RandomFilename() + ".torrent")
  if err != nil {
    w.WriteHeader(500)
    fmt.Println("error creating file")
    panic(err)
  }
  defer out.Close()

  // verify file data
  _, err = io.Copy(out, file)
  if err != nil {
    panic(err)
  }

  w.WriteHeader(201)
}

func AddMagnet(w http.ResponseWriter, r *http.Request) {
  if !Validate(r.Header.Get("Authorization")) {
    w.WriteHeader(401)
    return
  }

  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 500000))
  if err != nil {
    panic(err)
  }
  defer r.Body.Close()

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

  out, err := os.Create("./torrents/" + RandomFilename() + ".magnet")
  if err != nil {
    w.WriteHeader(500)
    fmt.Println("error creating file")
    panic(err)
  }
  defer out.Close()

  _, err = out.WriteString(magnet.Magnet)
  if err != nil {
    fmt.Println("error writing to magnet file")
    panic(err)
  }

  w.WriteHeader(http.StatusCreated)
}

// sanitize input by generating random filename
func RandomFilename() string {
  const CHAR_LENGTH = 10
  const chars = "abcdefghijklmnopqrstuvwxyz" +
                "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
                "0123456789"

  rand.Seed(time.Now().UTC().UnixNano())
  result := make([]byte, CHAR_LENGTH)

  for i := 0; i<CHAR_LENGTH; i++ {
    result[i] = chars[rand.Intn(len(chars))]
  }

  return string(result)
}

// delete a torrent file from FS
func DeleteTorrent(w http.ResponseWriter, r *http.Request) {
    if !Validate(r.Header.Get("Authorization")) {
    w.WriteHeader(401)
    return
  }

  filename := r.URL.Query().Get("torrent")

  if err := os.Remove("./torrents/" + filename); err != nil {
    if os.IsNotExist(err) {
      w.WriteHeader(404)
      return
    }

    w.WriteHeader(400)
    panic(err)
  }

  w.WriteHeader(200)
}