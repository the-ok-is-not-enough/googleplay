package gplayapi

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "google.golang.org/protobuf/proto"
   "io"
   "net/http"
   "net/url"
   "os"
   "strconv"
   "strings"
   "time"
   gpproto "github.com/89z/googleplay/research/ignore"
)

type (
	App struct {
		PackageName        string
		AppInfo            *AppInfo
		CategoryImage      *gpproto.Image
		CategoryID         int
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

	AppInfo struct {
		AppInfoMap map[string]string
	}
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

func (client *GooglePlayClient) toc() (_ *gpproto.TocResponse, err error) {
	r, _ := http.NewRequest("GET", UrlToc, nil)
	payload, err := client.doAuthedReq(r)
	if err != nil {
		return
	}
	tocResp := payload.TocResponse
	if tocResp.TosContent != nil && tocResp.TosToken != nil {
		err = client.acceptTos(*tocResp.TosToken)
		if err != nil {
			return
		}
	}
	if tocResp.Cookie != nil {
		client.AuthData.DFECookie = *tocResp.Cookie
	}
	return
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

func (client *GooglePlayClient) setDefaultHeaders(r *http.Request) {
	data := client.AuthData
	r.Header.Set("Authorization", "Bearer "+data.AuthToken)
	r.Header.Set("User-Agent", client.DeviceInfo.GetUserAgent())
	r.Header.Set("X-DFE-Device-Id", data.GsfID)
	r.Header.Set("Accept-Language", "en-GB")
	r.Header.Set(
		"X-DFE-Encoded-Targets",
		"CAESN/qigQYC2AMBFfUbyA7SM5Ij/CvfBoIDgxHqGP8R3xzIBvoQtBKFDZ4HAY4FrwSVMasHBO0O2Q8akgYRAQECAQO7AQEpKZ0CnwECAwRrAQYBr9PPAoK7sQMBAQMCBAkIDAgBAwEDBAICBAUZEgMEBAMLAQEBBQEBAcYBARYED+cBfS8CHQEKkAEMMxcBIQoUDwYHIjd3DQ4MFk0JWGYZEREYAQOLAYEBFDMIEYMBAgICAgICOxkCD18LGQKEAcgDBIQBAgGLARkYCy8oBTJlBCUocxQn0QUBDkkGxgNZQq0BZSbeAmIDgAEBOgGtAaMCDAOQAZ4BBIEBKUtQUYYBQscDDxPSARA1oAEHAWmnAsMB2wFyywGLAxol+wImlwOOA80CtwN26A0WjwJVbQEJPAH+BRDeAfkHK/ABASEBCSAaHQemAzkaRiu2Ad8BdXeiAwEBGBUBBN4LEIABK4gB2AFLfwECAdoENq0CkQGMBsIBiQEtiwGgA1zyAUQ4uwS8AwhsvgPyAcEDF27vApsBHaICGhl3GSKxAR8MC6cBAgItmQYG9QIeywLvAeYBDArLAh8HASI4ELICDVmVBgsY/gHWARtcAsMBpALiAdsBA7QBpAJmIArpByn0AyAKBwHTARIHAX8D+AMBcRIBBbEDmwUBMacCHAciNp0BAQF0OgQLJDuSAh54kwFSP0eeAQQ4M5EBQgMEmwFXywFo0gFyWwMcapQBBugBPUW2AVgBKmy3AR6PAbMBGQxrUJECvQR+8gFoWDsYgQNwRSczBRXQAgtRswEW0ALMAREYAUEBIG6yATYCRE8OxgER8gMBvQEDRkwLc8MBTwHZAUOnAXiiBakDIbYBNNcCIUmuArIBSakBrgFHKs0EgwV/G3AD0wE6LgECtQJ4xQFwFbUCjQPkBS6vAQqEAUZF3QIM9wEhCoYCQhXsBCyZArQDugIziALWAdIBlQHwBdUErQE6qQaSA4EEIvYBHir9AQVLmgMCApsCKAwHuwgrENsBAjNYswEVmgIt7QJnN4wDEnta+wGfAcUBxgEtEFXQAQWdAUAeBcwBAQM7rAEJATJ0LENrdh73A6UBhAE+qwEeASxLZUMhDREuH0CGARbd7K0GlQo",
	)
	r.Header.Set(
		"X-DFE-Phenotype",
		"H4sIAAAAAAAAAB3OO3KjMAAA0KRNuWXukBkBQkAJ2MhgAZb5u2GCwQZbCH_EJ77QHmgvtDtbv-Z9_H63zXXU0NVPB1odlyGy7751Q3CitlPDvFd8lxhz3tpNmz7P92CFw73zdHU2Ie0Ad2kmR8lxhiErTFLt3RPGfJQHSDy7Clw10bg8kqf2owLokN4SecJTLoSwBnzQSd652_MOf2d1vKBNVedzg4ciPoLz2mQ8efGAgYeLou-l-PXn_7Sna1MfhHuySxt-4esulEDp8Sbq54CPPKjpANW-lkU2IZ0F92LBI-ukCKSptqeq1eXU96LD9nZfhKHdtjSWwJqUm_2r6pMHOxk01saVanmNopjX3YxQafC4iC6T55aRbC8nTI98AF_kItIQAJb5EQxnKTO7TZDWnr01HVPxelb9A2OWX6poidMWl16K54kcu_jhXw-JSBQkVcD_fPsLSZu6joIBAAA",
	)
	r.Header.Set("X-DFE-Client-Id", "am-android-google")
	r.Header.Set("X-DFE-Network-Type", "4")
	r.Header.Set("X-DFE-Content-Filters", "")
	r.Header.Set("X-Limit-Ad-Tracking-Enabled", "false")
	r.Header.Set("X-Ad-Id", "LawadaMera")
	r.Header.Set("X-DFE-UserLanguages", "en_GB")
	r.Header.Set("X-DFE-Request-Params", "timeoutMs=4000")

	if data.DeviceCheckInConsistencyToken != "" {
		r.Header.Set("X-DFE-Device-Checkin-Consistency-Token", data.DeviceCheckInConsistencyToken)
	}

	if data.DeviceConfigToken != "" {
		r.Header.Set("X-DFE-Device-Config-Token", data.DeviceConfigToken)
	}

	if data.DFECookie != "" {
		r.Header.Set("X-DFE-Cookie", data.DFECookie)
	}

	if client.DeviceInfo.SimOperator != "" {
		r.Header.Set("X-DFE-MCCMNC", client.DeviceInfo.SimOperator)
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
	params.Set("app", "com.android.vending")
	params.Set("client_sig", "38918a453d07199354f8b19af05ec6562ced5788")
	params.Set("callerPkg", "com.google.android.gms")
	params.Set("Token", client.AuthData.AASToken)
	params.Set("oauth2_foreground", "1")
	params.Set("token_request_options", "CAA4AVAB")
	params.Set("check_email", "1")
	params.Set("system_partition", "1")
}

type (
	DeviceInfo struct {
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

	DeviceBuildInfo struct {
		*gpproto.AndroidBuildProto
		VersionRelease int
	}

	DeviceInfoScreen struct {
		Density int32
		Width   int32
		Height  int32
	}
)

func (i *DeviceInfo) GetAuthUserAgent() string {
	return fmt.Sprintf("GoogleAuth/1.4 %s %s", i.Build.GetDevice(), i.Build.GetId())
}

func (i *DeviceInfo) GetUserAgent() string {
	params := []string{
		"api=3",
		"versionCode=81582300",
		"sdk=28",
		"device=" + i.Build.GetDevice(),
		"hardware=" + i.Build.GetDevice(),
		"product=" + i.Build.GetProduct(),
		"platformVersionRelease=" + strconv.Itoa(i.Build.VersionRelease),
		"model=" + i.Build.GetModel(),
		"buildId=" + i.Build.GetId(),
		"isWideScreen=0",
		"supportedAbis=" + strings.Join(i.Platforms, ";"),
	}
	return "Android-Finsky/15.8.23-all [0] [PR] 259261889 (" + strings.Join(params, ",") + ")"
}

func (i *DeviceInfo) GenerateAndroidCheckInRequest() *gpproto.AndroidCheckinRequest {
	var int0 int64 = 0
	timestamp := time.Now().Unix()
	i.Build.Timestamp = &timestamp

	return &gpproto.AndroidCheckinRequest{
		Id: &int0,
		Checkin: &gpproto.AndroidCheckinProto{
			Build:           i.Build.AndroidBuildProto,
			LastCheckinMsec: &int0,
			CellOperator:    &i.CellOperator,
			SimOperator:     &i.SimOperator,
			Roaming:         &i.Roaming,
			UserNumber:      ptrInt32(0),
		},
		Locale:              ptrStr("en_GB"),
		TimeZone:            &i.TimeZone,
		Version:             ptrInt32(3),
		DeviceConfiguration: i.GetDeviceConfigProto(),
		Fragment:            ptrInt32(0),
	}
}

func (i *DeviceInfo) GetDeviceConfigProto() *gpproto.DeviceConfigurationProto {
	var mem int64 = 8589935000
	return &gpproto.DeviceConfigurationProto{
		TouchScreen:            &i.TouchScreen,
		Keyboard:               &i.Keyboard,
		Navigation:             &i.Navigation,
		ScreenLayout:           &i.ScreenLayout,
		HasHardKeyboard:        ptrBool(false),
		HasFiveWayNavigation:   ptrBool(false),
		ScreenDensity:          &i.Screen.Density,
		GlEsVersion:            &i.GLVersion,
		SystemSharedLibrary:    i.SharedLibraries,
		SystemAvailableFeature: i.Features,
		NativePlatform:         i.Platforms,
		ScreenWidth:            &i.Screen.Width,
		ScreenHeight:           &i.Screen.Height,
		SystemSupportedLocale:  i.Locales,
		GlExtension:            i.GLExtensions,
		DeviceClass:            ptrInt32(0),
		LowRamDevice:           ptrInt32(0),
		TotalMemoryBytes:       &mem,
		MaxNumOf_CPUCores:      ptrInt32(8),
		DeviceFeature:          i.GetDeviceFeatures(),
	}
}

func (i *DeviceInfo) GetDeviceFeatures() (ret []*gpproto.DeviceFeature) {
	var int0 int32 = 0
	for _, f := range i.Features {
		name := f
		ret = append(ret, &gpproto.DeviceFeature{Name: &name, Value: &int0})
	}
	return ret
}
var (
	unknown               = ptrStr("unknown")
	androidGoogle         = ptrStr("android-google")
	falsePtr              = ptrBool(false)
	roaming               = "mobile-notroaming"
	glVersion       int32 = 196610
	glExtensions          = []string{"GL_AMD_compressed_ATC_texture", "GL_AMD_performance_monitor", "GL_ANDROID_extension_pack_es31a", "GL_APPLE_texture_2D_limited_npot", "GL_ARB_vertex_buffer_object", "GL_ARM_shader_framebuffer_fetch_depth_stencil", "GL_EXT_EGL_image_array", "GL_EXT_EGL_image_external_wrap_modes", "GL_EXT_EGL_image_storage", "GL_EXT_YUV_target", "GL_EXT_blend_func_extended", "GL_EXT_blit_framebuffer_params", "GL_EXT_buffer_storage", "GL_EXT_clip_control", "GL_EXT_clip_cull_distance", "GL_EXT_color_buffer_float", "GL_EXT_color_buffer_half_float", "GL_EXT_copy_image", "GL_EXT_debug_label", "GL_EXT_debug_marker", "GL_EXT_discard_framebuffer", "GL_EXT_disjoint_timer_query", "GL_EXT_draw_buffers_indexed", "GL_EXT_external_buffer", "GL_EXT_geometry_shader", "GL_EXT_gpu_shader5", "GL_EXT_memory_object", "GL_EXT_memory_object_fd", "GL_EXT_multisampled_render_to_texture", "GL_EXT_multisampled_render_to_texture2", "GL_EXT_primitive_bounding_box", "GL_EXT_protected_textures", "GL_EXT_robustness", "GL_EXT_sRGB", "GL_EXT_sRGB_write_control", "GL_EXT_shader_framebuffer_fetch", "GL_EXT_shader_io_blocks", "GL_EXT_shader_non_constant_global_initializers", "GL_EXT_tessellation_shader", "GL_EXT_texture_border_clamp", "GL_EXT_texture_buffer", "GL_EXT_texture_cube_map_array", "GL_EXT_texture_filter_anisotropic", "GL_EXT_texture_format_BGRA8888", "GL_EXT_texture_format_sRGB_override", "GL_EXT_texture_norm16", "GL_EXT_texture_sRGB_R8", "GL_EXT_texture_sRGB_decode", "GL_EXT_texture_type_2_10_10_10_REV", "GL_KHR_blend_equation_advanced", "GL_KHR_blend_equation_advanced_coherent", "GL_KHR_debug", "GL_KHR_no_error", "GL_KHR_robust_buffer_access_behavior", "GL_KHR_texture_compression_astc_hdr", "GL_KHR_texture_compression_astc_ldr", "GL_NV_shader_noperspective_interpolation", "GL_OES_EGL_image", "GL_OES_EGL_image_external", "GL_OES_EGL_image_external_essl3", "GL_OES_EGL_sync", "GL_OES_blend_equation_separate", "GL_OES_blend_func_separate", "GL_OES_blend_subtract", "GL_OES_compressed_ETC1_RGB8_texture", "GL_OES_compressed_paletted_texture", "GL_OES_depth24", "GL_OES_depth_texture", "GL_OES_depth_texture_cube_map", "GL_OES_draw_texture", "GL_OES_element_index_uint", "GL_OES_framebuffer_object", "GL_OES_get_program_binary", "GL_OES_matrix_palette", "GL_OES_packed_depth_stencil", "GL_OES_point_size_array", "GL_OES_point_sprite", "GL_OES_read_format", "GL_OES_rgb8_rgba8", "GL_OES_sample_shading", "GL_OES_sample_variables", "GL_OES_shader_image_atomic", "GL_OES_shader_multisample_interpolation", "GL_OES_standard_derivatives", "GL_OES_stencil_wrap", "GL_OES_surfaceless_context", "GL_OES_texture_3D", "GL_OES_texture_compression_astc", "GL_OES_texture_cube_map", "GL_OES_texture_env_crossbar", "GL_OES_texture_float", "GL_OES_texture_float_linear", "GL_OES_texture_half_float", "GL_OES_texture_half_float_linear", "GL_OES_texture_mirrored_repeat", "GL_OES_texture_npot", "GL_OES_texture_stencil8", "GL_OES_texture_storage_multisample_2d_array", "GL_OES_vertex_array_object", "GL_OES_vertex_half_float", "GL_OVR_multiview", "GL_OVR_multiview2", "GL_OVR_multiview_multisampled_render_to_texture", "GL_QCOM_alpha_test", "GL_QCOM_extended_get", "GL_QCOM_shader_framebuffer_fetch_noncoherent", "GL_QCOM_texture_foveated", "GL_QCOM_tiled_rendering"}
	sharedLibraries       = []string{"android.ext.services", "android.ext.shared", "android.test.base", "android.test.mock", "android.test.runner", "com.android.future.usb.accessory", "com.android.ims.rcsmanager", "com.android.location.provider", "com.android.media.remotedisplay", "com.android.mediadrm.signer", "com.google.android.camera.experimental2018", "com.google.android.dialer.support", "com.google.android.gms", "com.google.android.hardwareinfo", "com.google.android.lowpowermonitordevicefactory", "com.google.android.lowpowermonitordeviceinterface", "com.google.android.maps", "com.google.android.poweranomalydatafactory", "com.google.android.poweranomalydatamodeminterface", "com.qti.snapdragon.sdk.display", "com.qualcomm.embmslibrary", "com.qualcomm.qcrilhook", "com.qualcomm.qti.QtiTelephonyServicelibrary", "com.qualcomm.qti.imscmservice@1.0-java", "com.qualcomm.qti.lpa.uimlpalibrary", "com.qualcomm.qti.ltedirectdiscoverylibrary", "com.qualcomm.qti.remoteSimlock.uimremotesimlocklibrary", "com.qualcomm.qti.uim.uimservicelibrary", "com.quicinc.cne", "com.quicinc.cneapiclient", "com.verizon.embms", "com.verizon.provider", "com.vzw.apnlib", "javax.obex", "org.apache.http.legacy"}
	features              = []string{"android.hardware.audio.low_latency", "android.hardware.audio.output", "android.hardware.audio.pro", "android.hardware.bluetooth", "android.hardware.bluetooth_le", "android.hardware.camera", "android.hardware.camera.any", "android.hardware.camera.autofocus", "android.hardware.camera.capability.manual_post_processing", "android.hardware.camera.capability.manual_sensor", "android.hardware.camera.capability.raw", "android.hardware.camera.flash", "android.hardware.camera.front", "android.hardware.camera.level.full", "android.hardware.faketouch", "android.hardware.fingerprint", "android.hardware.location", "android.hardware.location.gps", "android.hardware.location.network", "android.hardware.microphone", "android.hardware.nfc", "android.hardware.nfc.any", "android.hardware.nfc.hce", "android.hardware.nfc.hcef", "android.hardware.opengles.aep", "android.hardware.ram.normal", "android.hardware.screen.landscape", "android.hardware.screen.portrait", "android.hardware.sensor.accelerometer", "android.hardware.sensor.assist", "android.hardware.sensor.barometer", "android.hardware.sensor.compass", "android.hardware.sensor.gyroscope", "android.hardware.sensor.hifi_sensors", "android.hardware.sensor.light", "android.hardware.sensor.proximity", "android.hardware.sensor.stepcounter", "android.hardware.sensor.stepdetector", "android.hardware.strongbox_keystore", "android.hardware.telephony", "android.hardware.telephony.carrierlock", "android.hardware.telephony.cdma", "android.hardware.telephony.euicc", "android.hardware.telephony.gsm", "android.hardware.touchscreen", "android.hardware.touchscreen.multitouch", "android.hardware.touchscreen.multitouch.distinct", "android.hardware.touchscreen.multitouch.jazzhand", "android.hardware.usb.accessory", "android.hardware.usb.host", "android.hardware.vulkan.compute", "android.hardware.vulkan.level", "android.hardware.vulkan.version", "android.hardware.wifi", "android.hardware.wifi.aware", "android.hardware.wifi.direct", "android.hardware.wifi.passpoint", "android.hardware.wifi.rtt", "android.software.activities_on_secondary_displays", "android.software.app_widgets", "android.software.autofill", "android.software.backup", "android.software.cant_save_state", "android.software.companion_device_setup", "android.software.connectionservice", "android.software.cts", "android.software.device_admin", "android.software.device_id_attestation", "android.software.file_based_encryption", "android.software.home_screen", "android.software.input_methods", "android.software.live_wallpaper", "android.software.managed_users", "android.software.midi", "android.software.picture_in_picture", "android.software.print", "android.software.securely_removes_users", "android.software.sip", "android.software.sip.voip", "android.software.verified_boot", "android.software.voice_recognizers", "android.software.webview", "com.google.android.apps.dialer.SUPPORTED", "com.google.android.apps.photos.PIXEL_2019_MIDYEAR_PRELOAD", "com.google.android.feature.EXCHANGE_6_2", "com.google.android.feature.GOOGLE_BUILD", "com.google.android.feature.GOOGLE_EXPERIENCE", "com.google.android.feature.PIXEL_2017_EXPERIENCE", "com.google.android.feature.PIXEL_2018_EXPERIENCE", "com.google.android.feature.PIXEL_2019_MIDYEAR_EXPERIENCE", "com.google.android.feature.PIXEL_EXPERIENCE", "com.google.android.feature.TURBO_PRELOAD", "com.google.android.feature.WELLBEING", "com.google.android.feature.ZERO_TOUCH", "com.google.hardware.camera.easel_2018", "com.verizon.hardware.telephony.ehrpd", "com.verizon.hardware.telephony.lte"}
	locales               = []string{"af", "af_ZA", "am", "am_ET", "ar", "ar_EG", "ar_SA", "ar_XB", "as", "ast", "az", "be", "be_BY", "bg", "bg_BG", "bh_IN", "bn", "bs", "ca", "ca_ES", "cs", "cs_CZ", "cy_GB", "da", "da_DK", "de", "de_DE", "el", "el_GR", "en", "en_AU", "en_CA", "en_GB", "en_IN", "en_US", "en_XA", "es", "es_ES", "es_US", "et", "et_EE", "eu", "fa", "fa_IR", "fi", "fi_FI", "fil", "fil_PH", "fr", "fr_CA", "fr_FR", "gl", "gl_ES", "gu", "hi", "hi_IN", "hr", "hr_HR", "hu", "hu_HU", "hy", "in", "in_ID", "is", "it", "it_IT", "iw", "iw_IL", "ja", "ja_JP", "ka", "kab_DZ", "kk", "km", "kn", "ko", "ko_KR", "ky", "lo", "lt", "lt_LT", "lv", "lv_LV", "mk", "ml", "mn", "mr", "ms", "ms_MY", "my", "nb", "nb_NO", "ne", "nl", "nl_NL", "or", "pa", "pa_IN", "pl", "pl_PL", "pt", "pt_BR", "pt_PT", "ro", "ro_RO", "ru", "ru_RU", "sc_IT", "si", "sk", "sk_SK", "sl", "sl_SI", "sq", "sr", "sr_Latn", "sr_RS", "sv", "sv_SE", "sw", "sw_TZ", "ta", "te", "th", "th_TH", "tr", "tr_TR", "uk", "uk_UA", "ur", "uz", "vi", "vi_VN", "zh_CN", "zh_HK", "zh_TW", "zu", "zu_ZA"}

	sargo = ptrStr("sargo")

	// Pixel3a is default device which uses arm64 and Android 9 / SDK 28
	Pixel3a = &DeviceInfo{
		Build: &DeviceBuildInfo{
			AndroidBuildProto: &gpproto.AndroidBuildProto{
				BuildProduct:   sargo,
				Radio:          ptrStr("g670-00011-190411-B-5457439"),
				Bootloader:     ptrStr("b4s4-0.1-5613380"),
				Carrier:        ptrStr("google"),
				Device:         sargo,
				Model:          ptrStr("Pixel 3a"),
				Manufacturer:   ptrStr("google"),
				Product:        sargo,
				Id:             ptrStr("PQ3B.190705.003"),
				SdkVersion:     ptrInt32(28),
				Client:         androidGoogle,
				GoogleServices: ptrInt32(203615028),
				OtaInstalled:   falsePtr,
			},
			VersionRelease: 9,
		},
		SimOperator:  "20815",
		Platforms:    []string{"arm64-v8a", "armeabi-v7a", "armeabi"},
		CellOperator: "334050",
		Roaming:      roaming,
		TimeZone:     "America/Mexico_City",
		TouchScreen:  3,
		Keyboard:     1,
		Navigation:   1,
		ScreenLayout: 2,
		Screen: &DeviceInfoScreen{
			Density: 490,
			Width:   1080,
			Height:  2073,
		},
		GLVersion:       glVersion,
		GLExtensions:    glExtensions,
		SharedLibraries: sharedLibraries,
		Features:        features,
		Locales:         locales,
	}

	santoni = ptrStr("santoni")

	// Redmi4 is device which uses arm and Android 10 / SDK 29
	Redmi4 = &DeviceInfo{
		Build: &DeviceBuildInfo{
			AndroidBuildProto: &gpproto.AndroidBuildProto{
				BuildProduct:   santoni,
				Radio:          ptrStr("MPSS.TA.2.3.c1-00395-8953_GEN_PACK-1_V048"),
				Bootloader:     unknown,
				Carrier:        ptrStr("Xiaomi"),
				Device:         santoni,
				Model:          ptrStr("Redmi 4"),
				Manufacturer:   ptrStr("Xiaomi"),
				Product:        santoni,
				Id:             ptrStr("QQ3A.200805.001"),
				SdkVersion:     ptrInt32(29),
				Client:         androidGoogle,
				GoogleServices: ptrInt32(203315024),
				OtaInstalled:   falsePtr,
			},
			VersionRelease: 10,
		},
		SimOperator:  "21601",
		Platforms:    []string{"armeabi-v7a", "armeabi"},
		CellOperator: "21601",
		Roaming:      roaming,
		TimeZone:     "Europe/Budapest",
		TouchScreen:  3,
		Keyboard:     1,
		Navigation:   1,
		ScreenLayout: 3,
		Screen: &DeviceInfoScreen{
			Density: 224,
			Width:   720,
			Height:  1280,
		},
		GLVersion:       glVersion,
		GLExtensions:    glExtensions,
		SharedLibraries: sharedLibraries,
		Features:        features,
		Locales:         locales,
	}

	glExtensionsEmulator = []string{"ANDROID_EMU_CHECKSUM_HELPER_v1", "ANDROID_EMU_YUV_Cache", "ANDROID_EMU_deferred_vulkan_commands", "ANDROID_EMU_direct_mem", "ANDROID_EMU_dma_v1", "ANDROID_EMU_gles_max_version_3_0", "ANDROID_EMU_has_shared_slots_host_memory_allocator", "ANDROID_EMU_host_composition_v1", "ANDROID_EMU_host_composition_v2", "ANDROID_EMU_host_side_tracing", "ANDROID_EMU_native_sync_v2", "ANDROID_EMU_native_sync_v3", "ANDROID_EMU_native_sync_v4", "ANDROID_EMU_read_color_buffer_dma", "ANDROID_EMU_sync_buffer_data", "ANDROID_EMU_vulkan", "ANDROID_EMU_vulkan_async_queue_submit", "ANDROID_EMU_vulkan_create_resources_with_requirements", "ANDROID_EMU_vulkan_free_memory_sync", "ANDROID_EMU_vulkan_ignored_handles", "ANDROID_EMU_vulkan_null_optional_strings", "ANDROID_EMU_vulkan_shader_float16_int8", "GL_APPLE_texture_format_BGRA8888", "GL_EXT_color_buffer_float", "GL_EXT_color_buffer_half_float", "GL_EXT_debug_marker", "GL_EXT_robustness", "GL_EXT_texture_format_BGRA8888", "GL_KHR_texture_compression_astc_ldr", "GL_OES_EGL_image", "GL_OES_EGL_image_external", "GL_OES_EGL_image_external_essl3", "GL_OES_EGL_sync", "GL_OES_blend_equation_separate", "GL_OES_blend_func_separate", "GL_OES_blend_subtract", "GL_OES_byte_coordinates", "GL_OES_compressed_ETC1_RGB8_texture", "GL_OES_compressed_paletted_texture", "GL_OES_depth24", "GL_OES_depth32", "GL_OES_depth_texture", "GL_OES_draw_texture", "GL_OES_element_index_uint", "GL_OES_fbo_render_mipmap", "GL_OES_framebuffer_object", "GL_OES_packed_depth_stencil", "GL_OES_point_size_array", "GL_OES_point_sprite", "GL_OES_rgb8_rgba8", "GL_OES_single_precision", "GL_OES_stencil1", "GL_OES_stencil4", "GL_OES_stencil8", "GL_OES_stencil_wrap", "GL_OES_texture_cube_map", "GL_OES_texture_env_crossbar", "GL_OES_texture_float", "GL_OES_texture_float_linear", "GL_OES_texture_half_float", "GL_OES_texture_half_float_linear", "GL_OES_texture_mirored_repeat", "GL_OES_texture_npot", "GL_OES_vertex_array_object", "GL_OES_vertex_half_float"}
	android              = ptrStr("Android")

	// Emulator_x86_64 is device which uses x86_64 and Android 11 / SDK 30
	Emulator_x86_64 = &DeviceInfo{
		Build: &DeviceBuildInfo{
			AndroidBuildProto: &gpproto.AndroidBuildProto{
				BuildProduct:   ptrStr("sdk_phone_x86_64"),
				Radio:          ptrStr("1.0.0.0"),
				Bootloader:     unknown,
				Carrier:        android,
				Device:         ptrStr("generic_x86_64"),
				Model:          ptrStr("Android SDK built for x86_64"),
				Manufacturer:   unknown,
				Product:        ptrStr("sdk_phone_x86_64"),
				Id:             ptrStr("RSR1.210722.012"),
				SdkVersion:     ptrInt32(30),
				Client:         androidGoogle,
				GoogleServices: ptrInt32(203019037),
				OtaInstalled:   falsePtr,
			},
			VersionRelease: 11,
		},
		SimOperator:  "38",
		Platforms:    []string{"x86_64", "x86"},
		CellOperator: "310",
		Roaming:      roaming,
		TimeZone:     "UTC-10",
		TouchScreen:  3,
		Keyboard:     2,
		Navigation:   2,
		ScreenLayout: 2,
		Screen: &DeviceInfoScreen{
			Density: 420,
			Width:   1080,
			Height:  1794,
		},
		GLVersion:       glVersion,
		GLExtensions:    glExtensionsEmulator,
		SharedLibraries: sharedLibraries,
		Features:        features,
		Locales:         locales,
	}

	// Emulator_x86 is device which uses x86 and Android 11 / SDK 30
	Emulator_x86 = &DeviceInfo{
		Build: &DeviceBuildInfo{
			AndroidBuildProto: &gpproto.AndroidBuildProto{
				BuildProduct:   ptrStr("sdk_phone_x86"),
				Radio:          ptrStr("1.0.0.0"),
				Bootloader:     unknown,
				Carrier:        android,
				Device:         ptrStr("generic_x86"),
				Model:          ptrStr("Android SDK built for x86"),
				Manufacturer:   unknown,
				Product:        ptrStr("sdk_phone_x86"),
				Id:             ptrStr("RSR1.210722.012"),
				SdkVersion:     ptrInt32(30),
				Client:         androidGoogle,
				GoogleServices: ptrInt32(203019037),
				OtaInstalled:   falsePtr,
			},
			VersionRelease: 11,
		},
		SimOperator:  "38",
		Platforms:    []string{"x86"},
		CellOperator: "310",
		Roaming:      roaming,
		TimeZone:     "UTC-10",
		TouchScreen:  3,
		Keyboard:     2,
		Navigation:   2,
		ScreenLayout: 2,
		Screen: &DeviceInfoScreen{
			Density: 420,
			Width:   1080,
			Height:  1794,
		},
		GLVersion:       glVersion,
		GLExtensions:    glExtensionsEmulator,
		SharedLibraries: sharedLibraries,
		Features:        features,
		Locales:         locales,
	}
)
type GooglePlayClient struct {
	AuthData   *AuthData
	DeviceInfo *DeviceInfo

	// SessionFile if SessionFile is set then session will be saved to it after modification
	SessionFile string
}

var (
   ErrTokenExpired = errors.New("unauthorized, gp token expired")
   httpClient = &http.Client{}
)

func NewClient(email, aasToken string) (*GooglePlayClient, error) {
   return NewClientWithDeviceInfo(email, aasToken, Pixel3a)
}

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

	_, err = client.toc()
	return
}

func (client *GooglePlayClient) SaveSession(file string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(client.AuthData)
}

func LoadSession(file string) (*GooglePlayClient, error) {
	return LoadSessionWithDeviceInfo(file, Pixel3a)
}

func LoadSessionWithDeviceInfo(file string, deviceInfo *DeviceInfo) (client *GooglePlayClient, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	client = &GooglePlayClient{DeviceInfo: deviceInfo}
	err = json.NewDecoder(f).Decode(&client.AuthData)
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
	UrlCheckIn            = UrlBase + "/checkin"
	UrlDetails            = UrlFdfe + "/details"
	UrlDelivery           = UrlFdfe + "/delivery"
	UrlPurchase           = UrlFdfe + "/purchase"
	UrlToc                = UrlFdfe + "/toc"
	UrlTosAccept          = UrlFdfe + "/acceptTos"
	UrlUploadDeviceConfig = UrlFdfe + "/uploadDeviceConfig"
)

var ErrNilPayload = errors.New("got nil payload from google play")
func ptrBool(b bool) *bool {
	return &b
}

func ptrStr(str string) *string {
	return &str
}

func ptrInt32(i int32) *int32 {
	return &i
}

func doReq(r *http.Request) ([]byte, int, error) {
	res, err := httpClient.Do(r)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	return b, res.StatusCode, err
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
		return nil, ErrTokenExpired
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
	if err == ErrTokenExpired {
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
