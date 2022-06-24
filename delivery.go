package googleplay

import (
   "errors"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

func (h Header) Delivery(app string, ver uint64) (*Delivery, error) {
   req, err := http.NewRequest(
      "GET", "https://play-fe.googleapis.com/fdfe/delivery", nil,
   )
   if err != nil {
      return nil, err
   }
   h.Set_Agent(req.Header)
   h.Set_Auth(req.Header) // needed for single APK
   h.Set_Device(req.Header)
   req.URL.RawQuery = url.Values{
      "doc": {app},
      "vc": {strconv.FormatUint(ver, 10)},
   }.Encode()
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
   // .payload.deliveryResponse.status
   status, err := response_wrapper.Get(1).Get(21).Get_Varint(1)
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
   app_data := response_wrapper.Get(1).Get(21).Get(2)
   var del Delivery
   // .downloadUrl
   del.Download_URL, err = app_data.Get_String(3)
   if err != nil {
      return nil, err
   }
   del.Package_Name = app
   del.Version_Code = ver
   // .splitDeliveryData
   for _, data := range app_data.Get_Messages(15) {
      var split Split_Data
      // .id
      split.ID, err = data.Get_String(1)
      if err != nil {
         return nil, err
      }
      // .downloadUrl
      split.Download_URL, err = data.Get_String(5)
      if err != nil {
         return nil, err
      }
      del.Split_Data = append(del.Split_Data, split)
   }
   // .additionalFile
   for _, file := range app_data.Get_Messages(4) {
      var app File_Metadata
      // .fileType
      app.File_Type, err = file.Get_Varint(1)
      if err != nil {
         return nil, err
      }
      // .downloadUrl
      app.Download_URL, err = file.Get_String(4)
      if err != nil {
         return nil, err
      }
      del.Additional_File = append(del.Additional_File, app)
   }
   return &del, nil
}

type Split_Data struct {
   ID string
   Download_URL string
}

type File_Metadata struct {
   Download_URL string
   File_Type uint64
}

type Delivery struct {
   Additional_File []File_Metadata
   Download_URL string
   Package_Name string
   Split_Data []Split_Data
   Version_Code uint64
}

func (d Delivery) Additional(typ uint64) string {
   var buf []byte
   if typ == 0 {
      buf = append(buf, "main"...)
   } else {
      buf = append(buf, "patch"...)
   }
   buf = append(buf, '.')
   buf = strconv.AppendUint(buf, d.Version_Code, 10)
   buf = append(buf, '.')
   buf = append(buf, d.Package_Name...)
   buf = append(buf, ".obb"...)
   return string(buf)
}

func (d Delivery) Download() string {
   var buf []byte
   buf = append(buf, d.Package_Name...)
   buf = append(buf, '-')
   buf = strconv.AppendUint(buf, d.Version_Code, 10)
   buf = append(buf, ".apk"...)
   return string(buf)
}

func (d Delivery) Split(id string) string {
   var buf []byte
   buf = append(buf, d.Package_Name...)
   buf = append(buf, '-')
   buf = append(buf, id...)
   buf = append(buf, '-')
   buf = strconv.AppendUint(buf, d.Version_Code, 10)
   buf = append(buf, ".apk"...)
   return string(buf)
}
