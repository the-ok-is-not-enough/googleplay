package googleplay

import (
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

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
   // User-Agent is only needed with "/fdfe/details" for some apps, example:
   // com.xiaomi.smarthome
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

func (h Header) Category(cat string) ([]Document, error) {
   // You can also use "/fdfe/browse" or "/fdfe/homeV2", but they do Prefetch,
   // and seem to ignore the X-DFE-No-Prefetch:true header. You can also use
   // "/fdfe/list", but it requires subcategory.
   req, err := http.NewRequest("GET", origin + "/fdfe/getHomeStream", nil)
   if err != nil {
      return nil, err
   }
   req.Header = h.Header
   // You can change this to "4", but it will fail with newer versionCode:
   req.URL.RawQuery = "c=3&cat=" + url.QueryEscape(cat)
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
   child := responseWrapper.Get(1, "payload").
      Get(1, "listResponse").
      Get(2, "doc").
      GetMessages(11, "child")
   var docs []Document
   for _, element := range child {
      switch element.GetString(5, "title") {
      case "Based on your recent activity", "Recommended for you":
      default:
         for _, element := range element.GetMessages(11, "child") {
            var doc Document
            doc.ID = element.GetString(1, "docID")
            doc.Title = element.GetString(5, "title")
            doc.Creator = element.GetString(6, "creator")
            docs = append(docs, doc)
         }
      }
   }
   return docs, nil
}

func (h Header) Details(app string) (*Details, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/details", nil)
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
   docV2 := responseWrapper.Get(1, "payload").
      Get(2, "detailsResponse").
      Get(4, "docV2")
   var det Details
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

func (h Header) Delivery(app string, ver int64) (*Delivery, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/delivery", nil)
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
      return nil, errorString("Regional lockout")
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

// Purchase app. Only needs to be done once per Google account.
func (h Header) Purchase(app string) error {
   query := "doc=" + url.QueryEscape(app)
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/purchase", strings.NewReader(query),
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
   return res.Body.Close()
}
