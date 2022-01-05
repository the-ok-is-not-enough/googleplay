package gplayapi

import (
   "bytes"
   "fmt"
   "github.com/89z/format/protobuf"
   "github.com/Juby210/gplayapi-go/gpproto"
   "google.golang.org/protobuf/proto"
   "net/http"
)

var zombo = protobuf.Message{
   protobuf.Tag{Number:4, String:""}:protobuf.Message{
      protobuf.Tag{Number:1, String:""}:protobuf.Message{
         protobuf.Tag{Number:10, String:""}:uint64(29),
      },
   },
   protobuf.Tag{Number:14, String:""}:uint64(3),
   protobuf.Tag{Number:18, String:""}:protobuf.Message{
      protobuf.Tag{Number:1, String:"touchScreen"}:uint64(3),
      protobuf.Tag{Number:2, String:"keyboard"}:uint64(0),
      protobuf.Tag{Number:3, String:"navigation"}:uint64(0),
      protobuf.Tag{Number:4, String:"screenLayout"}:uint64(0),
      protobuf.Tag{Number:5, String:"hasHardKeyboard"}:uint64(0),
      protobuf.Tag{Number:6, String:"hasFiveWayNavigation"}:uint64(0),
      protobuf.Tag{Number:7, String:"screenDensity"}:uint64(0),
      protobuf.Tag{Number:8, String:"glEsVersion"}:uint64(0x3_0001),
      protobuf.Tag{Number:11, String:"nativePlatform"}:[]string{
         "x86",
         "armeabi-v7a",
      },
      protobuf.Tag{Number:12, String:"screenWidth"}:uint64(1),
      protobuf.Tag{Number:26, String:""}:[]protobuf.Message{
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.bluetooth",
         protobuf.Tag{Number:2, String:""}:uint64(0)},
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.bluetooth_le",
         protobuf.Tag{Number:2, String:""}:uint64(0)},
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.camera",
         protobuf.Tag{Number:2, String:""}:uint64(0)},
         protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
         protobuf.Tag{Number:1, String:""}:"android.hardware.camera.autofocus"},
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.location",
         protobuf.Tag{Number:2, String:""}:uint64(0)},
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.location.gps",
         protobuf.Tag{Number:2, String:""}:uint64(0)},
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.microphone",
         protobuf.Tag{Number:2, String:""}:uint64(0)},
         protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
         protobuf.Tag{Number:1, String:""}:"android.hardware.screen.landscape"},
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.screen.portrait",
         protobuf.Tag{Number:2, String:""}:uint64(0)},
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.sensor.accelerometer",
         protobuf.Tag{Number:2, String:""}:uint64(0)},
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.telephony",
         protobuf.Tag{Number:2, String:""}:uint64(0)},
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.touchscreen",
         protobuf.Tag{Number:2, String:""}:uint64(0)},
         protobuf.Message{protobuf.Tag{Number:2, String:""}:uint64(0),
         protobuf.Tag{Number:1, String:""}:"android.hardware.usb.host"},
         protobuf.Message{protobuf.Tag{Number:1, String:""}:"android.hardware.wifi",
         protobuf.Tag{Number:2, String:""}:uint64(0)},
      },
   },
}

func (client *GooglePlayClient) checkIn(req *gpproto.AndroidCheckinRequest) (resp *gpproto.AndroidCheckinResponse, err error) {
   b := zombo.Marshal()
   r, _ := http.NewRequest("POST", UrlCheckIn, bytes.NewReader(b))
   client.setAuthHeaders(r)
   r.Header.Set("Content-Type", "application/x-protobuffer")
   r.Header.Set("Host", "android.clients.google.com")
   b, _, err = doReq(r)
   if err != nil {
   return
   }
   resp = &gpproto.AndroidCheckinResponse{}
   err = proto.Unmarshal(b, resp)
   return
}

type AuthData struct {
	Email                         string
	AASToken                      string
	AuthToken                     string
	GsfID                         string
	DeviceCheckInConsistencyToken string
	DeviceConfigToken             string
	DFECookie                     string
	Locale                        string
}

func (client *GooglePlayClient) GenerateGsfID() (gsfID string, err error) {
	req := client.DeviceInfo.GenerateAndroidCheckInRequest()
	checkInResp, err := client.checkIn(req)
	if err != nil {
		return
	}
	gsfID = fmt.Sprintf("%x", checkInResp.GetAndroidId())
	client.AuthData.GsfID = gsfID
	client.AuthData.DeviceCheckInConsistencyToken = checkInResp.GetDeviceCheckinConsistencyToken()
	return
}

func (client *GooglePlayClient) setAuthHeaders(r *http.Request) {
   r.Header.Set("app", "com.google.android.gms")
   r.Header.Set("User-Agent", client.DeviceInfo.GetAuthUserAgent())
   if client.AuthData.GsfID != "" {
      r.Header.Set("device", client.AuthData.GsfID)
   }
}

func (client *GooglePlayClient) setDefaultHeaders(r *http.Request) {
   data := client.AuthData
   r.Header.Set("Authorization", "Bearer "+data.AuthToken)
   r.Header.Set("X-DFE-Device-Id", data.GsfID)
   agent := "Android-Finsky (sdk=99,versionCode=99999999)"
   r.Header.Set("User-Agent", agent)
}
