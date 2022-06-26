package googleplay

import (
   "fmt"
   "os"
   "testing"
)

func Test_Delivery(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   auth, err := Open_Auth(home + "/googleplay/auth.txt")
   if err != nil {
      t.Fatal(err)
   }
   device, err := Open_Device(home + "/googleplay/x86.bin")
   if err != nil {
      t.Fatal(err)
   }
   id, err := device.ID()
   if err != nil {
      t.Fatal(err)
   }
   head, err := auth.Header(id, false)
   if err != nil {
      t.Fatal(err)
   }
   del, err := head.Delivery("com.google.android.youtube", 1524221376)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", del)
}
