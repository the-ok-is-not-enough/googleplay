package play

import (
   "fmt"
   "testing"
)

//pass
const app = "com.google.android.projection.gearhead.phonescreen"

//fail
//const app = "com.google.android.youtube"

func TestDiscord(t *testing.T) {
   det, err := details(app)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(det)
}
