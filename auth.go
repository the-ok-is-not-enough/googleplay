package googleplay

import (
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func itemNotFound(responseWrapper protobuf.Message) error {
   message := responseWrapper.Get(2, "commands").
      GetString(2, "displayErrorMessage")
   if message != "" {
      return errorString(message)
   }
   return nil
}

func itemNotOwned(responseWrapper protobuf.Message) error {
   status := responseWrapper.Get(1, "payload").
      Get(21, "deliveryResponse").
      GetUint64(1, "status")
   switch status {
   case 3:
      return errorString("Item not owned.")
   case 5:
      return errorString("Version not found.")
   }
   return nil
}

type Auth struct {
   Auth string
}

func (a Auth) Delivery(dev *Device, app string, ver int64) (*Delivery, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/delivery", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Auth},
      "User-Agent": {agent},
      "X-DFE-Device-ID": {dev.String()},
   }
   req.URL.RawQuery = url.Values{
      "doc": {app},
      "vc": {strconv.FormatInt(ver, 10)},
   }.Encode()
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   if err := itemNotOwned(responseWrapper); err != nil {
      return nil, err
   }
   var del Delivery
   deliveryData := responseWrapper.Get(1, "payload").
      Get(21, "deliveryResponse").
      Get(2, "appDeliveryData")
   del.DownloadURL = deliveryData.GetString(3, "downloadUrl")
   for _, split := range deliveryData.GetMessages(15, "splitDeliveryData") {
      var dSplit SplitDeliveryData
      dSplit.ID = split.GetString(1, "id")
      dSplit.DownloadURL = split.GetString(5, "downloadUrl")
      del.SplitDeliveryData = append(del.SplitDeliveryData, dSplit)
   }
   return &del, nil
}

func (a Auth) Details(dev *Device, app string) (*Details, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Auth},
      // This is needed for some apps, for example:
      // com.xiaomi.smarthome
      "User-Agent": {agent},
      "X-Dfe-Device-ID": {dev.String()},
   }
   req.URL.RawQuery = "doc=" + url.QueryEscape(app)
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   if err := itemNotFound(responseWrapper); err != nil {
      return nil, err
   }
   var det Details
   docV2 := responseWrapper.Get(1, "payload").
      Get(2, "detailsResponse").
      Get(4, "docV2")
   det.CurrencyCode = docV2.Get(8, "offer").GetString(2, "currencyCode")
   det.Micros = docV2.Get(8, "offer").GetUint64(1, "micros")
   det.NumDownloads = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetUint64(70, "numDownloads")
   // The shorter path 13,1,9 returns wrong size for some packages:
   // com.riotgames.league.wildriftvn
   det.Size = docV2.Get(13, "details").
      Get(1, "appDetails").
      Get(34, "installDetails").
      GetUint64(2, "size")
   det.Title = docV2.GetString(5, "title")
   det.UploadDate = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetString(16, "uploadDate")
   det.VersionCode = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetUint64(3, "versionCode")
   det.VersionString = docV2.Get(13, "details").
      Get(1, "appDetails").
      GetString(4, "versionString")
   return &det, nil
}

// Purchase app. Only needs to be done once per Google account.
func (a Auth) Purchase(dev *Device, app string) error {
   query := "doc=" + url.QueryEscape(app)
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/purchase", strings.NewReader(query),
   )
   if err != nil {
      return err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + a.Auth},
      "Content-Type": {"application/x-www-form-urlencoded"},
      "User-Agent": {agent},
      "X-DFE-Device-ID": {dev.String()},
   }
   format.Log.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   return res.Body.Close()
}

