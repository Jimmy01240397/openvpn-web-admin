package user

import (
    "log"
    "fmt"
    "time"
    "strings"
    "encoding/json"

    "github.com/go-errors/errors"
    "golang.org/x/crypto/ssh/terminal"

    "openvpn-web-admin/utils/database"
    "openvpn-web-admin/utils/password"
    "openvpn-web-admin/models/group"
)

type User struct {
    Username string
    Password password.Password
    Online bool
    Enable bool
    Startdate *time.Time
    Enddate *time.Time
}

func (user *User) ToMap() map[string]any {
    var usermap map[string]any
    userjson, _ := json.Marshal(user)
    json.Unmarshal(userjson, &usermap)
    usermap["Password"] = user.Password.String()
    return usermap
}

func (user *User) GetGroups() []*group.Group {
    var resultgroups []*group.Group
    groups, err := group.GetGroups()
    if err != nil {
        log.Panicln(err)
        return nil
    }
    for _, nowgroup := range groups {
        if nowgroup.Contain(user.Username) {
            resultgroups = append(resultgroups, nowgroup)
        }
    }
    return resultgroups
}

func (user *User) IsAdmin() bool {
    admins, err := group.GetGroup("Administrators")
    if err != nil {
        log.Panicln(err)
        return false
    }
    return admins.Contain(user.Username)
}

const tablename = "user"

func init() {
    _, err := database.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            username varchar(255) NOT NULL PRIMARY KEY,
            password varchar(255),
            online tinyint(1) NOT NULL DEFAULT 0,
            enable tinyint(1) NOT NULL DEFAULT 1,
            startdate date DEFAULT NULL,
            enddate date DEFAULT NULL
        )
    `, tablename))
    if err != nil {
        log.Panicln(err)
        return
    }
    admin, err := GetUser("admin")
    if err != nil {
        log.Panicln(err)
        return
    }
    if admin == nil {
        var (
            adminpass any
            verifyadminpass any
        )
        fmt.Print("Set admin password: ")
        adminpass, err = terminal.ReadPassword(0)
        fmt.Println()
        if err != nil {
            log.Panicln(err)
            return
        }
        adminpass = string(adminpass.([]byte))
        fmt.Print("Verifying - Set admin password: ")
        verifyadminpass, err = terminal.ReadPassword(0)
        fmt.Println()
        if err != nil {
            log.Panicln(err)
            return
        }
        verifyadminpass = string(verifyadminpass.([]byte))
        if adminpass != verifyadminpass {
            log.Panicln("Verify failure")
            return
        }
        admin = &User{
            Username: "admin",
            Password: password.New(adminpass.(string)),
            Online: false,
            Enable: true,
        }
        AddUser(admin)
        admins, err := group.GetGroup("Administrators")
        if err != nil {
            log.Panicln(err)
            return
        }
        admins.AddMembers(admin.Username)
        //_, err = database.Exec("UPDATE groups SET members='[\"admin\"]' where groupname='Administrators'")
        err = group.UpdateGroup(admins.Groupname, admins, "Members")
        if err != nil {
            log.Panicln(err)
            return
        }
    }
}

func GetUser(username string) (*User, error) {
    rows, err := database.Query(fmt.Sprintf("SELECT * FROM %s WHERE username=? LIMIT 1", tablename), username)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    if !rows.Next() {
        return nil, nil
    }
    result := new(User)
    rows.Scan(&result.Username, &result.Password, &result.Online, &result.Enable, &result.Startdate, &result.Enddate)
    return result, nil
}

func GetUsers() ([]*User, error) {
    rows, err := database.Query(fmt.Sprintf("SELECT * FROM %s", tablename))
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var users []*User
    for rows.Next() {
        result := new(User)
        rows.Scan(&result.Username, &result.Password, &result.Online, &result.Enable, &result.Startdate, &result.Enddate)
        users = append(users, result)
    }
    return users, nil
}

func AddUser(user *User) error {
    testuser, err := GetUser(user.Username)
    if err != nil {
        return err
    }

    if testuser != nil {
        return errors.New("User is exist")
    }

    _, err = database.Exec(fmt.Sprintf("INSERT INTO %s (username, password, enable, startdate, enddate) VALUES (?,?,?,?,?)", tablename), user.Username, user.Password.String(), user.Enable, user.Startdate, user.Enddate)
    return err
}

func DeleteUser(username string) error {
    testuser, err := GetUser(username)
    if err != nil {
        return err
    }

    if testuser == nil {
        return errors.New("User not exist")
    }
    
    _, err = database.Exec(fmt.Sprintf("DELETE FROM %s where username=?", tablename), username)
    return err
}

func UpdateUser(username string, user *User, fields ...string) error {
    testuser, err := GetUser(username)
    if err != nil {
        return err
    }

    if testuser == nil {
        return errors.New("User not exist")
    }

    usermap := user.ToMap()

    data := []any{}
    for index := 0; index < len(fields); index++ {
        if fields[index] == "Username" {
            fields = append(fields[:index], fields[index+1:]...)
            index--
        } else {
            data = append(data, usermap[fields[index]])
        }
    }
    data = append(data, username)

    if len(fields) <= 0 {
        return nil
    }
    _, err = database.Exec(fmt.Sprintf("UPDATE %s SET %s where username=?", tablename, strings.ToLower(strings.Join(fields[:], "=?, ")) + "=?"), data...)
    return err
}
