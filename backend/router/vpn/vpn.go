package vpn

import (
    "log"

    "github.com/gin-gonic/gin"
//    "github.com/gin-contrib/sessions"

    "openvpn-web-admin/middlewares/auth"
    "openvpn-web-admin/models/vpn"
    "openvpn-web-admin/utils/errutil"
//    "net/http"
)

var router *gin.RouterGroup

func Init(r *gin.RouterGroup) {
    router = r
    router.POST("/get", auth.CheckSignIn, getvpn)
    router.GET("/getall", auth.CheckSignIn, auth.CheckIsAdmin, getvpns)
//    router.POST("/login", login)
    router.POST("/add", auth.CheckSignIn, auth.CheckIsAdmin, addvpn)
    router.POST("/delete", auth.CheckSignIn, auth.CheckIsAdmin, deletevpn)
    router.POST("/startstop", auth.CheckSignIn, auth.CheckIsAdmin, startstopvpn)
//    router.GET("/logout", logout)
//    router.GET("/issignin", issignin)
}

func getvpn(c *gin.Context) {
    postdata := make(map[string]any)
    c.BindJSON(&postdata)
    vpndata, err := vpn.GetVPN(postdata["vpnname"].(string))
    if err != nil {
        log.Panicln(err)
        return
    }
    if vpndata == nil {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 403,
            Msg: "vpn not exist",
        })
        return
    }
    c.JSON(200, vpndata)
}

func getvpns(c *gin.Context) {
    vpndatas, err := vpn.GetVPNs()
    if err != nil {
        log.Panicln(err)
        return
    }
    c.JSON(200, vpndatas)
}

func addvpn(c *gin.Context) {
    postdata := make(map[string]any)
    c.BindJSON(&postdata)
    if postdata["vpnname"].(string) == "" {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 403,
            Msg: "vpnname can't be empty",
        })
        return
    }
    newvpn := &vpn.VPN{
        VPNname: postdata["vpnname"].(string),
        Enable: postdata["enable"].(bool),
        AllowUsers: []string{},
        AllowGroups: []string{},
    }
    err := vpn.AddVPN(newvpn)
    if err != nil {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 403,
            Msg: err.Error(),
        })
        return
    }
    c.String(200, "Add vpn success")
}

func deletevpn(c *gin.Context) {
    postdata := make(map[string]any)
    c.BindJSON(&postdata)
    if postdata["vpnname"].(string) == "" {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 403,
            Msg: "vpnname can't be empty",
        })
        return
    }
    err := vpn.DeleteVPN(postdata["vpnname"].(string))
    if err != nil {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 403,
            Msg: err.Error(),
        })
        return
    }
    c.String(200, "Delete vpn success")
}

func startstopvpn(c *gin.Context) {
    postdata := make(map[string]any)
    c.BindJSON(&postdata)
    if postdata["vpnname"].(string) == "" {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 403,
            Msg: "vpnname can't be empty",
        })
        return
    }
    vpndata, err := vpn.GetVPN(postdata["vpnname"].(string))
    if err != nil {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 403,
            Msg: err.Error(),
        })
        return
    }
    vpndata.Active = postdata["active"].(bool)
    err = vpn.UpdateVPN(postdata["vpnname"].(string), vpndata, "Active")
    if err != nil {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 403,
            Msg: err.Error(),
        })
        return
    }
    c.String(200, "Start or Stop vpn success")
}
