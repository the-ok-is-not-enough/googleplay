package googleplay

import (
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

type Header struct {
   http.Header
}

func (h Header) Delivery(app string, ver uint64) (*Delivery, error) {
   req, err := http.NewRequest(
      "GET", "https://play-fe.googleapis.com/fdfe/delivery", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header = h.Header
   req.URL.RawQuery = url.Values{
      "doc": {app},
      "vc": {strconv.FormatUint(ver, 10)},
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
   status := responseWrapper.Get(/* payload */ 1).
      Get(/* deliveryResponse */ 21).
      GetUint64(1)
   switch status {
   case 2:
      return nil, errorString("Geo-blocking")
   case 3:
      return nil, errorString("Purchase required")
   case 5:
      return nil, errorString("Invalid version")
   }
   appData := responseWrapper.Get(/* payload */ 1).
      Get(/* deliveryResponse */ 21).
      Get(/* appDeliveryData */ 2)
   var del Delivery
   del.DownloadURL = appData.GetString(3)
   for _, data := range appData.GetMessages(/* splitDeliveryData */ 15) {
      var split SplitDeliveryData
      split.ID = data.GetString(1)
      split.DownloadURL = data.GetString(5)
      del.SplitDeliveryData = append(del.SplitDeliveryData, split)
   }
   return &del, nil
}

func (h Header) Details(app string) (*Details, error) {
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/details", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header = h.Header
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
   docV2 := responseWrapper.Get(/* payload */ 1).
      Get(/* detailsResponse */ 2).
      Get(4)
   var det Details
   det.CurrencyCode = docV2.Get(/* offer */ 8).GetString(2)
   file := docV2.Get(/* details */ 13).Get(/* appDetails */ 1).GetMessages(17)
   det.Files = len(file)
   for _, mes := range docV2.GetMessages(/* image */ 10) {
      var image Image
      image.Type = mes.GetUint64(/* imageType */ 1)
      image.URL = mes.GetString(/* imageUrl */ 5)
      det.Images = append(det.Images, image)
   }
   det.Micros = docV2.Get(/* offer */ 8).GetUint64(1)
   det.NumDownloads = docV2.Get(/* details */ 13).
      Get(/* appDetails */ 1).
      GetUint64(70)
   det.Size = docV2.Get(/* details */ 13).
      Get(/* appDetails */ 1).
      Get(/* installDetails */ 34).
      GetUint64(2)
   det.Title = docV2.GetString(5)
   det.UploadDate = docV2.Get(/* details */ 13).
      Get(/* appDetails */ 1).
      GetString(16)
   det.VersionCode = docV2.Get(/* details */ 13).
      Get(/* appDetails */ 1).
      GetUint64(3)
   det.VersionString = docV2.Get(/* details */ 13).
      Get(/* appDetails */ 1).
      GetString(4)
   return &det, nil
}

// Purchase app. Only needs to be done once per Google account.
func (h Header) Purchase(app string) error {
   query := "doc=" + url.QueryEscape(app)
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/fdfe/purchase",
      strings.NewReader(query),
   )
   if err != nil {
      return err
   }
   h.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header = h.Header
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return errorString(res.Status)
   }
   return nil
}
