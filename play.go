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
   "strconv"
   "strings"
   "time"
)

func tag(number float64, name string) float64 {
   return number
}

const (
   Sleep = 4 * time.Second
   origin = "https://android.clients.google.com"
)

var LogLevel format.LogLevel

type Auth struct {
   Auth string
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
   buf := []byte("Title: ")
   buf = append(buf, d.Title...)
   buf = append(buf, "\nUploadDate: "...)
   buf = append(buf, d.UploadDate...)
   buf = append(buf, "\nVersionString: "...)
   buf = append(buf, d.VersionString...)
   buf = append(buf, "\nVersionCode: "...)
   buf = strconv.AppendUint(buf, d.VersionCode, 10)
   buf = append(buf, "\nNumDownloads: "...)
   buf = append(buf, format.Number.GetUint64(d.NumDownloads)...)
   buf = append(buf, "\nSize: "...)
   buf = append(buf, format.Size.GetUint64(d.Size)...)
   buf = append(buf, "\nOffer: "...)
   buf = strconv.AppendFloat(buf, float64(d.Micros)/1e6, 'f', 2, 64)
   buf = append(buf, ' ')
   buf = append(buf, d.CurrencyCode...)
   return string(buf)
}

type Document struct {
   ID string
   Title string
   Creator string
}

func (d Document) String() string {
   var buf strings.Builder
   buf.WriteString("ID: ")
   buf.WriteString(d.ID)
   buf.WriteString("\nTitle: ")
   buf.WriteString(d.Title)
   buf.WriteString("\nCreator: ")
   buf.WriteString(d.Creator)
   return buf.String()
}

type Review struct {
   Author string
   Comment string
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
      // Instead of the following two, you can use this:
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
   LogLevel.Dump(req)
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
   LogLevel.Dump(req)
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
   value string
}

func (n notFound) Error() string {
   return strconv.Quote(n.value) + " not found"
}


var DefaultConfig = Config{
   DeviceFeature: []string{
      // com.google.android.apps.walletnfcrel
      "android.software.device_admin",
      // com.google.android.youtube
      "android.hardware.touchscreen",
      "android.hardware.wifi",
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
      "android.hardware.bluetooth",
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
      // com.bbca.bbcafullepisodes
      "org.apache.http.legacy",
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

func (a Auth) Header(dev *Device) Header {
   return a.headerVersion(dev, 9999_9999)
}

func (a Auth) SingleAPK(dev *Device) Header {
   return a.headerVersion(dev, 8091_9999)
}

func (a Auth) headerVersion(dev *Device, version int64) Header {
   var val Header
   val.Header = make(http.Header)
   val.Set("Authorization", "Bearer " + a.Auth)
   // User-Agent is only needed with "/fdfe/details" for some apps, example:
   // com.xiaomi.smarthome
   buf := []byte("Android-Finsky (sdk=9,versionCode=")
   buf = strconv.AppendInt(buf, version, 10)
   val.Set("User-Agent", string(buf))
   id := strconv.FormatUint(dev.AndroidID, 16)
   val.Set("X-DFE-Device-ID", id)
   return val
}

type Header struct {
   http.Header
}

// Purchase app. Only needs to be done once per Google account.
func (h Header) Purchase(app string) error {
   query := "doc=" + url.QueryEscape(app)
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/purchase", strings.NewReader(query),
   )
   if err != nil {
      return err
   }
   h.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header = h.Header
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}
