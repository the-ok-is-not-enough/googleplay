package googleplay

import (
   "bytes"
   "github.com/89z/format/protobuf"
   "net/http"
)

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

var DefaultConfig = Config{
   DeviceFeature: []string{
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

// A Sleep is needed after this.
func (c Config) Checkin() (*Device, error) {
   checkinRequest := protobuf.Message{
      /* checkin */ 4: protobuf.Message{
         /* build */ 1: protobuf.Message{
            /* sdkVersion */ 10: uint64(29),
         },
      },
      /* version */ 14: uint64(3),
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
         11: c.NativePlatform,
         15: c.GLextension,
      },
   }
   deviceConfiguration := checkinRequest.Get(18)
   for _, name := range c.DeviceFeature {
      deviceFeature := protobuf.Message{1: name}
      deviceConfiguration.Add(26, deviceFeature)
   }
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/checkin",
      bytes.NewReader(checkinRequest.Marshal()),
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
   checkin, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   var dev Device
   dev.AndroidID = checkin.GetUint64(7)
   return &dev, nil
}

type Device struct {
   AndroidID uint64
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
