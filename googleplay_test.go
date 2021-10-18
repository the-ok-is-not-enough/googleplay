package googleplay

import (
   "net/url"
   "os"
   "testing"
)

func TestRead(t *testing.T) {
   txt, err := os.ReadFile("ac2dm.txt")
   if err != nil {
      t.Fatal(err)
   }
   q, err := url.ParseQuery(string(txt))
   if err != nil {
      t.Fatal(err)
   }
   a := Ac2dm{q}
   o, err := a.OAuth2()
   if err != nil {
      t.Fatal(err)
   }
   b, err := o.Details("38B5418D8683ADBB", "com.google.android.youtube")
   if err != nil {
      t.Fatal(err)
   }
   os.Stdout.Write(b)
}

func TestWrite(t *testing.T) {
   a, err := NewAc2dm("oauth2_4/0AX4XfWhSw5Vdq96VAC28j7ZEdoryJMWUbLibNDbgC9NsIk_3WHAqkDRi7PfwEa7kEDsDnw")
   if err != nil {
      t.Fatal(err)
   }
   f, err := os.Create("ac2dm.txt")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   f.WriteString(a.Encode())
}
