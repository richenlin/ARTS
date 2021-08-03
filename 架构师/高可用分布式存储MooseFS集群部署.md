# 高可用分布式存储MooseFS集群部署

## MooseFS
  MooseFS（mfs）被称为对象存储，提供了强大的扩展性、高可靠性和持久性。它可以将文件分布存储于不一样的物理机器上，对外却提供的是一个透明的接口的存储资源池。它还具备在线扩展（这是个很大的好处）、文件切块存储、读写效率高等特色。MFS分布式文件系统由元数据服务器(Master Server)、元数据日志服务器(Metalogger Server)、数据存储服务器(Chunk Server)、客户端(Client)组成。

* 元数据服务器：MFS系统中的核心组成部分，存储每一个文件的元数据，负责文件的读写调度、空间回收和在多个chunk server之间的数据拷贝等。目前MFS仅支持一个元数据服务器，所以可能会出现单点故障。针对此问题咱们须要用一台性能很稳定的服务器来做为咱们的元数据服务器，这样能够下降出现单点故障的几率。

* 元数据日志服务器：元数据服务器的备份节点，按照指定的周期从元数据服务器上将保存元数据、更新日志和会话信息的文件下载到本地目录下。当元数据服务器出现故障时，咱们能够从该服务器的文件中拿到相关的必要的信息对整个系统进行恢复。


* 数据存储服务器：负责链接元数据管理服务器，遵从元数据服务器的调度，提供存储空间，并为客户端提供数据传输，MooseFS提供一个手动指定每一个目录的备份个数。假设个数为n，那么咱们在向系统写入文件时，系统会将切分好的文件块在不一样的chunk server上复制n份。备份数的增长不会影响系统的写性能，可是能够提升系统的读性能和可用性，这能够说是一种以存储容量换取写性能和可用性的策略。

* 客户端：使用mfsmount的方式经过FUSE内核接口挂接远程管理服务器上管理的数据存储服务器到本地目录上，而后就能够像使用本地文件同样来使用咱们的MFS文件系统了。

  MooseFS虽然利用元数据进行了备份，但这只是一种常规的日志备份手段，Master节点（元数据服务器）仍然存在单点问题，在某些状况下并不能完美的接管业务，甚至造成数据丢失。因此我们需要引入Keepalived + Unison 机制来保证MooseFS Master节点双机热备，满足高可用。
  
  
## 资源规划


![1609830219500nc3cyR](https://upic-1258482165.cos.ap-chengdu.myqcloud.com/2021-01-05/1609830219500nc3cyR.png)


| 服务器 | 别名 | CPU | 内存 | 挂载点 | 作用 | 说明 | 
| --- | --- | --- |  --- |  --- | --- | --- |
| 192.168.5.101 | cml1 | >=2核 | >= 16G |  /data/mfs  | Master Servers1 | 元数据节点1，复制调度和管理元数据
| 192.168.5.102 | cml2 | >=2核 | >= 16G | /data/mfs  | Master Servers2 | 元数据节点2，复制调度和管理元数据
| 192.168.5.103 | cml3 | >=2核 | >= 4G | /data/mfs | Metalogger Servers | 用于备份 master 的元数据和日志
| 192.168.5.104 | cml4 | >=2核 | >= 4G | /data/mfs | Chunk Servers | 数据节点，数据实际存放的节点
| 192.168.5.105 | cml5 | >=2核 | >= 4G | /data/mfs | Chunk Servers | 数据节点，数据实际存放的节点
| 192.168.5.129 | cml6 | >=1核 | >= 1G | /data/mfs | Client | 客户端，使用存储的节点

    
系统版本: Centos7

Yum源：http://mirrors.aliyun.com/repo/

## MooseFS集群配置

### 基础配置
在所有集群节点机器上执行：

```bash
[root@cml1 ~]# vi /etc/hosts
...
192.168.5.200 mfsmaster

```
设置时钟同步

```bash
[root@cml1 ~]#crontab -l
*/5 * * * *ntpdate cn.pool.ntp.org
```
设置源

```bash
```bash
[root@cml1 ~]# curl "https://ppa.moosefs.com/RPM-GPG-KEY-MooseFS" > /etc/pki/rpm-gpg/RPM-GPG-KEY-MooseFS
[root@cml1 ~]# curl "http://ppa.moosefs.com/MooseFS-3-el7.repo" > /etc/yum.repos.d/MooseFS.repo
```
```

### Master Server

MFS安装：

```bash
[root@cml1 ~]# yum install moosefs-master moosefs-climoosefs-cgi moosefs-cgiserv
```
开机启动

```bash
[root@cml1 ~]# systemctl enable moosefs-master
[root@cml1 ~]# systemctl start moosefs-master
[root@cml1 ~]# systemctl enable moosefs-cgiserv
[root@cml1 ~]# systemctl start moosefs-cgiserv
```

###  Metalogger Servers

```bash
[root@cml3 ~]# yum install moosefs-metalogger
[root@cml3 ~]# systemctl enable moosefs-metalogger
[root@cml3 ~]# systemctl start moosefs-metalogger
```

###  Chunk Servers


安装

```bash
[root@cml4 ~]# yum install moosefs-chunkserver
[root@cml4 ~]# systemctl enable moosefs-chunkserver
```

创建mfs用户以及用户组

```bash 
[root@cml4 ~]# mkdir -p /data/mfs 
[root@cml4 ~]# chown -R mfs:mfs /data/mfs 
[root@cml4 ~]# echo "/data/mfs " >> /etc/mfs/mfshdd.cfg
[root@cml4 ~]#  systemctl start moosefs-chunkserver
```


### Client端

```bash
[root@cml4 ~]# yum install moosefs-client
```

创建mfs用户以及用户组

```bash 
[root@cml6 ~]# mkdir -p /data/mfs 
[root@cml6 ~]# groupadd mfs      #添加用户组
[root@cml6 ~]# useradd -g mfs mfs -s /sbin/nologin  #添加用户
[root@cml6 ~]# chown -R mfs:mfs /data/mfs 
```
挂载：

```bash
#挂载根目录（/）
[root@cml6 ~]# mfsmount -o nonempty /mnt/mfs-cli/              

#挂载目录（/logs）                                
[root@cml6 ~]# mfsmount mfsmaster:9421:/logs -o nonempty /mnt/mfs-logs/            
```


## 高可用配置

### Keepalived

需要再cml1、cml2分别安装

* 安装

```bash
[root@cml1 ~]# yum install keepalived
[root@cml1 ~]# systemctl enable keepalived
```

* cml1配置文件：

```bash
[root@cml1 ~]# cp /etc/keepalived/keepalived.conf /etc/keepalived/keepalived.conf-bak
[root@cml1 ~]# vi /etc/keepalived/keepalived.conf

! Configuration File for keepalived
global_defs {
  notification_email {
    root@localhost
    }
  
notification_email_from keepalived@localhost
smtp_server 127.0.0.1
smtp_connect_timeout 30
router_id MFS_HA_MASTER
}
  
vrrp_script chk_mfs {                           
  script "/usr/local/mfs/keepalived_check_mfsmaster.sh"
  interval 2
  weight 2
}
  
vrrp_instance VI_1 {
  state MASTER
  interface eth0
  virtual_router_id 51
  priority 100
  advert_int 1
  authentication {
    auth_type PASS
    auth_pass 1111
}
  track_script {
    chk_mfs
}
virtual_ipaddress {
    192.168.5.200
}
notify_master "/etc/keepalived/clean_arp.sh 192.168.5.200"
}
```

* 监控脚本:

```bash
[root@cml1 ~]# vim /usr/local/mfs/keepalived_check_mfsmaster.sh
#!/bin/bash
A=`ps -C mfsmaster --no-header | wc -l`
if [ $A -eq 0 ];then
/etc/init.d/mfsmaster start
sleep 3
   if [ `ps -C mfsmaster --no-header | wc -l ` -eq 0 ];then
      /usr/bin/killall -9 mfscgiserv
      /usr/bin/killall -9 keepalived
   fi
fi

[root@cml1 ~]# chmod 755 /usr/local/mfs/keepalived_check_mfsmaster.sh
```

设置更新虚拟服务器（VIP）地址的arp记录到网关脚本

```bash
[root@cml1 ~]# vim /etc/keepalived/clean_arp.sh 
#!/bin/sh
VIP=$1
GATEWAY=192.168.5.1                                        //这个是网关地址                       
/sbin/arping -I eth0 -c 5 -s $VIP $GATEWAY &>/dev/null

[root@cml1 ~]# chmod 755 /etc/keepalived/clean_arp.sh
```

启动keepalived（确保Keepalived_MASTER机器的mfs master服务和Keepalived服务都要启动）

```
[root@cml1 ~]# /etc/init.d/keepalived start
[root@cml1 ~]# ps -ef|grep keepalived

#查看vip
[root@cml1 ~]# ip addr
```

* Keepalived_BACKUP（mfs master）配置

```bash
[root@cml2 ~]# cp /etc/keepalived/keepalived.conf /etc/keepalived/keepalived.conf-bak
[root@cml2 ~]# vi /etc/keepalived/keepalived.conf
! Configuration File for keepalived
global_defs {
  notification_email {
    root@localhost
    }
  
notification_email_from keepalived@localhost
smtp_server 127.0.0.1
smtp_connect_timeout 30
router_id MFS_HA_BACKUP
}
  
vrrp_script chk_mfs {                           
  script "/usr/local/mfs/keepalived_check_mfsmaster.sh"
  interval 2
  weight 2
}
  
vrrp_instance VI_1 {
  state BACKUP
  interface eth0
  virtual_router_id 51
  priority 99
  advert_int 1
  authentication {
    auth_type PASS
    auth_pass 1111
}
  track_script {
    chk_mfs
}
virtual_ipaddress {
    192.168.5.200
}
notify_master "/etc/keepalived/clean_arp.sh 192.168.5.200"
}
```
监控脚本

```
[root@cml2 ~]# vim /usr/local/mfs/keepalived_notify.sh
#!/bin/bash
A=`ps -C mfsmaster --no-header | wc -l`
if [ $A -eq 0 ];then
/etc/init.d/mfsmaster start
sleep 3
   if [ `ps -C mfsmaster --no-header | wc -l ` -eq 0 ];then
      /usr/bin/killall -9 mfscgiserv
      /usr/bin/killall -9 keepalived
   fi
fi

[root@cml2 ~]# chmod 755 /usr/local/mfs/keepalived_notify.sh
```

设置更新虚拟服务器（VIP）地址的arp记录到网关脚本

```bash
[root@cml2 ~]# vim /etc/keepalived/clean_arp.sh 
#!/bin/sh
VIP=$1
GATEWAY=192.168.5.1                                        
/sbin/arping -I eth0 -c 5 -s $VIP $GATEWAY &>/dev/null
```
启动keepalived

```bash
[root@cml2 ~]# /etc/init.d/keepalived start
```

* iptales防火墙设置

如果开启了iptables防火墙功能，则需要在两台机器的iptables里配置如下：

```bash
[root@cml1 ~]# vim /etc/sysconfig/iptables
........
-A INPUT -s 192.168.5.0/24 -d 224.0.0.18 -j ACCEPT       #允许组播地址通信。注意设置这两行，就会在Keepalived_MASTER故障恢复后，将VIP资源从Keepalived_BACK那里再转移回来
-A INPUT -s 182.148.15.0/24 -p vrrp -j ACCEPT             #允许VRRP（虚拟路由器冗余协）通信
-A INPUT -m state --state NEW -m tcp -p tcp --dport 9419 -j ACCEPT 
-A INPUT -m state --state NEW -m tcp -p tcp --dport 9420 -j ACCEPT  
-A INPUT -m state --state NEW -m tcp -p tcp --dport 9421 -j ACCEPT
-A INPUT -m state --state NEW -m tcp -p tcp --dport 9425 -j ACCEPT 

[root@cml1 ~]# /etc/init.d/iptables start
```
### 数据同步

Unison是Windows、Linux以及其他Unix平台下都可以使用的文件同步工具，它能使两个文件夹（本地或网络上的）保持内容的一致。

* ssh互信

```bash
[root@cml1 ~]#vi /etc/ssh/sshd_config
...
RSAAuthentication yes
PubkeyAuthentication yes
```
拷贝授权

```bash
[root@cml1 ~]#ssh-keygen
[root@cml1 ~]#ssh-copy-id cml2
```

cml2上做相同拷贝授权操作

* 安装

在cml1、cml2上执行安装

```bash
[root@cml1 ~]# yum install  -y  ocaml inotify-tools
[root@cml1 ~]# yum install  -y  unison
```

* Unison配置

```bash
## cml1机器上执行
[root@cml1 ~]# vim /home/mfs/.unison/default.prf
#Unison preferences file
#force=/home/unison/test  local->remote
root=/data/mfs
root=ssh://mfs@192.168.5.102//data/mfs
#path=aaa
#path=bbb
batch=true
#maxthreads=300
owner=true
group=true
perms=-1
fastcheck=false
rsync=false
sshargs=-C
xferbycopying=true
confirmbigdel=false
log=true
logfile=/home/mfs/.unison/unison.log

## cml2机器上执行
[root@cml2 ~]# vim /home/mfs/.unison/default.prf
#Unison preferences file
#force=/home/unison/test  local->remote
root=/data/mfs
root=ssh://mfs@192.168.5.101//data/mfs
#path=aaa
#path=bbb
batch=true
#maxthreads=300
owner=true
group=true
perms=-1
fastcheck=false
rsync=false
sshargs=-C
xferbycopying=true
confirmbigdel=false
log=true
logfile=/home/mfs/.unison/unison.log
```

* inotify实时监控脚本

```bash
[root@cml2 ~]# vim /usr/local/mfs/inotify.sh
#!/bin/bash
SRCDIR=/data/mfs
/usr/local/inotify/bin/inotifywait-mrq  --timefmt '%d/%m/%y %H:%M' --format'%T %w %f'  -e create,delete,modify,move$SRCDIR | while read files
do
unison-servercmd=/usr/local/bin/unison
done
```

开机启动

```bash
[root@cml2 ~]# vi /etc/rc.local
nohup bash /usr/local/mfs/inotify.sh > /var/log/vip_watch.log 2>&1 &

```

