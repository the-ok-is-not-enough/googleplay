package googleplay

import (
   "encoding/json"
   "github.com/89z/format"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

const (
   Sleep = 16 * time.Second
   agent = "Android-Finsky (sdk=99,versionCode=99999999)"
   origin = "https://android.clients.google.com"
)

var LogLevel format.LogLevel

// Purchase app. Only needs to be done once per Google account.
func (a Auth) Purchase(dev *Device, app string) error {
   query := "doc=" + url.QueryEscape(app)
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/purchase", strings.NewReader(query),
   )
   if err != nil {
      return err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Auth},
      "Content-Type": {"application/x-www-form-urlencoded"},
      "User-Agent": {agent},
      "X-DFE-Device-ID": {dev.String()},
   }
   LogLevel.Dump(req)
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
   LogLevel.Dump(req)
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
func (d *Device) Decode(src io.Reader) error {
   return json.NewDecoder(src).Decode(d)
}

// Write Device to file.
func (d Device) Encode(dst io.Writer) error {
   enc := json.NewEncoder(dst)
   enc.SetIndent("", " ")
   return enc.Encode(d)
}

func (d Device) String() string {
   return strconv.FormatInt(d.Android_ID, 16)
}

type NumDownloads struct {
   Value uint64
}

func (n NumDownloads) String() string {
   return format.Number.LabelUint(n.Value)
}

type Size struct {
   Value uint64
}

func (i Size) String() string {
   return format.Size.LabelUint(i.Value)
}

type response struct {
   *http.Response
}

func (r response) Error() string {
   return strconv.Itoa(r.StatusCode) + " " + r.Status
}

var purchaseRequired = response{
   &http.Response{StatusCode: 3, Status: "purchase required"},
}
