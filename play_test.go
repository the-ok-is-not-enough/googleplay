package googleplay

import (
   "os"
   "testing"
   "time"
)

func Test_Token(t *testing.T) {
   token, err := New_Token(email, password)
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   if err := token.Create(home + "/googleplay/token.txt"); err != nil {
      t.Fatal(err)
   }
}

func Test_Header(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   token, err := Open_Token(home + "/googleplay/token.txt")
   if err != nil {
      t.Fatal(err)
   }
   for i := 0; i < 9; i++ {
      head, err := token.Header(0, false)
      if err != nil {
         t.Fatal(err)
      }
      if head.Auth == "" {
         t.Fatalf("%+v", head)
      }
      time.Sleep(time.Second)
   }
}
