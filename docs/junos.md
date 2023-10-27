# JunOS

## Configure

* JTI - listen to the UDP socket and accept protobuf from juniper
* GNMI - connecting to switches via GRPC and receiving metrics (not implemented yet)

### Using JTI (Juniper Telemetry Interface)

#### Configuring the JTI collector

```shell
set services analytics streaming-server jsa-collector remote-address 192.0.2.22
set services analytics streaming-server jsa-collector remote-port 21000
```

```shell
set services analytics export-profile jsa-export local-address 192.0.2.1
set services analytics export-profile jsa-export reporting-rate 10
set services analytics export-profile jsa-export payload-size 9000
set services analytics export-profile jsa-export format gpb
set services analytics export-profile jsa-export transport udp
set services analytics export-profile jsa-export forwarding-class network-control
```

#### Exporting firewall statistics

```shell
set services analytics sensor jsa-system-linecard-firewall server-name jsa-collector
set services analytics sensor jsa-system-linecard-firewall export-name jsa-export
set services analytics sensor jsa-system-linecard-firewall resource /junos/system/linecard/firewall/
```

#### Exporting interface statistics

Already include:

* `/junos/system/linecard/interface/queue/`

```shell
set services analytics sensor jsa-system-linecard-interface server-name jsa-collector
set services analytics sensor jsa-system-linecard-interface export-name jsa-export
set services analytics sensor jsa-system-linecard-interface resource /junos/system/linecard/interface/
```
