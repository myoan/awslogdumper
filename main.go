package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/pingcap/errors"
)

type LogConfig struct {
	LogGroupName  string
	LogStreamName string
	StartTime     int64
	EndTime       int64
}

func GetStreamInfo(grp, strm string) (int64, int64, error) {
	result, err := client.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(grp),
		LogStreamNamePrefix: aws.String(strm),
	})
	if err != nil {
		return 0, 0, err
	}
	if len(result.LogStreams) == 0 {
		return 0, 0, errors.Errorf("stream %s not found", strm)
	}

	stream := *result.LogStreams[0]
	return *stream.FirstEventTimestamp, *stream.LastEventTimestamp, nil
}

func UnixtimeMilli(t string) (int64, error) {
	tm, err := time.Parse("2006/01/02/15:04:05", t)
	if err != nil {
		return 0, err
	}
	return tm.UnixMilli(), nil
}

func NewLogConfig(grp, strm, start, end string) (*LogConfig, error) {
	starttime, endtime, err := GetStreamInfo(grp, strm)
	if err != nil {
		return nil, err
	}

	if start != "" {
		stunix, err := UnixtimeMilli(start)
		if err != nil {
			return nil, err
		}
		if starttime < stunix {
			starttime = stunix
		}
	}

	if end != "" {
		etunix, err := UnixtimeMilli(end)
		if err != nil {
			return nil, err
		}
		if etunix < endtime {
			endtime = etunix
		}
	}

	return &LogConfig{
		LogGroupName:  grp,
		LogStreamName: strm,
		StartTime:     starttime,
		EndTime:       endtime,
	}, nil
}

var client *cloudwatchlogs.CloudWatchLogs

func Run(args []string) error {
	var (
		profile   string
		loggrp    string
		logstream string
		start     string
		end       string
	)

	flag.StringVar(&profile, "profile", "", "aws profile name")
	flag.StringVar(&loggrp, "g", "", "cloudwatch log-group-name")
	flag.StringVar(&logstream, "s", "", "cloudwatch log-stream-name")
	flag.StringVar(&start, "start", "", "start time ('2006/01/02/15:04:05')")
	flag.StringVar(&end, "end", "", "end time ('2006/01/02/15:04:05')")
	flag.Parse()

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")

	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("ap-northeast-1")},
		Profile: profile,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	client = cloudwatchlogs.New(sess)

	config, err := NewLogConfig(loggrp, logstream, start, end)
	if err != nil {
		return err
	}

	log.Println(config)

	output, err := client.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(config.LogGroupName),
		LogStreamName: aws.String(config.LogStreamName),
		StartTime:     aws.Int64(config.StartTime),
		EndTime:       aws.Int64(config.EndTime),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	for _, e := range output.Events {
		fmt.Println(*e.Message)
	}

	return nil
}

func main() {
	if err := Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
