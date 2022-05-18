package googleplay

import (
   "bytes"
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
)

const (
   // com.kakaogames.twodin
   Arm64 String = "arm64-v8a"
   // com.miui.weather2
   Armeabi String = "armeabi-v7a"
   // com.google.android.youtube
   X86 String = "x86"
)

// These can use default values, but they must all be included
type Config struct {
   DeviceFeature []String
   GlEsVersion Varint
   GlExtension String
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
      // br.com.rodrigokolb.realdrum
      "android.software.midi",
      // com.clearchannel.iheartradio.controller
      "android.hardware.microphone",
      // com.google.android.apps.walletnfcrel
      "android.software.device_admin",
      // com.google.android.youtube
      "android.hardware.touchscreen",
      "android.hardware.wifi",
      // com.madhead.tos.zh
      "android.hardware.sensor.accelerometer",
      // com.pinterest
      "android.hardware.camera",
      "android.hardware.location",
      "android.hardware.screen.portrait",
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
   // com.axis.drawingdesk.v3
   GlEsVersion: 0x0003_0001,
   // com.kakaogames.twodin
   GlExtension: "GL_KHR_texture_compression_astc_ldr",
   SystemSharedLibrary: []String{
      // com.amctve.amcfullepisodes
      "org.apache.http.legacy",
      // com.miui.weather2
      "global-miui11-empty.jar",
   },
   // com.valvesoftware.android.steam.community
   TouchScreen: 3,
}

// A Sleep is needed after this.
func (c Config) Checkin(platform String) (*Device, error) {
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
         11: platform, // nativePlatform
         15: c.GlExtension, // glExtension
      },
   }
   for _, library := range c.SystemSharedLibrary {
      checkin.Get(18).AddString(9, library)
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
   // .timeMsec
   dev.TimeMsec, err = checkinResponse.GetVarint(3)
   if err != nil {
      return nil, err
   }
   return &dev, nil
}

type Device struct {
   AndroidID Fixed64
   TimeMsec Varint
}

func OpenDevice(elem ...string) (*Device, error) {
   return format.Open[Device](elem...)
}

func (d Device) Create(elem ...string) error {
   return format.Create(d, elem...)
}
