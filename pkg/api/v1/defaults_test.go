/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	versioned "k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/util/intstr"
)

func roundTrip(t *testing.T, obj runtime.Object) runtime.Object {
	codec := api.Codecs.LegacyCodec(metav1.SchemeGroupVersion)
	data, err := runtime.Encode(codec, obj)
	if err != nil {
		t.Errorf("%v\n %#v", err, obj)
		return nil
	}
	obj2, err := runtime.Decode(codec, data)
	if err != nil {
		t.Errorf("%v\nData: %s\nSource: %#v", err, string(data), obj)
		return nil
	}
	obj3 := reflect.New(reflect.TypeOf(obj).Elem()).Interface().(runtime.Object)
	err = api.Scheme.Convert(obj2, obj3, nil)
	if err != nil {
		t.Errorf("%v\nSource: %#v", err, obj2)
		return nil
	}
	return obj3
}

func TestSetDefaultReplicationController(t *testing.T) {
	tests := []struct {
		rc             *metav1.ReplicationController
		expectLabels   bool
		expectSelector bool
	}{
		{
			rc: &metav1.ReplicationController{
				Spec: metav1.ReplicationControllerSpec{
					Template: &metav1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			expectLabels:   true,
			expectSelector: true,
		},
		{
			rc: &metav1.ReplicationController{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"bar": "foo",
					},
				},
				Spec: metav1.ReplicationControllerSpec{
					Template: &metav1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			expectLabels:   false,
			expectSelector: true,
		},
		{
			rc: &metav1.ReplicationController{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"bar": "foo",
					},
				},
				Spec: metav1.ReplicationControllerSpec{
					Selector: map[string]string{
						"some": "other",
					},
					Template: &metav1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			expectLabels:   false,
			expectSelector: false,
		},
		{
			rc: &metav1.ReplicationController{
				Spec: metav1.ReplicationControllerSpec{
					Selector: map[string]string{
						"some": "other",
					},
					Template: &metav1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			expectLabels:   true,
			expectSelector: false,
		},
	}

	for _, test := range tests {
		rc := test.rc
		obj2 := roundTrip(t, runtime.Object(rc))
		rc2, ok := obj2.(*metav1.ReplicationController)
		if !ok {
			t.Errorf("unexpected object: %v", rc2)
			t.FailNow()
		}
		if test.expectSelector != reflect.DeepEqual(rc2.Spec.Selector, rc2.Spec.Template.Labels) {
			if test.expectSelector {
				t.Errorf("expected: %v, got: %v", rc2.Spec.Template.Labels, rc2.Spec.Selector)
			} else {
				t.Errorf("unexpected equality: %v", rc.Spec.Selector)
			}
		}
		if test.expectLabels != reflect.DeepEqual(rc2.Labels, rc2.Spec.Template.Labels) {
			if test.expectLabels {
				t.Errorf("expected: %v, got: %v", rc2.Spec.Template.Labels, rc2.Labels)
			} else {
				t.Errorf("unexpected equality: %v", rc.Labels)
			}
		}
	}
}

func newInt(val int32) *int32 {
	p := new(int32)
	*p = val
	return p
}

func TestSetDefaultReplicationControllerReplicas(t *testing.T) {
	tests := []struct {
		rc             metav1.ReplicationController
		expectReplicas int32
	}{
		{
			rc: metav1.ReplicationController{
				Spec: metav1.ReplicationControllerSpec{
					Template: &metav1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			expectReplicas: 1,
		},
		{
			rc: metav1.ReplicationController{
				Spec: metav1.ReplicationControllerSpec{
					Replicas: newInt(0),
					Template: &metav1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			expectReplicas: 0,
		},
		{
			rc: metav1.ReplicationController{
				Spec: metav1.ReplicationControllerSpec{
					Replicas: newInt(3),
					Template: &metav1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			expectReplicas: 3,
		},
	}

	for _, test := range tests {
		rc := &test.rc
		obj2 := roundTrip(t, runtime.Object(rc))
		rc2, ok := obj2.(*metav1.ReplicationController)
		if !ok {
			t.Errorf("unexpected object: %v", rc2)
			t.FailNow()
		}
		if rc2.Spec.Replicas == nil {
			t.Errorf("unexpected nil Replicas")
		} else if test.expectReplicas != *rc2.Spec.Replicas {
			t.Errorf("expected: %d replicas, got: %d", test.expectReplicas, *rc2.Spec.Replicas)
		}
	}
}

func TestSetDefaultReplicationControllerImagePullPolicy(t *testing.T) {
	containersWithoutPullPolicy, _ := json.Marshal([]map[string]interface{}{
		{
			"name":  "install",
			"image": "busybox:latest",
		},
	})

	containersWithPullPolicy, _ := json.Marshal([]map[string]interface{}{
		{
			"name":            "install",
			"imagePullPolicy": "IfNotPresent",
		},
	})

	tests := []struct {
		rc               metav1.ReplicationController
		expectPullPolicy metav1.PullPolicy
	}{
		{
			rc: metav1.ReplicationController{
				Spec: metav1.ReplicationControllerSpec{
					Template: &metav1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								"pod.beta.kubernetes.io/init-containers": string(containersWithoutPullPolicy),
							},
						},
					},
				},
			},
			expectPullPolicy: metav1.PullAlways,
		},
		{
			rc: metav1.ReplicationController{
				Spec: metav1.ReplicationControllerSpec{
					Template: &metav1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Annotations: map[string]string{
								"pod.beta.kubernetes.io/init-containers": string(containersWithPullPolicy),
							},
						},
					},
				},
			},
			expectPullPolicy: metav1.PullIfNotPresent,
		},
	}

	for _, test := range tests {
		rc := &test.rc
		obj2 := roundTrip(t, runtime.Object(rc))
		rc2, ok := obj2.(*metav1.ReplicationController)
		if !ok {
			t.Errorf("unexpected object: %v", rc2)
			t.FailNow()
		}
		if test.expectPullPolicy != rc2.Spec.Template.Spec.InitContainers[0].ImagePullPolicy {
			t.Errorf("expected ImagePullPolicy: %s, got: %s",
				test.expectPullPolicy,
				rc2.Spec.Template.Spec.InitContainers[0].ImagePullPolicy,
			)
		}
	}
}

func TestSetDefaultService(t *testing.T) {
	svc := &metav1.Service{}
	obj2 := roundTrip(t, runtime.Object(svc))
	svc2 := obj2.(*metav1.Service)
	if svc2.Spec.SessionAffinity != metav1.ServiceAffinityNone {
		t.Errorf("Expected default session affinity type:%s, got: %s", metav1.ServiceAffinityNone, svc2.Spec.SessionAffinity)
	}
	if svc2.Spec.Type != metav1.ServiceTypeClusterIP {
		t.Errorf("Expected default type:%s, got: %s", metav1.ServiceTypeClusterIP, svc2.Spec.Type)
	}
}

func TestSetDefaultSecretVolumeSource(t *testing.T) {
	s := metav1.PodSpec{}
	s.Volumes = []metav1.Volume{
		{
			VolumeSource: metav1.VolumeSource{
				Secret: &metav1.SecretVolumeSource{},
			},
		},
	}
	pod := &metav1.Pod{
		Spec: s,
	}
	output := roundTrip(t, runtime.Object(pod))
	pod2 := output.(*metav1.Pod)
	defaultMode := pod2.Spec.Volumes[0].VolumeSource.Secret.DefaultMode
	expectedMode := metav1.SecretVolumeSourceDefaultMode

	if defaultMode == nil || *defaultMode != expectedMode {
		t.Errorf("Expected secret DefaultMode %v, got %v", expectedMode, defaultMode)
	}
}

func TestSetDefaultConfigMapVolumeSource(t *testing.T) {
	s := metav1.PodSpec{}
	s.Volumes = []metav1.Volume{
		{
			VolumeSource: metav1.VolumeSource{
				ConfigMap: &metav1.ConfigMapVolumeSource{},
			},
		},
	}
	pod := &metav1.Pod{
		Spec: s,
	}
	output := roundTrip(t, runtime.Object(pod))
	pod2 := output.(*metav1.Pod)
	defaultMode := pod2.Spec.Volumes[0].VolumeSource.ConfigMap.DefaultMode
	expectedMode := metav1.ConfigMapVolumeSourceDefaultMode

	if defaultMode == nil || *defaultMode != expectedMode {
		t.Errorf("Expected ConfigMap DefaultMode %v, got %v", expectedMode, defaultMode)
	}
}

func TestSetDefaultDownwardAPIVolumeSource(t *testing.T) {
	s := metav1.PodSpec{}
	s.Volumes = []metav1.Volume{
		{
			VolumeSource: metav1.VolumeSource{
				DownwardAPI: &metav1.DownwardAPIVolumeSource{},
			},
		},
	}
	pod := &metav1.Pod{
		Spec: s,
	}
	output := roundTrip(t, runtime.Object(pod))
	pod2 := output.(*metav1.Pod)
	defaultMode := pod2.Spec.Volumes[0].VolumeSource.DownwardAPI.DefaultMode
	expectedMode := metav1.DownwardAPIVolumeSourceDefaultMode

	if defaultMode == nil || *defaultMode != expectedMode {
		t.Errorf("Expected DownwardAPI DefaultMode %v, got %v", expectedMode, defaultMode)
	}
}

func TestSetDefaultSecret(t *testing.T) {
	s := &metav1.Secret{}
	obj2 := roundTrip(t, runtime.Object(s))
	s2 := obj2.(*metav1.Secret)

	if s2.Type != metav1.SecretTypeOpaque {
		t.Errorf("Expected secret type %v, got %v", metav1.SecretTypeOpaque, s2.Type)
	}
}

func TestSetDefaultPersistentVolume(t *testing.T) {
	pv := &metav1.PersistentVolume{}
	obj2 := roundTrip(t, runtime.Object(pv))
	pv2 := obj2.(*metav1.PersistentVolume)

	if pv2.Status.Phase != metav1.VolumePending {
		t.Errorf("Expected volume phase %v, got %v", metav1.VolumePending, pv2.Status.Phase)
	}
	if pv2.Spec.PersistentVolumeReclaimPolicy != metav1.PersistentVolumeReclaimRetain {
		t.Errorf("Expected pv reclaim policy %v, got %v", metav1.PersistentVolumeReclaimRetain, pv2.Spec.PersistentVolumeReclaimPolicy)
	}
}

func TestSetDefaultPersistentVolumeClaim(t *testing.T) {
	pvc := &metav1.PersistentVolumeClaim{}
	obj2 := roundTrip(t, runtime.Object(pvc))
	pvc2 := obj2.(*metav1.PersistentVolumeClaim)

	if pvc2.Status.Phase != metav1.ClaimPending {
		t.Errorf("Expected claim phase %v, got %v", metav1.ClaimPending, pvc2.Status.Phase)
	}
}

func TestSetDefaulEndpointsProtocol(t *testing.T) {
	in := &metav1.Endpoints{Subsets: []metav1.EndpointSubset{
		{Ports: []metav1.EndpointPort{{}, {Protocol: "UDP"}, {}}},
	}}
	obj := roundTrip(t, runtime.Object(in))
	out := obj.(*metav1.Endpoints)

	for i := range out.Subsets {
		for j := range out.Subsets[i].Ports {
			if in.Subsets[i].Ports[j].Protocol == "" {
				if out.Subsets[i].Ports[j].Protocol != metav1.ProtocolTCP {
					t.Errorf("Expected protocol %s, got %s", metav1.ProtocolTCP, out.Subsets[i].Ports[j].Protocol)
				}
			} else {
				if out.Subsets[i].Ports[j].Protocol != in.Subsets[i].Ports[j].Protocol {
					t.Errorf("Expected protocol %s, got %s", in.Subsets[i].Ports[j].Protocol, out.Subsets[i].Ports[j].Protocol)
				}
			}
		}
	}
}

func TestSetDefaulServiceTargetPort(t *testing.T) {
	in := &metav1.Service{Spec: metav1.ServiceSpec{Ports: []metav1.ServicePort{{Port: 1234}}}}
	obj := roundTrip(t, runtime.Object(in))
	out := obj.(*metav1.Service)
	if out.Spec.Ports[0].TargetPort != intstr.FromInt(1234) {
		t.Errorf("Expected TargetPort to be defaulted, got %v", out.Spec.Ports[0].TargetPort)
	}

	in = &metav1.Service{Spec: metav1.ServiceSpec{Ports: []metav1.ServicePort{{Port: 1234, TargetPort: intstr.FromInt(5678)}}}}
	obj = roundTrip(t, runtime.Object(in))
	out = obj.(*metav1.Service)
	if out.Spec.Ports[0].TargetPort != intstr.FromInt(5678) {
		t.Errorf("Expected TargetPort to be unchanged, got %v", out.Spec.Ports[0].TargetPort)
	}
}

func TestSetDefaultServicePort(t *testing.T) {
	// Unchanged if set.
	in := &metav1.Service{Spec: metav1.ServiceSpec{
		Ports: []metav1.ServicePort{
			{Protocol: "UDP", Port: 9376, TargetPort: intstr.FromString("p")},
			{Protocol: "UDP", Port: 8675, TargetPort: intstr.FromInt(309)},
		},
	}}
	out := roundTrip(t, runtime.Object(in)).(*metav1.Service)
	if out.Spec.Ports[0].Protocol != metav1.ProtocolUDP {
		t.Errorf("Expected protocol %s, got %s", metav1.ProtocolUDP, out.Spec.Ports[0].Protocol)
	}
	if out.Spec.Ports[0].TargetPort != intstr.FromString("p") {
		t.Errorf("Expected port %v, got %v", in.Spec.Ports[0].Port, out.Spec.Ports[0].TargetPort)
	}
	if out.Spec.Ports[1].Protocol != metav1.ProtocolUDP {
		t.Errorf("Expected protocol %s, got %s", metav1.ProtocolUDP, out.Spec.Ports[1].Protocol)
	}
	if out.Spec.Ports[1].TargetPort != intstr.FromInt(309) {
		t.Errorf("Expected port %v, got %v", in.Spec.Ports[1].Port, out.Spec.Ports[1].TargetPort)
	}

	// Defaulted.
	in = &metav1.Service{Spec: metav1.ServiceSpec{
		Ports: []metav1.ServicePort{
			{Protocol: "", Port: 9376, TargetPort: intstr.FromString("")},
			{Protocol: "", Port: 8675, TargetPort: intstr.FromInt(0)},
		},
	}}
	out = roundTrip(t, runtime.Object(in)).(*metav1.Service)
	if out.Spec.Ports[0].Protocol != metav1.ProtocolTCP {
		t.Errorf("Expected protocol %s, got %s", metav1.ProtocolTCP, out.Spec.Ports[0].Protocol)
	}
	if out.Spec.Ports[0].TargetPort != intstr.FromInt(int(in.Spec.Ports[0].Port)) {
		t.Errorf("Expected port %v, got %v", in.Spec.Ports[0].Port, out.Spec.Ports[0].TargetPort)
	}
	if out.Spec.Ports[1].Protocol != metav1.ProtocolTCP {
		t.Errorf("Expected protocol %s, got %s", metav1.ProtocolTCP, out.Spec.Ports[1].Protocol)
	}
	if out.Spec.Ports[1].TargetPort != intstr.FromInt(int(in.Spec.Ports[1].Port)) {
		t.Errorf("Expected port %v, got %v", in.Spec.Ports[1].Port, out.Spec.Ports[1].TargetPort)
	}
}

func TestSetDefaultNamespace(t *testing.T) {
	s := &metav1.Namespace{}
	obj2 := roundTrip(t, runtime.Object(s))
	s2 := obj2.(*metav1.Namespace)

	if s2.Status.Phase != metav1.NamespaceActive {
		t.Errorf("Expected phase %v, got %v", metav1.NamespaceActive, s2.Status.Phase)
	}
}

func TestSetDefaultPodSpecHostNetwork(t *testing.T) {
	portNum := int32(8080)
	s := metav1.PodSpec{}
	s.HostNetwork = true
	s.Containers = []metav1.Container{
		{
			Ports: []metav1.ContainerPort{
				{
					ContainerPort: portNum,
				},
			},
		},
	}
	pod := &metav1.Pod{
		Spec: s,
	}
	obj2 := roundTrip(t, runtime.Object(pod))
	pod2 := obj2.(*metav1.Pod)
	s2 := pod2.Spec

	hostPortNum := s2.Containers[0].Ports[0].HostPort
	if hostPortNum != portNum {
		t.Errorf("Expected container port to be defaulted, was made %d instead of %d", hostPortNum, portNum)
	}
}

func TestSetDefaultNodeExternalID(t *testing.T) {
	name := "node0"
	n := &metav1.Node{}
	n.Name = name
	obj2 := roundTrip(t, runtime.Object(n))
	n2 := obj2.(*metav1.Node)
	if n2.Spec.ExternalID != name {
		t.Errorf("Expected default External ID: %s, got: %s", name, n2.Spec.ExternalID)
	}
	if n2.Spec.ProviderID != "" {
		t.Errorf("Expected empty default Cloud Provider ID, got: %s", n2.Spec.ProviderID)
	}
}

func TestSetDefaultNodeStatusAllocatable(t *testing.T) {
	capacity := metav1.ResourceList{
		metav1.ResourceCPU:    resource.MustParse("1000m"),
		metav1.ResourceMemory: resource.MustParse("10G"),
	}
	allocatable := metav1.ResourceList{
		metav1.ResourceCPU:    resource.MustParse("500m"),
		metav1.ResourceMemory: resource.MustParse("5G"),
	}
	tests := []struct {
		capacity            metav1.ResourceList
		allocatable         metav1.ResourceList
		expectedAllocatable metav1.ResourceList
	}{{ // Everything set, no defaulting.
		capacity:            capacity,
		allocatable:         allocatable,
		expectedAllocatable: allocatable,
	}, { // Allocatable set, no defaulting.
		capacity:            nil,
		allocatable:         allocatable,
		expectedAllocatable: allocatable,
	}, { // Capacity set, allocatable defaults to capacity.
		capacity:            capacity,
		allocatable:         nil,
		expectedAllocatable: capacity,
	}, { // Nothing set, allocatable "defaults" to capacity.
		capacity:            nil,
		allocatable:         nil,
		expectedAllocatable: nil,
	}}

	copyResourceList := func(rl metav1.ResourceList) metav1.ResourceList {
		if rl == nil {
			return nil
		}
		copy := make(metav1.ResourceList, len(rl))
		for k, v := range rl {
			copy[k] = *v.Copy()
		}
		return copy
	}

	resourceListsEqual := func(a metav1.ResourceList, b metav1.ResourceList) bool {
		if len(a) != len(b) {
			return false
		}
		for k, v := range a {
			vb, found := b[k]
			if !found {
				return false
			}
			if v.Cmp(vb) != 0 {
				return false
			}
		}
		return true
	}

	for i, testcase := range tests {
		node := metav1.Node{
			Status: metav1.NodeStatus{
				Capacity:    copyResourceList(testcase.capacity),
				Allocatable: copyResourceList(testcase.allocatable),
			},
		}
		node2 := roundTrip(t, runtime.Object(&node)).(*metav1.Node)
		actual := node2.Status.Allocatable
		expected := testcase.expectedAllocatable
		if !resourceListsEqual(expected, actual) {
			t.Errorf("[%d] Expected NodeStatus.Allocatable: %+v; Got: %+v", i, expected, actual)
		}
	}
}

func TestSetDefaultObjectFieldSelectorAPIVersion(t *testing.T) {
	s := metav1.PodSpec{
		Containers: []metav1.Container{
			{
				Env: []metav1.EnvVar{
					{
						ValueFrom: &metav1.EnvVarSource{
							FieldRef: &metav1.ObjectFieldSelector{},
						},
					},
				},
			},
		},
	}
	pod := &metav1.Pod{
		Spec: s,
	}
	obj2 := roundTrip(t, runtime.Object(pod))
	pod2 := obj2.(*metav1.Pod)
	s2 := pod2.Spec

	apiVersion := s2.Containers[0].Env[0].ValueFrom.FieldRef.APIVersion
	if apiVersion != "v1" {
		t.Errorf("Expected default APIVersion v1, got: %v", apiVersion)
	}
}

func TestSetMinimumScalePod(t *testing.T) {
	// verify we default if limits are specified (and that request=0 is preserved)
	s := metav1.PodSpec{}
	s.Containers = []metav1.Container{
		{
			Resources: metav1.ResourceRequirements{
				Requests: metav1.ResourceList{
					metav1.ResourceMemory: resource.MustParse("1n"),
				},
				Limits: metav1.ResourceList{
					metav1.ResourceCPU: resource.MustParse("2n"),
				},
			},
		},
	}
	pod := &metav1.Pod{
		Spec: s,
	}
	metav1.SetObjectDefaults_Pod(pod)

	if expect := resource.MustParse("1m"); expect.Cmp(pod.Spec.Containers[0].Resources.Requests[metav1.ResourceMemory]) != 0 {
		t.Errorf("did not round resources: %#v", pod.Spec.Containers[0].Resources)
	}
}

func TestSetDefaultRequestsPod(t *testing.T) {
	// verify we default if limits are specified (and that request=0 is preserved)
	s := metav1.PodSpec{}
	s.Containers = []metav1.Container{
		{
			Resources: metav1.ResourceRequirements{
				Requests: metav1.ResourceList{
					metav1.ResourceMemory: resource.MustParse("0"),
				},
				Limits: metav1.ResourceList{
					metav1.ResourceCPU:    resource.MustParse("100m"),
					metav1.ResourceMemory: resource.MustParse("1Gi"),
				},
			},
		},
	}
	pod := &metav1.Pod{
		Spec: s,
	}
	output := roundTrip(t, runtime.Object(pod))
	pod2 := output.(*metav1.Pod)
	defaultRequest := pod2.Spec.Containers[0].Resources.Requests
	if requestValue := defaultRequest[metav1.ResourceCPU]; requestValue.String() != "100m" {
		t.Errorf("Expected request cpu: %s, got: %s", "100m", requestValue.String())
	}
	if requestValue := defaultRequest[metav1.ResourceMemory]; requestValue.String() != "0" {
		t.Errorf("Expected request memory: %s, got: %s", "0", requestValue.String())
	}

	// verify we do nothing if no limits are specified
	s = metav1.PodSpec{}
	s.Containers = []metav1.Container{{}}
	pod = &metav1.Pod{
		Spec: s,
	}
	output = roundTrip(t, runtime.Object(pod))
	pod2 = output.(*metav1.Pod)
	defaultRequest = pod2.Spec.Containers[0].Resources.Requests
	if requestValue := defaultRequest[metav1.ResourceCPU]; requestValue.String() != "0" {
		t.Errorf("Expected 0 request value, got: %s", requestValue.String())
	}
}

func TestDefaultRequestIsNotSetForReplicationController(t *testing.T) {
	s := metav1.PodSpec{}
	s.Containers = []metav1.Container{
		{
			Resources: metav1.ResourceRequirements{
				Limits: metav1.ResourceList{
					metav1.ResourceCPU: resource.MustParse("100m"),
				},
			},
		},
	}
	rc := &metav1.ReplicationController{
		Spec: metav1.ReplicationControllerSpec{
			Replicas: newInt(3),
			Template: &metav1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"foo": "bar",
					},
				},
				Spec: s,
			},
		},
	}
	output := roundTrip(t, runtime.Object(rc))
	rc2 := output.(*metav1.ReplicationController)
	defaultRequest := rc2.Spec.Template.Spec.Containers[0].Resources.Requests
	requestValue := defaultRequest[metav1.ResourceCPU]
	if requestValue.String() != "0" {
		t.Errorf("Expected 0 request value, got: %s", requestValue.String())
	}
}

func TestSetDefaultLimitRangeItem(t *testing.T) {
	limitRange := &metav1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-defaults",
		},
		Spec: metav1.LimitRangeSpec{
			Limits: []metav1.LimitRangeItem{{
				Type: metav1.LimitTypeContainer,
				Max: metav1.ResourceList{
					metav1.ResourceCPU: resource.MustParse("100m"),
				},
				Min: metav1.ResourceList{
					metav1.ResourceMemory: resource.MustParse("100Mi"),
				},
				Default:        metav1.ResourceList{},
				DefaultRequest: metav1.ResourceList{},
			}},
		},
	}

	output := roundTrip(t, runtime.Object(limitRange))
	limitRange2 := output.(*metav1.LimitRange)
	defaultLimit := limitRange2.Spec.Limits[0].Default
	defaultRequest := limitRange2.Spec.Limits[0].DefaultRequest

	// verify that default cpu was set to the max
	defaultValue := defaultLimit[metav1.ResourceCPU]
	if defaultValue.String() != "100m" {
		t.Errorf("Expected default cpu: %s, got: %s", "100m", defaultValue.String())
	}
	// verify that default request was set to the limit
	requestValue := defaultRequest[metav1.ResourceCPU]
	if requestValue.String() != "100m" {
		t.Errorf("Expected request cpu: %s, got: %s", "100m", requestValue.String())
	}
	// verify that if a min is provided, it will be the default if no limit is specified
	requestMinValue := defaultRequest[metav1.ResourceMemory]
	if requestMinValue.String() != "100Mi" {
		t.Errorf("Expected request memory: %s, got: %s", "100Mi", requestMinValue.String())
	}
}

func TestSetDefaultProbe(t *testing.T) {
	originalProbe := metav1.Probe{}
	expectedProbe := metav1.Probe{
		InitialDelaySeconds: 0,
		TimeoutSeconds:      1,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}

	pod := &metav1.Pod{
		Spec: metav1.PodSpec{
			Containers: []metav1.Container{{LivenessProbe: &originalProbe}},
		},
	}

	output := roundTrip(t, runtime.Object(pod)).(*metav1.Pod)
	actualProbe := *output.Spec.Containers[0].LivenessProbe
	if actualProbe != expectedProbe {
		t.Errorf("Expected probe: %+v\ngot: %+v\n", expectedProbe, actualProbe)
	}
}
