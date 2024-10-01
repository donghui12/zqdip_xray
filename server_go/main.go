package main

import (
    "encoding/json"
    "fmt"
    "time"
    "net/http"

    "golang.org/x/crypto/ssh"
    "github.com/gin-gonic/gin"
)

type HostInfo struct {
    IP   string `json:"ip"`
    User string `json:"user"`
    Pwd  string `json:"pwd"`
    Port string `json:"port"`
}

type Request struct {
    Hosts []HostInfo `json:"hosts"`
}

type Result struct {
    IP     string                 `json:"ip"`
    User   string                 `json:"user"`
    Pass   string                 `json:"pass"`
    Port   int                    `json:"port"`
    Status string                 `json:"status"`
}

type Inbound struct {
    Listen        string `json:"listen"`
    Port          int    `json:"port"`
    Protocol      string `json:"protocol"`
    Settings      struct {
        Auth     string `json:"auth"`
        Accounts []struct {
            User string `json:"user"`
            Pass string `json:"pass"`
        } `json:"accounts"`
        UDP bool   `json:"udp"`
        IP  string `json:"ip"`
    } `json:"settings"`
}

func main() {
    r := gin.Default()

    // 定义批量执行任务的接口
    r.POST("/batch_execute", func(c *gin.Context) {
        var req Request
        if err := c.BindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
            return
        }

        results := []Result{}

        for _, host := range req.Hosts {
            user, pass, port, err := executeCommandsAndGetConfig(host.IP, host.User, host.Pwd, host.Port)
            status := "Success"
            if err != nil {
                status = fmt.Sprintf("Failed: %v", err)
            }

            results = append(results, Result{
                IP:     host.IP,
                User:   user,
                Pass:   pass,
                Port:   port,
                Status: status,
            })
        }

        c.JSON(http.StatusOK, gin.H{
            "hosts": results,
        })
    })

    // 启动服务
    r.Run(":8080")
}

// 执行命令并获取 inbounds 中的相关字段
func executeCommandsAndGetConfig(ip, user, pwd, port string) (string, string, int, error) {
    config := &ssh.ClientConfig{
        User: user,
        Auth: []ssh.AuthMethod{
            ssh.Password(pwd),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout:         10 * time.Second,
    }

    addr := fmt.Sprintf("%s:22", ip)
    client, err := ssh.Dial("tcp", addr, config)
    if err != nil {
        return "", "", 0, fmt.Errorf("failed to dial: %v", err)
    }
    defer client.Close()

    // 创建一个新的SSH会话
    session, err := client.NewSession()
    if err != nil {
        return "", "", 0, fmt.Errorf("failed to create session: %v", err)
    }
    defer session.Close()

    // 执行下载、解压和运行脚本的命令
    // cmd := `curl -O http://82.157.189.81/install/zqdip_xray.zip && unzip -q zqdip_xray.zip && cd zqdip_xray && bash prepare.sh run`
    cmd := `ls`
    if err := session.Run(cmd); err != nil {
        return "", "", 0, fmt.Errorf("failed to run setup command: %v", err)
    }

    // 第二步：读取 config.json 文件
    session, err = client.NewSession()
    if err != nil {
        return "", "", 0, fmt.Errorf("failed to create new session for config: %v", err)
    }
    defer session.Close()

    cmd = "cat /usr/local/etc/xray/config.json"
    output, err := session.CombinedOutput(cmd)
    if err != nil {
        return "", "", 0, fmt.Errorf("failed to read config.json: %v", err)
    }

    // 解析 config.json 文件
    var configData map[string]interface{}
    if err := json.Unmarshal(output, &configData); err != nil {
        return "", "", 0, fmt.Errorf("failed to parse JSON: %v", err)
    }

    // 获取 routing.inbounds[0]
    inbounds, ok := configData["inbounds"].([]interface{})
    if !ok || len(inbounds) == 0 {
        return "", "", 0, fmt.Errorf("failed to find inbounds in config")
    }

    inbound0, ok := inbounds[0].(map[string]interface{})
    if !ok {
        return "", "", 0, fmt.Errorf("failed to parse inbounds[0]")
    }

    settings, ok := inbound0["settings"].(map[string]interface{})
    if !ok {
        return "", "", 0, fmt.Errorf("failed to find settings in inbounds[0]")
    }

    accounts, ok := settings["accounts"].([]interface{})
    if !ok || len(accounts) == 0 {
        return "", "", 0, fmt.Errorf("failed to find accounts in settings")
    }

    account0, ok := accounts[0].(map[string]interface{})
    if !ok {
        return "", "", 0, fmt.Errorf("failed to parse accounts[0]")
    }

    user, _ = account0["user"].(string)
    pass, _ := account0["pass"].(string)
    portNum, _ := inbound0["port"].(float64) // JSON numbers are float64

    return user, pass, int(portNum), nil
}