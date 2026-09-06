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
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"

	"d7y.io/dragonfly/v2/datacontroller/apis/data/v1alpha1"
	"d7y.io/dragonfly/v2/datacontroller/dragonfly"
	managertypes "d7y.io/dragonfly/v2/manager/types"
	"d7y.io/dragonfly/v2/pkg/idgen"
	pkghttp "d7y.io/dragonfly/v2/pkg/net/http"
)

// Secret keys read by the controller.
const (
	SecretKeyUsername        = "username"
	SecretKeyPassword        = "password"
	SecretKeyToken           = "token"
	SecretKeyAccessKeyID     = "accessKeyID"
	SecretKeyAccessKeySecret = "accessKeySecret"
	SecretKeySessionToken    = "sessionToken"
	SecretKeySecurityToken   = "securityToken"
	SecretKeyCredentialPath  = "credentialPath"
	SecretKeyDelegationToken = "delegationToken"
)

// operationError is a failure of an operation that retrying will not fix, such as a
// rejected request. It is recorded on the source instead of being returned from the
// reconciler.
type operationError struct {
	err error
}

func (e *operationError) Error() string { return e.err.Error() }
func (e *operationError) Unwrap() error { return e.err }

// isRetryable reports whether an error talking to the manager is worth retrying on
// the reconciler's backoff instead of counting as an operation failure.
func isRetryable(err error) bool {
	var apiErr *dragonfly.APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests ||
			apiErr.StatusCode == http.StatusRequestTimeout ||
			apiErr.StatusCode >= http.StatusInternalServerError
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	return errors.Is(err, context.DeadlineExceeded)
}

// secretReader reads Secrets in the namespace of a dataset.
type secretReader struct {
	client    client.Reader
	namespace string
	cache     map[string]*corev1.Secret
}

func newSecretReader(c client.Reader, namespace string) *secretReader {
	return &secretReader{client: c, namespace: namespace, cache: map[string]*corev1.Secret{}}
}

func (s *secretReader) get(ctx context.Context, name string) (*corev1.Secret, error) {
	if secret, ok := s.cache[name]; ok {
		return secret, nil
	}

	var secret corev1.Secret
	if err := s.client.Get(ctx, types.NamespacedName{Namespace: s.namespace, Name: name}, &secret); err != nil {
		return nil, fmt.Errorf("get secret %s/%s: %w", s.namespace, name, err)
	}

	s.cache[name] = &secret
	return &secret, nil
}

// value returns a key of a secret, or an error when the key is missing and required.
func (s *secretReader) value(ctx context.Context, name, key string, required bool) (string, error) {
	secret, err := s.get(ctx, name)
	if err != nil {
		return "", err
	}

	data, ok := secret.Data[key]
	if !ok || len(data) == 0 {
		if required {
			return "", fmt.Errorf("secret %s/%s has no %q key", s.namespace, name, key)
		}

		return "", nil
	}

	return strings.TrimSpace(string(data)), nil
}

// sourceURLs returns the URLs a source refers to.
func sourceURLs(source *v1alpha1.DataSource) []string {
	if source.Type == v1alpha1.DataSourceTypeImage {
		return []string{source.URL}
	}

	urls := make([]string, 0, len(source.URLs)+1)
	if source.URL != "" {
		urls = append(urls, source.URL)
	}

	for _, u := range source.URLs {
		if u != "" {
			urls = append(urls, u)
		}
	}

	return uniqueStrings(urls)
}

// validateSource checks the parts of a source that the CRD schema cannot express.
func validateSource(source *v1alpha1.DataSource) error {
	urls := sourceURLs(source)
	switch source.Type {
	case v1alpha1.DataSourceTypeImage:
		if source.URL == "" {
			return errors.New("image sources require url")
		}
	case v1alpha1.DataSourceTypeFile, "":
		if len(urls) == 0 {
			return errors.New("file sources require url or urls")
		}
	default:
		return fmt.Errorf("unknown source type %q", source.Type)
	}

	return nil
}

// filteredQueryParams returns the filtered query params used for task IDs, falling
// back to the manager default so IDs computed here match the ones the schedulers use.
func filteredQueryParams(source *v1alpha1.DataSource) string {
	if source.FilteredQueryParams != "" {
		return source.FilteredQueryParams
	}

	return pkghttp.RawDefaultFilteredQueryParams
}

// taskID computes the task ID of a URL the same way the schedulers do when preheating.
func taskID(source *v1alpha1.DataSource, url string) (string, error) {
	return idgen.TaskIDV2(url, source.PieceLength, source.Tag, source.Application,
		idgen.ParseFilteredQueryParams(filteredQueryParams(source)), "", source.EnableTaskIDBasedBlobDigest)
}

// buildPreheatRequest builds the preheat job request of a source.
func buildPreheatRequest(ctx context.Context, secrets *secretReader, dataset *v1alpha1.Dataset, source *v1alpha1.DataSource) (managertypes.CreatePreheatJobRequest, error) {
	distribution := dataset.Spec.Distribution
	args := managertypes.PreheatArgs{
		Type:                        string(source.Type),
		Tag:                         source.Tag,
		Application:                 source.Application,
		Platform:                    source.Platform,
		PieceLength:                 source.PieceLength,
		FilteredQueryParams:         source.FilteredQueryParams,
		Scope:                       string(distribution.Scope),
		IPs:                         distribution.IPs,
		Percentage:                  distribution.Percentage,
		Count:                       distribution.Count,
		EnableTaskIDBasedBlobDigest: source.EnableTaskIDBasedBlobDigest,
	}

	if args.Type == "" {
		args.Type = string(v1alpha1.DataSourceTypeFile)
	}

	if source.Type == v1alpha1.DataSourceTypeImage {
		args.URL = source.URL
	} else {
		args.URLs = sourceURLs(source)
	}

	if distribution.ConcurrentTaskCount != nil {
		args.ConcurrentTaskCount = *distribution.ConcurrentTaskCount
	}

	if distribution.ConcurrentPeerCount != nil {
		args.ConcurrentPeerCount = *distribution.ConcurrentPeerCount
	}

	if distribution.Timeout != nil {
		args.Timeout = distribution.Timeout.Duration
	}

	headers, err := resolveHeaders(ctx, secrets, source)
	if err != nil {
		return managertypes.CreatePreheatJobRequest{}, err
	}
	args.Headers = headers

	if source.AuthSecretRef != nil {
		username, err := secrets.value(ctx, source.AuthSecretRef.Name, SecretKeyUsername, true)
		if err != nil {
			return managertypes.CreatePreheatJobRequest{}, err
		}

		password, err := secrets.value(ctx, source.AuthSecretRef.Name, SecretKeyPassword, true)
		if err != nil {
			return managertypes.CreatePreheatJobRequest{}, err
		}

		args.Username = username
		args.Password = password
	}

	if source.ObjectStorage != nil {
		objectStorage, err := resolveObjectStorage(ctx, secrets, source.ObjectStorage)
		if err != nil {
			return managertypes.CreatePreheatJobRequest{}, err
		}

		args.ObjectStorage = objectStorage
	}

	if source.HDFS != nil {
		hdfs, err := resolveHDFS(ctx, secrets, source.HDFS)
		if err != nil {
			return managertypes.CreatePreheatJobRequest{}, err
		}

		args.Hdfs = hdfs
	}

	return managertypes.CreatePreheatJobRequest{
		BIO:                 jobDescription(dataset, source, v1alpha1.OperationTypeDistribute),
		Type:                dragonfly.JobTypePreheat,
		Args:                args,
		SchedulerClusterIDs: distribution.SchedulerClusterIDs,
	}, nil
}

// buildGetTaskRequest builds the get_task job request of one task of a source.
func buildGetTaskRequest(dataset *v1alpha1.Dataset, source *v1alpha1.DataSource, task *v1alpha1.TaskStatus) (managertypes.CreateGetTaskJobRequest, error) {
	id := task.ID
	if id == "" {
		var err error
		if id, err = taskID(source, task.URL); err != nil {
			return managertypes.CreateGetTaskJobRequest{}, fmt.Errorf("task id of %s: %w", task.URL, err)
		}
	}

	distribution := dataset.Spec.Distribution
	args := managertypes.GetTaskArgs{TaskID: id}
	if distribution.ConcurrentPeerCount != nil {
		args.ConcurrentPeerCount = *distribution.ConcurrentPeerCount
	}

	if distribution.Timeout != nil {
		args.Timeout = distribution.Timeout.Duration
	}

	return managertypes.CreateGetTaskJobRequest{
		BIO:                 jobDescription(dataset, source, v1alpha1.OperationTypeVerify),
		Type:                dragonfly.JobTypeGetTask,
		Args:                args,
		SchedulerClusterIDs: distribution.SchedulerClusterIDs,
	}, nil
}

// buildDeleteTaskRequest builds the delete_task job request of one task of a source.
func buildDeleteTaskRequest(dataset *v1alpha1.Dataset, source *v1alpha1.DataSource, task *v1alpha1.TaskStatus) (managertypes.CreateDeleteTaskJobRequest, error) {
	id := task.ID
	if id == "" {
		var err error
		if id, err = taskID(source, task.URL); err != nil {
			return managertypes.CreateDeleteTaskJobRequest{}, fmt.Errorf("task id of %s: %w", task.URL, err)
		}
	}

	distribution := dataset.Spec.Distribution
	args := managertypes.DeleteTaskArgs{TaskID: id}
	if distribution.Timeout != nil {
		args.Timeout = distribution.Timeout.Duration
	}

	return managertypes.CreateDeleteTaskJobRequest{
		BIO:                 jobDescription(dataset, source, v1alpha1.OperationTypePurge),
		Type:                dragonfly.JobTypeDeleteTask,
		Args:                args,
		SchedulerClusterIDs: distribution.SchedulerClusterIDs,
	}, nil
}

func jobDescription(dataset *v1alpha1.Dataset, source *v1alpha1.DataSource, op v1alpha1.OperationType) string {
	return fmt.Sprintf("dataset %s/%s source %s: %s", dataset.Namespace, dataset.Name, source.Name, strings.ToLower(string(op)))
}

func resolveHeaders(ctx context.Context, secrets *secretReader, source *v1alpha1.DataSource) (map[string]string, error) {
	headers := map[string]string{}
	for k, v := range source.Headers {
		headers[k] = v
	}

	if source.HeadersSecretRef != nil {
		secret, err := secrets.get(ctx, source.HeadersSecretRef.Name)
		if err != nil {
			return nil, err
		}

		for k, v := range secret.Data {
			headers[k] = string(v)
		}
	}

	if len(headers) == 0 {
		return nil, nil
	}

	return headers, nil
}

func resolveObjectStorage(ctx context.Context, secrets *secretReader, spec *v1alpha1.ObjectStorageSpec) (*commonv2.ObjectStorage, error) {
	objectStorage := &commonv2.ObjectStorage{}
	if spec.Region != "" {
		objectStorage.Region = &spec.Region
	}

	if spec.Endpoint != "" {
		objectStorage.Endpoint = &spec.Endpoint
	}

	if spec.PredefinedACL != "" {
		objectStorage.PredefinedAcl = &spec.PredefinedACL
	}

	if spec.InsecureSkipVerify {
		objectStorage.InsecureSkipVerify = &spec.InsecureSkipVerify
	}

	if spec.CredentialsSecretRef != nil {
		name := spec.CredentialsSecretRef.Name
		optional := []struct {
			key string
			dst **string
		}{
			{SecretKeyAccessKeyID, &objectStorage.AccessKeyId},
			{SecretKeyAccessKeySecret, &objectStorage.AccessKeySecret},
			{SecretKeySessionToken, &objectStorage.SessionToken},
			{SecretKeySecurityToken, &objectStorage.SecurityToken},
			{SecretKeyCredentialPath, &objectStorage.CredentialPath},
		}

		for _, field := range optional {
			value, err := secrets.value(ctx, name, field.key, false)
			if err != nil {
				return nil, err
			}

			if value != "" {
				v := value
				*field.dst = &v
			}
		}

		if objectStorage.AccessKeyId == nil && objectStorage.CredentialPath == nil {
			return nil, fmt.Errorf("secret %s/%s needs %q or %q", secrets.namespace, name, SecretKeyAccessKeyID, SecretKeyCredentialPath)
		}
	}

	return objectStorage, nil
}

func resolveHDFS(ctx context.Context, secrets *secretReader, spec *v1alpha1.HDFSSpec) (*commonv2.HDFS, error) {
	hdfs := &commonv2.HDFS{}
	if spec.DelegationTokenSecretRef != nil {
		token, err := secrets.value(ctx, spec.DelegationTokenSecretRef.Name, SecretKeyDelegationToken, true)
		if err != nil {
			return nil, err
		}

		hdfs.DelegationToken = &token
	}

	return hdfs, nil
}

// distributeOutcome is what a finished preheat job means for a source.
type distributeOutcome struct {
	tasks   []v1alpha1.TaskStatus
	failed  bool
	message string
}

// evaluateDistribute turns a terminal preheat job into tasks and a verdict.
func evaluateDistribute(job *dragonfly.Job, source *v1alpha1.DataSource) distributeOutcome {
	if job.State == dragonfly.JobStateFailure {
		_, errs, _ := dragonfly.PreheatResults(job)
		return distributeOutcome{failed: true, message: joinMessages("preheat job failed", errs)}
	}

	results, errs, err := dragonfly.PreheatResults(job)
	if err != nil {
		return distributeOutcome{failed: true, message: fmt.Sprintf("preheat job %d: %v", job.ID, err)}
	}

	var (
		successURLs []string
		failureURLs []string
		failures    []string
	)

	for _, result := range results {
		for _, task := range result.SuccessTasks {
			successURLs = append(successURLs, task.URL)
		}

		for _, task := range result.FailureTasks {
			failureURLs = append(failureURLs, task.URL)
			failures = append(failures, fmt.Sprintf("%s on %s: %s", task.URL, task.IP, task.Description))
		}
	}

	successURLs = uniqueStrings(successURLs)
	failureURLs = uniqueStrings(failureURLs)

	if len(successURLs) == 0 {
		if len(failureURLs) > 0 || len(errs) > 0 {
			return distributeOutcome{failed: true, message: joinMessages("no task was preheated", append(errs, failures...))}
		}

		// Older schedulers do not report the preheated tasks; fall back to the URLs of the spec.
		successURLs = sourceURLs(source)
	}

	tasks := make([]v1alpha1.TaskStatus, 0, len(successURLs))
	for _, u := range successURLs {
		task := v1alpha1.TaskStatus{URL: u}
		if id, err := taskID(source, u); err == nil {
			task.ID = id
		}

		tasks = append(tasks, task)
	}

	message := fmt.Sprintf("%d task(s) distributed", len(tasks))
	if len(failureURLs) > 0 {
		message = fmt.Sprintf("%s, %d task(s) failed on some peers: %s", message, len(failureURLs), truncate(strings.Join(failures, "; "), 512))
	}

	return distributeOutcome{tasks: tasks, message: message}
}

// evaluateVerify updates the peer count of a task from a terminal get_task job. It
// returns false when the job failed.
func evaluateVerify(job *dragonfly.Job, task *v1alpha1.TaskStatus) (bool, string) {
	if job.State == dragonfly.JobStateFailure {
		_, errs, _ := dragonfly.GetTaskResults(job)
		task.PeerCount = nil
		return false, joinMessages(fmt.Sprintf("verification of %s failed", task.URL), errs)
	}

	results, _, err := dragonfly.GetTaskResults(job)
	if err != nil {
		task.PeerCount = nil
		return false, fmt.Sprintf("verification of %s: %v", task.URL, err)
	}

	var count int32
	for _, result := range results {
		count += int32(len(result.Peers))
	}

	task.PeerCount = &count
	return true, ""
}

// evaluatePurge reports whether a terminal delete_task job removed the task.
func evaluatePurge(job *dragonfly.Job, task *v1alpha1.TaskStatus) (bool, string) {
	if job.State == dragonfly.JobStateFailure {
		_, errs, _ := dragonfly.DeleteTaskResults(job)
		return false, joinMessages(fmt.Sprintf("purge of %s failed", task.URL), errs)
	}

	results, _, err := dragonfly.DeleteTaskResults(job)
	if err != nil {
		return false, fmt.Sprintf("purge of %s: %v", task.URL, err)
	}

	var failures []string
	for _, result := range results {
		for _, peer := range result.FailureTasks {
			failures = append(failures, fmt.Sprintf("%s: %s", peer.IP, peer.Description))
		}
	}

	if len(failures) > 0 {
		return false, joinMessages(fmt.Sprintf("purge of %s failed on %d peer(s)", task.URL, len(failures)), failures)
	}

	return true, ""
}

func joinMessages(prefix string, details []string) string {
	if len(details) == 0 {
		return prefix
	}

	return truncate(prefix+": "+strings.Join(details, "; "), 1024)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}

	return s[:n] + "..."
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}

		seen[v] = struct{}{}
		result = append(result, v)
	}

	return result
}

// retryBackoff returns the delay before the n-th retry (1-based).
func retryBackoff(base, max time.Duration, n int32) time.Duration {
	if n < 1 {
		n = 1
	}

	delay := base
	for i := int32(1); i < n; i++ {
		delay *= 2
		if delay >= max {
			return max
		}
	}

	if delay > max {
		return max
	}

	return delay
}

// sortedTaskURLs returns the URLs of the tasks in a stable order.
func sortedTaskURLs(tasks []v1alpha1.TaskStatus) []string {
	urls := make([]string, 0, len(tasks))
	for _, task := range tasks {
		urls = append(urls, task.URL)
	}

	sort.Strings(urls)
	return urls
}
