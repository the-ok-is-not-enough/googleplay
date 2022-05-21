package googleplay

import (
   "bytes"
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
   "strconv"
)

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
   var buf []byte
   buf = append(buf, "nativePlatform"...)
   for key, val := range n {
      buf = append(buf, '\n')
      buf = strconv.AppendInt(buf, key, 10)
      buf = append(buf, ": "...)
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
