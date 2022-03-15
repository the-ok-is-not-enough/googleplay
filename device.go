package googleplay

import (
   "bytes"
   "github.com/89z/format/protobuf"
   "net/http"
)

var tag = protobuf.NewTag

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
   body := message{
      tag(4, "checkin"): message{
         tag(1, "build"): message{
            tag(10, "sdkVersion"): uint64(29),
         },
      },
      tag(14, "version"): uint64(3),
      tag(18, "deviceConfiguration"): message{
         tag(1, "touchScreen"): c.TouchScreen,
         tag(2, "keyboard"): c.Keyboard,
         tag(3, "navigation"): c.Navigation,
         tag(4, "screenLayout"): c.ScreenLayout,
         tag(5, "hasHardKeyboard"): c.HasHardKeyboard,
         tag(6, "hasFiveWayNavigation"): c.HasFiveWayNavigation,
         tag(7, "screenDensity"): c.ScreenDensity,
         tag(8, "glEsVersion"): c.GLESversion,
         tag(9, "systemSharedLibrary"): c.SystemSharedLibrary,
         tag(11, "nativePlatform"): c.NativePlatform,
         tag(15, "glExtension"): c.GLextension,
      },
   }
   config := body.Get(18, "deviceConfiguration")
   for _, name := range c.DeviceFeature {
      feature := message{
         tag(1, "name"): name,
      }
      config.Add(26, "deviceFeature", feature)
   }
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/checkin",
      bytes.NewReader(body.Marshal()),
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
   dev.AndroidID = checkinResponse.GetFixed64(7, "androidId")
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

type message = protobuf.Message
