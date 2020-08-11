package cloudwatchlogs

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	cwl "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CloudWatchLogs struct {
	client            *cwl.CloudWatchLogs
	logGroupName      string
	logStreamName     string
	regionName        string
	nextSequenceToken string
}

func NewCloudWatchLogs(logGroupName, logStreamName, regionName string) *CloudWatchLogs {
	mySession := session.Must(session.NewSession())
	client := cloudwatchlogs.New(mySession, aws.NewConfig().WithRegion(regionName))

	logGroupInput := cloudwatchlogs.CreateLogGroupInput{LogGroupName: &logGroupName}
	client.CreateLogGroup(&logGroupInput)
	logStreamInput := cloudwatchlogs.CreateLogStreamInput{LogGroupName: &logGroupName, LogStreamName: &logStreamName}
	_, err := client.CreateLogStream(&logStreamInput)

	token := ""
	if err != nil {
		streams, err := client.DescribeLogStreams(&cwl.DescribeLogStreamsInput{
			LogGroupName:        aws.String(logGroupName),
			Descending:          aws.Bool(true),
			LogStreamNamePrefix: aws.String(logStreamName),
		})

		if err != nil {
			log.Fatal(err)
		}
		for _, stream := range streams.LogStreams {
			if *stream.LogStreamName == logStreamName {
				token = *stream.UploadSequenceToken
				break
			}
		}
	}

	return &CloudWatchLogs{
		client:            client,
		logGroupName:      logGroupName,
		logStreamName:     logStreamName,
		regionName:        regionName,
		nextSequenceToken: token,
	}
}

func (c *CloudWatchLogs) PutLogEvents(event *v1.Event) error {
	logevents := make([]*cloudwatchlogs.InputLogEvent, 0)

	sample_json, _ := json.Marshal(event)
	logevents = append(logevents, &cloudwatchlogs.InputLogEvent{
		Message:   aws.String(string(sample_json)),
		Timestamp: aws.Int64(int64(translateTimestamp(event.LastTimestamp))),
	})

	var p cloudwatchlogs.PutLogEventsInput
	if len(c.nextSequenceToken) == 0 {
		p = cloudwatchlogs.PutLogEventsInput{
			LogEvents:     logevents,
			LogGroupName:  aws.String(c.logGroupName),
			LogStreamName: aws.String(c.logStreamName)}
	} else {
		p = cloudwatchlogs.PutLogEventsInput{
			LogEvents:     logevents,
			LogGroupName:  aws.String(c.logGroupName),
			LogStreamName: aws.String(c.logStreamName),
			SequenceToken: aws.String(c.nextSequenceToken)}
	}

	resp, err := c.client.PutLogEvents(&p)
	if err != nil {
		panic(err)
	}
	if resp.NextSequenceToken != nil {
		c.nextSequenceToken = *resp.NextSequenceToken
	}
	fmt.Print("Next Token: {}", resp.NextSequenceToken)
	return err
}

func translateTimestamp(timestamp metav1.Time) int64 {
	if timestamp.IsZero() {
		return time.Now().UnixNano() / 1000000
	}

	return timestamp.UnixNano() / 1000000
}
