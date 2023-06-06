# snmp-agent-akvorado

snmp-agent-akvorado は [Akvorado](https://github.com/akvorado/akvorado) 用の SNMP エージェントです。snmp-agent-akvorado は snmpd の代替ではありません。

## 使い方

本リポジトリの [Release](https://github.com/cuteip/snmp-agent-akvorado/releases) からダウンロードする、もしくは git clone してビルドします。以降は Release からダウンロードして利用する方法について説明します。

```shell
VERSION=v0.2.0
wget https://github.com/cuteip/snmp-agent-akvorado/releases/download/${VERSION}/snmp-agent-akvorado_${VERSION}_linux_x86_64.tar.gz
tar xzvf snmp-agent-akvorado_${VERSION}_linux_x86_64.tar.gz
./snmp-agent-akvorado run --listen-addr 192.0.2.1:161
```

## systemd service ファイル例

```ini
[Unit]
Description=https://github.com/cuteip/snmp-agent-akvorado
After=network.target

[Service]
Type=simple
User=nobody
Group=nogroup
Restart=on-failure
RestartSec=10
ExecStart=/usr/local/bin/snmp-agent-akvorado run \
    --listen-addr 192.0.2.1:161

SyslogIdentifier=snmp-agent-akvorado

# 161/UDP のため
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
```
