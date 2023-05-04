package group

import (
    "log"

    "github.com/gin-gonic/gin"

    "openvpn-web-admin/middlewares/auth"
    "openvpn-web-admin/models/group"
    "openvpn-web-admin/utils/errutil"
//    "net/http"
)

var router *gin.RouterGroup

func Init(r *gin.RouterGroup) {
    router = r
//    router.GET("/", auth.CheckSignIn, userinfo)
    router.POST("/get", auth.CheckSignIn, getgroup)
    router.GET("/getall", auth.CheckSignIn, getgroups)
//    router.POST("/adduser", adduser)
}

func getgroup(c *gin.Context) {
    postdata := make(map[string]any)
    c.BindJSON(&postdata)
    groupdata, err := group.GetGroup(postdata["groupname"].(string))
    if err != nil {
        log.Panicln(err)
        return
    }
    if groupdata == nil {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 403,
            Msg: "group not exist",
        })
        return
    }
    c.JSON(200, groupdata)
}

func getgroups(c *gin.Context) {
    groupdatas, err := group.GetGroups()
    if err != nil {
        log.Panicln(err)
        return
    }
    c.JSON(200, groupdatas)
}
