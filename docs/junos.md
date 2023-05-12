# JunOS

## Configure

```shell
set services analytics streaming-server JSA-Server remote-address 192.0.2.22
set services analytics streaming-server JSA-Server remote-port 21000
```

```shell
set services analytics export-profile JSA-Export local-address 192.0.2.1
set services analytics export-profile JSA-Export reporting-rate 30
set services analytics export-profile JSA-Export payload-size 1480
set services analytics export-profile JSA-Export format gpb
set services analytics export-profile JSA-Export transport udp
```

```shell
set services analytics sensor JSA-system-linecard-firewall server-name JSA-Server
set services analytics sensor JSA-system-linecard-firewall export-name JSA-Export
set services analytics sensor JSA-system-linecard-firewall resource /junos/system/linecard/firewall/
set services analytics sensor JSA-system-linecard-firewall resource-filter uplinks.*
```

```shell
set services analytics sensor JSA-system-linecard-interface server-name JSA-Server
set services analytics sensor JSA-system-linecard-interface export-name JSA-Export
set services analytics sensor JSA-system-linecard-interface resource /junos/system/linecard/interface/
```
