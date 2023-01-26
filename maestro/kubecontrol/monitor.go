package kubecontrol

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	maestroconfig "github.com/ns1labs/orb/maestro/config"
	maestroredis "github.com/ns1labs/orb/maestro/redis"
	sinkspb "github.com/ns1labs/orb/sinks/pb"
	"go.uber.org/zap"
	"io"
	k8scorev1 "k8s.io/api/core/v1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
	"time"
)

const MonitorFixedDuration = 1 * time.Minute
const TimeDiffActiveIdle = 5 * time.Minute

func NewMonitorService(logger *zap.Logger, sinksClient *sinkspb.SinkServiceClient, redisClient *redis.Client, kubecontrol *Service) MonitorService {
	return &monitorService{
		logger:      logger,
		sinksClient: *sinksClient,
		redisClient: redisClient,
		kubecontrol: *kubecontrol,
	}
}

type MonitorService interface {
	Start(ctx context.Context, cancelFunc context.CancelFunc) error
	GetRunningPods(ctx context.Context) ([]string, error)
}

type monitorService struct {
	logger      *zap.Logger
	sinksClient sinkspb.SinkServiceClient
	redisClient *redis.Client
	kubecontrol Service
}

func (svc *monitorService) Start(ctx context.Context, cancelFunc context.CancelFunc) error {
	go func(ctx context.Context, cancelFunc context.CancelFunc) {
		ticker := time.NewTicker(MonitorFixedDuration)
		svc.logger.Info("start monitor routine", zap.Any("routine", ctx))
		defer func() {
			cancelFunc()
			svc.logger.Info("stopping monitor routine")
		}()
		for {
			select {
			case <-ctx.Done():
				cancelFunc()
				return
			case _ = <-ticker.C:
				svc.logger.Info("monitoring sinks")
				svc.monitorSinks(ctx)
			}
		}
	}(ctx, cancelFunc)
	return nil
}

func (svc *monitorService) getPodLogs(ctx context.Context, pod k8scorev1.Pod) ([]string, error) {
	maxTailLines := int64(10)
	podLogOpts := k8scorev1.PodLogOptions{TailLines: &maxTailLines}
	config, err := rest.InClusterConfig()
	if err != nil {
		svc.logger.Error("error on get cluster config", zap.Error(err))
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		svc.logger.Error("error on get client", zap.Error(err))
		return nil, err
	}
	req := clientSet.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		svc.logger.Error("error on get logs", zap.Error(err))
		return nil, err
	}
	defer func(podLogs io.ReadCloser) {
		err := podLogs.Close()
		if err != nil {
			svc.logger.Error("error closing log stream", zap.Error(err))
		}
	}(podLogs)

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		svc.logger.Error("error on copying buffer", zap.Error(err))
		return nil, err
	}
	str := buf.String()
	splitLogs := strings.Split(str, "\n")
	svc.logger.Info("logs length", zap.Int("amount line logs", len(splitLogs)))
	return splitLogs, nil
}

func (svc *monitorService) GetRunningPods(ctx context.Context) ([]string, error) {
	pods, err := svc.getRunningPods(ctx)
	if err != nil {
		svc.logger.Error("error getting running collectors")
		return nil, err
	}
	runningSinks := make([]string, len(pods))
	if len(pods) > 0 {
		for i, pod := range pods {
			runningSinks[i] = strings.TrimPrefix(pod.Name, "otel-")
		}
		return runningSinks, nil
	}
	return nil, nil
}

func (svc *monitorService) getRunningPods(ctx context.Context) ([]k8scorev1.Pod, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		svc.logger.Error("error on get cluster config", zap.Error(err))
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		svc.logger.Error("error on get client", zap.Error(err))
		return nil, err
	}
	pods, err := clientSet.CoreV1().Pods(namespace).List(ctx, k8smetav1.ListOptions{})
	return pods.Items, err
}

func (svc *monitorService) monitorSinks(ctx context.Context) {
	runningCollectors, err := svc.getRunningPods(ctx)
	if err != nil {
		svc.logger.Error("error getting running pods on namespace", zap.Error(err))
		return
	}
	if len(runningCollectors) == 0 {
		svc.logger.Info("skipping, no running collectors")
		return
	}
	sinksRes, err := svc.sinksClient.RetrieveSinks(ctx, &sinkspb.SinksFilterReq{OtelEnabled: "enabled"})

	if err != nil {
		svc.logger.Error("error collecting sinks", zap.Error(err))
		return
	}
	svc.logger.Info("reading logs from collectors", zap.Int("collectors_length", len(sinksRes.Sinks)))
	for _, sink := range sinksRes.Sinks {
		var sinkCollector *k8scorev1.Pod
		for _, collector := range runningCollectors {
			if strings.Contains(collector.Name, sink.Id) {
				sinkCollector = &collector
				break
			}
		}
		if sinkCollector == nil {
			svc.logger.Warn("collector not found for sink, skipping", zap.String("sinkID", sink.Id))
			continue
		}
		var data maestroconfig.SinkData
		if err := json.Unmarshal(sink.Config, &data); err != nil {
			svc.logger.Warn("failed to unmarshal sink, skipping", zap.String("sink-id", sink.Id))
			continue
		}
		if data.LastRemoteWrite.After(time.Now().Add(-TimeDiffActiveIdle)) {
			svc.logger.Warn("collector recently updated, skipping", zap.String("sink-id", sink.Id))
			continue
		}
		data.SinkID = sink.Id
		data.OwnerID = sink.OwnerID
		data.LastRemoteWrite = time.Now()
		logs, err := svc.getPodLogs(ctx, *sinkCollector)
		if err != nil {
			svc.logger.Error("error on getting logs, skipping", zap.Error(err))
			continue
		}
		status, logsErr := svc.analyzeLogs(logs)
		if status == "fail" {
			svc.logger.Error("error during analyze logs", zap.Error(logsErr))
			continue
		}
		if data.State.String() != status {
			if err != nil {
				svc.logger.Info("updating status", zap.Any("before", sink.GetState()), zap.String("new status", status), zap.String("error_message (opt)", err.Error()))
			} else {
				svc.logger.Info("updating status", zap.Any("before", sink.GetState()), zap.String("new status", status))
			}
			svc.publishSinkStateChange(sink, status, logsErr, err)
		}
	}

}

func (svc *monitorService) publishSinkStateChange(sink *sinkspb.SinkRes, status string, logsErr error, err error) {
	streamID := "orb.sinker"
	logMessage := ""
	if logsErr != nil {
		logMessage = logsErr.Error()
	}
	event := maestroredis.SinkerUpdateEvent{
		SinkID:    sink.Id,
		Owner:     sink.OwnerID,
		State:     status,
		Msg:       logMessage,
		Timestamp: time.Now(),
	}

	record := &redis.XAddArgs{
		Stream: streamID,
		Values: event.Encode(),
	}
	err = svc.redisClient.XAdd(context.Background(), record).Err()
	if err != nil {
		svc.logger.Error("error sending event to event store", zap.Error(err))
	}
}

// analyzeLogs, will check for errors in exporter, and will return as follows
//
//		for active, the timestamp should not be longer than 5 minutes of the last metric export
//		for errors 479 will send a "warning" state, plus message of too many requests
//		for any other errors, will add error and message
//		if no error message on exporter, will log as active
//	 logs from otel-collector are coming in the standard from https://pkg.go.dev/log,
//
// TODO changing the logs from otel-collector to a json format that we can read and check for errors, will affect this
func (svc *monitorService) analyzeLogs(logEntry []string) (status string, err error) {
	var lastTimeStamp string
	for _, logLine := range logEntry {
		if len(logLine) > 24 {
			lastTimeStamp = logLine[0:24]
			if strings.Contains(logLine, "error") {
				errStringLog := strings.TrimRight(logLine, "error")
				if len(errStringLog) > 4 {
					jsonError := strings.Split(errStringLog, "\t")[4]
					errorJson := make(map[string]interface{})
					err := json.Unmarshal([]byte(jsonError), &errorJson)
					if err != nil {
						return "fail", err
					}
					if errorJson != nil && errorJson["error"] != nil {
						errorMessage := errorJson["error"].(string)
						if strings.Contains(errorMessage, "429") {
							return "warning", errors.New(errorMessage)
						} else {
							return "error", errors.New(errorMessage)
						}
					}
				}
			}
		}
	}
	lastLogTime, err := time.Parse(time.RFC3339, lastTimeStamp)
	if err != nil {
		return "fail", err
	}
	if lastLogTime.After(time.Now().Add(-TimeDiffActiveIdle)) {
		return "idle", nil
	} else {
		return "active", nil
	}
}
