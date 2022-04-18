package main

import (
   "fmt"
   "github.com/89z/format"
   "io"
   "net/http"
   "os"
   "time"
   gp "github.com/89z/googleplay"
)

func doDelivery(head *gp.Header, app string, ver uint64) error {
   del, err := head.Delivery(app, ver)
   if err != nil {
      return err
   }
   for _, data := range del.Data() {
      fmt.Println("GET", data.DownloadURL)
      res, err := http.Get(string(data.DownloadURL))
      if err != nil {
         return err
      }
      file, err := os.Create(data.Name(app, ver))
      if err != nil {
         return err
      }
      pro := format.ProgressBytes(file, res.ContentLength)
      if _, err := io.Copy(pro, res.Body); err != nil {
         return err
      }
      if err := res.Body.Close(); err != nil {
         return err
      }
      if err := file.Close(); err != nil {
         return err
      }
   }
   return nil
}

func doDevice(armeabi, arm64 bool) error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   var (
      elem = "googleplay/x86.json"
      platform = gp.X86
   )
   if armeabi {
      elem = "googleplay/armeabi.json"
      platform = gp.Armeabi
   } else if arm64 {
      elem = "googleplay/arm64.json"
      platform = gp.Arm64
   }
   device, err := gp.Phone.Checkin(platform)
   if err != nil {
      return err
   }
   fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
   time.Sleep(gp.Sleep)
   return device.Create(cache, elem)
}

func doToken(email, password string) error {
   tok, err := gp.NewToken(email, password)
   if err != nil {
      return err
   }
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   return tok.Create(cache, "googleplay/token.json")
}

func newHeader(armeabi, arm64, single bool) (*gp.Header, error) {
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, err
   }
   token, err := gp.OpenToken(cache, "googleplay/token.json")
   if err != nil {
      return nil, err
   }
   elem := "googleplay/x86.json"
   if armeabi {
      elem = "googleplay/armeabi.json"
   } else if arm64 {
      elem = "googleplay/arm64.json"
   }
   device, err := gp.OpenDevice(cache, elem)
   if err != nil {
      return nil, err
   }
   if single {
      return token.SingleAPK(device)
   }
   return token.Header(device)
}
