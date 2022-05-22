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

func doHeader(platform string, agentID int64) (*gp.Header, error) {
   cache, err := os.UserHomeDir()
   if err != nil {
      return nil, err
   }
   token, err := gp.OpenToken(cache, "googleplay/token.json")
   if err != nil {
      return nil, err
   }
   device, err := gp.OpenDevice(cache, "googleplay", platform + ".json")
   if err != nil {
      return nil, err
   }
   return token.Header(device.AndroidID, agentID)
}

func doDetails(head *gp.Header, app string) error {
   detail, err := head.Details(app)
   if err != nil {
      return err
   }
   date, err := time.Parse(gp.DateInput, string(detail.UploadDate))
   if err == nil {
      detail.UploadDate = gp.String(date.Format(gp.DateOutput))
   }
   fmt.Println(detail)
   return nil
}

func doDelivery(head *gp.Header, app string, ver uint64) error {
   download := func(addr gp.String, name string) error {
      fmt.Println("GET", addr)
      res, err := http.Get(string(addr))
      if err != nil {
         return err
      }
      defer res.Body.Close()
      file, err := os.Create(name)
      if err != nil {
         return err
      }
      defer file.Close()
      pro := format.ProgressBytes(file, res.ContentLength)
      if _, err := io.Copy(pro, res.Body); err != nil {
         return err
      }
      return nil
   }
   del, err := head.Delivery(app, ver)
   if err != nil {
      return err
   }
   for _, split := range del.SplitDeliveryData {
      err := download(split.DownloadURL, del.Split(split.ID))
      if err != nil {
         return err
      }
   }
   for _, file := range del.AdditionalFile {
      err := download(file.DownloadURL, del.Additional(file.FileType))
      if err != nil {
         return err
      }
   }
   return download(del.DownloadURL, del.Download())
}

func doToken(email, password string) error {
   tok, err := gp.NewToken(email, password)
   if err != nil {
      return err
   }
   cache, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   return tok.Create(cache, "googleplay/token.json")
}

func doDevice(platform string) error {
   cache, err := os.UserHomeDir()
   if err != nil {
      return err
   }
   device, err := gp.Phone.Checkin(platform)
   if err != nil {
      return err
   }
   fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
   time.Sleep(gp.Sleep)
   return device.Create(cache, "googleplay", platform + ".json")
}
