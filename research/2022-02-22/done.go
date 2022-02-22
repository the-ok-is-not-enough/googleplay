package googleplay

import (
   "github.com/89z/format"
   "net/http"
)

type Token struct {
   Token string
}

type Device struct {
   AndroidID uint64
}

type Header struct {
   http.Header
}

var LogLevel format.LogLevel
