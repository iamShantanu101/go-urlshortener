## Simple load testing with Apache Bench (ab)

### Prerequisites
1. [ab](https://httpd.apache.org/docs/2.4/programs/ab.html)

### Examples:

```
ab -k -c 35 -n 2000 http://url_shortener_ip/a
```

This will send 2000 requests/sec with 35 concurrent connections.
