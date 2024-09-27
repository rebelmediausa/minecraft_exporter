# Systemd Unit

The unit files (`*.service` and `*.socket`) in this directory are to be put into `/etc/systemd/system`.
It needs a user named `minecraft_exporter`, whose shell should be `/sbin/nologin` and should not have any special privileges.
It needs a sysconfig file in `/etc/sysconfig/minecraft_exporter`.
A sample file can be found in `sysconfig.minecraft_exporter`.
