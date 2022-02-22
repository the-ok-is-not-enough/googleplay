package googleplay

import (
   "github.com/89z/format/protobuf"
   "github.com/89z/format"
   "net/http"
   "net/url"
   "strings"
)

var LogLevel format.LogLevel

type Device struct {
   AndroidID uint64
}

type Document struct {
   ID string
   Title string
   Creator string
   NextPageURL string
   Child []Document
}

type Token struct {
   Token string
}

type Header struct {
   http.Header
}

func (h Header) Category(cat string) (*Document, error) {
   return h.getCategory(cat, "")
}

func (h Header) CategoryNext(nextPage string) (*Document, error) {
   return h.getCategory("", nextPage)
}

func (h Header) getCategory(cat, nextPage string) (*Document, error) {
   var buf strings.Builder
   buf.WriteString("https://android.clients.google.com/fdfe/")
   if cat != "" {
      buf.WriteString("list?")
      val := url.Values{
         "c": {"3"},
         "cat": {cat},
         "ctr": {"apps_topselling_free"},
      }.Encode()
      buf.WriteString(val)
   } else {
      buf.WriteString(nextPage)
   }
   req, err := http.NewRequest("GET", buf.String(), nil)
   if err != nil {
      return nil, err
   }
   req.Header = h.Header
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
   docV2 := responseWrapper.Get(1, "payload").
      Get(1, "listResponse").
      Get(2, "doc")
   var doc Document
   for _, elem := range docV2.GetMessages(11, "child") {
      var child Document
      child.ID = elem.GetString(1, "docID")
      child.Title = elem.GetString(5, "title")
      child.Creator = elem.GetString(6, "creator")
      doc.Child = append(doc.Child, child)
   }
   doc.NextPageURL = docV2.Get(12, "containerMetadata").
      GetString(2, "nextPageUrl")
   return &doc, nil
}
