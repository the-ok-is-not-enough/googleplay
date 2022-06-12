package googleplay

import (
   "fmt"
   "os"
   "testing"
   "time"
)

const sleep = time.Minute

func TestHeader(t *testing.T) {
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   token, err := OpenToken(home, "googleplay/token.json")
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(sleep)
   for i := 9; i >= 1; i-- {
      time.Sleep(sleep)
      head, err := token.Header(0, false)
      if err != nil {
         t.Fatal(err)
      }
      if head.Auth == "" {
         t.Fatalf("%+v", head)
      }
   }
}

func TestToken(t *testing.T) {
   token, err := NewToken(email, password)
   if err != nil {
      t.Fatal(err)
   }
   home, err := os.UserHomeDir()
   if err != nil {
      t.Fatal(err)
   }
   if err := token.Create(home, "googleplay/token.json"); err != nil {
      t.Fatal(err)
   }
}
