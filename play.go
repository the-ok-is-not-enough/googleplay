package googleplay

import (
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

func Open_Auth(name string) (*Auth, error) {
   query, err := os.ReadFile(name)
   if err != nil {
      return nil, err
   }
   var auth Auth
   auth.Values = parse_query(string(query))
   return &auth, nil
}

func (a Auth) Create(name string) error {
   query := format_query(a.Values)
   return format.WriteFile(name, []byte(query))
}

// this beats "io.Reader", and also "bytes.Fields"
func parse_query(query string) url.Values {
   vals := make(url.Values)
   for _, field := range strings.Fields(query) {
      key, val, ok := strings.Cut(field, "=")
      if ok {
         vals.Add(key, val)
      }
   }
   return vals
}

func format_query(vals url.Values) string {
   var buf strings.Builder
   for key := range vals {
      val := vals.Get(key)
      buf.WriteString(key)
      buf.WriteByte('=')
      buf.WriteString(val)
      buf.WriteByte('\n')
   }
   return buf.String()
}

const Sleep = 4 * time.Second

var Client = http.Default_Client

// You can also use host "android.clients.google.com", but it also uses
// TLS fingerprinting.
func New_Auth(email, password string) (*Auth, error) {
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
   query, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   var auth Auth
   auth.Values = parse_query(string(query))
   return &auth, nil
}

type Auth struct {
   url.Values
}

func (a *Auth) Exchange() error {
   // these values take from Android API 28
   body := url.Values{
      "Token": {a.Token()},
      "app": {"com.android.vending"},
      "client_sig": {"38918a453d07199354f8b19af05ec6562ced5788"},
      "service": {"oauth2:https://www.googleapis.com/auth/googleplay"},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/auth", strings.NewReader(body),
   )
   if err != nil {
      return err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := Client.Do(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   query, err := io.ReadAll(res.Body)
   if err != nil {
      return err
   }
   a.Values = parse_query(string(query))
   return nil
}

func (a Auth) Token() string {
   return a.Get("Token")
}

func (a Auth) Auth() string {
   return a.Get("Auth")
}

type Header struct {
   http.Header
}

func (a Auth) Header(device_id uint64, single bool) Header {
   var (
      buf []byte
      head Header
   )
   head.Header = make(http.Header)
   // Authorization
   head.Set("Authorization", "Bearer " + a.Auth())
   // X-DFE-Device-ID
   head.Set("X-DFE-Device-ID", strconv.FormatUint(device_id, 16))
   // User-Agent
   buf = append(buf, "Android-Finsky (sdk=9,versionCode="...)
   if single {
      buf = strconv.AppendInt(buf, 8091_9999, 10)
   } else {
      buf = strconv.AppendInt(buf, 9999_9999, 10)
   }
   buf = append(buf, ')')
   head.Set("User-Agent", string(buf))
   return head
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
   req.Header = h.Header
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := Client.Do(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}
