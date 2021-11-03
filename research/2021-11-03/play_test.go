package googleplay

import (
   "fmt"
   "os"
   "testing"
)

func TestUpload(t *testing.T) {
   buf, err := upload(auth, device)
   if err != nil {
      t.Fatal(err)
   }
   os.Stdout.Write(buf)
}

func TestCheck(t *testing.T) {
   check, err := newCheckin()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(check)
}
