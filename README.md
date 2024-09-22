# zqdip_xray
xray 配置 &amp;&amp; 安装

### 系统准备阶段
cd zqdip_xray && bash prepare.sh run

### 安装
python3 main.py install

### 配置
python3 main.py config_init --name 555

```shell
usage: main.py [-h] [--list]
               {install,config_init,uninstall,status,show_config} ...

站群服务器隧道管理脚本

positional arguments:
  {install,config_init,uninstall,status,show_config}
                        选择进入子菜单功能
    install             完整安装Xray【不包含配置】
    config_init         进行配置初始化并重载内核设置
    uninstall           从这个电脑上完全移除站群管理服务
    status              查看xray运行状态
    show_config         查看文件中的配置

optional arguments:
  -h, --help            show this help message and exit
  --list, -L            列出站群服务器内的所有节点
```

### 修改配置文件
/usr/local/etc/xray/config.json