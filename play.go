package googleplay

import (
   "bufio"
   "github.com/89z/format"
   "github.com/89z/format/crypto"
   "github.com/89z/format/http"
   "io"
   "net/url"
   "os"
   "strconv"
   "strings"
   "time"
)

func decode_query(r io.Reader) (url.Values, error) {
   buf := bufio.NewReader(r)
   vals := make(url.Values)
   for {
      key, err := buf.ReadString('=')
      if err == io.EOF {
         break
      } else if err != nil {
         return nil, err
      }
      val, err := buf.ReadString('\n')
      key = strings.TrimSuffix(key, "=")
      val = strings.TrimSuffix(val, "\n")
      vals.Add(key, val)
      if err == io.EOF {
         break
      } else if err != nil {
         return nil, err
      }
   }
   return vals, nil
}

func encode_query(w io.Writer, v url.Values) error {
   for key := range v {
      val := v.Get(key)
      if _, err := io.WriteString(w, key); err != nil {
         return err
      }
      if _, err := io.WriteString(w, "="); err != nil {
         return err
      }
      if _, err := io.WriteString(w, val); err != nil {
         return err
      }
      if _, err := io.WriteString(w, "\n"); err != nil {
         return err
      }
   }
   return nil
}

const Sleep = 4 * time.Second

var Client = http.Default_Client

func (t Token) Create(name string) error {
   file, err := format.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   return encode_query(file, t.Values)
}

func Open_Token(name string) (*Token, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   var tok Token
   tok.Values, err = decode_query(file)
   if err != nil {
      return nil, err
   }
   return &tok, nil
}

func (t Token) Header(device_ID uint64, single bool) (*Header, error) {
   // these values take from Android API 28
   body := url.Values{
      "Token": {t.Token()},
      "app": {"com.android.vending"},
      "client_sig": {"38918a453d07199354f8b19af05ec6562ced5788"},
      "service": {"oauth2:https://www.googleapis.com/auth/googleplay"},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/auth", strings.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   var head Header
   head.SDK = 9
   head.Device_ID = device_ID
   if single {
      head.Version_Code = 8091_9999 // single APK
   } else {
      head.Version_Code = 9999_9999
   }
   val, err := decode_query(res.Body)
   if err != nil {
      return nil, err
   }
   head.Auth = val.Get("Auth")
   return &head, nil
}

type Header struct {
   Device_ID uint64 // X-DFE-Device-ID
   SDK int64 // User-Agent
   Version_Code int64 // User-Agent
   Auth string // Authorization
}

func (h Header) Set_Agent(head http.Header) {
   var buf []byte
   buf = append(buf, "Android-Finsky (sdk="...)
   buf = strconv.AppendInt(buf, h.SDK, 10)
   buf = append(buf, ",versionCode="...)
   buf = strconv.AppendInt(buf, h.Version_Code, 10)
   buf = append(buf, ')')
   head.Set("User-Agent", string(buf))
}

func (h Header) Set_Auth(head http.Header) {
   head.Set("Authorization", "Bearer " + h.Auth)
}

func (h Header) Set_Device(head http.Header) {
   device := strconv.FormatUint(h.Device_ID, 16)
   head.Set("X-DFE-Device-ID", device)
}

// Purchase app. Only needs to be done once per Google account.
func (h Header) Purchase(app string) error {
   body := make(url.Values)
   body.Set("doc", app)
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/fdfe/purchase",
      strings.NewReader(body.Encode()),
   )
   if err != nil {
      return err
   }
   h.Set_Auth(req.Header)
   h.Set_Device(req.Header)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := Client.Do(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}

func (t Token) Token() string {
   return t.Get("Token")
}

// You can also use host "android.clients.google.com", but it also uses
// TLS fingerprinting.
func New_Token(email, password string) (*Token, error) {
   body := url.Values{
      "Email": {email},
      "Passwd": {password},
      "client_sig": {""},
      "droidguard_results": {""},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/auth", strings.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   hello, err := crypto.Parse_JA3(crypto.Android_API_26)
   if err != nil {
      return nil, err
   }
   tr := crypto.Transport(hello)
   res, err := Client.Transport(tr).Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   var tok Token
   tok.Values, err = decode_query(res.Body)
   if err != nil {
      return nil, err
   }
   return &tok, nil
}

type Token struct {
   url.Values
}
