package gplayapi

import (
   "fmt"
   "testing"
)

func TestPlay(t *testing.T) {
   client, err := NewClientWithDeviceInfo("srpen6@gmail.com", token, Pixel3a)
   if err != nil {
      t.Fatal(err)
   }
   det, err := client.GetAppDetails("com.kakaogames.twodin")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", det)
}
