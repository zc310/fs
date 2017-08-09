# HelloWorld server example

* Displays Hello, world!.
* delay.



# How to run

```
./go run main.go -addr=:7070
```


# hey test

```
hey -c 1000  -n 100000 http://127.0.0.1:7070/
20648 requests done.
50909 requests done.
80757 requests done.
All requests done.

Summary:
  Total:	1.8158 secs
  Slowest:	0.1645 secs
  Fastest:	0.0001 secs
  Average:	0.0175 secs
  Requests/sec:	55072.1283
  Total data:	1400000 bytes
  Size/request:	14 bytes

Status code distribution:
  [200]	100000 responses

Response time histogram:
  0.000 [1]	|
  0.017 [67796]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.033 [29922]	|∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  0.049 [1382]	|∎
  0.066 [182]	|
  0.082 [22]	|
  0.099 [156]	|
  0.115 [185]	|
  0.132 [298]	|
  0.148 [43]	|
  0.164 [13]	|

Latency distribution:
  10% in 0.0118 secs
  25% in 0.0133 secs
  50% in 0.0148 secs
  75% in 0.0185 secs
  90% in 0.0265 secs
  95% in 0.0294 secs
  99% in 0.0429 secs
```