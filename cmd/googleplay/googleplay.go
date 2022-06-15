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

func doHeader(dir, platform string, single bool) (*gp.Header, error) {
   token, err := gp.OpenToken(dir + "/token.txt")
   if err != nil {
      return nil, err
   }
   device, err := gp.OpenDevice(dir + "/" + platform + ".txt")
   if err != nil {
      return nil, err
   }
   id, err := device.AndroidID()
   if err != nil {
      return nil, err
   }
   return token.Header(id, single)
}

func doDevice(dir, platform string) error {
   device, err := gp.Phone.Checkin(platform)
   if err != nil {
      return err
   }
   fmt.Printf("Sleeping %v for server to process\n", gp.Sleep)
   time.Sleep(gp.Sleep)
   file, err := format.Create(dir + "/" + platform + ".txt")
   if err != nil {
      return err
   }
   defer file.Close()
   if _, err := device.WriteTo(file); err != nil {
      return err
   }
   return nil
}

func doToken(dir, email, password string) error {
   token, err := gp.NewToken(email, password)
   if err != nil {
      return err
   }
   file, err := format.Create(dir + "/token.txt")
   if err != nil {
      return err
   }
   defer file.Close()
   if _, err := token.WriteTo(file); err != nil {
      return err
   }
   return nil
}

func doDetails(head *gp.Header, app string, parse bool) error {
   detail, err := head.Details(app)
   if err != nil {
      return err
   }
   if parse {
      date, err := detail.Time()
      if err != nil {
         return err
      }
      detail.UploadDate = date.String()
   }
   fmt.Println(detail)
   return nil
}

func doDelivery(head *gp.Header, app string, ver uint64) error {
   download := func(addr, name string) error {
      fmt.Println("GET", addr)
      res, err := http.Get(addr)
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
