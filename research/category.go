package googleplay

import (
   "bufio"
   "github.com/89z/format"
   "github.com/89z/format/protobuf"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

var LogLevel format.LogLevel

func parseQuery(query io.Reader) url.Values {
   val := make(url.Values)
   buf := bufio.NewScanner(query)
   for buf.Scan() {
      var key string
      for i, elem := range strings.SplitN(buf.Text(), "=", 2) {
         switch i {
         case 0:
            key = elem
         case 1:
            val.Add(key, elem)
         }
      }
   }
   return val
}

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

func Category(cat string) (*Document, error) {
   return getCategory(cat, "")
}

func CategoryNext(nextPage string) (*Document, error) {
   return getCategory("", nextPage)
}

func getCategory(cat, nextPage string) (*Document, error) {
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
   req.Header = http.Header{
      "Authorization": {bearer},
      "X-Dfe-Device-ID": {"3588cd1e2b1781ee"},
   }
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

func (d Document) String() string {
   var buf strings.Builder
   if d.Child != nil {
      for i, doc := range d.Child {
         if i >= 1 {
            buf.WriteString("\n\n")
         }
         buf.WriteString(doc.String())
      }
      buf.WriteString("\nNextPageURL: ")
      buf.WriteString(d.NextPageURL)
   } else {
      buf.WriteString("ID: ")
      buf.WriteString(d.ID)
      buf.WriteString("\nTitle: ")
      buf.WriteString(d.Title)
      buf.WriteString("\nCreator: ")
      buf.WriteString(d.Creator)
   }
   return buf.String()
}

type Header struct {
   http.Header
}

type Token struct {
   Token string
}

func (t Token) Header(dev *Device) (*Header, error) {
   return t.headerVersion(dev, 9999_9999)
}

func (t Token) SingleAPK(dev *Device) (*Header, error) {
   return t.headerVersion(dev, 8091_9999)
}

func (t Token) headerVersion(dev *Device, version int64) (*Header, error) {
   val := url.Values{
      "Token": {t.Token},
      "service": {"oauth2:https://www.googleapis.com/auth/googleplay"},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://android.googleapis.com/auth", strings.NewReader(val),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   LogLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   var head Header
   head.Header = make(http.Header)
   auth := parseQuery(res.Body).Get("Auth")
   if auth != "" {
      head.Set("Authorization", "Bearer " + auth)
   }
   buf := []byte("Android-Finsky (sdk=9,versionCode=")
   buf = strconv.AppendInt(buf, version, 10)
   head.Set("User-Agent", string(buf))
   head.Set("X-DFE-Device-ID", strconv.FormatUint(dev.AndroidID, 16))
   return &head, nil
}
