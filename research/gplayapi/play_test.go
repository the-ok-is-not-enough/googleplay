package gplayapi

import (
   "fmt"
   "testing"
   "time"
)

var names = []string{
   //pass
   //"com.google.android.projection.gearhead.phonescreen",
   //"com.google.android.youtube",
   //fail?
   "com.xiaomi.smarthome",
}

func TestPlay(t *testing.T) {
   client, err := NewClient("srpen6@gmail.com", token)
   if err != nil {
      t.Fatal(err)
   }
   time.Sleep(4 * time.Second)
   for _, name := range names {
      app, err := client.GetAppDetails(name)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println("VersionCode", app.VersionCode)
      time.Sleep(time.Second)
   }
}
