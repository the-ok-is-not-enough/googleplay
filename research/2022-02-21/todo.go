package googleplay

import (
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
)

func (h Header) Delivery(app string, ver int64) (*Delivery, error) {
   req, err := http.NewRequest(
      "GET", "https://play-fe.googleapis.com/fdfe/delivery", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header = h.Header
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
   status := responseWrapper.Get(1, "payload").
      Get(21, "deliveryResponse").
      GetVarint(1, "status")
   switch status {
   case 2:
      return nil, errorString("Geo-blocking")
   case 3:
      return nil, errorString("Purchase required")
   case 5:
      return nil, errorString("Invalid version")
   }
   appData := responseWrapper.Get(1, "payload").
      Get(21, "deliveryResponse").
      Get(2, "appDeliveryData")
   var del Delivery
   del.DownloadURL = appData.GetString(3, "downloadUrl")
   for _, data := range appData.GetMessages(15, "splitDeliveryData") {
      var split SplitDeliveryData
      split.ID = data.GetString(1, "id")
      split.DownloadURL = data.GetString(5, "downloadUrl")
      del.SplitDeliveryData = append(del.SplitDeliveryData, split)
   }
   return &del, nil
}

func (a Auth) Header(dev *Device) Header {
   return a.headerVersion(dev, 9999_9999)
}

func (a Auth) SingleAPK(dev *Device) Header {
   return a.headerVersion(dev, 8091_9999)
}

func (a Auth) headerVersion(dev *Device, version int64) Header {
   var val Header
   val.Header = make(http.Header)
   val.Set("Authorization", "Bearer " + a.Auth)
   buf := []byte("Android-Finsky (sdk=9,versionCode=")
   buf = strconv.AppendInt(buf, version, 10)
   val.Set("User-Agent", string(buf))
   id := strconv.FormatUint(dev.AndroidID, 16)
   val.Set("X-DFE-Device-ID", id)
   return val
}

type Header struct {
   http.Header
}
