Unfortunately, it seems quite difficult to get grafana to set up
data sources and dashboards from configuration.  Instead it wants
to treat them as data, which means they must be imperatively imported.

Since this is annoying to automate, for now, do this as a manual step,
following the instructions
[here](https://www.cockroachlabs.com/docs/stable/monitor-cockroachdb-with-prometheus.html#step-4-visualize-metrics-in-grafana).
`docker-compose` is set up to persist grafana's data in our local
`./data/` directory, so this configuration persists as long as you
preserve that.

For convenience, the dashboard definitions have been copied here from
[the originals](https://github.com/cockroachdb/cockroach/tree/db5496ba6e48b6bec9bba9bf5013e2e9455244be/monitoring/grafana-dashboards).
Perhaps they should be periodically synced.
