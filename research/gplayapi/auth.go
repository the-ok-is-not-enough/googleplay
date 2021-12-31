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

func (client *GooglePlayClient) checkIn(req *gpproto.AndroidCheckinRequest) (resp *gpproto.AndroidCheckinResponse, err error) {
	b, err := proto.Marshal(req)
	if err != nil {
		return
	}
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
	client.setAuthHeaders(r)
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

func (client *GooglePlayClient) acceptTos(tosToken string) error {
	r, _ := http.NewRequest("POST", UrlTosAccept+"?toscme=false&tost="+tosToken, nil)
	_, err := client.doAuthedReq(r)
	return err
}

func (client *GooglePlayClient) setAuthHeaders(r *http.Request) {
	r.Header.Set("app", "com.google.android.gms")
	r.Header.Set("User-Agent", client.DeviceInfo.GetAuthUserAgent())
	if client.AuthData.GsfID != "" {
		r.Header.Set("device", client.AuthData.GsfID)
	}
}

func (client *GooglePlayClient) setDefaultAuthParams(params *url.Values) {
   if client.AuthData.GsfID != "" {
      params.Set("androidId", client.AuthData.GsfID)
   }
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
   r.Header.Set("X-DFE-Content-Filters", "")
   r.Header.Set("Authorization", "Bearer "+data.AuthToken)
   r.Header.Set("X-DFE-Device-Id", data.GsfID)
}




