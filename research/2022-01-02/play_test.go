package play

import (
   "fmt"
   "testing"
)

type config struct {
   app string
   upload bool
}

var cons = []config{
   {"com.google.android.projection.gearhead.phonescreen", false},
   {"com.google.android.youtube", true},
}

func TestPlay(t *testing.T) {
   for _, con := range cons {
      det, err := details(con.app, con.upload)
      if err != nil {
         t.Fatal(err)
      }
      if det == 0 {
         t.Fatal(det)
      }
      fmt.Println(det)
   }
}
