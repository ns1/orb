// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Adapted for Orb project, modifications licensed under MPL v. 2.0:
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package grpc

import (
	"context"
	"github.com/orb-community/orb/policies"
	"github.com/orb-community/orb/policies/pb"
	"go.uber.org/zap"

	kitot "github.com/go-kit/kit/tracing/opentracing"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.PolicyServiceServer = (*grpcServer)(nil)

type grpcServer struct {
	pb.UnimplementedPolicyServiceServer
	retrievePolicy           kitgrpc.Handler
	retrievePoliciesByGroups kitgrpc.Handler
	retrieveDataset          kitgrpc.Handler
	retrieveDatasetsByGroups kitgrpc.Handler
}

// NewServer returns new PolicyServiceServer instance.
func NewServer(tracer opentracing.Tracer, svc policies.Service) pb.PolicyServiceServer {
	return &grpcServer{
		retrievePolicy: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_policy")(retrievePolicyEndpoint(svc)),
			decodeRetrievePolicyRequest,
			encodePolicyResponse,
		),
		retrievePoliciesByGroups: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_policies_by_groups")(retrievePoliciesByGroupsEndpoint(svc)),
			decodeRetrievePoliciesByGroupRequest,
			encodePolicyInDSListResponse,
		),
		retrieveDataset: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_dataset")(retrieveDatasetEnpoint(svc)),
			decodeRetrieveDatasetRequest,
			encodeDatasetResponse,
		),
		retrieveDatasetsByGroups: kitgrpc.NewServer(
			kitot.TraceServer(tracer, "retrieve_datasets_by_groups")(retrieveDatasetsByGroupsEndpoint(svc)),
			decodeRetrieveDatasetsByGroupRequest,
			encodeDatasetListResponse,
		),
	}
}

func (gs *grpcServer) RetrievePoliciesByGroups(ctx context.Context, req *pb.PoliciesByGroupsReq) (*pb.PolicyInDSListRes, error) {
	_, res, err := gs.retrievePoliciesByGroups.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.PolicyInDSListRes), nil
}

func (gs *grpcServer) RetrievePolicy(ctx context.Context, req *pb.PolicyByIDReq) (*pb.PolicyRes, error) {
	_, res, err := gs.retrievePolicy.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.PolicyRes), nil
}

func (gs *grpcServer) RetrieveDataset(ctx context.Context, req *pb.DatasetByIDReq) (*pb.DatasetRes, error) {
	_, res, err := gs.retrieveDataset.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.DatasetRes), nil
}

func (gs *grpcServer) RetrieveDatasetsByGroups(ctx context.Context, req *pb.DatasetsByGroupsReq) (*pb.DatasetsRes, error) {
	_, res, err := gs.retrieveDatasetsByGroups.ServeGRPC(ctx, req)
	if err != nil {
		return nil, encodeError(err)
	}

	return res.(*pb.DatasetsRes), nil
}

func decodeRetrievePolicyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.PolicyByIDReq)
	return accessByIDReq{PolicyID: req.PolicyID, OwnerID: req.OwnerID}, nil
}

func decodeRetrievePoliciesByGroupRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.PoliciesByGroupsReq)
	return accessByGroupIDReq{GroupIDs: req.GroupIDs, OwnerID: req.OwnerID}, nil
}

func decodeRetrieveDatasetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DatasetByIDReq)
	return accessDatasetByIDReq{
		datasetID: req.DatasetID,
		ownerID:   req.OwnerID,
	}, nil
}

func decodeRetrieveDatasetsByGroupRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DatasetsByGroupsReq)
	return accessByGroupIDReq{GroupIDs: req.GroupIDs, OwnerID: req.OwnerID}, nil
}

func encodePolicyResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(policyRes)
	return &pb.PolicyRes{
		Id:      res.id,
		Name:    res.name,
		Backend: res.backend,
		Version: res.version,
		Data:    res.data,
		Format:  res.format,
	}, nil
}

func encodePolicyInDSListResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(policyInDSListRes)
	plist := make([]*pb.PolicyInDSRes, len(res.policies))
	l, _ := zap.NewDevelopment()
	for i, p := range res.policies {
		l.Debug("policy format", zap.String("format", p.format), zap.String("policy_id", p.id))
		plist[i] = &pb.PolicyInDSRes{Id: p.id,
			Name:         p.name,
			Data:         p.data,
			Backend:      p.backend,
			Version:      p.version,
			DatasetId:    p.datasetID,
			AgentGroupId: p.agentGroupID,
			Format:       p.format,
		}
	}
	return &pb.PolicyInDSListRes{Policies: plist}, nil
}

func encodeDatasetResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(datasetRes)

	return &pb.DatasetRes{
		Id:           res.id,
		AgentGroupId: res.agentGroupID,
		PolicyId:     res.policyID,
		SinkIds:      res.sinkIDs,
	}, nil
}

func encodeDatasetListResponse(_ context.Context, grpcRes interface{}) (interface{}, error) {
	res := grpcRes.(datasetListRes)

	dsList := make([]*pb.DatasetRes, len(res.datasets))
	for i, ds := range res.datasets {
		dsList[i] = &pb.DatasetRes{Id: ds.id, PolicyId: ds.policyID, AgentGroupId: ds.agentGroupID, SinkIds: ds.sinkIDs}
	}
	return &pb.DatasetsRes{DatasetList: dsList}, nil
}

func encodeError(err error) error {
	switch err {
	case nil:
		return nil
	case policies.ErrMalformedEntity:
		return status.Error(codes.InvalidArgument, "received invalid can access request")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
