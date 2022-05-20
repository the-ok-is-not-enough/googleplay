package googleplay

import (
   "bufio"
   "bytes"
   "errors"
   "github.com/89z/format"
   "github.com/89z/format/crypto"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

const (
   DateInput = "Jan 2, 2006"
   DateOutput = "2006-01-02"
)

type Details struct {
   Title String
   Creator String
   UploadDate String // Jun 1, 2021
   VersionString String
   VersionCode Varint
   NumDownloads Varint
   Size Varint
   File []Varint
   Micros Varint
   CurrencyCode String
}

func (d Details) String() string {
   var buf []byte
   buf = append(buf, "Title: "...)
   buf = append(buf, d.Title...)
   buf = append(buf, "\nCreator: "...)
   buf = append(buf, d.Creator...)
   buf = append(buf, "\nUploadDate: "...)
   buf = append(buf, d.UploadDate...)
   buf = append(buf, "\nVersionString: "...)
   buf = append(buf, d.VersionString...)
   buf = append(buf, "\nVersionCode: "...)
   buf = strconv.AppendUint(buf, uint64(d.VersionCode), 10)
   buf = append(buf, "\nNumDownloads: "...)
   buf = append(buf, format.LabelNumber(d.NumDownloads)...)
   buf = append(buf, "\nSize: "...)
   buf = append(buf, format.LabelSize(d.Size)...)
   buf = append(buf, "\nFile:"...)
   for _, file := range d.File {
      if file == 0 {
         buf = append(buf, " APK"...)
      } else {
         buf = append(buf, " OBB"...)
      }
   }
   buf = append(buf, "\nOffer: "...)
   buf = strconv.AppendUint(buf, uint64(d.Micros), 10)
   buf = append(buf, ' ')
   buf = append(buf, d.CurrencyCode...)
   return string(buf)
}

type deviceConfiguration struct {
   app string
}

func (d deviceConfiguration) Error() string {
   var buf []byte
   buf = append(buf, "bad DeviceConfiguration for "...)
   buf = strconv.AppendQuote(buf, d.app)
   return string(buf)
}
type AppFileMetadata struct {
   FileType Varint
   DownloadURL String
}

func (d Delivery) Additional(typ Varint) string {
   var buf []byte
   if typ == 0 {
      buf = append(buf, "main"...)
   } else {
      buf = append(buf, "patch"...)
   }
   buf = append(buf, '.')
   buf = strconv.AppendUint(buf, d.VersionCode, 10)
   buf = append(buf, '.')
   buf = append(buf, d.PackageName...)
   buf = append(buf, ".obb"...)
   return string(buf)
}

type Delivery struct {
   DownloadURL String
   PackageName string
   SplitDeliveryData []SplitDeliveryData
   VersionCode uint64
   AdditionalFile []AppFileMetadata
}

type SplitDeliveryData struct {
   ID String
   DownloadURL String
}

func (d Delivery) Split(id String) string {
   var buf []byte
   buf = append(buf, d.PackageName...)
   buf = append(buf, '-')
   buf = append(buf, id...)
   buf = append(buf, '-')
   buf = strconv.AppendUint(buf, d.VersionCode, 10)
   buf = append(buf, ".apk"...)
   return string(buf)
}

func (d Delivery) Download() string {
   var buf []byte
   buf = append(buf, d.PackageName...)
   buf = append(buf, '-')
   buf = strconv.AppendUint(buf, d.VersionCode, 10)
   buf = append(buf, ".apk"...)
   return string(buf)
}

type Message = protobuf.Message

type String = protobuf.String

type Varint = protobuf.Varint

type Fixed64 = protobuf.Fixed64

const Sleep = 4 * time.Second

var LogLevel format.LogLevel

func parseQuery(query io.Reader) url.Values {
   vals := make(url.Values)
   buf := bufio.NewScanner(query)
   for buf.Scan() {
      key, val, ok := strings.Cut(buf.Text(), "=")
      if ok {
         vals.Add(key, val)
      }
   }
   return vals
}

type Token struct {
   Services string
   Token string
}

// You can also use host "android.clients.google.com", but it also uses
// TLS fingerprinting.
func NewToken(email, password string) (*Token, error) {
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
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   val := parseQuery(res.Body)
   var tok Token
   tok.Services = val.Get("services")
   tok.Token = val.Get("Token")
   return &tok, nil
}

func OpenToken(elem ...string) (*Token, error) {
   return format.Open[Token](elem...)
}

func (t Token) Create(elem ...string) error {
   return format.Create(t, elem...)
}

type NativePlatform map[int64]string

var Platforms = NativePlatform{
   // com.google.android.youtube
   0: "x86",
   // com.miui.weather2
   1: "armeabi-v7a",
   // com.kakaogames.twodin
   2: "arm64-v8a",
}

func (n NativePlatform) String() string {
   first := true
   var buf []byte
   for key, val := range n {
      if first {
         first = false
      } else {
         buf = append(buf, '\n')
      }
      buf = strconv.AppendInt(buf, key, 10)
      buf = append(buf, ' ')
      buf = append(buf, val...)
   }
   return string(buf)
}

type Device struct {
   AndroidID Fixed64
}

func OpenDevice(elem ...string) (*Device, error) {
   return format.Open[Device](elem...)
}

func (d Device) Create(elem ...string) error {
   return format.Create(d, elem...)
}

// These can use default values, but they must all be included
type Config struct {
   DeviceFeature []String
   GlEsVersion Varint
   GlExtension []String
   HasFiveWayNavigation Varint
   HasHardKeyboard Varint
   Keyboard Varint
   Navigation Varint
   ScreenDensity Varint
   ScreenLayout Varint
   SystemSharedLibrary []String
   TouchScreen Varint
}

var Phone = Config{
   DeviceFeature: []String{
      // app.source.getcontact
      "android.hardware.location.gps",
      // br.com.rodrigokolb.realdrum
      "android.software.midi",
      // com.app.xt
      "android.hardware.camera.front",
      // com.clearchannel.iheartradio.controller
      "android.hardware.microphone",
      // com.google.android.apps.walletnfcrel
      "android.software.device_admin",
      // com.google.android.youtube
      "android.hardware.touchscreen",
      "android.hardware.wifi",
      // com.illumix.fnafar
      "android.hardware.sensor.gyroscope",
      // com.madhead.tos.zh
      "android.hardware.sensor.accelerometer",
      // com.miHoYo.GenshinImpact
      "android.hardware.opengles.aep",
      // com.pinterest
      "android.hardware.camera",
      "android.hardware.location",
      "android.hardware.screen.portrait",
      // com.sygic.aura
      "android.hardware.location.network",
      // com.xiaomi.smarthome
      "android.hardware.bluetooth",
      "android.hardware.bluetooth_le",
      "android.hardware.camera.autofocus",
      "android.hardware.usb.host",
      // kr.sira.metal
      "android.hardware.sensor.compass",
      // org.thoughtcrime.securesms
      "android.hardware.telephony",
      // org.videolan.vlc
      "android.hardware.screen.landscape",
   },
   SystemSharedLibrary: []String{
      // com.amctve.amcfullepisodes
      "org.apache.http.legacy",
      // com.binance.dev
      "android.test.runner",
      // com.miui.weather2
      "global-miui11-empty.jar",
   },
   // com.axis.drawingdesk.v3
   GlEsVersion: 0x9_9999,
   GlExtension: []String{
      // com.instagram.android
      "GL_OES_compressed_ETC1_RGB8_texture",
      // com.kakaogames.twodin
      "GL_KHR_texture_compression_astc_ldr",
   },
   // com.valvesoftware.android.steam.community
   TouchScreen: 3,
}

// A Sleep is needed after this.
func (c Config) Checkin(platform string) (*Device, error) {
   checkin := Message{
      4: Message{ // checkin
         1: Message{ // build
            10: Varint(29), // sdkVersion
         },
      },
      14: Varint(3), // version
      18: Message{ // deviceConfiguration
         1: c.TouchScreen, // touchScreen
         2: c.Keyboard, // keyboard
         3: c.Navigation, // navigation
         4: c.ScreenLayout, // screenLayout
         5: c.HasHardKeyboard, // hasHardKeyboard
         6: c.HasFiveWayNavigation, // hasFiveWayNavigation
         7: c.ScreenDensity, // screenDensity
         8: c.GlEsVersion, // glEsVersion
         11: String(platform), // nativePlatform
      },
   }
   for _, library := range c.SystemSharedLibrary {
      checkin.Get(18).AddString(9, library)
   }
   for _, extension := range c.GlExtension {
      checkin.Get(18).AddString(15, extension)
   }
   for _, name := range c.DeviceFeature {
      // .deviceConfiguration.deviceFeature
      checkin.Get(18).Add(26, Message{1: name})
   }
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/checkin",
      bytes.NewReader(checkin.Marshal()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-protobuffer")
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   checkinResponse, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   var dev Device
   // .androidId
   dev.AndroidID, err = checkinResponse.GetFixed64(7)
   if err != nil {
      return nil, err
   }
   return &dev, nil
}
