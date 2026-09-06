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
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/events"
	clocktesting "k8s.io/utils/clock/testing"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"d7y.io/dragonfly/v2/datacontroller/apis/data/v1alpha1"
	"d7y.io/dragonfly/v2/datacontroller/config"
	"d7y.io/dragonfly/v2/datacontroller/dragonfly"
	managertypes "d7y.io/dragonfly/v2/manager/types"
)

// fakeManager is an in-memory Dragonfly manager job API.
type fakeManager struct {
	mu        sync.Mutex
	nextID    uint
	jobs      map[uint]*dragonfly.Job
	requests  []any
	createErr error
}

func newFakeManager() *fakeManager {
	return &fakeManager{jobs: map[uint]*dragonfly.Job{}}
}

func (f *fakeManager) create(jobType string, req any) (*dragonfly.Job, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.createErr != nil {
		return nil, f.createErr
	}

	f.nextID++
	job := &dragonfly.Job{ID: f.nextID, Type: jobType, State: dragonfly.JobStatePending}
	f.jobs[job.ID] = job
	f.requests = append(f.requests, req)
	return job, nil
}

func (f *fakeManager) CreatePreheatJob(_ context.Context, req managertypes.CreatePreheatJobRequest) (*dragonfly.Job, error) {
	return f.create(dragonfly.JobTypePreheat, req)
}

func (f *fakeManager) CreateGetTaskJob(_ context.Context, req managertypes.CreateGetTaskJobRequest) (*dragonfly.Job, error) {
	return f.create(dragonfly.JobTypeGetTask, req)
}

func (f *fakeManager) CreateDeleteTaskJob(_ context.Context, req managertypes.CreateDeleteTaskJobRequest) (*dragonfly.Job, error) {
	return f.create(dragonfly.JobTypeDeleteTask, req)
}

func (f *fakeManager) GetJob(_ context.Context, id uint) (*dragonfly.Job, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	job, ok := f.jobs[id]
	if !ok {
		return nil, dragonfly.ErrNotFound
	}

	copied := *job
	return &copied, nil
}

// finish moves a job to a terminal state with the given per-scheduler results.
func (f *fakeManager) finish(id uint, state string, results ...any) {
	f.mu.Lock()
	defer f.mu.Unlock()

	job := f.jobs[id]
	job.State = state
	job.Result = map[string]any{
		"group_uuid": fmt.Sprintf("group_%d", id),
		"state":      state,
		"job_states": []any{map[string]any{"task_name": job.Type, "state": state, "results": results}},
	}
}

// finishAll moves every non-terminal job of a type to a terminal state.
func (f *fakeManager) finishAll(jobType, state string, results ...any) []uint {
	f.mu.Lock()
	ids := []uint{}
	for id, job := range f.jobs {
		if job.Type == jobType && !dragonfly.IsTerminalState(job.State) {
			ids = append(ids, id)
		}
	}
	f.mu.Unlock()

	for _, id := range ids {
		f.finish(id, state, results...)
	}

	return ids
}

func (f *fakeManager) countJobs(jobType string) int {
	f.mu.Lock()
	defer f.mu.Unlock()

	count := 0
	for _, job := range f.jobs {
		if job.Type == jobType {
			count++
		}
	}

	return count
}

func preheatSuccess(urls ...string) map[string]any {
	tasks := []any{}
	for _, u := range urls {
		tasks = append(tasks, map[string]any{"url": u, "hostname": "seed", "ip": "10.0.0.1"})
	}

	return map[string]any{"success_tasks": tasks, "failure_tasks": []any{}, "scheduler_cluster_id": 1}
}

func getTaskPeers(n int) map[string]any {
	peers := []any{}
	for i := 0; i < n; i++ {
		peers = append(peers, map[string]any{"id": fmt.Sprintf("peer-%d", i), "ip": fmt.Sprintf("10.0.0.%d", i+1)})
	}

	return map[string]any{"peers": peers, "scheduler_cluster_id": 1}
}

func deleteTaskSuccess() map[string]any {
	return map[string]any{"success_tasks": []any{map[string]any{"ip": "10.0.0.1"}}, "failure_tasks": []any{}, "scheduler_cluster_id": 1}
}

// harness wires a reconciler to a fake API server and a fake manager.
type harness struct {
	t          *testing.T
	client     client.Client
	manager    *fakeManager
	clock      *clocktesting.FakePassiveClock
	recorder   *events.FakeRecorder
	reconciler *DatasetReconciler
	cfg        *config.Config
}

func newHarness(t *testing.T, objects ...client.Object) *harness {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, v1alpha1.AddToScheme(scheme))

	c := fake.NewClientBuilder().
		WithScheme(scheme).
		WithStatusSubresource(&v1alpha1.Dataset{}, &v1alpha1.DataLifecyclePolicy{}).
		WithObjects(objects...).
		Build()

	cfg := config.New()
	cfg.Manager.Endpoint = "http://manager:8080"
	cfg.Manager.Token = "token"
	cfg.Controller.ResyncInterval = time.Hour

	manager := newFakeManager()
	h := &harness{
		t:        t,
		client:   c,
		manager:  manager,
		clock:    clocktesting.NewFakePassiveClock(time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)),
		recorder: events.NewFakeRecorder(100),
		cfg:      cfg,
	}

	h.reconciler = &DatasetReconciler{
		Client:   c,
		Scheme:   scheme,
		Recorder: h.recorder,
		Config:   cfg,
		Clock:    h.clock,
		NewManagerClient: func(opts dragonfly.Options) (dragonfly.Client, error) {
			if opts.Endpoint == "" || opts.Token == "" {
				return nil, errors.New("missing endpoint or token")
			}

			return manager, nil
		},
	}

	return h
}

func (h *harness) reconcile(name string) (ctrl.Result, error) {
	h.t.Helper()
	return h.reconciler.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: name}})
}

// reconcileOK reconciles and fails the test on error.
func (h *harness) reconcileOK(name string) ctrl.Result {
	h.t.Helper()
	result, err := h.reconcile(name)
	require.NoError(h.t, err)
	return result
}

func (h *harness) get(name string) *v1alpha1.Dataset {
	h.t.Helper()
	var dataset v1alpha1.Dataset
	require.NoError(h.t, h.client.Get(context.Background(), types.NamespacedName{Namespace: "default", Name: name}, &dataset))
	return &dataset
}

func (h *harness) advance(d time.Duration) {
	h.clock.SetTime(h.clock.Now().Add(d))
}

func (h *harness) drainEvents() []string {
	var events []string
	for {
		select {
		case e := <-h.recorder.Events:
			events = append(events, e)
		default:
			return events
		}
	}
}

func fileDataset(name string, urls ...string) *v1alpha1.Dataset {
	return &v1alpha1.Dataset{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "default", Labels: map[string]string{"tier": "hot"}},
		Spec: v1alpha1.DatasetSpec{
			Sources: []v1alpha1.DataSource{{
				Name: "files",
				Type: v1alpha1.DataSourceTypeFile,
				URLs: urls,
				Tag:  "v1",
			}},
			Distribution: v1alpha1.DistributionPolicy{Scope: v1alpha1.DistributionScopeAllSeedPeers},
		},
	}
}

func sourceStatus(dataset *v1alpha1.Dataset, name string) *v1alpha1.DataSourceStatus {
	for i := range dataset.Status.Sources {
		if dataset.Status.Sources[i].Name == name {
			return &dataset.Status.Sources[i]
		}
	}

	return nil
}

// assertTime compares instants regardless of location, which a round trip through
// the API server changes.
func assertTime(t *testing.T, expected time.Time, actual *metav1.Time, msgAndArgs ...any) {
	t.Helper()
	require.NotNil(t, actual, msgAndArgs...)
	assert.True(t, expected.Equal(actual.Time), "expected %s, got %s", expected, actual.Time)
}

func condition(dataset *v1alpha1.Dataset, conditionType string) *metav1.Condition {
	return meta.FindStatusCondition(dataset.Status.Conditions, conditionType)
}

// distribute takes a fresh dataset through finalizer, preheat and success.
func (h *harness) distribute(name string, urls ...string) {
	h.t.Helper()

	// First reconcile adds the finalizer.
	h.reconcileOK(name)
	require.Contains(h.t, h.get(name).Finalizers, v1alpha1.DatasetFinalizer)

	// Second reconcile creates the preheat job.
	result := h.reconcileOK(name)
	assert.Equal(h.t, h.cfg.Controller.PollInterval, result.RequeueAfter)

	ds := h.get(name)
	assert.Equal(h.t, v1alpha1.DatasetPhaseDistributing, ds.Status.Phase)
	require.NotNil(h.t, sourceStatus(ds, "files").Operation)
	assert.Equal(h.t, v1alpha1.OperationTypeDistribute, sourceStatus(ds, "files").Operation.Type)

	ids := h.manager.finishAll(dragonfly.JobTypePreheat, dragonfly.JobStateSuccess, preheatSuccess(urls...))
	require.Len(h.t, ids, 1)

	h.reconcileOK(name)
}

func TestDatasetReconciler_Distribute(t *testing.T) {
	h := newHarness(t, fileDataset("ds", "https://example.com/a", "https://example.com/b"))
	h.distribute("ds", "https://example.com/a", "https://example.com/b")

	ds := h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseReady, ds.Status.Phase)
	assert.Equal(t, ds.Generation, ds.Status.ObservedGeneration)
	assertTime(t, h.clock.Now(), ds.Status.LastDistributionTime)
	assert.Nil(t, ds.Status.ExpirationTime, "no ttl configured")
	assert.Nil(t, ds.Status.NextRefreshTime, "no refresh configured")

	ready := condition(ds, v1alpha1.ConditionReady)
	require.NotNil(t, ready)
	assert.Equal(t, metav1.ConditionTrue, ready.Status)
	assert.Equal(t, metav1.ConditionTrue, condition(ds, v1alpha1.ConditionDistributed).Status)
	assert.Nil(t, condition(ds, v1alpha1.ConditionVerified), "no verification configured")
	assert.Nil(t, condition(ds, v1alpha1.ConditionExpired), "no ttl configured")

	source := sourceStatus(ds, "files")
	assert.Equal(t, v1alpha1.DataSourceStateReady, source.State)
	assert.Equal(t, ReasonDistributed, source.Reason)
	assert.Nil(t, source.Operation)
	require.Len(t, source.Tasks, 2)
	assert.Equal(t, "https://example.com/a", source.Tasks[0].URL)
	assert.NotEmpty(t, source.Tasks[0].ID, "task ids are computed locally")
	assert.NotEqual(t, source.Tasks[0].ID, source.Tasks[1].ID)

	// The preheat request carries the spec.
	require.Len(t, h.manager.requests, 1)
	req := h.manager.requests[0].(managertypes.CreatePreheatJobRequest)
	assert.Equal(t, "file", req.Args.Type)
	assert.Equal(t, []string{"https://example.com/a", "https://example.com/b"}, req.Args.URLs)
	assert.Equal(t, "v1", req.Args.Tag)
	assert.Equal(t, managertypes.AllSeedPeersScope, req.Args.Scope)

	// The reconcile that finished the job requeues right away so the next step
	// starts from the recorded outcome; a steady dataset then only resyncs.
	h.reconcileOK("ds")
	result := h.reconcileOK("ds")
	assert.Equal(t, h.cfg.Controller.ResyncInterval, result.RequeueAfter)
	assert.Equal(t, 1, h.manager.countJobs(dragonfly.JobTypePreheat), "no new job on a steady dataset")
}

func TestDatasetReconciler_SpecChangeRedistributes(t *testing.T) {
	h := newHarness(t, fileDataset("ds", "https://example.com/a"))
	h.distribute("ds", "https://example.com/a")

	ds := h.get("ds")
	ds.Spec.Sources[0].Tag = "v2"
	require.NoError(t, h.client.Update(context.Background(), ds))

	h.reconcileOK("ds")
	ds = h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseDistributing, ds.Status.Phase)
	assert.Equal(t, 2, h.manager.countJobs(dragonfly.JobTypePreheat))
	assert.Equal(t, v1alpha1.DataSourceStateDistributing, sourceStatus(ds, "files").State)
}

func TestDatasetReconciler_RetryThenFail(t *testing.T) {
	ds := fileDataset("ds", "https://example.com/a")
	retryLimit := int32(1)
	ds.Spec.Distribution.RetryLimit = &retryLimit
	h := newHarness(t, ds)

	h.reconcileOK("ds")
	h.reconcileOK("ds")
	h.manager.finishAll(dragonfly.JobTypePreheat, dragonfly.JobStateFailure)

	h.reconcileOK("ds")
	result := h.reconcileOK("ds")
	got := h.get("ds")
	source := sourceStatus(got, "files")
	assert.Equal(t, v1alpha1.DataSourceStatePending, source.State)
	assert.Equal(t, ReasonRetryScheduled, source.Reason)
	assert.Equal(t, int32(1), source.Retries)
	require.NotNil(t, source.NextRetryTime)
	assert.Equal(t, h.cfg.Controller.RetryBackoff, result.RequeueAfter)
	assert.Equal(t, v1alpha1.DatasetPhasePending, got.Status.Phase)

	// Not due yet: nothing happens.
	h.advance(h.cfg.Controller.RetryBackoff / 2)
	h.reconcileOK("ds")
	assert.Equal(t, 1, h.manager.countJobs(dragonfly.JobTypePreheat))

	// Due: a new preheat job is created.
	h.advance(h.cfg.Controller.RetryBackoff)
	h.reconcileOK("ds")
	assert.Equal(t, 2, h.manager.countJobs(dragonfly.JobTypePreheat))

	h.manager.finishAll(dragonfly.JobTypePreheat, dragonfly.JobStateFailure)
	h.reconcileOK("ds")
	result = h.reconcileOK("ds")
	got = h.get("ds")
	source = sourceStatus(got, "files")
	assert.Equal(t, v1alpha1.DataSourceStateFailed, source.State)
	assert.Equal(t, ReasonDistributionFailed, source.Reason)
	assert.Nil(t, source.NextRetryTime)
	assert.Equal(t, v1alpha1.DatasetPhaseFailed, got.Status.Phase)
	assert.Equal(t, metav1.ConditionFalse, condition(got, v1alpha1.ConditionReady).Status)
	assert.Equal(t, h.cfg.Controller.ResyncInterval, result.RequeueAfter)

	// Failed sources stay failed until the spec changes.
	h.reconcileOK("ds")
	assert.Equal(t, 2, h.manager.countJobs(dragonfly.JobTypePreheat))
}

func TestDatasetReconciler_RefreshInterval(t *testing.T) {
	ds := fileDataset("ds", "https://example.com/a")
	ds.Spec.Lifecycle = &v1alpha1.LifecycleSpec{RefreshInterval: duration(30 * time.Minute)}
	h := newHarness(t, ds)
	h.distribute("ds", "https://example.com/a")

	got := h.get("ds")
	assertTime(t, h.clock.Now().Add(30*time.Minute), got.Status.NextRefreshTime)

	result := h.reconcileOK("ds")
	assert.Equal(t, 30*time.Minute, result.RequeueAfter)

	h.advance(30 * time.Minute)
	h.reconcileOK("ds")
	assert.Equal(t, 2, h.manager.countJobs(dragonfly.JobTypePreheat), "refresh creates a new preheat job")
	assert.Equal(t, v1alpha1.DatasetPhaseDistributing, h.get("ds").Status.Phase)

	h.manager.finishAll(dragonfly.JobTypePreheat, dragonfly.JobStateSuccess, preheatSuccess("https://example.com/a"))
	h.reconcileOK("ds")
	got = h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseReady, got.Status.Phase)
	assertTime(t, h.clock.Now().Add(30*time.Minute), got.Status.NextRefreshTime, "next refresh moves forward")
}

func TestDatasetReconciler_ExpirePurges(t *testing.T) {
	ds := fileDataset("ds", "https://example.com/a", "https://example.com/b")
	ds.Spec.Lifecycle = &v1alpha1.LifecycleSpec{TTL: duration(time.Hour)}
	h := newHarness(t, ds)
	h.distribute("ds", "https://example.com/a", "https://example.com/b")

	got := h.get("ds")
	assertTime(t, h.clock.Now().Add(time.Hour), got.Status.ExpirationTime)
	assert.Equal(t, metav1.ConditionFalse, condition(got, v1alpha1.ConditionExpired).Status)

	h.advance(time.Hour)
	result := h.reconcileOK("ds")
	assert.Equal(t, h.cfg.Controller.PollInterval, result.RequeueAfter)
	assert.Equal(t, 2, h.manager.countJobs(dragonfly.JobTypeDeleteTask), "one delete_task job per task")

	got = h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhasePurging, got.Status.Phase)
	assert.Equal(t, v1alpha1.DataSourceStatePurging, sourceStatus(got, "files").State)

	// The delete requests address the tasks by the locally computed ID.
	for _, req := range h.manager.requests[1:] {
		deleteReq := req.(managertypes.CreateDeleteTaskJobRequest)
		assert.NotEmpty(t, deleteReq.Args.TaskID)
		assert.Empty(t, deleteReq.Args.URL)
	}

	h.manager.finishAll(dragonfly.JobTypeDeleteTask, dragonfly.JobStateSuccess, deleteTaskSuccess())
	h.reconcileOK("ds")
	got = h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseExpired, got.Status.Phase)
	assert.Equal(t, metav1.ConditionTrue, condition(got, v1alpha1.ConditionExpired).Status)
	source := sourceStatus(got, "files")
	assert.Equal(t, v1alpha1.DataSourceStateExpired, source.State)
	assert.Empty(t, source.Tasks, "purged tasks are forgotten")

	// Expired data stays expired.
	h.reconcileOK("ds")
	assert.Equal(t, 2, h.manager.countJobs(dragonfly.JobTypeDeleteTask))
	assert.Equal(t, 1, h.manager.countJobs(dragonfly.JobTypePreheat))
}

func TestDatasetReconciler_ExpireRetains(t *testing.T) {
	ds := fileDataset("ds", "https://example.com/a")
	retain := v1alpha1.ExpireActionRetain
	ds.Spec.Lifecycle = &v1alpha1.LifecycleSpec{TTL: duration(time.Hour), ExpireAction: &retain}
	h := newHarness(t, ds)
	h.distribute("ds", "https://example.com/a")

	h.advance(time.Hour)
	h.reconcileOK("ds")
	got := h.get("ds")
	assert.Equal(t, 0, h.manager.countJobs(dragonfly.JobTypeDeleteTask))
	assert.Equal(t, v1alpha1.DatasetPhaseExpired, got.Status.Phase)
	assert.Len(t, sourceStatus(got, "files").Tasks, 1, "retained tasks are remembered")
}

func TestDatasetReconciler_RefreshRevivesExpired(t *testing.T) {
	ds := fileDataset("ds", "https://example.com/a")
	ds.Spec.Lifecycle = &v1alpha1.LifecycleSpec{TTL: duration(time.Hour), RefreshInterval: duration(2 * time.Hour)}
	h := newHarness(t, ds)
	h.distribute("ds", "https://example.com/a")

	h.advance(time.Hour)
	h.reconcileOK("ds")
	h.manager.finishAll(dragonfly.JobTypeDeleteTask, dragonfly.JobStateSuccess, deleteTaskSuccess())
	h.reconcileOK("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseExpired, h.get("ds").Status.Phase)

	h.advance(time.Hour)
	h.reconcileOK("ds")
	assert.Equal(t, 2, h.manager.countJobs(dragonfly.JobTypePreheat))
	h.manager.finishAll(dragonfly.JobTypePreheat, dragonfly.JobStateSuccess, preheatSuccess("https://example.com/a"))
	h.reconcileOK("ds")
	got := h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseReady, got.Status.Phase)
	assert.Equal(t, metav1.ConditionFalse, condition(got, v1alpha1.ConditionExpired).Status)
}

func TestDatasetReconciler_Verify(t *testing.T) {
	ds := fileDataset("ds", "https://example.com/a", "https://example.com/b")
	ds.Spec.Lifecycle = &v1alpha1.LifecycleSpec{VerifyInterval: duration(10 * time.Minute)}
	h := newHarness(t, ds)
	h.distribute("ds", "https://example.com/a", "https://example.com/b")

	got := h.get("ds")
	assert.Equal(t, metav1.ConditionUnknown, condition(got, v1alpha1.ConditionVerified).Status)
	require.NotNil(t, got.Status.NextVerificationTime)

	h.advance(10 * time.Minute)
	h.reconcileOK("ds")
	assert.Equal(t, 2, h.manager.countJobs(dragonfly.JobTypeGetTask))
	got = h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseReady, got.Status.Phase, "verification does not change the phase")
	assert.Equal(t, v1alpha1.OperationTypeVerify, sourceStatus(got, "files").Operation.Type)

	// One task is present, the other is gone from every peer.
	ids := h.manager.finishAll(dragonfly.JobTypeGetTask, dragonfly.JobStateSuccess, getTaskPeers(3))
	require.Len(t, ids, 2)
	source := sourceStatus(got, "files")
	var missingURL string
	for _, ref := range source.Operation.Jobs {
		if ref.ID == ids[0] {
			h.manager.finish(ref.ID, dragonfly.JobStateSuccess, getTaskPeers(0))
			missingURL = ref.URL
		}
	}

	h.reconcileOK("ds")
	got = h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseDegraded, got.Status.Phase)
	assert.Equal(t, metav1.ConditionFalse, condition(got, v1alpha1.ConditionVerified).Status)
	assert.Equal(t, ReasonDataMissing, condition(got, v1alpha1.ConditionVerified).Reason)
	source = sourceStatus(got, "files")
	assert.Equal(t, v1alpha1.DataSourceStateDegraded, source.State)
	require.NotNil(t, source.LastVerificationTime)
	for _, task := range source.Tasks {
		require.NotNil(t, task.PeerCount)
		if task.URL == missingURL {
			assert.Equal(t, int32(0), *task.PeerCount)
		} else {
			assert.Equal(t, int32(3), *task.PeerCount)
		}
	}

	// The next verification finds everything.
	h.advance(10 * time.Minute)
	h.reconcileOK("ds")
	h.manager.finishAll(dragonfly.JobTypeGetTask, dragonfly.JobStateSuccess, getTaskPeers(2))
	h.reconcileOK("ds")
	got = h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseReady, got.Status.Phase)
	assert.Equal(t, metav1.ConditionTrue, condition(got, v1alpha1.ConditionVerified).Status)
}

func TestDatasetReconciler_DeleteWithPurgePolicy(t *testing.T) {
	purge := v1alpha1.DeletionPolicyPurge
	p := policy("purge-on-delete", 0, map[string]string{"tier": "hot"}, v1alpha1.LifecycleSpec{DeletionPolicy: &purge})
	h := newHarness(t, fileDataset("ds", "https://example.com/a"), &p)
	h.distribute("ds", "https://example.com/a")

	got := h.get("ds")
	assert.Equal(t, "purge-on-delete", got.Status.AppliedPolicy)

	require.NoError(t, h.client.Delete(context.Background(), got))
	got = h.get("ds")
	require.False(t, got.DeletionTimestamp.IsZero(), "finalizer keeps the object")

	result := h.reconcileOK("ds")
	assert.Equal(t, h.cfg.Controller.PollInterval, result.RequeueAfter)
	assert.Equal(t, 1, h.manager.countJobs(dragonfly.JobTypeDeleteTask))
	got = h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhasePurging, got.Status.Phase)

	h.manager.finishAll(dragonfly.JobTypeDeleteTask, dragonfly.JobStateSuccess, deleteTaskSuccess())
	h.reconcileOK("ds")

	var dataset v1alpha1.Dataset
	err := h.client.Get(context.Background(), types.NamespacedName{Namespace: "default", Name: "ds"}, &dataset)
	assert.True(t, apierrors.IsNotFound(err), "finalizer removed after the purge")
}

func TestDatasetReconciler_DeleteWithRetainPolicy(t *testing.T) {
	h := newHarness(t, fileDataset("ds", "https://example.com/a"))
	h.distribute("ds", "https://example.com/a")

	require.NoError(t, h.client.Delete(context.Background(), h.get("ds")))
	h.reconcileOK("ds")

	var dataset v1alpha1.Dataset
	err := h.client.Get(context.Background(), types.NamespacedName{Namespace: "default", Name: "ds"}, &dataset)
	assert.True(t, apierrors.IsNotFound(err))
	assert.Equal(t, 0, h.manager.countJobs(dragonfly.JobTypeDeleteTask))
}

func TestDatasetReconciler_RemovedSourceIsPurged(t *testing.T) {
	purge := v1alpha1.DeletionPolicyPurge
	ds := fileDataset("ds", "https://example.com/a")
	ds.Spec.Sources = append(ds.Spec.Sources, v1alpha1.DataSource{Name: "extra", Type: v1alpha1.DataSourceTypeFile, URL: "https://example.com/extra"})
	ds.Spec.Lifecycle = &v1alpha1.LifecycleSpec{DeletionPolicy: &purge}
	h := newHarness(t, ds)

	h.reconcileOK("ds")
	h.reconcileOK("ds")
	assert.Equal(t, 2, h.manager.countJobs(dragonfly.JobTypePreheat), "one preheat job per source")
	for id, job := range h.manager.jobs {
		req := h.manager.requests[id-1].(managertypes.CreatePreheatJobRequest)
		h.manager.finish(job.ID, dragonfly.JobStateSuccess, preheatSuccess(req.Args.URLs...))
	}
	h.reconcileOK("ds")
	got := h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseReady, got.Status.Phase)
	require.Len(t, got.Status.Sources, 2)

	got.Spec.Sources = got.Spec.Sources[:1]
	require.NoError(t, h.client.Update(context.Background(), got))

	h.reconcileOK("ds")
	assert.Equal(t, 1, h.manager.countJobs(dragonfly.JobTypeDeleteTask))
	got = h.get("ds")
	require.Len(t, got.Status.Sources, 2, "the removed source stays until purged")
	assert.Equal(t, v1alpha1.DataSourceStatePurging, sourceStatus(got, "extra").State)
	assert.Equal(t, v1alpha1.DatasetPhasePurging, got.Status.Phase)

	h.manager.finishAll(dragonfly.JobTypeDeleteTask, dragonfly.JobStateSuccess, deleteTaskSuccess())
	h.reconcileOK("ds")
	h.reconcileOK("ds")
	got = h.get("ds")
	assert.Len(t, got.Status.Sources, 1)
	assert.Nil(t, sourceStatus(got, "extra"))
	assert.Equal(t, v1alpha1.DatasetPhaseReady, got.Status.Phase)
}

func TestDatasetReconciler_Suspend(t *testing.T) {
	ds := fileDataset("ds", "https://example.com/a")
	ds.Spec.Suspend = true
	h := newHarness(t, ds)

	h.reconcileOK("ds")
	h.reconcileOK("ds")
	assert.Equal(t, 0, h.manager.countJobs(dragonfly.JobTypePreheat))
	got := h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseSuspended, got.Status.Phase)
	assert.Equal(t, ReasonSuspended, condition(got, v1alpha1.ConditionReady).Reason)

	got.Spec.Suspend = false
	require.NoError(t, h.client.Update(context.Background(), got))
	h.reconcileOK("ds")
	assert.Equal(t, 1, h.manager.countJobs(dragonfly.JobTypePreheat))
}

func TestDatasetReconciler_ManagerRefAndSecrets(t *testing.T) {
	ds := fileDataset("ds", "https://example.com/a")
	ds.Spec.ManagerRef = &v1alpha1.ManagerRef{
		Endpoint:       "https://other-manager",
		TokenSecretRef: &corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "manager"}},
	}
	ds.Spec.Sources[0].AuthSecretRef = &corev1.LocalObjectReference{Name: "origin"}
	ds.Spec.Sources[0].HeadersSecretRef = &corev1.LocalObjectReference{Name: "headers"}
	ds.Spec.Sources[0].Headers = map[string]string{"X-Static": "1", "X-Override": "inline"}

	secrets := []client.Object{
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "manager", Namespace: "default"}, Data: map[string][]byte{"token": []byte("pat\n")}},
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "origin", Namespace: "default"}, Data: map[string][]byte{"username": []byte("u"), "password": []byte("p")}},
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "headers", Namespace: "default"}, Data: map[string][]byte{"X-Override": []byte("secret")}},
	}

	var seen dragonfly.Options
	h := newHarness(t, append(secrets, ds)...)
	h.cfg.Manager.Endpoint = ""
	h.cfg.Manager.Token = ""
	manager := h.manager
	h.reconciler.NewManagerClient = func(opts dragonfly.Options) (dragonfly.Client, error) {
		seen = opts
		return manager, nil
	}

	h.reconcileOK("ds")
	h.reconcileOK("ds")
	assert.Equal(t, "https://other-manager", seen.Endpoint)
	assert.Equal(t, "pat", seen.Token)

	require.Len(t, manager.requests, 1)
	req := manager.requests[0].(managertypes.CreatePreheatJobRequest)
	assert.Equal(t, "u", req.Args.Username)
	assert.Equal(t, "p", req.Args.Password)
	assert.Equal(t, map[string]string{"X-Static": "1", "X-Override": "secret"}, req.Args.Headers)
}

func TestDatasetReconciler_MissingSecretFailsSource(t *testing.T) {
	ds := fileDataset("ds", "https://example.com/a")
	ds.Spec.Sources[0].AuthSecretRef = &corev1.LocalObjectReference{Name: "missing"}
	retryLimit := int32(0)
	ds.Spec.Distribution.RetryLimit = &retryLimit
	h := newHarness(t, ds)

	h.reconcileOK("ds")
	h.reconcileOK("ds")
	got := h.get("ds")
	assert.Equal(t, 0, h.manager.countJobs(dragonfly.JobTypePreheat))
	assert.Equal(t, v1alpha1.DatasetPhaseFailed, got.Status.Phase)
	source := sourceStatus(got, "files")
	assert.Equal(t, v1alpha1.DataSourceStateFailed, source.State)
	assert.Contains(t, source.Message, "missing")
}

func TestDatasetReconciler_ManagerUnavailable(t *testing.T) {
	h := newHarness(t, fileDataset("ds", "https://example.com/a"))
	h.cfg.Manager.Endpoint = ""
	h.cfg.Manager.Token = ""

	h.reconcileOK("ds")
	_, err := h.reconcile("ds")
	require.Error(t, err)

	got := h.get("ds")
	ready := condition(got, v1alpha1.ConditionReady)
	require.NotNil(t, ready)
	assert.Equal(t, ReasonManagerUnavailable, ready.Reason)
	assert.Contains(t, h.drainEvents()[0], ReasonManagerUnavailable)
}

func TestDatasetReconciler_TransientManagerErrorIsRetried(t *testing.T) {
	h := newHarness(t, fileDataset("ds", "https://example.com/a"))
	h.reconcileOK("ds")

	h.manager.createErr = &dragonfly.APIError{StatusCode: 429, Message: "rate limit exceeded"}
	_, err := h.reconcile("ds")
	require.Error(t, err)
	got := h.get("ds")
	source := sourceStatus(got, "files")
	assert.Equal(t, v1alpha1.DataSourceStateDistributing, source.State)
	assert.Equal(t, int32(0), source.Retries, "a transient error does not consume a retry")
	require.NotNil(t, source.Operation)
	assert.Empty(t, source.Operation.Jobs)

	h.manager.createErr = nil
	h.reconcileOK("ds")
	assert.Equal(t, 1, h.manager.countJobs(dragonfly.JobTypePreheat))
}

func TestDatasetReconciler_InvalidSpec(t *testing.T) {
	ds := fileDataset("ds")
	h := newHarness(t, ds)

	h.reconcileOK("ds")
	h.reconcileOK("ds")
	got := h.get("ds")
	assert.Equal(t, v1alpha1.DatasetPhaseFailed, got.Status.Phase)
	assert.Equal(t, ReasonInvalidSpec, condition(got, v1alpha1.ConditionReady).Reason)
	assert.Equal(t, 0, h.manager.countJobs(dragonfly.JobTypePreheat))
}

func TestDatasetReconciler_InvalidLifecycle(t *testing.T) {
	ds := fileDataset("ds", "https://example.com/a")
	ds.Spec.Lifecycle = &v1alpha1.LifecycleSpec{RefreshSchedule: stringPtr("every other tuesday")}
	h := newHarness(t, ds)

	h.reconcileOK("ds")
	_, err := h.reconcile("ds")
	require.Error(t, err)
	got := h.get("ds")
	assert.Equal(t, ReasonInvalidLifecycle, condition(got, v1alpha1.ConditionReady).Reason)
}

func TestDatasetReconciler_ImageSource(t *testing.T) {
	ds := &v1alpha1.Dataset{
		ObjectMeta: metav1.ObjectMeta{Name: "img", Namespace: "default"},
		Spec: v1alpha1.DatasetSpec{
			Sources: []v1alpha1.DataSource{{
				Name:                        "image",
				Type:                        v1alpha1.DataSourceTypeImage,
				URL:                         "https://registry.example.com/v2/library/nginx/manifests/latest",
				Platform:                    "linux/amd64",
				EnableTaskIDBasedBlobDigest: true,
			}},
			Lifecycle: &v1alpha1.LifecycleSpec{TTL: duration(time.Hour)},
		},
	}
	h := newHarness(t, ds)

	h.reconcileOK("img")
	h.reconcileOK("img")
	require.Len(t, h.manager.requests, 1)
	req := h.manager.requests[0].(managertypes.CreatePreheatJobRequest)
	assert.Equal(t, "image", req.Args.Type)
	assert.Equal(t, ds.Spec.Sources[0].URL, req.Args.URL)
	assert.Equal(t, "linux/amd64", req.Args.Platform)
	assert.True(t, req.Args.EnableTaskIDBasedBlobDigest)

	layers := []string{
		"https://registry.example.com/v2/library/nginx/blobs/sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"https://registry.example.com/v2/library/nginx/blobs/sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	}
	h.manager.finishAll(dragonfly.JobTypePreheat, dragonfly.JobStateSuccess, preheatSuccess(layers...))
	h.reconcileOK("img")
	got := h.get("img")
	source := sourceStatus(got, "image")
	require.Len(t, source.Tasks, 2, "layers reported by the preheat become tasks")
	assert.Equal(t, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", source.Tasks[0].ID, "blob digest based task id")

	// Expiry purges every layer by its digest-based ID.
	h.advance(time.Hour)
	h.reconcileOK("img")
	assert.Equal(t, 2, h.manager.countJobs(dragonfly.JobTypeDeleteTask))
	deleteReq := h.manager.requests[1].(managertypes.CreateDeleteTaskJobRequest)
	assert.Equal(t, source.Tasks[0].ID, deleteReq.Args.TaskID)
}

func TestDatasetReconciler_PurgeBatchesJobs(t *testing.T) {
	urls := make([]string, 0, maxJobsPerReconcile+5)
	for i := 0; i < maxJobsPerReconcile+5; i++ {
		urls = append(urls, fmt.Sprintf("https://example.com/%d", i))
	}

	ds := fileDataset("ds", urls...)
	ds.Spec.Lifecycle = &v1alpha1.LifecycleSpec{TTL: duration(time.Hour)}
	h := newHarness(t, ds)
	h.distribute("ds", urls...)

	h.advance(time.Hour)
	result := h.reconcileOK("ds")
	assert.Equal(t, maxJobsPerReconcile, h.manager.countJobs(dragonfly.JobTypeDeleteTask), "first batch")
	assert.Equal(t, minRequeue, result.RequeueAfter, "the rest is created right away")

	h.reconcileOK("ds")
	assert.Equal(t, len(urls), h.manager.countJobs(dragonfly.JobTypeDeleteTask), "second batch")
}

func TestDataLifecyclePolicyReconciler(t *testing.T) {
	purge := v1alpha1.DeletionPolicyPurge
	hot := policy("hot", 10, map[string]string{"tier": "hot"}, v1alpha1.LifecycleSpec{TTL: duration(time.Hour), DeletionPolicy: &purge})
	all := policy("all", 0, nil, v1alpha1.LifecycleSpec{TTL: duration(2 * time.Hour)})
	bad := policy("bad", 0, nil, v1alpha1.LifecycleSpec{RefreshSchedule: stringPtr("nope")})

	cold := fileDataset("cold", "https://example.com/c")
	cold.Labels = map[string]string{"tier": "cold"}

	h := newHarness(t, fileDataset("ds", "https://example.com/a"), cold, &hot, &all, &bad)
	r := &DataLifecyclePolicyReconciler{Client: h.client, Scheme: h.reconciler.Scheme}

	for _, name := range []string{"hot", "all", "bad"} {
		_, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: name}})
		require.NoError(t, err)
	}

	get := func(name string) *v1alpha1.DataLifecyclePolicy {
		var p v1alpha1.DataLifecyclePolicy
		require.NoError(t, h.client.Get(context.Background(), types.NamespacedName{Namespace: "default", Name: name}, &p))
		return &p
	}

	assert.Equal(t, int32(1), get("hot").Status.MatchedDatasets, "hot governs the hot dataset")
	assert.Equal(t, int32(1), get("all").Status.MatchedDatasets, "all governs the cold dataset only")
	assert.Equal(t, metav1.ConditionTrue, meta.FindStatusCondition(get("hot").Status.Conditions, ConditionPolicyValid).Status)

	badStatus := get("bad").Status
	assert.Equal(t, int32(0), badStatus.MatchedDatasets)
	assert.Equal(t, metav1.ConditionFalse, meta.FindStatusCondition(badStatus.Conditions, ConditionPolicyValid).Status)
}
