# Level 7 Load Balancer

Simple application level load balancer.

Can use the following algorithms:

- Round robin
- Least connections

Active health checks in the background for hosts, so requests can be distributed across active hosts.

## Demo

Lower standard deviation is better, because it means the load is being spread more evenly.

### Round Robin

![Round robin demonstration gif](./public/rr.gif)

### Least Connections

![Least connections demonstration gif](./public/lc.gif)

## Usage

```
> l7-load balancer <hosts>
starting load balancer: port=8080 algo=lc
```
