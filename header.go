package googleplay

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strings"
)

type File struct {
   Size Varint
   VersionCode Varint
}

type Header struct {
   http.Header
}

func (h Header) Delivery(app string, ver Varint) (*Delivery, error) {
   req, err := http.NewRequest(
      "GET", "https://play-fe.googleapis.com/fdfe/delivery", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header = h.Header
   req.URL.RawQuery = url.Values{
      "doc": {app},
      "vc": {fmt.Sprint(ver)},
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
      GetVarint(1)
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
   det.Micros = docV2.Get(/* offer */ 8).GetVarint(1)
   det.NumDownloads = docV2.Get(/* details */ 13).
      Get(/* appDetails */ 1).
      GetVarint(70)
   // The shorter path 13,1,9 returns wrong size for some packages:
   // com.riotgames.league.wildriftvn
   det.Size = docV2.Get(/* details */ 13).
      Get(/* appDetails */ 1).
      Get(/* installDetails */ 34).
      GetVarint(2)
   det.Title = docV2.GetString(5)
   det.UploadDate = docV2.Get(/* details */ 13).
      Get(/* appDetails */ 1).
      GetString(16)
   det.VersionCode = docV2.Get(/* details */ 13).
      Get(/* appDetails */ 1).
      GetVarint(3)
   det.VersionString = docV2.Get(/* details */ 13).
      Get(/* appDetails */ 1).
      GetString(4)
   // file
   file := docV2.Get(/* details */ 13).Get(/* appDetails */ 1).GetMessages(17)
   det.Files = len(file)
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
