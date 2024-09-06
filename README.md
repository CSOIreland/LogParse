# LogParse
Parses .NET8 application log files

To build the executable type go build

To run the executable type LogParse.exe \<directory of log files\> where \<directory of log files\> is the directory where the log files are located.

This will produce an output JSON file with the date when the LogParse executable was run(e.g., results2024-01-03.json).

This file shows the errors in priority of frequency of occurance, in descending order. The log file name, the time of the first occurance of the error and the first
line of the error are also shown e.g.:

```
[
 {
  "LogFile": "log3.log",
  "Time": "2024-02-19 13:01:41,868 [8]",
  "Origin": "CSO-WEB-TD-S01",
  "CorrelationId": "e8427826041bf9fea98b39724e052bc9",
  "LineNumber": 22,
  "Error": "ERROR Enyim.Caching.Memcached.MemcachedNode+InternalPoolImpl.InitPool:0 - Failed to put PooledSocket 1 in Pool",
  "Frequency": 382
 },
 {
  "LogFile": "log5.log",
  "Time": "2024-02-19 12:55:41,290 [8]",
  "Origin": "CSO-WEB-TD-S01",
  "CorrelationId": "99477ac95486dedc7cf86d2de39eebdd",
  "LineNumber": 31,
  "Error": "ERROR Enyim.Caching.Memcached.MemcachedNode.CreateSocket:0 - Create PooledSocket",
  "Frequency": 382
 },
 {
  "LogFile": "log6.log",
  "Time": "2024-02-21 11:16:09,505 [508]",
  "Origin": "CSO-WEB-TD-S01",
  "CorrelationId": "(null)",
  "LineNumber": 27,
  "Error": "FATAL API.RESTful+\u003cProcessRequest\u003ed__3.MoveNext:0 - System.ArgumentNullException: Value cannot be null. (Parameter 'text')",
  "Frequency": 364
 },
 ...
]
```
