package gplayapi

import (
   "fmt"
   "testing"
)

const packageName = "com.google.android.projection.gearhead.phonescreen"

func TestPlay(t *testing.T) {
   client, err := NewClientWithDeviceInfo("srpen6@gmail.com", "", Pixel3a)
   if err != nil {
      t.Fatal(err)
   }
   app, err := client.GetAppDetails(packageName)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println("VersionCode", app.VersionCode)
}
