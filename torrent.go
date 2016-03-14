package main

type Torrent struct {
  Magnet string `json:"magnet"`
}

type Torrents []Torrent