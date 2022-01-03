package play

import (
   "fmt"
   "testing"
   "time"
)

var apps = []string{
   "com.google.android.projection.gearhead.phonescreen",
   "com.google.android.youtube",
}

func TestPlay(t *testing.T) {
   for _, app := range apps {
      det, err := details(app)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(det)
      time.Sleep(time.Second)
   }
}
