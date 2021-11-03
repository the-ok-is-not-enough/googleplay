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

    def getAuthSubToken(self, email, passwd, master_token):
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
      second_round_token = self.getSecondRoundToken(master_token, requestParams)
      print(second_round_token)
      self.setAuthSubToken(second_round_token)

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
   'gsf.version': '12685052', 'vending.version': '81031200', 'vending.versionstring': '10.3.12-all [0] [PR] 198814133', 'celloperator': '311480',
   'hasfivewaynavigation': 'false', 'gl.version': '196610', 'screen.density': '420', 'screen.width': '1794', 'screen.height': '1080',
   'platforms': 'arm64-v8a,armeabi-v7a,armeabi',
   'roaming': 'mobile-notroaming', 'client': 'android-google',
   'simoperator': '311480', 'timezone': 'America/Los_Angeles',
   'userreadablename': 'Google Pixel (api27)', 'build.hardware': 'sailfish', 'build.radio': 'unknown', 'build.bootloader': '8996-012001-1711291800',
   # NEED THIS
   'features': ','.join([
      'android.hardware.audio.low_latency,android.hardware.audio.output,android.hardware.audio.pro,android.hardware.bluetooth,android.hardware.bluetooth_le',
      'android.hardware.camera,android.hardware.camera.any,android.hardware.camera.autofocus,android.hardware.camera.capability.manual_post_processing',
      'android.hardware.camera.capability.manual_sensor,android.hardware.camera.capability.raw,android.hardware.camera.flash,android.hardware.camera.front',
      'android.hardware.camera.level.full,android.hardware.faketouch,android.hardware.fingerprint,android.hardware.location,android.hardware.location.gps',
      'android.hardware.location.network,android.hardware.microphone,android.hardware.nfc,android.hardware.nfc.any,android.hardware.nfc.hce,android.hardware.nfc.hcef',
      'android.hardware.opengles.aep,android.hardware.ram.normal,android.hardware.screen.landscape,android.hardware.screen.portrait',
      'android.hardware.sensor.accelerometer,android.hardware.sensor.barometer,android.hardware.sensor.compass,android.hardware.sensor.gyroscope',
      'android.hardware.sensor.hifi_sensors,android.hardware.sensor.light,android.hardware.sensor.proximity,android.hardware.sensor.stepcounter',
      'android.hardware.sensor.stepdetector,android.hardware.telephony,android.hardware.telephony.carrierlock,android.hardware.telephony.cdma',
      'android.hardware.telephony.gsm,android.hardware.touchscreen,android.hardware.touchscreen.multitouch,android.hardware.touchscreen.multitouch.distinct',
      'android.hardware.touchscreen.multitouch.jazzhand,android.hardware.usb.accessory,android.hardware.usb.host,android.hardware.vr.headtracking',
      'android.hardware.vr.high_performance,android.hardware.vulkan.compute,android.hardware.vulkan.level,android.hardware.vulkan.version,android.hardware.wifi',
      'android.hardware.wifi.direct,android.hardware.wifi.passpoint,android.software.activities_on_secondary_displays,android.software.app_widgets',
      'android.software.autofill,android.software.backup,android.software.companion_device_setup,android.software.connectionservice,android.software.cts',
      'android.software.device_admin,android.software.file_based_encryption,android.software.home_screen,android.software.input_methods',
      'android.software.live_wallpaper,android.software.managed_users,android.software.midi,android.software.picture_in_picture,android.software.print',
      'android.software.securely_removes_users,android.software.sip,android.software.sip.voip,android.software.voice_recognizers,android.software.vr.mode',
      'android.software.webview,com.google.android.apps.dialer.SUPPORTED,com.google.android.apps.photos.NEXUS_PRELOAD,com.google.android.apps.photos.nexus_preload',
      'com.google.android.feature.EXCHANGE_6_2,com.google.android.feature.GOOGLE_BUILD,com.google.android.feature.GOOGLE_EXPERIENCE',
      'com.google.android.feature.PIXEL_EXPERIENCE,com.nxp.mifare,com.verizon.hardware.telephony.ehrpd,com.verizon.hardware.telephony.lte,org.lineageos.android',
      'org.lineageos.audio,org.lineageos.hardware,org.lineageos.livedisplay,org.lineageos.performance,org.lineageos.profiles,org.lineageos.style',
      'org.lineageos.weather,projekt.substratum.theme',
   ])
}

class DeviceBuilder(object):

    def __init__(self, device):
      self.device = {}
    
    def getBaseHeaders(self):
      agent = "Android-Finsky/8.4.19.V-all [0] [FP] 175058788 (api=3,"
      agent += 'versionCode=' + selfDevice.get('vending.version') + ','
      agent += 'sdk=' + selfDevice.get('build.version.sdk_int') + ','
      agent += 'device=' + selfDevice.get('build.device') + ','
      agent += 'hardware=' + selfDevice.get('build.hardware') + ','
      agent += 'product=' + selfDevice.get('build.product') + ','
      agent += 'platformVersionRelease=' + selfDevice.get('build.version.release') + ','
      agent += 'model=' + selfDevice.get('build.model') + ','
      agent += 'buildId=' + selfDevice.get('build.id') + ','
      agent += 'isWideScreen=0,'
      agent += 'supportedAbis=' + selfDevice.get('platforms').replace(',', ';') + ')'
      return {'User-Agent': agent}

    def getDeviceUploadHeaders(self):
      headers = {}
      headers["X-DFE-Enabled-Experiments"] = "cl:billing.select_add_instrument_by_default"
      headers["X-DFE-Unsupported-Experiments"] = (
         "nocache:billing.use_charging_poller,"
         "market_emails,buyer_currency,prod_baseline,checkin.set_asset_paid_app_field,"
         "shekel_test,content_ratings,buyer_currency_in_app,nocache:encrypted_apk,recent_changes"
      )
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
        request.version = 3
        request.deviceConfiguration.CopyFrom(self.getDeviceConfig())
        request.fragment = 0
        return request

    def getDeviceConfig(self):
      featureList = selfDevice['features'].split(",")
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
      for x in featureList:
          deviceConfig.systemAvailableFeature.append(x)
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
