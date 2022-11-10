# JunOS

## Configure

```shell
set services analytics streaming-server JSA-Server remote-address 192.0.2.22
set services analytics streaming-server JSA-Server remote-port 21000
set services analytics export-profile JSA-Export local-address 192.0.2.1
set services analytics export-profile JSA-Export reporting-rate 60
set services analytics export-profile JSA-Export payload-size 1480
set services analytics export-profile JSA-Export format gpb
set services analytics export-profile JSA-Export transport udp
set services analytics sensor JSA-metrics server-name JSA-Server
set services analytics sensor JSA-metrics export-name JSA-Export
set services analytics sensor JSA-metrics resource /junos/system/linecard/firewall/
set services analytics sensor JSA-metrics resource-filter uplinks.*
```
