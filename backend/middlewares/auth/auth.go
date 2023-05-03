package auth

import (
//    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    "openvpn-web-admin/utils/errutil"
)

func CheckSignIn(c *gin.Context) {
    if isSignIn, exist := c.Get("isSignIn"); !exist || !isSignIn.(bool) {
        errutil.AbortAndStatus(c, 401)
    }
}

func AddMeta(c *gin.Context) {
    session := sessions.Default(c)
    user := session.Get("user")
    if user != nil {
        c.Set("isSignIn", true)
    } else {
        c.Set("isSignIn", false)
    }
}
