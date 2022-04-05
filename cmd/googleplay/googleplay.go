package main

import (
   "fmt"
   "github.com/89z/format"
   "net/http"
   "os"
   "time"
   gp "github.com/89z/googleplay"
)

func doDevice(tv, tablet bool) error {
   cache, err := os.UserCacheDir()
   if err != nil {
      return err
   }
   var (
      config = gp.Phone
      elem = "googleplay/phone.json"
   )
   if tv {
      config = gp.TV
      elem = "googleplay/tv.json"
   } else if tablet {
      config = gp.Tablet
      elem = "googleplay/tablet.json"
   }
   device, err := config.Checkin()
   if err != nil {
      return err
   }
   fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
   time.Sleep(gp.Sleep)
   return device.Create(cache, elem)
}

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
      defer res.Body.Close()
      file, err := os.Create(data.Name(app, ver))
      if err != nil {
         return err
      }
      defer file.Close()
      pro := format.NewProgress(res)
      if _, err := file.ReadFrom(pro); err != nil {
         return err
      }
   }
   return nil
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

func newHeader(tv, tablet, single bool) (*gp.Header, error) {
   cache, err := os.UserCacheDir()
   if err != nil {
      return nil, err
   }
   token, err := gp.OpenToken(cache, "googleplay/token.json")
   if err != nil {
      return nil, err
   }
   elem := "googleplay/phone.json"
   if tv {
      elem = "googleplay/tv.json"
   } else if tablet {
      elem = "googleplay/tablet.json"
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
