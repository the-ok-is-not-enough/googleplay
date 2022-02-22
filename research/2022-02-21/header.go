package googleplay

import (
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

type Auth struct {
   Auth string
}

type Delivery struct {
   DownloadURL string
   SplitDeliveryData []SplitDeliveryData
}

var LogLevel format.LogLevel

type Details struct {
   Title string
   UploadDate string
   VersionString string
   VersionCode uint64
   NumDownloads uint64
   Size uint64
   Micros uint64
   CurrencyCode string
}

type Device struct {
   AndroidID uint64
}

type SplitDeliveryData struct {
   ID string
   DownloadURL string
}

type errorString string

func (e errorString) Error() string {
   return string(e)
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
