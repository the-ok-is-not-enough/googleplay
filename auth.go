package googleplay

import (
   "bytes"
   "github.com/segmentio/encoding/proto"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "time"
)

const (
   Sleep = 16 * time.Second
   agent = "Android-Finsky (sdk=99,versionCode=99999999)"
)

type Auth struct {
   url.Values
}

func (a Auth) Delivery(dev *Device, app string, ver int) (*Delivery, error) {
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
   var wrap responseWrapper
   if err := proto.Unmarshal(buf, &wrap); err != nil {
      return nil, err
   }
   return &wrap.Payload.DeliveryResponse, nil
}

func (a Auth) Details(dev *Device, app string) (*Details, error) {
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
   var wrap responseWrapper
   if err := proto.Unmarshal(buf, &wrap); err != nil {
      return nil, err
   }
   return &wrap.Payload.DetailsResponse, nil
}

// This seems to return `StatusOK`, even with invalid requests, and the response
// body only contains a token, that doesnt seem to indicate success or failure.
// Only way I know to check, it to try the `deviceID` with a `details` request
// or similar. Also, after the POST, you need to wait at least 16 seconds
// before the `deviceID` can be used.
func (a Auth) Upload(dev *Device, con Config) error {
   buf, err := proto.Marshal(con)
   if err != nil {
      return err
   }
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

type Delivery struct {
   AppDeliveryData struct {
      DownloadUrl string `protobuf:"bytes,3"`
   } `protobuf:"bytes,2"`
}

type Details struct {
   DocV2 struct {
      DocumentDetails struct {
         AppDetails struct {
            DeveloperName string `protobuf:"bytes,1"`
            VersionCode int `protobuf:"varint,3"`
            Version string `protobuf:"bytes,4"`
            InstallationSize int `protobuf:"varint,9"`
            DeveloperEmail string `protobuf:"bytes,11"`
         } `protobuf:"bytes,1"`
      } `protobuf:"bytes,13"`
   } `protobuf:"bytes,4"`
}

type responseWrapper struct {
   Payload struct {
      DetailsResponse Details `protobuf:"bytes,2"`
      DeliveryResponse Delivery `protobuf:"bytes,21"`
   } `protobuf:"bytes,1"`
}
