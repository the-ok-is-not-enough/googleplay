package googleplay

import (
   "bytes"
   "crypto/rsa"
   "crypto/sha1"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "github.com/89z/parse/crypto"
   "github.com/89z/parse/net"
   "github.com/segmentio/encoding/proto"
   "io"
   "math/big"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

const AndroidJA3 =
   "769,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-49172-" +
   "51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0"

const androidKey =
   "AAAAgMom/1a/v0lblO2Ubrt60J2gcuXSljGFQXgcyZWveWLEwo6prwgi3iJIZdodyhKZQrNWp" +
   "5nKJ3srRXcUW+F1BD3baEVGcmEgqaLZUNBjm057pKRI16kB0YppeGx5qIQ5QjKzsR8ETQbKLN" +
   "WgRY0QRNVz34kMJR3P/LgHax/6rmf5AAAAAwEAAQ=="

const origin = "https://android.clients.google.com"

func Signature(email, password string) (string, error) {
   data, err := base64.StdEncoding.DecodeString(androidKey)
   if err != nil {
      return "", err
   }
   var key rsa.PublicKey
   buf := crypto.NewBuffer(data)
   // modulus_length | modulus | exponent_length | exponent
   _, mod, ok := buf.ReadUint32LengthPrefixed()
   if ok {
      key.N = new(big.Int).SetBytes(mod)
   }
   _, exp, ok := buf.ReadUint32LengthPrefixed()
   if ok {
      exp := new(big.Int).SetBytes(exp).Int64()
      key.E = int(exp)
   }
   var (
      msg bytes.Buffer
      zero devZero
   )
   msg.WriteString(email)
   msg.WriteByte(0)
   msg.WriteString(password)
   login, err := rsa.EncryptOAEP(
      sha1.New(), zero, &key, msg.Bytes(), nil,
   )
   if err != nil {
      return "", err
   }
   hash := sha1.Sum(data)
   msg.Reset()
   msg.WriteByte(0)
   msg.Write(hash[:4])
   msg.Write(login)
   return base64.URLEncoding.EncodeToString(msg.Bytes()), nil
}

type Checkin struct {
   Checkin struct{} `json:"checkin"`
   Version int `json:"version"`
}

func NewCheckin() Checkin {
   return Checkin{Version: 3}
}

func (c Checkin) Post() (*CheckinResponse, error) {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(c)
   if err != nil {
      return nil, err
   }
   res, err := http.Post(origin + "/checkin", "application/json", buf)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("status %q", res.Status)
   }
   check := new(CheckinResponse)
   if err := json.NewDecoder(res.Body).Decode(check); err != nil {
      return nil, err
   }
   return check, nil
}

type CheckinResponse struct {
   Android_ID int64
}

func (c CheckinResponse) String() string {
   return strconv.FormatInt(c.Android_ID, 16)
}

type Device struct {
   Configuration struct {
      GlEsVersion            int32   `protobuf:"varint,8"`
      HasFiveWayNavigation   bool    `protobuf:"varint,6"`
      HasHardKeyboard        bool    `protobuf:"varint,5"`
      Keyboard               int32   `protobuf:"varint,2"`
      NativePlatform         []string `protobuf:"bytes,11"`
      Navigation             int32   `protobuf:"varint,3"`
      ScreenDensity          int32   `protobuf:"varint,7"`
      ScreenLayout           int32   `protobuf:"varint,4"`
      SystemAvailableFeature []string `protobuf:"bytes,10"`
      TouchScreen            int32   `protobuf:"varint,1"`
   } `protobuf:"bytes,1"`
}

func NewDevice() Device {
   var d Device
   // OpenGL ES version
   // developer.android.com/guide/topics/manifest/uses-feature-element
   d.Configuration.GlEsVersion = 0x0002_0000
   d.Configuration.HasFiveWayNavigation = true
   d.Configuration.HasHardKeyboard = true
   d.Configuration.Keyboard = 1
   d.Configuration.NativePlatform = []string{"armeabi-v7a"}
   d.Configuration.Navigation = 1
   d.Configuration.ScreenDensity = 1
   d.Configuration.ScreenLayout = 1
   d.Configuration.SystemAvailableFeature = []string{
      "android.hardware.touchscreen",
      "android.hardware.wifi",
   }
   d.Configuration.TouchScreen = 1
   return d
}

func (d Device) Upload(check *CheckinResponse, oauth OAuth) error {
   buf, err := proto.Marshal(d)
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", origin + "/fdfe/uploadDeviceConfig", bytes.NewReader(buf),
   )
   if err != nil {
      return err
   }
   req.Header = http.Header{
      "Authorization": {"Bearer " + oauth.Get("Auth")},
      "User-Agent": {"Android-Finsky (versionCode=81031200,sdk=27)"},
      "X-DFE-Device-Id": {check.String()},
   }
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return fmt.Errorf("status %q", res.Status)
   }
   return nil
}

type OAuth struct {
   url.Values
}

// device is Google Service Framework.
func (a OAuth) Details(device, app string) ([]byte, error) {
   req, err := http.NewRequest("GET", origin + "/fdfe/details", nil)
   if err != nil {
      return nil, err
   }
   val := req.URL.Query()
   val.Set("doc", app)
   req.URL.RawQuery = val.Encode()
   req.Header.Set("Authorization", "Bearer " + a.Get("Auth"))
   req.Header.Set("X-DFE-Device-ID", device)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("status %q", res.Status)
   }
   return io.ReadAll(res.Body)
}

type Token struct {
   url.Values
}

// Request refresh token.
func NewToken(email, password string) (*Token, error) {
   hello, err := crypto.ParseJA3(AndroidJA3)
   if err != nil {
      return nil, err
   }
   sig, err := Signature(email, password)
   if err != nil {
      return nil, err
   }
   val := make(url.Values)
   val.Set("Email", email)
   val.Set("EncryptedPasswd", sig)
   val.Set("sdk_version", "16")
   req, err := http.NewRequest(
      "POST", origin + "/auth", strings.NewReader(val.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := crypto.NewTransport(hello.ClientHelloSpec).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   query, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   if val := net.ParseQuery(query); val != nil {
      return &Token{val}, nil
   }
   return nil, fmt.Errorf("parseQuery %q", query)
}

// Exchange refresh token for access token.
func (t Token) OAuth() (*OAuth, error) {
   val := make(url.Values)
   val.Set("Token", t.Get("Token"))
   val.Set("service", "oauth2:https://www.googleapis.com/auth/googleplay")
   req, err := http.NewRequest(
      "POST", origin + "/auth", strings.NewReader(val.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   query, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   if val := net.ParseQuery(query); val != nil {
      return &OAuth{val}, nil
   }
   return nil, fmt.Errorf("parseQuery %q", query)
}

type devZero struct{}

func (devZero) Read(b []byte) (int, error) {
   return len(b), nil
}
