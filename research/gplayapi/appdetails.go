package gplayapi

import (
   "bytes"
   "errors"
   "fmt"
   "google.golang.org/protobuf/proto"
   "gplayapi/gpproto"
   "io"
   "net/http"
   "net/http/httputil"
   "strings"
)

func NewClientWithDeviceInfo(email, aasToken string, deviceInfo *DeviceInfo) (*GooglePlayClient, error) {
   authData := &AuthData{
      AASToken: aasToken,
      Email:    email,
      Locale:   "en_GB",
   }
   client := GooglePlayClient{AuthData: authData, DeviceInfo: deviceInfo}
   checkInResp, err := client.checkIn()
   if err != nil {
      return nil, err
   }
   client.AuthData.GsfID = fmt.Sprintf("%x", checkInResp.GetAndroidId())
   authData.AuthToken = token
   return &client, nil
}

func (client *GooglePlayClient) checkIn() (*gpproto.AndroidCheckinResponse, error) {
   req := &gpproto.AndroidCheckinRequest{
      Version: ptrInt32(3),
      DeviceConfiguration: client.DeviceInfo.GetDeviceConfigProto(),
      Checkin: &gpproto.AndroidCheckinProto{
         Build: &gpproto.AndroidBuildProto{
            SdkVersion: ptrInt32(29),
         },
      },
   }
   b, err := proto.Marshal(req)
   if err != nil {
      return nil, err
   }
   r, err := http.NewRequest("POST", UrlBase + "/checkin", bytes.NewReader(b))
   if err != nil {
      return nil, err
   }
   r.Header.Set("Content-Type", "application/x-protobuffer")
   b, _, err = doReq(r)
   if err != nil {
      return nil, err
   }
   check := new(gpproto.AndroidCheckinResponse)
   if err := proto.Unmarshal(b, check); err != nil {
      return nil, err
   }
   return check, nil
}

func (client *GooglePlayClient) GetAppDetails(packageName string) (*App, error) {
   r, _ := http.NewRequest("GET", UrlDetails+"?doc="+packageName, nil)
   payload, err := client.doAuthedReq(r)
   if err != nil {
      return nil, err
   }
   return buildAppFromItem(payload.DetailsResponse.Item), nil
   
}

func (i *DeviceInfo) GetDeviceConfigProto() *gpproto.DeviceConfigurationProto {
   return &gpproto.DeviceConfigurationProto{
      GlEsVersion:            &i.GLVersion,
      GlExtension:            i.GLExtensions,
      HasFiveWayNavigation:   ptrBool(false),
      HasHardKeyboard:        ptrBool(false),
      Keyboard:               &i.Keyboard,
      Navigation:             &i.Navigation,
      ScreenDensity:          &i.Screen.Density,
      ScreenLayout:           &i.ScreenLayout,
      SystemAvailableFeature: i.Features,
      SystemSharedLibrary:    i.SharedLibraries,
      TouchScreen:            &i.TouchScreen,
      DeviceFeature:          i.GetDeviceFeatures(),
   }
}

type App struct {
   AppInfo            *AppInfo
   CategoryID         int
   CategoryImage      *gpproto.Image
   CategoryName       string
   Changes            string
   ContainsAds        bool
   CoverImage         *gpproto.Image
   Description        string
   DeveloperName      string
   DisplayName        string
   DownloadString     string
   EarlyAccess        bool
   IconImage          *gpproto.Image
   InstantAppLink     string
   IsFree             bool
   IsSystem           bool
   LiveStreamUrl      string
   OfferDetails       map[string]string
   OfferType          int32
   PackageName        string
   Price              string
   PromotionStreamUrl string
   Screenshots        []*gpproto.Image
   ShareUrl           string
   ShortDescription   string
   Size               int64
   TargetSdk          int
   UpdatedOn          string
   VersionCode        int
   VersionName        string
   Video              *gpproto.Image
}

type AppInfo struct {
       AppInfoMap map[string]string
}

func buildAppFromItem(item *gpproto.Item) *App {
   details := item.Details.AppDetails
   app := &App{
      CategoryID:       int(item.GetCategoryId()),
      CategoryName:     details.GetCategoryName(),
      Changes:          details.GetRecentChangesHtml(),
      ContainsAds:      details.InstallNotes != nil,
      Description:      item.GetDescriptionHtml(),
      DeveloperName:    details.GetDeveloperName(),
      DisplayName:      item.GetTitle(),
      DownloadString:   details.GetDownloadLabelAbbreviated(),
      EarlyAccess:      details.EarlyAccessInfo != nil,
      PackageName:      *item.Id,
      ShareUrl:         item.GetShareUrl(),
      ShortDescription: item.GetPromotionalDescription(),
      Size:             details.GetInfoDownloadSize(),
      TargetSdk:        int(details.GetTargetSdkVersion()),
      UpdatedOn:        details.GetInfoUpdatedOn(),
      VersionCode:      int(details.GetVersionCode()),
      VersionName:      details.GetVersionString(),
   }
   return app
}

var Pixel3a = &DeviceInfo{
   Features:        []string{
      "android.hardware.faketouch",
      "android.hardware.screen.portrait",
   },
   GLVersion:       196610,
   Platforms:    []string{"arm64-v8a", "armeabi-v7a", "armeabi"},
   Screen: &DeviceInfoScreen{},
}

type DeviceInfo struct {
   Build           *DeviceBuildInfo
   SimOperator     string
   Platforms       []string
   OtaInstalled    bool
   CellOperator    string
   Roaming         string
   TimeZone        string
   TouchScreen     int32
   Keyboard        int32
   Navigation      int32
   ScreenLayout    int32
   Screen          *DeviceInfoScreen
   GLVersion       int32
   GLExtensions    []string
   SharedLibraries []string
   Features        []string
   Locales         []string
}

type DeviceBuildInfo struct {
   *gpproto.AndroidBuildProto
   VersionRelease int
}

type DeviceInfoScreen struct {
   Density int32
   Width   int32
   Height  int32
}

func (i *DeviceInfo) GetDeviceFeatures() (ret []*gpproto.DeviceFeature) {
   var int0 int32
   for _, f := range i.Features {
      name := f
      ret = append(ret, &gpproto.DeviceFeature{Name: &name, Value: &int0})
   }
   return ret
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

func (client *GooglePlayClient) doAuthedReq(r *http.Request) (*gpproto.Payload, error) {
   data := client.AuthData
   r.Header.Set("Authorization", "Bearer "+data.AuthToken)
   r.Header.Set("X-DFE-Device-Id", data.GsfID)
   b, status, err := doReq(r)
   if err != nil {
      return nil, err
   }
   if status == 401 {
      return nil, GPTokenExpired
   }
   resp := &gpproto.ResponseWrapper{}
   if err := proto.Unmarshal(b, resp); err != nil {
      return nil, err
   }
   return resp.Payload, nil
}

const (
   UrlBase               = "https://android.clients.google.com"
   UrlFdfe               = UrlBase + "/fdfe"
   UrlAuth               = UrlBase + "/auth"
   UrlDetails            = UrlFdfe + "/details"
)

type GooglePlayClient struct {
   AuthData   *AuthData
   DeviceInfo *DeviceInfo
}

var (
   GPTokenExpired = errors.New("unauthorized, gp token expired")
   httpClient = &http.Client{}
)

type AuthData struct {
   AASToken                      string
   AuthToken                     string
   DFECookie                     string
   DeviceConfigToken             string
   Email                         string
   GsfID                         string
   Locale                        string
}
