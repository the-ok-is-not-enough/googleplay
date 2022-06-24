package googleplay

import (
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

func (h Header) Details(app string) (*Details, error) {
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/details", nil,
   )
   if err != nil {
      return nil, err
   }
   // half of the apps I test require User-Agent,
   // so just set it for all of them
   h.Set_Agent(req.Header)
   h.Set_Auth(req.Header)
   h.Set_Device(req.Header)
   req.URL.RawQuery = "doc=" + url.QueryEscape(app)
   res, err := Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   buf, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   response_wrapper, err := protobuf.Unmarshal(buf)
   if err != nil {
      return nil, err
   }
   // .payload.detailsResponse.docV2
   doc_V2 := response_wrapper.Get(1).Get(2).Get(4)
   var det Details
   // The following fields will fail with wrong ABI, so try them first. If the
   // first one passes, then use native error for the rest.
   // .details.appDetails.versionCode
   det.Version_Code, err = doc_V2.Get(13).Get(1).Get_Varint(3)
   if err != nil {
      return nil, version_error{app}
   }
   // .details.appDetails.versionString
   det.Version, err = doc_V2.Get(13).Get(1).Get_String(4)
   if err != nil {
      return nil, err
   }
   // .details.appDetails.installationSize
   det.Size, err = doc_V2.Get(13).Get(1).Get_Varint(9)
   if err != nil {
      return nil, err
   }
   // .details.appDetails.uploadDate
   det.Upload_Date, err = doc_V2.Get(13).Get(1).Get_String(16)
   if err != nil {
      return nil, err
   }
   // .details.appDetails.file
   for _, file := range doc_V2.Get(13).Get(1).Get_Messages(17) {
      // .fileType
      typ, err := file.Get_Varint(1)
      if err != nil {
         return nil, err
      }
      det.File = append(det.File, typ)
   }
   // The following fields should work with any ABI.
   // .title
   det.Title, err = doc_V2.Get_String(5)
   if err != nil {
      return nil, err
   }
   // .creator
   det.Creator, err = doc_V2.Get_String(6)
   if err != nil {
      return nil, err
   }
   // .offer.micros
   det.Micros, err = doc_V2.Get(8).Get_Varint(1)
   if err != nil {
      return nil, err
   }
   // .offer.currencyCode
   det.Currency_Code, err = doc_V2.Get(8).Get_String(2)
   if err != nil {
      return nil, err
   }
   // I dont know the name of field 70
   // .details.appDetails
   det.Downloads, err = doc_V2.Get(13).Get(1).Get_Varint(70)
   if err != nil {
      return nil, err
   }
   return &det, nil
}

type version_error struct {
   app string
}

func (v version_error) Error() string {
   var buf strings.Builder
   buf.WriteString(v.app)
   buf.WriteString(" versionCode missing\n")
   buf.WriteString("Check nativePlatform")
   return buf.String()
}

type Details struct {
   Creator string
   Currency_Code string
   Downloads uint64
   File []uint64
   Micros uint64
   Size uint64
   Title string
   Upload_Date string // Jun 1, 2021
   Version string
   Version_Code uint64
}

func (d Details) String() string {
   var buf []byte
   buf = append(buf, "Title: "...)
   buf = append(buf, d.Title...)
   buf = append(buf, "\nCreator: "...)
   buf = append(buf, d.Creator...)
   buf = append(buf, "\nDate: "...)
   buf = append(buf, d.Upload_Date...)
   buf = append(buf, "\nVersion: "...)
   buf = append(buf, d.Version...)
   buf = append(buf, "\nVersion code: "...)
   buf = strconv.AppendUint(buf, d.Version_Code, 10)
   buf = append(buf, "\nDownloads: "...)
   buf = append(buf, format.Label_Number(d.Downloads)...)
   buf = append(buf, "\nSize: "...)
   buf = append(buf, format.Label_Size(d.Size)...)
   buf = append(buf, "\nFile:"...)
   for _, file := range d.File {
      if file == 0 {
         buf = append(buf, " APK"...)
      } else {
         buf = append(buf, " OBB"...)
      }
   }
   buf = append(buf, "\nOffer: "...)
   buf = strconv.AppendUint(buf, d.Micros, 10)
   buf = append(buf, ' ')
   buf = append(buf, d.Currency_Code...)
   return string(buf)
}

// This only works with English. You can force English with:
// Accept-Language: en
func (d Details) Time() (time.Time, error) {
   return time.Parse("Jan 2, 2006", d.Upload_Date)
}
