// SPDX-License-Identifier: GPL-3.0-or-later

package nats

import (
	"fmt"
	"maps"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
)

const (
	prioServerTraffic = module.Priority + iota
	prioServerMessages
	prioServerConnections
	prioServerConnectionsRate
	prioHttpEndpointRequests
	prioServerHealthProbeStatus
	prioServerCpuUsage
	prioServerMemoryUsage
	prioServerUptime

	prioAccountTraffic
	prioAccountMessages
	prioAccountConnections
	prioAccountConnectionsRate
	prioAccountSubscriptions
	prioAccountSlowConsumers
	prioAccountLeafNodes

	prioRouteTraffic
	prioRouteMessages
	prioRouteSubscriptions

	prioGatewayConnTraffic
	prioGatewayConnMessages
	prioGatewayConnSubscriptions
	prioGatewayConnUptime
)

var serverCharts = func() module.Charts {
	charts := module.Charts{
		chartServerConnectionsCurrent.Copy(),
		chartServerConnectionsRate.Copy(),
		chartServerTraffic.Copy(),
		chartServerMessages.Copy(),
		chartServerHealthProbeStatus.Copy(),
		chartServerCpuUsage.Copy(),
		chartServerMemUsage.Copy(),
		chartServerUptime.Copy(),
	}
	charts = append(charts, httpEndpointsCharts()...)
	return charts
}()

var (
	chartServerTraffic = module.Chart{
		ID:       "server_traffic",
		Title:    "Server Traffic",
		Units:    "bytes/s",
		Fam:      "traffic",
		Ctx:      "nats.server_traffic",
		Priority: prioServerTraffic,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "varz_srv_in_bytes", Name: "received", Algo: module.Incremental},
			{ID: "varz_srv_out_bytes", Name: "sent", Mul: -1, Algo: module.Incremental},
		},
	}
	chartServerMessages = module.Chart{
		ID:       "server_messages",
		Title:    "Server Messages",
		Units:    "messages/s",
		Fam:      "traffic",
		Ctx:      "nats.server_messages",
		Priority: prioServerMessages,
		Dims: module.Dims{
			{ID: "varz_srv_in_msgs", Name: "received", Algo: module.Incremental},
			{ID: "varz_srv_out_msgs", Name: "sent", Mul: -1, Algo: module.Incremental},
		},
	}
	chartServerConnectionsCurrent = module.Chart{
		ID:       "server_connections",
		Title:    "Server Active Connections",
		Units:    "connections",
		Fam:      "connections",
		Ctx:      "nats.server_connections",
		Priority: prioServerConnections,
		Dims: module.Dims{
			{ID: "varz_srv_connections", Name: "active"},
		},
	}
	chartServerConnectionsRate = module.Chart{
		ID:       "server_connections_rate",
		Title:    "Server Connections",
		Units:    "connections/s",
		Fam:      "connections",
		Ctx:      "nats.server_connections_rate",
		Priority: prioServerConnectionsRate,
		Dims: module.Dims{
			{ID: "varz_srv_total_connections", Name: "connections", Algo: module.Incremental},
		},
	}
	chartServerHealthProbeStatus = module.Chart{
		ID:       "server_health_probe_status",
		Title:    "Server Health Probe Status",
		Units:    "status",
		Fam:      "health",
		Ctx:      "nats.server_health_probe_status",
		Priority: prioServerHealthProbeStatus,
		Dims: module.Dims{
			{ID: "varz_srv_healthz_status_ok", Name: "ok"},
			{ID: "varz_srv_healthz_status_error", Name: "error"},
		},
	}
	chartServerCpuUsage = module.Chart{
		ID:       "server_cpu_usage",
		Title:    "Server CPU Usage",
		Units:    "percent",
		Fam:      "rusage",
		Ctx:      "nats.server_cpu_usage",
		Priority: prioServerCpuUsage,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "varz_srv_cpu", Name: "used"},
		},
	}
	chartServerMemUsage = module.Chart{
		ID:       "server_mem_usage",
		Title:    "Server Memory Usage",
		Units:    "bytes",
		Fam:      "rusage",
		Ctx:      "nats.server_mem_usage",
		Priority: prioServerMemoryUsage,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "varz_srv_mem", Name: "used"},
		},
	}
	chartServerUptime = module.Chart{
		ID:       "server_uptime",
		Title:    "Server Uptime",
		Units:    "seconds",
		Fam:      "uptime",
		Ctx:      "nats.server_uptime",
		Priority: prioServerUptime,
		Dims: module.Dims{
			{ID: "varz_srv_uptime", Name: "uptime"},
		},
	}
)

func httpEndpointsCharts() module.Charts {
	var charts module.Charts

	for _, path := range httpEndpoints {
		chart := httpEndpointRequestsChartTmpl.Copy()

		chart.ID = fmt.Sprintf(chart.ID, path)
		chart.Labels = []module.Label{
			{Key: "http_endpoint", Value: path},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, path)
		}
		charts = append(charts, chart)
	}

	return charts
}

var httpEndpointRequestsChartTmpl = module.Chart{
	ID:       "http_endpoint_%s_requests",
	Title:    "HTTP Endpoint Requests",
	Units:    "requests/s",
	Fam:      "http requests",
	Ctx:      "nats.http_endpoint_requests",
	Priority: prioHttpEndpointRequests,
	Dims: module.Dims{
		{ID: "varz_http_endpoint_%s_req", Name: "requests", Algo: module.Incremental},
	},
}

var accountChartsTmpl = module.Charts{
	accountTrafficTmpl.Copy(),
	accountMessagesTmpl.Copy(),
	accountConnectionsCurrentTmpl.Copy(),
	accountConnectionsRateTmpl.Copy(),
	accountSubscriptionsTmpl.Copy(),
	accountSlowConsumersTmpl.Copy(),
	accountLeadNodesTmpl.Copy(),
}

var (
	accountTrafficTmpl = module.Chart{
		ID:       "account_%s_traffic",
		Title:    "Account Traffic",
		Units:    "bytes/s",
		Fam:      "acc traffic",
		Ctx:      "nats.account_traffic",
		Priority: prioAccountTraffic,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "accstatz_acc_%s_received_bytes", Name: "received", Algo: module.Incremental},
			{ID: "accstatz_acc_%s_sent_bytes", Name: "sent", Mul: -1, Algo: module.Incremental},
		},
	}
	accountMessagesTmpl = module.Chart{
		ID:       "account_%s_messages",
		Title:    "Account Messages",
		Units:    "messages/s",
		Fam:      "acc traffic",
		Ctx:      "nats.account_messages",
		Priority: prioAccountMessages,
		Type:     module.Line,
		Dims: module.Dims{
			{ID: "accstatz_acc_%s_received_msgs", Name: "received", Algo: module.Incremental},
			{ID: "accstatz_acc_%s_sent_msgs", Name: "sent", Mul: -1, Algo: module.Incremental},
		},
	}
	accountConnectionsCurrentTmpl = module.Chart{
		ID:       "account_%s_connections",
		Title:    "Account Active Connections",
		Units:    "connections",
		Fam:      "acc connections",
		Ctx:      "nats.account_connections",
		Priority: prioAccountConnections,
		Type:     module.Line,
		Dims: module.Dims{
			{ID: "accstatz_acc_%s_conns", Name: "active"},
		},
	}
	accountConnectionsRateTmpl = module.Chart{
		ID:       "account_%s_connections_rate",
		Title:    "Account Connections",
		Units:    "connections/s",
		Fam:      "acc connections",
		Ctx:      "nats.account_connections_rate",
		Priority: prioAccountConnectionsRate,
		Type:     module.Line,
		Dims: module.Dims{
			{ID: "accstatz_acc_%s_total_conns", Name: "connections", Algo: module.Incremental},
		},
	}
	accountSubscriptionsTmpl = module.Chart{
		ID:       "account_%s_subscriptions",
		Title:    "Account Active Subscriptions",
		Units:    "subscriptions",
		Fam:      "acc subscriptions",
		Ctx:      "nats.account_subscriptions",
		Priority: prioAccountSubscriptions,
		Type:     module.Line,
		Dims: module.Dims{
			{ID: "accstatz_acc_%s_num_subs", Name: "active"},
		},
	}
	accountSlowConsumersTmpl = module.Chart{
		ID:       "account_%s_slow_consumers",
		Title:    "Account Slow Consumers",
		Units:    "consumers/s",
		Fam:      "acc consumers",
		Ctx:      "nats.account_slow_consumers",
		Priority: prioAccountSlowConsumers,
		Type:     module.Line,
		Dims: module.Dims{
			{ID: "accstatz_acc_%s_slow_consumers", Name: "slow", Algo: module.Incremental},
		},
	}
	accountLeadNodesTmpl = module.Chart{
		ID:       "account_%s_leaf_nodes",
		Title:    "Account Leaf Nodes",
		Units:    "servers",
		Fam:      "acc leaf nodes",
		Ctx:      "nats.account_leaf_nodes",
		Priority: prioAccountLeafNodes,
		Type:     module.Line,
		Dims: module.Dims{
			{ID: "accstatz_acc_%s_leaf_nodes", Name: "leafnode"},
		},
	}
)

var routeChartsTmpl = module.Charts{
	routeTrafficTmpl.Copy(),
	routeMessagesTmpl.Copy(),
	routeSubscriptionsTmpl.Copy(),
}

var (
	routeTrafficTmpl = module.Chart{
		ID:       "route_%d_traffic",
		Title:    "Route Traffic",
		Units:    "bytes/s",
		Fam:      "route traffic",
		Ctx:      "nats.route_traffic",
		Priority: prioRouteTraffic,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "routez_route_id_%d_in_bytes", Name: "in", Algo: module.Incremental},
			{ID: "routez_route_id_%d_out_bytes", Name: "out", Mul: -1, Algo: module.Incremental},
		},
	}
	routeMessagesTmpl = module.Chart{
		ID:       "route_%d_messages",
		Title:    "Route Messages",
		Units:    "messages/s",
		Fam:      "route traffic",
		Ctx:      "nats.route_messages",
		Priority: prioRouteMessages,
		Type:     module.Line,
		Dims: module.Dims{
			{ID: "routez_route_id_%d_in_msgs", Name: "in", Algo: module.Incremental},
			{ID: "routez_route_id_%d_out_msgs", Name: "out", Mul: -1, Algo: module.Incremental},
		},
	}
	routeSubscriptionsTmpl = module.Chart{
		ID:       "route_%d_subscriptions",
		Title:    "Route Active Subscriptions",
		Units:    "subscriptions",
		Fam:      "route subscriptions",
		Ctx:      "nats.route_subscriptions",
		Priority: prioRouteSubscriptions,
		Type:     module.Line,
		Dims: module.Dims{
			{ID: "routez_route_id_%d_num_subs", Name: "active"},
		},
	}
)

var gatewayConnChartsTmpl = module.Charts{
	gatewayConnTrafficTmpl.Copy(),
	gatewayConnMessagesTmpl.Copy(),
	gatewayConnSubscriptionsTmpl.Copy(),
	gatewayConnUptime.Copy(),
}

var (
	gatewayConnTrafficTmpl = module.Chart{
		ID:       "%s_gw_%s_cid_%d_traffic",
		Title:    "%s Gateway Traffic",
		Units:    "bytes/s",
		Fam:      "gw traffic",
		Ctx:      "nats.%s_gateway_conn_traffic",
		Priority: prioGatewayConnTraffic,
		Type:     module.Area,
		Dims: module.Dims{
			{ID: "gatewayz_%s_gw_%s_cid_%d_in_bytes", Name: "in", Algo: module.Incremental},
			{ID: "gatewayz_%s_gw_%s_cid_%d_out_bytes", Name: "out", Mul: -1, Algo: module.Incremental},
		},
	}
	gatewayConnMessagesTmpl = module.Chart{
		ID:       "%s_gw_%s_cid_%d_messages",
		Title:    "%s Gateway Messages",
		Units:    "messages/s",
		Fam:      "gw traffic",
		Ctx:      "nats.%s_gateway_conn_messages",
		Priority: prioGatewayConnMessages,
		Type:     module.Line,
		Dims: module.Dims{
			{ID: "gatewayz_%s_gw_%s_cid_%d_in_msgs", Name: "in", Algo: module.Incremental},
			{ID: "gatewayz_%s_gw_%s_cid_%d_out_msgs", Name: "out", Mul: -1, Algo: module.Incremental},
		},
	}
	gatewayConnSubscriptionsTmpl = module.Chart{
		ID:       "%s_gw_%s_cid_%d_subscriptions",
		Title:    "%s Gateway Active Subscriptions",
		Units:    "subscriptions",
		Fam:      "gw subscriptions",
		Ctx:      "nats.%s_gateway_conn_subscriptions",
		Priority: prioGatewayConnSubscriptions,
		Type:     module.Line,
		Dims: module.Dims{
			{ID: "gatewayz_%s_gw_%s_cid_%d_num_subs", Name: "active"},
		},
	}
	gatewayConnUptime = module.Chart{
		ID:       "%s_gw_%s_cid_%d_uptime",
		Title:    "%s Gateway Connection Uptime",
		Units:    "seconds",
		Fam:      "gw uptime",
		Ctx:      "nats.%s_gateway_conn_uptime",
		Priority: prioGatewayConnUptime,
		Dims: module.Dims{
			{ID: "gatewayz_%s_gw_%s_cid_%d_uptime", Name: "uptime"},
		},
	}
)

func (c *Collector) updateCharts() {
	maps.DeleteFunc(c.cache.accounts, func(_ string, acc *accCacheEntry) bool {
		if !acc.updated {
			c.removeAccountCharts(acc)
			return true
		}
		if !acc.hasCharts {
			acc.hasCharts = true
			c.addAccountCharts(acc)
		}
		return false
	})
	maps.DeleteFunc(c.cache.routes, func(_ uint64, route *routeCacheEntry) bool {
		if !route.updated {
			c.removeRouteCharts(route)
			return true
		}
		if !route.hasCharts {
			route.hasCharts = true
			c.addRouteCharts(route)
		}
		return false
	})
	maps.DeleteFunc(c.cache.inGateways, func(_ string, igw *gwCacheEntry) bool {
		maps.DeleteFunc(igw.conns, func(_ uint64, inConn *gwConnCacheEntry) bool {
			if !inConn.updated {
				c.removeGatewayConnCharts(inConn, true)
				return true
			}
			if !inConn.hasCharts {
				inConn.hasCharts = true
				c.addGatewayConnCharts(inConn, true)
			}
			return false
		})
		return false
	})
	maps.DeleteFunc(c.cache.outGateways, func(_ string, ogw *gwCacheEntry) bool {
		maps.DeleteFunc(ogw.conns, func(_ uint64, outConn *gwConnCacheEntry) bool {
			if !outConn.updated {
				c.removeGatewayConnCharts(outConn, false)
				return true
			}
			if !outConn.hasCharts {
				outConn.hasCharts = true
				c.addGatewayConnCharts(outConn, false)
			}
			return false
		})
		return false
	})
}

func (c *Collector) addAccountCharts(acc *accCacheEntry) {
	charts := accountChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, acc.accName)
		chart.Labels = []module.Label{
			{Key: "account", Value: acc.accName},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, acc.accName)
		}
	}

	if err := c.Charts().Add(*charts...); err != nil {
		c.Warningf("failed to add charts for account %s: %s", acc.accName, err)
	}
}

func (c *Collector) removeAccountCharts(acc *accCacheEntry) {
	px := fmt.Sprintf("account_%s_", acc.accName)
	c.removeCharts(px)
}

func (c *Collector) addRouteCharts(route *routeCacheEntry) {
	charts := routeChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, route.rid)
		chart.Labels = []module.Label{
			{Key: "route_id", Value: strconv.FormatUint(route.rid, 10)},
			{Key: "remote_id", Value: route.remoteId},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, route.rid)
		}
	}

	if err := c.Charts().Add(*charts...); err != nil {
		c.Warningf("failed to add charts for route id %d: %s", route.rid, err)
	}
}

func (c *Collector) removeRouteCharts(route *routeCacheEntry) {
	px := fmt.Sprintf("route_%d_", route.rid)
	c.removeCharts(px)
}

func (c *Collector) addGatewayConnCharts(gwConn *gwConnCacheEntry, isInbound bool) {
	direction := "outbound"
	if isInbound {
		direction = "inbound"
	}

	charts := gatewayConnChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, direction, gwConn.rgwName, gwConn.cid)
		chart.Title = fmt.Sprintf(chart.Title, cases.Title(language.English, cases.Compact).String(direction))
		chart.Ctx = fmt.Sprintf(chart.Ctx, direction)
		chart.Labels = []module.Label{
			{Key: "gateway", Value: gwConn.gwName},
			{Key: "remote_gateway", Value: gwConn.rgwName},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, direction, gwConn.rgwName, gwConn.cid)
		}
	}

	if err := c.Charts().Add(*charts...); err != nil {
		c.Warningf("failed to add charts for gateway %s %s %d: %s", direction, gwConn.rgwName, gwConn.cid, err)
	}
}

func (c *Collector) removeGatewayConnCharts(gwConn *gwConnCacheEntry, isInbound bool) {
	direction := "outbound"
	if isInbound {
		direction = "inbound"
	}
	px := fmt.Sprintf("%s_gw_%s_cid_%d_", direction, gwConn.rgwName, gwConn.cid)
	c.removeCharts(px)
}

func (c *Collector) removeCharts(prefix string) {
	for _, chart := range *c.Charts() {
		if strings.HasPrefix(chart.ID, prefix) {
			chart.MarkRemove()
			chart.MarkNotCreated()
		}
	}
}