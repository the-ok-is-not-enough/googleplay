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

func doDevice() error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   device, err := gp.Phone.Checkin()
   if err != nil {
      return err
   }
   fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
   time.Sleep(gp.Sleep)
   return device.Create(cache, "googleplay/phone.json")
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

func newHeader(single bool) (*gp.Header, error) {
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, err
   }
   token, err := gp.OpenToken(cache, "googleplay/token.json")
   if err != nil {
      return nil, err
   }
   device, err := gp.OpenDevice(cache, "googleplay/phone.json")
   if err != nil {
      return nil, err
   }
   if single {
      return token.SingleAPK(device)
   }
   return token.Header(device)
}
