/*
 *     Copyright 2026 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/robfig/cron/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"d7y.io/dragonfly/v2/datacontroller/apis/data/v1alpha1"
)

// EffectiveLifecycle is the lifecycle in effect for a dataset after merging the
// inline lifecycle with the matching policy and applying defaults.
type EffectiveLifecycle struct {
	TTL             *time.Duration
	ExpireAction    v1alpha1.ExpireAction
	RefreshInterval *time.Duration
	RefreshSchedule cron.Schedule
	VerifyInterval  *time.Duration
	DeletionPolicy  v1alpha1.DeletionPolicy

	// PolicyName is the name of the DataLifecyclePolicy that contributed, if any.
	PolicyName string
}

// PolicyMatches reports whether a policy selects a dataset.
func PolicyMatches(policy *v1alpha1.DataLifecyclePolicy, dataset *v1alpha1.Dataset) (bool, error) {
	if policy.Namespace != dataset.Namespace {
		return false, nil
	}

	selector := labels.Everything()
	if policy.Spec.DatasetSelector != nil {
		var err error
		selector, err = metav1.LabelSelectorAsSelector(policy.Spec.DatasetSelector)
		if err != nil {
			return false, fmt.Errorf("policy %s/%s: invalid dataset selector: %w", policy.Namespace, policy.Name, err)
		}
	}

	return selector.Matches(labels.Set(dataset.Labels)), nil
}

// SelectPolicy picks the policy governing a dataset: the highest priority among the
// matching policies, ties broken by the lexically smallest name. It returns nil when
// no policy matches.
func SelectPolicy(policies []v1alpha1.DataLifecyclePolicy, dataset *v1alpha1.Dataset) (*v1alpha1.DataLifecyclePolicy, error) {
	var matched []*v1alpha1.DataLifecyclePolicy
	for i := range policies {
		policy := &policies[i]
		ok, err := PolicyMatches(policy, dataset)
		if err != nil {
			return nil, err
		}

		if ok {
			matched = append(matched, policy)
		}
	}

	if len(matched) == 0 {
		return nil, nil
	}

	sort.SliceStable(matched, func(i, j int) bool {
		if matched[i].Spec.Priority != matched[j].Spec.Priority {
			return matched[i].Spec.Priority > matched[j].Spec.Priority
		}

		return matched[i].Name < matched[j].Name
	})

	return matched[0], nil
}

// ResolveLifecycle merges the inline lifecycle of a dataset over the lifecycle of the
// policy and applies defaults. Refresh settings are taken as a pair: when the dataset
// sets either RefreshInterval or RefreshSchedule, the policy's refresh settings are
// ignored entirely.
func ResolveLifecycle(inline *v1alpha1.LifecycleSpec, policy *v1alpha1.DataLifecyclePolicy) (*EffectiveLifecycle, error) {
	merged := v1alpha1.LifecycleSpec{}
	if policy != nil {
		merged = *policy.Spec.Lifecycle.DeepCopy()
	}

	if inline != nil {
		if inline.TTL != nil {
			merged.TTL = inline.TTL
		}

		if inline.ExpireAction != nil {
			merged.ExpireAction = inline.ExpireAction
		}

		if inline.RefreshInterval != nil || inline.RefreshSchedule != nil {
			merged.RefreshInterval = inline.RefreshInterval
			merged.RefreshSchedule = inline.RefreshSchedule
		}

		if inline.VerifyInterval != nil {
			merged.VerifyInterval = inline.VerifyInterval
		}

		if inline.DeletionPolicy != nil {
			merged.DeletionPolicy = inline.DeletionPolicy
		}
	}

	effective := &EffectiveLifecycle{
		ExpireAction:   v1alpha1.ExpireActionPurge,
		DeletionPolicy: v1alpha1.DeletionPolicyRetain,
	}

	if policy != nil {
		effective.PolicyName = policy.Name
	}

	if merged.TTL != nil {
		if merged.TTL.Duration <= 0 {
			return nil, fmt.Errorf("lifecycle: ttl must be positive, got %s", merged.TTL.Duration)
		}

		ttl := merged.TTL.Duration
		effective.TTL = &ttl
	}

	if merged.ExpireAction != nil {
		effective.ExpireAction = *merged.ExpireAction
	}

	if merged.RefreshInterval != nil && merged.RefreshSchedule != nil {
		return nil, fmt.Errorf("lifecycle: refreshInterval and refreshSchedule are mutually exclusive")
	}

	if merged.RefreshInterval != nil {
		if merged.RefreshInterval.Duration <= 0 {
			return nil, fmt.Errorf("lifecycle: refreshInterval must be positive, got %s", merged.RefreshInterval.Duration)
		}

		interval := merged.RefreshInterval.Duration
		effective.RefreshInterval = &interval
	}

	if merged.RefreshSchedule != nil {
		schedule, err := cron.ParseStandard(*merged.RefreshSchedule)
		if err != nil {
			return nil, fmt.Errorf("lifecycle: invalid refreshSchedule %q: %w", *merged.RefreshSchedule, err)
		}

		effective.RefreshSchedule = schedule
	}

	if merged.VerifyInterval != nil {
		if merged.VerifyInterval.Duration <= 0 {
			return nil, fmt.Errorf("lifecycle: verifyInterval must be positive, got %s", merged.VerifyInterval.Duration)
		}

		interval := merged.VerifyInterval.Duration
		effective.VerifyInterval = &interval
	}

	if merged.DeletionPolicy != nil {
		effective.DeletionPolicy = *merged.DeletionPolicy
	}

	return effective, nil
}

// NextRefresh returns when the data should be re-distributed after its last distribution,
// or nil when no refresh is configured or the data was never distributed.
func (l *EffectiveLifecycle) NextRefresh(lastDistribution *time.Time) *time.Time {
	if lastDistribution == nil {
		return nil
	}

	switch {
	case l.RefreshInterval != nil:
		next := lastDistribution.Add(*l.RefreshInterval)
		return &next
	case l.RefreshSchedule != nil:
		next := l.RefreshSchedule.Next(*lastDistribution)
		return &next
	default:
		return nil
	}
}

// Expiration returns when the data expires after its last distribution, or nil when no
// TTL is configured or the data was never distributed.
func (l *EffectiveLifecycle) Expiration(lastDistribution *time.Time) *time.Time {
	if lastDistribution == nil || l.TTL == nil {
		return nil
	}

	expiration := lastDistribution.Add(*l.TTL)
	return &expiration
}

// NextVerification returns when the data should be verified next. The interval counts
// from the last verification, or from the last distribution when the data was never
// verified.
func (l *EffectiveLifecycle) NextVerification(lastDistribution, lastVerification *time.Time) *time.Time {
	if l.VerifyInterval == nil {
		return nil
	}

	base := lastVerification
	if base == nil {
		base = lastDistribution
	}

	if base == nil {
		return nil
	}

	next := base.Add(*l.VerifyInterval)
	return &next
}

// SourceSpecHash hashes the parts of the spec that, when changed, require the source to
// be distributed again.
func SourceSpecHash(source *v1alpha1.DataSource, distribution *v1alpha1.DistributionPolicy) string {
	data, err := json.Marshal(struct {
		Source       *v1alpha1.DataSource         `json:"source"`
		Distribution *v1alpha1.DistributionPolicy `json:"distribution"`
	}{
		Source:       source,
		Distribution: distribution,
	})
	if err != nil {
		// The types are plain data and always marshal; fall back to a hash of the
		// formatted value so a change is still detected.
		data = []byte(fmt.Sprintf("%#v%#v", source, distribution))
	}

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:8])
}
