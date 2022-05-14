package googleplay

import (
   "fmt"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
)


type Delivery struct {
   AdditionalFile String
   DownloadURL String
   SplitDeliveryData []SplitDeliveryData
}

func (d Delivery) Data() []SplitDeliveryData {
   datas := d.SplitDeliveryData
   data := SplitDeliveryData{DownloadURL: d.DownloadURL}
   return append(datas, data)
}

type SplitDeliveryData struct {
   ID String
   DownloadURL String
}

func (s SplitDeliveryData) Name(app string, ver uint64) string {
   if s.ID != "" {
      return fmt.Sprint(app, "-", s.ID, "-", ver, ".apk")
   }
   return fmt.Sprint(app, "-", ver, ".apk")
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
   // .payload.deliveryResponse.status
   status, err := responseWrapper.Get(1).Get(21).GetVarint(1)
   if err != nil {
      return nil, err
   }
   switch status {
   case 2:
      return nil, errorString("Geo-blocking")
   case 3:
      return nil, errorString("Purchase required")
   case 5:
      return nil, errorString("Invalid version")
   }
   // .payload.deliveryResponse.appDeliveryData
   appData := responseWrapper.Get(1).Get(21).Get(2)
   var del Delivery
   // downloadUrl
   del.DownloadURL, err = appData.GetString(3)
   if err != nil {
      return nil, err
   }
   // .splitDeliveryData
   for _, data := range appData.GetMessages(15) {
      var split SplitDeliveryData
      // .id
      split.ID, err = data.GetString(1)
      if err != nil {
         return nil, err
      }
      // .downloadUrl
      split.DownloadURL, err = data.GetString(5)
      if err != nil {
         return nil, err
      }
      del.SplitDeliveryData = append(del.SplitDeliveryData, split)
   }
   return &del, nil
}
