package googleplay

import (
   "encoding/json"
   "github.com/89z/format"
   "net/http"
   "net/url"
   "os"
   "strconv"
   "strings"
   "time"
)

const (
   Sleep = 4 * time.Second
   agent = "Android-Finsky (sdk=99,versionCode=99999999)"
)

var purchaseRequired = response{
   &http.Response{StatusCode: 3, Status: "purchase required"},
}

type Auth struct {
   Auth string
}

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
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}

type Delivery struct {
   DownloadURL string
   SplitDeliveryData []SplitDeliveryData
}

type Device struct {
   AndroidID uint64
}

func OpenDevice(name string) (*Device, error) {
   file, err := os.Open(name)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   dev := new(Device)
   if err := json.NewDecoder(file).Decode(dev); err != nil {
      return nil, err
   }
   return dev, nil
}

func (d Device) Create(name string) error {
   file, err := os.Create(name)
   if err != nil {
      return err
   }
   defer file.Close()
   enc := json.NewEncoder(file)
   enc.SetIndent("", " ")
   return enc.Encode(d)
}

func (d Device) String() string {
   return strconv.FormatUint(d.AndroidID, 16)
}

type SplitDeliveryData struct {
   ID string
   DownloadURL string
}

type notFound struct {
   input string
}

func (n notFound) Error() string {
   return strconv.Quote(n.input) + " not found"
}

type response struct {
   *http.Response
}

func (r response) Error() string {
   // Status includes both
   return r.Status
}
