# cockroach-test-drive
Goal: demo a sample golang service using a cockroach cluster, testing and monitoring rolling database restarts.

Steps:
- [ ] docker-compose cluster with containers:
  - [x] 3 identical cockroach nodes with persistent volumes on host
  - [ ] golang microservice using `lib/pq` doing continuous reads and writes
  - [x] prometheus monitoring database and service nodes
  - [x] grafana graphing everything
- [ ] control mechanism to stop/start db nodes on demand
  - [ ] sometimes regaining persistent volumes
  - [ ] sometimes losing persistent volumes
  - [ ] measure recovery time in both cases
  - [ ] measure application reaction

## Questions

### How easy is it to write unit/integration tests against Cockroach?

Check out github.com/cockroachdb/cockroach-go/testserver

### How do we make the Postgres client aware of multiple servers?

- https://www.cockroachlabs.com/docs/stable/manual-deployment.html#step-5-set-up-haproxy-load-balancers
- https://github.com/cockroachdb/docs/issues/978
