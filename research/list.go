package googleplay

import (
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
   "net/url"
   "strings"
)

var LogLevel format.LogLevel

func (d document) String() string {
   var buf strings.Builder
   if d.child != nil {
      for i, doc := range d.child {
         if i >= 1 {
            buf.WriteString("\n\n")
         }
         buf.WriteString(doc.String())
      }
      buf.WriteString("\nnextPageURL: ")
      buf.WriteString(d.nextPageURL)
   } else {
      buf.WriteString("id: ")
      buf.WriteString(d.id)
      buf.WriteString("\ntitle: ")
      buf.WriteString(d.title)
      buf.WriteString("\ncreator: ")
      buf.WriteString(d.creator)
   }
   return buf.String()
}

type document struct {
   child []document
   id string
   title string
   creator string
   nextPageURL string
}

func list(category string) (*document, error) {
   return document{}.getList(category)
}

func (d document) list() (*document, error) {
   return d.getList("")
}

func (d document) getList(category string) (*document, error) {
   req, err := http.NewRequest(
      "GET", "https://android.clients.google.com/fdfe/", nil,
   )
   if err != nil {
      return nil, err
   }
   if d.nextPageURL != "" {
      req.URL.Path += d.nextPageURL
   } else {
      req.URL.Path += "list"
      req.URL.RawQuery = url.Values{
         "c": {"3"},
         "cat": {category},
         "ctr": {"apps_topselling_free"},
      }.Encode()
   }
   req.Header = http.Header{
      "Authorization": {bearer},
      "X-Dfe-Device-Id": {"3588cd1e2b1781ee"},
   }
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, errorString(res.Status)
   }
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, err
   }
   docV2 := responseWrapper.Get(1, "payload").
      Get(1, "listResponse").
      Get(2, "doc")
   var doc document
   for _, elem := range docV2.GetMessages(11, "child") {
      var child document
      child.id = elem.GetString(1, "docID")
      child.title = elem.GetString(5, "title")
      child.creator = elem.GetString(6, "creator")
      doc.child = append(doc.child, child)
   }
   doc.nextPageURL = docV2.Get(12, "containerMetadata").
      GetString(2, "nextPageUrl")
   return &doc, nil
}

type errorString string

func (e errorString) Error() string {
   return string(e)
}
