package gplayapi

import (
   "encoding/json"
   "fmt"
   "os"
   "testing"
)

const packageName = "com.google.android.projection.gearhead.phonescreen"

func TestPlay(t *testing.T) {
   cache, err := os.UserCacheDir()
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Open(cache + "/googleplay/token.json")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   var tok struct {
      Token string
   }
   if err := json.NewDecoder(file).Decode(&tok); err != nil {
      t.Fatal(err)
   }
   client, err := NewClientWithDeviceInfo(
      "srpen6@gmail.com", tok.Token, Pixel3a,
   )
   if err != nil {
      t.Fatal(err)
   }
   app, err := client.GetAppDetails(packageName)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println("VersionCode", app.VersionCode)
}
