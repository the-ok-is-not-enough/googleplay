package googleplay

import (
   "encoding/json"
   "fmt"
   "github.com/89z/format"
   "io"
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
   origin = "https://android.clients.google.com"
)

const androidKey =
   "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp" +
   "5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLN" +
   "WgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="

var DefaultConfig = Config{
   DeviceFeature: []string{
      // com.instagram.android
      "android.hardware.bluetooth",
      // com.xiaomi.smarthome
      "android.hardware.bluetooth_le",
      // com.pinterest
      "android.hardware.camera",
      // com.xiaomi.smarthome
      "android.hardware.camera.autofocus",
      // com.pinterest
      "android.hardware.location",
      // com.smarty.voomvoom
      "android.hardware.location.gps",
      // com.vimeo.android.videoapp
      "android.hardware.microphone",
      // org.videolan.vlc
      "android.hardware.screen.landscape",
      // com.pinterest
      "android.hardware.screen.portrait",
      // com.smarty.voomvoom
      "android.hardware.sensor.accelerometer",
      // org.thoughtcrime.securesms
      "android.hardware.telephony",
      // com.google.android.youtube
      "android.hardware.touchscreen",
      // com.xiaomi.smarthome
      "android.hardware.usb.host",
      // com.google.android.youtube
      "android.hardware.wifi",
   },
   // com.axis.drawingdesk.v3
   GLESversion: 0x0003_0001,
   GLextension: []string{
      // com.instagram.android
      "GL_OES_compressed_ETC1_RGB8_texture",
   },
   NativePlatform: []string{
      // com.vimeo.android.videoapp
      "x86",
      // com.axis.drawingdesk.v3
      "armeabi-v7a",
   },
   ScreenWidth: 1,
   SystemSharedLibrary: []string{
      // com.miui.weather2
      "global-miui11-empty.jar",
   },
   // com.valvesoftware.android.steam.community
   TouchScreen: 3,
}

var Log = format.Log{Writer: os.Stdout}

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
   Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}

type Config struct {
   DeviceFeature []string
   GLESversion uint64
   GLextension []string
   // this can be 0, but it must be included:
   HasFiveWayNavigation uint64
   // this can be 0, but it must be included:
   HasHardKeyboard uint64
   // this can be 0, but it must be included:
   Keyboard uint64
   NativePlatform []string
   // this can be 0, but it must be included:
   Navigation uint64
   // this can be 0, but it must be included:
   ScreenDensity uint64
   // this can be 0, but it must be included:
   ScreenLayout uint64
   ScreenWidth uint64
   SystemSharedLibrary []string
   // this can be 0, but it must be included:
   TouchScreen uint64
}

type Delivery struct {
   DownloadURL string
   SplitDeliveryData []SplitDeliveryData
}

type Details struct {
   Title string
   VersionString string
   VersionCode uint64
   NumDownloads uint64
   Size uint64
   Micros uint64
   CurrencyCode string
}

func (d Details) String() string {
   str := new(strings.Builder)
   fmt.Fprintln(str, "Title:", d.Title)
   fmt.Fprintln(str, "VersionString:", d.VersionString)
   fmt.Fprintln(str, "VersionCode:", d.VersionCode)
   fmt.Fprint(str, "NumDownloads: ")
   format.Number.LabelUint64(str, d.NumDownloads)
   fmt.Fprintln(str)
   fmt.Fprint(str, "Size: ")
   format.Size.LabelUint64(str, d.Size)
   fmt.Fprintln(str)
   fmt.Fprintf(str, "Offer: %.2f ", float64(d.Micros)/1_000_000)
   fmt.Fprint(str, d.CurrencyCode)
   return str.String()
}

type Device struct {
   AndroidID uint64
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
   return strconv.FormatUint(d.AndroidID, 16)
}

type SplitDeliveryData struct {
   ID string
   DownloadURL string
}

type nopSource struct{}

func (nopSource) Read(buf []byte) (int, error) {
   return len(buf), nil
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
   return strconv.Itoa(r.StatusCode) + " " + r.Status
}
