package main

import (
    "flag"
    "os"
    "fmt"
    "http"
    "log"
    "io/ioutil"
    "json"
    )

var ip = flag.String("address", "cube.leppoc.net", "cuboboplis http server")
var port = flag.String("port", "10234", "http service port")

type ChunkData [16][16][]string

type ChunkInfo struct {
  X int
  Y int
  Data *ChunkData
}

func GetChunk(chunkX, chunkY int) (*ChunkInfo, os.Error) {
  address := fmt.Sprintf("http://%s:%s/a/r?cy=%d&cx=%d", *ip, *port, chunkY, chunkX)
  fmt.Println("Making request " + address)
  r, _, err := http.Get(address)
  if err != nil {
    log.Panic(err)
  }
  defer r.Body.Close()
  content, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Panic(err)
  }
  info := &ChunkInfo{chunkX, chunkY, &ChunkData{}}
  err = json.Unmarshal(content, &info.Data)
  if err != nil {
    fmt.Println(err.String())
  }
  for y := 0; y < 16; y++ {
    for x := 0; x < 16; x++ {
      if len(info.Data[y][x]) == 0 {
        info.Data[y][x] = append(info.Data[y][x], "0")
      } else if info.Data[y][x][0] == "" {
        info.Data[y][x][0] = "0"
      }
    }
  }
  fmt.Println(info.Data)
  return info, err
}

