# My Configuration

> Why not use [collectd](https://collectd.org/)?
>
> Collectd is fine, but I only need occasional CPU, MEM and DISK usages.
> Also for the lulz.

## Commands

Create the necessary folders
```bash
mkdir $HOME/.empty  # to measure space of /home partition
sudo mkdir /var/lib/influxdb /var/log/grafana  # keep data after removing containers
sudo chmod 0777 /var/lib/influxdb /var/log/grafana
```

Command to launch `sysmon`
```bash
docker run --rm -v $HOME/.empty:/home sysmon
```

Command to launch `influxdb`
```bash
docker run --rm --name influxc \
    -p 172.17.0.1:8086:8086 \
    -v /var/lib/influxdb:/var/lib/influxdb \
    influxdb
```

Command to launch `grafana`
```bash
docker run --rm --name grafanac \
    -p 172.17.0.1:3000:3000 \
    -v /var/lib/grafana:/var/lib/grafana \
    grafana/grafana
```

---

## Services

> In order of requirement, files located at `/etc/systemd/system`

`influxd.service`
```bash
[Unit]
Description=Running InfluxDB in a container
Requires=docker.service
After=docker.service
[Service]
User=mikel
ExecStart=/bin/bash -c "docker run --rm --name influxc \
    -p 172.17.0.1:8086:8086 \
    -v /var/lib/influxdb:/var/lib/influxdb \
    influxdb"
ExecStop=/bin/bash -c "docker stop influxc"
[Install]
WantedBy=multi-user.target
```

`grafanad.service`
```bash
[Unit]
Description=Running Grafana in a container
BindsTo=docker.service influxd.service
After=docker.service influxd.service
[Service]
User=mikel
ExecStart=/bin/bash -c "docker run --rm --name grafanac \
    -p 172.17.0.1:3000:3000 \
    -v /var/lib/grafana:/var/lib/grafana \
    grafana/grafana"
ExecStop=/bin/bash -c "docker stop grafanac"
[Install]
WantedBy=multi-user.target
```

`sysmond.service`
```bash
[Unit]
Description=Monitoring CPU, MEM and DISK usage
BindsTo=docker.service influxd.service
After=docker.service influxd.service
[Service]
User=mikel
ExecStart=/bin/bash -c "docker run --rm --name sysmonc \
    -v /home/mikel/.empty:/home sysmon"
ExecStop=/bin/bash -c "docker stop sysmonc"
[Install]
WantedBy=multi-user.target
```
---
To enable all of them
```bash
sudo systemctl enable grafanad influxd sysmond
```
