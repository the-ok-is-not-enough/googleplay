package research

import (
   "net/http"
   "net/url"
)

func auth() (*http.Response, error) {
   req := new(http.Request)
   req.Header = make(http.Header)
   req.Method = "POST"
   req.URL = new(url.URL)
   req.URL.Host = "android.clients.google.com"
   req.URL.Path = "/auth"
   req.URL.Scheme = "https"
   val := make(url.Values)
   val["Token"] = []string{token}
   val["service"] = []string{"oauth2:https://www.googleapis.com/auth/googleplay"}
   // Error=BAD_REQUEST
   val["client_sig"] = []string{"38918a453d07199354f8b19af05ec6562ced5788"}
   // Error=ServerError
   val["app"] = []string{"com.google.android.gms"}
   req.URL.RawQuery = val.Encode()
   return new(http.Transport).RoundTrip(req)
}
