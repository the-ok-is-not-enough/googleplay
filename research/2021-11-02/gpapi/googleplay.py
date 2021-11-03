from . import googleplay_pb2
from base64 import b64decode, urlsafe_b64encode
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives.asymmetric.utils import encode_dss_signature
from cryptography.hazmat.primitives.serialization import load_der_public_key
from datetime import datetime
from os import path
from re import match
from sys import version_info
from time import time
import requests

BASE = "https://android.clients.google.com/"
FDFE = BASE + "fdfe/"

ACCOUNT = "HOSTED_OR_GOOGLE"
AUTH_URL = BASE + "auth"
BROWSE_URL = FDFE + "browse"
BULK_URL = FDFE + "bulkDetails"
CHECKIN_URL = BASE + "checkin"
CONTENT_TYPE_PROTO = "application/x-protobuf"
CONTENT_TYPE_URLENC = "application/x-www-form-urlencoded; charset=UTF-8"
DELIVERY_URL = FDFE + "delivery"
DETAILS_URL = FDFE + "details"
HOME_URL = FDFE + "homeV2"
LIST_URL = FDFE + "list"
PURCHASE_URL = FDFE + "purchase"
REVIEWS_URL = FDFE + "rev"
SEARCH_URL = FDFE + "search"
UPLOAD_URL = FDFE + "uploadDeviceConfig"
ssl_verify = True

class LoginError(Exception):
    def __init__(self, value):
        self.value = value

    def __str__(self):
        return repr(self.value)

class GooglePlayAPI(object):
    def __init__(self, locale="en_US", timezone="UTC", device_codename="bacon",
                 proxies_config=None):
        self.authSubToken = None
        self.gsfId = None
        self.device_config_token = None
        self.deviceCheckinConsistencyToken = None
        self.dfeCookie = None
        self.proxies_config = proxies_config
        self.deviceBuilder = DeviceBuilder(device_codename)

    def setAuthSubToken(self, authSubToken):
        self.authSubToken = authSubToken

    def getHeaders(self, upload_fields=False):
        # NEED THIS
        headers = self.deviceBuilder.getBaseHeaders()
        if self.gsfId is not None:
            headers["X-DFE-Device-Id"] = "{0:x}".format(self.gsfId)
        if self.authSubToken is not None:
            headers["Authorization"] = "GoogleLogin auth=%s" % self.authSubToken
        if self.device_config_token is not None:
            headers["X-DFE-Device-Config-Token"] = self.device_config_token
        if self.deviceCheckinConsistencyToken is not None:
            headers["X-DFE-Device-Checkin-Consistency-Token"] = self.deviceCheckinConsistencyToken
        if self.dfeCookie is not None:
            headers["X-DFE-Cookie"] = self.dfeCookie
        return headers

    def checkin(self, email, ac2dmToken):
        headers = self.getHeaders()
        headers["Content-Type"] = CONTENT_TYPE_PROTO
        request = self.deviceBuilder.getAndroidCheckinRequest()
        stringRequest = request.SerializeToString()
        res = requests.post(CHECKIN_URL, data=stringRequest,
                            headers=headers, verify=ssl_verify,
                            proxies=self.proxies_config)
        print(res.status_code, res.url)
        response = googleplay_pb2.AndroidCheckinResponse()
        response.ParseFromString(res.content)
        self.deviceCheckinConsistencyToken = response.deviceCheckinConsistencyToken
        # checkin again to upload gfsid
        request.id = response.androidId
        request.securityToken = response.securityToken
        request.accountCookie.append("[" + email + "]")
        request.accountCookie.append(ac2dmToken)
        stringRequest = request.SerializeToString()
        res = requests.post(CHECKIN_URL,
                      data=stringRequest,
                      headers=headers,
                      verify=ssl_verify,
                      proxies=self.proxies_config)
        print(res.status_code, res.url)
        return response.androidId

    def uploadDeviceConfig(self):
        upload = googleplay_pb2.UploadDeviceConfigRequest()
        upload.deviceConfiguration.CopyFrom(self.deviceBuilder.getDeviceConfig())
        headers = self.getHeaders(upload_fields=True)
        stringRequest = upload.SerializeToString()
        response = requests.post(UPLOAD_URL, data=stringRequest,
            headers=headers,
            verify=ssl_verify,
            timeout=60,
            proxies=self.proxies_config)
        print(response.status_code, response.url)
        response = googleplay_pb2.ResponseWrapper.FromString(response.content)
        try:
            if response.payload.HasField('uploadDeviceConfigResponse'):
                self.device_config_token = response.payload.uploadDeviceConfigResponse
                self.device_config_token = self.device_config_token.uploadDeviceConfigToken
        except ValueError:
            pass

    def getAuthSubToken(self, email, passwd):
        requestParams = {
            "Email": email,
            "EncryptedPasswd": passwd,
            "droidguard_results": "dummy123",
            "add_account": "1",
            "callerSig": "38918a453d07199354f8b19af05ec6562ced5788",
            "client_sig": "38918a453d07199354f8b19af05ec6562ced5788",
            "has_permission": "1",
            "source": "android",
               }
        requestParams['service'] = 'androidmarket'
        requestParams['app'] = 'com.android.vending'
        headers = self.deviceBuilder.getAuthHeaders(self.gsfId)
        headers['app'] = 'com.android.vending'
        response = requests.post(AUTH_URL,
            data=requestParams,
            verify=ssl_verify,
            headers=headers,
            proxies=self.proxies_config)
        print(response.status_code, response.url)
        data = response.text.split()
        params = {}
        for d in data:
            if "=" not in d:
                continue
            k, v = d.split("=", 1)
            params[k.strip().lower()] = v.strip()
        if "token" in params:
            master_token = params["token"]
            second_round_token = self.getSecondRoundToken(master_token, requestParams)
            self.setAuthSubToken(second_round_token)
        elif "error" in params:
            raise LoginError("server says: " + params["error"])
        else:
            raise LoginError("auth token not found.")

    def getSecondRoundToken(self, first_token, params):
        if self.gsfId is not None:
            params['androidId'] = "{0:x}".format(self.gsfId)
        params['Token'] = first_token
        params['check_email'] = '1'
        params['token_request_options'] = 'CAA4AQ=='
        params['system_partition'] = '1'
        params['_opt_is_called_from_account_manager'] = '1'
        params.pop('Email')
        params.pop('EncryptedPasswd')
        headers = self.deviceBuilder.getAuthHeaders(self.gsfId)
        headers['app'] = 'com.android.vending'
        response = requests.post(AUTH_URL,
            data=params,
            headers=headers,
            verify=ssl_verify,
            proxies=self.proxies_config)
        print(response.status_code, response.url)
        data = response.text.split()
        params = {}
        for d in data:
            if "=" not in d:
                continue
            k, v = d.split("=", 1)
            params[k.strip().lower()] = v.strip()
        if "auth" in params:
            return params["auth"]
        elif "error" in params:
            raise LoginError("server says: " + params["error"])
        else:
            raise LoginError("Auth token not found.")

selfDevice = {
   'build.fingerprint': 'google/sailfish/sailfish:8.1.0/OPM4.171019.016.B1/4720843:user/release-keys', 'build.brand': 'Google', 'build.device': 'sailfish',
   'build.version.release': '8.1.0', 'touchscreen': '3', 'keyboard': '1', 'navigation': '1', 'screenlayout': '2', 'hashardkeyboard': 'false',
   'build.version.sdk_int': '27', 'build.model': 'Pixel', 'build.manufacturer': 'Google', 'build.product': 'sailfish', 'build.id': 'OPM2.171019.029.B1',
   'features': 'android.hardware.audio.low_latency,android.hardware.audio.output,android.hardware.audio.pro,android.hardware.bluetooth,android.hardware.bluetooth_le,android.hardware.camera,android.hardware.camera.any,android.hardware.camera.autofocus,android.hardware.camera.capability.manual_post_processing,android.hardware.camera.capability.manual_sensor,android.hardware.camera.capability.raw,android.hardware.camera.flash,android.hardware.camera.front,android.hardware.camera.level.full,android.hardware.faketouch,android.hardware.fingerprint,android.hardware.location,android.hardware.location.gps,android.hardware.location.network,android.hardware.microphone,android.hardware.nfc,android.hardware.nfc.any,android.hardware.nfc.hce,android.hardware.nfc.hcef,android.hardware.opengles.aep,android.hardware.ram.normal,android.hardware.screen.landscape,android.hardware.screen.portrait,android.hardware.sensor.accelerometer,android.hardware.sensor.barometer,android.hardware.sensor.compass,android.hardware.sensor.gyroscope,android.hardware.sensor.hifi_sensors,android.hardware.sensor.light,android.hardware.sensor.proximity,android.hardware.sensor.stepcounter,android.hardware.sensor.stepdetector,android.hardware.telephony,android.hardware.telephony.carrierlock,android.hardware.telephony.cdma,android.hardware.telephony.gsm,android.hardware.touchscreen,android.hardware.touchscreen.multitouch,android.hardware.touchscreen.multitouch.distinct,android.hardware.touchscreen.multitouch.jazzhand,android.hardware.usb.accessory,android.hardware.usb.host,android.hardware.vr.headtracking,android.hardware.vr.high_performance,android.hardware.vulkan.compute,android.hardware.vulkan.level,android.hardware.vulkan.version,android.hardware.wifi,android.hardware.wifi.direct,android.hardware.wifi.passpoint,android.software.activities_on_secondary_displays,android.software.app_widgets,android.software.autofill,android.software.backup,android.software.companion_device_setup,android.software.connectionservice,android.software.cts,android.software.device_admin,android.software.file_based_encryption,android.software.home_screen,android.software.input_methods,android.software.live_wallpaper,android.software.managed_users,android.software.midi,android.software.picture_in_picture,android.software.print,android.software.securely_removes_users,android.software.sip,android.software.sip.voip,android.software.voice_recognizers,android.software.vr.mode,android.software.webview,com.google.android.apps.dialer.SUPPORTED,com.google.android.apps.photos.NEXUS_PRELOAD,com.google.android.apps.photos.nexus_preload,com.google.android.feature.EXCHANGE_6_2,com.google.android.feature.GOOGLE_BUILD,com.google.android.feature.GOOGLE_EXPERIENCE,com.google.android.feature.PIXEL_EXPERIENCE,com.nxp.mifare,com.verizon.hardware.telephony.ehrpd,com.verizon.hardware.telephony.lte,org.lineageos.android,org.lineageos.audio,org.lineageos.hardware,org.lineageos.livedisplay,org.lineageos.performance,org.lineageos.profiles,org.lineageos.style,org.lineageos.weather,projekt.substratum.theme',
   'gl.extensions': 'GL_AMD_compressed_ATC_texture,GL_AMD_performance_monitor,GL_ANDROID_extension_pack_es31a,GL_APPLE_texture_2D_limited_npot,GL_ARB_vertex_buffer_object,GL_ARM_shader_framebuffer_fetch_depth_stencil,GL_EXT_EGL_image_array,GL_EXT_YUV_target,GL_EXT_blit_framebuffer_params,GL_EXT_buffer_storage,GL_EXT_clip_cull_distance,GL_EXT_color_buffer_float,GL_EXT_color_buffer_half_float,GL_EXT_copy_image,GL_EXT_debug_label,GL_EXT_debug_marker,GL_EXT_discard_framebuffer,GL_EXT_disjoint_timer_query,GL_EXT_draw_buffers_indexed,GL_EXT_external_buffer,GL_EXT_geometry_shader,GL_EXT_gpu_shader5,GL_EXT_multisampled_render_to_texture,GL_EXT_multisampled_render_to_texture2,GL_EXT_primitive_bounding_box,GL_EXT_protected_textures,GL_EXT_robustness,GL_EXT_sRGB,GL_EXT_sRGB_write_control,GL_EXT_shader_framebuffer_fetch,GL_EXT_shader_io_blocks,GL_EXT_shader_non_constant_global_initializers,GL_EXT_tessellation_shader,GL_EXT_texture_border_clamp,GL_EXT_texture_buffer,GL_EXT_texture_cube_map_array,GL_EXT_texture_filter_anisotropic,GL_EXT_texture_format_BGRA8888,GL_EXT_texture_norm16,GL_EXT_texture_sRGB_R8,GL_EXT_texture_sRGB_decode,GL_EXT_texture_type_2_10_10_10_REV,GL_KHR_blend_equation_advanced,GL_KHR_blend_equation_advanced_coherent,GL_KHR_debug,GL_KHR_no_error,GL_KHR_texture_compression_astc_hdr,GL_KHR_texture_compression_astc_ldr,GL_NV_shader_noperspective_interpolation,GL_OES_EGL_image,GL_OES_EGL_image_external,GL_OES_EGL_image_external_essl3,GL_OES_EGL_sync,GL_OES_blend_equation_separate,GL_OES_blend_func_separate,GL_OES_blend_subtract,GL_OES_compressed_ETC1_RGB8_texture,GL_OES_compressed_paletted_texture,GL_OES_depth24,GL_OES_depth_texture,GL_OES_depth_texture_cube_map,GL_OES_draw_texture,GL_OES_element_index_uint,GL_OES_framebuffer_object,GL_OES_get_program_binary,GL_OES_matrix_palette,GL_OES_packed_depth_stencil,GL_OES_point_size_array,GL_OES_point_sprite,GL_OES_read_format,GL_OES_rgb8_rgba8,GL_OES_sample_shading,GL_OES_sample_variables,GL_OES_shader_image_atomic,GL_OES_shader_multisample_interpolation,GL_OES_standard_derivatives,GL_OES_stencil_wrap,GL_OES_surfaceless_context,GL_OES_texture_3D,GL_OES_texture_compression_astc,GL_OES_texture_cube_map,GL_OES_texture_env_crossbar,GL_OES_texture_float,GL_OES_texture_float_linear,GL_OES_texture_half_float,GL_OES_texture_half_float_linear,GL_OES_texture_mirrored_repeat,GL_OES_texture_npot,GL_OES_texture_stencil8,GL_OES_texture_storage_multisample_2d_array,GL_OES_vertex_array_object,GL_OES_vertex_half_float,GL_OVR_multiview,GL_OVR_multiview2,GL_OVR_multiview_multisampled_render_to_texture,GL_QCOM_alpha_test,GL_QCOM_extended_get,GL_QCOM_framebuffer_foveated,GL_QCOM_shader_framebuffer_fetch_noncoherent,GL_QCOM_tiled_rendering',
   'gsf.version': '12685052', 'vending.version': '81031200', 'vending.versionstring': '10.3.12-all [0] [PR] 198814133', 'celloperator': '311480',
   'hasfivewaynavigation': 'false', 'gl.version': '196610', 'screen.density': '420', 'screen.width': '1794', 'screen.height': '1080',
   'locales': 'af,af_ZA,am,am_ET,ar,ar_EG,ar_XB,ast,az,az_AZ,be,be_BY,bg,bg_BG,bn,bs,ca,ca_ES,cs,cs_CZ,da,da_DK,de,de_AT,de_CH,de_DE,de_LI,el,el_GR,en,en_AU,en_CA,en_GB,en_IN,en_NZ,en_SG,en_US,en_XA,en_XC,eo,es,es_ES,es_US,et,et_EE,eu,eu_ES,fa,fa_IR,fi,fi_FI,fil,fil_PH,fr,fr_BE,fr_CA,fr_CH,fr_FR,gl,gl_ES,gu,gu_IN,hi,hi_IN,hr,hr_HR,hu,hu_HU,hy,in,in_ID,is,it,it_CH,it_IT,iw,iw_IL,ja,ja_JP,ka,kk,km,kn,kn_IN,ko,ko_KR,ky,lo,lt,lt_LT,lv,lv_LV,mk,ml,ml_IN,mn,mr,mr_IN,ms,ms_MY,my,nb,nb_NO,ne,nl,nl_BE,nl_NL,pa,pl,pl_PL,pt,pt_BR,pt_PT,ro,ro_RO,ru,ru_RU,si,sk,sk_SK,sl,sl_SI,sq,sq_AL,sr,sr_Latn,sr_RS,sv,sv_SE,sw,sw_TZ,ta,ta_IN,te,te_IN,th,th_TH,tr,tr_TR,uk,uk_UA,ur,uz,vi,vi_VN,zh,zh_CN,zh_HK,zh_TW,zu,zu_ZA',
   'platforms': 'arm64-v8a,armeabi-v7a,armeabi',
   'roaming': 'mobile-notroaming', 'client': 'android-google',
   'sharedlibraries': 'android.ext.services,android.ext.shared,android.test.mock,android.test.runner,com.android.future.usb.accessory,com.android.ims.rcsmanager,com.android.location.provider,com.android.media.remotedisplay,com.android.mediadrm.signer,com.google.android.camera.experimental2016,com.google.android.dialer.support,com.google.android.gms,com.google.android.maps,com.google.android.media.effects,com.google.widevine.software.drm,com.qti.vzw.ims.internal,com.qualcomm.embmslibrary,com.qualcomm.qcrilhook,com.qualcomm.qti.QtiTelephonyServicelibrary,com.qualcomm.qti.rcsservice,com.verizon.embms,com.verizon.provider,com.vzw.apnlib,javax.obex,org.apache.http.legacy,org.lineageos.hardware,org.lineageos.platform',
   'simoperator': '311480', 'timezone': 'America/Los_Angeles',
   'userreadablename': 'Google Pixel (api27)', 'build.hardware': 'sailfish', 'build.radio': 'unknown', 'build.bootloader': '8996-012001-1711291800',
}

class DeviceBuilder(object):

    def __init__(self, device):
      self.device = {}
    
    def getBaseHeaders(self):
      return {
         "User-Agent": "Android-Finsky/{versionString} (api=3,versionCode={versionCode},sdk={sdk},device={device},hardware={hardware},product={product},platformVersionRelease={platform_v},model={model},buildId={build_id},isWideScreen=0,supportedAbis={supported_abis})".format(
            build_id=selfDevice.get('build.id'),
            device=selfDevice.get('build.device'),
            hardware=selfDevice.get('build.hardware'),
            model=selfDevice.get('build.model'),
            platform_v=selfDevice.get('build.version.release'),
            product=selfDevice.get('build.product'),
            sdk=selfDevice.get('build.version.sdk_int'),
            supported_abis=selfDevice.get('platforms').replace(',', ';'),
            versionCode=selfDevice.get('vending.version'),
            versionString = '8.4.19.V-all [0] [FP] 175058788',
         ),
      }

    def getDeviceUploadHeaders(self):
        headers = {}
        headers["X-DFE-Enabled-Experiments"] = "cl:billing.select_add_instrument_by_default"
        headers["X-DFE-Unsupported-Experiments"] = ("nocache:billing.use_charging_poller,"
            "market_emails,buyer_currency,prod_baseline,checkin.set_asset_paid_app_field,"
            "shekel_test,content_ratings,buyer_currency_in_app,nocache:encrypted_apk,recent_changes")
        headers["X-DFE-SmallestScreenWidthDp"] = "320"
        headers["X-DFE-Filter-Level"] = "3"
        return headers

    def getAuthHeaders(self, gsfid):
      headers = {
         "User-Agent": ("GoogleAuth/1.4 ({device} {id})").format(
            device=selfDevice.get('build.device'),
            id=selfDevice.get('build.id')
         )
      }
      if gsfid is not None:
            headers['device'] = "{0:x}".format(gsfid)
      return headers
    
    # NEED THIS
    def getAndroidCheckinRequest(self):
        request = googleplay_pb2.AndroidCheckinRequest()
        request.id = 0
        request.checkin.CopyFrom(self.getAndroidCheckin())
        #request.locale = self.locale
        request.version = 3
        request.deviceConfiguration.CopyFrom(self.getDeviceConfig())
        request.fragment = 0
        return request

    def getDeviceConfig(self):
        libList = selfDevice['sharedlibraries'].split(",")
        featureList = selfDevice['features'].split(",")
        localeList = selfDevice['locales'].split(",")
        glList = selfDevice['gl.extensions'].split(",")
        platforms = selfDevice['platforms'].split(",")
        hasFiveWayNavigation = (selfDevice['hasfivewaynavigation'] == 'true')
        hasHardKeyboard = (selfDevice['hashardkeyboard'] == 'true')
        deviceConfig = googleplay_pb2.DeviceConfigurationProto()
        deviceConfig.touchScreen = int(selfDevice['touchscreen'])
        deviceConfig.keyboard = int(selfDevice['keyboard'])
        deviceConfig.navigation = int(selfDevice['navigation'])
        deviceConfig.screenLayout = int(selfDevice['screenlayout'])
        deviceConfig.hasHardKeyboard = hasHardKeyboard
        deviceConfig.hasFiveWayNavigation = hasFiveWayNavigation
        deviceConfig.screenDensity = int(selfDevice['screen.density'])
        deviceConfig.screenWidth = int(selfDevice['screen.width'])
        deviceConfig.screenHeight = int(selfDevice['screen.height'])
        deviceConfig.glEsVersion = int(selfDevice['gl.version'])
        for x in platforms:
            deviceConfig.nativePlatform.append(x)
        for x in libList:
            deviceConfig.systemSharedLibrary.append(x)
        for x in featureList:
            deviceConfig.systemAvailableFeature.append(x)
        for x in localeList:
            deviceConfig.systemSupportedLocale.append(x)
        for x in glList:
            deviceConfig.glExtension.append(x)
        return deviceConfig

    def getAndroidBuild(self):
        androidBuild = googleplay_pb2.AndroidBuildProto()
        androidBuild.id = selfDevice['build.fingerprint']
        androidBuild.product = selfDevice['build.hardware']
        androidBuild.carrier = selfDevice['build.brand']
        androidBuild.radio = selfDevice['build.radio']
        androidBuild.bootloader = selfDevice['build.bootloader']
        androidBuild.device = selfDevice['build.device']
        androidBuild.sdkVersion = int(selfDevice['build.version.sdk_int'])
        androidBuild.model = selfDevice['build.model']
        androidBuild.manufacturer = selfDevice['build.manufacturer']
        androidBuild.buildProduct = selfDevice['build.product']
        androidBuild.client = selfDevice['client']
        androidBuild.otaInstalled = False
        androidBuild.timestamp = int(time()/1000)
        androidBuild.googleServices = int(selfDevice['gsf.version'])
        return androidBuild

    def getAndroidCheckin(self):
        androidCheckin = googleplay_pb2.AndroidCheckinProto()
        androidCheckin.build.CopyFrom(self.getAndroidBuild())
        androidCheckin.lastCheckinMsec = 0
        androidCheckin.cellOperator = selfDevice['celloperator']
        androidCheckin.simOperator = selfDevice['simoperator']
        androidCheckin.roaming = selfDevice['roaming']
        androidCheckin.userNumber = 0
        return androidCheckin
