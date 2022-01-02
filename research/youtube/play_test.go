package play

import (
   "fmt"
   "testing"
)

//pass
const app = "com.google.android.youtube"

//fail
//const app = "com.google.android.projection.gearhead.phonescreen"

func TestDiscord(t *testing.T) {
   det, err := details(app)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(det)
}
