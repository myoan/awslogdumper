# awslogdumper

`awslogdumper` easily get AWS clouswatch log-events.

## install

```
go install github.com/myoan/awslogdumper
```

## how to use it

```
❯ awslogdumper -h                                                                                                                                                
Usage of awslogdumper:
  -end string
        end time ('2006/01/02/15:04:05')
  -g string
        cloudwatch log-group-name
  -profile string
        aws profile name
  -s string
        cloudwatch log-stream-name
  -start string
        start time ('2006/01/02/15:04:05')
```

When you omit `-start` option, awslogdumper gets log-stream information and uses `stream.FirstEventTimestamp`.
If your requested start-time is earlier than `stream.FirstEventTimestamp`, awslogdumper uses `stream.FirstEventTimestamp`.
When you omit `-end` option, awslogdumper gets stream info and uses `stream.LastEventTimestamp`.
If your requested end-time is later than `stream.LastEventTimestamp`, awslogdumper uses `stream.LastEventTimestamp`.