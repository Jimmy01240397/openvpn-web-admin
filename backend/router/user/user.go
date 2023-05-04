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
    router.GET("/get", auth.CheckSignIn, userinfo)
    router.GET("/getall", auth.CheckSignIn, auth.CheckIsAdmin, alluserinfo)
    router.POST("/login", login)
    router.POST("/add", auth.CheckSignIn, auth.CheckIsAdmin, adduser)
    router.GET("/logout", logout)
    router.GET("/issignin", issignin)
}

func issignin(c *gin.Context) {
    isSignIn, exist := c.Get("isSignIn")
    c.JSON(200, exist && isSignIn.(bool))
}

func alluserinfo(c *gin.Context) {
    userdatas, err := user.GetUsers()
    if err != nil {
        log.Panicln(err)
        return
    }
    c.JSON(200, userdatas)
}

func userinfo(c *gin.Context) {
    session := sessions.Default(c)
    username := session.Get("user").(string)
    userdata, err := user.GetUser(username)
    if err != nil {
        log.Panicln(err)
        return
    }
    if userdata == nil {
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

func adduser(c *gin.Context) {
    postdata := make(map[string]any)
    c.BindJSON(&postdata)
    if postdata["username"].(string) == "" || postdata["password"].(string) == "" {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 403,
            Msg: "username or password can't be empty",
        })
        return
    }
    newuser := &user.User{
        Username: postdata["username"].(string),
        Password: password.New(postdata["password"].(string)),
        Online: false,
        Enable: true,
    }
    err := user.AddUser(newuser)
    if err != nil {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 403,
            Msg: err.Error(),
        })
        return
    }
    c.String(200, "Add user success")
}
