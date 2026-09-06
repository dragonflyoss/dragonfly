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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"d7y.io/dragonfly/v2/datacontroller/apis/data/v1alpha1"
)

func duration(d time.Duration) *metav1.Duration {
	return &metav1.Duration{Duration: d}
}

func stringPtr(s string) *string { return &s }

func policy(name string, priority int32, selector map[string]string, lifecycle v1alpha1.LifecycleSpec) v1alpha1.DataLifecyclePolicy {
	p := v1alpha1.DataLifecyclePolicy{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default"},
		Spec:       v1alpha1.DataLifecyclePolicySpec{Priority: priority, Lifecycle: lifecycle},
	}

	if selector != nil {
		p.Spec.DatasetSelector = &metav1.LabelSelector{MatchLabels: selector}
	}

	return p
}

func TestSelectPolicy(t *testing.T) {
	dataset := &v1alpha1.Dataset{ObjectMeta: metav1.ObjectMeta{Name: "ds", Namespace: "default", Labels: map[string]string{"tier": "hot"}}}

	tests := []struct {
		name     string
		policies []v1alpha1.DataLifecyclePolicy
		expect   string
		wantErr  bool
	}{
		{
			name:     "no policies",
			policies: nil,
			expect:   "",
		},
		{
			name: "empty selector matches everything",
			policies: []v1alpha1.DataLifecyclePolicy{
				policy("all", 0, nil, v1alpha1.LifecycleSpec{}),
			},
			expect: "all",
		},
		{
			name: "selector mismatch is skipped",
			policies: []v1alpha1.DataLifecyclePolicy{
				policy("cold", 10, map[string]string{"tier": "cold"}, v1alpha1.LifecycleSpec{}),
				policy("hot", 0, map[string]string{"tier": "hot"}, v1alpha1.LifecycleSpec{}),
			},
			expect: "hot",
		},
		{
			name: "highest priority wins",
			policies: []v1alpha1.DataLifecyclePolicy{
				policy("low", 1, nil, v1alpha1.LifecycleSpec{}),
				policy("high", 5, map[string]string{"tier": "hot"}, v1alpha1.LifecycleSpec{}),
			},
			expect: "high",
		},
		{
			name: "ties broken by name",
			policies: []v1alpha1.DataLifecyclePolicy{
				policy("b", 1, nil, v1alpha1.LifecycleSpec{}),
				policy("a", 1, nil, v1alpha1.LifecycleSpec{}),
			},
			expect: "a",
		},
		{
			name: "other namespace is ignored",
			policies: []v1alpha1.DataLifecyclePolicy{
				{ObjectMeta: metav1.ObjectMeta{Name: "other", Namespace: "other"}},
			},
			expect: "",
		},
		{
			name: "invalid selector is an error",
			policies: []v1alpha1.DataLifecyclePolicy{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "bad", Namespace: "default"},
					Spec: v1alpha1.DataLifecyclePolicySpec{DatasetSelector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{{Key: "tier", Operator: "Bogus"}},
					}},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			selected, err := SelectPolicy(tc.policies, dataset)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tc.expect == "" {
				assert.Nil(t, selected)
				return
			}

			require.NotNil(t, selected)
			assert.Equal(t, tc.expect, selected.Name)
		})
	}
}

func TestResolveLifecycle(t *testing.T) {
	purge := v1alpha1.DeletionPolicyPurge
	retain := v1alpha1.ExpireActionRetain

	t.Run("defaults", func(t *testing.T) {
		lc, err := ResolveLifecycle(nil, nil)
		require.NoError(t, err)
		assert.Nil(t, lc.TTL)
		assert.Nil(t, lc.RefreshInterval)
		assert.Nil(t, lc.RefreshSchedule)
		assert.Nil(t, lc.VerifyInterval)
		assert.Equal(t, v1alpha1.ExpireActionPurge, lc.ExpireAction)
		assert.Equal(t, v1alpha1.DeletionPolicyRetain, lc.DeletionPolicy)
		assert.Empty(t, lc.PolicyName)
	})

	t.Run("policy applies", func(t *testing.T) {
		p := policy("p", 0, nil, v1alpha1.LifecycleSpec{
			TTL:             duration(time.Hour),
			RefreshSchedule: stringPtr("0 2 * * *"),
			VerifyInterval:  duration(10 * time.Minute),
			DeletionPolicy:  &purge,
		})

		lc, err := ResolveLifecycle(nil, &p)
		require.NoError(t, err)
		assert.Equal(t, time.Hour, *lc.TTL)
		assert.NotNil(t, lc.RefreshSchedule)
		assert.Nil(t, lc.RefreshInterval)
		assert.Equal(t, 10*time.Minute, *lc.VerifyInterval)
		assert.Equal(t, v1alpha1.DeletionPolicyPurge, lc.DeletionPolicy)
		assert.Equal(t, "p", lc.PolicyName)
	})

	t.Run("inline overrides policy and refresh is taken as a pair", func(t *testing.T) {
		p := policy("p", 0, nil, v1alpha1.LifecycleSpec{
			TTL:             duration(time.Hour),
			RefreshSchedule: stringPtr("0 2 * * *"),
			DeletionPolicy:  &purge,
		})

		inline := &v1alpha1.LifecycleSpec{
			TTL:             duration(2 * time.Hour),
			ExpireAction:    &retain,
			RefreshInterval: duration(30 * time.Minute),
		}

		lc, err := ResolveLifecycle(inline, &p)
		require.NoError(t, err)
		assert.Equal(t, 2*time.Hour, *lc.TTL)
		assert.Equal(t, v1alpha1.ExpireActionRetain, lc.ExpireAction)
		assert.Equal(t, 30*time.Minute, *lc.RefreshInterval)
		assert.Nil(t, lc.RefreshSchedule, "policy schedule must be dropped when the dataset sets an interval")
		assert.Equal(t, v1alpha1.DeletionPolicyPurge, lc.DeletionPolicy, "unset inline fields keep the policy value")
	})

	t.Run("invalid values", func(t *testing.T) {
		_, err := ResolveLifecycle(&v1alpha1.LifecycleSpec{TTL: duration(-time.Second)}, nil)
		assert.Error(t, err)

		_, err = ResolveLifecycle(&v1alpha1.LifecycleSpec{RefreshInterval: duration(0)}, nil)
		assert.Error(t, err)

		_, err = ResolveLifecycle(&v1alpha1.LifecycleSpec{RefreshSchedule: stringPtr("not a cron")}, nil)
		assert.Error(t, err)

		_, err = ResolveLifecycle(&v1alpha1.LifecycleSpec{RefreshInterval: duration(time.Minute), RefreshSchedule: stringPtr("* * * * *")}, nil)
		assert.Error(t, err)

		_, err = ResolveLifecycle(&v1alpha1.LifecycleSpec{VerifyInterval: duration(-time.Minute)}, nil)
		assert.Error(t, err)
	})
}

func TestEffectiveLifecycleTimes(t *testing.T) {
	last := time.Date(2026, 1, 1, 1, 30, 0, 0, time.UTC)
	verified := last.Add(20 * time.Minute)

	lc, err := ResolveLifecycle(&v1alpha1.LifecycleSpec{
		TTL:             duration(time.Hour),
		RefreshSchedule: stringPtr("0 2 * * *"),
		VerifyInterval:  duration(15 * time.Minute),
	}, nil)
	require.NoError(t, err)

	assert.Nil(t, lc.NextRefresh(nil))
	assert.Nil(t, lc.Expiration(nil))
	assert.Nil(t, lc.NextVerification(nil, nil))

	assert.Equal(t, time.Date(2026, 1, 1, 2, 0, 0, 0, time.UTC), *lc.NextRefresh(&last))
	assert.Equal(t, last.Add(time.Hour), *lc.Expiration(&last))
	assert.Equal(t, last.Add(15*time.Minute), *lc.NextVerification(&last, nil))
	assert.Equal(t, verified.Add(15*time.Minute), *lc.NextVerification(&last, &verified))

	interval, err := ResolveLifecycle(&v1alpha1.LifecycleSpec{RefreshInterval: duration(45 * time.Minute)}, nil)
	require.NoError(t, err)
	assert.Equal(t, last.Add(45*time.Minute), *interval.NextRefresh(&last))
}

func TestSourceSpecHash(t *testing.T) {
	source := v1alpha1.DataSource{Name: "a", Type: v1alpha1.DataSourceTypeFile, URL: "https://example.com/a"}
	distribution := v1alpha1.DistributionPolicy{Scope: v1alpha1.DistributionScopeAllSeedPeers}

	base := SourceSpecHash(&source, &distribution)
	assert.Len(t, base, 16)
	assert.Equal(t, base, SourceSpecHash(&source, &distribution), "hash is deterministic")

	changedSource := source
	changedSource.Tag = "v2"
	assert.NotEqual(t, base, SourceSpecHash(&changedSource, &distribution))

	changedDistribution := distribution
	changedDistribution.Scope = v1alpha1.DistributionScopeAllPeers
	assert.NotEqual(t, base, SourceSpecHash(&source, &changedDistribution))
}

func TestRetryBackoff(t *testing.T) {
	assert.Equal(t, 30*time.Second, retryBackoff(30*time.Second, 10*time.Minute, 0))
	assert.Equal(t, 30*time.Second, retryBackoff(30*time.Second, 10*time.Minute, 1))
	assert.Equal(t, time.Minute, retryBackoff(30*time.Second, 10*time.Minute, 2))
	assert.Equal(t, 4*time.Minute, retryBackoff(30*time.Second, 10*time.Minute, 4))
	assert.Equal(t, 10*time.Minute, retryBackoff(30*time.Second, 10*time.Minute, 10))
	assert.Equal(t, 10*time.Minute, retryBackoff(30*time.Second, 10*time.Minute, 60), "does not overflow")
}
