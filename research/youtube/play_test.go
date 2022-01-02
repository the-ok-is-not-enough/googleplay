package play

import (
   "fmt"
   "testing"
)

const app = "com.google.android.youtube"

func TestDiscord(t *testing.T) {
   det, err := details(app)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(det)
}
