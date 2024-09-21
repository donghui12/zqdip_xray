# zqdip_xray
xray 配置 &amp;&amp; 安装

 apt-get update &&  apt-get install -y unzip && curl -O http://138.2.231.114/zqdip_xray.zip && unzip zqdip_xray.zip  下载并且解压缩
进入操做系统准备阶段
cd zqdip_xray && bash prepare.sh run
安装xray
apt-get install python3-pip -y # 如果您的系统是基于 Debian/Ubuntu 的

pip3 install psutil

python3 main.py install
python3 main.py config_init --name 555
手动调整参数
(venv01) [root@monther test_psutil]# python3 main.py --help
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
修改配置文件
/usr/local/etc/xray/config.json
重启xray
sudo systemctl restart xray
