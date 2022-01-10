package k8sresources

import (
	"kubectlfzf/pkg/util"

	appsv1 "k8s.io/api/apps/v1"
)

const StatefulSetHeader = "Cluster Namespace Name Replicas Selector Age Labels\n"

// StatefulSet is the summary of a kubernetes statefulset
type StatefulSet struct {
	ResourceMeta
	currentReplicas int
	replicas        int
	selectors       []string
}

// NewStatefulSetFromRuntime builds a k8sresource from informer result
func NewStatefulSetFromRuntime(obj interface{}, config CtorConfig) K8sResource {
	p := &StatefulSet{}
	p.FromRuntime(obj, config)
	return p
}

// FromRuntime builds object from the informer's result
func (s *StatefulSet) FromRuntime(obj interface{}, config CtorConfig) {
	statefulset := obj.(*appsv1.StatefulSet)
	s.FromObjectMeta(statefulset.ObjectMeta, config)
	s.currentReplicas = int(statefulset.Status.CurrentReplicas)
	s.replicas = int(statefulset.Status.Replicas)
	s.selectors = util.JoinStringMap(statefulset.Spec.Selector.MatchLabels, ExcludedLabels, "=")
}

// HasChanged returns true if the resource's dump needs to be updated
func (s *StatefulSet) HasChanged(k K8sResource) bool {
	oldSts := k.(*StatefulSet)
	return (s.currentReplicas != oldSts.currentReplicas ||
		s.replicas != oldSts.replicas ||
		util.StringSlicesEqual(s.selectors, oldSts.selectors) ||
		util.StringMapsEqual(s.Labels, oldSts.Labels))
}
