package googleplay

import (
   "errors"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func (h Header) SetAuth(head http.Header) {
   head.Set("Authorization", "Bearer " + h.Auth)
}

func (h Header) SetDevice(head http.Header) {
   device := strconv.FormatUint(h.AndroidID, 16)
   head.Set("X-DFE-Device-ID", device)
}

func (h Header) Delivery(app string, ver uint64) (*Delivery, error) {
   req, err := http.NewRequest(
      "GET", "https://play-fe.googleapis.com/fdfe/delivery", nil,
   )
   if err != nil {
      return nil, err
   }
   h.SetAgent(req.Header)
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

type Header struct {
   AndroidID uint64
   SDK int64
   VersionCode uint64
   Auth string
}

func (h Header) Details(app string) (*Details, error) {
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/details", nil,
   )
   if err != nil {
      return nil, err
   }
   h.SetAgent(req.Header) // app.source.getcontact
   h.SetAuth(req.Header)
   h.SetDevice(req.Header)
   req.URL.RawQuery = "doc=" + url.QueryEscape(app)
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   // .payload.detailsResponse.docV2
   docV2 := responseWrapper.Get(1).Get(2).Get(4)
   var det Details
   // The following fields will fail with wrong ABI, so try them first. If the
   // first one passes, then use native error for the rest.
   // .details.appDetails.versionCode
   det.VersionCode, err = docV2.Get(13).Get(1).GetVarint(3)
   if err != nil {
      return nil, deviceConfiguration{app}
   }
   // .details.appDetails.versionString
   det.VersionString, err = docV2.Get(13).Get(1).GetString(4)
   if err != nil {
      return nil, err
   }
   // .details.appDetails.installationSize
   det.Size, err = docV2.Get(13).Get(1).GetVarint(9)
   if err != nil {
      return nil, err
   }
   // .details.appDetails.uploadDate
   det.UploadDate, err = docV2.Get(13).Get(1).GetString(16)
   if err != nil {
      return nil, err
   }
   // .details.appDetails.file
   for _, file := range docV2.Get(13).Get(1).GetMessages(17) {
      // .fileType
      typ, err := file.GetVarint(1)
      if err != nil {
         return nil, err
      }
      det.File = append(det.File, typ)
   }
   // The following fields should work with any ABI.
   // .title
   det.Title, err = docV2.GetString(5)
   if err != nil {
      return nil, err
   }
   // .creator
   det.Creator, err = docV2.GetString(6)
   if err != nil {
      return nil, err
   }
   // .offer.micros
   det.Micros, err = docV2.Get(8).GetVarint(1)
   if err != nil {
      return nil, err
   }
   // .offer.currencyCode
   det.CurrencyCode, err = docV2.Get(8).GetString(2)
   if err != nil {
      return nil, err
   }
   // I dont know the name of field 70
   // .details.appDetails
   det.NumDownloads, err = docV2.Get(13).Get(1).GetVarint(70)
   if err != nil {
      return nil, err
   }
   return &det, nil
}

func (t Token) headerVersion(androidID, agent Fixed64) (*Header, error) {
   val := url.Values{
      "Token": {t.Token},
      "service": {"oauth2:https://www.googleapis.com/auth/googleplay"},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/auth", strings.NewReader(val),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errors.New(res.Status)
   }
   var head Header
   head.AndroidID = uint64(androidID)
   head.Auth = parseQuery(res.Body).Get("Auth")
   head.SDK = 9
   head.VersionCode = uint64(agent)
   return &head, nil
}

func (t Token) Header(androidID Fixed64) (*Header, error) {
   return t.headerVersion(androidID, 9999_9999)
}

func (t Token) SingleAPK(androidID Fixed64) (*Header, error) {
   return t.headerVersion(androidID, 8091_9999)
}

func (h Header) SetAgent(head http.Header) {
   var buf []byte
   buf = append(buf, "Android-Finsky (sdk="...)
   buf = strconv.AppendInt(buf, h.SDK, 10)
   buf = append(buf, ",versionCode="...)
   buf = strconv.AppendUint(buf, h.VersionCode, 10)
   head.Set("User-Agent", string(buf))
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
   h.SetAuth(req.Header)
   h.SetDevice(req.Header)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return errors.New(res.Status)
   }
   return nil
}
