package configmap

// Functions for parsing out common ConfigMaps.

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/astro-walker/tilt/pkg/apis/core/v1alpha1"
	"github.com/astro-walker/tilt/pkg/model"
)

const TriggerQueueName = "tilt-trigger-queue"

func TriggerQueue(ctx context.Context, client client.Client) (*v1alpha1.ConfigMap, error) {
	var cm v1alpha1.ConfigMap
	err := client.Get(ctx, types.NamespacedName{Name: TriggerQueueName}, &cm)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	}

	return &cm, nil
}

func NamesInTriggerQueue(cm *v1alpha1.ConfigMap) []string {
	result := make([]string, 0, len(cm.Data)/2)
	for k, v := range cm.Data {
		if !strings.HasSuffix(k, "-name") {
			continue
		}

		result = append(result, v)
	}
	return result
}

func InTriggerQueue(cm *v1alpha1.ConfigMap, nn types.NamespacedName) bool {
	name := nn.Name
	for k, v := range cm.Data {
		if !strings.HasSuffix(k, "-name") {
			continue
		}

		if v == name {
			return true
		}
	}
	return false
}

func TriggerQueueReason(cm *v1alpha1.ConfigMap, nn types.NamespacedName) model.BuildReason {
	name := nn.Name
	for k, v := range cm.Data {
		if !strings.HasSuffix(k, "-name") {
			continue
		}

		if v != name {
			continue
		}

		cur := strings.TrimSuffix(k, "-name")
		reasonCode := cm.Data[fmt.Sprintf("%s-reason-code", cur)]
		i, err := strconv.Atoi(reasonCode)
		if err != nil {
			return model.BuildReasonFlagTriggerUnknown
		}
		return model.BuildReason(i)
	}
	return model.BuildReasonNone
}

type TriggerQueueEntry struct {
	Name   model.ManifestName
	Reason model.BuildReason
}

func TriggerQueueCreate(entries []TriggerQueueEntry) v1alpha1.ConfigMap {
	cm := v1alpha1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: TriggerQueueName,
		},
		Data: make(map[string]string, len(entries)),
	}

	for i, entry := range entries {
		cm.Data[fmt.Sprintf("%d-name", i)] = entry.Name.String()
		reason := entry.Reason
		if !reason.HasTrigger() {
			reason = model.BuildReasonFlagTriggerUnknown
		}
		cm.Data[fmt.Sprintf("%d-reason-code", i)] = fmt.Sprintf("%d", reason)
	}
	return cm
}
