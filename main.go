package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var (
	subjectTemplate = "ALARM: %s"
	bodyTemplate    = `{
	"AlarmName": "%s",
	"AlarmDescription": "%s",
	"AWSAccountId": "%s",
	"NewStateValue": "ALARM",
	"NewStateReason": "Threshold Crossed: 1 datapoint wanted to test the waters",
	"StateChangeTime": "%s",
	"Region": "%s",
	"OldStateValue": "OK",
	"Trigger": {
		"MetricName": "Testing",
		"Namespace": "TestSpace",
		"Statistic": "SUM",
		"Unit": null,
		"Dimensions": [],
		"Period": 300,
		"EvaluationPeriods": 1,
		"ComparisonOperator": "GreaterThanOrEqualToThreshold",
		"Threshold": 1.0
	}
}`
)

func main() {
	ctx := context.Background()
	now := time.Now().UTC().Format("2006-01-02T15:04:05.000-0700")

	topicArn := ""
	alarmName := "testing-alarm"
	alarmDescription := "you can just chill!"
	verbose := false

	flag.StringVar(&topicArn, "topic", "", "arn to the sns topic to send the alarm to")
	flag.StringVar(&alarmName, "name", alarmName, "the alarm name to send")
	flag.StringVar(&alarmDescription, "description", alarmDescription, "the alarm description to send")
	flag.BoolVar(&verbose, "verbose", verbose, "print the subject and body being sent to sns")
	flag.Parse()

	partsArn := strings.Split(topicArn, ":")
	if len(partsArn) != 6 {
		flag.Usage()
		if topicArn != "" {
			fmt.Printf("\nerror: topic must be an valid arn\n")
		}
		os.Exit(2)
	}
	accountId := partsArn[4]
	region := partsArn[3]

	subject := fmt.Sprintf(subjectTemplate, alarmName)
	body := fmt.Sprintf(bodyTemplate, alarmName, alarmDescription, accountId, now, region)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}

	if verbose {
		fmt.Printf("Topic:\n%s\n\n", topicArn)
		fmt.Printf("Subject:\n%s\n\n", subject)
		fmt.Printf("Body:\n%s\n\n", body)
	}

	client := sns.NewFromConfig(cfg)
	_, err = client.Publish(ctx, &sns.PublishInput{
		TopicArn: &topicArn,
		Subject:  &subject,
		Message:  &body,
	})
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}

	fmt.Printf("alarm %q sent to %q\n", alarmName, topicArn)
}
