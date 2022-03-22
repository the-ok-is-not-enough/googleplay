package googleplay

import (
   "bytes"
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
)

// These can use default values, but they must all be included
type Config struct {
   DeviceFeature []String
   GLESversion Varint
   GLextension String
   HasFiveWayNavigation Varint
   HasHardKeyboard Varint
   Keyboard Varint
   NativePlatform []String
   Navigation Varint
   ScreenDensity Varint
   ScreenLayout Varint
   SystemSharedLibrary String
   TouchScreen Varint
}

var DefaultConfig = Config{
   DeviceFeature: []String{
      // com.google.android.GoogleCamera
      "android.hardware.camera.level.full",
      "com.google.android.feature.GOOGLE_EXPERIENCE",
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
   // com.instagram.android
   GLextension: "GL_OES_compressed_ETC1_RGB8_texture",
   NativePlatform: []String{
      // com.vimeo.android.videoapp
      "x86",
      // com.axis.drawingdesk.v3
      "armeabi-v7a",
      // com.exnoa.misttraingirls
      "arm64-v8a",
   },
   // com.miui.weather2
   SystemSharedLibrary: "global-miui11-empty.jar",
   // com.valvesoftware.android.steam.community
   TouchScreen: 3,
}

// A Sleep is needed after this.
func (c Config) Checkin() (*Device, error) {
   checkin := Message{
      /* checkin */ 4: Message{
         /* build */ 1: Message{
            /* sdkVersion */ 10: Varint(29),
         },
      },
      /* version */ 14: Varint(3),
      /* deviceConfiguration */ 18: Message{
         1: c.TouchScreen,
         2: c.Keyboard,
         3: c.Navigation,
         4: c.ScreenLayout,
         5: c.HasHardKeyboard,
         6: c.HasFiveWayNavigation,
         7: c.ScreenDensity,
         8: c.GLESversion,
         9: c.SystemSharedLibrary,
         15: c.GLextension,
      },
   }
   for _, platform := range c.NativePlatform {
      checkin.Get(18).AddString(11, platform)
   }
   for _, name := range c.DeviceFeature {
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
   dev.AndroidID = checkinResponse.GetFixed64(7)
   dev.TimeMsec = checkinResponse.GetVarint(3)
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

type Fixed64 = protobuf.Fixed64

type Message = protobuf.Message

type String = protobuf.String

type Varint = protobuf.Varint
