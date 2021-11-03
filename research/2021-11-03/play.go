package googleplay

import (
   "encoding/json"
   "net/http"
   "net/http/httputil"
   "strconv"
   "strings"
)

const origin = "https://android.clients.google.com"

const config = `
{"deviceConfiguration":{"touchScreen":3, "keyboard":1, "navigation":1,
"screenLayout":2, "hasHardKeyboard":false, "hasFiveWayNavigation":false,
"screenDensity":420, "glEsVersion":196610,
"systemAvailableFeature":["android.hardware.touchscreen",
"android.hardware.wifi"], "nativePlatform":["arm64-v8a,armeabi-v7a",
"armeabi"]}}
`

func upload(auth, deviceID string) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/uploadDeviceConfig", strings.NewReader(config),
   )
   if err != nil {
      return nil, err
   }
   req.Header = http.Header{
      "Accept": {"application/json"},
      "Authorization": {"Bearer " + auth},
      "Content-Type": {"application/json"},
      "User-Agent": {"Android-Finsky (versionCode=81031200,sdk=27)"},
      "X-DFE-Device-ID": {deviceID},
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return httputil.DumpResponse(res, true)
}

type checkin struct {
   Android_ID int64
}

func newCheckin() (*checkin, error) {
   req, err := http.NewRequest(
      "POST", origin + "/checkin",
      strings.NewReader(`{"checkin":{}, "version":3}`),
   )
   if err != nil {
      return nil, err
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   check := new(checkin)
   if err := json.NewDecoder(res.Body).Decode(check); err != nil {
      return nil, err
   }
   return check, nil
}

func (c checkin) String() string {
   return strconv.FormatInt(c.Android_ID, 16)
}
