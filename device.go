package googleplay

import (
   "bytes"
   "github.com/89z/format/protobuf"
   "net/http"
)

// These can use default values, but they must all be included
type Config struct {
   DeviceFeature []protobuf.String
   GLESversion protobuf.Uint64
   GLextension protobuf.String
   HasFiveWayNavigation protobuf.Uint64
   HasHardKeyboard protobuf.Uint64
   Keyboard protobuf.Uint64
   NativePlatform []protobuf.String
   Navigation protobuf.Uint64
   ScreenDensity protobuf.Uint64
   ScreenLayout protobuf.Uint64
   SystemSharedLibrary protobuf.String
   TouchScreen protobuf.Uint64
}

var DefaultConfig = Config{
   DeviceFeature: []protobuf.String{
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
   NativePlatform: []protobuf.String{
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
   checkin := protobuf.Message{
      /* checkin */ 4: protobuf.Message{
         /* build */ 1: protobuf.Message{
            /* sdkVersion */ 10: protobuf.Uint64(29),
         },
      },
      /* version */ 14: protobuf.Uint64(3),
      /* deviceConfiguration */ 18: protobuf.Message{
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
   for _, each := range c.NativePlatform {
      checkin.Get(18).AddString(11, each)
   }
   for _, each := range c.DeviceFeature {
      checkin.Get(18).Add(26, protobuf.Message{1: each})
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
   dev.AndroidID = checkinResponse.GetUint64(7)
   return &dev, nil
}

type Device struct {
   AndroidID protobuf.Uint64
}

func OpenDevice(elem ...string) (*Device, error) {
   dev := new(Device)
   err := decode(dev, elem...)
   if err != nil {
      return nil, err
   }
   return dev, nil
}

func (d Device) Create(elem ...string) error {
   return encode(d, elem...)
}
