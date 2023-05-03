package router
import (
    "github.com/gin-gonic/gin"
//    "net/http"

    "openvpn-web-admin/router/user"
    "openvpn-web-admin/middlewares/auth"
//    "openvpn-web-admin/utils/error"
)

var router *gin.RouterGroup

func Init(r *gin.RouterGroup) {
    router = r
    router.GET("/status", auth.CheckSignIn, status)
    user.Init(router.Group("/user"))
}

func status(c *gin.Context) {
//    panic("dead")
    /*error.AbortAndError(c, &error.Err{
        Code: 401,
        Msg: "test bad",
    })*/
    c.String(200, "test2")
}
