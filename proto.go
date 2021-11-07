package googleplay

import (
   "bytes"
   "github.com/89z/parse/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

func (a Auth) Delivery(dev *Device, app string, ver int) (protobuf.Message, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/delivery", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Get("Auth")},
      "User-Agent": {agent},
      "X-DFE-Device-ID": {dev.String()},
   }
   req.URL.RawQuery = url.Values{
      "doc": {app},
      "vc": {strconv.Itoa(ver)},
   }.Encode()
   res, err := roundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return protobuf.Parse(buf), nil
}

// This seems to return `StatusOK`, even with invalid requests, and the response
// body only contains a token, that doesnt seem to indicate success or failure.
// Only way I know to check, it to try the `deviceID` with a `details` request
// or similar. Also, after the POST, you need to wait at least 16 seconds
// before the `deviceID` can be used.
func (a Auth) Upload(dev *Device, config protobuf.Message) error {
   buf := config.Marshal()
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/uploadDeviceConfig", bytes.NewReader(buf),
   )
   if err != nil {
      return err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Get("Auth")},
      "User-Agent": {agent},
      "X-DFE-Device-ID": {dev.String()},
   }
   res, err := roundTrip(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}

var DefaultConfig = protobuf.Message{
   deviceConfiguration: protobuf.Message{
      touchScreen: int32(1),
      keyboard: int32(1),
      navigation: int32(1),
      screenLayout: int32(1),
      hasHardKeyboard: true,
      hasFiveWayNavigation: true,
      screenDensity: int32(1),
      // developer.android.com/guide/topics/manifest/uses-feature-element
      glEsVersion: int32(0x0009_0000),
      // developer.android.com/guide/topics/manifest/uses-feature-element
      systemAvailableFeature: protobuf.Repeated{
         // com.pinterest
         "android.hardware.camera",
         // com.pinterest
         "android.hardware.faketouch",
         // com.pinterest
         "android.hardware.location",
         // com.pinterest
         "android.hardware.screen.portrait",
         // com.google.android.youtube
         "android.hardware.touchscreen",
         // com.google.android.youtube
         "android.hardware.wifi",
      },
      // developer.android.com/ndk/guides/abis
      nativePlatform: protobuf.Repeated{
         "armeabi-v7a",
      },
   },
}

func (a Auth) Details(dev *Device, app string) (protobuf.Message, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Get("Auth")},
      "X-DFE-Device-ID": {dev.String()},
   }
   req.URL.RawQuery = url.Values{
      "doc": {app},
   }.Encode()
   res, err := roundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return protobuf.Parse(buf), nil
}

type responseWrapper struct {
   protobuf.Message
}

func (r responseWrapper) payload() protobuf.Message {
   return r.Message.Message(1)
}

type payload struct {
   protobuf.Message
}

func (p payload) detailsResponse() protobuf.Message {
   return p.Message.Message(2)
}

type detailsResponse struct {
   protobuf.Message
}

func (d detailsResponse) docV2() protobuf.Message {
   return d.Message.Message(4)
}

type docV2 struct {
   protobuf.Message
}

func (d docV2) details() protobuf.Message {
   return d.Message.Message(13)
}

type documentDetails struct {
   protobuf.Message
}

func (d documentDetails) appDetails() protobuf.Message {
   return d.Message.Message(1)
}
