package gplayapi

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "google.golang.org/protobuf/proto"
   "gplayapi/gpproto"
   "io"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strconv"
   "strings"
)

func (client *GooglePlayClient) GenerateGsfID() (string, error) {
   req := &gpproto.AndroidCheckinRequest{
      Version:             ptrInt32(3),
      Checkin: &gpproto.AndroidCheckinProto{
         Build:           client.DeviceInfo.Build.AndroidBuildProto,
      },
      DeviceConfiguration: client.DeviceInfo.GetDeviceConfigProto(),
   }
   checkInResp, err := client.checkIn(req)
   if err != nil {
      return "", err
   }
   client.AuthData.GsfID = fmt.Sprintf("%x", checkInResp.GetAndroidId())
   return client.AuthData.GsfID, nil
}

func doReq(r *http.Request) ([]byte, int, error) {
   buf, err := httputil.DumpRequest(r, true)
   if err != nil {
      return nil, 0, err
   }
   fmt.Printf("%q\n\n", buf)
   res, err := httpClient.Do(r)
   if err != nil {
      return nil, 0, err
   }
   defer res.Body.Close()
   b, err := io.ReadAll(res.Body)
   return b, res.StatusCode, err
}

func ptrBool(b bool) *bool {
	return &b
}

func ptrStr(str string) *string {
	return &str
}

func ptrInt32(i int32) *int32 {
	return &i
}

func parseResponse(res string) map[string]string {
	ret := map[string]string{}
	for _, ln := range strings.Split(res, "\n") {
		keyVal := strings.SplitN(ln, "=", 2)
		if len(keyVal) >= 2 {
			ret[keyVal[0]] = keyVal[1]
		}
	}
	return ret
}

func (client *GooglePlayClient) _doAuthedReq(r *http.Request) (_ *gpproto.Payload, err error) {
	client.setDefaultHeaders(r)
	b, status, err := doReq(r)
	if err != nil {
		return
	}
	if status == 401 {
		return nil, GPTokenExpired
	}
	resp := &gpproto.ResponseWrapper{}
	err = proto.Unmarshal(b, resp)
	if err != nil {
		return
	}
	return resp.Payload, nil
}

func (client *GooglePlayClient) doAuthedReq(r *http.Request) (res *gpproto.Payload, err error) {
	res, err = client._doAuthedReq(r)
	if err == GPTokenExpired {
		err = client.RegenerateGPToken()
		if err != nil {
			return
		}
		if client.SessionFile != "" {
			client.SaveSession(client.SessionFile)
		}
		res, err = client._doAuthedReq(r)
	}
	return
}

func (client *GooglePlayClient) RegenerateGPToken() (err error) {
	client.AuthData.AuthToken, err = client.GenerateGPToken()
	return
}

const (
	ImageTypeAppScreenshot = iota + 1
	ImageTypePlayStorePageBackground
	ImageTypeYoutubeVideoLink
	ImageTypeAppIcon
	ImageTypeCategoryIcon
	ImageTypeYoutubeVideoThumbnail = 13

	UrlBase               = "https://android.clients.google.com"
	UrlFdfe               = UrlBase + "/fdfe"
	UrlAuth               = UrlBase + "/auth"
	UrlDetails            = UrlFdfe + "/details"
	UrlDelivery           = UrlFdfe + "/delivery"
	UrlPurchase           = UrlFdfe + "/purchase"
	UrlToc                = UrlFdfe + "/toc"
	UrlTosAccept          = UrlFdfe + "/acceptTos"
	UrlUploadDeviceConfig = UrlFdfe + "/uploadDeviceConfig"
)

type GooglePlayClient struct {
	AuthData   *AuthData
	DeviceInfo *DeviceInfo

	// SessionFile if SessionFile is set then session will be saved to it after modification
	SessionFile string
}

var (
   GPTokenExpired = errors.New("unauthorized, gp token expired")
   httpClient = &http.Client{}
)

func NewClientWithDeviceInfo(email, aasToken string, deviceInfo *DeviceInfo) (client *GooglePlayClient, err error) {
   authData := &AuthData{
   Email:    email,
   AASToken: aasToken,
   Locale:   "en_GB",
   }
   client = &GooglePlayClient{AuthData: authData, DeviceInfo: deviceInfo}
   _, err = client.GenerateGsfID()
   if err != nil {
   return
   }
   deviceConfigRes, err := client.uploadDeviceConfig()
   if err != nil {
   return
   }
   authData.DeviceConfigToken = deviceConfigRes.GetUploadDeviceConfigToken()
   token, err := client.GenerateGPToken()
   if err != nil {
   return
   }
   authData.AuthToken = token
   return
}

func (client *GooglePlayClient) SaveSession(file string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(client.AuthData)
}

type AuthData struct {
   AASToken                      string
   AuthToken                     string
   DFECookie                     string
   DeviceConfigToken             string
   Email                         string
   GsfID                         string
   Locale                        string
}

func (client *GooglePlayClient) uploadDeviceConfig() (*gpproto.UploadDeviceConfigResponse, error) {
   b, err := proto.Marshal(&gpproto.UploadDeviceConfigRequest{DeviceConfiguration: client.DeviceInfo.GetDeviceConfigProto()})
   if err != nil {
      return nil, err
   }
   r, err := http.NewRequest("POST", UrlUploadDeviceConfig, bytes.NewReader(b))
   if err != nil {
      return nil, err
   }
   payload, err := client.doAuthedReq(r)
   if err != nil {
      return nil, err
   }
   return payload.UploadDeviceConfigResponse, nil
}

func (client *GooglePlayClient) GenerateGPToken() (string, error) {
   params := &url.Values{}
   params.Set("sdk_version", strconv.Itoa(int(client.DeviceInfo.Build.GetSdkVersion())))
   params.Set("email", client.AuthData.Email)
   params.Set("google_play_services_version", strconv.Itoa(int(client.DeviceInfo.Build.GetGoogleServices())))
   params.Set("callerSig", "38918a453d07199354f8b19af05ec6562ced5788")
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
