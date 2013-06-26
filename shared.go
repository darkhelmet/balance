package main

import (
    "net/http"
    "strings"
)

const (
    colon   = ":"
    XRealIP = "X-Real-IP"
)

func RealIP(req *http.Request) string {
    parts := strings.Split(req.RemoteAddr, colon)
    if len(parts) == 0 {
        return ""
    }
    return parts[0]
}
