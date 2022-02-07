package googleplay

import (
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func deliveryResponse(responseWrapper protobuf.Message) error {
   status := responseWrapper.Get(1, "payload").
      Get(21, "deliveryResponse").
      GetVarint(1, "status")
   switch status {
   case 2:
      return errorString("Regional lockout")
   case 3:
      return errorString("Purchase required")
   case 5:
      return errorString("Invalid version")
   }
   return nil
}

type Auth struct {
   Auth string
}

func (a Auth) AndroidPackage(dev *Device) Header {
   return Header{
      "Authorization": {"Bearer " + a.Auth},
      "User-Agent": {"Android-Finsky (sdk=9,versionCode=80919999)"},
      "X-DFE-Device-ID": {dev.String()},
   }
}

func (a Auth) AppBundle(dev *Device) Header {
   return Header{
      "Authorization": {"Bearer " + a.Auth},
      // "/fdfe/details" only needs User-Agent for some apps, for example:
      // com.xiaomi.smarthome
      "User-Agent": {"Android-Finsky (sdk=99,versionCode=99999999)"},
      "X-DFE-Device-ID": {dev.String()},
   }
}

type Header http.Header

func (h Header) Delivery(app string, ver int64) (*Delivery, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/delivery", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header(h)
   req.URL.RawQuery = url.Values{
      "doc": {app},
      "vc": {strconv.FormatInt(ver, 10)},
   }.Encode()
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   if err := deliveryResponse(responseWrapper); err != nil {
      return nil, err
   }
   var del Delivery
   deliveryData := responseWrapper.Get(1, "payload").
      Get(21, "deliveryResponse").
      Get(2, "appDeliveryData")
   del.DownloadURL = deliveryData.GetString(3, "downloadUrl")
   for _, split := range deliveryData.GetMessages(15, "splitDeliveryData") {
      var data SplitDeliveryData
      data.ID = split.GetString(1, "id")
      data.DownloadURL = split.GetString(5, "downloadUrl")
      del.SplitDeliveryData = append(del.SplitDeliveryData, data)
   }
   return &del, nil
}

type Delivery struct {
   DownloadURL string
   SplitDeliveryData []SplitDeliveryData
}

func (h Header) Details(app string) (*Details, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header(h)
   req.URL.RawQuery = "doc=" + url.QueryEscape(app)
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != http.StatusOK {
      return nil, errorString(res.Status)
   }
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   var det Details
   docV2 := responseWrapper.Get(1, "payload").
      Get(2, "detailsResponse").
      Get(4, "docV2")
   det.CurrencyCode = docV2.Get(8, "offer").GetString(2, "currencyCode")
   det.Micros = docV2.Get(8, "offer").GetVarint(1, "micros")
   det.NumDownloads = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetVarint(70, "numDownloads")
   // The shorter path 13,1,9 returns wrong size for some packages:
   // com.riotgames.league.wildriftvn
   det.Size = docV2.Get(13, "details").
      Get(1, "appDetails").
      Get(34, "installDetails").
      GetVarint(2, "size")
   det.Title = docV2.GetString(5, "title")
   det.UploadDate = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetString(16, "uploadDate")
   det.VersionCode = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetVarint(3, "versionCode")
   det.VersionString = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetString(4, "versionString")
   return &det, nil
}

// Purchase app. Only needs to be done once per Google account.
func (h Header) Purchase(app string) error {
   query := "doc=" + url.QueryEscape(app)
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/purchase", strings.NewReader(query),
   )
   if err != nil {
      return err
   }
   h["Content-Type"] = []string{"application/x-www-form-urlencoded"}
   req.Header = http.Header(h)
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}
