package group

import (
    "log"
    "fmt"
    "strings"
    "encoding/json"

    "github.com/go-errors/errors"

    "openvpn-web-admin/utils/database"
)

type Group struct {
    Groupname string
    Members []string
}

func (group *Group) ToMap() map[string]any {
    var groupmap map[string]any
    groupjson, _ := json.Marshal(group)
    json.Unmarshal(groupjson, &groupmap)
    members, _ := json.Marshal(group.Members)
    groupmap["Members"] = string(members)
    return groupmap
}

func (group *Group) AddMembers(members ...string) {
    var appendmembers []string
    for _, member := range members {
        if group.Contain(member) {
            continue
        }
        appendmembers = append(appendmembers, member)
    }
    group.Members = append(group.Members, appendmembers...)
}

func (group *Group) Contain(member string) bool {
    for _, groupmember := range group.Members {
        if groupmember == member {
            return true
        }
    }
    return false
}

const tablename = "groups"

func init() {
    _, err := database.Exec(fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            groupname varchar(255) NOT NULL PRIMARY KEY,
            members string NOT NULL DEFAULT "[]"
        )
    `, tablename))
    if err != nil {
        log.Panicln(err)
    }
    admins, err := GetGroup("Administrators")
    if err != nil {
        log.Panicln(err)
    }
    if admins == nil {
        AddGroup(&Group{
            Groupname: "Administrators",
            Members: []string{},
        })
    }
}

func GetGroup(groupname string) (*Group, error) {
    rows, err := database.Query(fmt.Sprintf("SELECT * FROM %s WHERE groupname=? LIMIT 1", tablename), groupname)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    if !rows.Next() {
        return nil, nil
    }
    result := new(Group);
    var membersjson string
    rows.Scan(&result.Groupname, &membersjson)
    json.Unmarshal([]byte(membersjson), &result.Members)
    return result, nil
}

func GetGroups() ([]*Group, error) {
    rows, err := database.Query(fmt.Sprintf("SELECT * FROM %s", tablename))
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var groups []*Group
    for rows.Next() {
        result := new(Group);
        var membersjson string
        rows.Scan(&result.Groupname, &membersjson)
        json.Unmarshal([]byte(membersjson), &result.Members)
        groups = append(groups, result)
    }
    return groups, nil
}

func AddGroup(group *Group) error {
    testgroup, err := GetGroup(group.Groupname)
    if err != nil {
        return err
    }

    if testgroup != nil {
        return errors.New("Group is exist")
    }

    members, _ := json.Marshal(group.Members)
    _, err = database.Exec(fmt.Sprintf("INSERT INTO %s (groupname, members) VALUES (?,?)", tablename), group.Groupname, string(members))
    return err
}

func DeleteGroup(groupname string) error {
    testgroup, err := GetGroup(groupname)
    if err != nil {
        return err
    }

    if testgroup == nil {
        return errors.New("Group not exist")
    }

    _, err = database.Exec(fmt.Sprintf("DELETE FROM %s where groupname=?", tablename), groupname)
    return err
}

func UpdateGroup(groupname string, group *Group, fields ...string) error {
    testgroup, err := GetGroup(groupname)
    if err != nil {
        return err
    }

    if testgroup == nil {
        return errors.New("Group not exist")
    }

    groupmap := group.ToMap()

    data := []any{}
    for index := 0; index < len(fields); index++ {
        if fields[index] == "Groupname" {
            fields = append(fields[:index], fields[index+1:]...)
            index--
        } else {
            data = append(data, groupmap[fields[index]])
        }
    }
    data = append(data, groupname)

    if len(fields) <= 0 {
        return nil
    }
    _, err = database.Exec(fmt.Sprintf("UPDATE %s SET %s where groupname=?", tablename, strings.ToLower(strings.Join(fields[:], "=?, ")) + "=?"), data...)
    return err
}
