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
	flag.StringVar(&start, "start", "", "start time ('2006/01/02-15:04:05')")
	flag.StringVar(&end, "end", "", "end time ('2006/01/02-15:04:05')")
	flag.Parse()

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")

	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("ap-northeast-1")},
		Profile: profile,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	client := cloudwatchlogs.New(sess)

	starttime, err := time.Parse("2006/01/02-15:04:05", start)
	if err != nil {
		return errors.WithStack(err)
	}

	endtime, err := time.Parse("2006/01/02-15:04:05", end)
	if err != nil {
		return errors.WithStack(err)
	}

	output, err := client.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(loggrp),
		LogStreamName: aws.String(logstream),
		StartTime:     aws.Int64(starttime.UnixMilli()),
		EndTime:       aws.Int64(endtime.UnixMilli()),
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
