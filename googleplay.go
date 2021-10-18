package googleplay

import (
   "bytes"
   "github.com/89z/mech"
   "io"
   "net/http"
   "net/url"
   "strings"
)

const Origin = "https://android.clients.google.com"

var Verbose = mech.Verbose

// text/plain encoding algorithm
// html.spec.whatwg.org/multipage/form-control-infrastructure.html
func ParseQuery(query []byte) url.Values {
   res := make(url.Values)
   for _, pair := range bytes.Split(query, []byte{'\n'}) {
      nv := bytes.SplitN(pair, []byte{'='}, 2)
      res.Add(string(nv[0]), string(nv[1]))
   }
   return res
}

type Ac2dm struct {
   url.Values
}

// Exchange embedded token (oauth2_4) for refresh token (aas_et).
// accounts.google.com/EmbeddedSetup
func NewAc2dm(token string) (*Ac2dm, error) {
   req, err := http.NewRequest("POST", Origin + "/auth", nil)
   if err != nil {
      return nil, err
   }
   val := url.Values{
      "ACCESS_TOKEN": {"1"}, "Token": {token}, "service": {"ac2dm"},
   }
   req.URL.RawQuery = val.Encode()
   res, err := mech.RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   query, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return &Ac2dm{
      ParseQuery(query),
   }, nil
}

// Exchange refresh token (aas_et) for access token (Auth).
func (a Ac2dm) OAuth2() (*OAuth2, error) {
   val := url.Values{
      "Token": {
         a.Get("Token"),
      },
      "service": {"oauth2:https://www.googleapis.com/auth/googleplay"},
   }
   req, err := http.NewRequest(
      "POST", Origin + "/auth", strings.NewReader(val.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Content-Type": {"application/x-www-form-urlencoded"},
   }
   res, err := mech.RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   query, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   return &OAuth2{
      ParseQuery(query),
   }, nil
}

type OAuth2 struct {
   url.Values
}

// device is Google Service Framework.
func (o OAuth2) Details(device, app string) ([]byte, error) {
   req, err := http.NewRequest("GET", Origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Authorization": {
         "Bearer " + o.Get("Auth"),
      },
      "X-DFE-Device-Id": {device},
   }
   val := url.Values{
      "doc": {app},
   }
   req.URL.RawQuery = val.Encode()
   res, err := mech.RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return io.ReadAll(res.Body)
}
