# GPB files

Download:

```shell
wget -O telemetry_top.proto https://raw.githubusercontent.com/Juniper/telemetry/master/22.1/22.1R1/protos/junos-telemetry-interface/telemetry_top.proto
wget -O firewall.proto https://raw.githubusercontent.com/Juniper/telemetry/master/22.1/22.1R1/protos/junos-telemetry-interface/firewall.proto
wget -O port.proto https://raw.githubusercontent.com/Juniper/telemetry/master/22.1/22.1R1/protos/junos-telemetry-interface/port.proto
```

Add to eatch file:

```proto
option go_package = "../jti"; // !! Custom
```

Generate:

```shell
protoc --proto_path=./internal/protos --go_out=./internal/jti internal/protos/telemetry_top.proto
protoc --proto_path=./internal/protos --go_out=./internal/jti internal/protos/firewall.proto
protoc --proto_path=./internal/protos --go_out=./internal/jti internal/protos/port.proto
```
