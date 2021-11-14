package main

import (
   "fmt"
   "github.com/xi2/xz"
   "net/http"
   "os"
   "path/filepath"
   "strings"
)

const frida =
   "https://github.com/frida/frida/releases/download/" +
   "15.1.10/frida-server-15.1.10-android-x86.xz"

func download(src, dst string) error {
   fmt.Println("GET", src)
   rHTTP, err := http.Get(src)
   if err != nil {
      return err
   }
   defer rHTTP.Body.Close()
   rXZ, err := xz.NewReader(rHTTP.Body, 0)
   if err != nil {
      return err
   }
   file, err := os.Create(dst)
   if err != nil {
      return err
   }
   defer file.Close()
   if _, err := file.ReadFrom(rXZ); err != nil {
      return err
   }
   return nil
}

func stem(s string) string {
   low := strings.LastIndexByte(s, '/')
   if low == -1 {
      return s
   }
   s = s[low+1:]
   high := strings.LastIndexByte(s, '.')
   if high == -1 {
      return s
   }
   return s[:high]
}

func main() {
   cache, err := os.UserCacheDir()
   if err != nil {
      panic(err)
   }
   cache = filepath.Join(cache, "googleplay")
   os.Mkdir(cache, os.ModeDir)
   cache = filepath.Join(cache, stem(frida))
   fmt.Println("Stat", cache)
   if _, err := os.Stat(cache); err != nil {
      err := download(frida, cache)
      if err != nil {
         panic(err)
      }
   }
}
