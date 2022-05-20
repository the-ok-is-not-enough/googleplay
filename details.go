package googleplay

import (
   "errors"
   "fmt"
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
)

type Details struct {
   Title String
   Creator String
   UploadDate String
   VersionString String
   VersionCode Varint
   NumDownloads Varint
   Size Varint
   Files int
   Micros Varint
   CurrencyCode String
}

func (d Details) Format(f fmt.State, verb rune) {
   fmt.Fprintln(f, "Title:", d.Title)
   fmt.Fprintln(f, "Creator:", d.Creator)
   fmt.Fprintln(f, "UploadDate:", d.UploadDate)
   fmt.Fprintln(f, "VersionString:", d.VersionString)
   fmt.Fprintln(f, "VersionCode:", d.VersionCode)
   fmt.Fprintln(f, "NumDownloads:", format.LabelNumber(d.NumDownloads))
   fmt.Fprintln(f, "Size:", format.LabelSize(d.Size))
   fmt.Fprintln(f, "Files:", d.Files)
   fmt.Fprint(f, "Offer: ", d.Micros, " ", d.CurrencyCode)
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
      return nil, errors.New("wrong platform")
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
   files := docV2.Get(13).Get(1).GetMessages(17)
   det.Files = len(files)
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
