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
   id, err := checkinProto()
   if err != nil {
      t.Fatal(err)
   }
   time.Sleep(4 * time.Second)
   for _, app := range apps {
      det, err := details(id, app)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(det)
      time.Sleep(time.Second)
   }
}
