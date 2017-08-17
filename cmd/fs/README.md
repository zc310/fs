# proxy 

```
hey -c 100 -n 10000 http://localhost:8080/api
All requests done.

Summary:
  Total:	0.3076 secs
  Slowest:	0.0182 secs
  Fastest:	0.0001 secs
  Average:	0.0029 secs
  Requests/sec:	32505.2056
  Total data:	140000 bytes
  Size/request:	14 bytes

Status code distribution:
  [200]	10000 responses

Response time histogram:
  0.000 [1]	|
  0.002 [3023]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.004 [4562]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.006 [1628]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.007 [527]	|∎∎∎∎∎
  0.009 [173]	|∎∎
  0.011 [62]	|∎
  0.013 [15]	|
  0.015 [6]	|
  0.016 [2]	|
  0.018 [1]	|

Latency distribution:
  10% in 0.0012 secs
  25% in 0.0018 secs
  50% in 0.0025 secs
  75% in 0.0037 secs
  90% in 0.0052 secs
  95% in 0.0061 secs
  99% in 0.0089 secs

  ```