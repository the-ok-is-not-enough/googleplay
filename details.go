package googleplay

import (
   "errors"
   "github.com/89z/rosso/protobuf"
   "github.com/89z/rosso/strconv"
   "io"
   "net/http"
   "net/url"
   "time"
)

func (self Details) Upload_Date() (string, error) {
   // .details.appDetails.uploadDate
   date, err := self.Get(13).Get(1).Get_String(16)
   if err != nil {
      return "", errors.New("uploadDate not found, try another platform")
   }
   return date, nil
}

func (self Header) Details(app string) (*Details, error) {
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/details", nil,
   )
   if err != nil {
      return nil, err
   }
   // half of the apps I test require User-Agent,
   // so just set it for all of them
   self.Set_Agent(req.Header)
   self.Set_Auth(req.Header)
   self.Set_Device(req.Header)
   req.URL.RawQuery = "doc=" + url.QueryEscape(app)
   res, err := Client.Do(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   body, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   // ResponseWrapper
   response_wrapper, err := protobuf.Unmarshal(body)
   if err != nil {
      return nil, err
   }
   var det Details
   // .payload.detailsResponse.docV2
   det.Message = response_wrapper.Get(1).Get(2).Get(4)
   return &det, nil
}

type Details struct {
   protobuf.Message
}

// will fail with wrong ABI
func (self Details) Version_Code() (uint64, error) {
   // .details.appDetails.versionCode
   return self.Get(13).Get(1).Get_Varint(3)
}

// will fail with wrong ABI
func (self Details) Version() (string, error) {
   // .details.appDetails.versionString
   return self.Get(13).Get(1).Get_String(4)
}

// will fail with wrong ABI
func (self Details) Installation_Size() (uint64, error) {
   // .details.appDetails.installationSize
   return self.Get(13).Get(1).Get_Varint(9)
}

// should work with any ABI
func (self Details) Title() (string, error) {
   // .title
   return self.Get_String(5)
}

// should work with any ABI
func (self Details) Creator() (string, error) {
   // .creator
   return self.Get_String(6)
}

// should work with any ABI
func (self Details) Micros() (uint64, error) {
   // .offer.micros
   return self.Get(8).Get_Varint(1)
}

// should work with any ABI
func (self Details) Currency_Code() (string, error) {
   // .offer.currencyCode
   return self.Get(8).Get_String(2)
}

// should work with any ABI
func (self Details) Num_Downloads() (uint64, error) {
   // .details.appDetails
   // I dont know the name of field 70, but the similar field 13 is called
   // .numDownloads
   return self.Get(13).Get(1).Get_Varint(70)
}

// will fail with wrong ABI
func (self Details) File() []File_Metadata {
   var files []File_Metadata
   // .details.appDetails.file
   for _, file := range self.Get(13).Get(1).Get_Messages(17) {
      files = append(files, File_Metadata{file})
   }
   return files
}

// This only works with English. You can force English with:
// Accept-Language: en
func (self Details) Time() (time.Time, error) {
   date, err := self.Upload_Date()
   if err != nil {
      return time.Time{}, err
   }
   return time.Parse("Jan 2, 2006", date)
}

func (self Details) MarshalText() ([]byte, error) {
   var b []byte
   b = append(b, "Title: "...)
   if v, err := self.Title(); err != nil {
      return nil, err
   } else {
      b = append(b, v...)
   }
   b = append(b, "\nCreator: "...)
   if v, err := self.Creator(); err != nil {
      return nil, err
   } else {
      b = append(b, v...)
   }
   b = append(b, "\nUpload Date: "...)
   if v, err := self.Upload_Date(); err != nil {
      return nil, err
   } else {
      b = append(b, v...)
   }
   b = append(b, "\nVersion: "...)
   if v, err := self.Version(); err != nil {
      return nil, err
   } else {
      b = append(b, v...)
   }
   b = append(b, "\nVersion Code: "...)
   if v, err := self.Version_Code(); err != nil {
      return nil, err
   } else {
      b = strconv.AppendUint(b, v, 10)
   }
   b = append(b, "\nNum Downloads: "...)
   if v, err := self.Num_Downloads(); err != nil {
      return nil, err
   } else {
      b = append(b, strconv.Number(v)...)
   }
   b = append(b, "\nInstallation Size: "...)
   if v, err := self.Installation_Size(); err != nil {
      return nil, err
   } else {
      b = append(b, strconv.Size(v)...)
   }
   b = append(b, "\nFile:"...)
   for _, file := range self.File() {
      if v, err := file.File_Type(); err != nil {
         return nil, err
      } else if v >= 1 {
         b = append(b, " OBB"...)
      } else {
         b = append(b, " APK"...)
      }
   }
   b = append(b, "\nOffer: "...)
   if v, err := self.Micros(); err != nil {
      return nil, err
   } else {
      b = strconv.AppendUint(b, v, 10)
   }
   b = append(b, ' ')
   if v, err := self.Currency_Code(); err != nil {
      return nil, err
   } else {
      b = append(b, v...)
   }
   b = append(b, '\n')
   return b, nil
}
