package googleplay
// github.com/89z

import (
   "errors"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
)

type AppFileMetadata struct {
   FileType uint64
   DownloadURL string
}

type Delivery struct {
   DownloadURL string
   PackageName string
   SplitDeliveryData []SplitDeliveryData
   VersionCode uint64
   AdditionalFile []AppFileMetadata
}

func (d Delivery) Additional(typ uint64) string {
   var buf []byte
   if typ == 0 {
      buf = append(buf, "main"...)
   } else {
      buf = append(buf, "patch"...)
   }
   buf = append(buf, '.')
   buf = strconv.AppendUint(buf, d.VersionCode, 10)
   buf = append(buf, '.')
   buf = append(buf, d.PackageName...)
   buf = append(buf, ".obb"...)
   return string(buf)
}

func (d Delivery) Download() string {
   var buf []byte
   buf = append(buf, d.PackageName...)
   buf = append(buf, '-')
   buf = strconv.AppendUint(buf, d.VersionCode, 10)
   buf = append(buf, ".apk"...)
   return string(buf)
}

func (d Delivery) Split(id string) string {
   var buf []byte
   buf = append(buf, d.PackageName...)
   buf = append(buf, '-')
   buf = append(buf, id...)
   buf = append(buf, '-')
   buf = strconv.AppendUint(buf, d.VersionCode, 10)
   buf = append(buf, ".apk"...)
   return string(buf)
}

func (h Header) Delivery(app string, ver uint64) (*Delivery, error) {
   req, err := http.NewRequest(
      "GET", "https://play-fe.googleapis.com/fdfe/delivery", nil,
   )
   if err != nil {
      return nil, err
   }
   h.SetAgent(req.Header)
   h.SetAuth(req.Header) // needed for single APK
   h.SetDevice(req.Header)
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
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
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
      return nil, errors.New("geo-blocking")
   case 3:
      return nil, errors.New("purchase required")
   case 5:
      return nil, errors.New("invalid version")
   }
   // .payload.deliveryResponse.appDeliveryData
   appData := responseWrapper.Get(1).Get(21).Get(2)
   var del Delivery
   // .downloadUrl
   del.DownloadURL, err = appData.GetString(3)
   if err != nil {
      return nil, err
   }
   del.PackageName = app
   del.VersionCode = ver
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
   // .additionalFile
   for _, file := range appData.GetMessages(4) {
      var app AppFileMetadata
      // .fileType
      app.FileType, err = file.GetVarint(1)
      if err != nil {
         return nil, err
      }
      // .downloadUrl
      app.DownloadURL, err = file.GetString(4)
      if err != nil {
         return nil, err
      }
      del.AdditionalFile = append(del.AdditionalFile, app)
   }
   return &del, nil
}

type SplitDeliveryData struct {
   ID string
   DownloadURL string
}
