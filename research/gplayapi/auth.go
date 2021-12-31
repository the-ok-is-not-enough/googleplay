package gplayapi

import (
   "bytes"
   "errors"
   "fmt"
   "gplayapi/gpproto"
   "google.golang.org/protobuf/proto"
   "net/http"
   "net/url"
   "strconv"
)

func (client *GooglePlayClient) GenerateGsfID() (string, error) {
   req := client.DeviceInfo.GenerateAndroidCheckInRequest()
   checkInResp, err := client.checkIn(req)
   if err != nil {
      return "", err
   }
   client.AuthData.GsfID = fmt.Sprintf("%x", checkInResp.GetAndroidId())
   return client.AuthData.GsfID, nil
}

func (client *GooglePlayClient) checkIn(req *gpproto.AndroidCheckinRequest) (resp *gpproto.AndroidCheckinResponse, err error) {
   b, err := proto.Marshal(req)
   if err != nil {
      return
   }
   r, _ := http.NewRequest("POST", UrlCheckIn, bytes.NewReader(b))
   r.Header.Set("Content-Type", "application/x-protobuffer")
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

func (client *GooglePlayClient) uploadDeviceConfig() (*gpproto.UploadDeviceConfigResponse, error) {
	b, err := proto.Marshal(&gpproto.UploadDeviceConfigRequest{DeviceConfiguration: client.DeviceInfo.GetDeviceConfigProto()})
	if err != nil {
		return nil, err
	}
	r, _ := http.NewRequest("POST", UrlUploadDeviceConfig, bytes.NewReader(b))
	payload, err := client.doAuthedReq(r)
	if err != nil {
		return nil, err
	}
	return payload.UploadDeviceConfigResponse, nil
}

func (client *GooglePlayClient) GenerateGPToken() (string, error) {
	params := &url.Values{}
	client.setDefaultAuthParams(params)
	client.setAuthParams(params)

	params.Set("app", "com.google.android.gms")
	params.Set("service", "oauth2:https://www.googleapis.com/auth/googleplay")

	r, _ := http.NewRequest("POST", UrlAuth+"?"+params.Encode(), nil)
	b, _, err := doReq(r)
	if err != nil {
		return "", nil
	}
	resp := parseResponse(string(b))
	token, ok := resp["Auth"]
	if !ok {
		return "", errors.New("authentication failed: could not generate oauth token")
	}
	return token, nil
}

func (client *GooglePlayClient) setDefaultAuthParams(params *url.Values) {
   params.Set("sdk_version", strconv.Itoa(int(client.DeviceInfo.Build.GetSdkVersion())))
   params.Set("email", client.AuthData.Email)
   params.Set("google_play_services_version", strconv.Itoa(int(client.DeviceInfo.Build.GetGoogleServices())))
   params.Set("device_country", "us")
   params.Set("lang", "en-gb")
   params.Set("callerSig", "38918a453d07199354f8b19af05ec6562ced5788")
}

func (client *GooglePlayClient) setAuthParams(params *url.Values) {
   params.Set("Token", client.AuthData.AASToken)
   params.Set("check_email", "1")
   params.Set("client_sig", "38918a453d07199354f8b19af05ec6562ced5788")
}

func (client *GooglePlayClient) setDefaultHeaders(r *http.Request) {
   data := client.AuthData
   r.Header.Set("Authorization", "Bearer "+data.AuthToken)
   r.Header.Set("X-DFE-Device-Id", data.GsfID)
}
