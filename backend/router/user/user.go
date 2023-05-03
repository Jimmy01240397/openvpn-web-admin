package user

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    "openvpn-web-admin/middlewares/auth"
    "openvpn-web-admin/models/user"
    "openvpn-web-admin/utils/errutil"
    "openvpn-web-admin/utils/password"
//    "net/http"
)

var router *gin.RouterGroup

func Init(r *gin.RouterGroup) {
    router = r
    router.GET("/", auth.CheckSignIn, userinfo)
    router.POST("/login", login)
    router.GET("/logout", logout)
    router.GET("/issignin", issignin)
}

func issignin(c *gin.Context) {
    isSignIn, exist := c.Get("isSignIn")
    c.JSON(200, exist && isSignIn.(bool))
}

func userinfo(c *gin.Context) {
    session := sessions.Default(c)
    username := session.Get("user").(string)
    userdata, err := user.GetUser(username)
    if err != nil {
        log.Panicln(err)
        return
    }
    if err != nil {
        errutil.AbortAndStatus(c, 401)
        return
    }
    userdata.Password = password.Password{}
    c.JSON(200, userdata)
}

func logout(c *gin.Context) {
    session := sessions.Default(c)
    session.Clear()
    session.Save()
    c.JSON(200, true)
}

func login(c *gin.Context) {
    postdata := make(map[string]any)
    c.BindJSON(&postdata)
    userdata, err := user.GetUser(postdata["username"].(string))
    if err != nil {
        log.Panicln(err)
        return
    }
    nowpass := password.New(postdata["password"].(string))
    if userdata == nil || userdata.Password != nowpass {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 401,
            Msg: "username or password incorrect",
        })
        return
    }
    session := sessions.Default(c)
    session.Set("user", postdata["username"])
    session.Save()
    c.String(200, "Login success")
}
