package gplayapi

import (
   "encoding/json"
   "os"
   "testing"
)

func TestPlay(t *testing.T) {
   client, err := NewClient("srpen6@gmail.com", token)
   if err != nil {
      t.Fatal(err)
   }
   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", " ")
   enc.Encode(client.AuthData)
}
