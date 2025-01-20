package bridgeservice

import (
	"context"
	"fmt"
	"github.com/orb-community/orb/sinker/redis/producer"
	sinkspb "github.com/orb-community/orb/sinks/pb"
	"sort"
	"time"

	"github.com/go-kit/kit/metrics"
	fleetpb "github.com/orb-community/orb/fleet/pb"
	policiespb "github.com/orb-community/orb/policies/pb"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type BridgeService interface {
	ExtractAgent(ctx context.Context, channelID string) (*fleetpb.AgentInfoRes, error)
	GetPolicyName(ctx context.Context, policyId, ownerId string) (*policiespb.PolicyRes, error)
	GetDataSetsFromAgentGroups(ctx context.Context, mfOwnerId string, agentGroupIds []string) (map[string]string, error)
	NotifyActiveSink(ctx context.Context, mfOwnerId, sinkId, state, message string) error
	GetSinkIdsFromPolicyID(ctx context.Context, mfOwnerId string, policyID string) (map[string]string, error)
	IncrementMessageCounter(publisher, subtopic, channel, protocol string)
}

func NewBridgeService(logger *zap.Logger,
	defaultCacheExpiration time.Duration,
	sinkActivity producer.SinkActivityProducer,
	policiesClient policiespb.PolicyServiceClient,
	sinksClient sinkspb.SinkServiceClient,
	fleetClient fleetpb.FleetServiceClient, messageInputCounter metrics.Counter) SinkerOtelBridgeService {
	return SinkerOtelBridgeService{
		defaultCacheExpiration: defaultCacheExpiration,
		inMemoryCache:          *cache.New(defaultCacheExpiration, defaultCacheExpiration*2),
		logger:                 logger,
		sinkerActivitySvc:      sinkActivity,
		policiesClient:         policiesClient,
		fleetClient:            fleetClient,
		sinksClient:            sinksClient,
		messageInputCounter:    messageInputCounter,
	}
}

type SinkerOtelBridgeService struct {
	inMemoryCache          cache.Cache
	defaultCacheExpiration time.Duration
	logger                 *zap.Logger
	sinkerActivitySvc      producer.SinkActivityProducer
	policiesClient         policiespb.PolicyServiceClient
	fleetClient            fleetpb.FleetServiceClient
	sinksClient            sinkspb.SinkServiceClient
	messageInputCounter    metrics.Counter
}

// IncrementMessageCounter add to our metrics the number of messages received
func (bs *SinkerOtelBridgeService) IncrementMessageCounter(publisher, subtopic, channel, protocol string) {
	labels := []string{
		"method", "handleMsgFromAgent",
		"agent_id", publisher,
		"subtopic", subtopic,
		"channel", channel,
		"protocol", protocol,
	}
	bs.messageInputCounter.With(labels...).Add(1)
}

// NotifyActiveSink notify the sinker that a sink is active
func (bs *SinkerOtelBridgeService) NotifyActiveSink(ctx context.Context, mfOwnerId, sinkId, size string) error {
	cacheKey := fmt.Sprintf("active_sink-%s-%s", mfOwnerId, sinkId)
	_, found := bs.inMemoryCache.Get(cacheKey)
	if !found {
		bs.logger.Debug("notifying active sink", zap.String("sink_id", sinkId), zap.String("owner_id", mfOwnerId),
			zap.String("payload_size", size))
		event := producer.SinkActivityEvent{
			OwnerID:   mfOwnerId,
			SinkID:    sinkId,
			State:     "active",
			Size:      size,
			Timestamp: time.Now(),
		}
		err := bs.sinkerActivitySvc.PublishSinkActivity(ctx, event)
		if err != nil {
			bs.logger.Error("error publishing sink activity", zap.Error(err))
		}
		bs.inMemoryCache.Set(cacheKey, true, cache.DefaultExpiration)
	} else {
		bs.logger.Debug("active sink already notified", zap.String("sink_id", sinkId), zap.String("owner_id", mfOwnerId),
			zap.String("payload_size", size))
	}

	return nil
}

// ExtractAgent retrieve agent info from fleet, or cache
func (bs *SinkerOtelBridgeService) ExtractAgent(ctx context.Context, channelID string) (*fleetpb.AgentInfoRes, error) {
	cacheKey := fmt.Sprintf("agent-%s", channelID)
	value, found := bs.inMemoryCache.Get(cacheKey)
	if !found {
		agentPb, err := bs.fleetClient.RetrieveAgentInfoByChannelID(ctx, &fleetpb.AgentInfoByChannelIDReq{Channel: channelID})
		if err != nil {
			return nil, err
		}
		bs.inMemoryCache.Set(cacheKey, agentPb, cache.DefaultExpiration)
		return agentPb, nil
	}
	return value.(*fleetpb.AgentInfoRes), nil
}

// GetPolicyName retrieve policy info from policies service, or cache.
func (bs *SinkerOtelBridgeService) GetPolicyName(ctx context.Context, policyId, ownerID string) (*policiespb.PolicyRes, error) {
	cacheKey := fmt.Sprintf("policy-%s", policyId)
	value, found := bs.inMemoryCache.Get(cacheKey)
	if !found {
		policyPb, err := bs.policiesClient.RetrievePolicy(ctx, &policiespb.PolicyByIDReq{PolicyID: policyId, OwnerID: ownerID})
		if err != nil {
			return nil, err
		}
		bs.inMemoryCache.Set(cacheKey, policyPb, cache.DefaultExpiration)
		return policyPb, nil
	}
	return value.(*policiespb.PolicyRes), nil
}

// GetSinkIdsFromDatasetIDs retrieve sink_ids from datasets from policies service, or cache
func (bs *SinkerOtelBridgeService) GetSinkIdsFromDatasetIDs(ctx context.Context, mfOwnerId string, datasetIDs []string) (map[string]string, error) {
	// Here needs to retrieve datasets
	mapSinkIdPolicy := make(map[string]string)
	sort.Strings(datasetIDs)
	for i := 0; i < len(datasetIDs); i++ {
		datasetID := datasetIDs[i]
		cacheKey := fmt.Sprintf("ds-%s-%s", mfOwnerId, datasetID)
		value, found := bs.inMemoryCache.Get(cacheKey)
		if !found {
			datasetRes, err := bs.policiesClient.RetrieveDataset(ctx, &policiespb.DatasetByIDReq{
				DatasetID: datasetID,
				OwnerID:   mfOwnerId,
			})
			if err != nil {
				bs.logger.Info("unable to retrieve datasets from policy")
				return nil, err
			}
			value = datasetRes.SinkIds
			bs.inMemoryCache.Set(cacheKey, value, cache.DefaultExpiration)
		}
		for _, sinkId := range value.([]string) {
			mapSinkIdPolicy[sinkId] = "active"
		}
	}
	return mapSinkIdPolicy, nil
}
