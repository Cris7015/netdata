<!--startmeta
custom_edit_url: "https://github.com/netdata/netdata/edit/master/src/go/plugin/go.d/modules/haproxy/README.md"
meta_yaml: "https://github.com/netdata/netdata/edit/master/src/go/plugin/go.d/modules/haproxy/metadata.yaml"
sidebar_label: "HAProxy"
learn_status: "Published"
learn_rel_path: "Collecting Metrics/Web Servers and Web Proxies"
most_popular: False
message: "DO NOT EDIT THIS FILE DIRECTLY, IT IS GENERATED BY THE COLLECTOR'S metadata.yaml FILE"
endmeta-->

# HAProxy


<img src="https://netdata.cloud/img/haproxy.svg" width="150"/>


Plugin: go.d.plugin
Module: haproxy

<img src="https://img.shields.io/badge/maintained%20by-Netdata-%2300ab44" />

## Overview

This collector monitors HAProxy servers.




This collector is supported on all platforms.

This collector supports collecting metrics from multiple instances of this integration, including remote instances.


### Default Behavior

#### Auto-Detection

This integration doesn't support auto-detection.

#### Limits

The default configuration for this integration does not impose any limits on data collection.

#### Performance Impact

The default configuration for this integration is not expected to impose a significant performance impact on the system.


## Metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.



### Per HAProxy instance

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric | Dimensions | Unit |
|:------|:----------|:----|
| haproxy.backend_current_sessions | a dimension per proxy | sessions |
| haproxy.backend_sessions | a dimension per proxy | sessions/s |
| haproxy.backend_response_time_average | a dimension per proxy | milliseconds |
| haproxy.backend_queue_time_average | a dimension per proxy | milliseconds |
| haproxy.backend_current_queue | a dimension per proxy | requests |

### Per proxy

These metrics refer to the Proxy.

This scope has no labels.

Metrics:

| Metric | Dimensions | Unit |
|:------|:----------|:----|
| haproxy.backend_http_responses | 1xx, 2xx, 3xx, 4xx, 5xx, other | responses/s |
| haproxy.backend_network_io | in, out | bytes/s |



## Alerts

There are no alerts configured by default for this integration.


## Setup

### Prerequisites

#### Enable PROMEX addon.

To enable PROMEX addon, follow the [official documentation](https://github.com/haproxy/haproxy/tree/master/addons/promex).



### Configuration

#### File

The configuration file name for this integration is `go.d/haproxy.conf`.


You can edit the configuration file using the [`edit-config`](https://github.com/netdata/netdata/blob/master/docs/netdata-agent/configuration/README.md#edit-a-configuration-file-using-edit-config) script from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/netdata-agent/configuration/README.md#the-netdata-config-directory).

```bash
cd /etc/netdata 2>/dev/null || cd /opt/netdata/etc/netdata
sudo ./edit-config go.d/haproxy.conf
```
#### Options

The following options can be defined globally: update_every, autodetection_retry.


<details open><summary>Config options</summary>

| Name | Description | Default | Required |
|:----|:-----------|:-------|:--------:|
| update_every | Data collection frequency. | 1 | no |
| autodetection_retry | Recheck interval in seconds. Zero means no recheck will be scheduled. | 0 | no |
| url | Server URL. | http://127.0.0.1 | yes |
| timeout | HTTP request timeout. | 1 | no |
| username | Username for basic HTTP authentication. |  | no |
| password | Password for basic HTTP authentication. |  | no |
| proxy_url | Proxy URL. |  | no |
| proxy_username | Username for proxy basic HTTP authentication. |  | no |
| proxy_password | Password for proxy basic HTTP authentication. |  | no |
| method | HTTP request method. | GET | no |
| body | HTTP request body. |  | no |
| headers | HTTP request headers. |  | no |
| not_follow_redirects | Redirect handling policy. Controls whether the client follows redirects. | no | no |
| tls_skip_verify | Server certificate chain and hostname validation policy. Controls whether the client performs this check. | no | no |
| tls_ca | Certification authority that the client uses when verifying the server's certificates. |  | no |
| tls_cert | Client TLS certificate. |  | no |
| tls_key | Client TLS key. |  | no |

</details>

#### Examples

##### Basic

A basic example configuration.

<details open><summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8404/metrics

```
</details>

##### HTTP authentication

Basic HTTP authentication.

<details open><summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8404/metrics
    username: username
    password: password

```
</details>

##### HTTPS with self-signed certificate

NGINX Plus with enabled HTTPS and self-signed certificate.

<details open><summary>Config</summary>

```yaml
jobs:
  - name: local
    url: https://127.0.0.1:8404/metrics
    tls_skip_verify: yes

```
</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Collecting metrics from local and remote instances.


<details open><summary>Config</summary>

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:8404/metrics

  - name: remote
    url: http://192.0.2.1:8404/metrics

```
</details>



## Troubleshooting

### Debug Mode

**Important**: Debug mode is not supported for data collection jobs created via the UI using the Dyncfg feature.

To troubleshoot issues with the `haproxy` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m haproxy
  ```

### Getting Logs

If you're encountering problems with the `haproxy` collector, follow these steps to retrieve logs and identify potential issues:

- **Run the command** specific to your system (systemd, non-systemd, or Docker container).
- **Examine the output** for any warnings or error messages that might indicate issues.  These messages should provide clues about the root cause of the problem.

#### System with systemd

Use the following command to view logs generated since the last Netdata service restart:

```bash
journalctl _SYSTEMD_INVOCATION_ID="$(systemctl show --value --property=InvocationID netdata)" --namespace=netdata --grep haproxy
```

#### System without systemd

Locate the collector log file, typically at `/var/log/netdata/collector.log`, and use `grep` to filter for collector's name:

```bash
grep haproxy /var/log/netdata/collector.log
```

**Note**: This method shows logs from all restarts. Focus on the **latest entries** for troubleshooting current issues.

#### Docker Container

If your Netdata runs in a Docker container named "netdata" (replace if different), use this command:

```bash
docker logs netdata 2>&1 | grep haproxy
```

