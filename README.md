# awslogdumper

`awslogdumper` easily get AWS clouswatch log-events.

## install

```
go install github.com/myoan/awslogdumper@latest
```

## how to use it

```
❯ ./awslogdumper -h                                                                                                                                                
Usage of ./awslogdumper:
  -end string
        end time ('2006/01/02-15:04:05')
  -g string
        cloudwatch log-group-name
  -profile string
        aws profile name
  -s string
        cloudwatch log-stream-name
  -start string
        start time ('2006/01/02-15:04:05')
```
