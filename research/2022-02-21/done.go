package googleplay

import (
   "github.com/89z/format"
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
