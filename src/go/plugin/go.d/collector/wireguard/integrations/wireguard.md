<!--startmeta
custom_edit_url: "https://github.com/netdata/netdata/edit/master/src/go/plugin/go.d/modules/wireguard/README.md"
meta_yaml: "https://github.com/netdata/netdata/edit/master/src/go/plugin/go.d/modules/wireguard/metadata.yaml"
sidebar_label: "WireGuard"
learn_status: "Published"
learn_rel_path: "Collecting Metrics/VPNs"
most_popular: False
message: "DO NOT EDIT THIS FILE DIRECTLY, IT IS GENERATED BY THE COLLECTOR'S metadata.yaml FILE"
endmeta-->

# WireGuard


<img src="https://netdata.cloud/img/wireguard.svg" width="150"/>


Plugin: go.d.plugin
Module: wireguard

<img src="https://img.shields.io/badge/maintained%20by-Netdata-%2300ab44" />

## Overview

This collector monitors WireGuard VPN devices and peers traffic.


It connects to the local WireGuard instance using [wireguard-go client](https://github.com/WireGuard/wireguard-go).


This collector is supported on all platforms.

This collector supports collecting metrics from multiple instances of this integration, including remote instances.

This collector requires the CAP_NET_ADMIN capability, but it is set automatically during installation, so no manual configuration is needed.


### Default Behavior

#### Auto-Detection

It automatically detects instances running on localhost.


#### Limits

Doesn't work if Netdata or WireGuard is installed in the container.


#### Performance Impact

The default configuration for this integration is not expected to impose a significant performance impact on the system.


## Metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.



### Per device

These metrics refer to the VPN network interface.

Labels:

| Label      | Description     |
|:-----------|:----------------|
| device | VPN network interface |

Metrics:

| Metric | Dimensions | Unit |
|:------|:----------|:----|
| wireguard.device_network_io | receive, transmit | B/s |
| wireguard.device_peers | peers | peers |

### Per peer

These metrics refer to the VPN peer.

Labels:

| Label      | Description     |
|:-----------|:----------------|
| device | VPN network interface |
| public_key | Public key of a peer |

Metrics:

| Metric | Dimensions | Unit |
|:------|:----------|:----|
| wireguard.peer_network_io | receive, transmit | B/s |
| wireguard.peer_latest_handshake_ago | time | seconds |



## Alerts

There are no alerts configured by default for this integration.


## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name for this integration is `go.d/wireguard.conf`.


You can edit the configuration file using the [`edit-config`](https://github.com/netdata/netdata/blob/master/docs/netdata-agent/configuration/README.md#edit-a-configuration-file-using-edit-config) script from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/netdata-agent/configuration/README.md#the-netdata-config-directory).

```bash
cd /etc/netdata 2>/dev/null || cd /opt/netdata/etc/netdata
sudo ./edit-config go.d/wireguard.conf
```
#### Options

The following options can be defined globally: update_every, autodetection_retry.


<details open><summary>Config options</summary>

| Name | Description | Default | Required |
|:----|:-----------|:-------|:--------:|
| update_every | Data collection frequency. | 1 | no |
| autodetection_retry | Recheck interval in seconds. Zero means no recheck will be scheduled. | 0 | no |

</details>

#### Examples
There are no configuration examples.



## Troubleshooting

### Debug Mode

**Important**: Debug mode is not supported for data collection jobs created via the UI using the Dyncfg feature.

To troubleshoot issues with the `wireguard` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

- Navigate to the `plugins.d` directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on
  your system, open `netdata.conf` and look for the `plugins` setting under `[directories]`.

  ```bash
  cd /usr/libexec/netdata/plugins.d/
  ```

- Switch to the `netdata` user.

  ```bash
  sudo -u netdata -s
  ```

- Run the `go.d.plugin` to debug the collector:

  ```bash
  ./go.d.plugin -d -m wireguard
  ```

### Getting Logs

If you're encountering problems with the `wireguard` collector, follow these steps to retrieve logs and identify potential issues:

- **Run the command** specific to your system (systemd, non-systemd, or Docker container).
- **Examine the output** for any warnings or error messages that might indicate issues.  These messages should provide clues about the root cause of the problem.

#### System with systemd

Use the following command to view logs generated since the last Netdata service restart:

```bash
journalctl _SYSTEMD_INVOCATION_ID="$(systemctl show --value --property=InvocationID netdata)" --namespace=netdata --grep wireguard
```

#### System without systemd

Locate the collector log file, typically at `/var/log/netdata/collector.log`, and use `grep` to filter for collector's name:

```bash
grep wireguard /var/log/netdata/collector.log
```

**Note**: This method shows logs from all restarts. Focus on the **latest entries** for troubleshooting current issues.

#### Docker Container

If your Netdata runs in a Docker container named "netdata" (replace if different), use this command:

```bash
docker logs netdata 2>&1 | grep wireguard
```

