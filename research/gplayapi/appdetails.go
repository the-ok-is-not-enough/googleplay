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
   "net/url"
   "strconv"
   "strings"
)

func (client *GooglePlayClient) GenerateGPToken() (string, error) {
   params := &url.Values{}
   params.Set("sdk_version", strconv.Itoa(int(client.DeviceInfo.Build.GetSdkVersion())))
   params.Set("email", client.AuthData.Email)
   params.Set("google_play_services_version", strconv.Itoa(int(client.DeviceInfo.Build.GetGoogleServices())))
   params.Set("callerSig", "38918a453d07199354f8b19af05ec6562ced5788")
   params.Set("Token", client.AuthData.AASToken)
   params.Set("check_email", "1")
   params.Set("client_sig", "38918a453d07199354f8b19af05ec6562ced5788")
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

func (client *GooglePlayClient) checkIn(req *gpproto.AndroidCheckinRequest) (resp *gpproto.AndroidCheckinResponse, err error) {
   b, err := proto.Marshal(req)
   if err != nil {
      return
   }
   r, _ := http.NewRequest("POST", UrlBase + "/checkin", bytes.NewReader(b))
   r.Header.Set("Content-Type", "application/x-protobuffer")
   b, _, err = doReq(r)
   if err != nil {
      return
   }
   resp = &gpproto.AndroidCheckinResponse{}
   err = proto.Unmarshal(b, resp)
   return
}

func (client *GooglePlayClient) GetAppDetails(packageName string) (*App, error) {
   r, _ := http.NewRequest("GET", UrlDetails+"?doc="+packageName, nil)
   payload, err := client.doAuthedReq(r)
   if err != nil {
      return nil, err
   }
   return buildAppFromItem(payload.DetailsResponse.Item), nil
   
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
      PackageName:      *item.Id,
      CategoryID:       int(item.GetCategoryId()),
      DisplayName:      item.GetTitle(),
      Description:      item.GetDescriptionHtml(),
      ShortDescription: item.GetPromotionalDescription(),
      ShareUrl:         item.GetShareUrl(),
      VersionName:      details.GetVersionString(),
      VersionCode:      int(details.GetVersionCode()),
      CategoryName:     details.GetCategoryName(),
      Size:             details.GetInfoDownloadSize(),
      DownloadString:   details.GetDownloadLabelAbbreviated(),
      Changes:          details.GetRecentChangesHtml(),
      ContainsAds:      details.InstallNotes != nil,
      EarlyAccess:      details.EarlyAccessInfo != nil,
      DeveloperName:    details.GetDeveloperName(),
      TargetSdk:        int(details.GetTargetSdkVersion()),
      UpdatedOn:        details.GetInfoUpdatedOn(),
   }
   if len(item.Offer) > 0 {
   offer := item.Offer[0]
   app.OfferType = offer.GetOfferType()
   app.IsFree = offer.GetMicros() == 0
   app.Price = offer.GetFormattedAmount()
   }
   if app.DeveloperName == "" {
   app.DeveloperName = item.GetCreator()
   }
   if details.InstantLink != nil {
   app.InstantAppLink = details.GetInstantLink()
   }
   parseAppInfo(app, item)
   parseStreamUrls(app, item)
   parseImages(app, item)
   return app
}

func parseAppInfo(app *App, item *gpproto.Item) {
	if item.AppInfo != nil {
		app.AppInfo = &AppInfo{map[string]string{}}
		for _, s := range item.AppInfo.Section {
			if s.Label != nil && s.Container != nil && s.GetContainer().Description != nil {
				app.AppInfo.AppInfoMap[s.GetLabel()] = s.GetContainer().GetDescription()
			}
		}
	}
}

func parseStreamUrls(app *App, item *gpproto.Item) {
	if item.Annotations != nil {
		app.LiveStreamUrl = item.Annotations.GetLiveStreamUrl()
		app.PromotionStreamUrl = item.Annotations.GetPromotionStreamUrl()
	}
}

func parseImages(app *App, item *gpproto.Item) {
	for _, image := range item.Image {
		switch image.GetImageType() {
		case ImageTypeCategoryIcon:
			app.CategoryImage = image
		case ImageTypeAppIcon:
			app.IconImage = image
		case ImageTypeYoutubeVideoThumbnail:
			app.Video = image
		case ImageTypePlayStorePageBackground:
			app.CoverImage = image
		case ImageTypeAppScreenshot:
			app.Screenshots = append(app.Screenshots, image)
		}
	}

	if len(app.Screenshots) == 0 {
		if item.Annotations != nil && item.Annotations.SectionImage != nil {
			for _, imageContainer := range item.Annotations.SectionImage.ImageContainer {
				app.Screenshots = append(app.Screenshots, imageContainer.Image)
			}
		}
	}
}

var Pixel3a = &DeviceInfo{
   Build: &DeviceBuildInfo{
      AndroidBuildProto: &gpproto.AndroidBuildProto{
         Radio:          ptrStr("g670-00011-190411-B-5457439"),
         Bootloader:     ptrStr("b4s4-0.1-5613380"),
         SdkVersion:     ptrInt32(29),
      },
   },
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

func (i *DeviceInfo) GetDeviceFeatures() (ret []*gpproto.DeviceFeature) {
   var int0 int32
   for _, f := range i.Features {
      name := f
      ret = append(ret, &gpproto.DeviceFeature{Name: &name, Value: &int0})
   }
   return ret
}

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
   data := client.AuthData
   r.Header.Set("Authorization", "Bearer "+data.AuthToken)
   r.Header.Set("X-DFE-Device-Id", data.GsfID)
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

func (client *GooglePlayClient) doAuthedReq(r *http.Request) (*gpproto.Payload, error) {
   res, err := client._doAuthedReq(r)
   if err == GPTokenExpired {
      tok, err := client.GenerateGPToken()
      if err != nil {
         return nil, err
      }
      client.AuthData.AuthToken = tok
      return client._doAuthedReq(r)
   } else if err != nil {
      return nil, err
   }
   return res, nil
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

type AuthData struct {
   AASToken                      string
   AuthToken                     string
   DFECookie                     string
   DeviceConfigToken             string
   Email                         string
   GsfID                         string
   Locale                        string
}
