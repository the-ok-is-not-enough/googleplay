package research

import (
   "fmt"
   "net/http"
   "net/http/httputil"
   "testing"
   "time"
)

func TestAuth(t *testing.T) {
   for i := 0; i < 9; i++ {
      res, err := auth()
      if err != nil {
         t.Fatal(err)
      }
      buf, err := httputil.DumpResponse(res, true)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(string(buf))
      if res.StatusCode != http.StatusOK {
         t.Fatal(res.Status)
      }
      if err := res.Body.Close(); err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
