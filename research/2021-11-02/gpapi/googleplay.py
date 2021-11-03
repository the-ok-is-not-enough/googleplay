from . import googleplay_pb2, config, utils
from base64 import b64decode, urlsafe_b64encode
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives.asymmetric.utils import encode_dss_signature
from cryptography.hazmat.primitives.serialization import load_der_public_key
from datetime import datetime
import requests

BASE = "https://android.clients.google.com/"
FDFE = BASE + "fdfe/"

ACCEPT_TOS_URL = FDFE + "acceptTos"
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
LOG_URL = FDFE + "log"
PURCHASE_URL = FDFE + "purchase"
REVIEWS_URL = FDFE + "rev"
SEARCH_URL = FDFE + "search"
TOC_URL = FDFE + "toc"
UPLOAD_URL = FDFE + "uploadDeviceConfig"
ssl_verify = True

class LoginError(Exception):
    def __init__(self, value):
        self.value = value

    def __str__(self):
        return repr(self.value)

class RequestError(Exception):
    def __init__(self, value):
        self.value = value

    def __str__(self):
        return repr(self.value)

class SecurityCheckError(Exception):
    def __init__(self, value):
        self.value = value

    def __str__(self):
        return repr(self.value)


class GooglePlayAPI(object):
    """Google Play Unofficial API Class

    Usual APIs methods are login(), search(), details(), bulkDetails(),
    download(), browse(), reviews() and list()."""

    def __init__(self, locale="en_US", timezone="UTC", device_codename="bacon",
                 proxies_config=None):
        self.authSubToken = None
        self.gsfId = None
        self.device_config_token = None
        self.deviceCheckinConsistencyToken = None
        self.dfeCookie = None
        self.proxies_config = proxies_config
        self.deviceBuilder = config.DeviceBuilder(device_codename)
        self.setLocale(locale)
        self.setTimezone(timezone)

    def setLocale(self, locale):
        self.deviceBuilder.setLocale(locale)

    def setTimezone(self, timezone):
        self.deviceBuilder.setTimezone(timezone)

    def setAuthSubToken(self, authSubToken):
        self.authSubToken = authSubToken

    def getHeaders(self, upload_fields=False):
        """Return the default set of request headers, which
        can later be expanded, based on the request type"""

        if upload_fields:
            headers = self.deviceBuilder.getDeviceUploadHeaders()
        else:
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
        """Upload the device configuration of the fake device
        selected in the __init__ methodi to the google account."""

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
        requestParams = self.deviceBuilder.getLoginParams(email, passwd)
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

    def _deliver_data(self, url, cookies):
        headers = self.getHeaders()
        response = requests.get(url, headers=headers,
                                cookies=cookies, verify=ssl_verify,
                                stream=True, timeout=60,
                                proxies=self.proxies_config)
        print(response.status_code, response.url)
        total_size = response.headers.get('content-length')
        chunk_size = 32 * (1 << 10)
        return {'data': response.iter_content(chunk_size=chunk_size),
                'total_size': total_size,
                'chunk_size': chunk_size}

    def log(self, docid):
        log_request = googleplay_pb2.LogRequest()
        log_request.downloadConfirmationQuery = "confirmFreeDownload?doc=" + docid
        timestamp = int(datetime.now().timestamp())
        log_request.timestamp = timestamp

        string_request = log_request.SerializeToString()
        response = requests.post(LOG_URL,
                                 data=string_request,
                                 headers=self.getHeaders(),
                                 verify=ssl_verify,
                                 timeout=60,
                                 proxies=self.proxies_config)
        print(response.status_code, response.url)
        response = googleplay_pb2.ResponseWrapper.FromString(response.content)
        if response.commands.displayErrorMessage != "":
            raise RequestError(response.commands.displayErrorMessage)

    def toc(self):
        response = requests.get(TOC_URL,
                               headers=self.getHeaders(),
                               verify=ssl_verify,
                               timeout=60,
                               proxies=self.proxies_config)
        print(response.status_code, response.url)
        data = googleplay_pb2.ResponseWrapper.FromString(response.content)
        tocResponse = data.payload.tocResponse
        if utils.hasTosContent(tocResponse) and utils.hasTosToken(tocResponse):
            self.acceptTos(tocResponse.tosToken)
        if utils.hasCookie(tocResponse):
            self.dfeCookie = tocResponse.cookie
        return utils.parseProtobufObj(tocResponse)


    def acceptTos(self, tosToken):
        params = {
            "tost": tosToken,
            "toscme": "false"
        }
        response = requests.get(ACCEPT_TOS_URL,
                               headers=self.getHeaders(),
                               params=params,
                               verify=ssl_verify,
                               timeout=60,
                               proxies=self.proxies_config)
        print(response.status_code, response.url)
        data = googleplay_pb2.ResponseWrapper.FromString(response.content)
        return utils.parseProtobufObj(data.payload.acceptTosResponse)

    @staticmethod
    def getDevicesCodenames():
        return config.getDevicesCodenames()

    @staticmethod
    def getDevicesReadableNames():
        return config.getDevicesReadableNames()
