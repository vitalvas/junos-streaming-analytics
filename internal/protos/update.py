#!/usr/bin/env python3

import os
import urllib.request
import subprocess

VERSION = "22.2"

NAMES = [
    "telemetry_top.proto",
    "firewall.proto",
    "port.proto",
]

curdir = os.path.dirname(os.path.realpath(__file__))
relpath = []

for dirpath, dirnames, filenames in os.walk(curdir):    
    for file in filenames:
        if file.endswith((".proto", ".pb.go")):
            os.remove(os.path.join(dirpath, file))

for row in NAMES:
    with urllib.request.urlopen(f"https://raw.githubusercontent.com/Juniper/telemetry/master/{VERSION}/{VERSION}R1/protos/junos-telemetry-interface/{row}") as f:
        payload = f.read().decode('utf-8')
        print(f"Updating {row}")
        
        file_path = os.path.join(curdir, row)

        with open(file_path, "w") as f2:
            f2.write(payload)

        relpath.append((curdir[len(os.getcwd())+1:], row))

# os.chdir(curdir)


build_args = []
for mdir, mfile in relpath:
    build_args.append(f"--go_opt=M{mfile}={mdir}")

for mdir, mfile in relpath:
    exec_args = ["protoc", f"--go_out={mdir}/jti", "--go_opt=paths=source_relative", f"--proto_path={mdir}"] + build_args + [mfile]
    
    print(' '.join(exec_args))
    subprocess.check_call(exec_args, shell=False)
    