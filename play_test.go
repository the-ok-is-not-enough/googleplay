package googleplay

import (
   "os"
   "testing"
   "time"
)

func Test_Auth(t *testing.T) {
   auth, err := New_Auth(email, password)
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   if err := auth.Create(home + "/googleplay/auth.txt"); err != nil {
      t.Fatal(err)
   }
}

func Test_Header(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   auth, err := Open_Auth(home + "/googleplay/auth.txt")
   if err != nil {
      t.Fatal(err)
   }
   for i := 0; i < 9; i++ {
      head, err := auth.Header(0, false)
      if err != nil {
         t.Fatal(err)
      }
      if head.Auth == "" {
         t.Fatalf("%+v", head)
      }
      time.Sleep(time.Second)
   }
}
