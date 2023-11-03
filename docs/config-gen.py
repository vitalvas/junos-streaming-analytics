#!/usr/bin/env python3

import argparse


class ConfigGenerator:
    def __init__(self, config: dict) -> None:
        self.config = config

    def _streaming_server(self) -> list:
        return [
            f"set services analytics streaming-server {self.config['streaming_server_name']} remote-address {self.config['collector_address']}",
            f"set services analytics streaming-server {self.config['streaming_server_name']} remote-port {self.config['collector_port']}",
        ]

    def _export_profile(self) -> list:
        data = [
            f"set services analytics export-profile {self.config['export_profile_name']} reporting-rate {self.config['export_reporting_rate']}",
            f"set services analytics export-profile {self.config['export_profile_name']} payload-size 3000",
            f"set services analytics export-profile {self.config['export_profile_name']} format gpb",
            f"set services analytics export-profile {self.config['export_profile_name']} transport udp",
            f"set services analytics export-profile {self.config['export_profile_name']} forwarding-class network-control",
        ]

        if self.config.get('device_address'):
            data.append(f"set services analytics export-profile {self.config['export_profile_name']} local-address {self.config['device_address']}")
        
        return data
    
    def _sensor(self) -> list:
        resources = [
            "/junos/system/linecard/firewall/",
            "/junos/system/linecard/interface/",
            "/junos/system/linecard/interface/logical/family/ipv4/usage/",
            "/junos/system/linecard/interface/logical/family/ipv6/usage/",
            "/junos/system/linecard/interface/logical/usage/",
            "/junos/system/linecard/interface/queue/",
            "/junos/system/linecard/npu/memory/",
            "/junos/system/linecard/optics/",
        ]

        sorted(resources)
        data = []

        for resource in resources:
            sensor_suffix = resource.replace("/", "-")[1:-1]
            sensor_name = f"{self.config['sensor_prefix']}-{sensor_suffix}"

            data.extend([
                f"set services analytics sensor {sensor_name} server-name {self.config['streaming_server_name']}",
                f"set services analytics sensor {sensor_name} export-name {self.config['export_profile_name']}",
                f"set services analytics sensor {sensor_name} resource {resource}"
            ])
        
        return data

    def generate(self):
        rows = []
        rows.extend(self._streaming_server())
        rows.extend(self._export_profile())
        rows.extend(self._sensor())
        print("\n".join(rows))


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Generate config for Juniper Telemetry Streaming')
    parser.add_argument('--collector-address', type=str, required=True)
    parser.add_argument('--collector-port', type=int, default=21000)
    parser.add_argument('--device-address', type=str)
    parser.add_argument('--export-profile-name', type=str, default='jsa-export')
    parser.add_argument('--export-reporting-rate', type=int, default=10)
    parser.add_argument('--sensor-prefix', type=str, default='jsa')
    parser.add_argument('--streaming-server-name', type=str, default='jsa-collector')
    
    args = parser.parse_args()

    ConfigGenerator(config=args.__dict__).generate()
