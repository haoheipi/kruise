/*
Copyright 2021 The Kruise Authors.

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

package features

import (
	"fmt"
	"os"
	"strings"

	utilfeature "github.com/openkruise/kruise/pkg/util/feature"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/component-base/featuregate"
)

const (
	// KruiseDaemon enables the features relied on kruise-daemon, such as image pulling and container restarting.
	KruiseDaemon featuregate.Feature = "KruiseDaemon"

	// PodWebhook enables webhook for Pods creations. This is also related to SidecarSet.
	PodWebhook featuregate.Feature = "PodWebhook"

	// CloneSetShortHash enables CloneSet controller only set revision hash name to pod label.
	CloneSetShortHash featuregate.Feature = "CloneSetShortHash"

	// KruisePodReadinessGate enables Kruise webhook to inject 'KruisePodReady' readiness-gate to
	// all Pods during creation.
	// Otherwise, it will only be injected to Pods created by Kruise workloads.
	KruisePodReadinessGate featuregate.Feature = "KruisePodReadinessGate"
)

var defaultFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	PodWebhook:             {Default: true, PreRelease: featuregate.Beta},
	KruiseDaemon:           {Default: true, PreRelease: featuregate.Beta},
	CloneSetShortHash:      {Default: false, PreRelease: featuregate.Alpha},
	KruisePodReadinessGate: {Default: false, PreRelease: featuregate.Alpha},
}

func init() {
	compatibleEnv()
	runtime.Must(utilfeature.DefaultMutableFeatureGate.Add(defaultFeatureGates))
}

// Make it compatible with the old CUSTOM_RESOURCE_ENABLE gate in env.
func compatibleEnv() {
	str := strings.TrimSpace(os.Getenv("CUSTOM_RESOURCE_ENABLE"))
	if len(str) == 0 {
		return
	}
	limits := sets.NewString(strings.Split(str, ",")...)
	if !limits.Has("SidecarSet") {
		defaultFeatureGates[PodWebhook] = featuregate.FeatureSpec{Default: false, PreRelease: featuregate.Beta}
	}
}

func ValidateFeatureGates() error {
	if utilfeature.DefaultFeatureGate.Enabled(KruisePodReadinessGate) && !utilfeature.DefaultFeatureGate.Enabled(PodWebhook) {
		return fmt.Errorf("can not enable feature-gate %s because of %s disabled", KruisePodReadinessGate, PodWebhook)
	}
	return nil
}
