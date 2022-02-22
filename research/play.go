package googleplay

import (
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strings"
)

var LogLevel format.LogLevel

func list(category string) ([]Document, error) {
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/list", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {bearer},
      "X-Dfe-Device-Id": {"3588cd1e2b1781ee"},
   }
   req.URL.RawQuery = url.Values{
      "c": {"3"},
      "cat": {category},
      "ctr": {"apps_topselling_free"},
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
   child := responseWrapper.Get(1, "payload").
      Get(1, "listResponse").
      Get(2, "doc"). // DocV2
      GetMessages(11, "child") // DocV2
   var docs []Document
   for _, element := range child {
      var doc Document
      doc.ID = element.GetString(1, "docID")
      doc.Title = element.GetString(5, "title")
      doc.Creator = element.GetString(6, "creator")
      docs = append(docs, doc)
   }
   return docs, nil
}

type docV2 struct {
   docID string
   title string
   creator string
   child []docV2
   containerMetadata struct {
   }
}

type responseWrapper struct {
   payload struct {
      listResponse struct {
         doc docV2
      }
   }
}
