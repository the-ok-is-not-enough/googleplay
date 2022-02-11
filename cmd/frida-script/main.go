package main

import (
   "fmt"
   "os"
   "os/exec"
   "path/filepath"
)

func main() {
   _, err := exec.LookPath("frida")
   if err != nil {
      panic("pip install frida-tools")
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      panic(err)
   }
   cache = filepath.Join(cache, "googleplay")
   os.Mkdir(cache, os.ModeDir)
   cacheScript := filepath.Join(cache, filepath.Base(script))
   if len(os.Args) == 2 {
      pack := os.Args[1]
      if err := downloadScript(cacheScript); err != nil {
         panic(err)
      }
      server := newServer(version)
      cacheServer := filepath.Join(cache, stem(server))
      if err := downloadServer(server, cacheServer); err != nil {
         panic(err)
      }
      if err := runCommands(cacheServer, cacheScript, pack); err != nil {
         panic(err)
      }
   } else {
      fmt.Println("frida-script [package]")
   }
}
