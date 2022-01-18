package googleplay

import (
   "encoding/json"
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
   "os"
   "strconv"
)

type Device struct {
   AndroidID uint64
}

// A Sleep is needed after this.
func NewDevice(con Config) (*Device, error) {
   checkinRequest := protobuf.Message{
      {4, "checkin"}: protobuf.Message{
         {1, "build"}: protobuf.Message{
            {10, "sdkVersion"}: uint64(29),
         },
      },
      {14, "version"}: uint64(3),
      {18, "deviceConfiguration"}: protobuf.Message{
         {1, "touchScreen"}: con.TouchScreen,
         {2, "keyboard"}: con.Keyboard,
         {3, "navigation"}: con.Navigation,
         {4, "screenLayout"}: con.ScreenLayout,
         {5, "hasHardKeyboard"}: con.HasHardKeyboard,
         {6, "hasFiveWayNavigation"}: con.HasFiveWayNavigation,
         {7, "screenDensity"}: con.ScreenDensity,
         {8, "glEsVersion"}: con.GLESversion,
         {9, "systemSharedLibrary"}: con.SystemSharedLibrary,
         {11, "nativePlatform"}: con.NativePlatform,
         {15, "glExtension"}: con.GLextension,
      },
   }
   for _, feature := range con.DeviceFeature {
      checkinRequest.Get(18, "deviceConfiguration").
      Add(26, "deviceFeature", protobuf.Message{
         {1, "name"}: feature,
      })
   }
   req, err := http.NewRequest(
      "POST", origin + "/checkin", checkinRequest.Encode(),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-protobuffer")
   format.Log.Dump(req)
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
   dev.AndroidID = checkinResponse.GetUint64(7, "androidId")
   return &dev, nil
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

func (d Device) String() string {
   return strconv.FormatUint(d.AndroidID, 16)
}

