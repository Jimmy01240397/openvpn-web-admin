package vpn

import (
    "log"
    "fmt"
    "strings"
    "sync"
    "os"
    "os/exec"
    "io/ioutil"
    "encoding/json"
    "path/filepath"

    "github.com/go-errors/errors"

    "openvpn-web-admin/utils/database"
)

type VPN struct {
    VPNname string
    Enable bool
    Active bool
    AllowUsers []string
    AllowGroups []string
}

var lock *sync.RWMutex

func (vpn *VPN) ToMap() map[string]any {
    var vpnmap map[string]any
    vpnjson, _ := json.Marshal(vpn)
    json.Unmarshal(vpnjson, &vpnmap)
    allowuser, _ := json.Marshal(vpn.AllowUsers)
    vpnmap["AllowUsers"] = string(allowuser)
    allowgroup, _ := json.Marshal(vpn.AllowGroups)
    vpnmap["AllowGroups"] = string(allowgroup)
    return vpnmap
}

func (vpn *VPN) AddUsers(users ...string) {
    var appendusers []string
    for _, user := range users {
        if vpn.ContainUser(user) {
            continue
        }
        appendusers = append(appendusers, user)
    }
    vpn.AllowUsers = append(vpn.AllowUsers, appendusers...)
}

func (vpn *VPN) AddGroups(groups ...string) {
    var appendgroups []string
    for _, group := range groups {
        if vpn.ContainGroup(group) {
            continue
        }
        appendgroups = append(appendgroups, group)
    }
    vpn.AllowGroups = append(vpn.AllowGroups, appendgroups...)
}

func (vpn *VPN) ContainUser(username string) bool {
    for _, allowuser := range vpn.AllowUsers {
        if allowuser == username {
            return true
        }
    }
    return false
}

func (vpn *VPN) ContainGroup(groupname string) bool {
    for _, allowgroup := range vpn.AllowGroups {
        if allowgroup == groupname {
            return true
        }
    }
    return false
}

func (vpn *VPN) execAdd() {
    lock.Lock()
    defer lock.Unlock()
    os.MkdirAll("server", 0755)
    if _, err := os.Stat(filepath.Join("server", vpn.VPNname + ".conf")); err != nil {
        if _, err := os.Stat(vpn.VPNname + ".conf"); err == nil {
            os.Rename(vpn.VPNname + ".conf", filepath.Join("server", vpn.VPNname + ".conf"))
        } else {
            bytesRead, _ := ioutil.ReadFile("/usr/share/doc/openvpn/examples/sample-config-files/server.conf")
            ioutil.WriteFile(filepath.Join("server", vpn.VPNname + ".conf"), bytesRead, 0644)
        }
        if vpn.Enable {
            lock.Unlock()
            vpn.execEnable()
            lock.Lock()
        }
    } else {
        os.Remove(vpn.VPNname + ".conf")
    }
}

func (vpn *VPN) execEnable() {
    if vpn.Enable {
        vpn.execAdd()
        lock.Lock()
        defer lock.Unlock()
        os.Symlink(filepath.Join("server", vpn.VPNname + ".conf"), vpn.VPNname + ".conf")
        exec.Command("sudo", "systemctl", "enable", "openvpn@" + vpn.VPNname + ".service").Run()
    } else {
        lock.Lock()
        defer lock.Unlock()
        exec.Command("sudo", "systemctl", "stop", "openvpn@" + vpn.VPNname + ".service").Run()
        exec.Command("sudo", "systemctl", "disable", "openvpn@" + vpn.VPNname + ".service").Run()
        os.Remove(vpn.VPNname + ".conf")
    }
}

const tablename = "vpn"

func init() {
    lock = new(sync.RWMutex)
    _, err := database.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            vpnname varchar(255) NOT NULL PRIMARY KEY,
            enable tinyint(1) NOT NULL DEFAULT 1,
            allowusers string NOT NULL DEFAULT "[]",
            allowgroups string NOT NULL DEFAULT "[]"
        )
    `, tablename))
    if err != nil {
        log.Panicln(err)
    }
}

func GetVPN(vpnname string) (*VPN, error) {
    rows, err := database.Query(fmt.Sprintf("SELECT * FROM %s WHERE vpnname=? LIMIT 1", tablename), vpnname)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    if !rows.Next() {
        return nil, nil
    }
    result := new(VPN);
    var allowusersjson string
    var allowgroupsjson string
    rows.Scan(&result.VPNname, &result.Enable, &allowusersjson, &allowgroupsjson)
    json.Unmarshal([]byte(allowusersjson), &result.AllowUsers)
    json.Unmarshal([]byte(allowgroupsjson), &result.AllowGroups)
    lock.RLock()
    defer lock.RUnlock()
    active, _ := exec.Command("sudo", "systemctl", "is-active", "openvpn@" + vpnname + ".service").Output()
    result.Active = (strings.TrimSpace(string(active)) == "active")
    return result, nil
}

func GetVPNs() ([]*VPN, error) {
    rows, err := database.Query(fmt.Sprintf("SELECT * FROM %s", tablename))
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var vpns []*VPN
    lock.RLock()
    defer lock.RUnlock()
    for rows.Next() {
        result := new(VPN);
        var allowusersjson string
        var allowgroupsjson string
        rows.Scan(&result.VPNname, &result.Enable, &allowusersjson, &allowgroupsjson)
        json.Unmarshal([]byte(allowusersjson), &result.AllowUsers)
        json.Unmarshal([]byte(allowgroupsjson), &result.AllowGroups)
        active, _ := exec.Command("sudo", "systemctl", "is-active", "openvpn@" + result.VPNname + ".service").Output()
        result.Active = (strings.TrimSpace(string(active)) == "active")
        vpns = append(vpns, result)
    }
    return vpns, nil
}

func AddVPN(vpn *VPN) error {
    testvpn, err := GetVPN(vpn.VPNname)
    if err != nil {
        return err
    }

    if testvpn != nil {
        return errors.New("VPN is exist")
    }

    vpn.execAdd()

    allowusers, _ := json.Marshal(vpn.AllowUsers)
    allowgroups, _ := json.Marshal(vpn.AllowGroups)
    _, err = database.Exec(fmt.Sprintf("INSERT INTO %s (vpnname, enable, allowusers, allowgroups) VALUES (?,?,?,?)", tablename), vpn.VPNname, vpn.Enable, string(allowusers), string(allowgroups))
    return err
}

func DeleteVPN(vpnname string) error {
    vpn, err := GetVPN(vpnname)
    if err != nil {
        return err
    }

    if vpn == nil {
        return errors.New("VPN not exist")
    }

    vpn.Enable = false
    vpn.execEnable()
    lock.Lock()
    defer lock.Unlock()
    os.Remove(filepath.Join("server", vpn.VPNname + ".conf"))
    exec.Command("sudo", "systemctl", "daemon-reload").Run()
    exec.Command("sudo", "systemctl", "reset-failed").Run()

    _, err = database.Exec(fmt.Sprintf("DELETE FROM %s where vpnname=?", tablename), vpnname)
    return err
}

func UpdateVPN(vpnname string, vpn *VPN, fields ...string) error {
    testvpn, err := GetVPN(vpnname)
    if err != nil {
        return err
    }

    if testvpn == nil {
        return errors.New("VPN not exist")
    }

    vpnmap := vpn.ToMap()

    data := []any{}
    for index := 0; index < len(fields); index++ {
        if fields[index] == "VPNname" {
            fields = append(fields[:index], fields[index+1:]...)
            index--
        } else if fields[index] == "Active" {
            lock.Lock()
            if vpn.Active {
                exec.Command("sudo", "systemctl", "start", "openvpn@" + vpn.VPNname + ".service").Run()
            } else {
                exec.Command("sudo", "systemctl", "stop", "openvpn@" + vpn.VPNname + ".service").Run()
            }
            lock.Unlock()
            fields = append(fields[:index], fields[index+1:]...)
            index--
        } else {
            data = append(data, vpnmap[fields[index]])
            if fields[index] == "Enable" {
                vpn.execEnable()
            }
        }
    }
    data = append(data, vpnname)

    if len(fields) <= 0 {
        return nil
    }
    _, err = database.Exec(fmt.Sprintf("UPDATE %s SET %s where vpnname=?", tablename, strings.ToLower(strings.Join(fields[:], "=?, ")) + "=?"), data...)
    return err
}

