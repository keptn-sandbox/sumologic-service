package main

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/SumoLogic-Labs/sumologic-go-sdk/service/cip"
	"github.com/SumoLogic-Labs/sumologic-go-sdk/service/cip/types"
	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const (
	sliFile                        = "sumologic-service/sli.yaml"
	defaultSleepBeforeAPIInSeconds = 30
)

// We have to put a min of 30s of sleep for the Sumo Logic API to reflect the data correctly
var sleepBeforeAPIInSeconds int

func init() {
	var err error
	sleepBeforeAPIInSeconds, err = strconv.Atoi(strings.TrimSpace(os.Getenv("SLEEP_BEFORE_API_IN_SECONDS")))
	if err != nil || sleepBeforeAPIInSeconds < defaultSleepBeforeAPIInSeconds {
		log.Infof("defaulting SLEEP_BEFORE_API_IN_SECONDS to 30s because it was set to '%v' which is less than the min allowed value of 30s", sleepBeforeAPIInSeconds)
		sleepBeforeAPIInSeconds = defaultSleepBeforeAPIInSeconds
	}
}

/**
* Here are all the handler functions for the individual event
* See https://github.com/keptn/spec/blob/0.8.0-alpha/cloudevents.md for details on the payload
**/

// GenericLogKeptnCloudEventHandler is a generic handler for Keptn Cloud Events that logs the CloudEvent
func GenericLogKeptnCloudEventHandler(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data interface{}) error {
	log.Printf("Handling %s Event: %s", incomingEvent.Type(), incomingEvent.Context.GetID())
	log.Printf("CloudEvent %T: %v", data, data)

	return nil
}

// OldHandleConfigureMonitoringEvent handles old configure-monitoring events
// TODO: add in your handler code
func OldHandleConfigureMonitoringEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptn.ConfigureMonitoringEventData) error {
	log.Printf("Handling old configure-monitoring Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleConfigureMonitoringTriggeredEvent handles configure-monitoring.triggered events
// TODO: add in your handler code
func HandleConfigureMonitoringTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ConfigureMonitoringTriggeredEventData) error {
	log.Printf("Handling configure-monitoring.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleDeploymentTriggeredEvent handles deployment.triggered events
// TODO: add in your handler code
func HandleDeploymentTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.DeploymentTriggeredEventData) error {
	log.Printf("Handling deployment.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleTestTriggeredEvent handles test.triggered events
// TODO: add in your handler code
func HandleTestTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.TestTriggeredEventData) error {
	log.Printf("Handling test.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleApprovalTriggeredEvent handles approval.triggered events
// TODO: add in your handler code
func HandleApprovalTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ApprovalTriggeredEventData) error {
	log.Printf("Handling approval.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleEvaluationTriggeredEvent handles evaluation.triggered events
// TODO: add in your handler code
func HandleEvaluationTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.EvaluationTriggeredEventData) error {
	log.Printf("Handling evaluation.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleReleaseTriggeredEvent handles release.triggered events
// TODO: add in your handler code
func HandleReleaseTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ReleaseTriggeredEventData) error {
	log.Printf("Handling release.triggered Event: %s", incomingEvent.Context.GetID())

	return nil
}

// HandleGetSliTriggeredEvent handles get-sli.triggered events if SLIProvider == sumologic-service
// This function acts as an example showing how to handle get-sli events by sending .started and .finished events
// TODO: adapt handler code to your needs
func HandleGetSliTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.GetSLITriggeredEventData) error {
	log.Printf("Handling get-sli.triggered Event: %s", incomingEvent.Context.GetID())

	// Step 1 - Do we need to do something?
	// Lets make sure we are only processing an event that really belongs to our SLI Provider
	if data.GetSLI.SLIProvider != "sumologic" {
		log.Printf("Not handling get-sli event as it is meant for %s", data.GetSLI.SLIProvider)
		return nil
	}

	// Step 2 - Send out a get-sli.started CloudEvent
	// The get-sli.started cloud-event is new since Keptn 0.8.0 and is required to be send when the task is started
	_, err := myKeptn.SendTaskStartedEvent(data, ServiceName)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task started CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	start, err := parseUnixTimestamp(data.GetSLI.Start)
	if err != nil {
		log.Error("unable to parse sli start timestamp: %v", err)
		return err
	}
	end, err := parseUnixTimestamp(data.GetSLI.End)
	if err != nil {
		log.Errorf("unable to parse sli end timestamp: %v", err)
		return err
	}

	// Step 4 - prep-work
	// Get any additional input / configuration data
	// - Labels: get the incoming labels for potential config data and use it to pass more labels on result, e.g: links
	// - SLI.yaml: if your service uses SLI.yaml to store query definitions for SLIs get that file from Keptn
	labels := data.Labels
	if labels == nil {
		labels = make(map[string]string)
	}

	// Step 5 - get SLI Config File
	// Get SLI File from sumologic-service subdirectory of the config repo - to add the file use:
	//   keptn add-resource --project=PROJECT --stage=STAGE --service=SERVICE --resource=my-sli-config.yaml  --resourceUri=sumologic-service/sli.yaml
	sliConfig, err := myKeptn.GetSLIConfiguration(data.Project, data.Stage, data.Service, sliFile)
	log.Debugf("SLI config: %v", sliConfig)

	// FYI you do not need to "fail" if sli.yaml is missing, you can also assume smart defaults like we do
	// in keptn-contrib/dynatrace-service and keptn-contrib/prometheus-service
	if err != nil {
		// failed to fetch sli config file
		errMsg := fmt.Sprintf("Failed to fetch SLI file %s from config repo: %s", sliFile, err.Error())
		log.Error(errMsg)
		// send a get-sli.finished event with status=error and result=failed back to Keptn

		_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status: keptnv2.StatusErrored,
			Result: keptnv2.ResultFailed,
			Labels: labels,
		}, ServiceName)

		return err
	}

	// Step 6 - do your work - iterate through the list of requested indicators and return their values
	// Indicators: this is the list of indicators as requested in the SLO.yaml
	// SLIResult: this is the array that will receive the results
	indicators := data.GetSLI.Indicators
	sliResults := []*keptnv2.SLIResult{}

	client := cip.APIClient{
		Cfg: &cip.Configuration{
			Authentication: cip.BasicAuth{
				AccessId:  env.AccessId,
				AccessKey: env.AccessKey,
			},
			BasePath:   env.SumoEndPt,
			HTTPClient: &http.Client{},
		},
	}

	// default values
	getSliFinishedEventData := &keptnv2.GetSLIFinishedEventData{
		EventData: keptnv2.EventData{
			Status: keptnv2.StatusSucceeded,
			Result: keptnv2.ResultPass,
		},
		GetSLI: keptnv2.GetSLIFinished{
			Start: data.GetSLI.Start,
			End:   data.GetSLI.End,
		},
	}
	var sliResult *keptnv2.SLIResult

	for _, indicatorName := range indicators {
		// Pulling the data from Sumo Logic api immediately gives incorrect data in the api response
		// we have to wait for some time for the correct data to be reflected in the api response
		log.Debugf("waiting for %vs so that the metrics data is reflected correctly in the api", sleepBeforeAPIInSeconds)
		time.Sleep(time.Second * time.Duration(sleepBeforeAPIInSeconds))
		query := replaceQueryParameters(data, sliConfig[indicatorName], start, end)
		log.Debugf("actual query sent to sumologic: %v, from: %v, to: %v", query, start.Unix(), end.Unix())

		formattedQuery, quantizeDuration, quantizeRollup, err := processQuery(query)
		if err != nil {
			log.Error(err)
			return err
		}

		// It takes some time until the metrics
		// start reflecting in the SumoLogic API results
		time.Sleep(time.Second * 30)

		req := types.MetricsQueryRequest{
			Queries: []types.MetricsQueryRow{
				types.MetricsQueryRow{
					Query:        formattedQuery,
					RowId:        "A",
					Quantization: quantizeDuration,
					Rollup:       quantizeRollup,
				},
			},
			TimeRange: &types.ResolvableTimeRange{
				Type_: "BeginBoundedTimeRange",
				From: types.TimeRangeBoundary{
					Type:        "EpochTimeRangeBoundary",
					EpochMillis: start.UnixMilli(),
					RangeName:   "from",
				},
				To: types.TimeRangeBoundary{
					Type:        "EpochTimeRangeBoundary",
					EpochMillis: end.UnixMilli(),
					RangeName:   "to",
				},
			},
		}
		mRes, hRes, err := client.RunMetricsQueries(req)
		if err != nil {
			log.Debugf("metrics query response: %v", mRes)
			log.Debugf("http response: %v", *hRes)
			log.Error(err)
			getSliFinishedEventData.EventData.Status = keptnv2.StatusErrored
			getSliFinishedEventData.EventData.Result = keptnv2.ResultFailed
		} else {
			sliResult = &keptnv2.SLIResult{
				Metric: indicatorName,
				Value:  mRes.QueryResult[0].TimeSeriesList.TimeSeries[0].Points.Values[0],
			}
			sliResults = append(sliResults, sliResult)
		}

	}

	getSliFinishedEventData.GetSLI.IndicatorValues = sliResults

	_, err = myKeptn.SendTaskFinishedEvent(getSliFinishedEventData, ServiceName)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task finished CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	return nil
}

// HandleProblemEvent handles two problem events:
// - ProblemOpenEventType = "sh.keptn.event.problem.open"
// - ProblemEventType = "sh.keptn.events.problem"
// TODO: add in your handler code
func HandleProblemEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptn.ProblemEventData) error {
	log.Printf("Handling Problem Event: %s", incomingEvent.Context.GetID())

	// Deprecated since Keptn 0.7.0 - use the HandleActionTriggeredEvent instead

	return nil
}

// HandleActionTriggeredEvent handles action.triggered events
// TODO: add in your handler code
func HandleActionTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.ActionTriggeredEventData) error {
	log.Printf("Handling Action Triggered Event: %s", incomingEvent.Context.GetID())
	log.Printf("Action=%s\n", data.Action.Action)

	// check if action is supported
	if data.Action.Action == "action-xyz" {
		// -----------------------------------------------------
		// 1. Send Action.Started Cloud-Event
		// -----------------------------------------------------
		myKeptn.SendTaskStartedEvent(data, ServiceName)

		// -----------------------------------------------------
		// 2. Implement your remediation action here
		// -----------------------------------------------------
		time.Sleep(5 * time.Second) // Example: Wait 5 seconds. Maybe the problem fixes itself.

		// -----------------------------------------------------
		// 3. Send Action.Finished Cloud-Event
		// -----------------------------------------------------
		myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status:  keptnv2.StatusSucceeded, // alternative: keptnv2.StatusErrored
			Result:  keptnv2.ResultPass,      // alternative: keptnv2.ResultFailed
			Message: "Successfully sleeped!",
		}, ServiceName)

	} else {
		log.Printf("Retrieved unknown action %s, skipping...", data.Action.Action)
		return nil
	}
	return nil
}

func parseUnixTimestamp(timestamp string) (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err == nil {
		return parsedTime, nil
	}

	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	unix := time.Unix(timestampInt, 0)
	return unix, nil
}

func replaceQueryParameters(data *keptnv2.GetSLITriggeredEventData, query string, start, end time.Time) string {
	query = strings.Replace(query, "$PROJECT", data.Project, -1)
	query = strings.Replace(query, "$STAGE", data.Stage, -1)
	query = strings.Replace(query, "$SERVICE", data.Service, -1)
	query = strings.Replace(query, "$project", data.Project, -1)
	query = strings.Replace(query, "$stage", data.Stage, -1)
	query = strings.Replace(query, "$service", data.Service, -1)
	durationString := strconv.FormatInt(getDurationInSeconds(start, end), 10)
	query = strings.Replace(query, "$DURATION", durationString, -1)
	return query
}

func getDurationInSeconds(start, end time.Time) int64 {

	seconds := end.Sub(start).Seconds()
	return int64(math.Ceil(seconds))
}

// processQuery takes the query, parses it and returns the
// 1. formattedQuery after removing `quantize` operator (because it is
// 	not supported by the API <- syntactic sugar added by us)
// 2. quantizeVal which is in milliseconds (int64)
// 3. rollup which is either of Avg, Min, Max, Sum, or Count.
// Check https://help.sumologic.com/Metrics/Metric-Queries-and-Alerts/07Metrics_Operators/quantize#quantize-syntax
// for more info
func processQuery(query string) (string, int64, string, error) {

	// Ensure there is only one `quantize` in the query
	re := regexp.MustCompile(`quantize`)
	matches := re.FindAllString(query, -1)
	if len(matches) != 1 {
		return "", 0, "", errors.New("please specify 1 `quantize` in the query")
	}

	// Parse the quantize duration and Roll up type (e.g., avg, sum) from the query
	qRe := regexp.MustCompile(`quantize\s+to\s+(\d+[a-z])\s+using\s+([a-z]+)\s*\|?`)
	quantizePart := qRe.FindAllString(query, -1)[0]
	if len(quantizePart) == 0 {
		return "", 0, "", errors.New(fmt.Sprintf("`quantize` part of the query should match the regex `%s`", qRe.String()))
	}

	// Output of FindAllStringSubmatch is of the form [["actual", "result", "goes", "here"]]
	submatches := qRe.FindAllStringSubmatch(quantizePart, -1)[0]

	query = strings.ReplaceAll(query, quantizePart, " ")

	quantizeVal, err := time.ParseDuration(submatches[1])
	if err != nil {
		return "", 0, "", errors.New("couldn't parse the value for quantize interval/duration")
	}

	formattedQuery := strings.TrimSpace(query)

	if formattedQuery[len(formattedQuery)-1] == '|' {
		formattedQuery = formattedQuery[:len(formattedQuery)-1]
	}

	c := cases.Title(language.English)

	return formattedQuery, quantizeVal.Milliseconds(), c.String(submatches[2]), nil
}
