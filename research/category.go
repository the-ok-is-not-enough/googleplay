package googleplay

import (
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "net/http"
   "strings"
   "net/url"
)

var LogLevel format.LogLevel

type Device struct {
   AndroidID uint64
}

type Document struct {
   ID string
   Title string
   Creator string
}

type Header struct {
   http.Header
}

func (h Header) Category(cat string, length int) ([]Document, error) {
   var (
      docs []Document
      done int
      next string
   )
   for done < length {
      var (
         doct []Document
         err error
      )
      if done == 0 {
         doct, next, err = h.documents(cat, "")
      } else {
         doct, next, err = h.documents("", next)
      }
      if err != nil {
         return nil, err
      }
      docs = append(docs, doct...)
      done += len(doct)
   }
   return docs, nil
}

func (h Header) documents(cat, next string) ([]Document, string, error) {
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
      buf.WriteString(next)
   }
   req, err := http.NewRequest("GET", buf.String(), nil)
   if err != nil {
      return nil, "", err
   }
   req.Header = h.Header
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, "", err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, "", errorString(res.Status)
   }
   responseWrapper, err := protobuf.Decode(res.Body)
   if err != nil {
      return nil, "", err
   }
   docV2 := responseWrapper.Get(1, "payload").
      Get(1, "listResponse").
      Get(2, "doc")
   var docs []Document
   for _, child := range docV2.GetMessages(11, "child") {
      var doc Document
      doc.ID = child.GetString(1, "docID")
      doc.Title = child.GetString(5, "title")
      doc.Creator = child.GetString(6, "creator")
      docs = append(docs, doc)
   }
   next = docV2.Get(12, "containerMetadata").GetString(2, "nextPageUrl")
   return docs, next, nil
}

type Token struct {
   Token string
}

type errorString string

func (e errorString) Error() string {
   return string(e)
}
