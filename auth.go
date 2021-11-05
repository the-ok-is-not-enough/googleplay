package googleplay

import (
   "bytes"
   "github.com/segmentio/encoding/proto"
   "io"
   "net/http"
   "net/url"
)

type AppDetails struct {
   DeveloperName string `protobuf:"bytes,1"`
   VersionCode int32 `protobuf:"varint,3"`
   Version string `protobuf:"bytes,4"`
   InstallationSize int64 `protobuf:"varint,9"`
   DeveloperEmail string `protobuf:"bytes,11"`
}

type Auth struct {
   url.Values
}

// deviceID is Google Service Framework.
func (a Auth) Details(deviceID, app string) (*AppDetails, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Get("Auth")},
      "X-DFE-Device-ID": {deviceID},
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
   return wrap.appDetails(), nil
}

// This seems to return `StatusOK`, even with invalid requests, and the response
// body only contains a token, that doesnt seem to indicate success or failure.
// Only way I know to check, it to try the `deviceID` with a `details` request
// or similar. Also, after the POST, you need to wait at least 16 seconds
// before the `deviceID` can be used.
func (a Auth) Upload(deviceID string, dev Device) error {
   buf, err := proto.Marshal(dev)
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
      "User-Agent": {"Android-Finsky (sdk=99,versionCode=99999999)"},
      "X-DFE-Device-ID": {deviceID},
   }
   res, err := roundTrip(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}

type responseWrapper struct {
   Payload struct {
      DetailsResponse struct {
         DocV2 struct {
            DocumentDetails struct {
               AppDetails AppDetails `protobuf:"bytes,1"`
            } `protobuf:"bytes,13"`
         } `protobuf:"bytes,4"`
      } `protobuf:"bytes,2"`
   } `protobuf:"bytes,1"`
}

func (r responseWrapper) appDetails() *AppDetails {
   return &r.Payload.DetailsResponse.DocV2.DocumentDetails.AppDetails
}
