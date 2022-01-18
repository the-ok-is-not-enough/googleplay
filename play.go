package googleplay

import (
   "bufio"
   "encoding/json"
   "fmt"
   "github.com/89z/format"
   "github.com/89z/format/crypto"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
   "strconv"
   "strings"
   "time"
)

const (
   Sleep = 4 * time.Second
   agent = "Android-Finsky (sdk=99,versionCode=99999999)"
   origin = "https://android.clients.google.com"
)

var DefaultConfig = Config{
   DeviceFeature: []string{
      // com.google.android.apps.walletnfcrel
      "android.software.device_admin",
      // com.google.android.youtube
      "android.hardware.touchscreen",
      "android.hardware.wifi",
      // com.instagram.android
      "android.hardware.bluetooth",
      // com.pinterest
      "android.hardware.camera",
      "android.hardware.location",
      "android.hardware.screen.portrait",
      // com.smarty.voomvoom
      "android.hardware.location.gps",
      "android.hardware.sensor.accelerometer",
      // com.tgc.sky.android
      "android.hardware.touchscreen.multitouch",
      "android.hardware.touchscreen.multitouch.distinct",
      "android.hardware.vulkan.level",
      "android.hardware.vulkan.version",
      // org.videolan.vlc
      "android.hardware.screen.landscape",
      // com.vimeo.android.videoapp
      "android.hardware.microphone",
      // com.xiaomi.smarthome
      "android.hardware.bluetooth_le",
      "android.hardware.camera.autofocus",
      "android.hardware.usb.host",
      // org.thoughtcrime.securesms
      "android.hardware.telephony",
      // se.pax.calima
      "android.hardware.location.network",
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
      // com.exnoa.misttraingirls
      "arm64-v8a",
   },
   SystemSharedLibrary: []string{
      // com.miui.weather2
      "global-miui11-empty.jar",
   },
   // com.valvesoftware.android.steam.community
   TouchScreen: 3,
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
   UploadDate string
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
   fmt.Fprintln(str, "UploadDate:", d.UploadDate)
   fmt.Fprintln(str, "VersionString:", d.VersionString)
   fmt.Fprintln(str, "VersionCode:", d.VersionCode)
   fmt.Fprint(str, "NumDownloads: ")
   format.Number.Uint64(str, d.NumDownloads)
   fmt.Fprintln(str)
   fmt.Fprint(str, "Size: ")
   format.Size.Uint64(str, d.Size)
   fmt.Fprintln(str)
   fmt.Fprintf(str, "Offer: %.2f ", float64(d.Micros)/1_000_000)
   fmt.Fprint(str, d.CurrencyCode)
   return str.String()
}

type SplitDeliveryData struct {
   ID string
   DownloadURL string
}

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
      return nil, errorString(res.Status)
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

type errorString string

func (e errorString) Error() string {
   return string(e)
}

type notFound struct {
   input string
}

func (n notFound) Error() string {
   return strconv.Quote(n.input) + " not found"
}
