package googleplay

import (
   "bufio"
   "encoding/json"
   "github.com/89z/format"
   "github.com/89z/format/crypto"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
   "strings"
)

type Token struct {
   Token string
}

// Request refresh token.
func NewToken(email, password string) (*Token, error) {
   val := url.Values{
      "Email": {email},
      "Passwd": {password},
      // Instead of the following two, you can instead use this:
      // sdk_version=20
      // but I couldnt get newer versions to work, so I think this is the
      // better option.
      "client_sig": {""},
      "droidguard_results": {""},
   }.Encode()
   req, err := http.NewRequest(
      "POST", origin + "/auth", strings.NewReader(val),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   hello, err := crypto.ParseJA3(crypto.AndroidAPI26)
   if err != nil {
      return nil, err
   }
   format.Log.Dump(req)
   res, err := crypto.Transport(hello).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf := bufio.NewScanner(res.Body)
   for buf.Scan() {
      kv := strings.SplitN(buf.Text(), "=", 2)
      if len(kv) == 2 && kv[0] == "Token" {
         var tok Token
         tok.Token = kv[1]
         return &tok, nil
      }
   }
   return nil, notFound{"Token"}
}

func OpenToken(name string) (*Token, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   tok := new(Token)
   if err := json.NewDecoder(file).Decode(tok); err != nil {
      return nil, err
   }
   return tok, nil
}

// Exchange refresh token for access token.
func (t Token) Auth() (*Auth, error) {
   val := url.Values{
      "Token": {t.Token},
      "service": {"oauth2:https://www.googleapis.com/auth/googleplay"},
   }.Encode()
   req, err := http.NewRequest(
      "POST", origin + "/auth", strings.NewReader(val),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, response{res}
   }
   buf := bufio.NewScanner(res.Body)
   for buf.Scan() {
      kv := strings.SplitN(buf.Text(), "=", 2)
      if len(kv) == 2 && kv[0] == "Auth" {
         var auth Auth
         auth.Auth = kv[1]
         return &auth, nil
      }
   }
   return nil, notFound{"Auth"}
}

func (t Token) Create(name string) error {
   err := os.MkdirAll(filepath.Dir(name), os.ModeDir)
   if err != nil {
      return err
   }
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   enc := json.NewEncoder(file)
   enc.SetIndent("", " ")
   return enc.Encode(t)
}
