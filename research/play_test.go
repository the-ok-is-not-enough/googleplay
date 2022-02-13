package gplayapi

import (
   "fmt"
   "time"
   "testing"
)

var names = []string{
   "com.google.android.youtube",
   "com.bbca.bbcafullepisodes",
}

func TestPlay(t *testing.T) {
   client, err := NewClientWithDeviceInfo("srpen6@gmail.com", token, Pixel3a)
   if err != nil {
      panic(err)
   }
   for _, name := range names {
      app, err := client.GetAppDetails(name)
      if err != nil {
         panic(err)
      }
      fmt.Println("VersionCode", app.VersionCode)
      time.Sleep(time.Second)
   }
}
