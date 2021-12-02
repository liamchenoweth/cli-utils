// Copyright 2020 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package statusreaders

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cli-utils/pkg/kstatus/polling/engine"
	"sigs.k8s.io/cli-utils/pkg/kstatus/polling/event"
)

func NewReplicaSetStatusReader(mapper meta.RESTMapper, podStatusReader resourceTypeStatusReader) engine.StatusReader {
	return &baseStatusReader{
		mapper: mapper,
		resourceStatusReader: &replicaSetStatusReader{
			mapper:          mapper,
			podStatusReader: podStatusReader,
		},
	}
}

// replicaSetStatusReader is an engine that can fetch ReplicaSet resources
// from the cluster, knows how to find any Pods belonging to the ReplicaSet,
// and compute status for the ReplicaSet.
type replicaSetStatusReader struct {
	mapper meta.RESTMapper

	podStatusReader resourceTypeStatusReader
}

var _ resourceTypeStatusReader = &replicaSetStatusReader{}

func (r *replicaSetStatusReader) ReadStatusForObject(ctx context.Context, reader engine.ClusterReader, rs *unstructured.Unstructured) *event.ResourceStatus {
	return newPodControllerStatusReader(r.mapper, r.podStatusReader).readStatus(ctx, reader, rs)
}
