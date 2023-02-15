package bridgeservice

import (
	"context"
	fleetpb "github.com/ns1labs/orb/fleet/pb"
	policiespb "github.com/ns1labs/orb/policies/pb"
	"github.com/ns1labs/orb/sinker/config"
	"go.uber.org/zap"
	"sort"
	"time"
)

type BridgeService interface {
	ExtractAgent(ctx context.Context, channelID string) (*fleetpb.AgentInfoRes, error)
	GetDataSetsFromAgentGroups(ctx context.Context, mfOwnerId string, agentGroupIds []string) (map[string]string, error)
	NotifyActiveSink(ctx context.Context, mfOwnerId, sinkId, state, message string) error
	GetSinkIdsFromPolicyID(ctx context.Context, mfOwnerId string, policyID string) (map[string]string, error)
}

func NewBridgeService(logger *zap.Logger,
	sinkerCache config.ConfigRepo,
	policiesClient policiespb.PolicyServiceClient,
	fleetClient fleetpb.FleetServiceClient) SinkerOtelBridgeService {
	return SinkerOtelBridgeService{
		logger:         logger,
		sinkerCache:    sinkerCache,
		policiesClient: policiesClient,
		fleetClient:    fleetClient,
	}
}

type SinkerOtelBridgeService struct {
	logger         *zap.Logger
	sinkerCache    config.ConfigRepo
	policiesClient policiespb.PolicyServiceClient
	fleetClient    fleetpb.FleetServiceClient
}

func (bs *SinkerOtelBridgeService) NotifyActiveSink(ctx context.Context, mfOwnerId, sinkId, newState, message string) error {
	cfgRepo, err := bs.sinkerCache.Get(mfOwnerId, sinkId)
	if err != nil {
		bs.logger.Error("unable to retrieve the sink config", zap.Error(err))
		return err
	}
	if cfgRepo.State == config.Idle || cfgRepo.State == config.Unknown || cfgRepo.State == config.Active {
		err = cfgRepo.State.SetFromString(newState)
		if err != nil {
			bs.logger.Error("unable to set state", zap.String("new_state", newState), zap.Error(err))
			return err
		}
		err = bs.sinkerCache.DeployCollector(ctx, cfgRepo)
		if err != nil {
			bs.logger.Error("error during update sink cache", zap.String("sinkId", sinkId), zap.Error(err))
			return err
		}
		if cfgRepo.State == config.Active {
			cfgRepo.LastRemoteWrite = time.Now()
			bs.logger.Info("sink is already active, skipping")
			err = bs.sinkerCache.AddActivity(mfOwnerId, sinkId)
			if err != nil {
				bs.logger.Error("error during update last remote write", zap.String("sinkId", sinkId), zap.Error(err))
				return err
			}
		}
		bs.logger.Info("notified active sink", zap.String("sinkID", sinkId), zap.String("newState", newState))
	} else if cfgRepo.State == config.Error {
		cfgRepo.Msg = message
	}

	return nil
}

func (bs *SinkerOtelBridgeService) ExtractAgent(ctx context.Context, channelID string) (*fleetpb.AgentInfoRes, error) {
	agentPb, err := bs.fleetClient.RetrieveAgentInfoByChannelID(ctx, &fleetpb.AgentInfoByChannelIDReq{Channel: channelID})
	if err != nil {
		return nil, err
	}
	return agentPb, nil
}

func (bs *SinkerOtelBridgeService) GetSinkIdsFromDatasetIDs(ctx context.Context, mfOwnerId string, datasetIDs []string) (map[string]string, error) {
	// Here needs to retrieve datasets
	mapSinkIdPolicy := make(map[string]string)
	sort.Strings(datasetIDs)
	for i := 0; i < len(datasetIDs); i++ {
		datasetRes, err := bs.policiesClient.RetrieveDataset(ctx, &policiespb.DatasetByIDReq{
			DatasetID: datasetIDs[i],
			OwnerID:   mfOwnerId,
		})
		if err != nil {
			bs.logger.Info("unable to retrieve datasets from policy")
			return nil, err
		}
		for _, sinkId := range datasetRes.SinkIds {
			mapSinkIdPolicy[sinkId] = "active"
		}
	}
	return mapSinkIdPolicy, nil
}
