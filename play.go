package googleplay

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strconv"
   "strings"
   "time"
)

const (
   Sleep = 16 * time.Second
   agent = "Android-Finsky (sdk=99,versionCode=99999999)"
   origin = "https://android.clients.google.com"
)

var LogLevel logLevel

func numberFormat(val float64, metric []string) string {
   var key int
   for val >= 1000 {
      val /= 1000
      key++
   }
   if key >= len(metric) {
      return ""
   }
   return strconv.FormatFloat(val, 'f', 3, 64) + metric[key]
}

// Purchase app. Only needs to be done once per Google account.
func (a Auth) Purchase(dev *Device, app string) error {
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/purchase", values{"doc": app}.reader(),
   )
   if err != nil {
      return err
   }
   val := make(values)
   val["Authorization"] = "Bearer " + a.Auth
   val["Content-Type"] = "application/x-www-form-urlencoded"
   val["User-Agent"] = agent
   val["X-DFE-Device-ID"] = dev.String()
   req.Header = val.header()
   LogLevel.dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}

type Device struct {
   Android_ID int64
}

func NewDevice() (*Device, error) {
   req, err := http.NewRequest(
      "POST", origin + "/checkin",
      strings.NewReader(`{"checkin":{},"version":3}`),
   )
   if err != nil {
      return nil, err
   }
   LogLevel.dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   dev := new(Device)
   if err := json.NewDecoder(res.Body).Decode(dev); err != nil {
      return nil, err
   }
   return dev, nil
}

// Read Device from file.
func (d *Device) Decode(r io.Reader) error {
   return json.NewDecoder(r).Decode(d)
}

// Write Device to file.
func (d Device) Encode(w io.Writer) error {
   enc := json.NewEncoder(w)
   enc.SetIndent("", " ")
   return enc.Encode(d)
}

func (d Device) String() string {
   return strconv.FormatInt(d.Android_ID, 16)
}

type logLevel int

func (l logLevel) dump(req *http.Request) error {
   switch l {
   case 0:
      fmt.Println(req.Method, req.URL)
   case 1:
      buf, err := httputil.DumpRequest(req, true)
      if err != nil {
         return err
      }
      os.Stdout.Write(buf)
      if !bytes.HasSuffix(buf, []byte{'\n'}) {
         os.Stdout.WriteString("\n")
      }
   case 2:
      buf, err := httputil.DumpRequestOut(req, true)
      if err != nil {
         return err
      }
      os.Stdout.Write(buf)
   }
   return nil
}

type response struct {
   code uint64
   status string
}

func (r response) Error() string {
   return strconv.FormatUint(r.code, 10) + " " + r.status
}

type values map[string]string

func (v values) encode() string {
   vals := make(url.Values)
   for key, val := range v {
      vals.Set(key, val)
   }
   return vals.Encode()
}

func (v values) header() http.Header {
   vals := make(http.Header)
   for key, val := range v {
      vals.Set(key, val)
   }
   return vals
}

func (v values) reader() io.Reader {
   enc := v.encode()
   return strings.NewReader(enc)
}
