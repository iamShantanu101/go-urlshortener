## Overview of URL shortener service

1.  Programming language: go
2.  Fast in-memory K/V store: [BBoltDB](https://github.com/coreos/bbolt)
3.  HTTP router: [httprouter](https://github.com/julienschmidt/httprouter)
4.  Logging library: [logrus](https://github.com/sirupsen/logrus)

### Inspirations
1. https://github.com/pantrif/url-shortener
2. https://github.com/pankajkhairnar/goShort
3. https://github.com/xcoulon/go-url-shortener

### REST endpoints

1.  [POST] - http://localhost:8080/url : It accepts POST form data with parameter "url" and returns json response with short URL and the original URL.
2.  [GET] - http://localhost:8080/{SHORT_CODE}/ : If SHORT_CODE is valid and found in BBoltDB request will be redirected to original URL.
3.  All other non-existent endpoints will return `404`.
4.  In case invalid json payload gets supplied, it'll return `422`.

### Routing

| Endpoint      | Method   | Payload                               | Response        | Feature                          |
| ------------- |:--------:|--------------------------------------:|-----------------| ---------------------------------|
| /url          | POST     |   {"url": "https://google.com"}       | JSON            | url validation, mandatory param  |
| /{short_code} | GET      |                                       |301 Redirection  | redirects to original URL        |
| /invalid_ep   | GET      |                                       | 404             | endpoint validation              |
| /url          | POST     |   {"uaarl": "https:///google.com"}    | 422             | Invalid payload                  |

### Design decisions made while selecting libs:

1. [logrus](https://github.com/sirupsen/logrus):
   1. Structural design of logrus makes user to think really hard about the important areas of the application where logging is absolutely required.
   2. Added benefits while doing log analysis with tools like `ELK stack`.
2. [BBoltDB](https://github.com/coreos/bbolt):
   1. For an app like URL shortener it made sense to use fast in-memory KV store instead of going for relational way.
   2. BBoltDB is the fork of BoltDB and both of them are widely used.
   3. Ease of use.

2. [httprouter](https://github.com/julienschmidt/httprouter):
   1. Lightweight high performance HTTP request router.
   2. Scales better.
   3. Best Performance. Check [Benchmarks](https://github.com/julienschmidt/go-http-routing-benchmark).

### Benchmarking:
1. Tools used: [vegeta](https://github.com/tsenart/vegeta),[ab](https://httpd.apache.org/docs/2.4/programs/ab.html)
2. Head over to `benchmark` directory in the repository for exploring the above two options with detailed guide on running them.    
