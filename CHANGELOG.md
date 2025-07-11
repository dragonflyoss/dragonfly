<a name="unreleased"></a>
## [Unreleased]

### Chore
- update SECURITY-INSIGHTS with contribution policy details 📜 ([#4198](https://github.com/CooooolFrog/dragonfly/issues/4198))
- **SECURITY-INSIGHTS.yml:** change the address of dragonfly repository ([#4199](https://github.com/CooooolFrog/dragonfly/issues/4199))
- **deps:** bump aquasecurity/trivy-action from 0.31.0 to 0.32.0 ([#4188](https://github.com/CooooolFrog/dragonfly/issues/4188))
- **deps:** bump kentaro-m/auto-assign-action from a558e3a7a389eae6d782d4f5ad0aaa5e36ec2e39 to db77befe126df3163e8ee988af621569f7f2b82f ([#4192](https://github.com/CooooolFrog/dragonfly/issues/4192))
- **deps:** bump github.com/go-playground/validator/v10 from 10.26.0 to 10.27.0 ([#4191](https://github.com/CooooolFrog/dragonfly/issues/4191))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.60.0 to 0.62.0 ([#4190](https://github.com/CooooolFrog/dragonfly/issues/4190))
- **deps:** bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc from 1.36.0 to 1.37.0 ([#4189](https://github.com/CooooolFrog/dragonfly/issues/4189))

### Docs
- add security insights and update community details 📚 ([#4193](https://github.com/CooooolFrog/dragonfly/issues/4193))


<a name="v2.3.1-beta.0"></a>
## [v2.3.1-beta.0] - 2025-07-08
### Chore
- bump console version to v0.3.0 ([#4196](https://github.com/CooooolFrog/dragonfly/issues/4196))
- bump client-rs to v1.0.1 ([#4174](https://github.com/CooooolFrog/dragonfly/issues/4174))

### Docs
- update README with community engagement details 📚 ([#4184](https://github.com/CooooolFrog/dragonfly/issues/4184))

### Feat
- Add support for selecting specific peer IPs for preheating in job scheduler ([#4194](https://github.com/CooooolFrog/dragonfly/issues/4194))
- support preheat for multiple urls ([#4185](https://github.com/CooooolFrog/dragonfly/issues/4185))
- Change Result column type to longtext in Job model to fix data too long error ([#4181](https://github.com/CooooolFrog/dragonfly/issues/4181))
- add gc task to clean up inactive seed peers ([#4176](https://github.com/CooooolFrog/dragonfly/issues/4176))
- add the gc task to clean up the inactive schedulers ([#4175](https://github.com/CooooolFrog/dragonfly/issues/4175))

### Fix
- revert the remove of ip and port in mamager rpc server ([#4180](https://github.com/CooooolFrog/dragonfly/issues/4180))
- populate the User from UserID when get application and job ([#4178](https://github.com/CooooolFrog/dragonfly/issues/4178))


<a name="v2.3.0"></a>
## [v2.3.0] - 2025-07-01
### Chore
- **deps:** bump client-rs to v1.0.0 and bump console to v0.2.8 ([#4172](https://github.com/CooooolFrog/dragonfly/issues/4172))
- **deps:** bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp from 1.36.0 to 1.37.0 ([#4167](https://github.com/CooooolFrog/dragonfly/issues/4167))
- **deps:** bump github/codeql-action from 3.28.19 to 3.29.2 ([#4165](https://github.com/CooooolFrog/dragonfly/issues/4165))
- **deps:** bump go.opentelemetry.io/otel/exporters/otlp/otlptrace from 1.36.0 to 1.37.0 ([#4169](https://github.com/CooooolFrog/dragonfly/issues/4169))
- **deps:** bump github.com/redis/go-redis/v9 from 9.10.0 to 9.11.0 ([#4170](https://github.com/CooooolFrog/dragonfly/issues/4170))
- **deps:** bump client-rs to v0.2.42 ([#4164](https://github.com/CooooolFrog/dragonfly/issues/4164))

### Fix
- correct the task id of job gc ([#4163](https://github.com/CooooolFrog/dragonfly/issues/4163))


<a name="v2.2.4-rc.2"></a>
## [v2.2.4-rc.2] - 2025-06-25
### Chore
- **deps:** bump mheap/github-action-required-labels from 5.5.0 to 5.5.1 ([#4149](https://github.com/CooooolFrog/dragonfly/issues/4149))
- **deps:** bump sigstore/cosign-installer from 3.8.2 to 3.9.1 ([#4150](https://github.com/CooooolFrog/dragonfly/issues/4150))
- **deps:** bump kentaro-m/auto-assign-action from e6c4932f995626505d2a5d85401b0319753caa11 to a558e3a7a389eae6d782d4f5ad0aaa5e36ec2e39 ([#4152](https://github.com/CooooolFrog/dragonfly/issues/4152))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.60.0 to 0.61.0 ([#4154](https://github.com/CooooolFrog/dragonfly/issues/4154))
- **deps:** bump docker/setup-buildx-action from 3.11.0 to 3.11.1 ([#4155](https://github.com/CooooolFrog/dragonfly/issues/4155))
- **goreleaser:** add release footer with contributors and changelog links ([#4161](https://github.com/CooooolFrog/dragonfly/issues/4161))

### Feat
- support HTTPS protocol for otel ([#4159](https://github.com/CooooolFrog/dragonfly/issues/4159))


<a name="v2.2.4-rc.1"></a>
## [v2.2.4-rc.1] - 2025-06-24
### Chore
- **client-rs:** bump client-rs to v0.2.37 and bump console to v0.2.6 ([#4143](https://github.com/CooooolFrog/dragonfly/issues/4143))
- **deps:** bump google.golang.org/grpc from 1.72.2 to 1.73.0 ([#4151](https://github.com/CooooolFrog/dragonfly/issues/4151))
- **deps:** bump github.com/jellydator/ttlcache/v3 from 3.3.0 to 3.4.0 ([#4153](https://github.com/CooooolFrog/dragonfly/issues/4153))
- **deps:** bump golang.org/x/time from 0.11.0 to 0.12.0 ([#4148](https://github.com/CooooolFrog/dragonfly/issues/4148))

### Docs
- **README.md:** add FOSSA Status badge ([#4144](https://github.com/CooooolFrog/dragonfly/issues/4144))

### Feat
- add logLevel and remove verbose ([#4157](https://github.com/CooooolFrog/dragonfly/issues/4157))
- remove unique index uk_casbin_rule from CasbinRule model to resolve permission issue ([#4156](https://github.com/CooooolFrog/dragonfly/issues/4156))
- add --log-level support for manager and scheduler ([#4146](https://github.com/CooooolFrog/dragonfly/issues/4146))


<a name="v2.2.4-rc.0"></a>
## [v2.2.4-rc.0] - 2025-06-19
### Chore
- **client-rs:** bump client-rs to v0.2.37 and bump console to v0.2.6
- **deps:** bump docker/setup-buildx-action from 3.10.0 to 3.11.0 ([#4134](https://github.com/CooooolFrog/dragonfly/issues/4134))
- **deps:** bump github.com/redis/go-redis/v9 from 9.7.3 to 9.10.0 ([#4132](https://github.com/CooooolFrog/dragonfly/issues/4132))
- **deps:** bump golang.org/x/sync from 0.14.0 to 0.15.0 ([#4120](https://github.com/CooooolFrog/dragonfly/issues/4120))
- **deps:** bump github/codeql-action from 3.28.18 to 3.28.19 ([#4118](https://github.com/CooooolFrog/dragonfly/issues/4118))
- **deps:** bump aquasecurity/trivy-action from 0.30.0 to 0.31.0 ([#4117](https://github.com/CooooolFrog/dragonfly/issues/4117))
- **deps:** bump kentaro-m/auto-assign-action from 7ae38e468e64dec0af17820972bc4915aa511ec2 to e6c4932f995626505d2a5d85401b0319753caa11 ([#4116](https://github.com/CooooolFrog/dragonfly/issues/4116))

### Docs
- Update OWNERS.md to reference emeritus_members.yaml for emeritus maintainers ([#4128](https://github.com/CooooolFrog/dragonfly/issues/4128))
- Add Cheng Ming to Community Managers ([#4127](https://github.com/CooooolFrog/dragonfly/issues/4127))

### Feat
- support http protocol for otel exporter ([#4142](https://github.com/CooooolFrog/dragonfly/issues/4142))
- support get info of the preheating image ([#4139](https://github.com/CooooolFrog/dragonfly/issues/4139))
- add headers for tracing ([#4129](https://github.com/CooooolFrog/dragonfly/issues/4129))
- add the type in the gc job args ([#4126](https://github.com/CooooolFrog/dragonfly/issues/4126))


<a name="v2.2.4-beta.1"></a>
## [v2.2.4-beta.1] - 2025-06-06
### Chore
- **client-rs:** bump client-rs to v0.2.34 ([#4114](https://github.com/CooooolFrog/dragonfly/issues/4114))
- **deps:** bump ossf/scorecard-action from 2.4.1 to 2.4.2 ([#4106](https://github.com/CooooolFrog/dragonfly/issues/4106))
- **deps:** bump docker/build-push-action from 6.16.0 to 6.18.0 ([#4105](https://github.com/CooooolFrog/dragonfly/issues/4105))
- **deps:** bump github.com/gin-gonic/gin from 1.10.0 to 1.10.1 ([#4110](https://github.com/CooooolFrog/dragonfly/issues/4110))
- **deps:** bump golang.org/x/crypto from 0.37.0 to 0.38.0 ([#4109](https://github.com/CooooolFrog/dragonfly/issues/4109))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.35.0 to 1.36.0 ([#4111](https://github.com/CooooolFrog/dragonfly/issues/4111))
- **deps:** bump github.com/go-jose/go-jose/v4 from 4.0.4 to 4.0.5 ([#4104](https://github.com/CooooolFrog/dragonfly/issues/4104))

### Feat
- update validate for job's handler ([#4113](https://github.com/CooooolFrog/dragonfly/issues/4113))


<a name="v2.2.4-beta.0"></a>
## [v2.2.4-beta.0] - 2025-05-29
### Chore
- bump client to v0.2.23 ([#4103](https://github.com/CooooolFrog/dragonfly/issues/4103))
- **README.md:** add maintainer google groups for communication channels and remove discussion group ([#4093](https://github.com/CooooolFrog/dragonfly/issues/4093))
- **deps:** bump google.golang.org/grpc from 1.71.1 to 1.72.2 ([#4097](https://github.com/CooooolFrog/dragonfly/issues/4097))
- **deps:** bump github.com/huaweicloud/huaweicloud-sdk-go-obs from 3.25.3+incompatible to 3.25.4+incompatible ([#4101](https://github.com/CooooolFrog/dragonfly/issues/4101))
- **deps:** bump github.com/gomodule/redigo from 1.8.10-0.20230511231101-78e255f9bd2a to 1.9.2 ([#4098](https://github.com/CooooolFrog/dragonfly/issues/4098))

### Docs
- **host.go:** Update NewHost comment for clarity ([#4092](https://github.com/CooooolFrog/dragonfly/issues/4092))

### Feat
- **manager:** optimizing tracing configuration ([#4102](https://github.com/CooooolFrog/dragonfly/issues/4102))


<a name="v2.2.3"></a>
## [v2.2.3] - 2025-05-21
### Chore
- **deps:** bump codecov/codecov-action from 5.4.2 to 5.4.3 ([#4077](https://github.com/CooooolFrog/dragonfly/issues/4077))
- **deps:** bump github/codeql-action from 3.28.17 to 3.28.18 ([#4076](https://github.com/CooooolFrog/dragonfly/issues/4076))

### Feat
- set poolSize and poolTimeout for manager ([#4089](https://github.com/CooooolFrog/dragonfly/issues/4089))
- add the API for trigger the GC manually ([#4074](https://github.com/CooooolFrog/dragonfly/issues/4074))
- update Redis connection pool configuration ([#4086](https://github.com/CooooolFrog/dragonfly/issues/4086))
- update machinery to v1.10.14 and adjust Redis config ([#4075](https://github.com/CooooolFrog/dragonfly/issues/4075))
- **manager:** add ratelimiter for job oapi ([#4088](https://github.com/CooooolFrog/dragonfly/issues/4088))


<a name="v2.2.3-rc.4"></a>
## [v2.2.3-rc.4] - 2025-05-16
### Comments
- **router:** add Init function to initialize gin engine ([#4071](https://github.com/CooooolFrog/dragonfly/issues/4071))

### Feat
- Update machinery to v1.10.13 and optimize Redis connection pool for weak network conditions ([#4072](https://github.com/CooooolFrog/dragonfly/issues/4072))
- implement distributed rate limiting by redis_rate/v10 ([#4068](https://github.com/CooooolFrog/dragonfly/issues/4068))
- add the storage pkg for stat task ([#4064](https://github.com/CooooolFrog/dragonfly/issues/4064))
- **manager:** add cleanup of associated tables in job garbage collection ([#4069](https://github.com/CooooolFrog/dragonfly/issues/4069))

### Fix
- add the authorization for get configs API ([#4067](https://github.com/CooooolFrog/dragonfly/issues/4067))


<a name="v2.2.3-rc.3"></a>
## [v2.2.3-rc.3] - 2025-05-14
### Feat
- support preheat by peer count ([#4066](https://github.com/CooooolFrog/dragonfly/issues/4066))
- add task type to ID generation for Standard and PersistentCache tasks ([#4065](https://github.com/CooooolFrog/dragonfly/issues/4065))


<a name="v2.2.3-rc.2"></a>
## [v2.2.3-rc.2] - 2025-05-13
### Chore
- **client-rs:** bump client-rs to v0.2.27 ([#4055](https://github.com/CooooolFrog/dragonfly/issues/4055))
- **deps:** bump golang.org/x/sync from 0.13.0 to 0.14.0 ([#4059](https://github.com/CooooolFrog/dragonfly/issues/4059))
- **deps:** bump actions/setup-go from 5.4.0 to 5.5.0 ([#4054](https://github.com/CooooolFrog/dragonfly/issues/4054))
- **deps:** bump golang.org/x/oauth2 from 0.29.0 to 0.30.0 ([#4058](https://github.com/CooooolFrog/dragonfly/issues/4058))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.22.2 to 2.23.4 ([#4056](https://github.com/CooooolFrog/dragonfly/issues/4056))

### Feat
- **manager:** change refresh interval of the job rate limiter ([#4063](https://github.com/CooooolFrog/dragonfly/issues/4063))
- **manager:** support preheating by percentage ([#4053](https://github.com/CooooolFrog/dragonfly/issues/4053))
- **middlewares:** add CreateJobRateLimiter middleware ([#4061](https://github.com/CooooolFrog/dragonfly/issues/4061))

### Refactor
- **manager:** redis proxy to use bufio.Reader for protocol detection ([#4062](https://github.com/CooooolFrog/dragonfly/issues/4062))


<a name="v2.2.3-rc.1"></a>
## [v2.2.3-rc.1] - 2025-05-08
### Chore
- **deps:** bump sigstore/cosign-installer from 3.8.1 to 3.8.2 ([#4034](https://github.com/CooooolFrog/dragonfly/issues/4034))
- **deps:** bump go.uber.org/mock from 0.5.1 to 0.5.2 ([#4032](https://github.com/CooooolFrog/dragonfly/issues/4032))
- **deps:** bump github.com/aws/aws-sdk-go from 1.55.6 to 1.55.7 ([#4031](https://github.com/CooooolFrog/dragonfly/issues/4031))
- **deps:** bump golang.org/x/sys from 0.32.0 to 0.33.0 ([#4033](https://github.com/CooooolFrog/dragonfly/issues/4033))
- **deps:** bump github/codeql-action from 3.28.16 to 3.28.17 ([#4029](https://github.com/CooooolFrog/dragonfly/issues/4029))
- **deps:** bump github.com/huaweicloud/huaweicloud-sdk-go-obs from 3.24.6+incompatible to 3.25.3+incompatible ([#4036](https://github.com/CooooolFrog/dragonfly/issues/4036))
- **deps:** bump actions/stale from 9.0.0 to 9.1.0 ([#4010](https://github.com/CooooolFrog/dragonfly/issues/4010))
- **deps:** bump kentaro-m/auto-assign-action from 248761c4feb3917c1b0444e33fad1a50093b9847 to 7ae38e468e64dec0af17820972bc4915aa511ec2 ([#4011](https://github.com/CooooolFrog/dragonfly/issues/4011))
- **deps:** bump actions/download-artifact from 4.2.1 to 4.3.0 ([#4015](https://github.com/CooooolFrog/dragonfly/issues/4015))
- **deps:** bump docker/build-push-action from 6.15.0 to 6.16.0 ([#4017](https://github.com/CooooolFrog/dragonfly/issues/4017))
- **deps:** bump github.com/docker/docker from 28.0.0+incompatible to 28.1.1+incompatible ([#4018](https://github.com/CooooolFrog/dragonfly/issues/4018))
- **deps:** bump github/codeql-action from 3.28.15 to 3.28.16 ([#4016](https://github.com/CooooolFrog/dragonfly/issues/4016))

### Comments
- **dferrors:** add comments for IsEndOfStream function to check … ([#4008](https://github.com/CooooolFrog/dragonfly/issues/4008))
- **dferrors:** add comments for IsGRPCDfError function to check … ([#4009](https://github.com/CooooolFrog/dragonfly/issues/4009))

### Feat
- bump client to v0.2.26 and add CRC32FromStrings func ([#4024](https://github.com/CooooolFrog/dragonfly/issues/4024))
- add unit tests for persistent cache task handlers ([#4023](https://github.com/CooooolFrog/dragonfly/issues/4023))
- **manager:** add redis proxy for manager and support redis username ([#4041](https://github.com/CooooolFrog/dragonfly/issues/4041))

### Fix
- set the pat to context for audit ([#4007](https://github.com/CooooolFrog/dragonfly/issues/4007))

### Refactor
- pass context to run gc task ([#4042](https://github.com/CooooolFrog/dragonfly/issues/4042))


<a name="v2.2.3-rc.0"></a>
## [v2.2.3-rc.0] - 2025-04-28
### Chore
- bump console from 0.2.1 to 0.2.3 ([#4001](https://github.com/CooooolFrog/dragonfly/issues/4001))
- add auto-assign workflows for pr assign ([#3998](https://github.com/CooooolFrog/dragonfly/issues/3998))

### Feat
- add the audit gc ([#3996](https://github.com/CooooolFrog/dragonfly/issues/3996))
- support to preheat multiple files ([#3997](https://github.com/CooooolFrog/dragonfly/issues/3997))
- add content for calculating task id when preheat file and handle task ([#3994](https://github.com/CooooolFrog/dragonfly/issues/3994))
- add the default gc config ([#3987](https://github.com/CooooolFrog/dragonfly/issues/3987))
- resolve statsview hanging issue by improving port handling and shutdown logic ([#3986](https://github.com/CooooolFrog/dragonfly/issues/3986))

### Fix
- remove the gorm tag for CreatedAt and UpdatedAt ([#3990](https://github.com/CooooolFrog/dragonfly/issues/3990))

### Refactor
- remove legacy job gc config ([#3995](https://github.com/CooooolFrog/dragonfly/issues/3995))
- migrate the job gc to mananger gc server ([#3991](https://github.com/CooooolFrog/dragonfly/issues/3991))
- **peer_manager:** remove unnecessary log in RunGC ([#3984](https://github.com/CooooolFrog/dragonfly/issues/3984))
- **task_manager:** Remove unnecessary error log in RunGC ([#3985](https://github.com/CooooolFrog/dragonfly/issues/3985))


<a name="v2.2.2"></a>
## [v2.2.2] - 2025-04-22
### Build
- **client-rs:** update subproject commit and add comment in dfget ([#3961](https://github.com/CooooolFrog/dragonfly/issues/3961))

### Chore
- bump golang version to 1.23.8 ([#3960](https://github.com/CooooolFrog/dragonfly/issues/3960))
- bump client-rs to v0.2.23 ([#3955](https://github.com/CooooolFrog/dragonfly/issues/3955))
- bump client-rs to v0.2.24 ([#3983](https://github.com/CooooolFrog/dragonfly/issues/3983))
- **README.md:** change OpenSSF badge ([#3959](https://github.com/CooooolFrog/dragonfly/issues/3959))
- **client-rs:** bump client-rs from 0.2.22 to 0.2.23 ([#3953](https://github.com/CooooolFrog/dragonfly/issues/3953))
- **deps:** bump github.com/gin-contrib/zap from 1.1.3 to 1.1.5 ([#3978](https://github.com/CooooolFrog/dragonfly/issues/3978))
- **deps:** bump d7y.io/api/v2 from 2.1.30 to 2.1.36 ([#3976](https://github.com/CooooolFrog/dragonfly/issues/3976))
- **deps:** bump codecov/codecov-action from 5.4.0 to 5.4.2 ([#3973](https://github.com/CooooolFrog/dragonfly/issues/3973))
- **deps:** bump google.golang.org/api from 0.219.0 to 0.229.0 ([#3977](https://github.com/CooooolFrog/dragonfly/issues/3977))
- **deps:** bump github.com/jarcoal/httpmock from 1.3.1 to 1.4.0 ([#3974](https://github.com/CooooolFrog/dragonfly/issues/3974))
- **deps:** bump github.com/gin-contrib/gzip from 1.2.2 to 1.2.3 ([#3951](https://github.com/CooooolFrog/dragonfly/issues/3951))
- **deps:** bump github.com/go-sql-driver/mysql from 1.7.0 to 1.9.2 ([#3950](https://github.com/CooooolFrog/dragonfly/issues/3950))
- **deps:** bump golang.org/x/oauth2 from 0.27.0 to 0.29.0 ([#3949](https://github.com/CooooolFrog/dragonfly/issues/3949))
- **deps:** bump google.golang.org/protobuf from 1.36.5 to 1.36.6 ([#3948](https://github.com/CooooolFrog/dragonfly/issues/3948))
- **deps:** bump github.com/gin-contrib/static from 1.1.2 to 1.1.5 ([#3947](https://github.com/CooooolFrog/dragonfly/issues/3947))
- **deps:** bump github/codeql-action from 3.28.14 to 3.28.15 ([#3946](https://github.com/CooooolFrog/dragonfly/issues/3946))

### Comments
- add initDfpath function to initialize dfpath ([#3968](https://github.com/CooooolFrog/dragonfly/issues/3968))
- **genDocCommand:** add runDoc method to generate markdown docs ([#3967](https://github.com/CooooolFrog/dragonfly/issues/3967))
- **rpcserver:** add Stop method to gracefully stop servers ([#3965](https://github.com/CooooolFrog/dragonfly/issues/3965))

### Feat
- support query parameters for list audit API ([#3970](https://github.com/CooooolFrog/dragonfly/issues/3970))
- adds the audit middleware to records the system audit events ([#3958](https://github.com/CooooolFrog/dragonfly/issues/3958))
- change PersistentCacheTask to value type and ensure empty slice response ([#3945](https://github.com/CooooolFrog/dragonfly/issues/3945))
- **manager:** add default value for operated_at ([#3966](https://github.com/CooooolFrog/dragonfly/issues/3966))

### Fix
- **genDocCommand:** handle directory creation error in runDoc
- **manager:** fix peer ID logging in loadAllPeersByTaskID ([#3954](https://github.com/CooooolFrog/dragonfly/issues/3954))

### Refactor
- replace custom math functions with built-in max/min ([#3972](https://github.com/CooooolFrog/dragonfly/issues/3972))
- modernize the code by new features provided by golang stand… ([#3963](https://github.com/CooooolFrog/dragonfly/issues/3963))
- **lint:** update golangci-lint to v2 and configure .golangci.yml ([#3956](https://github.com/CooooolFrog/dragonfly/issues/3956))
- **root.go:** simplify error handling in checkAndSpawnDaemon ([#3962](https://github.com/CooooolFrog/dragonfly/issues/3962))

### Style
- Fix spacing in ComputePieceSize comment ([#3969](https://github.com/CooooolFrog/dragonfly/issues/3969))

### Test
- **dragonfly-client-storage:** add e2e test for cache upload/download. ([#3941](https://github.com/CooooolFrog/dragonfly/issues/3941))
- **e2e:** optimize testing for preheating in cache ([#3952](https://github.com/CooooolFrog/dragonfly/issues/3952))


<a name="v2.2.2-rc.0"></a>
## [v2.2.2-rc.0] - 2025-04-11
### Chore
- update client-rs version to 0.2.21 ([#3918](https://github.com/CooooolFrog/dragonfly/issues/3918))
- update client version to 0.2.18 ([#3887](https://github.com/CooooolFrog/dragonfly/issues/3887))
- bump console to v0.2.2 and bump client-rs to v0.2.22 ([#3944](https://github.com/CooooolFrog/dragonfly/issues/3944))
- **client-rs:** update client-rs version to 0.2.20 ([#3917](https://github.com/CooooolFrog/dragonfly/issues/3917))
- **deps:** bump github.com/appleboy/gin-jwt/v2 from 2.9.2 to 2.10.3 ([#3910](https://github.com/CooooolFrog/dragonfly/issues/3910))
- **deps:** bump github.com/redis/go-redis/v9 from 9.7.0 to 9.7.3 ([#3904](https://github.com/CooooolFrog/dragonfly/issues/3904))
- **deps:** bump github.com/fsouza/fake-gcs-server from 1.52.0 to 1.52.2 ([#3931](https://github.com/CooooolFrog/dragonfly/issues/3931))
- **deps:** bump github/codeql-action from 3.28.11 to 3.28.13 ([#3914](https://github.com/CooooolFrog/dragonfly/issues/3914))
- **deps:** bump actions/upload-artifact from 4.6.1 to 4.6.2 ([#3913](https://github.com/CooooolFrog/dragonfly/issues/3913))
- **deps:** bump actions/setup-go from 5.3.0 to 5.4.0 ([#3912](https://github.com/CooooolFrog/dragonfly/issues/3912))
- **deps:** bump github/codeql-action from 3.28.13 to 3.28.14 ([#3930](https://github.com/CooooolFrog/dragonfly/issues/3930))
- **deps:** bump goreleaser/goreleaser-action from 6.2.1 to 6.3.0 ([#3924](https://github.com/CooooolFrog/dragonfly/issues/3924))
- **deps:** bump github.com/golang-jwt/jwt/v4 from 4.5.1 to 4.5.2 ([#3906](https://github.com/CooooolFrog/dragonfly/issues/3906))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.59.0 to 0.60.0 ([#3893](https://github.com/CooooolFrog/dragonfly/issues/3893))
- **deps:** bump golang.org/x/sys from 0.30.0 to 0.31.0 ([#3896](https://github.com/CooooolFrog/dragonfly/issues/3896))
- **deps:** bump docker/login-action from 3.3.0 to 3.4.0 ([#3898](https://github.com/CooooolFrog/dragonfly/issues/3898))
- **deps:** bump aquasecurity/trivy-action from 0.29.0 to 0.30.0 ([#3900](https://github.com/CooooolFrog/dragonfly/issues/3900))
- **deps:** bump github.com/containerd/containerd from 1.6.26 to 1.6.38 ([#3901](https://github.com/CooooolFrog/dragonfly/issues/3901))
- **deps:** bump golang.org/x/crypto from 0.36.0 to 0.37.0 ([#3935](https://github.com/CooooolFrog/dragonfly/issues/3935))
- **deps:** bump golang.org/x/net from 0.35.0 to 0.36.0 ([#3886](https://github.com/CooooolFrog/dragonfly/issues/3886))

### Enhance
- add comments for key scanning and context error handling

### Feat
- add log for searcher plugin ([#3942](https://github.com/CooooolFrog/dragonfly/issues/3942))
- enhance PAT validation in middleware ([#3888](https://github.com/CooooolFrog/dragonfly/issues/3888))
- add swagger annotations for open API router ([#3885](https://github.com/CooooolFrog/dragonfly/issues/3885))
- **client:** removing hostname resolution fallback when the IP address check fails ([#3936](https://github.com/CooooolFrog/dragonfly/issues/3936))
- **manager:** optimize service of the persistent cache task ([#3929](https://github.com/CooooolFrog/dragonfly/issues/3929))
- **manager:** an empty permissions list grants all permissions ([#3889](https://github.com/CooooolFrog/dragonfly/issues/3889))

### Fix
- replace builds with ids in config ([#3943](https://github.com/CooooolFrog/dragonfly/issues/3943))
- update peer and task manager methods to use correct redis key format

### Refactor
- use slices.Contains to check the permission ([#3940](https://github.com/CooooolFrog/dragonfly/issues/3940))


<a name="v2.2.1"></a>
## [v2.2.1] - 2025-03-12
### Chore
- update image versions to v2.2.1-rc.3 in compatibility e2e test ([#3882](https://github.com/CooooolFrog/dragonfly/issues/3882))
- **deps:** bump github/codeql-action from 3.28.10 to 3.28.11 ([#3881](https://github.com/CooooolFrog/dragonfly/issues/3881))
- **deps:** bump docker/setup-qemu-action from 3.4.0 to 3.6.0 ([#3880](https://github.com/CooooolFrog/dragonfly/issues/3880))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.20.0 to 1.22.0 ([#3879](https://github.com/CooooolFrog/dragonfly/issues/3879))
- **deps:** bump codecov/codecov-action from 5.3.1 to 5.4.0 ([#3878](https://github.com/CooooolFrog/dragonfly/issues/3878))
- **deps:** bump golang.org/x/time from 0.10.0 to 0.11.0 ([#3875](https://github.com/CooooolFrog/dragonfly/issues/3875))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.59.0 to 0.60.0 ([#3874](https://github.com/CooooolFrog/dragonfly/issues/3874))


<a name="v2.2.1-rc.3"></a>
## [v2.2.1-rc.3] - 2025-03-10
### Chore
- **client-rs:** update client version to 0.2.17 ([#3872](https://github.com/CooooolFrog/dragonfly/issues/3872))
- **deps:** bump github/codeql-action from 3.28.9 to 3.28.10 ([#3867](https://github.com/CooooolFrog/dragonfly/issues/3867))
- **deps:** bump golang.org/x/oauth2 from 0.25.0 to 0.27.0 ([#3862](https://github.com/CooooolFrog/dragonfly/issues/3862))
- **deps:** bump golang.org/x/crypto from 0.34.0 to 0.35.0 ([#3864](https://github.com/CooooolFrog/dragonfly/issues/3864))
- **deps:** bump actions/cache from 4.2.1 to 4.2.2 ([#3865](https://github.com/CooooolFrog/dragonfly/issues/3865))
- **deps:** bump docker/setup-buildx-action from 3.9.0 to 3.10.0 ([#3866](https://github.com/CooooolFrog/dragonfly/issues/3866))
- **deps:** bump docker/build-push-action from 6.14.0 to 6.15.0 ([#3868](https://github.com/CooooolFrog/dragonfly/issues/3868))
- **deps:** bump slsa-framework/slsa-github-generator from 2.0.0 to 2.1.0 ([#3869](https://github.com/CooooolFrog/dragonfly/issues/3869))
- **deps:** bump golang.org/x/crypto from 0.33.0 to 0.34.0 ([#3858](https://github.com/CooooolFrog/dragonfly/issues/3858))
- **deps:** bump github.com/docker/docker from 27.3.1+incompatible to 28.0.0+incompatible ([#3857](https://github.com/CooooolFrog/dragonfly/issues/3857))
- **deps:** bump k8s.io/component-base from 0.32.1 to 0.32.2 ([#3856](https://github.com/CooooolFrog/dragonfly/issues/3856))
- **deps:** bump github.com/gin-contrib/gzip from 1.0.1 to 1.2.2 ([#3855](https://github.com/CooooolFrog/dragonfly/issues/3855))
- **deps:** bump github.com/go-playground/validator/v10 from 10.24.0 to 10.25.0 ([#3854](https://github.com/CooooolFrog/dragonfly/issues/3854))
- **deps:** bump ossf/scorecard-action from 2.4.0 to 2.4.1 ([#3853](https://github.com/CooooolFrog/dragonfly/issues/3853))
- **deps:** bump sigstore/cosign-installer from 3.8.0 to 3.8.1 ([#3852](https://github.com/CooooolFrog/dragonfly/issues/3852))
- **deps:** bump docker/build-push-action from 6.13.0 to 6.14.0 ([#3851](https://github.com/CooooolFrog/dragonfly/issues/3851))
- **deps:** bump actions/upload-artifact from 4.6.0 to 4.6.1 ([#3850](https://github.com/CooooolFrog/dragonfly/issues/3850))
- **deps:** bump actions/cache from 4.2.0 to 4.2.1 ([#3849](https://github.com/CooooolFrog/dragonfly/issues/3849))

### Feat
- update helm charts and remove storage config in scheduler ([#3873](https://github.com/CooooolFrog/dragonfly/issues/3873))
- support to set piece length for preheat ([#3848](https://github.com/CooooolFrog/dragonfly/issues/3848))


<a name="v2.2.1-rc.1"></a>
## [v2.2.1-rc.1] - 2025-02-20
### Chore
- add trivy action for vulnerability scanner ([#3837](https://github.com/CooooolFrog/dragonfly/issues/3837))
- **client-rs:** update client rs submodule to v0.2.11 ([#3831](https://github.com/CooooolFrog/dragonfly/issues/3831))
- **deps:** bump github.com/spf13/cobra from 1.8.1 to 1.9.1 ([#3844](https://github.com/CooooolFrog/dragonfly/issues/3844))
- **deps:** bump d7y.io/api/v2 from 2.1.25 to 2.1.27 ([#3843](https://github.com/CooooolFrog/dragonfly/issues/3843))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.17.1 to 3.18.0 ([#3842](https://github.com/CooooolFrog/dragonfly/issues/3842))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.33.0 to 1.34.0 ([#3841](https://github.com/CooooolFrog/dragonfly/issues/3841))
- **deps:** bump goreleaser/goreleaser-action from 6.1.0 to 6.2.1 ([#3840](https://github.com/CooooolFrog/dragonfly/issues/3840))

### Feat
- replace Pipeliner with lua ([#3846](https://github.com/CooooolFrog/dragonfly/issues/3846))
- **scheduler:** handle DownloadPieceFinishedRequest and DownloadPieceFailedRequest asynchronously ([#3847](https://github.com/CooooolFrog/dragonfly/issues/3847))


<a name="v2.2.1-rc.0"></a>
## [v2.2.1-rc.0] - 2025-02-12
### Chore
- remove archives.format in goreleaser ([#3828](https://github.com/CooooolFrog/dragonfly/issues/3828))
- update client-rs and console version ([#3827](https://github.com/CooooolFrog/dragonfly/issues/3827))
- update client version to v0.2.2 ([#3754](https://github.com/CooooolFrog/dragonfly/issues/3754))
- **deps:** bump github.com/go-playground/validator/v10 from 10.23.0 to 10.24.0 ([#3795](https://github.com/CooooolFrog/dragonfly/issues/3795))
- **deps:** bump google.golang.org/api from 0.205.0 to 0.216.0 ([#3767](https://github.com/CooooolFrog/dragonfly/issues/3767))
- **deps:** bump docker/setup-qemu-action from 3.3.0 to 3.4.0 ([#3824](https://github.com/CooooolFrog/dragonfly/issues/3824))
- **deps:** bump sigstore/cosign-installer from 3.6.0 to 3.8.0 ([#3825](https://github.com/CooooolFrog/dragonfly/issues/3825))
- **deps:** bump golangci/golangci-lint-action from 6.2.0 to 6.3.2 ([#3826](https://github.com/CooooolFrog/dragonfly/issues/3826))
- **deps:** bump google.golang.org/protobuf from 1.36.4 to 1.36.5 ([#3820](https://github.com/CooooolFrog/dragonfly/issues/3820))
- **deps:** bump actions/setup-go from 5.2.0 to 5.3.0 ([#3792](https://github.com/CooooolFrog/dragonfly/issues/3792))
- **deps:** bump github.com/aws/aws-sdk-go from 1.55.5 to 1.55.6 ([#3803](https://github.com/CooooolFrog/dragonfly/issues/3803))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.58.0 to 0.59.0 ([#3801](https://github.com/CooooolFrog/dragonfly/issues/3801))
- **deps:** bump github.com/hashicorp/memberlist from 0.5.1 to 0.5.3 ([#3804](https://github.com/CooooolFrog/dragonfly/issues/3804))
- **deps:** bump google.golang.org/api from 0.216.0 to 0.219.0 ([#3805](https://github.com/CooooolFrog/dragonfly/issues/3805))
- **deps:** bump github/codeql-action from 3.28.5 to 3.28.8 ([#3806](https://github.com/CooooolFrog/dragonfly/issues/3806))
- **deps:** bump google.golang.org/grpc from 1.69.4 to 1.70.0 ([#3797](https://github.com/CooooolFrog/dragonfly/issues/3797))
- **deps:** bump docker/setup-buildx-action from 3.6.1 to 3.9.0 ([#3822](https://github.com/CooooolFrog/dragonfly/issues/3822))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.22.1 to 2.22.2 ([#3748](https://github.com/CooooolFrog/dragonfly/issues/3748))
- **deps:** bump github/codeql-action from 3.28.8 to 3.28.9 ([#3823](https://github.com/CooooolFrog/dragonfly/issues/3823))
- **deps:** bump golang.org/x/crypto from 0.32.0 to 0.33.0 ([#3817](https://github.com/CooooolFrog/dragonfly/issues/3817))
- **deps:** bump docker/build-push-action from 6.12.0 to 6.13.0 ([#3791](https://github.com/CooooolFrog/dragonfly/issues/3791))
- **deps:** bump golangci/golangci-lint-action from 6.1.1 to 6.2.0 ([#3778](https://github.com/CooooolFrog/dragonfly/issues/3778))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.58.0 to 0.59.0 ([#3775](https://github.com/CooooolFrog/dragonfly/issues/3775))
- **deps:** bump docker/build-push-action from 6.11.0 to 6.12.0 ([#3779](https://github.com/CooooolFrog/dragonfly/issues/3779))
- **deps:** bump cloud.google.com/go/storage from 1.49.0 to 1.50.0 ([#3777](https://github.com/CooooolFrog/dragonfly/issues/3777))
- **deps:** bump golang.org/x/sys from 0.28.0 to 0.29.0 ([#3776](https://github.com/CooooolFrog/dragonfly/issues/3776))
- **deps:** bump golang.org/x/crypto from 0.31.0 to 0.32.0 ([#3773](https://github.com/CooooolFrog/dragonfly/issues/3773))
- **deps:** bump k8s.io/component-base from 0.32.0 to 0.32.1 ([#3774](https://github.com/CooooolFrog/dragonfly/issues/3774))
- **deps:** bump github/codeql-action from 3.28.0 to 3.28.1 ([#3759](https://github.com/CooooolFrog/dragonfly/issues/3759))
- **deps:** bump docker/setup-qemu-action from 3.2.0 to 3.3.0 ([#3762](https://github.com/CooooolFrog/dragonfly/issues/3762))
- **deps:** bump docker/build-push-action from 6.10.0 to 6.11.0 ([#3760](https://github.com/CooooolFrog/dragonfly/issues/3760))
- **deps:** bump actions/upload-artifact from 4.5.0 to 4.6.0 ([#3761](https://github.com/CooooolFrog/dragonfly/issues/3761))
- **deps:** bump codecov/codecov-action from 5.1.2 to 5.3.1 ([#3793](https://github.com/CooooolFrog/dragonfly/issues/3793))
- **deps:** bump golang.org/x/time from 0.9.0 to 0.10.0 ([#3821](https://github.com/CooooolFrog/dragonfly/issues/3821))
- **deps:** bump github/codeql-action from 3.28.1 to 3.28.5 ([#3794](https://github.com/CooooolFrog/dragonfly/issues/3794))

### Docs
- remove FOSSA in README.md ([#3829](https://github.com/CooooolFrog/dragonfly/issues/3829))
- add OpenSSF Scorecard badge ([#3788](https://github.com/CooooolFrog/dragonfly/issues/3788))
- add FOSSA embed status badge ([#3789](https://github.com/CooooolFrog/dragonfly/issues/3789))

### Feat
- bump api from v2.1.23 to v2.1.25 and add ttl for StatPersistentCacheTask ([#3830](https://github.com/CooooolFrog/dragonfly/issues/3830))
- add option to cancel idle peer task ([#3810](https://github.com/CooooolFrog/dragonfly/issues/3810))
- replicate the persistent cache task when delete host ([#3787](https://github.com/CooooolFrog/dragonfly/issues/3787))
- delete persistent cache task from peer ([#3786](https://github.com/CooooolFrog/dragonfly/issues/3786))
- replicate persistent cache task when task needs persistent replicas ([#3784](https://github.com/CooooolFrog/dragonfly/issues/3784))
- handle the error for scheduling persistent cache task ([#3781](https://github.com/CooooolFrog/dragonfly/issues/3781))
- add pieceLength for creating persistent cache task ([#3772](https://github.com/CooooolFrog/dragonfly/issues/3772))
- implement FindCandidatePersistentCacheParents for scheduling persistent cache task ([#3770](https://github.com/CooooolFrog/dragonfly/issues/3770))
- support google cloud storage ([#3756](https://github.com/CooooolFrog/dragonfly/issues/3756))

### Fix
- update client version in e2e testing ([#3755](https://github.com/CooooolFrog/dragonfly/issues/3755))

### Refactor
- refactor the e2e file server ([#3741](https://github.com/CooooolFrog/dragonfly/issues/3741))

### Test
- add the e2e test cases for download & prefetch rate limit ([#3757](https://github.com/CooooolFrog/dragonfly/issues/3757))


<a name="v2.2.0"></a>
## [v2.2.0] - 2024-12-31
### Chore
- update console version to 0.1.43 ([#3732](https://github.com/CooooolFrog/dragonfly/issues/3732))
- update client version to 0.2.0 ([#3735](https://github.com/CooooolFrog/dragonfly/issues/3735))
- update console version to v0.1.42 ([#3729](https://github.com/CooooolFrog/dragonfly/issues/3729))
- **deps:** bump github.com/mennanov/limiters from 1.9.0 to 1.11.0 ([#3723](https://github.com/CooooolFrog/dragonfly/issues/3723))
- **deps:** bump k8s.io/component-base from 0.31.2 to 0.32.0 ([#3737](https://github.com/CooooolFrog/dragonfly/issues/3737))
- **deps:** bump github.com/onsi/gomega from 1.36.1 to 1.36.2 ([#3740](https://github.com/CooooolFrog/dragonfly/issues/3740))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.57.0 to 0.58.0 ([#3736](https://github.com/CooooolFrog/dragonfly/issues/3736))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.19.1 to 1.20.0 ([#3724](https://github.com/CooooolFrog/dragonfly/issues/3724))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.22.0 to 2.22.1 ([#3726](https://github.com/CooooolFrog/dragonfly/issues/3726))
- **deps:** bump google.golang.org/protobuf from 1.35.2 to 1.36.1 ([#3725](https://github.com/CooooolFrog/dragonfly/issues/3725))
- **deps:** bump github.com/go-playground/validator/v10 from 10.22.1 to 10.23.0 ([#3722](https://github.com/CooooolFrog/dragonfly/issues/3722))
- **deps:** bump codecov/codecov-action from 5.1.1 to 5.1.2 ([#3721](https://github.com/CooooolFrog/dragonfly/issues/3721))
- **deps:** bump actions/upload-artifact from 4.4.3 to 4.5.0 ([#3720](https://github.com/CooooolFrog/dragonfly/issues/3720))
- **deps:** bump github/codeql-action from 3.27.9 to 3.28.0 ([#3719](https://github.com/CooooolFrog/dragonfly/issues/3719))
- **deps:** bump helm/kind-action from 1.10.0 to 1.12.0 ([#3718](https://github.com/CooooolFrog/dragonfly/issues/3718))

### Feat
- add cacher for gorm to reduce the load of db ([#3734](https://github.com/CooooolFrog/dragonfly/issues/3734))
- reduce the cache size of the lfu in manager ([#3730](https://github.com/CooooolFrog/dragonfly/issues/3730))

### Fix
- skip the sync peer that does not belong to the scheduler cluster ([#3731](https://github.com/CooooolFrog/dragonfly/issues/3731))


<a name="v2.1.67"></a>
## [v2.1.67] - 2024-12-23
### Chore
- update client-rs version to v0.1.126 ([#3717](https://github.com/CooooolFrog/dragonfly/issues/3717))
- rename repo from Dragonfly2 to dragonfly ([#3715](https://github.com/CooooolFrog/dragonfly/issues/3715))

### Docs
- change the order of the badges ([#3716](https://github.com/CooooolFrog/dragonfly/issues/3716))


<a name="v2.1.66"></a>
## [v2.1.66] - 2024-12-17
### Chore
- update console verison to v0.1.41 ([#3700](https://github.com/CooooolFrog/dragonfly/issues/3700))
- update client version to v0.1.125 ([#3687](https://github.com/CooooolFrog/dragonfly/issues/3687))
- **deps:** bump golang.org/x/crypto from 0.28.0 to 0.30.0 ([#3691](https://github.com/CooooolFrog/dragonfly/issues/3691))
- **deps:** bump codecov/codecov-action from 5.0.7 to 5.1.1 ([#3696](https://github.com/CooooolFrog/dragonfly/issues/3696))
- **deps:** bump github.com/onsi/gomega from 1.35.1 to 1.36.1 ([#3704](https://github.com/CooooolFrog/dragonfly/issues/3704))
- **deps:** bump go.opentelemetry.io/otel/trace from 1.32.0 to 1.33.0 ([#3705](https://github.com/CooooolFrog/dragonfly/issues/3705))
- **deps:** bump github/codeql-action from 3.27.6 to 3.27.9 ([#3708](https://github.com/CooooolFrog/dragonfly/issues/3708))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.18.0 to 1.19.1 ([#3701](https://github.com/CooooolFrog/dragonfly/issues/3701))
- **deps:** bump golang.org/x/crypto from 0.30.0 to 0.31.0 ([#3710](https://github.com/CooooolFrog/dragonfly/issues/3710))
- **deps:** bump github.com/go-redsync/redsync/v4 from 4.8.1 to 4.13.0 ([#3702](https://github.com/CooooolFrog/dragonfly/issues/3702))
- **deps:** bump github/codeql-action from 3.27.5 to 3.27.6 ([#3695](https://github.com/CooooolFrog/dragonfly/issues/3695))
- **deps:** bump actions/cache from 4.1.2 to 4.2.0 ([#3694](https://github.com/CooooolFrog/dragonfly/issues/3694))
- **deps:** bump docker/build-push-action from 6.7.0 to 6.10.0 ([#3693](https://github.com/CooooolFrog/dragonfly/issues/3693))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.16.0 to 1.18.0 ([#3692](https://github.com/CooooolFrog/dragonfly/issues/3692))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.29.0 to 1.32.0 ([#3689](https://github.com/CooooolFrog/dragonfly/issues/3689))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.20.1 to 2.22.0 ([#3688](https://github.com/CooooolFrog/dragonfly/issues/3688))
- **deps:** bump actions/setup-go from 5.1.0 to 5.2.0 ([#3707](https://github.com/CooooolFrog/dragonfly/issues/3707))

### Feat
- when the redis is disabled, AnnounceHost need to skip store redis ([#3712](https://github.com/CooooolFrog/dragonfly/issues/3712))
- add client version for MakeSchedulersKeyForPeerInManager ([#3711](https://github.com/CooooolFrog/dragonfly/issues/3711))
- add AllSeedPeersScope for preheating ([#3698](https://github.com/CooooolFrog/dragonfly/issues/3698))
- load empty return false in persistentcache ([#3697](https://github.com/CooooolFrog/dragonfly/issues/3697))


<a name="v2.1.65"></a>
## [v2.1.65] - 2024-12-06
### Chore
- update rust client version ([#3685](https://github.com/CooooolFrog/dragonfly/issues/3685))
- **deps:** bump codecov/codecov-action from 5.0.2 to 5.0.7 ([#3671](https://github.com/CooooolFrog/dragonfly/issues/3671))
- **deps:** bump github/codeql-action from 3.27.4 to 3.27.5 ([#3670](https://github.com/CooooolFrog/dragonfly/issues/3670))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.13.0 to 1.16.0 ([#3667](https://github.com/CooooolFrog/dragonfly/issues/3667))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.17.0 to 3.17.1 ([#3665](https://github.com/CooooolFrog/dragonfly/issues/3665))

### Docs
- update SECURITY.md for Dragonfly ([#3680](https://github.com/CooooolFrog/dragonfly/issues/3680))

### Enhance
- support syncpeers by service and optimize the merge logic ([#3637](https://github.com/CooooolFrog/dragonfly/issues/3637))

### Feat
- reuse connections and limit the number of connections for preheating ([#3683](https://github.com/CooooolFrog/dragonfly/issues/3683))
- delete jobs in batches ([#3682](https://github.com/CooooolFrog/dragonfly/issues/3682))
- optimize implement of the sync peers ([#3677](https://github.com/CooooolFrog/dragonfly/issues/3677))
- remove deploy without docker compose ([#3672](https://github.com/CooooolFrog/dragonfly/issues/3672))
- support CRC-32-Castagnoli algorithm ([#3664](https://github.com/CooooolFrog/dragonfly/issues/3664))


<a name="v2.1.64"></a>
## [v2.1.64] - 2024-11-22
### Chore
- update client submodule ([#3661](https://github.com/CooooolFrog/dragonfly/issues/3661))
- **deps:** bump github.com/gammazero/deque from 0.2.1 to 1.0.0 ([#3657](https://github.com/CooooolFrog/dragonfly/issues/3657))
- **deps:** bump github/codeql-action from 3.27.1 to 3.27.4 ([#3658](https://github.com/CooooolFrog/dragonfly/issues/3658))
- **deps:** bump codecov/codecov-action from 4.6.0 to 5.0.2 ([#3659](https://github.com/CooooolFrog/dragonfly/issues/3659))
- **deps:** bump github.com/swaggo/swag from 1.16.3 to 1.16.4 ([#3655](https://github.com/CooooolFrog/dragonfly/issues/3655))
- **deps:** bump golang.org/x/sync from 0.8.0 to 0.9.0 ([#3654](https://github.com/CooooolFrog/dragonfly/issues/3654))
- **deps:** bump google.golang.org/protobuf from 1.35.1 to 1.35.2 ([#3653](https://github.com/CooooolFrog/dragonfly/issues/3653))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.56.0 to 0.57.0 ([#3644](https://github.com/CooooolFrog/dragonfly/issues/3644))
- **deps:** bump k8s.io/component-base from 0.29.2 to 0.31.2 ([#3649](https://github.com/CooooolFrog/dragonfly/issues/3649))
- **deps:** bump goreleaser/goreleaser-action from 6.0.0 to 6.1.0 ([#3648](https://github.com/CooooolFrog/dragonfly/issues/3648))
- **deps:** bump golang.org/x/oauth2 from 0.23.0 to 0.24.0 ([#3647](https://github.com/CooooolFrog/dragonfly/issues/3647))
- **deps:** bump github/codeql-action from 3.27.0 to 3.27.1 ([#3646](https://github.com/CooooolFrog/dragonfly/issues/3646))
- **deps:** bump google.golang.org/api from 0.199.0 to 0.205.0 ([#3645](https://github.com/CooooolFrog/dragonfly/issues/3645))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.56.0 to 0.57.0 ([#3643](https://github.com/CooooolFrog/dragonfly/issues/3643))

### Feat
- optimize api for shceduling ([#3660](https://github.com/CooooolFrog/dragonfly/issues/3660))
- add disk bandwidth information for host ([#3652](https://github.com/CooooolFrog/dragonfly/issues/3652))
- add garbage collection for persistent cache host ([#3642](https://github.com/CooooolFrog/dragonfly/issues/3642))
- store persistent cache host by announce host api ([#3640](https://github.com/CooooolFrog/dragonfly/issues/3640))


<a name="v2.1.63"></a>
## [v2.1.63] - 2024-11-11
### Chore
- update console submodule ([#3641](https://github.com/CooooolFrog/dragonfly/issues/3641))
- update console and rust client version ([#3639](https://github.com/CooooolFrog/dragonfly/issues/3639))
- update golang version to v1.23.0 ([#3609](https://github.com/CooooolFrog/dragonfly/issues/3609))
- **deps:** bump actions/checkout from 4.2.1 to 4.2.2 ([#3618](https://github.com/CooooolFrog/dragonfly/issues/3618))
- **deps:** bump go.uber.org/mock from 0.4.0 to 0.5.0 ([#3630](https://github.com/CooooolFrog/dragonfly/issues/3630))
- **deps:** bump golang.org/x/time from 0.6.0 to 0.7.0 ([#3627](https://github.com/CooooolFrog/dragonfly/issues/3627))
- **deps:** bump github.com/onsi/gomega from 1.34.1 to 1.35.1 ([#3631](https://github.com/CooooolFrog/dragonfly/issues/3631))
- **deps:** bump github/codeql-action from 3.26.12 to 3.27.0 ([#3617](https://github.com/CooooolFrog/dragonfly/issues/3617))
- **deps:** bump actions/cache from 4.1.0 to 4.1.2 ([#3616](https://github.com/CooooolFrog/dragonfly/issues/3616))
- **deps:** bump actions/setup-go from 5.0.2 to 5.1.0 ([#3615](https://github.com/CooooolFrog/dragonfly/issues/3615))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.53.0 to 0.56.0 ([#3612](https://github.com/CooooolFrog/dragonfly/issues/3612))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.14.6 to 3.17.0 ([#3610](https://github.com/CooooolFrog/dragonfly/issues/3610))
- **deps:** bump github.com/golang-jwt/jwt/v4 from 4.5.0 to 4.5.1 ([#3632](https://github.com/CooooolFrog/dragonfly/issues/3632))

### Feat
- add rate limit for job open api by cluster ([#3638](https://github.com/CooooolFrog/dragonfly/issues/3638))
- implement delete peer and task in persistent cache ([#3623](https://github.com/CooooolFrog/dragonfly/issues/3623))
- implement upload persistent cache task ([#3620](https://github.com/CooooolFrog/dragonfly/issues/3620))
- optimize error message of preheating ([#3622](https://github.com/CooooolFrog/dragonfly/issues/3622))
- add filtered query params of the containerd ([#3621](https://github.com/CooooolFrog/dragonfly/issues/3621))
- implement delete persistent cache task in scheduler ([#3619](https://github.com/CooooolFrog/dragonfly/issues/3619))

### Fix
- generate wrong sql with gorm ([#3626](https://github.com/CooooolFrog/dragonfly/issues/3626))


<a name="v2.1.62"></a>
## [v2.1.62] - 2024-10-28
### Chore
- update client version to v0.1.113 ([#3593](https://github.com/CooooolFrog/dragonfly/issues/3593))
- **deps:** bump actions/checkout from 4.2.0 to 4.2.1 ([#3600](https://github.com/CooooolFrog/dragonfly/issues/3600))
- **deps:** bump actions/upload-artifact from 4.3.6 to 4.4.3 ([#3599](https://github.com/CooooolFrog/dragonfly/issues/3599))
- **deps:** bump github.com/prometheus/client_golang from 1.20.4 to 1.20.5 ([#3598](https://github.com/CooooolFrog/dragonfly/issues/3598))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.54.0 to 0.56.0 ([#3594](https://github.com/CooooolFrog/dragonfly/issues/3594))

### Feat
- add DefaultFilteredQueryParams for job ([#3608](https://github.com/CooooolFrog/dragonfly/issues/3608))
- support searching task by url for GetTask and DeleteTask ([#3607](https://github.com/CooooolFrog/dragonfly/issues/3607))
- implement StatPersistentCachePeerRequest and StatPersistentCacheTaskRequest for persistent cache ([#3603](https://github.com/CooooolFrog/dragonfly/issues/3603))
- add interface of the persistent cache resource ([#3602](https://github.com/CooooolFrog/dragonfly/issues/3602))
- add peer manager for persistent cache task ([#3592](https://github.com/CooooolFrog/dragonfly/issues/3592))
- change job gc interval in manager ([#3591](https://github.com/CooooolFrog/dragonfly/issues/3591))
- add E2E tests for cases that peers going offline ([#3524](https://github.com/CooooolFrog/dragonfly/issues/3524))
- add CreateJobCount and CreateJobSuccessCount metrics ([#3588](https://github.com/CooooolFrog/dragonfly/issues/3588))

### Test
- refactor the tests for the list hosts and delete host ([#3605](https://github.com/CooooolFrog/dragonfly/issues/3605))
- refactor e2e testings of the client leaving ([#3590](https://github.com/CooooolFrog/dragonfly/issues/3590))


<a name="v2.1.61"></a>
## [v2.1.61] - 2024-10-17
### Chore
- generate sbom for release artifacts ([#3587](https://github.com/CooooolFrog/dragonfly/issues/3587))
- generate SBOM for release artifacts ([#3585](https://github.com/CooooolFrog/dragonfly/issues/3585))
- update go version to v1.22.4 ([#3580](https://github.com/CooooolFrog/dragonfly/issues/3580))
- update rust client and console submodule ([#3567](https://github.com/CooooolFrog/dragonfly/issues/3567))
- **deps:** bump github.com/docker/docker from 27.1.1+incompatible to 27.3.1+incompatible ([#3556](https://github.com/CooooolFrog/dragonfly/issues/3556))
- **deps:** bump google.golang.org/api from 0.197.0 to 0.199.0 ([#3554](https://github.com/CooooolFrog/dragonfly/issues/3554))
- **deps:** bump codecov/codecov-action from 4.5.0 to 4.6.0 ([#3559](https://github.com/CooooolFrog/dragonfly/issues/3559))
- **deps:** bump actions/cache from 4.0.2 to 4.1.0 ([#3561](https://github.com/CooooolFrog/dragonfly/issues/3561))
- **deps:** bump github/codeql-action from 3.26.8 to 3.26.12 ([#3565](https://github.com/CooooolFrog/dragonfly/issues/3565))
- **deps:** bump actions/checkout from 4.1.7 to 4.2.0 ([#3549](https://github.com/CooooolFrog/dragonfly/issues/3549))

### Feat
- add self-signed certs for mTLS ([#3583](https://github.com/CooooolFrog/dragonfly/issues/3583))
- support set self-signed cert for service ([#3568](https://github.com/CooooolFrog/dragonfly/issues/3568))
- add fsm for persistent cache peer ([#3563](https://github.com/CooooolFrog/dragonfly/issues/3563))


<a name="v2.1.60"></a>
## [v2.1.60] - 2024-10-09
### Chore
- update client-rs version ([#3562](https://github.com/CooooolFrog/dragonfly/issues/3562))

### Feat
- add downloadRate and uploadRate for host ([#3548](https://github.com/CooooolFrog/dragonfly/issues/3548))
- removed network topology ([#3547](https://github.com/CooooolFrog/dragonfly/issues/3547))
- add host manager for persistent cache ([#3546](https://github.com/CooooolFrog/dragonfly/issues/3546))
- add persistent cache task for scheduler ([#3545](https://github.com/CooooolFrog/dragonfly/issues/3545))
- increase interval of the preheat polling ([#3544](https://github.com/CooooolFrog/dragonfly/issues/3544))
- add auto switch scheduler e2e test ([#3486](https://github.com/CooooolFrog/dragonfly/issues/3486))
- rename `scheduler/resource` to `scheduler/resource/standard` ([#3542](https://github.com/CooooolFrog/dragonfly/issues/3542))
- update new task type(TaskType_STANDARD, TaskType_PERSISTENT, TaskType_PERSISTENT_CACHE) ([#3540](https://github.com/CooooolFrog/dragonfly/issues/3540))

### Test
- remove scheduler e2e test ([#3543](https://github.com/CooooolFrog/dragonfly/issues/3543))


<a name="v2.1.59"></a>
## [v2.1.59] - 2024-09-26
### Chore
- update api version to v2.0.158 and update helm chart ([#3527](https://github.com/CooooolFrog/dragonfly/issues/3527))
- **deps:** bump github.com/prometheus/client_golang from 1.19.0 to 1.20.4 ([#3532](https://github.com/CooooolFrog/dragonfly/issues/3532))
- **deps:** bump github/codeql-action from 3.26.2 to 3.26.8 ([#3530](https://github.com/CooooolFrog/dragonfly/issues/3530))
- **deps:** bump actions/checkout from 4.1.1 to 4.1.7 ([#3529](https://github.com/CooooolFrog/dragonfly/issues/3529))
- **deps:** bump sigstore/cosign-installer from 3.5.0 to 3.6.0 ([#3528](https://github.com/CooooolFrog/dragonfly/issues/3528))

### Feat
- support preheat with self-signed certs ([#3541](https://github.com/CooooolFrog/dragonfly/issues/3541))
- add metrics for grpc api of the cache task ([#3539](https://github.com/CooooolFrog/dragonfly/issues/3539))
- support set max threads ([#3537](https://github.com/CooooolFrog/dragonfly/issues/3537))
- fixed lint in manager sync_peers.go ([#3536](https://github.com/CooooolFrog/dragonfly/issues/3536))
- seed max concurrent ([#3482](https://github.com/CooooolFrog/dragonfly/issues/3482))

### Fix
- make e2e test  ([#3487](https://github.com/CooooolFrog/dragonfly/issues/3487))
- update get and delete task unit test and e2e test. ([#3525](https://github.com/CooooolFrog/dragonfly/issues/3525))
- **dfget:** Change file path ([#3519](https://github.com/CooooolFrog/dragonfly/issues/3519))

### Test
- generate mode task id for testings ([#3526](https://github.com/CooooolFrog/dragonfly/issues/3526))


<a name="v2.1.58"></a>
## [v2.1.58] - 2024-09-20
### Chore
- fixed Pinned-Dependencies in actions ([#3521](https://github.com/CooooolFrog/dragonfly/issues/3521))
- add sbom and provenance for docker build

### Docs
- fixed the badge uri of the OpenSSF Scorecard ([#3520](https://github.com/CooooolFrog/dragonfly/issues/3520))

### Fix
- update delete task rpc and create e2e test. ([#3447](https://github.com/CooooolFrog/dragonfly/issues/3447))

### Refactor
- get task job and delete task job ([#3522](https://github.com/CooooolFrog/dragonfly/issues/3522))


<a name="v2.1.57"></a>
## [v2.1.57] - 2024-09-18
### Chore
- remove verify base image
- fixed Token-Permissions for Actions ([#3517](https://github.com/CooooolFrog/dragonfly/issues/3517))
- remove distribution in goreleaser
- fixed Pinned-Dependencies for actions ([#3518](https://github.com/CooooolFrog/dragonfly/issues/3518))
- remove sboms in goreleaser
- add -y option to cosign
- add upload-tag-name for release
- generate SLSA3 provenance for GoReleaser ([#3516](https://github.com/CooooolFrog/dragonfly/issues/3516))
- add COSIGN_PUBLIC_KEY for cosign verify
- add GITHUB_TOKEN for ghcr to login ([#3515](https://github.com/CooooolFrog/dragonfly/issues/3515))
- signed release for container image ([#3514](https://github.com/CooooolFrog/dragonfly/issues/3514))
- update client version ([#3496](https://github.com/CooooolFrog/dragonfly/issues/3496))
- update console submodule version ([#3505](https://github.com/CooooolFrog/dragonfly/issues/3505))
- **deps:** bump google.golang.org/api from 0.189.0 to 0.197.0 ([#3510](https://github.com/CooooolFrog/dragonfly/issues/3510))
- **deps:** bump github.com/jellydator/ttlcache/v3 from 3.2.0 to 3.3.0 ([#3494](https://github.com/CooooolFrog/dragonfly/issues/3494))
- **deps:** bump golang.org/x/oauth2 from 0.21.0 to 0.23.0 ([#3491](https://github.com/CooooolFrog/dragonfly/issues/3491))

### Feat
- change gc config for manager job ([#3507](https://github.com/CooooolFrog/dragonfly/issues/3507))
- when peer disabled shared, scheduler will skip peer in filterCandidateParents ([#3506](https://github.com/CooooolFrog/dragonfly/issues/3506))
- add scope for preheating job ([#3497](https://github.com/CooooolFrog/dragonfly/issues/3497))
- clean up expired jobs to prevent performance problems ([#3504](https://github.com/CooooolFrog/dragonfly/issues/3504))
- remove order by created_at for selecting job ([#3502](https://github.com/CooooolFrog/dragonfly/issues/3502))
- change timeout and max refresh in jwt token ([#3501](https://github.com/CooooolFrog/dragonfly/issues/3501))
- response server real ip in X-Server-IP header ([#3500](https://github.com/CooooolFrog/dragonfly/issues/3500))
- add job feature for scheduler's announcer ([#3489](https://github.com/CooooolFrog/dragonfly/issues/3489))
- add scheduler features api for manager ([#3488](https://github.com/CooooolFrog/dragonfly/issues/3488))

### Fix
- fixed the token-permission and pinned dependencies ([#3513](https://github.com/CooooolFrog/dragonfly/issues/3513))
- add sentinel authentication settings ([#3484](https://github.com/CooooolFrog/dragonfly/issues/3484))


<a name="v2.1.56"></a>
## [v2.1.56] - 2024-09-06
### Chore
- mute some reclaim check log ([#3469](https://github.com/CooooolFrog/dragonfly/issues/3469))
- udpate client-rs version ([#3461](https://github.com/CooooolFrog/dragonfly/issues/3461))
- fix 'as' keyword should match the case of the 'from' keyword ([#3445](https://github.com/CooooolFrog/dragonfly/issues/3445))
- update console and client version ([#3438](https://github.com/CooooolFrog/dragonfly/issues/3438))
- **deps:** bump github.com/zeebo/blake3 from 0.2.3 to 0.2.4 ([#3466](https://github.com/CooooolFrog/dragonfly/issues/3466))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.19.0 to 2.20.1 ([#3463](https://github.com/CooooolFrog/dragonfly/issues/3463))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.28.0 to 1.29.0 ([#3462](https://github.com/CooooolFrog/dragonfly/issues/3462))
- **deps:** bump golang.org/x/crypto from 0.25.0 to 0.26.0 ([#3464](https://github.com/CooooolFrog/dragonfly/issues/3464))
- **deps:** bump ossf/scorecard-action from 2.3.1 to 2.4.0 ([#3453](https://github.com/CooooolFrog/dragonfly/issues/3453))
- **deps:** bump actions/upload-artifact from 3.pre.node20 to 4.3.6 ([#3454](https://github.com/CooooolFrog/dragonfly/issues/3454))

### Feat
- use context.Background for ScheduleCandidateParents to prevent stream breaking ([#3485](https://github.com/CooooolFrog/dragonfly/issues/3485))
- add rate limit middlewares for job ([#3481](https://github.com/CooooolFrog/dragonfly/issues/3481))
- add ratelimit for job in manager ([#3480](https://github.com/CooooolFrog/dragonfly/issues/3480))
- add proxy error metric ([#3470](https://github.com/CooooolFrog/dragonfly/issues/3470))
- add source http metric ([#3468](https://github.com/CooooolFrog/dragonfly/issues/3468))
- resource add rdb for cache task ([#3467](https://github.com/CooooolFrog/dragonfly/issues/3467))
- remove trainer and model
- daemon add back source metric ([#3460](https://github.com/CooooolFrog/dragonfly/issues/3460))
- add ScheduleDuration for recording duration of the scheduling ([#3444](https://github.com/CooooolFrog/dragonfly/issues/3444))
- added the scorecard github action and its badge ([#3440](https://github.com/CooooolFrog/dragonfly/issues/3440))

### Fix
- key/value format with whitespace separator should not be used ([#3459](https://github.com/CooooolFrog/dragonfly/issues/3459))
- fix null pointer dereference issue when parent task storage registration fails  ([#3458](https://github.com/CooooolFrog/dragonfly/issues/3458))
- AutoIssueCert loses control in dfdaemon([#3422](https://github.com/CooooolFrog/dragonfly/issues/3422)) ([#3423](https://github.com/CooooolFrog/dragonfly/issues/3423))


<a name="v2.1.55"></a>
## [v2.1.55] - 2024-08-16
### Chore
- rename name_template to version_template in goreleaser
- export content range in http source ([#3437](https://github.com/CooooolFrog/dragonfly/issues/3437))
- fix typo ([#3435](https://github.com/CooooolFrog/dragonfly/issues/3435))
- update rust client version ([#3424](https://github.com/CooooolFrog/dragonfly/issues/3424))
- **deps:** bump golang.org/x/sys from 0.22.0 to 0.24.0 ([#3428](https://github.com/CooooolFrog/dragonfly/issues/3428))

### Feat
- optimize GetTaskJob and DeleteTaskJob ([#3434](https://github.com/CooooolFrog/dragonfly/issues/3434))
- add accept ranges header ([#3433](https://github.com/CooooolFrog/dragonfly/issues/3433))
- add delete task and list tasks manager api with request type and service type. ([#3378](https://github.com/CooooolFrog/dragonfly/issues/3378))

### Fix
- reuse length check ([#3432](https://github.com/CooooolFrog/dragonfly/issues/3432))
- start with zero half-open interval range ([#3431](https://github.com/CooooolFrog/dragonfly/issues/3431))


<a name="v2.1.54"></a>
## [v2.1.54] - 2024-08-07
### Chore
- update client verison to v0.1.96 ([#3420](https://github.com/CooooolFrog/dragonfly/issues/3420))


<a name="v2.1.53"></a>
## [v2.1.53] - 2024-08-06
### Chore
- optimize hijack ca format ([#3418](https://github.com/CooooolFrog/dragonfly/issues/3418))
- update api verison to v2.0.141 ([#3384](https://github.com/CooooolFrog/dragonfly/issues/3384))
- optimize tls cert expire check ([#3394](https://github.com/CooooolFrog/dragonfly/issues/3394))
- build service local in Dockerfile ([#3411](https://github.com/CooooolFrog/dragonfly/issues/3411))
- add self-hosted runner for ci ([#3409](https://github.com/CooooolFrog/dragonfly/issues/3409))
- update rust client version ([#3403](https://github.com/CooooolFrog/dragonfly/issues/3403))
- update client version ([#3395](https://github.com/CooooolFrog/dragonfly/issues/3395))
- **deps:** bump github.com/docker/docker from 27.0.3+incompatible to 27.1.1+incompatible ([#3397](https://github.com/CooooolFrog/dragonfly/issues/3397))
- **deps:** bump github.com/redis/go-redis/v9 from 9.5.1 to 9.6.1 ([#3399](https://github.com/CooooolFrog/dragonfly/issues/3399))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.52.0 to 0.53.0 ([#3398](https://github.com/CooooolFrog/dragonfly/issues/3398))
- **deps:** bump google.golang.org/api from 0.183.0 to 0.189.0 ([#3400](https://github.com/CooooolFrog/dragonfly/issues/3400))
- **deps:** bump golang.org/x/sync from 0.7.0 to 0.8.0 ([#3417](https://github.com/CooooolFrog/dragonfly/issues/3417))
- **deps:** bump github.com/gin-contrib/gzip from 0.0.6 to 1.0.1 ([#3390](https://github.com/CooooolFrog/dragonfly/issues/3390))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.52.0 to 0.53.0 ([#3389](https://github.com/CooooolFrog/dragonfly/issues/3389))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.14.4 to 3.14.6 ([#3413](https://github.com/CooooolFrog/dragonfly/issues/3413))
- **deps:** bump github.com/aws/aws-sdk-go from 1.54.15 to 1.54.19 ([#3380](https://github.com/CooooolFrog/dragonfly/issues/3380))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.25.0 to 1.28.0 ([#3379](https://github.com/CooooolFrog/dragonfly/issues/3379))

### Feat
- remove unused context and update api version ([#3412](https://github.com/CooooolFrog/dragonfly/issues/3412))
- use content.Background() to avoid stream cancel by dfdaemon ([#3402](https://github.com/CooooolFrog/dragonfly/issues/3402))

### Fix
- immediately flush data to client for event-stream response ([#3375](https://github.com/CooooolFrog/dragonfly/issues/3375))

### Refractor
- fix typo in object storage.go  ([#3374](https://github.com/CooooolFrog/dragonfly/issues/3374))


<a name="v2.1.52"></a>
## [v2.1.52] - 2024-07-15
### Chore
- update client and console version ([#3376](https://github.com/CooooolFrog/dragonfly/issues/3376))
- update chart version for rust client ([#3373](https://github.com/CooooolFrog/dragonfly/issues/3373))
- **deps:** bump google.golang.org/grpc from 1.64.0 to 1.64.1 ([#3371](https://github.com/CooooolFrog/dragonfly/issues/3371))


<a name="v2.1.51"></a>
## [v2.1.51] - 2024-07-09
### Chore
- avoid use default mux ([#3346](https://github.com/CooooolFrog/dragonfly/issues/3346))
- **deps:** bump golang.org/x/crypto from 0.24.0 to 0.25.0 ([#3365](https://github.com/CooooolFrog/dragonfly/issues/3365))
- **deps:** bump github.com/huaweicloud/huaweicloud-sdk-go-obs from 3.23.12+incompatible to 3.24.6+incompatible ([#3369](https://github.com/CooooolFrog/dragonfly/issues/3369))
- **deps:** bump github.com/aws/aws-sdk-go from 1.53.10 to 1.54.15 ([#3368](https://github.com/CooooolFrog/dragonfly/issues/3368))
- **deps:** bump go.opentelemetry.io/otel from 1.27.0 to 1.28.0 ([#3366](https://github.com/CooooolFrog/dragonfly/issues/3366))
- **deps:** bump github.com/hashicorp/memberlist from 0.5.0 to 0.5.1 ([#3356](https://github.com/CooooolFrog/dragonfly/issues/3356))
- **deps:** bump github.com/docker/docker from 26.1.1+incompatible to 27.0.3+incompatible ([#3358](https://github.com/CooooolFrog/dragonfly/issues/3358))

### Feat
- send config of the seed peer cluster for load limit ([#3370](https://github.com/CooooolFrog/dragonfly/issues/3370))
- delete the host after host leave peers during host manager gc ([#3361](https://github.com/CooooolFrog/dragonfly/issues/3361))
- scheduler detects peer survivability and clears offline peers' metadata ([#3353](https://github.com/CooooolFrog/dragonfly/issues/3353))

### Fix
- file task uid gid ([#3359](https://github.com/CooooolFrog/dragonfly/issues/3359))

### Refactor
- update comments and define size constant for human readable ([#3363](https://github.com/CooooolFrog/dragonfly/issues/3363))
- Preallocate pieceDigests slice in genMetadata for improved ([#3325](https://github.com/CooooolFrog/dragonfly/issues/3325))


<a name="v2.1.50"></a>
## [v2.1.50] - 2024-06-27
### Chore
- optimize calculate digest ([#3343](https://github.com/CooooolFrog/dragonfly/issues/3343))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.17.2 to 2.19.0 ([#3340](https://github.com/CooooolFrog/dragonfly/issues/3340))

### Docs
- add the responsibilities of contributor roles ([#3345](https://github.com/CooooolFrog/dragonfly/issues/3345))

### Feat
- reduce the scheduler return content of the parent ([#3349](https://github.com/CooooolFrog/dragonfly/issues/3349))
- optimize peertask watch dog ([#3327](https://github.com/CooooolFrog/dragonfly/issues/3327))


<a name="v2.1.49"></a>
## [v2.1.49] - 2024-06-20
### Chore
- update client rs verison to v0.1.82 ([#3336](https://github.com/CooooolFrog/dragonfly/issues/3336))
- **deps:** bump github.com/looplab/fsm from 1.0.1 to 1.0.2 ([#3329](https://github.com/CooooolFrog/dragonfly/issues/3329))
- **deps:** bump github.com/spf13/cobra from 1.7.0 to 1.8.1 ([#3330](https://github.com/CooooolFrog/dragonfly/issues/3330))
- **deps:** bump gorm.io/gorm from 1.25.7-0.20240204074919-46816ad31dde to 1.25.10 ([#3331](https://github.com/CooooolFrog/dragonfly/issues/3331))
- **deps:** bump docker/build-push-action from 5 to 6 ([#3333](https://github.com/CooooolFrog/dragonfly/issues/3333))

### Feat
- increase revice msg size and send msg size in grpc ([#3335](https://github.com/CooooolFrog/dragonfly/issues/3335))
- add cache peer to grpc api ([#3334](https://github.com/CooooolFrog/dragonfly/issues/3334))
- update api protoc version to v2.0.118 ([#3326](https://github.com/CooooolFrog/dragonfly/issues/3326))


<a name="v2.1.48"></a>
## [v2.1.48] - 2024-06-14

<a name="v2.1.47"></a>
## [v2.1.47] - 2024-06-14
### Chore
- update client-rs version ([#3321](https://github.com/CooooolFrog/dragonfly/issues/3321))
- update client and console verison ([#3319](https://github.com/CooooolFrog/dragonfly/issues/3319))
- update client version ([#3307](https://github.com/CooooolFrog/dragonfly/issues/3307))
- update client verison ([#3291](https://github.com/CooooolFrog/dragonfly/issues/3291))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.24.4 to 3.24.5 ([#3302](https://github.com/CooooolFrog/dragonfly/issues/3302))
- **deps:** bump golang.org/x/crypto from 0.23.0 to 0.24.0 ([#3314](https://github.com/CooooolFrog/dragonfly/issues/3314))
- **deps:** bump goreleaser/goreleaser-action from 5 to 6 ([#3313](https://github.com/CooooolFrog/dragonfly/issues/3313))
- **deps:** bump google.golang.org/api from 0.181.0 to 0.183.0 ([#3309](https://github.com/CooooolFrog/dragonfly/issues/3309))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.47.0 to 0.52.0 ([#3315](https://github.com/CooooolFrog/dragonfly/issues/3315))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.51.0 to 0.52.0 ([#3301](https://github.com/CooooolFrog/dragonfly/issues/3301))
- **deps:** bump github.com/onsi/gomega from 1.32.0 to 1.33.1 ([#3304](https://github.com/CooooolFrog/dragonfly/issues/3304))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.14.2 to 3.14.4 ([#3318](https://github.com/CooooolFrog/dragonfly/issues/3318))
- **deps:** bump go.opentelemetry.io/otel from 1.26.0 to 1.27.0 ([#3294](https://github.com/CooooolFrog/dragonfly/issues/3294))
- **deps:** bump github.com/gin-contrib/zap from 1.1.1 to 1.1.3 ([#3292](https://github.com/CooooolFrog/dragonfly/issues/3292))
- **deps:** bump github.com/aws/aws-sdk-go from 1.52.2 to 1.53.10 ([#3295](https://github.com/CooooolFrog/dragonfly/issues/3295))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.12 to 3.24.4 ([#3293](https://github.com/CooooolFrog/dragonfly/issues/3293))

### Feat
- change loadLimit validation to 50000 ([#3322](https://github.com/CooooolFrog/dragonfly/issues/3322))
- if host is seed client, set concurrentUploadLimit when dynconfig refreshes ([#3320](https://github.com/CooooolFrog/dragonfly/issues/3320))

### Fix
- release action unknown flag ([#3324](https://github.com/CooooolFrog/dragonfly/issues/3324))
- correct spelling typos ([#3312](https://github.com/CooooolFrog/dragonfly/issues/3312))
- clean running piece done ([#3308](https://github.com/CooooolFrog/dragonfly/issues/3308))


<a name="v2.1.46"></a>
## [v2.1.46] - 2024-05-24
### Chore
- update client and console ([#3290](https://github.com/CooooolFrog/dragonfly/issues/3290))
- update client verison ([#3289](https://github.com/CooooolFrog/dragonfly/issues/3289))
- optimize register timeout ([#3288](https://github.com/CooooolFrog/dragonfly/issues/3288))
- **deps:** bump google.golang.org/api from 0.180.0 to 0.181.0 ([#3287](https://github.com/CooooolFrog/dragonfly/issues/3287))
- **deps:** bump github.com/gin-gonic/gin from 1.9.1 to 1.10.0 ([#3284](https://github.com/CooooolFrog/dragonfly/issues/3284))
- **deps:** bump helm/kind-action from 1.9.0 to 1.10.0 ([#3282](https://github.com/CooooolFrog/dragonfly/issues/3282))


<a name="v2.1.45"></a>
## [v2.1.45] - 2024-05-17
### BUGFIX
- hardlink ([#3277](https://github.com/CooooolFrog/dragonfly/issues/3277))

### Chore
- update client version ([#3267](https://github.com/CooooolFrog/dragonfly/issues/3267))
- update pex replica clean logic ([#3272](https://github.com/CooooolFrog/dragonfly/issues/3272))
- fix typo in scheduler config ([#3273](https://github.com/CooooolFrog/dragonfly/issues/3273))
- transport proxy ttl cache ([#3268](https://github.com/CooooolFrog/dragonfly/issues/3268))
- update pex initial backoff ([#3266](https://github.com/CooooolFrog/dragonfly/issues/3266))
- update api verison to v2.0.112 ([#3256](https://github.com/CooooolFrog/dragonfly/issues/3256))
- **deps:** bump github.com/jackc/pgx/v5 from 5.4.3 to 5.5.4 ([#3125](https://github.com/CooooolFrog/dragonfly/issues/3125))
- **deps:** bump google.golang.org/api from 0.176.0 to 0.180.0 ([#3261](https://github.com/CooooolFrog/dragonfly/issues/3261))
- **deps:** bump github.com/gin-contrib/static from 1.1.1 to 1.1.2 ([#3262](https://github.com/CooooolFrog/dragonfly/issues/3262))
- **deps:** bump golangci/golangci-lint-action from 5 to 6 ([#3260](https://github.com/CooooolFrog/dragonfly/issues/3260))

### Fix
- use safe set in client daemon network_topology.go ([#3276](https://github.com/CooooolFrog/dragonfly/issues/3276))
- peer exchange option ([#3259](https://github.com/CooooolFrog/dragonfly/issues/3259))
- pass through the headers specified by user to preheat job ([#3257](https://github.com/CooooolFrog/dragonfly/issues/3257))
- DownloadPeerDuration observe cost by milliseconds ([#3255](https://github.com/CooooolFrog/dragonfly/issues/3255))

### Test
- add nydus e2e testing for API v2 ([#3234](https://github.com/CooooolFrog/dragonfly/issues/3234))


<a name="v2.1.44"></a>
## [v2.1.44] - 2024-05-07
### Chore
- release failed running piece ([#3214](https://github.com/CooooolFrog/dragonfly/issues/3214))
- update api version ([#3251](https://github.com/CooooolFrog/dragonfly/issues/3251))
- update client version ([#3238](https://github.com/CooooolFrog/dragonfly/issues/3238))
- update client submodule version ([#3229](https://github.com/CooooolFrog/dragonfly/issues/3229))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.49.0 to 0.51.0 ([#3245](https://github.com/CooooolFrog/dragonfly/issues/3245))
- **deps:** bump github.com/aws/aws-sdk-go from 1.51.11 to 1.52.2 ([#3250](https://github.com/CooooolFrog/dragonfly/issues/3250))
- **deps:** bump github.com/docker/docker from 26.0.2+incompatible to 26.1.1+incompatible ([#3249](https://github.com/CooooolFrog/dragonfly/issues/3249))
- **deps:** bump go.opentelemetry.io/otel/trace from 1.25.0 to 1.26.0 ([#3242](https://github.com/CooooolFrog/dragonfly/issues/3242))
- **deps:** bump golangci/golangci-lint-action from 4 to 5 ([#3240](https://github.com/CooooolFrog/dragonfly/issues/3240))
- **deps:** bump helm/kind-action from 1.9.0 to 1.10.0 ([#3239](https://github.com/CooooolFrog/dragonfly/issues/3239))

### Docs
- add iQIYI to ADOPTERS.md ([#3248](https://github.com/CooooolFrog/dragonfly/issues/3248))

### Feat
- set calculateDigest to false in dfdaemon, because it has bug in large files ([#3235](https://github.com/CooooolFrog/dragonfly/issues/3235))

### Fix
- client rpcserver subscriber hang ([#3246](https://github.com/CooooolFrog/dragonfly/issues/3246))

### Test
- change number of the concurrent in e2e testing ([#3236](https://github.com/CooooolFrog/dragonfly/issues/3236))
- delayed checking Seed Peer download content ([#3230](https://github.com/CooooolFrog/dragonfly/issues/3230))
- add unit test for manager handler ([#3222](https://github.com/CooooolFrog/dragonfly/issues/3222))


<a name="v2.1.43"></a>
## [v2.1.43] - 2024-04-24
### Chore
- update console submodule ([#3228](https://github.com/CooooolFrog/dragonfly/issues/3228))
- update client and helm charts submodule version ([#3227](https://github.com/CooooolFrog/dragonfly/issues/3227))
- change gitignore for vendor ([#3225](https://github.com/CooooolFrog/dragonfly/issues/3225))
- optimize e2e test config ([#3223](https://github.com/CooooolFrog/dragonfly/issues/3223))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.24.0 to 1.25.0 ([#3221](https://github.com/CooooolFrog/dragonfly/issues/3221))
- **deps:** bump google.golang.org/api from 0.172.0 to 0.176.0 ([#3220](https://github.com/CooooolFrog/dragonfly/issues/3220))
- **deps:** bump golang.org/x/sys from 0.18.0 to 0.19.0 ([#3217](https://github.com/CooooolFrog/dragonfly/issues/3217))

### Test
- add compatibility e2e testing for v2 ([#3216](https://github.com/CooooolFrog/dragonfly/issues/3216))


<a name="v2.1.42"></a>
## [v2.1.42] - 2024-04-22
### Chore
- **deps:** bump github.com/docker/docker from 26.0.0+incompatible to 26.0.2+incompatible ([#3204](https://github.com/CooooolFrog/dragonfly/issues/3204))

### Test
- add preheat e2e testing for v2 ([#3215](https://github.com/CooooolFrog/dragonfly/issues/3215))
- add concurrent download e2e testings ([#3212](https://github.com/CooooolFrog/dragonfly/issues/3212))
- add range download testings by prefetch proxy ([#3211](https://github.com/CooooolFrog/dragonfly/issues/3211))
- add unit tests for manager handlers ([#3202](https://github.com/CooooolFrog/dragonfly/issues/3202))
- add range request e2e testings for dfdaemon proxy ([#3208](https://github.com/CooooolFrog/dragonfly/issues/3208))
- add piece lenght, tag, application test cases for dfget downloading ([#3207](https://github.com/CooooolFrog/dragonfly/issues/3207))
- add download with dfget e2e test for API v2(Rust Client) ([#3205](https://github.com/CooooolFrog/dragonfly/issues/3205))
- add containerd e2e testings to API v2(rust client) ([#3203](https://github.com/CooooolFrog/dragonfly/issues/3203))


<a name="v2.1.41"></a>
## [v2.1.41] - 2024-04-17
### Chore
- update charts version ([#3190](https://github.com/CooooolFrog/dragonfly/issues/3190))
- update single piece storage ([#3186](https://github.com/CooooolFrog/dragonfly/issues/3186))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.49.0 to 0.50.0 ([#3194](https://github.com/CooooolFrog/dragonfly/issues/3194))
- **deps:** bump google.golang.org/api from 0.171.0 to 0.172.0 ([#3192](https://github.com/CooooolFrog/dragonfly/issues/3192))
- **deps:** upgrade 'github.com/RichardKnop/machinery' mod version from v1… ([#3197](https://github.com/CooooolFrog/dragonfly/issues/3197))

### Feat
- support preheating by v2 grpc protocol ([#3201](https://github.com/CooooolFrog/dragonfly/issues/3201))
- trigger download task by stream ([#3200](https://github.com/CooooolFrog/dragonfly/issues/3200))
- update charts and optimize charts values by image tag ([#3187](https://github.com/CooooolFrog/dragonfly/issues/3187))
- revise probedcount threshold ([#3189](https://github.com/CooooolFrog/dragonfly/issues/3189))

### Fix
- typo in manager database ([#3191](https://github.com/CooooolFrog/dragonfly/issues/3191))
- single long task keepalive ([#3184](https://github.com/CooooolFrog/dragonfly/issues/3184))


<a name="v2.1.40"></a>
## [v2.1.40] - 2024-04-10
### Chore
- update api version to v2.0.109 and fix rust e2e CI ([#3185](https://github.com/CooooolFrog/dragonfly/issues/3185))

### Feat
- peer exchange ([#3141](https://github.com/CooooolFrog/dragonfly/issues/3141))


<a name="v2.1.39"></a>
## [v2.1.39] - 2024-04-09
### Chore
- update go-redis verison to v9 ([#3182](https://github.com/CooooolFrog/dragonfly/issues/3182))
- **deps:** bump github.com/aws/aws-sdk-go from 1.50.25 to 1.51.11 ([#3163](https://github.com/CooooolFrog/dragonfly/issues/3163))
- **deps:** bump github.com/docker/docker from 25.0.3+incompatible to 26.0.0+incompatible ([#3159](https://github.com/CooooolFrog/dragonfly/issues/3159))
- **deps:** bump github.com/gin-contrib/zap from 1.1.0 to 1.1.1 ([#3161](https://github.com/CooooolFrog/dragonfly/issues/3161))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.16.0 to 2.17.1 ([#3162](https://github.com/CooooolFrog/dragonfly/issues/3162))

### Feat
- optimize reload storage ([#3167](https://github.com/CooooolFrog/dragonfly/issues/3167))
- remove otel in client ([#3169](https://github.com/CooooolFrog/dragonfly/issues/3169))
- add networkTopology e2e test ([#3124](https://github.com/CooooolFrog/dragonfly/issues/3124))

### Fix
- fix network topology e2e test ([#3171](https://github.com/CooooolFrog/dragonfly/issues/3171))
- daemon reuse file fd leak ([#3180](https://github.com/CooooolFrog/dragonfly/issues/3180))
- tunnel https ([#3153](https://github.com/CooooolFrog/dragonfly/issues/3153))
- tiny piece digest reader ([#3164](https://github.com/CooooolFrog/dragonfly/issues/3164))
- parent peertask conductor race condition ([#3154](https://github.com/CooooolFrog/dragonfly/issues/3154))


<a name="v2.1.38"></a>
## [v2.1.38] - 2024-03-28
### Chore
- optimize log ([#3151](https://github.com/CooooolFrog/dragonfly/issues/3151))
- rename test actions ([#3133](https://github.com/CooooolFrog/dragonfly/issues/3133))
- **deps:** bump gorm.io/gorm from 1.25.7-0.20240204074919-46816ad31dde to 1.25.8 ([#3147](https://github.com/CooooolFrog/dragonfly/issues/3147))
- **deps:** bump github.com/onsi/gomega from 1.31.1 to 1.32.0 ([#3149](https://github.com/CooooolFrog/dragonfly/issues/3149))
- **deps:** bump google.golang.org/api from 0.169.0 to 0.171.0 ([#3150](https://github.com/CooooolFrog/dragonfly/issues/3150))
- **deps:** bump github.com/gin-contrib/static from 1.1.0 to 1.1.1 ([#3146](https://github.com/CooooolFrog/dragonfly/issues/3146))
- **deps:** bump github.com/aliyun/aliyun-oss-go-sdk from 3.0.1+incompatible to 3.0.2+incompatible ([#3138](https://github.com/CooooolFrog/dragonfly/issues/3138))
- **deps:** bump github.com/gin-contrib/static from 0.0.1 to 1.1.0 ([#3137](https://github.com/CooooolFrog/dragonfly/issues/3137))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.14.1 to 3.14.2 ([#3136](https://github.com/CooooolFrog/dragonfly/issues/3136))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.15.0 to 2.16.0 ([#3135](https://github.com/CooooolFrog/dragonfly/issues/3135))

### Feat
- optimize graph based on sync.Map ([#3152](https://github.com/CooooolFrog/dragonfly/issues/3152))

### Fix
- scheduler storage buffer size limit ([#3142](https://github.com/CooooolFrog/dragonfly/issues/3142))
- make e2e-test ([#3143](https://github.com/CooooolFrog/dragonfly/issues/3143))
- storage buffer size ([#3130](https://github.com/CooooolFrog/dragonfly/issues/3130))


<a name="v2.1.37"></a>
## [v2.1.37] - 2024-03-15
### Chore
- update client-rs version ([#3131](https://github.com/CooooolFrog/dragonfly/issues/3131))
- update api verison to v2.0.106 ([#3129](https://github.com/CooooolFrog/dragonfly/issues/3129))
- remove tips in actions ([#3120](https://github.com/CooooolFrog/dragonfly/issues/3120))
- move cache for buildx in docker actions ([#3117](https://github.com/CooooolFrog/dragonfly/issues/3117))
- update client-rs submodule version ([#3116](https://github.com/CooooolFrog/dragonfly/issues/3116))
- add check-size to actions ([#3115](https://github.com/CooooolFrog/dragonfly/issues/3115))

### Feat
- storage write buffer size ([#3127](https://github.com/CooooolFrog/dragonfly/issues/3127))
- add e2e testings framework for rust client and add docker layer cache ([#3123](https://github.com/CooooolFrog/dragonfly/issues/3123))
- rename e2eutil to util ([#3122](https://github.com/CooooolFrog/dragonfly/issues/3122))


<a name="v2.1.36"></a>
## [v2.1.36] - 2024-03-12
### Chore
- add client-rs submodule to main ([#3114](https://github.com/CooooolFrog/dragonfly/issues/3114))
- **deps:** bump github.com/opencontainers/image-spec from 1.1.0-rc6 to 1.1.0 ([#3112](https://github.com/CooooolFrog/dragonfly/issues/3112))
- **deps:** bump golang.org/x/crypto from 0.20.0 to 0.21.0 ([#3111](https://github.com/CooooolFrog/dragonfly/issues/3111))
- **deps:** bump gorm.io/driver/postgres from 1.5.2 to 1.5.7 ([#3110](https://github.com/CooooolFrog/dragonfly/issues/3110))
- **deps:** bump google.golang.org/api from 0.165.0 to 0.169.0 ([#3113](https://github.com/CooooolFrog/dragonfly/issues/3113))
- **deps:** bump github.com/go-playground/validator/v10 from 10.18.0 to 10.19.0 ([#3109](https://github.com/CooooolFrog/dragonfly/issues/3109))


<a name="v2.1.35"></a>
## [v2.1.35] - 2024-03-08
### Chore
- update console verison to v1.0.24 ([#3106](https://github.com/CooooolFrog/dragonfly/issues/3106))
- update api version ([#3098](https://github.com/CooooolFrog/dragonfly/issues/3098))
- **deps:** bump github.com/prometheus/client_golang from 1.18.0 to 1.19.0 ([#3100](https://github.com/CooooolFrog/dragonfly/issues/3100))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.22.0 to 1.24.0 ([#3101](https://github.com/CooooolFrog/dragonfly/issues/3101))
- **deps:** bump github.com/gin-contrib/zap from 0.2.0 to 1.1.0 ([#3102](https://github.com/CooooolFrog/dragonfly/issues/3102))

### Feat
- change peer count limit for task ([#3107](https://github.com/CooooolFrog/dragonfly/issues/3107))
- change scheduling params ([#3097](https://github.com/CooooolFrog/dragonfly/issues/3097))

### Fix
- erorr middleware set error message return headers were already written ([#3105](https://github.com/CooooolFrog/dragonfly/issues/3105))
- typo in scheduler config ([#3095](https://github.com/CooooolFrog/dragonfly/issues/3095))


<a name="v2.1.34"></a>
## [v2.1.34] - 2024-02-28
### Chore
- update dragonfly-api version to v2.0.97 ([#3093](https://github.com/CooooolFrog/dragonfly/issues/3093))
- **deps:** bump gorm.io/driver/mysql from 1.5.2 to 1.5.4 ([#3089](https://github.com/CooooolFrog/dragonfly/issues/3089))
- **deps:** bump github.com/go-playground/validator/v10 from 10.17.0 to 10.18.0 ([#3088](https://github.com/CooooolFrog/dragonfly/issues/3088))
- **deps:** bump github.com/aws/aws-sdk-go from 1.49.17 to 1.50.25 ([#3090](https://github.com/CooooolFrog/dragonfly/issues/3090))
- **deps:** bump github.com/docker/docker from 24.0.7+incompatible to 25.0.3+incompatible ([#3087](https://github.com/CooooolFrog/dragonfly/issues/3087))

### Feat
- add ip to keepalive's params ([#3094](https://github.com/CooooolFrog/dragonfly/issues/3094))
- make log rotation settings configurable ([#3085](https://github.com/CooooolFrog/dragonfly/issues/3085))
- replace bitmap with bitset ([#3023](https://github.com/CooooolFrog/dragonfly/issues/3023))


<a name="v2.1.33"></a>
## [v2.1.33] - 2024-02-21
### Chore
- update console submodule version to v1.0.23 ([#3084](https://github.com/CooooolFrog/dragonfly/issues/3084))
- **deps:** bump golangci/golangci-lint-action from 3 to 4 ([#3077](https://github.com/CooooolFrog/dragonfly/issues/3077))
- **deps:** bump helm/kind-action from 1.8.0 to 1.9.0 ([#3076](https://github.com/CooooolFrog/dragonfly/issues/3076))
- **deps:** bump k8s.io/component-base from 0.29.0 to 0.29.2 ([#3081](https://github.com/CooooolFrog/dragonfly/issues/3081))
- **deps:** bump google.golang.org/api from 0.156.0 to 0.165.0 ([#3080](https://github.com/CooooolFrog/dragonfly/issues/3080))

### Feat
- cluster scopes add hostname regexes for client matchs cluster by hostname ([#3083](https://github.com/CooooolFrog/dragonfly/issues/3083))
- dfcache support customize socket path ([#3078](https://github.com/CooooolFrog/dragonfly/issues/3078))

### Fix
- dfcache daemon sock ([#3079](https://github.com/CooooolFrog/dragonfly/issues/3079))


<a name="v2.1.32"></a>
## [v2.1.32] - 2024-02-08
### Chore
- update helm charts submodule ([#3061](https://github.com/CooooolFrog/dragonfly/issues/3061))
- **deps:** bump codecov/codecov-action from 3 to 4 ([#3068](https://github.com/CooooolFrog/dragonfly/issues/3068))
- **deps:** bump github.com/appleboy/gin-jwt/v2 from 2.9.1 to 2.9.2 ([#3066](https://github.com/CooooolFrog/dragonfly/issues/3066))
- **deps:** bump github.com/opencontainers/image-spec from 1.1.0-rc5 to 1.1.0-rc6 ([#3065](https://github.com/CooooolFrog/dragonfly/issues/3065))
- **deps:** bump github.com/swaggo/swag from 1.16.1 to 1.16.3 ([#3063](https://github.com/CooooolFrog/dragonfly/issues/3063))

### Feat
- add leave host to dfdaemon ([#3070](https://github.com/CooooolFrog/dragonfly/issues/3070))
- add scan function in redis and neighbours function in network topology ([#3048](https://github.com/CooooolFrog/dragonfly/issues/3048))

### Fix
- comments in neighbours function ([#3069](https://github.com/CooooolFrog/dragonfly/issues/3069))


<a name="v2.1.31"></a>
## [v2.1.31] - 2024-01-30
### Chore
- update api version and console version ([#3059](https://github.com/CooooolFrog/dragonfly/issues/3059))
- update api version to v2.0.83 ([#3037](https://github.com/CooooolFrog/dragonfly/issues/3037))
- optimize dfget log ([#3034](https://github.com/CooooolFrog/dragonfly/issues/3034))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.46.1 to 0.47.0 ([#3057](https://github.com/CooooolFrog/dragonfly/issues/3057))
- **deps:** bump github.com/huaweicloud/huaweicloud-sdk-go-obs from 3.23.9+incompatible to 3.23.12+incompatible ([#3056](https://github.com/CooooolFrog/dragonfly/issues/3056))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.20.0 to 1.22.0 ([#3054](https://github.com/CooooolFrog/dragonfly/issues/3054))
- **deps:** bump go.opentelemetry.io/otel/trace from 1.21.0 to 1.22.0 ([#3043](https://github.com/CooooolFrog/dragonfly/issues/3043))
- **deps:** bump github.com/onsi/gomega from 1.30.0 to 1.31.1 ([#3042](https://github.com/CooooolFrog/dragonfly/issues/3042))
- **deps:** bump github.com/go-playground/validator/v10 from 10.16.0 to 10.17.0 ([#3041](https://github.com/CooooolFrog/dragonfly/issues/3041))
- **deps:** bump actions/cache from 3 to 4 ([#3039](https://github.com/CooooolFrog/dragonfly/issues/3039))

### Feat
- optimize network topology in evaluator ([#3053](https://github.com/CooooolFrog/dragonfly/issues/3053))
- evalutorBase adds networktopology ([#3015](https://github.com/CooooolFrog/dragonfly/issues/3015))
- optimize scheduling log when exceeded RetryBackToSourceLimit ([#3051](https://github.com/CooooolFrog/dragonfly/issues/3051))
- optimize scheduler default config ([#3050](https://github.com/CooooolFrog/dragonfly/issues/3050))
- replace filters with filteredQueryParams ([#3049](https://github.com/CooooolFrog/dragonfly/issues/3049))
- add gzip middleware for static file ([#3045](https://github.com/CooooolFrog/dragonfly/issues/3045))
- change request header in response ([#3038](https://github.com/CooooolFrog/dragonfly/issues/3038))
- CanAddPeerEdge is moved to the end of filter ([#3036](https://github.com/CooooolFrog/dragonfly/issues/3036))


<a name="v2.1.30"></a>
## [v2.1.30] - 2024-01-16
### Chore
- optimize dfget progressbar ([#3021](https://github.com/CooooolFrog/dragonfly/issues/3021))
- add yyzai384 to maintainers ([#3018](https://github.com/CooooolFrog/dragonfly/issues/3018))
- **deps:** bump google.golang.org/api from 0.154.0 to 0.156.0 ([#3026](https://github.com/CooooolFrog/dragonfly/issues/3026))
- **deps:** bump github.com/docker/go-connections from 0.4.0 to 0.5.0 ([#3025](https://github.com/CooooolFrog/dragonfly/issues/3025))

### Feat
- revise findProbedHosts function ([#2986](https://github.com/CooooolFrog/dragonfly/issues/2986))
- output path should be empty, prevent seed peer copy file to output path ([#3019](https://github.com/CooooolFrog/dragonfly/issues/3019))
- re-register when scheduler go away ([#3016](https://github.com/CooooolFrog/dragonfly/issues/3016))

### Test
- remove failed case in handleDownloadPeerStartedRequest ([#3017](https://github.com/CooooolFrog/dragonfly/issues/3017))


<a name="v2.1.29"></a>
## [v2.1.29] - 2024-01-10
### Chore
- update api verison to v2.0.76 ([#3013](https://github.com/CooooolFrog/dragonfly/issues/3013))
- add cache docker layers action to docker.yaml ([#2994](https://github.com/CooooolFrog/dragonfly/issues/2994))
- **deps:** bump github.com/aws/aws-sdk-go from 1.49.9 to 1.49.17 ([#3009](https://github.com/CooooolFrog/dragonfly/issues/3009))
- **deps:** bump github.com/deckarep/golang-set/v2 from 2.3.1 to 2.6.0 ([#3008](https://github.com/CooooolFrog/dragonfly/issues/3008))
- **deps:** bump github.com/prometheus/client_golang from 1.17.0 to 1.18.0 ([#3007](https://github.com/CooooolFrog/dragonfly/issues/3007))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.10.0 to 1.13.0 ([#3005](https://github.com/CooooolFrog/dragonfly/issues/3005))

### Feat
- prevent concurrent request to cause state switch failed ([#3014](https://github.com/CooooolFrog/dragonfly/issues/3014))
- update scheduler and seed peer by ip and port ([#3012](https://github.com/CooooolFrog/dragonfly/issues/3012))
- add block parent to peer if peer reschedule in downloading ([#3011](https://github.com/CooooolFrog/dragonfly/issues/3011))
- report bad parents and add bad parents to block parent set ([#3010](https://github.com/CooooolFrog/dragonfly/issues/3010))
- change state set when trigger seed peer ([#3003](https://github.com/CooooolFrog/dragonfly/issues/3003))
- remove maxScheduleCount and optimize register in v2 ([#3001](https://github.com/CooooolFrog/dragonfly/issues/3001))


<a name="v2.1.28"></a>
## [v2.1.28] - 2024-01-02
### Chore
- optimize stream task id ([#2983](https://github.com/CooooolFrog/dragonfly/issues/2983))
- add draft in goreleaser ([#2981](https://github.com/CooooolFrog/dragonfly/issues/2981))
- **deps:** bump github.com/opencontainers/image-spec from 1.1.0-rc2.0.20221005185240-3a7f492d3f1b to 1.1.0-rc5 ([#2991](https://github.com/CooooolFrog/dragonfly/issues/2991))
- **deps:** bump github.com/casbin/casbin/v2 from 2.77.2 to 2.81.0 ([#2990](https://github.com/CooooolFrog/dragonfly/issues/2990))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.10 to 3.23.12 ([#2988](https://github.com/CooooolFrog/dragonfly/issues/2988))
- **deps:** bump github.com/aws/aws-sdk-go from 1.48.16 to 1.49.9 ([#2977](https://github.com/CooooolFrog/dragonfly/issues/2977))
- **deps:** bump google.golang.org/api from 0.151.0 to 0.154.0 ([#2974](https://github.com/CooooolFrog/dragonfly/issues/2974))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.13.1 to 2.13.2 ([#2973](https://github.com/CooooolFrog/dragonfly/issues/2973))

### Feat
- remove object storage and model message ([#2992](https://github.com/CooooolFrog/dragonfly/issues/2992))
- add difference function for slices ([#2984](https://github.com/CooooolFrog/dragonfly/issues/2984))
- back to source with piece group(multiple pieces) ([#2826](https://github.com/CooooolFrog/dragonfly/issues/2826))
- if seed peer download task, return NeedBackToSource directly ([#2980](https://github.com/CooooolFrog/dragonfly/issues/2980))
- if trigger peer is seed peer, set NeedBackToSource flag true ([#2972](https://github.com/CooooolFrog/dragonfly/issues/2972))
- remove register seed peer in announce peer ([#2971](https://github.com/CooooolFrog/dragonfly/issues/2971))
- trigger download by task id ([#2970](https://github.com/CooooolFrog/dragonfly/issues/2970))

### Fix
- peer events state failed when register task ([#2979](https://github.com/CooooolFrog/dragonfly/issues/2979))

### Test
- add slice unit test ([#2985](https://github.com/CooooolFrog/dragonfly/issues/2985))


<a name="v2.1.27"></a>
## [v2.1.27] - 2023-12-22
### Chore
- grpc reverts to v1.60.0-dev version ([#2969](https://github.com/CooooolFrog/dragonfly/issues/2969))


<a name="v2.1.26"></a>
## [v2.1.26] - 2023-12-22
### Chore
- grpc_health_probe supports different architectures in dockerfile ([#2966](https://github.com/CooooolFrog/dragonfly/issues/2966))
- optimize sync pieces error report ([#2925](https://github.com/CooooolFrog/dragonfly/issues/2925))
- **deps:** bump github.com/containerd/containerd from 1.5.18 to 1.6.26 ([#2962](https://github.com/CooooolFrog/dragonfly/issues/2962))
- **deps:** bump github.com/google/uuid from 1.4.0 to 1.5.0 ([#2951](https://github.com/CooooolFrog/dragonfly/issues/2951))
- **deps:** bump golang.org/x/crypto from 0.16.0 to 0.17.0 ([#2955](https://github.com/CooooolFrog/dragonfly/issues/2955))
- **deps:** bump github/codeql-action from 2 to 3 ([#2954](https://github.com/CooooolFrog/dragonfly/issues/2954))
- **deps:** bump actions/upload-artifact from 3 to 4 ([#2953](https://github.com/CooooolFrog/dragonfly/issues/2953))
- **deps:** bump github.com/go-playground/validator/v10 from 10.15.5 to 10.16.0 ([#2952](https://github.com/CooooolFrog/dragonfly/issues/2952))
- **deps:** bump github.com/onsi/gomega from 1.29.0 to 1.30.0 ([#2950](https://github.com/CooooolFrog/dragonfly/issues/2950))
- **deps:** bump k8s.io/component-base from 0.28.3 to 0.29.0 ([#2949](https://github.com/CooooolFrog/dragonfly/issues/2949))
- **deps:** bump golang.org/x/time from 0.4.0 to 0.5.0 ([#2948](https://github.com/CooooolFrog/dragonfly/issues/2948))

### Feat
- support logging the range and header when registering a task ([#2967](https://github.com/CooooolFrog/dragonfly/issues/2967))
- add WithStatsHandler and StatsHandler replace deprecated func ([#2959](https://github.com/CooooolFrog/dragonfly/issues/2959))
- optimize trigger download task return empty ([#2958](https://github.com/CooooolFrog/dragonfly/issues/2958))
- implement trigger download task by seed peer in v2 ([#2957](https://github.com/CooooolFrog/dragonfly/issues/2957))
- optimize parameters in seed peer DownloadTask ([#2947](https://github.com/CooooolFrog/dragonfly/issues/2947))
- add log for resolver and dynconfig ([#2945](https://github.com/CooooolFrog/dragonfly/issues/2945))

### Fix
- change "ASK" to "ACK" in adopters.md ([#2960](https://github.com/CooooolFrog/dragonfly/issues/2960))
- add context timeout to health check ([#2943](https://github.com/CooooolFrog/dragonfly/issues/2943))


<a name="v2.1.25"></a>
## [v2.1.25] - 2023-12-13
### Chore
- **deps:** bump github.com/schollz/progressbar/v3 from 3.13.1 to 3.14.1 ([#2937](https://github.com/CooooolFrog/dragonfly/issues/2937))
- **deps:** bump github.com/aws/aws-sdk-go from 1.48.11 to 1.48.16 ([#2936](https://github.com/CooooolFrog/dragonfly/issues/2936))
- **deps:** bump actions/setup-go from 4 to 5 ([#2941](https://github.com/CooooolFrog/dragonfly/issues/2941))

### Feat
- remove concurrent_piece_count in scheduler ([#2942](https://github.com/CooooolFrog/dragonfly/issues/2942))
- change DefaultSchedulerFilterParentLimit to 15 ([#2932](https://github.com/CooooolFrog/dragonfly/issues/2932))

### Fix
- downgrade gorm.io/driver/postgres to 1.5.2 ([#2908](https://github.com/CooooolFrog/dragonfly/issues/2908))


<a name="v2.1.24"></a>
## [v2.1.24] - 2023-12-07
### Chore
- update api version to v2.0.59 ([#2923](https://github.com/CooooolFrog/dragonfly/issues/2923))
- **deps:** bump github.com/aws/aws-sdk-go from 1.47.9 to 1.48.11 ([#2919](https://github.com/CooooolFrog/dragonfly/issues/2919))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.45.0 to 0.46.1 ([#2916](https://github.com/CooooolFrog/dragonfly/issues/2916))
- **deps:** bump golang.org/x/crypto from 0.15.0 to 0.16.0 ([#2917](https://github.com/CooooolFrog/dragonfly/issues/2917))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.9 to 3.23.10 ([#2900](https://github.com/CooooolFrog/dragonfly/issues/2900))
- **deps:** bump google.golang.org/api from 0.138.0 to 0.151.0 ([#2902](https://github.com/CooooolFrog/dragonfly/issues/2902))
- **deps:** bump github.com/docker/distribution from 2.8.2+incompatible to 2.8.3+incompatible ([#2898](https://github.com/CooooolFrog/dragonfly/issues/2898))
- **deps:** bump golang.org/x/sys from 0.14.0 to 0.15.0 ([#2899](https://github.com/CooooolFrog/dragonfly/issues/2899))
- **deps:** bump golang.org/x/oauth2 from 0.13.0 to 0.14.0 ([#2889](https://github.com/CooooolFrog/dragonfly/issues/2889))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.45.0 to 0.46.1 ([#2888](https://github.com/CooooolFrog/dragonfly/issues/2888))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.13.0 to 2.13.1 ([#2890](https://github.com/CooooolFrog/dragonfly/issues/2890))
- **deps:** bump golang.org/x/sync from 0.4.0 to 0.5.0 ([#2887](https://github.com/CooooolFrog/dragonfly/issues/2887))
- **deps:** bump go.opentelemetry.io/otel from 1.20.0 to 1.21.0 ([#2886](https://github.com/CooooolFrog/dragonfly/issues/2886))

### Feat
- dfget sets data dir and cache dir ([#2931](https://github.com/CooooolFrog/dragonfly/issues/2931))
- add prefetch to trigger seed peer downloads entire task ([#2929](https://github.com/CooooolFrog/dragonfly/issues/2929))
- optimize service v2 log print ([#2922](https://github.com/CooooolFrog/dragonfly/issues/2922))
- optimize AnnouncePeer log and update api verison ([#2921](https://github.com/CooooolFrog/dragonfly/issues/2921))
- remove handleSyncPiecesFailedRequest when sync piece failed ([#2906](https://github.com/CooooolFrog/dragonfly/issues/2906))
- enable watch dog by default ([#2905](https://github.com/CooooolFrog/dragonfly/issues/2905))
- update api version to v2.0.56 ([#2891](https://github.com/CooooolFrog/dragonfly/issues/2891))

### Fix
- log parse error ([#2909](https://github.com/CooooolFrog/dragonfly/issues/2909))
- CANNOT connet to Redis with special password [#2893](https://github.com/CooooolFrog/dragonfly/issues/2893) ([#2894](https://github.com/CooooolFrog/dragonfly/issues/2894))
- digest and range validation in v2 ([#2892](https://github.com/CooooolFrog/dragonfly/issues/2892))


<a name="v2.1.23"></a>
## [v2.1.23] - 2023-11-16
### Chore
- update node version to 20-alpine in Dockerfile ([#2874](https://github.com/CooooolFrog/dragonfly/issues/2874))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.19.0 to 1.20.0 ([#2878](https://github.com/CooooolFrog/dragonfly/issues/2878))
- **deps:** bump github.com/aws/aws-sdk-go from 1.47.4 to 1.47.9 ([#2879](https://github.com/CooooolFrog/dragonfly/issues/2879))
- **deps:** bump k8s.io/component-base from 0.28.0 to 0.28.3 ([#2877](https://github.com/CooooolFrog/dragonfly/issues/2877))
- **deps:** bump golang.org/x/crypto from 0.14.0 to 0.15.0 ([#2876](https://github.com/CooooolFrog/dragonfly/issues/2876))

### Feat
- implement batch calculation of candidate parents scores ([#2853](https://github.com/CooooolFrog/dragonfly/issues/2853))
- avoid hot resolve in grpc ([#2884](https://github.com/CooooolFrog/dragonfly/issues/2884))
- optimize piece download failed handler ([#2883](https://github.com/CooooolFrog/dragonfly/issues/2883))
- add reschedule handler for schduler v2 ([#2882](https://github.com/CooooolFrog/dragonfly/issues/2882))
- remove TinyTaskResponse and SmallTaskResponse message ([#2881](https://github.com/CooooolFrog/dragonfly/issues/2881))
- change ratelimit's DefaultQPS and DefaultBurst ([#2872](https://github.com/CooooolFrog/dragonfly/issues/2872))


<a name="v2.1.22"></a>
## [v2.1.22] - 2023-11-09
### Feat
- add getObjectStorageMetadata to dfstore ([#2869](https://github.com/CooooolFrog/dragonfly/issues/2869))
- add MaxScheduleCount and GetSeedPeers ([#2868](https://github.com/CooooolFrog/dragonfly/issues/2868))


<a name="v2.1.21"></a>
## [v2.1.21] - 2023-11-08
### Chore
- **deps:** bump github.com/docker/distribution from 2.7.1+incompatible to 2.8.2+incompatible ([#2861](https://github.com/CooooolFrog/dragonfly/issues/2861))
- **deps:** bump github.com/containerd/containerd from 1.3.4 to 1.5.18 ([#2862](https://github.com/CooooolFrog/dragonfly/issues/2862))
- **deps:** bump github.com/docker/docker from 20.10.7+incompatible to 24.0.7+incompatible ([#2863](https://github.com/CooooolFrog/dragonfly/issues/2863))
- **deps:** bump github.com/aws/aws-sdk-go from 1.47.3 to 1.47.4 ([#2860](https://github.com/CooooolFrog/dragonfly/issues/2860))
- **deps:** bump golang.org/x/time from 0.3.0 to 0.4.0 ([#2855](https://github.com/CooooolFrog/dragonfly/issues/2855))
- **deps:** bump github.com/aws/aws-sdk-go from 1.46.6 to 1.47.3 ([#2858](https://github.com/CooooolFrog/dragonfly/issues/2858))
- **deps:** bump github.com/aliyun/aliyun-oss-go-sdk from 2.2.9+incompatible to 3.0.1+incompatible ([#2859](https://github.com/CooooolFrog/dragonfly/issues/2859))

### Feat
- implement ListSeedPeers feature in manager ([#2865](https://github.com/CooooolFrog/dragonfly/issues/2865))
- avoid hot reload when dynconfig refresh data ([#2866](https://github.com/CooooolFrog/dragonfly/issues/2866))
- if no matching manifest, then the log printing platform string ([#2867](https://github.com/CooooolFrog/dragonfly/issues/2867))
- preheat image supports authentication and parse multi manifest mediatype ([#2819](https://github.com/CooooolFrog/dragonfly/issues/2819))

### Refactor
- preheat with multi arch image layers ([#2864](https://github.com/CooooolFrog/dragonfly/issues/2864))


<a name="v2.1.20"></a>
## [v2.1.20] - 2023-11-03
### Chore
- update console verison to v1.0.20 ([#2850](https://github.com/CooooolFrog/dragonfly/issues/2850))


<a name="v2.1.19"></a>
## [v2.1.19] - 2023-11-01
### Chore
- update console and api verison ([#2847](https://github.com/CooooolFrog/dragonfly/issues/2847))
- split concurrentDownloadSource ([#2842](https://github.com/CooooolFrog/dragonfly/issues/2842))
- **deps:** bump gorm.io/driver/postgres from 1.5.3 to 1.5.4 ([#2837](https://github.com/CooooolFrog/dragonfly/issues/2837))
- **deps:** bump github.com/aws/aws-sdk-go from 1.45.25 to 1.46.6 ([#2841](https://github.com/CooooolFrog/dragonfly/issues/2841))
- **deps:** bump github.com/google/uuid from 1.3.1 to 1.4.0 ([#2840](https://github.com/CooooolFrog/dragonfly/issues/2840))
- **deps:** bump github.com/onsi/gomega from 1.28.0 to 1.29.0 ([#2839](https://github.com/CooooolFrog/dragonfly/issues/2839))
- **deps:** bump gorm.io/driver/mysql from 1.5.1 to 1.5.2 ([#2838](https://github.com/CooooolFrog/dragonfly/issues/2838))

### Feat
- use scheduler's cluster id in host resource ([#2844](https://github.com/CooooolFrog/dragonfly/issues/2844))
- list jobs order by created_at ([#2843](https://github.com/CooooolFrog/dragonfly/issues/2843))

### Fix
- getObject in oss ([#2846](https://github.com/CooooolFrog/dragonfly/issues/2846))


<a name="v2.1.18"></a>
## [v2.1.18] - 2023-10-30
### Feat
- print headers with debug log ([#2834](https://github.com/CooooolFrog/dragonfly/issues/2834))

### Fix
- if scheduler has no seed peer, return error in preheating ([#2835](https://github.com/CooooolFrog/dragonfly/issues/2835))
- filter Authorization header from preheat log ([#2817](https://github.com/CooooolFrog/dragonfly/issues/2817))


<a name="v2.1.17"></a>
## [v2.1.17] - 2023-10-26
### Chore
- update console version to v1.0.18 ([#2815](https://github.com/CooooolFrog/dragonfly/issues/2815))
- update console version to v1.0.15 ([#2813](https://github.com/CooooolFrog/dragonfly/issues/2813))
- update console verison to v1.0.14 ([#2811](https://github.com/CooooolFrog/dragonfly/issues/2811))
- update manager console verison to v1.0.13 ([#2810](https://github.com/CooooolFrog/dragonfly/issues/2810))
- **deps:** bump google.golang.org/grpc from 1.59.0-dev to 1.60.0-dev ([#2824](https://github.com/CooooolFrog/dragonfly/issues/2824))
- **deps:** bump gorm.io/driver/postgres from 1.5.2 to 1.5.3 ([#2823](https://github.com/CooooolFrog/dragonfly/issues/2823))
- **deps:** bump golang.org/x/sync from 0.3.0 to 0.4.0 ([#2822](https://github.com/CooooolFrog/dragonfly/issues/2822))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.8.0 to 1.10.0 ([#2821](https://github.com/CooooolFrog/dragonfly/issues/2821))
- **deps:** bump gorm.io/gorm from 1.25.4 to 1.25.5 ([#2809](https://github.com/CooooolFrog/dragonfly/issues/2809))
- **deps:** bump github.com/prometheus/client_golang from 1.16.0 to 1.17.0 ([#2807](https://github.com/CooooolFrog/dragonfly/issues/2807))
- **deps:** bump github.com/aws/aws-sdk-go from 1.45.6 to 1.45.25 ([#2806](https://github.com/CooooolFrog/dragonfly/issues/2806))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.43.0 to 0.45.0 ([#2808](https://github.com/CooooolFrog/dragonfly/issues/2808))

### Feat
- add sync peers log ([#2812](https://github.com/CooooolFrog/dragonfly/issues/2812))
- proxy graceful shutdown ([#2802](https://github.com/CooooolFrog/dragonfly/issues/2802))

### Fix
- console build failed ([#2814](https://github.com/CooooolFrog/dragonfly/issues/2814))


<a name="v2.1.16"></a>
## [v2.1.16] - 2023-10-16
### Chore
- add free disk space in e2e actions ([#2801](https://github.com/CooooolFrog/dragonfly/issues/2801))
- gofmt -w -r 'interface{} -> any' . ([#2790](https://github.com/CooooolFrog/dragonfly/issues/2790))
- optimize recursive download log ([#2785](https://github.com/CooooolFrog/dragonfly/issues/2785))
- **deps:** bump golang.org/x/net from 0.16.0 to 0.17.0 ([#2795](https://github.com/CooooolFrog/dragonfly/issues/2795))

### Feat
- change DownloadPeerDuration metric type to summary ([#2794](https://github.com/CooooolFrog/dragonfly/issues/2794))
- add content-length label to DownloadPeerDuration metric ([#2792](https://github.com/CooooolFrog/dragonfly/issues/2792))
- add comment about not checking object storage close ([#2789](https://github.com/CooooolFrog/dragonfly/issues/2789))

### Fix
- sync peer with unknow relation ([#2800](https://github.com/CooooolFrog/dragonfly/issues/2800))
- task length did not match range length ([#2787](https://github.com/CooooolFrog/dragonfly/issues/2787))


<a name="v2.1.15"></a>
## [v2.1.15] - 2023-10-10

<a name="v2.1.14"></a>
## [v2.1.14] - 2023-10-10
### Chore
- optimize recursive download log ([#2785](https://github.com/CooooolFrog/dragonfly/issues/2785))

### Fix
- task length did not match range length ([#2787](https://github.com/CooooolFrog/dragonfly/issues/2787))


<a name="v2.1.13"></a>
## [v2.1.13] - 2023-10-10
### Chore
- change gorm-adapter version to v3.5.0 ([#2774](https://github.com/CooooolFrog/dragonfly/issues/2774))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.17.0 to 1.19.0 ([#2782](https://github.com/CooooolFrog/dragonfly/issues/2782))
- **deps:** bump github.com/onsi/gomega from 1.27.10 to 1.28.0 ([#2779](https://github.com/CooooolFrog/dragonfly/issues/2779))
- **deps:** bump golang.org/x/oauth2 from 0.12.0 to 0.13.0 ([#2778](https://github.com/CooooolFrog/dragonfly/issues/2778))

### Feat
- change error handle in object storage ([#2777](https://github.com/CooooolFrog/dragonfly/issues/2777))

### Fix
- per page count in peer api ([#2775](https://github.com/CooooolFrog/dragonfly/issues/2775))


<a name="v2.1.12"></a>
## [v2.1.12] - 2023-10-08
### Chore
- update ginkgo version ([#2773](https://github.com/CooooolFrog/dragonfly/issues/2773))
- optimize dynconfig error log ([#2771](https://github.com/CooooolFrog/dragonfly/issues/2771))
- remove refs to deprecated io/ioutil ([#2754](https://github.com/CooooolFrog/dragonfly/issues/2754))
- update golang version to v1.21 ([#2753](https://github.com/CooooolFrog/dragonfly/issues/2753))
- **deps:** bump github.com/casbin/gorm-adapter/v3 from 3.5.0 to 3.20.0 ([#2766](https://github.com/CooooolFrog/dragonfly/issues/2766))
- **deps:** bump github.com/go-playground/validator/v10 from 10.15.4 to 10.15.5 ([#2765](https://github.com/CooooolFrog/dragonfly/issues/2765))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.43.0 to 0.45.0 ([#2764](https://github.com/CooooolFrog/dragonfly/issues/2764))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.8 to 3.23.9 ([#2763](https://github.com/CooooolFrog/dragonfly/issues/2763))
- **deps:** bump github.com/huaweicloud/huaweicloud-sdk-go-obs from 3.23.4+incompatible to 3.23.9+incompatible ([#2762](https://github.com/CooooolFrog/dragonfly/issues/2762))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.7 to 3.23.8 ([#2757](https://github.com/CooooolFrog/dragonfly/issues/2757))
- **deps:** bump github.com/gin-contrib/zap from 0.1.0 to 0.2.0 ([#2756](https://github.com/CooooolFrog/dragonfly/issues/2756))

### Feat
- update api to v2.0.29 with optional ([#2751](https://github.com/CooooolFrog/dragonfly/issues/2751))
- add CreatePeer to manager rest api ([#2749](https://github.com/CooooolFrog/dragonfly/issues/2749))

### Fix
- otelgrpc.UnaryClientInterceptor memory leak ([#2772](https://github.com/CooooolFrog/dragonfly/issues/2772))
- set download header log level debug ([#2770](https://github.com/CooooolFrog/dragonfly/issues/2770))
- add option to disable prepared statement for postgres ([#2768](https://github.com/CooooolFrog/dragonfly/issues/2768))
- overhead gc when DiskGCThreshold not set ([#2750](https://github.com/CooooolFrog/dragonfly/issues/2750))
- snapshot network topology id ([#2746](https://github.com/CooooolFrog/dragonfly/issues/2746))

### Test
- add unit test for GetLastModified ([#2747](https://github.com/CooooolFrog/dragonfly/issues/2747))
- add unit test for Notify ([#2748](https://github.com/CooooolFrog/dragonfly/issues/2748))


<a name="v2.1.11"></a>
## [v2.1.11] - 2023-09-19
### Chore
- **deps:** bump github.com/go-playground/validator/v10 from 10.15.1 to 10.15.4 ([#2741](https://github.com/CooooolFrog/dragonfly/issues/2741))
- **deps:** bump go.opentelemetry.io/otel/trace from 1.17.0 to 1.18.0 ([#2739](https://github.com/CooooolFrog/dragonfly/issues/2739))
- **deps:** bump golang.org/x/oauth2 from 0.11.0 to 0.12.0 ([#2738](https://github.com/CooooolFrog/dragonfly/issues/2738))
- **deps:** bump docker/setup-buildx-action from 2 to 3 ([#2737](https://github.com/CooooolFrog/dragonfly/issues/2737))
- **deps:** bump docker/login-action from 2 to 3 ([#2736](https://github.com/CooooolFrog/dragonfly/issues/2736))
- **deps:** bump docker/setup-qemu-action from 2 to 3 ([#2735](https://github.com/CooooolFrog/dragonfly/issues/2735))
- **deps:** bump goreleaser/goreleaser-action from 4 to 5 ([#2734](https://github.com/CooooolFrog/dragonfly/issues/2734))
- **deps:** bump docker/build-push-action from 4 to 5 ([#2733](https://github.com/CooooolFrog/dragonfly/issues/2733))

### Feat
- convert limit in dfstore ([#2745](https://github.com/CooooolFrog/dragonfly/issues/2745))

### Fix
- dfget couldn't download s3 directory correctly ([#2731](https://github.com/CooooolFrog/dragonfly/issues/2731))


<a name="v2.1.10"></a>
## [v2.1.10] - 2023-09-15
### Chore
- update console verison to v1.0.12 ([#2730](https://github.com/CooooolFrog/dragonfly/issues/2730))
- update helm charts verison to 1.1.6 ([#2726](https://github.com/CooooolFrog/dragonfly/issues/2726))

### Docs
- add security audit report ([#2729](https://github.com/CooooolFrog/dragonfly/issues/2729))

### Feat
- add InitProbedCount in AnnounceHost and DeleteHost in LeaveHost ([#2722](https://github.com/CooooolFrog/dragonfly/issues/2722))

### Fix
- if network topology is nil, LeaveHost panic in scheduler ([#2727](https://github.com/CooooolFrog/dragonfly/issues/2727))

### Test
- add unit test for Convert ([#2725](https://github.com/CooooolFrog/dragonfly/issues/2725))


<a name="v2.1.9"></a>
## [v2.1.9] - 2023-09-12
### Chore
- change timeout-minutes to 120 in docker actions ([#2724](https://github.com/CooooolFrog/dragonfly/issues/2724))
- update console submodule to v1.0.11 ([#2711](https://github.com/CooooolFrog/dragonfly/issues/2711))
- **deps:** bump github.com/aws/aws-sdk-go from 1.45.2 to 1.45.6 ([#2719](https://github.com/CooooolFrog/dragonfly/issues/2719))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.42.0 to 0.43.0 ([#2716](https://github.com/CooooolFrog/dragonfly/issues/2716))
- **deps:** bump github.com/casbin/casbin/v2 from 2.73.0 to 2.77.2 ([#2715](https://github.com/CooooolFrog/dragonfly/issues/2715))

### Feat
- add network timout to yarn install ([#2723](https://github.com/CooooolFrog/dragonfly/issues/2723))
- remove DefaultPersonalAccessTokenScopes in personalAccessToken ([#2708](https://github.com/CooooolFrog/dragonfly/issues/2708))
- add cluster rest api to open api ([#2707](https://github.com/CooooolFrog/dragonfly/issues/2707))

### Test
- add unit test for Items ([#2703](https://github.com/CooooolFrog/dragonfly/issues/2703))


<a name="v2.1.8"></a>
## [v2.1.8] - 2023-09-05
### Chore
- maintainer in alphabetical order of the company ([#2687](https://github.com/CooooolFrog/dragonfly/issues/2687))
- update api verison to v2.0.25 ([#2685](https://github.com/CooooolFrog/dragonfly/issues/2685))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.40.0 to 0.43.0 ([#2700](https://github.com/CooooolFrog/dragonfly/issues/2700))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.15.1 to 1.17.0 ([#2699](https://github.com/CooooolFrog/dragonfly/issues/2699))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.11.0 to 2.12.0 ([#2698](https://github.com/CooooolFrog/dragonfly/issues/2698))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.327 to 1.45.2 ([#2697](https://github.com/CooooolFrog/dragonfly/issues/2697))
- **deps:** bump github.com/google/uuid from 1.3.0 to 1.3.1 ([#2696](https://github.com/CooooolFrog/dragonfly/issues/2696))
- **deps:** bump actions/checkout from 3 to 4 ([#2695](https://github.com/CooooolFrog/dragonfly/issues/2695))
- **deps:** bump google.golang.org/api from 0.136.0 to 0.138.0 ([#2680](https://github.com/CooooolFrog/dragonfly/issues/2680))
- **deps:** bump github.com/go-playground/validator/v10 from 10.15.0 to 10.15.1 ([#2678](https://github.com/CooooolFrog/dragonfly/issues/2678))
- **deps:** bump go.opentelemetry.io/otel/trace from 1.16.0 to 1.17.0 ([#2677](https://github.com/CooooolFrog/dragonfly/issues/2677))
- **deps:** bump github.com/aliyun/aliyun-oss-go-sdk from 2.2.8+incompatible to 2.2.9+incompatible ([#2676](https://github.com/CooooolFrog/dragonfly/issues/2676))

### Docs
- remove china mobile in ADOPTERS.md ([#2690](https://github.com/CooooolFrog/dragonfly/issues/2690))
- add Akash HR to maintainers ([#2673](https://github.com/CooooolFrog/dragonfly/issues/2673))

### Feat
- change file mode when export task from local in dfcache ([#2701](https://github.com/CooooolFrog/dragonfly/issues/2701))
- use go-version-file: go.mod for go version detecting ([#2694](https://github.com/CooooolFrog/dragonfly/issues/2694))
- optimize GetObjectMetadatasInput validate in dfstore ([#2686](https://github.com/CooooolFrog/dragonfly/issues/2686))
- initialize object storage client in the New func ([#2682](https://github.com/CooooolFrog/dragonfly/issues/2682))
- remove peer index in databae table ([#2675](https://github.com/CooooolFrog/dragonfly/issues/2675))
- add idc and location to ListSchedulers in manager ([#2674](https://github.com/CooooolFrog/dragonfly/issues/2674))

### Fix
- Vertex.DeleteInEdges and Vertex.DeleteOutEdges functions are not thread safe ([#2688](https://github.com/CooooolFrog/dragonfly/issues/2688))

### Test
- add unit test for safeSocketControl ([#2684](https://github.com/CooooolFrog/dragonfly/issues/2684))
- add unit test for MarshalResponse ([#2683](https://github.com/CooooolFrog/dragonfly/issues/2683))


<a name="v2.1.7"></a>
## [v2.1.7] - 2023-08-25
### Chore
- **deps:** bump github.com/jarcoal/httpmock from 1.3.0 to 1.3.1 ([#2660](https://github.com/CooooolFrog/dragonfly/issues/2660))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.317 to 1.44.327 ([#2659](https://github.com/CooooolFrog/dragonfly/issues/2659))
- **deps:** bump gorm.io/gorm from 1.25.3 to 1.25.4 ([#2658](https://github.com/CooooolFrog/dragonfly/issues/2658))
- **deps:** bump k8s.io/component-base from 0.27.4 to 0.28.0 ([#2657](https://github.com/CooooolFrog/dragonfly/issues/2657))

### Docs
- add OWNERS.md to project ([#2671](https://github.com/CooooolFrog/dragonfly/issues/2671))
- add GOVERNANCE.md to dragonfly project ([#2669](https://github.com/CooooolFrog/dragonfly/issues/2669))

### Feat
- remove IsPrivate from safeSocketControl ([#2672](https://github.com/CooooolFrog/dragonfly/issues/2672))
- merge sync peer with peer table in manager ([#2668](https://github.com/CooooolFrog/dragonfly/issues/2668))
- add createSyncPeers to async job in manager ([#2664](https://github.com/CooooolFrog/dragonfly/issues/2664))
- add sync peer job for scheduler ([#2663](https://github.com/CooooolFrog/dragonfly/issues/2663))
- peer information is changed from being stored in metrics to being stored in mysql ([#2654](https://github.com/CooooolFrog/dragonfly/issues/2654))
- peer announces scheduler cluster id to scheduler ([#2652](https://github.com/CooooolFrog/dragonfly/issues/2652))

### Test
- add unit test for PieceDownloader.Error ([#2667](https://github.com/CooooolFrog/dragonfly/issues/2667))
- add unit test for DfgetConfig.Validate ([#2666](https://github.com/CooooolFrog/dragonfly/issues/2666))
- add unit test for Dfget.parseHeader ([#2665](https://github.com/CooooolFrog/dragonfly/issues/2665))


<a name="v2.1.6"></a>
## [v2.1.6] - 2023-08-21
### Chore
- update console version to v1.0.9 ([#2656](https://github.com/CooooolFrog/dragonfly/issues/2656))


<a name="v2.1.5"></a>
## [v2.1.5] - 2023-08-21
### Chore
- update console version to 1.0.8 ([#2655](https://github.com/CooooolFrog/dragonfly/issues/2655))
- update console version to v1.0.7 ([#2651](https://github.com/CooooolFrog/dragonfly/issues/2651))
- update console version to v1.0.6 ([#2650](https://github.com/CooooolFrog/dragonfly/issues/2650))

### Feat
- change max PerPage to 10000000 ([#2653](https://github.com/CooooolFrog/dragonfly/issues/2653))
- replace fmt.Sprintf with net.JoinHostPort ([#2642](https://github.com/CooooolFrog/dragonfly/issues/2642))

### Test
- add unit test for request.WithContext ([#2646](https://github.com/CooooolFrog/dragonfly/issues/2646))
- add unit tests for byte ([#2645](https://github.com/CooooolFrog/dragonfly/issues/2645))
- rename Test_Validate to TestDfstoreConfig_Validate ([#2644](https://github.com/CooooolFrog/dragonfly/issues/2644))
- add unit tests for dfstore ([#2643](https://github.com/CooooolFrog/dragonfly/issues/2643))


<a name="v2.1.4"></a>
## [v2.1.4] - 2023-08-16
### Chore
- update console to v1.0.5 ([#2640](https://github.com/CooooolFrog/dragonfly/issues/2640))


<a name="v2.1.3"></a>
## [v2.1.3] - 2023-08-15
### Chore
- update console to v1.0.4 ([#2639](https://github.com/CooooolFrog/dragonfly/issues/2639))
- skip export to exist file ([#2637](https://github.com/CooooolFrog/dragonfly/issues/2637))
- check task id ([#2636](https://github.com/CooooolFrog/dragonfly/issues/2636))

### Fix
- add traffic ([#2634](https://github.com/CooooolFrog/dragonfly/issues/2634))


<a name="v2.1.2"></a>
## [v2.1.2] - 2023-08-15
### Chore
- update console to v1.0.3 ([#2638](https://github.com/CooooolFrog/dragonfly/issues/2638))


<a name="v2.1.1"></a>
## [v2.1.1] - 2023-08-15
### Chore
- update console version to v1.0.2 ([#2635](https://github.com/CooooolFrog/dragonfly/issues/2635))
- update d7y.io/api version to v2.0.18 ([#2616](https://github.com/CooooolFrog/dragonfly/issues/2616))
- optimize unhandled file close error ([#2599](https://github.com/CooooolFrog/dragonfly/issues/2599))
- use subtle compare to verify proxy auth ([#2601](https://github.com/CooooolFrog/dragonfly/issues/2601))
- **deps:** bump gorm.io/gorm from 1.25.2 to 1.25.3 ([#2630](https://github.com/CooooolFrog/dragonfly/issues/2630))
- **deps:** bump go.uber.org/zap from 1.24.0 to 1.25.0 ([#2629](https://github.com/CooooolFrog/dragonfly/issues/2629))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.6 to 3.23.7 ([#2628](https://github.com/CooooolFrog/dragonfly/issues/2628))
- **deps:** bump google.golang.org/api from 0.134.0 to 0.136.0 ([#2626](https://github.com/CooooolFrog/dragonfly/issues/2626))
- **deps:** bump github.com/go-playground/validator/v10 from 10.14.1 to 10.15.0 ([#2608](https://github.com/CooooolFrog/dragonfly/issues/2608))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.312 to 1.44.317 ([#2606](https://github.com/CooooolFrog/dragonfly/issues/2606))
- **deps:** bump golang.org/x/oauth2 from 0.10.0 to 0.11.0 ([#2605](https://github.com/CooooolFrog/dragonfly/issues/2605))

### Feat
- change cache size in manager ([#2633](https://github.com/CooooolFrog/dragonfly/issues/2633))
- case insensitive string comparison in evaluater package of the manager ([#2632](https://github.com/CooooolFrog/dragonfly/issues/2632))
- change cache size in manager ([#2623](https://github.com/CooooolFrog/dragonfly/issues/2623))
- net.JoinHostPort replace fmt.Sprintf ([#2622](https://github.com/CooooolFrog/dragonfly/issues/2622))
- download tiny file with https scheme ([#2617](https://github.com/CooooolFrog/dragonfly/issues/2617))
- add lock to dag vertex and replace rand.Seed with rand.New ([#2614](https://github.com/CooooolFrog/dragonfly/issues/2614))
- add tls client config for preheat in manager ([#2612](https://github.com/CooooolFrog/dragonfly/issues/2612))
- add NewSafeDialer and fix ssrf in manager preheat api ([#2611](https://github.com/CooooolFrog/dragonfly/issues/2611))
- change default cluster name to cluster-1 ([#2604](https://github.com/CooooolFrog/dragonfly/issues/2604))

### Fix
- list personal access token with query string ([#2624](https://github.com/CooooolFrog/dragonfly/issues/2624))
- incorrect log message in the scheduler ([#2618](https://github.com/CooooolFrog/dragonfly/issues/2618))
- usage of architecture-dependent int type in the scheduler ([#2619](https://github.com/CooooolFrog/dragonfly/issues/2619))
- manager generates mTLS certificates for arbitrary IP addresses ([#2615](https://github.com/CooooolFrog/dragonfly/issues/2615))
- directories created via os.MkdirAll are not checked for permissions ([#2613](https://github.com/CooooolFrog/dragonfly/issues/2613))
- invalid error handling ([#2610](https://github.com/CooooolFrog/dragonfly/issues/2610))
- improper use strings.TrimLeft ([#2603](https://github.com/CooooolFrog/dragonfly/issues/2603))
- potential nil panic ([#2602](https://github.com/CooooolFrog/dragonfly/issues/2602))

### Refactor
- file close error


<a name="v2.1.0"></a>
## [v2.1.0] - 2023-08-04
### Chore
- release v2.1.0 ([#2597](https://github.com/CooooolFrog/dragonfly/issues/2597))


<a name="v2.1.0-rc.0"></a>
## [v2.1.0-rc.0] - 2023-08-03
### Chore
- update submodule version ([#2596](https://github.com/CooooolFrog/dragonfly/issues/2596))

### Feat
- use unscoped delete for resource in manager ([#2595](https://github.com/CooooolFrog/dragonfly/issues/2595))
- create seed peer with active state in manager ([#2593](https://github.com/CooooolFrog/dragonfly/issues/2593))
- change seed peer state to active in UpdateSeedPeer api ([#2592](https://github.com/CooooolFrog/dragonfly/issues/2592))
- implement DeleteSeedPeer api in manager ([#2591](https://github.com/CooooolFrog/dragonfly/issues/2591))


<a name="v2.1.0-beta.4"></a>
## [v2.1.0-beta.4] - 2023-08-01
### Chore
- clean temporary file when backsource error ([#2575](https://github.com/CooooolFrog/dragonfly/issues/2575))
- **deps:** bump github.com/aliyun/aliyun-oss-go-sdk from 2.2.7+incompatible to 2.2.8+incompatible ([#2588](https://github.com/CooooolFrog/dragonfly/issues/2588))
- **deps:** bump google.golang.org/api from 0.132.0 to 0.134.0 ([#2587](https://github.com/CooooolFrog/dragonfly/issues/2587))
- **deps:** bump github.com/onsi/gomega from 1.27.8 to 1.27.10 ([#2586](https://github.com/CooooolFrog/dragonfly/issues/2586))
- **deps:** bump github.com/casbin/casbin/v2 from 2.72.1 to 2.73.0 ([#2585](https://github.com/CooooolFrog/dragonfly/issues/2585))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.306 to 1.44.312 ([#2584](https://github.com/CooooolFrog/dragonfly/issues/2584))

### Feat
- add personal access token middleware to open api ([#2590](https://github.com/CooooolFrog/dragonfly/issues/2590))
- add personal access tokens api to rest server ([#2583](https://github.com/CooooolFrog/dragonfly/issues/2583))
- support tls in manager rest server ([#2580](https://github.com/CooooolFrog/dragonfly/issues/2580))
- provide support for JuiceFS objectStorage implementation ([#2578](https://github.com/CooooolFrog/dragonfly/issues/2578))
- update api version ([#2577](https://github.com/CooooolFrog/dragonfly/issues/2577))

### Refactor
- support for JuiceFS objectStorage implementation ([#2579](https://github.com/CooooolFrog/dragonfly/issues/2579))

### Test
- add unit test for seedPeerClient.Addrs ([#2589](https://github.com/CooooolFrog/dragonfly/issues/2589))
- add unit tests for PieceDownloader ([#2570](https://github.com/CooooolFrog/dragonfly/issues/2570))


<a name="v2.1.0-beta.3"></a>
## [v2.1.0-beta.3] - 2023-07-26
### Feat
- change per_page to 1000 ([#2576](https://github.com/CooooolFrog/dragonfly/issues/2576))


<a name="v2.1.0-beta.2"></a>
## [v2.1.0-beta.2] - 2023-07-25
### Chore
- change tainer address port from 9000 to 9090 in scheduler ([#2571](https://github.com/CooooolFrog/dragonfly/issues/2571))
- change trainer expose port from 8002 to 9090 in Dockerfile ([#2569](https://github.com/CooooolFrog/dragonfly/issues/2569))
- add fcgxz2003 to maintainer ([#2522](https://github.com/CooooolFrog/dragonfly/issues/2522))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.5 to 3.23.6 ([#2535](https://github.com/CooooolFrog/dragonfly/issues/2535))
- **deps:** bump google.golang.org/api from 0.129.0 to 0.130.0 ([#2533](https://github.com/CooooolFrog/dragonfly/issues/2533))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.258 to 1.44.306 ([#2565](https://github.com/CooooolFrog/dragonfly/issues/2565))
- **deps:** bump helm/kind-action from 1.7.0 to 1.8.0 ([#2553](https://github.com/CooooolFrog/dragonfly/issues/2553))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.7.0 to 1.8.0 ([#2551](https://github.com/CooooolFrog/dragonfly/issues/2551))
- **deps:** bump google.golang.org/api from 0.130.0 to 0.131.0 ([#2549](https://github.com/CooooolFrog/dragonfly/issues/2549))
- **deps:** bump github.com/MysteriousPotato/go-lockable from 0.2.0 to 1.0.0 ([#2548](https://github.com/CooooolFrog/dragonfly/issues/2548))
- **deps:** bump k8s.io/component-base from 0.27.2 to 0.27.4 ([#2562](https://github.com/CooooolFrog/dragonfly/issues/2562))
- **deps:** bump gorm.io/driver/postgres from 1.5.0 to 1.5.2 ([#2534](https://github.com/CooooolFrog/dragonfly/issues/2534))
- **deps:** bump google.golang.org/api from 0.131.0 to 0.132.0 ([#2564](https://github.com/CooooolFrog/dragonfly/issues/2564))
- **deps:** bump golang.org/x/oauth2 from 0.9.0 to 0.10.0 ([#2532](https://github.com/CooooolFrog/dragonfly/issues/2532))
- **deps:** bump github.com/casbin/casbin/v2 from 2.71.1 to 2.72.1 ([#2561](https://github.com/CooooolFrog/dragonfly/issues/2561))
- **deps:** bump gorm.io/gorm from 1.25.1 to 1.25.2 ([#2505](https://github.com/CooooolFrog/dragonfly/issues/2505))
- **deps:** bump github.com/MysteriousPotato/go-lockable from 0.1.0 to 0.2.0 ([#2504](https://github.com/CooooolFrog/dragonfly/issues/2504))
- **deps:** bump google.golang.org/api from 0.128.0 to 0.129.0 ([#2503](https://github.com/CooooolFrog/dragonfly/issues/2503))
- **deps:** bump google.golang.org/protobuf from 1.30.0 to 1.31.0 ([#2502](https://github.com/CooooolFrog/dragonfly/issues/2502))
- **deps:** bump github.com/casbin/casbin/v2 from 2.68.0 to 2.71.1 ([#2501](https://github.com/CooooolFrog/dragonfly/issues/2501))

### Feat
- add optional to build information ([#2567](https://github.com/CooooolFrog/dragonfly/issues/2567))
- update api version to 2.0.8 ([#2566](https://github.com/CooooolFrog/dragonfly/issues/2566))
- update api version to v2.0.7 ([#2560](https://github.com/CooooolFrog/dragonfly/issues/2560))
- update dy7.io/api to v2 ([#2558](https://github.com/CooooolFrog/dragonfly/issues/2558))
- add finished piece count element in types ([#2557](https://github.com/CooooolFrog/dragonfly/issues/2557))
- update api verison to v1.9.7 and compatible with tiny task response ([#2547](https://github.com/CooooolFrog/dragonfly/issues/2547))
- add training service ([#2543](https://github.com/CooooolFrog/dragonfly/issues/2543))
- implement Train grpc api in trainer ([#2541](https://github.com/CooooolFrog/dragonfly/issues/2541))
- add GNNModelIDV1 and MLPModelIDV1 in idgen ([#2540](https://github.com/CooooolFrog/dragonfly/issues/2540))
- add pieces element in download record ([#2531](https://github.com/CooooolFrog/dragonfly/issues/2531))
- remove ProbedAt function in network topology ([#2529](https://github.com/CooooolFrog/dragonfly/issues/2529))
- add update model rest api ([#2530](https://github.com/CooooolFrog/dragonfly/issues/2530))
- manager adds createModel function ([#2521](https://github.com/CooooolFrog/dragonfly/issues/2521))
- implement FindProbedHosts and add LoadRandomHosts to host manager ([#2519](https://github.com/CooooolFrog/dragonfly/issues/2519))
- set scan count for redis client ([#2508](https://github.com/CooooolFrog/dragonfly/issues/2508))
- replace keys with scan in redis client ([#2507](https://github.com/CooooolFrog/dragonfly/issues/2507))
- manager adds trainer config ([#2494](https://github.com/CooooolFrog/dragonfly/issues/2494))
- add inference client in grpc ([#2493](https://github.com/CooooolFrog/dragonfly/issues/2493))

### Fix
- remove archives.rlcp in .goreleaser.yaml refer to https://gorele… ([#2573](https://github.com/CooooolFrog/dragonfly/issues/2573))
- response of cluster rest api ([#2572](https://github.com/CooooolFrog/dragonfly/issues/2572))
- if condition judgment of clearing file in trainer service ([#2544](https://github.com/CooooolFrog/dragonfly/issues/2544))
- storage and announcer unit tests ([#2542](https://github.com/CooooolFrog/dragonfly/issues/2542))
- change model state in the same scheduler id ([#2537](https://github.com/CooooolFrog/dragonfly/issues/2537))
- scheduler.template.yaml comments ([#2526](https://github.com/CooooolFrog/dragonfly/issues/2526))
- also add ca to RootCAs ([#2516](https://github.com/CooooolFrog/dragonfly/issues/2516))

### Refactor
- store pieceRecords in download record ([#2539](https://github.com/CooooolFrog/dragonfly/issues/2539))
- create model grpc api in manager ([#2528](https://github.com/CooooolFrog/dragonfly/issues/2528))

### Test
- add unit test for Header.get ([#2568](https://github.com/CooooolFrog/dragonfly/issues/2568))
- add unit test for Header.has ([#2555](https://github.com/CooooolFrog/dragonfly/issues/2555))
- add unit test for Request.Context() ([#2554](https://github.com/CooooolFrog/dragonfly/issues/2554))
- trainer serivce unit tests ([#2545](https://github.com/CooooolFrog/dragonfly/issues/2545))
- add unit test for model and digest ([#2538](https://github.com/CooooolFrog/dragonfly/issues/2538))
- optimize TestDigest_HashFile ([#2515](https://github.com/CooooolFrog/dragonfly/issues/2515))
- improve hash file encode test case in digest test ([#2513](https://github.com/CooooolFrog/dragonfly/issues/2513))
- add unit tests for DownloadCount and NetworkTopologyCount ([#2512](https://github.com/CooooolFrog/dragonfly/issues/2512))
- replace keys with scan in redis client ([#2509](https://github.com/CooooolFrog/dragonfly/issues/2509))
- optimize config in scheduler ([#2511](https://github.com/CooooolFrog/dragonfly/issues/2511))


<a name="v2.1.0-beta.1"></a>
## [v2.1.0-beta.1] - 2023-06-28
### Chore
- update grpc api proto verison ([#1779](https://github.com/CooooolFrog/dragonfly/issues/1779))
- goreleaser remove cdn
- add pull request and issue templates ([#154](https://github.com/CooooolFrog/dragonfly/issues/154))
- create custom issue template ([#168](https://github.com/CooooolFrog/dragonfly/issues/168))
- change codeowners to dragonfly2's maintainers and reviewers ([#169](https://github.com/CooooolFrog/dragonfly/issues/169))
- change codeowners ([#179](https://github.com/CooooolFrog/dragonfly/issues/179))
- add SECURITY.md ([#181](https://github.com/CooooolFrog/dragonfly/issues/181))
- change manager swagger docs path and add makefile swagger command ([#183](https://github.com/CooooolFrog/dragonfly/issues/183))
- workflows remove main-rc branch ([#221](https://github.com/CooooolFrog/dragonfly/issues/221))
- update grpc proto version ([#2463](https://github.com/CooooolFrog/dragonfly/issues/2463))
- update dfget recursive log ([#2459](https://github.com/CooooolFrog/dragonfly/issues/2459))
- update grpc api definition to v1.9.0 ([#2444](https://github.com/CooooolFrog/dragonfly/issues/2444))
- remove manager netcat-openbsd ([#298](https://github.com/CooooolFrog/dragonfly/issues/298))
- docker building workflow ([#323](https://github.com/CooooolFrog/dragonfly/issues/323))
- remove build script's git operation ([#321](https://github.com/CooooolFrog/dragonfly/issues/321))
- update CI timeout ([#328](https://github.com/CooooolFrog/dragonfly/issues/328))
- remove protoc.sh ([#341](https://github.com/CooooolFrog/dragonfly/issues/341))
- change bash to sh ([#383](https://github.com/CooooolFrog/dragonfly/issues/383))
- add docs for dragonfly2.0 ([#234](https://github.com/CooooolFrog/dragonfly/issues/234))
- remove macos ci ([#404](https://github.com/CooooolFrog/dragonfly/issues/404))
- check grpc peer info for download service ([#2385](https://github.com/CooooolFrog/dragonfly/issues/2385))
- rename dfdaemon docker image ([#405](https://github.com/CooooolFrog/dragonfly/issues/405))
- remove goreleaser go generate ([#409](https://github.com/CooooolFrog/dragonfly/issues/409))
- custom charts template namespace ([#416](https://github.com/CooooolFrog/dragonfly/issues/416))
- set GOPROXY with default value ([#463](https://github.com/CooooolFrog/dragonfly/issues/463))
- change gorm-adaptor version to v3.5.0 ([#2370](https://github.com/CooooolFrog/dragonfly/issues/2370))
- optimize compute piece size function ([#528](https://github.com/CooooolFrog/dragonfly/issues/528))
- optimize grpc interceptor code ([#536](https://github.com/CooooolFrog/dragonfly/issues/536))
- optimize client rpc package name and other docs ([#541](https://github.com/CooooolFrog/dragonfly/issues/541))
- checkout code first in CI ([#2347](https://github.com/CooooolFrog/dragonfly/issues/2347))
- checkout code first in CI ([#2346](https://github.com/CooooolFrog/dragonfly/issues/2346))
- update redis config in docker compose and update helm chart version ([#2344](https://github.com/CooooolFrog/dragonfly/issues/2344))
- optimize peer task report function ([#543](https://github.com/CooooolFrog/dragonfly/issues/543))
- rename cdn server package to rpcserver ([#554](https://github.com/CooooolFrog/dragonfly/issues/554))
- add copyright ([#593](https://github.com/CooooolFrog/dragonfly/issues/593))
- add compatibility test workflow ([#594](https://github.com/CooooolFrog/dragonfly/issues/594))
- optimize app and tracer log ([#607](https://github.com/CooooolFrog/dragonfly/issues/607))
- update timeout in actions ([#2320](https://github.com/CooooolFrog/dragonfly/issues/2320))
- update submodule version ([#608](https://github.com/CooooolFrog/dragonfly/issues/608))
- update changelog ([#622](https://github.com/CooooolFrog/dragonfly/issues/622))
- skip workflows ([#624](https://github.com/CooooolFrog/dragonfly/issues/624))
- rename cdnsystem to cdn ([#626](https://github.com/CooooolFrog/dragonfly/issues/626))
- skip e2e ([#631](https://github.com/CooooolFrog/dragonfly/issues/631))
- compatibility with v2.0.0 test ([#639](https://github.com/CooooolFrog/dragonfly/issues/639))
- makefile typo
- add lucy-cl maintainer ([#645](https://github.com/CooooolFrog/dragonfly/issues/645))
- update oras error format ([#2282](https://github.com/CooooolFrog/dragonfly/issues/2282))
- update version ([#647](https://github.com/CooooolFrog/dragonfly/issues/647))
- change zzy987 maintainers email ([#649](https://github.com/CooooolFrog/dragonfly/issues/649))
- optimize advertise ip ([#652](https://github.com/CooooolFrog/dragonfly/issues/652))
- update build package config ([#653](https://github.com/CooooolFrog/dragonfly/issues/653))
- enable calculate digest ([#656](https://github.com/CooooolFrog/dragonfly/issues/656))
- add ChatGPT Code Review to workflows ([#2251](https://github.com/CooooolFrog/dragonfly/issues/2251))
- change timeout to 60m in docker workflows ([#2274](https://github.com/CooooolFrog/dragonfly/issues/2274))
- change dingtalk-group qrcode ([#2267](https://github.com/CooooolFrog/dragonfly/issues/2267))
- update dingtalk group qrcode ([#2262](https://github.com/CooooolFrog/dragonfly/issues/2262))
- export set log level ([#646](https://github.com/CooooolFrog/dragonfly/issues/646))
- e2e workflows remove goproxy ([#677](https://github.com/CooooolFrog/dragonfly/issues/677))
- change gorm-adaptor version to v3.5.0 ([#2247](https://github.com/CooooolFrog/dragonfly/issues/2247))
- add features swagger config ([#2246](https://github.com/CooooolFrog/dragonfly/issues/2246))
- remove skip-duplicate-actions ([#690](https://github.com/CooooolFrog/dragonfly/issues/690))
- workflows ignore paths ([#697](https://github.com/CooooolFrog/dragonfly/issues/697))
- release image to docker.pkg.github.com ([#703](https://github.com/CooooolFrog/dragonfly/issues/703))
- update config example ([#721](https://github.com/CooooolFrog/dragonfly/issues/721))
- change docker registry name ([#725](https://github.com/CooooolFrog/dragonfly/issues/725))
- update traffic shaper log ([#2227](https://github.com/CooooolFrog/dragonfly/issues/2227))
- repository name
- optimize span context for report ([#747](https://github.com/CooooolFrog/dragonfly/issues/747))
- check empty registry mirror ([#761](https://github.com/CooooolFrog/dragonfly/issues/761))
- optimize stream peer task ([#763](https://github.com/CooooolFrog/dragonfly/issues/763))
- update golang import lint ([#780](https://github.com/CooooolFrog/dragonfly/issues/780))
- format ci action
- add Mohammed Farooq to MAINTAINERS ([#2211](https://github.com/CooooolFrog/dragonfly/issues/2211))
- add markdown lint ([#779](https://github.com/CooooolFrog/dragonfly/issues/779))
- optimize client storage gc log ([#790](https://github.com/CooooolFrog/dragonfly/issues/790))
- add lint errcheck  and fix errcheck([#766](https://github.com/CooooolFrog/dragonfly/issues/766))
- unify binary directory ([#828](https://github.com/CooooolFrog/dragonfly/issues/828))
- upgrade to golang 1.17 and alpine 3.14 ([#861](https://github.com/CooooolFrog/dragonfly/issues/861))
- update changelog
- update UnknownSourceFileLen ([#888](https://github.com/CooooolFrog/dragonfly/issues/888))
- update nydus-snapshotter helm-charts to v0.0.4 ([#2188](https://github.com/CooooolFrog/dragonfly/issues/2188))
- migrate from k8s.gcr.io to registry.k8s.io ([#2186](https://github.com/CooooolFrog/dragonfly/issues/2186))
- change the compatibility testing version of manager and scheduler to v2.0.9 ([#2184](https://github.com/CooooolFrog/dragonfly/issues/2184))
- add build-man-page to makefile ([#2182](https://github.com/CooooolFrog/dragonfly/issues/2182))
- release v2.0.9 and generate changelog ([#2181](https://github.com/CooooolFrog/dragonfly/issues/2181))
- change codecov rules ([#2174](https://github.com/CooooolFrog/dragonfly/issues/2174))
- support multi daemons e2e test ([#896](https://github.com/CooooolFrog/dragonfly/issues/896))
- optimize back source update digest logic ([#950](https://github.com/CooooolFrog/dragonfly/issues/950))
- add version metric ([#954](https://github.com/CooooolFrog/dragonfly/issues/954))
- copy e2e proxy log to artifact ([#962](https://github.com/CooooolFrog/dragonfly/issues/962))
- change docker.pkg.github.com to ghcr.io ([#973](https://github.com/CooooolFrog/dragonfly/issues/973))
- clarify daemon interface ([#991](https://github.com/CooooolFrog/dragonfly/issues/991))
- parameterize tests in peer task ([#994](https://github.com/CooooolFrog/dragonfly/issues/994))
- sync docker-compose scheduler config ([#1001](https://github.com/CooooolFrog/dragonfly/issues/1001))
- add Garen Wen to maintainers ([#2136](https://github.com/CooooolFrog/dragonfly/issues/2136))
- workflow add test timeout ([#1011](https://github.com/CooooolFrog/dragonfly/issues/1011))
- optimize defer and test ([#1010](https://github.com/CooooolFrog/dragonfly/issues/1010))
- register to scheduler after updated running tasks ([#1016](https://github.com/CooooolFrog/dragonfly/issues/1016))
- remove unused MarkInvalid in daemon ([#2101](https://github.com/CooooolFrog/dragonfly/issues/2101))
- optimize metrics and trace in daemon ([#1022](https://github.com/CooooolFrog/dragonfly/issues/1022))
- update outdated log ([#1028](https://github.com/CooooolFrog/dragonfly/issues/1028))
- add piece task metrics in daemon ([#1030](https://github.com/CooooolFrog/dragonfly/issues/1030))
- upgrade to ginkgo v2 ([#1036](https://github.com/CooooolFrog/dragonfly/issues/1036))
- add missing pod log volumes in e2e ([#1037](https://github.com/CooooolFrog/dragonfly/issues/1037))
- use buildx to build docker images in e2e ([#1018](https://github.com/CooooolFrog/dragonfly/issues/1018))
- optimize https pass through ([#1054](https://github.com/CooooolFrog/dragonfly/issues/1054))
- add content length for fast stream peer task ([#1061](https://github.com/CooooolFrog/dragonfly/issues/1061))
- change e2e timeout ([#2062](https://github.com/CooooolFrog/dragonfly/issues/2062))
- add miHoYo to ADOPTERS.md ([#2054](https://github.com/CooooolFrog/dragonfly/issues/2054))
- enable range feature gate in e2e ([#1059](https://github.com/CooooolFrog/dragonfly/issues/1059))
- update gorelease ldflags ([#1086](https://github.com/CooooolFrog/dragonfly/issues/1086))
- init url meta in rpc server ([#1098](https://github.com/CooooolFrog/dragonfly/issues/1098))
- optimize reuse logic ([#1110](https://github.com/CooooolFrog/dragonfly/issues/1110))
- fast back source when get pieces task failed ([#1123](https://github.com/CooooolFrog/dragonfly/issues/1123))
- change scheduler config ([#1140](https://github.com/CooooolFrog/dragonfly/issues/1140))
- update issue templates ([#2041](https://github.com/CooooolFrog/dragonfly/issues/2041))
- change maintainers informations ([#2038](https://github.com/CooooolFrog/dragonfly/issues/2038))
- ignore configs generate with docker compose ([#2034](https://github.com/CooooolFrog/dragonfly/issues/2034))
- add makefile note ([#1155](https://github.com/CooooolFrog/dragonfly/issues/1155))
- update go mod ([#1156](https://github.com/CooooolFrog/dragonfly/issues/1156))
- always fallback to legacy get pieces ([#1180](https://github.com/CooooolFrog/dragonfly/issues/1180))
- optimize stream peer task ([#1186](https://github.com/CooooolFrog/dragonfly/issues/1186))
- change golangci-lint min-complexity value ([#1188](https://github.com/CooooolFrog/dragonfly/issues/1188))
- update workflows compatibility version ([#1192](https://github.com/CooooolFrog/dragonfly/issues/1192))
- report client back source error ([#1209](https://github.com/CooooolFrog/dragonfly/issues/1209))
- print client stream task error log ([#1210](https://github.com/CooooolFrog/dragonfly/issues/1210))
- update manager console commit ([#1219](https://github.com/CooooolFrog/dragonfly/issues/1219))
- fix workflows typo ([#2013](https://github.com/CooooolFrog/dragonfly/issues/2013))
- generate manager swagger ([#2009](https://github.com/CooooolFrog/dragonfly/issues/2009))
- generate change log
- update helm-charts commit
- update compatibility version to v2.0.2
- update pull request template ([#1251](https://github.com/CooooolFrog/dragonfly/issues/1251))
- update helm charts submodule ([#1997](https://github.com/CooooolFrog/dragonfly/issues/1997))
- optimize sync pieces ([#1253](https://github.com/CooooolFrog/dragonfly/issues/1253))
- add schedule cron with e2e testing ([#1262](https://github.com/CooooolFrog/dragonfly/issues/1262))
- add sync pieces trace and update sync pieces logic for done task ([#1263](https://github.com/CooooolFrog/dragonfly/issues/1263))
- optimize create synchronizer logic ([#1269](https://github.com/CooooolFrog/dragonfly/issues/1269))
- add target peer id in sync piece trace ([#1278](https://github.com/CooooolFrog/dragonfly/issues/1278))
- remove codecov patch feature ([#1977](https://github.com/CooooolFrog/dragonfly/issues/1977))
- update e2e timeout ([#1969](https://github.com/CooooolFrog/dragonfly/issues/1969))
- update charts version ([#1968](https://github.com/CooooolFrog/dragonfly/issues/1968))
- goreleaser set rlcp field ([#1967](https://github.com/CooooolFrog/dragonfly/issues/1967))
- update actions ([#1966](https://github.com/CooooolFrog/dragonfly/issues/1966))
- print e2e exec output ([#1963](https://github.com/CooooolFrog/dragonfly/issues/1963))
- change codecov coverage range ([#1965](https://github.com/CooooolFrog/dragonfly/issues/1965))
- check large files in pull request ([#1332](https://github.com/CooooolFrog/dragonfly/issues/1332))
- add check size action ([#1350](https://github.com/CooooolFrog/dragonfly/issues/1350))
- update content range for partial content ([#1357](https://github.com/CooooolFrog/dragonfly/issues/1357))
- release v2.0.3 ([#1360](https://github.com/CooooolFrog/dragonfly/issues/1360))
- add hack/gen-containerd-hosts.sh ([#1361](https://github.com/CooooolFrog/dragonfly/issues/1361))
- add check size workflows ([#1364](https://github.com/CooooolFrog/dragonfly/issues/1364))
- change dingtalk image ([#1954](https://github.com/CooooolFrog/dragonfly/issues/1954))
- build trainer binary and publish trainer docker image ([#2487](https://github.com/CooooolFrog/dragonfly/issues/2487))
- update submodule version
- release v2.0.4 ([#1425](https://github.com/CooooolFrog/dragonfly/issues/1425))
- create log dir ([#1947](https://github.com/CooooolFrog/dragonfly/issues/1947))
- optimize download log ([#1944](https://github.com/CooooolFrog/dragonfly/issues/1944))
- update codeql version ([#1428](https://github.com/CooooolFrog/dragonfly/issues/1428))
- exit when receive context done ([#1432](https://github.com/CooooolFrog/dragonfly/issues/1432))
- update docker compose ([#1431](https://github.com/CooooolFrog/dragonfly/issues/1431))
- upgrade kind node version ([#1433](https://github.com/CooooolFrog/dragonfly/issues/1433))
- update test/tools/no-content-length/main.go ([#1440](https://github.com/CooooolFrog/dragonfly/issues/1440))
- update helm charts version ([#1937](https://github.com/CooooolFrog/dragonfly/issues/1937))
- check header length before update ([#1445](https://github.com/CooooolFrog/dragonfly/issues/1445))
- dragonfly updates version to v2.0.5 ([#1498](https://github.com/CooooolFrog/dragonfly/issues/1498))
- add timestamp to stdout&stderr ([#1781](https://github.com/CooooolFrog/dragonfly/issues/1781))
- optimize source error log ([#1553](https://github.com/CooooolFrog/dragonfly/issues/1553))
- add priority to dfget man page ([#1917](https://github.com/CooooolFrog/dragonfly/issues/1917))
- add intel to ADOPTERS.md ([#1778](https://github.com/CooooolFrog/dragonfly/issues/1778))
- add e2e with nydus snapshotter ([#1860](https://github.com/CooooolFrog/dragonfly/issues/1860))
- update new manager ([#1597](https://github.com/CooooolFrog/dragonfly/issues/1597))
- add trainer to Makefile and shell ([#2488](https://github.com/CooooolFrog/dragonfly/issues/2488))
- add source error metrics ([#1560](https://github.com/CooooolFrog/dragonfly/issues/1560))
- fix macos build ([#1609](https://github.com/CooooolFrog/dragonfly/issues/1609))
- update debug info ([#1617](https://github.com/CooooolFrog/dragonfly/issues/1617))
- update api package verison ([#1893](https://github.com/CooooolFrog/dragonfly/issues/1893))
- optimize reregister ([#1888](https://github.com/CooooolFrog/dragonfly/issues/1888))
- releaser action disable cgo ([#1885](https://github.com/CooooolFrog/dragonfly/issues/1885))
- add Baidu to ADOPTERS.md ([#1884](https://github.com/CooooolFrog/dragonfly/issues/1884))
- release v2.0.8 ([#1877](https://github.com/CooooolFrog/dragonfly/issues/1877))
- workflows add tls e2e ([#1624](https://github.com/CooooolFrog/dragonfly/issues/1624))
- update tls e2e cert ([#1626](https://github.com/CooooolFrog/dragonfly/issues/1626))
- add Kuaishou to ADOPTERS.md ([#1866](https://github.com/CooooolFrog/dragonfly/issues/1866))
- release v2.0.6 version ([#1627](https://github.com/CooooolFrog/dragonfly/issues/1627))
- dependabot add github-actions ([#1629](https://github.com/CooooolFrog/dragonfly/issues/1629))
- gitignore add .run
- add disable seed peer action ([#1653](https://github.com/CooooolFrog/dragonfly/issues/1653))
- update helm-charts submodule version ([#1669](https://github.com/CooooolFrog/dragonfly/issues/1669))
- update dst peer log ([#1844](https://github.com/CooooolFrog/dragonfly/issues/1844))
- update download rpc check ([#1684](https://github.com/CooooolFrog/dragonfly/issues/1684))
- update e2e test ([#1839](https://github.com/CooooolFrog/dragonfly/issues/1839))
- remove unused code ([#1838](https://github.com/CooooolFrog/dragonfly/issues/1838))
- update api pkg ([#1700](https://github.com/CooooolFrog/dragonfly/issues/1700))
- change disk usage debug log format to decimal ([#1727](https://github.com/CooooolFrog/dragonfly/issues/1727))
- enable cache list metadata e2e ([#1829](https://github.com/CooooolFrog/dragonfly/issues/1829))
- daemon avoid alway open metadata files ([#1823](https://github.com/CooooolFrog/dragonfly/issues/1823))
- close out of use client grpc conn ([#1817](https://github.com/CooooolFrog/dragonfly/issues/1817))
- make lru cache safe ([#1737](https://github.com/CooooolFrog/dragonfly/issues/1737))
- change docker compose task ttl ([#1741](https://github.com/CooooolFrog/dragonfly/issues/1741))
- update console submodule ([#1748](https://github.com/CooooolFrog/dragonfly/issues/1748))
- update roundtrip log ([#1750](https://github.com/CooooolFrog/dragonfly/issues/1750))
- make SendMsg in doRecursiveDownload safe ([#1806](https://github.com/CooooolFrog/dragonfly/issues/1806))
- add list log in rpc download ([#1802](https://github.com/CooooolFrog/dragonfly/issues/1802))
- update console submodule ([#1755](https://github.com/CooooolFrog/dragonfly/issues/1755))
- update golang version to 1.19 ([#1760](https://github.com/CooooolFrog/dragonfly/issues/1760))
- check reuse file ([#1765](https://github.com/CooooolFrog/dragonfly/issues/1765))
- release v2.0.7 ([#1776](https://github.com/CooooolFrog/dragonfly/issues/1776))
- update helm-charts submodule
- upload nydus e2e logs to artifact ([#1909](https://github.com/CooooolFrog/dragonfly/issues/1909))
- **deps:** bump google.golang.org/api from 0.107.0 to 0.109.0 ([#2043](https://github.com/CooooolFrog/dragonfly/issues/2043))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.12 to 3.23.1 ([#2045](https://github.com/CooooolFrog/dragonfly/issues/2045))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.7 to 3.22.9 ([#1784](https://github.com/CooooolFrog/dragonfly/issues/1784))
- **deps:** bump github.com/gin-gonic/gin from 1.8.1 to 1.8.2 ([#1951](https://github.com/CooooolFrog/dragonfly/issues/1951))
- **deps:** bump gorm.io/driver/postgres from 1.3.7 to 1.3.8 ([#1503](https://github.com/CooooolFrog/dragonfly/issues/1503))
- **deps:** bump github.com/swaggo/gin-swagger from 1.5.1 to 1.5.2 ([#1604](https://github.com/CooooolFrog/dragonfly/issues/1604))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.2.0 to 2.4.0 ([#1787](https://github.com/CooooolFrog/dragonfly/issues/1787))
- **deps:** bump github.com/spf13/cobra from 1.5.0 to 1.6.1 ([#1796](https://github.com/CooooolFrog/dragonfly/issues/1796))
- **deps:** bump gorm.io/driver/postgres from 1.4.4 to 1.4.5 ([#1797](https://github.com/CooooolFrog/dragonfly/issues/1797))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.36.1 to 0.36.3 ([#1768](https://github.com/CooooolFrog/dragonfly/issues/1768))
- **deps:** bump github.com/swaggo/swag from 1.8.5 to 1.8.7 ([#1773](https://github.com/CooooolFrog/dragonfly/issues/1773))
- **deps:** bump k8s.io/component-base from 0.25.2 to 0.25.3 ([#1771](https://github.com/CooooolFrog/dragonfly/issues/1771))
- **deps:** bump github.com/swaggo/swag from 1.8.5 to 1.8.6 ([#1770](https://github.com/CooooolFrog/dragonfly/issues/1770))
- **deps:** bump github.com/casbin/casbin/v2 from 2.55.1 to 2.56.0 ([#1769](https://github.com/CooooolFrog/dragonfly/issues/1769))
- **deps:** bump go.opentelemetry.io/otel/trace from 1.10.0 to 1.11.0 ([#1767](https://github.com/CooooolFrog/dragonfly/issues/1767))
- **deps:** bump github.com/onsi/gomega from 1.22.1 to 1.23.0 ([#1798](https://github.com/CooooolFrog/dragonfly/issues/1798))
- **deps:** bump google.golang.org/api from 0.97.0 to 0.101.0 ([#1800](https://github.com/CooooolFrog/dragonfly/issues/1800))
- **deps:** bump gorm.io/driver/mysql from 1.4.1 to 1.4.3 ([#1799](https://github.com/CooooolFrog/dragonfly/issues/1799))
- **deps:** bump github.com/gammazero/deque from 0.2.0 to 0.2.1 ([#1810](https://github.com/CooooolFrog/dragonfly/issues/1810))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.36.1 to 0.36.4 ([#1811](https://github.com/CooooolFrog/dragonfly/issues/1811))
- **deps:** bump github.com/casbin/gorm-adapter/v3 from 3.5.0 to 3.11.0 ([#1745](https://github.com/CooooolFrog/dragonfly/issues/1745))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.8.7 to 3.11.0 ([#1746](https://github.com/CooooolFrog/dragonfly/issues/1746))
- **deps:** bump github.com/casbin/gorm-adapter/v3 from 3.11.0 to 3.12.1 ([#1786](https://github.com/CooooolFrog/dragonfly/issues/1786))
- **deps:** bump gorm.io/driver/postgres from 1.3.10 to 1.4.4 ([#1743](https://github.com/CooooolFrog/dragonfly/issues/1743))
- **deps:** bump google.golang.org/grpc from 1.49.0 to 1.50.0 ([#1742](https://github.com/CooooolFrog/dragonfly/issues/1742))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.9 to 3.22.10 ([#1812](https://github.com/CooooolFrog/dragonfly/issues/1812))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.10.0 to 1.11.1 ([#1813](https://github.com/CooooolFrog/dragonfly/issues/1813))
- **deps:** bump github.com/mdlayher/vsock from 1.1.1 to 1.2.0 ([#1834](https://github.com/CooooolFrog/dragonfly/issues/1834))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.95 to 1.44.114 ([#1725](https://github.com/CooooolFrog/dragonfly/issues/1725))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.3.0 to 1.3.3 ([#1722](https://github.com/CooooolFrog/dragonfly/issues/1722))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.9.0 to 1.10.0 ([#1720](https://github.com/CooooolFrog/dragonfly/issues/1720))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.34.0 to 0.36.1 ([#1719](https://github.com/CooooolFrog/dragonfly/issues/1719))
- **deps:** bump github.com/appleboy/gin-jwt/v2 from 2.8.0 to 2.9.0 ([#1718](https://github.com/CooooolFrog/dragonfly/issues/1718))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.1.6 to 2.2.0 ([#1705](https://github.com/CooooolFrog/dragonfly/issues/1705))
- **deps:** bump google.golang.org/api from 0.94.0 to 0.97.0 ([#1709](https://github.com/CooooolFrog/dragonfly/issues/1709))
- **deps:** bump k8s.io/component-base from 0.25.0 to 0.25.2 ([#1708](https://github.com/CooooolFrog/dragonfly/issues/1708))
- **deps:** bump gorm.io/gorm from 1.23.9 to 1.23.10 ([#1707](https://github.com/CooooolFrog/dragonfly/issues/1707))
- **deps:** bump github.com/casbin/casbin/v2 from 2.55.0 to 2.55.1 ([#1706](https://github.com/CooooolFrog/dragonfly/issues/1706))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.4.0 to 2.5.0 ([#1836](https://github.com/CooooolFrog/dragonfly/issues/1836))
- **deps:** bump gorm.io/gorm from 1.23.8 to 1.23.9 ([#1691](https://github.com/CooooolFrog/dragonfly/issues/1691))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.9.0 to 1.10.0 ([#1692](https://github.com/CooooolFrog/dragonfly/issues/1692))
- **deps:** bump gorm.io/driver/postgres from 1.3.9 to 1.3.10 ([#1690](https://github.com/CooooolFrog/dragonfly/issues/1690))
- **deps:** bump github.com/go-playground/validator/v10 from 10.11.0 to 10.11.1 ([#1689](https://github.com/CooooolFrog/dragonfly/issues/1689))
- **deps:** bump d7y.io/api from 1.1.4 to 1.1.6 ([#1688](https://github.com/CooooolFrog/dragonfly/issues/1688))
- **deps:** bump github.com/onsi/gomega from 1.23.0 to 1.24.1 ([#1832](https://github.com/CooooolFrog/dragonfly/issues/1832))
- **deps:** bump github.com/spf13/viper from 1.12.0 to 1.13.0 ([#1676](https://github.com/CooooolFrog/dragonfly/issues/1676))
- **deps:** bump github.com/casbin/casbin/v2 from 2.53.2 to 2.55.0 ([#1679](https://github.com/CooooolFrog/dragonfly/issues/1679))
- **deps:** bump k8s.io/component-base from 0.23.3 to 0.25.0 ([#1674](https://github.com/CooooolFrog/dragonfly/issues/1674))
- **deps:** bump github.com/swaggo/gin-swagger from 1.5.2 to 1.5.3 ([#1673](https://github.com/CooooolFrog/dragonfly/issues/1673))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.91 to 1.44.95 ([#1672](https://github.com/CooooolFrog/dragonfly/issues/1672))
- **deps:** bump k8s.io/component-base from 0.25.3 to 0.25.4 ([#1847](https://github.com/CooooolFrog/dragonfly/issues/1847))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.3.3 to 1.4.0 ([#1848](https://github.com/CooooolFrog/dragonfly/issues/1848))
- **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1650](https://github.com/CooooolFrog/dragonfly/issues/1650))
- **deps:** bump docker/login-action from 1 to 2 ([#1649](https://github.com/CooooolFrog/dragonfly/issues/1649))
- **deps:** bump docker/build-push-action from 2 to 3 ([#1648](https://github.com/CooooolFrog/dragonfly/issues/1648))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.44 to 1.44.91 ([#1647](https://github.com/CooooolFrog/dragonfly/issues/1647))
- **deps:** bump github.com/casbin/casbin/v2 from 2.52.2 to 2.53.2 ([#1644](https://github.com/CooooolFrog/dragonfly/issues/1644))
- **deps:** bump gorm.io/plugin/soft_delete from 1.1.0 to 1.2.0 ([#1643](https://github.com/CooooolFrog/dragonfly/issues/1643))
- **deps:** bump github.com/aliyun/aliyun-oss-go-sdk from 2.2.5+incompatible to 2.2.6+incompatible ([#1849](https://github.com/CooooolFrog/dragonfly/issues/1849))
- **deps:** bump go.uber.org/atomic from 1.9.0 to 1.10.0 ([#1639](https://github.com/CooooolFrog/dragonfly/issues/1639))
- **deps:** bump google.golang.org/api from 0.92.0 to 0.94.0 ([#1638](https://github.com/CooooolFrog/dragonfly/issues/1638))
- **deps:** bump github.com/onsi/gomega from 1.20.0 to 1.20.2 ([#1637](https://github.com/CooooolFrog/dragonfly/issues/1637))
- **deps:** bump github.com/swaggo/swag from 1.8.4 to 1.8.5 ([#1636](https://github.com/CooooolFrog/dragonfly/issues/1636))
- **deps:** bump go.uber.org/zap from 1.21.0 to 1.23.0 ([#1635](https://github.com/CooooolFrog/dragonfly/issues/1635))
- **deps:** bump docker/setup-buildx-action from 1 to 2 ([#1634](https://github.com/CooooolFrog/dragonfly/issues/1634))
- **deps:** bump actions/setup-go from 2 to 3 ([#1633](https://github.com/CooooolFrog/dragonfly/issues/1633))
- **deps:** bump actions/checkout from 2 to 3 ([#1631](https://github.com/CooooolFrog/dragonfly/issues/1631))
- **deps:** bump codecov/codecov-action from 1 to 3 ([#1630](https://github.com/CooooolFrog/dragonfly/issues/1630))
- **deps:** bump actions/upload-artifact from 2 to 3 ([#1632](https://github.com/CooooolFrog/dragonfly/issues/1632))
- **deps:** bump github.com/prometheus/client_golang from 1.13.0 to 1.14.0 ([#1851](https://github.com/CooooolFrog/dragonfly/issues/1851))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.121 to 1.44.143 ([#1853](https://github.com/CooooolFrog/dragonfly/issues/1853))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.5.0 to 2.5.1 ([#1871](https://github.com/CooooolFrog/dragonfly/issues/1871))
- **deps:** bump github.com/huaweicloud/huaweicloud-sdk-go-obs from 3.21.12+incompatible to 3.22.11+incompatible ([#1872](https://github.com/CooooolFrog/dragonfly/issues/1872))
- **deps:** bump github.com/go-sql-driver/mysql from 1.6.0 to 1.7.0 ([#1896](https://github.com/CooooolFrog/dragonfly/issues/1896))
- **deps:** bump github.com/swaggo/swag from 1.8.7 to 1.8.8 ([#1897](https://github.com/CooooolFrog/dragonfly/issues/1897))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.11.0 to 3.12.1 ([#1898](https://github.com/CooooolFrog/dragonfly/issues/1898))
- **deps:** bump gorm.io/driver/postgres from 1.3.8 to 1.3.9 ([#1608](https://github.com/CooooolFrog/dragonfly/issues/1608))
- **deps:** bump github.com/aliyun/aliyun-oss-go-sdk from 2.2.4+incompatible to 2.2.5+incompatible ([#1607](https://github.com/CooooolFrog/dragonfly/issues/1607))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.2.2 to 1.3.0 ([#1606](https://github.com/CooooolFrog/dragonfly/issues/1606))
- **deps:** bump github.com/gin-contrib/cors from 1.3.1 to 1.4.0 ([#1605](https://github.com/CooooolFrog/dragonfly/issues/1605))
- **deps:** bump github.com/casbin/casbin/v2 from 2.56.0 to 2.58.0 ([#1899](https://github.com/CooooolFrog/dragonfly/issues/1899))
- **deps:** bump go.uber.org/zap from 1.23.0 to 1.24.0 ([#1900](https://github.com/CooooolFrog/dragonfly/issues/1900))
- **deps:** bump github.com/casbin/casbin/v2 from 2.51.2 to 2.52.2 ([#1588](https://github.com/CooooolFrog/dragonfly/issues/1588))
- **deps:** bump github.com/swaggo/swag from 1.8.3 to 1.8.4 ([#1590](https://github.com/CooooolFrog/dragonfly/issues/1590))
- **deps:** bump k8s.io/apimachinery from 0.24.2 to 0.24.4 ([#1591](https://github.com/CooooolFrog/dragonfly/issues/1591))
- **deps:** bump gorm.io/driver/mysql from 1.3.4 to 1.3.6 ([#1567](https://github.com/CooooolFrog/dragonfly/issues/1567))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.33.0 to 0.34.0 ([#1566](https://github.com/CooooolFrog/dragonfly/issues/1566))
- **deps:** bump google.golang.org/api from 0.90.0 to 0.92.0 ([#1565](https://github.com/CooooolFrog/dragonfly/issues/1565))
- **deps:** bump github.com/prometheus/client_golang from 1.12.2 to 1.13.0 ([#1564](https://github.com/CooooolFrog/dragonfly/issues/1564))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.36.3 to 0.37.0 ([#1919](https://github.com/CooooolFrog/dragonfly/issues/1919))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.32.0 to 0.34.0 ([#1547](https://github.com/CooooolFrog/dragonfly/issues/1547))
- **deps:** bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 ([#1544](https://github.com/CooooolFrog/dragonfly/issues/1544))
- **deps:** bump github.com/jarcoal/httpmock from 1.0.8 to 1.2.0 ([#1542](https://github.com/CooooolFrog/dragonfly/issues/1542))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.8.0 to 1.9.0 ([#1541](https://github.com/CooooolFrog/dragonfly/issues/1541))
- **deps:** bump google.golang.org/grpc from 1.47.0 to 1.48.0 ([#1508](https://github.com/CooooolFrog/dragonfly/issues/1508))
- **deps:** bump github.com/casbin/casbin/v2 from 2.48.0 to 2.51.2 ([#1512](https://github.com/CooooolFrog/dragonfly/issues/1512))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.5 to 3.22.7 ([#1511](https://github.com/CooooolFrog/dragonfly/issues/1511))
- **deps:** bump google.golang.org/api from 0.86.0 to 0.90.0 ([#1510](https://github.com/CooooolFrog/dragonfly/issues/1510))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.32.0 to 0.33.0 ([#1509](https://github.com/CooooolFrog/dragonfly/issues/1509))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.12.1 to 3.12.2 ([#1920](https://github.com/CooooolFrog/dragonfly/issues/1920))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.7.0 to 1.8.0 ([#1506](https://github.com/CooooolFrog/dragonfly/issues/1506))
- **deps:** bump github.com/swaggo/gin-swagger from 1.5.0 to 1.5.1 ([#1505](https://github.com/CooooolFrog/dragonfly/issues/1505))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.8.6 to 3.8.7 ([#1502](https://github.com/CooooolFrog/dragonfly/issues/1502))
- **deps:** bump github.com/casbin/casbin/v2 from 2.58.0 to 2.60.0 ([#1921](https://github.com/CooooolFrog/dragonfly/issues/1921))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.11.1 to 1.11.2 ([#1922](https://github.com/CooooolFrog/dragonfly/issues/1922))
- **deps:** bump github.com/onsi/gomega from 1.24.1 to 1.24.2 ([#1931](https://github.com/CooooolFrog/dragonfly/issues/1931))
- **deps:** bump github.com/swaggo/swag from 1.8.8 to 1.8.9 ([#1932](https://github.com/CooooolFrog/dragonfly/issues/1932))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.10 to 3.22.11 ([#1935](https://github.com/CooooolFrog/dragonfly/issues/1935))
- **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1936](https://github.com/CooooolFrog/dragonfly/issues/1936))
- **deps:** bump github.com/mdlayher/vsock from 1.2.0 to 1.2.1 ([#2405](https://github.com/CooooolFrog/dragonfly/issues/2405))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.34.0 to 0.36.1 ([#1744](https://github.com/CooooolFrog/dragonfly/issues/1744))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.114 to 1.44.121 ([#1785](https://github.com/CooooolFrog/dragonfly/issues/1785))
- **deps:** bump go.opentelemetry.io/otel from 1.11.0 to 1.11.1 ([#1783](https://github.com/CooooolFrog/dragonfly/issues/1783))
- **deps:** bump google.golang.org/api from 0.101.0 to 0.105.0 ([#1952](https://github.com/CooooolFrog/dragonfly/issues/1952))
- **deps:** bump github.com/appleboy/gin-jwt/v2 from 2.9.0 to 2.9.1 ([#1949](https://github.com/CooooolFrog/dragonfly/issues/1949))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.167 to 1.44.171 ([#1958](https://github.com/CooooolFrog/dragonfly/issues/1958))
- **deps:** bump moul.io/zapgorm2 from 1.1.3 to 1.2.0 ([#1961](https://github.com/CooooolFrog/dragonfly/issues/1961))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.11 to 3.22.12 ([#1959](https://github.com/CooooolFrog/dragonfly/issues/1959))
- **deps:** bump gorm.io/driver/mysql from 1.4.4 to 1.4.5 ([#1962](https://github.com/CooooolFrog/dragonfly/issues/1962))
- **deps:** bump golang.org/x/time from 0.1.0 to 0.3.0 ([#1985](https://github.com/CooooolFrog/dragonfly/issues/1985))
- **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1986](https://github.com/CooooolFrog/dragonfly/issues/1986))
- **deps:** bump google.golang.org/api from 0.105.0 to 0.106.0 ([#1987](https://github.com/CooooolFrog/dragonfly/issues/1987))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.171 to 1.44.175 ([#1988](https://github.com/CooooolFrog/dragonfly/issues/1988))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.12.2 to 3.13.0 ([#1989](https://github.com/CooooolFrog/dragonfly/issues/1989))
- **deps:** bump gorm.io/driver/postgres from 1.4.5 to 1.4.6 ([#2002](https://github.com/CooooolFrog/dragonfly/issues/2002))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.6.1 to 2.7.0 ([#2003](https://github.com/CooooolFrog/dragonfly/issues/2003))
- **deps:** bump google.golang.org/api from 0.106.0 to 0.107.0 ([#2004](https://github.com/CooooolFrog/dragonfly/issues/2004))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.175 to 1.44.180 ([#2005](https://github.com/CooooolFrog/dragonfly/issues/2005))
- **deps:** bump gorm.io/gorm from 1.24.2 to 1.24.3 ([#2018](https://github.com/CooooolFrog/dragonfly/issues/2018))
- **deps:** bump github.com/spf13/viper from 1.13.0 to 1.15.0 ([#2019](https://github.com/CooooolFrog/dragonfly/issues/2019))
- **deps:** bump github.com/montanaflynn/stats from 0.6.6 to 0.7.0 ([#2020](https://github.com/CooooolFrog/dragonfly/issues/2020))
- **deps:** bump github.com/onsi/gomega from 1.24.2 to 1.25.0 ([#2021](https://github.com/CooooolFrog/dragonfly/issues/2021))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.180 to 1.44.184 ([#2022](https://github.com/CooooolFrog/dragonfly/issues/2022))
- **deps:** bump github.com/onsi/gomega from 1.25.0 to 1.26.0 ([#2024](https://github.com/CooooolFrog/dragonfly/issues/2024))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.7.0 to 2.7.1 ([#2028](https://github.com/CooooolFrog/dragonfly/issues/2028))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.184 to 1.44.189 ([#2029](https://github.com/CooooolFrog/dragonfly/issues/2029))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.11.2 to 1.12.0 ([#2030](https://github.com/CooooolFrog/dragonfly/issues/2030))
- **deps:** bump gorm.io/gorm from 1.24.3 to 1.24.5 ([#2042](https://github.com/CooooolFrog/dragonfly/issues/2042))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.143 to 1.44.167 ([#1948](https://github.com/CooooolFrog/dragonfly/issues/1948))
- **deps:** bump github.com/jarcoal/httpmock from 1.2.0 to 1.3.0 ([#2044](https://github.com/CooooolFrog/dragonfly/issues/2044))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.39.0 to 0.41.1 ([#2352](https://github.com/CooooolFrog/dragonfly/issues/2352))
- **deps:** bump docker/build-push-action from 3 to 4 ([#2047](https://github.com/CooooolFrog/dragonfly/issues/2047))
- **deps:** bump google.golang.org/grpc from 1.52.0 to 1.52.3 ([#2046](https://github.com/CooooolFrog/dragonfly/issues/2046))
- **deps:** bump github.com/looplab/fsm from 1.0.0 to 1.0.1 ([#2073](https://github.com/CooooolFrog/dragonfly/issues/2073))
- **deps:** bump go.opentelemetry.io/otel from 1.12.0 to 1.13.0 ([#2074](https://github.com/CooooolFrog/dragonfly/issues/2074))
- **deps:** bump github.com/casbin/casbin/v2 from 2.60.0 to 2.61.1 ([#2075](https://github.com/CooooolFrog/dragonfly/issues/2075))
- **deps:** bump github.com/go-playground/validator/v10 from 10.11.1 to 10.11.2 ([#2077](https://github.com/CooooolFrog/dragonfly/issues/2077))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.12.0 to 1.13.0 ([#2093](https://github.com/CooooolFrog/dragonfly/issues/2093))
- **deps:** bump golang.org/x/oauth2 from 0.4.0 to 0.5.0 ([#2094](https://github.com/CooooolFrog/dragonfly/issues/2094))
- **deps:** bump gorm.io/driver/mysql from 1.4.5 to 1.4.7 ([#2096](https://github.com/CooooolFrog/dragonfly/issues/2096))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.4.0 to 1.5.0 ([#2097](https://github.com/CooooolFrog/dragonfly/issues/2097))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.37.0 to 0.39.0 ([#2120](https://github.com/CooooolFrog/dragonfly/issues/2120))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.189 to 1.44.209 ([#2122](https://github.com/CooooolFrog/dragonfly/issues/2122))
- **deps:** bump github.com/casbin/casbin/v2 from 2.61.1 to 2.64.0 ([#2123](https://github.com/CooooolFrog/dragonfly/issues/2123))
- **deps:** bump github.com/stretchr/testify from 1.8.1 to 1.8.2 ([#2140](https://github.com/CooooolFrog/dragonfly/issues/2140))
- **deps:** bump gorm.io/driver/postgres from 1.4.6 to 1.4.8 ([#2142](https://github.com/CooooolFrog/dragonfly/issues/2142))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.13.0 to 1.14.0 ([#2144](https://github.com/CooooolFrog/dragonfly/issues/2144))
- **deps:** bump gorm.io/gorm from 1.24.5 to 1.24.6 ([#2143](https://github.com/CooooolFrog/dragonfly/issues/2143))
- **deps:** bump golang.org/x/crypto from 0.6.0 to 0.7.0 ([#2163](https://github.com/CooooolFrog/dragonfly/issues/2163))
- **deps:** bump github.com/casbin/casbin/v2 from 2.64.0 to 2.65.2 ([#2164](https://github.com/CooooolFrog/dragonfly/issues/2164))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.1 to 3.23.2 ([#2166](https://github.com/CooooolFrog/dragonfly/issues/2166))
- **deps:** bump moul.io/zapgorm2 from 1.2.0 to 1.3.0 ([#2167](https://github.com/CooooolFrog/dragonfly/issues/2167))
- **deps:** bump google.golang.org/protobuf from 1.29.0 to 1.29.1 ([#2195](https://github.com/CooooolFrog/dragonfly/issues/2195))
- **deps:** bump github.com/swaggo/swag from 1.8.9 to 1.8.10 ([#2197](https://github.com/CooooolFrog/dragonfly/issues/2197))
- **deps:** bump golang.org/x/oauth2 from 0.5.0 to 0.6.0 ([#2198](https://github.com/CooooolFrog/dragonfly/issues/2198))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.13.0 to 3.13.1 ([#2199](https://github.com/CooooolFrog/dragonfly/issues/2199))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.209 to 1.44.224 ([#2200](https://github.com/CooooolFrog/dragonfly/issues/2200))
- **deps:** bump google.golang.org/api from 0.109.0 to 0.114.0 ([#2201](https://github.com/CooooolFrog/dragonfly/issues/2201))
- **deps:** bump actions/setup-go from 3 to 4 ([#2202](https://github.com/CooooolFrog/dragonfly/issues/2202))
- **deps:** bump gorm.io/driver/postgres from 1.4.8 to 1.5.0 ([#2217](https://github.com/CooooolFrog/dragonfly/issues/2217))
- **deps:** bump github.com/aliyun/aliyun-oss-go-sdk from 2.2.6+incompatible to 2.2.7+incompatible ([#2218](https://github.com/CooooolFrog/dragonfly/issues/2218))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.37.0 to 0.40.0 ([#2219](https://github.com/CooooolFrog/dragonfly/issues/2219))
- **deps:** bump github.com/go-playground/validator/v10 from 10.11.2 to 10.12.0 ([#2220](https://github.com/CooooolFrog/dragonfly/issues/2220))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.224 to 1.44.229 ([#2221](https://github.com/CooooolFrog/dragonfly/issues/2221))
- **deps:** bump github.com/casbin/gorm-adapter/v3 from 3.5.0 to 3.15.0 ([#2237](https://github.com/CooooolFrog/dragonfly/issues/2237))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.2 to 3.23.3 ([#2239](https://github.com/CooooolFrog/dragonfly/issues/2239))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.229 to 1.44.234 ([#2240](https://github.com/CooooolFrog/dragonfly/issues/2240))
- **deps:** bump github.com/casbin/casbin/v2 from 2.65.2 to 2.66.1 ([#2238](https://github.com/CooooolFrog/dragonfly/issues/2238))
- **deps:** bump github.com/gin-gonic/gin from 1.8.2 to 1.9.0 ([#2241](https://github.com/CooooolFrog/dragonfly/issues/2241))
- **deps:** bump github.com/swaggo/gin-swagger from 1.5.3 to 1.6.0 ([#2256](https://github.com/CooooolFrog/dragonfly/issues/2256))
- **deps:** bump github.com/casbin/casbin/v2 from 2.66.1 to 2.66.3 ([#2260](https://github.com/CooooolFrog/dragonfly/issues/2260))
- **deps:** bump gorm.io/gorm from 1.24.7-0.20230306060331-85eaf9eeda11 to 1.25.0 ([#2277](https://github.com/CooooolFrog/dragonfly/issues/2277))
- **deps:** bump d7y.io/api from 1.8.6 to 1.8.7 ([#2278](https://github.com/CooooolFrog/dragonfly/issues/2278))
- **deps:** bump gorm.io/plugin/soft_delete from 1.2.0 to 1.2.1 ([#2279](https://github.com/CooooolFrog/dragonfly/issues/2279))
- **deps:** bump github.com/grpc-ecosystem/go-grpc-middleware from 1.3.0 to 1.4.0 ([#2280](https://github.com/CooooolFrog/dragonfly/issues/2280))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.234 to 1.44.244 ([#2281](https://github.com/CooooolFrog/dragonfly/issues/2281))
- **deps:** bump golang.org/x/sys from 0.6.0 to 0.7.0 ([#2297](https://github.com/CooooolFrog/dragonfly/issues/2297))
- **deps:** bump github.com/casbin/gorm-adapter/v3 from 3.5.0 to 3.16.1 ([#2298](https://github.com/CooooolFrog/dragonfly/issues/2298))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.5.0 to 1.7.0 ([#2300](https://github.com/CooooolFrog/dragonfly/issues/2300))
- **deps:** bump github.com/prometheus/client_golang from 1.14.0 to 1.15.0 ([#2299](https://github.com/CooooolFrog/dragonfly/issues/2299))
- **deps:** bump golang.org/x/crypto from 0.7.0 to 0.8.0 ([#2311](https://github.com/CooooolFrog/dragonfly/issues/2311))
- **deps:** bump golang.org/x/oauth2 from 0.6.0 to 0.7.0 ([#2310](https://github.com/CooooolFrog/dragonfly/issues/2310))
- **deps:** bump gorm.io/driver/mysql from 1.4.7 to 1.5.0 ([#2312](https://github.com/CooooolFrog/dragonfly/issues/2312))
- **deps:** bump go.opentelemetry.io/otel/trace from 1.14.0 to 1.15.0 ([#2313](https://github.com/CooooolFrog/dragonfly/issues/2313))
- **deps:** bump github.com/swaggo/swag from 1.8.12 to 1.16.1 ([#2331](https://github.com/CooooolFrog/dragonfly/issues/2331))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.7.1 to 2.9.4 ([#2332](https://github.com/CooooolFrog/dragonfly/issues/2332))
- **deps:** bump github.com/go-sql-driver/mysql from 1.7.0 to 1.7.1 ([#2333](https://github.com/CooooolFrog/dragonfly/issues/2333))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.244 to 1.44.258 ([#2334](https://github.com/CooooolFrog/dragonfly/issues/2334))
- **deps:** bump go.opentelemetry.io/otel/trace from 1.15.0 to 1.15.1 ([#2335](https://github.com/CooooolFrog/dragonfly/issues/2335))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.36.4 to 0.37.0 ([#1950](https://github.com/CooooolFrog/dragonfly/issues/1950))
- **deps:** bump github.com/casbin/gorm-adapter/v3 from 3.16.1 to 3.17.0 ([#2353](https://github.com/CooooolFrog/dragonfly/issues/2353))
- **deps:** bump golang.org/x/crypto from 0.8.0 to 0.9.0 ([#2355](https://github.com/CooooolFrog/dragonfly/issues/2355))
- **deps:** bump github.com/huaweicloud/huaweicloud-sdk-go-obs from 3.22.11+incompatible to 3.23.4+incompatible ([#2373](https://github.com/CooooolFrog/dragonfly/issues/2373))
- **deps:** bump golang.org/x/oauth2 from 0.7.0 to 0.8.0 ([#2372](https://github.com/CooooolFrog/dragonfly/issues/2372))
- **deps:** bump gorm.io/driver/mysql from 1.5.0 to 1.5.1 ([#2374](https://github.com/CooooolFrog/dragonfly/issues/2374))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.14.0 to 1.15.1 ([#2376](https://github.com/CooooolFrog/dragonfly/issues/2376))
- **deps:** bump go.uber.org/atomic from 1.10.0 to 1.11.0 ([#2404](https://github.com/CooooolFrog/dragonfly/issues/2404))
- **deps:** bump github.com/stretchr/testify from 1.8.2 to 1.8.3 ([#2408](https://github.com/CooooolFrog/dragonfly/issues/2408))
- **deps:** bump k8s.io/component-base from 0.25.4 to 0.26.0 ([#1934](https://github.com/CooooolFrog/dragonfly/issues/1934))
- **deps:** bump github.com/montanaflynn/stats from 0.7.0 to 0.7.1 ([#2407](https://github.com/CooooolFrog/dragonfly/issues/2407))
- **deps:** bump github.com/gin-gonic/gin from 1.9.0 to 1.9.1 ([#2419](https://github.com/CooooolFrog/dragonfly/issues/2419))
- **deps:** bump k8s.io/component-base from 0.26.0 to 0.27.2 ([#2432](https://github.com/CooooolFrog/dragonfly/issues/2432))
- **deps:** bump google.golang.org/grpc from 1.56.0-dev to 1.57.0-dev ([#2433](https://github.com/CooooolFrog/dragonfly/issues/2433))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.3 to 3.23.5 ([#2434](https://github.com/CooooolFrog/dragonfly/issues/2434))
- **deps:** bump golang.org/x/crypto from 0.9.0 to 0.10.0 ([#2474](https://github.com/CooooolFrog/dragonfly/issues/2474))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.41.1 to 0.42.0 ([#2475](https://github.com/CooooolFrog/dragonfly/issues/2475))
- **deps:** bump github.com/stretchr/testify from 1.8.3 to 1.8.4 ([#2476](https://github.com/CooooolFrog/dragonfly/issues/2476))
- **deps:** bump google.golang.org/api from 0.114.0 to 0.128.0 ([#2478](https://github.com/CooooolFrog/dragonfly/issues/2478))
- **deps:** bump github.com/prometheus/client_golang from 1.15.0 to 1.16.0 ([#2481](https://github.com/CooooolFrog/dragonfly/issues/2481))
- **deps:** bump github.com/go-playground/validator/v10 from 10.14.0 to 10.14.1 ([#2483](https://github.com/CooooolFrog/dragonfly/issues/2483))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.9.4 to 2.11.0 ([#2484](https://github.com/CooooolFrog/dragonfly/issues/2484))

### Daemon
- add add additional peer id for some logs ([#205](https://github.com/CooooolFrog/dragonfly/issues/205))
- create output parent directory if not exists ([#188](https://github.com/CooooolFrog/dragonfly/issues/188))
- update default timeout and add context for downloading piece ([#190](https://github.com/CooooolFrog/dragonfly/issues/190))
- record failed code when unfinished and event for scheduler ([#176](https://github.com/CooooolFrog/dragonfly/issues/176))

### Docs
- metrics configuration ([#816](https://github.com/CooooolFrog/dragonfly/issues/816))
- change dingtalk link
- add OpenSSF badge to README.md ([#2138](https://github.com/CooooolFrog/dragonfly/issues/2138))
- add public cloud providers Adopters.md ([#2137](https://github.com/CooooolFrog/dragonfly/issues/2137))
- change introduction in readem ([#2017](https://github.com/CooooolFrog/dragonfly/issues/2017))
- fix manager swag error ([#1982](https://github.com/CooooolFrog/dragonfly/issues/1982))
- optimize Community description in README.md ([#2255](https://github.com/CooooolFrog/dragonfly/issues/2255))
- add adopters.md ([#1759](https://github.com/CooooolFrog/dragonfly/issues/1759))
- add daemon-socket for daemon docs ([#1522](https://github.com/CooooolFrog/dragonfly/issues/1522))
- update CHANGELOG
- update CHANGELOG
- update readme system features ([#1359](https://github.com/CooooolFrog/dragonfly/issues/1359))
- readme typo
- readme add seed peer ([#1349](https://github.com/CooooolFrog/dragonfly/issues/1349))
- move document from /docs to d7y.io ([#1229](https://github.com/CooooolFrog/dragonfly/issues/1229))
- add slack and google groups ([#1203](https://github.com/CooooolFrog/dragonfly/issues/1203))
- add plugin builder ([#1101](https://github.com/CooooolFrog/dragonfly/issues/1101))
- add metrics document ([#1075](https://github.com/CooooolFrog/dragonfly/issues/1075))
- add containerd private registry configuration ([#1074](https://github.com/CooooolFrog/dragonfly/issues/1074))
- manager apis ([#814](https://github.com/CooooolFrog/dragonfly/issues/814))
- add docs about preheat console ([#1072](https://github.com/CooooolFrog/dragonfly/issues/1072))
- manager installation ([#1063](https://github.com/CooooolFrog/dragonfly/issues/1063))
- update plugin doc ([#951](https://github.com/CooooolFrog/dragonfly/issues/951))
- update plugin docs ([#921](https://github.com/CooooolFrog/dragonfly/issues/921))
- dir path ([#904](https://github.com/CooooolFrog/dragonfly/issues/904))
- add plugin guide ([#875](https://github.com/CooooolFrog/dragonfly/issues/875))
- keep alive ([#868](https://github.com/CooooolFrog/dragonfly/issues/868))
- add CODE_OF_CONDUCT.md ([#163](https://github.com/CooooolFrog/dragonfly/issues/163))
- update quick-start.md format ([#850](https://github.com/CooooolFrog/dragonfly/issues/850))
- add Volcano Engine to ADOPTERS.md ([#2169](https://github.com/CooooolFrog/dragonfly/issues/2169))
- add containerd private registry configuration ([#1073](https://github.com/CooooolFrog/dragonfly/issues/1073))
- add CHANGELOG.md
- manager api ([#774](https://github.com/CooooolFrog/dragonfly/issues/774))
- add maxConcurrency comment ([#755](https://github.com/CooooolFrog/dragonfly/issues/755))
- add troubleshooting guide ([#752](https://github.com/CooooolFrog/dragonfly/issues/752))
- update v0.1.0-beta changelog ([#387](https://github.com/CooooolFrog/dragonfly/issues/387))
- Add dfget man page ([#388](https://github.com/CooooolFrog/dragonfly/issues/388))
- add load limit ([#745](https://github.com/CooooolFrog/dragonfly/issues/745))
- update kubernetes docs ([#714](https://github.com/CooooolFrog/dragonfly/issues/714))
- add apis and preheat ([#712](https://github.com/CooooolFrog/dragonfly/issues/712))
- update kubernetes docs ([#705](https://github.com/CooooolFrog/dragonfly/issues/705))
- scheduler config ([#698](https://github.com/CooooolFrog/dragonfly/issues/698))
- update kubernetes docs ([#696](https://github.com/CooooolFrog/dragonfly/issues/696))
- scheduler config ([#654](https://github.com/CooooolFrog/dragonfly/issues/654))
- maintainers ([#636](https://github.com/CooooolFrog/dragonfly/issues/636))
- test guide link ([#635](https://github.com/CooooolFrog/dragonfly/issues/635))
- add manager preview ([#634](https://github.com/CooooolFrog/dragonfly/issues/634))
- install ([#628](https://github.com/CooooolFrog/dragonfly/issues/628))
- update document ([#625](https://github.com/CooooolFrog/dragonfly/issues/625))
- update docs/zh-CN/config/dfget.yaml ([#623](https://github.com/CooooolFrog/dragonfly/issues/623))
- Update documents ([#595](https://github.com/CooooolFrog/dragonfly/issues/595))
- update runtime guide in helm deploy ([#612](https://github.com/CooooolFrog/dragonfly/issues/612))
- rbac swagger comment
- **en:** upgrade docs ([#673](https://github.com/CooooolFrog/dragonfly/issues/673))
- **runtime:** upgrade containerd runtime ([#748](https://github.com/CooooolFrog/dragonfly/issues/748))
- **zh:** add zh docs ([#777](https://github.com/CooooolFrog/dragonfly/issues/777))
- **zh-CN:** refactor machine translation ([#783](https://github.com/CooooolFrog/dragonfly/issues/783))

### Feat
- remove cdn examples in grpc
- add network topology to daemon ([#2489](https://github.com/CooooolFrog/dragonfly/issues/2489))
- implement probe interface in client daemon ([#2473](https://github.com/CooooolFrog/dragonfly/issues/2473))
- add trainer cmd and trainer service ([#2479](https://github.com/CooooolFrog/dragonfly/issues/2479))
- add context key to SyncProbes ([#2485](https://github.com/CooooolFrog/dragonfly/issues/2485))
- daemon store file exclusive ([#2465](https://github.com/CooooolFrog/dragonfly/issues/2465))
- update lint action to solve cache conflict ([#2472](https://github.com/CooooolFrog/dragonfly/issues/2472))
- update kind action ([#2470](https://github.com/CooooolFrog/dragonfly/issues/2470))
- move probe interval from scheduler config to client config ([#2462](https://github.com/CooooolFrog/dragonfly/issues/2462))
- add Access-Control-Expose-Headers to headers ([#2467](https://github.com/CooooolFrog/dragonfly/issues/2467))
- support breakpoint resume for running tasks ([#2457](https://github.com/CooooolFrog/dragonfly/issues/2457))
- implement SyncProbes api in scheduler grpc service ([#2449](https://github.com/CooooolFrog/dragonfly/issues/2449))
- optimize dfpath format ([#2460](https://github.com/CooooolFrog/dragonfly/issues/2460))
- enable configuration of some directory modes for dfdaemon ([#2340](https://github.com/CooooolFrog/dragonfly/issues/2340))
- remove dirty file
- optimize announcer in scheduler and client ([#2445](https://github.com/CooooolFrog/dragonfly/issues/2445))
- change DefaultProbeInterval to 20 minute ([#2440](https://github.com/CooooolFrog/dragonfly/issues/2440))
- remove useless fields in network topology ([#2439](https://github.com/CooooolFrog/dragonfly/issues/2439))
- add storage to trainer ([#2431](https://github.com/CooooolFrog/dragonfly/issues/2431))
- support to collect and snapshot in network topology ([#2429](https://github.com/CooooolFrog/dragonfly/issues/2429))
- add ip to uk_scheduler index and uk_seed_peer index in manager ([#2426](https://github.com/CooooolFrog/dragonfly/issues/2426))
- change Dequeue to private func ([#2420](https://github.com/CooooolFrog/dragonfly/issues/2420))
- specify the version of golangci-lint as v1.52.2 ([#2421](https://github.com/CooooolFrog/dragonfly/issues/2421))
- remove redis Pipelined in network topology ([#2416](https://github.com/CooooolFrog/dragonfly/issues/2416))
- optimize network topology comment ([#2415](https://github.com/CooooolFrog/dragonfly/issues/2415))
- add ProbedAt to network topology ([#2413](https://github.com/CooooolFrog/dragonfly/issues/2413))
- implement Enqueue and AverageRTT in probes.go ([#2393](https://github.com/CooooolFrog/dragonfly/issues/2393))
- handle context in triggerSeedPeerTask ([#2392](https://github.com/CooooolFrog/dragonfly/issues/2392))
- optimize field name of ProbeConfig ([#2391](https://github.com/CooooolFrog/dragonfly/issues/2391))
- scheduler supports to disable redis ([#2389](https://github.com/CooooolFrog/dragonfly/issues/2389))
- add Reverse function to slice ([#2381](https://github.com/CooooolFrog/dragonfly/issues/2381))
- move redis key to pkg/redis package ([#2378](https://github.com/CooooolFrog/dragonfly/issues/2378))
- add network topology package ([#2364](https://github.com/CooooolFrog/dragonfly/issues/2364))
- add announceToTrainer in scheduler ([#2371](https://github.com/CooooolFrog/dragonfly/issues/2371))
- hide sensitive information in log ([#2369](https://github.com/CooooolFrog/dragonfly/issues/2369))
- replace net dial with grpc health check in client ([#2361](https://github.com/CooooolFrog/dragonfly/issues/2361))
- remove traffic_type in DownloadPeerDuration metric ([#2357](https://github.com/CooooolFrog/dragonfly/issues/2357))
- add traffic type of peer task download duration ([#2349](https://github.com/CooooolFrog/dragonfly/issues/2349))
- change DefaultServerPort to 9090 in trainer ([#2348](https://github.com/CooooolFrog/dragonfly/issues/2348))
- remove deprecated field in manager and scheduler ([#2345](https://github.com/CooooolFrog/dragonfly/issues/2345))
- add database config and move redis to it ([#2338](https://github.com/CooooolFrog/dragonfly/issues/2338))
- remove compatibility logic for manager config testing ([#2342](https://github.com/CooooolFrog/dragonfly/issues/2342))
- optimize job new in internal ([#2341](https://github.com/CooooolFrog/dragonfly/issues/2341))
- remove log of configuration ([#2322](https://github.com/CooooolFrog/dragonfly/issues/2322))
- rename createRecord to createDownloadRecord ([#2306](https://github.com/CooooolFrog/dragonfly/issues/2306))
- add CORS middleware to manager ([#2304](https://github.com/CooooolFrog/dragonfly/issues/2304))
-  add metrics for trainer ([#2293](https://github.com/CooooolFrog/dragonfly/issues/2293))
- add Access-Control-Allow-Credentials to rest api ([#2302](https://github.com/CooooolFrog/dragonfly/issues/2302))
- remove SyncNetworkTopology API ([#2296](https://github.com/CooooolFrog/dragonfly/issues/2296))
- move redis package to pkg dir ([#2294](https://github.com/CooooolFrog/dragonfly/issues/2294))
- optimize model rest api in manager ([#2291](https://github.com/CooooolFrog/dragonfly/issues/2291))
- add model operation api ([#2276](https://github.com/CooooolFrog/dragonfly/issues/2276))
- add network topology storage interface ([#2286](https://github.com/CooooolFrog/dragonfly/issues/2286))
- add cluster api in manager ([#2288](https://github.com/CooooolFrog/dragonfly/issues/2288))
- add network topology and probes storage structs ([#2254](https://github.com/CooooolFrog/dragonfly/issues/2254))
- remove security domain ([#2285](https://github.com/CooooolFrog/dragonfly/issues/2285))
- rename trainer config package to config ([#2283](https://github.com/CooooolFrog/dragonfly/issues/2283))
- add multi-arch container images to workflow ([#2270](https://github.com/CooooolFrog/dragonfly/issues/2270))
- rename Record to Download in storage ([#2253](https://github.com/CooooolFrog/dragonfly/issues/2253))
- local dynconfig notifies data in client ([#2264](https://github.com/CooooolFrog/dragonfly/issues/2264))
- update resource director ([#2243](https://github.com/CooooolFrog/dragonfly/issues/2243))
- add CreatedAt function ([#2244](https://github.com/CooooolFrog/dragonfly/issues/2244))
- add trainer configuration ([#2216](https://github.com/CooooolFrog/dragonfly/issues/2216))
- update d7y.io/api package and change cpu percent validation ([#2236](https://github.com/CooooolFrog/dragonfly/issues/2236))
- add authinfo injector ([#2149](https://github.com/CooooolFrog/dragonfly/issues/2149))
- when the cache is missing, change the error log to a warning log ([#2235](https://github.com/CooooolFrog/dragonfly/issues/2235))
-  if the scheduler feature is not in feature flags, then it will stop providing the featrue ([#2234](https://github.com/CooooolFrog/dragonfly/issues/2234))
- add train interval and trainer addresses ([#2223](https://github.com/CooooolFrog/dragonfly/issues/2223))
- add logger.CoreLogger to searcher plugin ([#2232](https://github.com/CooooolFrog/dragonfly/issues/2232))
- add log to searcher plugin ([#2231](https://github.com/CooooolFrog/dragonfly/issues/2231))
- add probes struct ([#2190](https://github.com/CooooolFrog/dragonfly/issues/2190))
- add trainer config in scheduler ([#2214](https://github.com/CooooolFrog/dragonfly/issues/2214))
- add tfserving service to rpc package ([#2210](https://github.com/CooooolFrog/dragonfly/issues/2210))
- add trainer service to rpc package ([#2209](https://github.com/CooooolFrog/dragonfly/issues/2209))
- rename security client file name ([#2208](https://github.com/CooooolFrog/dragonfly/issues/2208))
- add CreateModel func to manager grpc client ([#2207](https://github.com/CooooolFrog/dragonfly/issues/2207))
- rename SecurityService to Security ([#2206](https://github.com/CooooolFrog/dragonfly/issues/2206))
- rename HostName to Hostname ([#2205](https://github.com/CooooolFrog/dragonfly/issues/2205))
- remove model migration ([#2204](https://github.com/CooooolFrog/dragonfly/issues/2204))
- change default value of dynconfig cache ([#2203](https://github.com/CooooolFrog/dragonfly/issues/2203))
- add index uk_model to model table ([#2196](https://github.com/CooooolFrog/dragonfly/issues/2196))
- remove model api ([#2194](https://github.com/CooooolFrog/dragonfly/issues/2194))
- add inference model table in database ([#2192](https://github.com/CooooolFrog/dragonfly/issues/2192))
- rename manager/model to manager/models ([#2191](https://github.com/CooooolFrog/dragonfly/issues/2191))
- add advertisePort to manager ([#2189](https://github.com/CooooolFrog/dragonfly/issues/2189))
- add advertise port ([#2156](https://github.com/CooooolFrog/dragonfly/issues/2156))
- add error log to database in manager ([#2172](https://github.com/CooooolFrog/dragonfly/issues/2172))
- add auth config to manager ([#2161](https://github.com/CooooolFrog/dragonfly/issues/2161))
- add metrics to service v2 ([#2153](https://github.com/CooooolFrog/dragonfly/issues/2153))
- add SearchSchedulerClusterCount metric to manager ([#2152](https://github.com/CooooolFrog/dragonfly/issues/2152))
- implement announce peer ([#2150](https://github.com/CooooolFrog/dragonfly/issues/2150))
- add handleRegisterSeedPeerRequest to service v2 in scheduler ([#2148](https://github.com/CooooolFrog/dragonfly/issues/2148))
- add handleRegisterSeedPeerRequest to AnnouncePeer in service v2 ([#2147](https://github.com/CooooolFrog/dragonfly/issues/2147))
- change ScheduleCandidateParentsForNormalPeer implement ([#2133](https://github.com/CooooolFrog/dragonfly/issues/2133))
- enhance daemon health check ([#2130](https://github.com/CooooolFrog/dragonfly/issues/2130))
- implement v2 version of scheduler service ([#2125](https://github.com/CooooolFrog/dragonfly/issues/2125))
- update golang version to 1.20.1 ([#2117](https://github.com/CooooolFrog/dragonfly/issues/2117))
- correct grpc error code and implement StatPeer and LeavePeer ([#2115](https://github.com/CooooolFrog/dragonfly/issues/2115))
- add SyncNetworkTopology and SyncProbes to scheduler client ([#2114](https://github.com/CooooolFrog/dragonfly/issues/2114))
- add CIDR affinity to searcher ([#2111](https://github.com/CooooolFrog/dragonfly/issues/2111))
- remove Scopes and SecurityGroup in seed peer cluster ([#2110](https://github.com/CooooolFrog/dragonfly/issues/2110))
- dynconfig resolves addresses with host ([#2109](https://github.com/CooooolFrog/dragonfly/issues/2109))
- enable oss client download object concurrently. ([#2105](https://github.com/CooooolFrog/dragonfly/issues/2105))
- support reload scheduler addresses for local Dynconfig ([#2091](https://github.com/CooooolFrog/dragonfly/issues/2091))
- oss client supports STS access (set security token in header) ([#2103](https://github.com/CooooolFrog/dragonfly/issues/2103))
- don't GC task if expire time is 0 ([#2102](https://github.com/CooooolFrog/dragonfly/issues/2102))
- avoid checking dir existence before MkdirAll ([#2090](https://github.com/CooooolFrog/dragonfly/issues/2090))
- add host ttl to scheduler ([#2089](https://github.com/CooooolFrog/dragonfly/issues/2089))
- rename scheduler package to scheduling ([#2087](https://github.com/CooooolFrog/dragonfly/issues/2087))
- use v2 version of host id and add Addrs func to seed peer ([#2086](https://github.com/CooooolFrog/dragonfly/issues/2086))
- add networkTopology configuration to scheduler ([#2070](https://github.com/CooooolFrog/dragonfly/issues/2070))
- remove training configuration in scheduler ([#2081](https://github.com/CooooolFrog/dragonfly/issues/2081))
- change piece size to length ([#2079](https://github.com/CooooolFrog/dragonfly/issues/2079))
- set gorm log level ([#2063](https://github.com/CooooolFrog/dragonfly/issues/2063))
- change PeerCountLimitForTask to 1000 ([#2059](https://github.com/CooooolFrog/dragonfly/issues/2059))
- add v2 version of the idgen ([#2056](https://github.com/CooooolFrog/dragonfly/issues/2056))
- update task type from v1 to v2 ([#2053](https://github.com/CooooolFrog/dragonfly/issues/2053))
- add AnnouncePeers to task in resource ([#2051](https://github.com/CooooolFrog/dragonfly/issues/2051))
- add v2 version of dfdaemon client ([#2050](https://github.com/CooooolFrog/dragonfly/issues/2050))
- add DownloadTask to seed peer resource ([#2048](https://github.com/CooooolFrog/dragonfly/issues/2048))
- init AnnouncePeerStream of the peer ([#2040](https://github.com/CooooolFrog/dragonfly/issues/2040))
- update dingtalk qrcode ([#2016](https://github.com/CooooolFrog/dragonfly/issues/2016))
- update helm charts ([#2015](https://github.com/CooooolFrog/dragonfly/issues/2015))
- add directed graph to pkg ([#2014](https://github.com/CooooolFrog/dragonfly/issues/2014))
- change peer's piece type in resource ([#2012](https://github.com/CooooolFrog/dragonfly/issues/2012))
- support source client option ([#2008](https://github.com/CooooolFrog/dragonfly/issues/2008))
- change ok to loaded in loading func ([#2010](https://github.com/CooooolFrog/dragonfly/issues/2010))
- remove NetTopology from scheduler and manager ([#2007](https://github.com/CooooolFrog/dragonfly/issues/2007))
- add v2 verison of the grpc to scheduler ([#1999](https://github.com/CooooolFrog/dragonfly/issues/1999))
- add manager v2 api ([#1990](https://github.com/CooooolFrog/dragonfly/issues/1990))
- searcher can not found candidate scheduler clusters, return all scheduler clusters ([#1991](https://github.com/CooooolFrog/dragonfly/issues/1991))
- oras source client ([#1983](https://github.com/CooooolFrog/dragonfly/issues/1983))
- add fail_code in scheduler's DownloadFailureCount metric ([#1981](https://github.com/CooooolFrog/dragonfly/issues/1981))
- add udp ping to the net package ([#1979](https://github.com/CooooolFrog/dragonfly/issues/1979))
- add S3ForcePathStyle to object storage ([#1976](https://github.com/CooooolFrog/dragonfly/issues/1976))
- corrupt data check ([#1946](https://github.com/CooooolFrog/dragonfly/issues/1946))
- create synchronizers concurrently ([#1941](https://github.com/CooooolFrog/dragonfly/issues/1941))
- register reflection on grpc server ([#1943](https://github.com/CooooolFrog/dragonfly/issues/1943))
- remove legacy peers support ([#1939](https://github.com/CooooolFrog/dragonfly/issues/1939))
- update fsm stable api ([#1938](https://github.com/CooooolFrog/dragonfly/issues/1938))
- add IPAddresses and DNSNames to sans of the cert ([#1930](https://github.com/CooooolFrog/dragonfly/issues/1930))
- change yaml field type from string to net.IP ([#1929](https://github.com/CooooolFrog/dragonfly/issues/1929))
- random pieces download ([#1918](https://github.com/CooooolFrog/dragonfly/issues/1918))
- update version guage metrics ([#1927](https://github.com/CooooolFrog/dragonfly/issues/1927))
- remove callsystem and pattern ([#1925](https://github.com/CooooolFrog/dragonfly/issues/1925))
- client support 'priority' parameter ([#1911](https://github.com/CooooolFrog/dragonfly/issues/1911))
- handle application not found ([#1913](https://github.com/CooooolFrog/dragonfly/issues/1913))
- update priority api ([#1912](https://github.com/CooooolFrog/dragonfly/issues/1912))
- support redis sentinal ([#1910](https://github.com/CooooolFrog/dragonfly/issues/1910))
- update submodule console ([#1908](https://github.com/CooooolFrog/dragonfly/issues/1908))
- storage collects upload piece count, peer cost and error details ([#1907](https://github.com/CooooolFrog/dragonfly/issues/1907))
- priority of the register parameter is higher than parameter of the application ([#1906](https://github.com/CooooolFrog/dragonfly/issues/1906))
- trigger task with priority ([#1904](https://github.com/CooooolFrog/dragonfly/issues/1904))
- dynconfig adds list application in scheduler ([#1903](https://github.com/CooooolFrog/dragonfly/issues/1903))
- rename url priority struct and remove PriorityLevel constants ([#1902](https://github.com/CooooolFrog/dragonfly/issues/1902))
- add priority to application in manager ([#1901](https://github.com/CooooolFrog/dragonfly/issues/1901))
- remove relation of application ([#1894](https://github.com/CooooolFrog/dragonfly/issues/1894))
- add backSourceCount validation ([#1892](https://github.com/CooooolFrog/dragonfly/issues/1892))
- add health check to service ([#1889](https://github.com/CooooolFrog/dragonfly/issues/1889))
- add pieceDownloadTimeout to docker compose template ([#1881](https://github.com/CooooolFrog/dragonfly/issues/1881))
- add the timeout of downloading piece to scheduler ([#1880](https://github.com/CooooolFrog/dragonfly/issues/1880))
- change log rotate size ([#1879](https://github.com/CooooolFrog/dragonfly/issues/1879))
- support reregister peer task in client ([#1876](https://github.com/CooooolFrog/dragonfly/issues/1876))
- if the scheduler cannot find the peer, then return Code_SchedReregister to dfdaemon ([#1875](https://github.com/CooooolFrog/dragonfly/issues/1875))
- change announcer validation ([#1869](https://github.com/CooooolFrog/dragonfly/issues/1869))
- add mysql read and write timeout ([#1868](https://github.com/CooooolFrog/dragonfly/issues/1868))
- store parent information ([#1867](https://github.com/CooooolFrog/dragonfly/issues/1867))
- remove MainParent from peer and add IsPieceBackToSource to piece
- scheduler supports storage config ([#1864](https://github.com/CooooolFrog/dragonfly/issues/1864))
- store peer download information ([#1863](https://github.com/CooooolFrog/dragonfly/issues/1863))
- manager changes filterParentLimit value ([#1859](https://github.com/CooooolFrog/dragonfly/issues/1859))
- optimize gc package ([#1855](https://github.com/CooooolFrog/dragonfly/issues/1855))
- add announcer to scheduler ([#1854](https://github.com/CooooolFrog/dragonfly/issues/1854))
- add announcer to dfdameon ([#1852](https://github.com/CooooolFrog/dragonfly/issues/1852))
- when dfdaemon disable object storage, dynconfig can't fetch manager ([#1845](https://github.com/CooooolFrog/dragonfly/issues/1845))
- optimize manager log ([#1846](https://github.com/CooooolFrog/dragonfly/issues/1846))
- scheduler adds announce host handler ([#1843](https://github.com/CooooolFrog/dragonfly/issues/1843))
- call all nodes in consistent hashing and reuse grpc connection ([#1842](https://github.com/CooooolFrog/dragonfly/issues/1842))
- update concurrent-map version ([#1837](https://github.com/CooooolFrog/dragonfly/issues/1837))
- optimize scope size is error ([#1831](https://github.com/CooooolFrog/dragonfly/issues/1831))
- add timeout grpc and job ([#1830](https://github.com/CooooolFrog/dragonfly/issues/1830))
- optimize peer log ([#1828](https://github.com/CooooolFrog/dragonfly/issues/1828))
- optional save list metadata to p2p ([#1822](https://github.com/CooooolFrog/dragonfly/issues/1822))
- add s3 resource client and recursive e2e test ([#1826](https://github.com/CooooolFrog/dragonfly/issues/1826))
- optimize preheat log ([#1827](https://github.com/CooooolFrog/dragonfly/issues/1827))
- seed peer reuses traffic ([#1825](https://github.com/CooooolFrog/dragonfly/issues/1825))
- optimize preheat ([#1824](https://github.com/CooooolFrog/dragonfly/issues/1824))
- returns an scheduling error if the peer state is not PeerStateRunning ([#1821](https://github.com/CooooolFrog/dragonfly/issues/1821))
- optimize peer gc ([#1819](https://github.com/CooooolFrog/dragonfly/issues/1819))
- peer.UpdateAt needs to be updated during download process ([#1818](https://github.com/CooooolFrog/dragonfly/issues/1818))
- statistical the traffic of reused peer ([#1816](https://github.com/CooooolFrog/dragonfly/issues/1816))
- add workHome and pluginDir to configuration ([#1807](https://github.com/CooooolFrog/dragonfly/issues/1807))
- add otel trace in log ([#1804](https://github.com/CooooolFrog/dragonfly/issues/1804))
- add leave host logger ([#1801](https://github.com/CooooolFrog/dragonfly/issues/1801))
- scheduler's record adds ParentUploadCount and ParentUploadFailedCount ([#1795](https://github.com/CooooolFrog/dragonfly/issues/1795))
- support split running tasks ([#1794](https://github.com/CooooolFrog/dragonfly/issues/1794))
- add download header log ([#1793](https://github.com/CooooolFrog/dragonfly/issues/1793))
- grpc scheduler client dial options ([#1792](https://github.com/CooooolFrog/dragonfly/issues/1792))
- daemon call leaveHost when exit ([#1788](https://github.com/CooooolFrog/dragonfly/issues/1788))
- add calculateParentHostUploadSuccessScore to scheduler ([#1789](https://github.com/CooooolFrog/dragonfly/issues/1789))
- add LeaveHost handler ([#1780](https://github.com/CooooolFrog/dragonfly/issues/1780))
- grpc_retry removes WithPerRetryTimeout ([#1763](https://github.com/CooooolFrog/dragonfly/issues/1763))
- obs object storage support ([#1758](https://github.com/CooooolFrog/dragonfly/issues/1758))
- available peer includes state is PeerStatePending ([#1756](https://github.com/CooooolFrog/dragonfly/issues/1756))
- peer will back-to-source when task switch state failed ([#1754](https://github.com/CooooolFrog/dragonfly/issues/1754))
- change FilterParentRangeLimit validation ([#1752](https://github.com/CooooolFrog/dragonfly/issues/1752))
- task state is TaskStateRunning can be registered ([#1751](https://github.com/CooooolFrog/dragonfly/issues/1751))
- gin logger rotation ([#1749](https://github.com/CooooolFrog/dragonfly/issues/1749))
- overwrite task url and url meta ([#1740](https://github.com/CooooolFrog/dragonfly/issues/1740))
- update source temporary error logic ([#1739](https://github.com/CooooolFrog/dragonfly/issues/1739))
- add seed peer back source traffic ([#1738](https://github.com/CooooolFrog/dragonfly/issues/1738))
- http request content log ([#1736](https://github.com/CooooolFrog/dragonfly/issues/1736))
- remove task and host gc ttl ([#1735](https://github.com/CooooolFrog/dragonfly/issues/1735))
- add http request log ([#1734](https://github.com/CooooolFrog/dragonfly/issues/1734))
- add TaskStateLeave to task ([#1728](https://github.com/CooooolFrog/dragonfly/issues/1728))
- searcher calculates cluster type ([#1729](https://github.com/CooooolFrog/dragonfly/issues/1729))
- unregister failed task storage ([#1717](https://github.com/CooooolFrog/dragonfly/issues/1717))
- oss get metadata ([#1724](https://github.com/CooooolFrog/dragonfly/issues/1724))
- support concurrent recursive download ([#1714](https://github.com/CooooolFrog/dragonfly/issues/1714))
- add traffic shaper for download tasks ([#1654](https://github.com/CooooolFrog/dragonfly/issues/1654))
- async create a record ([#1711](https://github.com/CooooolFrog/dragonfly/issues/1711))
- optimize storage log ([#1703](https://github.com/CooooolFrog/dragonfly/issues/1703))
- remove ipv4 and ipv6 log ([#1699](https://github.com/CooooolFrog/dragonfly/issues/1699))
- enable ipv6 in unit test ([#1698](https://github.com/CooooolFrog/dragonfly/issues/1698))
- ipv6 support ([#1685](https://github.com/CooooolFrog/dragonfly/issues/1685))
- update docker compose image ([#1696](https://github.com/CooooolFrog/dragonfly/issues/1696))
- manager add advertiseIP ([#1695](https://github.com/CooooolFrog/dragonfly/issues/1695))
- empty file e2e ([#1687](https://github.com/CooooolFrog/dragonfly/issues/1687))
- support download empty file ([#1686](https://github.com/CooooolFrog/dragonfly/issues/1686))
- stop grpc client ([#1671](https://github.com/CooooolFrog/dragonfly/issues/1671))
- change event DownloadFromBackToSource ([#1670](https://github.com/CooooolFrog/dragonfly/issues/1670))
- dfget supports config file ([#1668](https://github.com/CooooolFrog/dragonfly/issues/1668))
- split concurrent back source e2e ([#1666](https://github.com/CooooolFrog/dragonfly/issues/1666))
- support redis cluster ([#1667](https://github.com/CooooolFrog/dragonfly/issues/1667))
- source changes ResponseHeaderTimeout and ExpectContinueTimeout ([#1662](https://github.com/CooooolFrog/dragonfly/issues/1662))
- change dfdaemon rate limit ([#1661](https://github.com/CooooolFrog/dragonfly/issues/1661))
- set created_at and updated_at to timestamp ([#1659](https://github.com/CooooolFrog/dragonfly/issues/1659))
- stat peer metrics with memory cache ([#1660](https://github.com/CooooolFrog/dragonfly/issues/1660))
- change storage strategy to simple ([#1658](https://github.com/CooooolFrog/dragonfly/issues/1658))
- add missing client version for ListSchedulers ([#1657](https://github.com/CooooolFrog/dragonfly/issues/1657))
- add MaxConnectionIdle to grpc keepalive ([#1655](https://github.com/CooooolFrog/dragonfly/issues/1655))
- change consistent hashing element name ([#1652](https://github.com/CooooolFrog/dragonfly/issues/1652))
- add cert spec to security configuration ([#1621](https://github.com/CooooolFrog/dragonfly/issues/1621))
- support mutate all proxy requests ([#1623](https://github.com/CooooolFrog/dragonfly/issues/1623))
- check whether scheduler is in the same cluster ([#1620](https://github.com/CooooolFrog/dragonfly/issues/1620))
- manager add cert spec ([#1619](https://github.com/CooooolFrog/dragonfly/issues/1619))
- add tls policy to scheduler grpc server ([#1616](https://github.com/CooooolFrog/dragonfly/issues/1616))
- set tls cert leaf ([#1615](https://github.com/CooooolFrog/dragonfly/issues/1615))
- resolver addr add ServerName ([#1614](https://github.com/CooooolFrog/dragonfly/issues/1614))
- refactor grpc credential ([#1612](https://github.com/CooooolFrog/dragonfly/issues/1612))
- add tls policy to manager grpc server ([#1611](https://github.com/CooooolFrog/dragonfly/issues/1611))
- add tls policy constants ([#1610](https://github.com/CooooolFrog/dragonfly/issues/1610))
- add grpc mux transport ([#1602](https://github.com/CooooolFrog/dragonfly/issues/1602))
- manager init cert for grpc server ([#1603](https://github.com/CooooolFrog/dragonfly/issues/1603))
- refactor peertask option ([#1600](https://github.com/CooooolFrog/dragonfly/issues/1600))
- add common serialize package ([#1601](https://github.com/CooooolFrog/dragonfly/issues/1601))
- add client grpc dial timeout ([#1599](https://github.com/CooooolFrog/dragonfly/issues/1599))
- support multiple certify cache ([#1598](https://github.com/CooooolFrog/dragonfly/issues/1598))
- PeerGauge adds version and commit labels ([#1596](https://github.com/CooooolFrog/dragonfly/issues/1596))
- daemon support auto issue certificate ([#1586](https://github.com/CooooolFrog/dragonfly/issues/1586))
- add default metrics address ([#1595](https://github.com/CooooolFrog/dragonfly/issues/1595))
- grpc dial adds context ([#1594](https://github.com/CooooolFrog/dragonfly/issues/1594))
- consistent hashing add picker log ([#1593](https://github.com/CooooolFrog/dragonfly/issues/1593))
- remove golang +build tag ([#1585](https://github.com/CooooolFrog/dragonfly/issues/1585))
- manager add certificate config ([#1583](https://github.com/CooooolFrog/dragonfly/issues/1583))
- manager implements issue certificate grpc ([#1577](https://github.com/CooooolFrog/dragonfly/issues/1577))
- dfdaemon add convert interceptor ([#1582](https://github.com/CooooolFrog/dragonfly/issues/1582))
- dynconfig refresh and notify listeners ([#1579](https://github.com/CooooolFrog/dragonfly/issues/1579))
- add grpc client error interceptor ([#1575](https://github.com/CooooolFrog/dragonfly/issues/1575))
- grpc removes MaxConnectionIdle ([#1574](https://github.com/CooooolFrog/dragonfly/issues/1574))
- grpc add ratelimit ([#1572](https://github.com/CooooolFrog/dragonfly/issues/1572))
- refresh dynconfig addresses when grpc requests unavailable ([#1571](https://github.com/CooooolFrog/dragonfly/issues/1571))
- manager adds model and model version grpc api ([#1569](https://github.com/CooooolFrog/dragonfly/issues/1569))
- dynconfig add refresh func ([#1563](https://github.com/CooooolFrog/dragonfly/issues/1563))
- manager client add context ([#1562](https://github.com/CooooolFrog/dragonfly/issues/1562))
- grpc add retry middleware ([#1561](https://github.com/CooooolFrog/dragonfly/issues/1561))
- grpc consistent hashing ([#1554](https://github.com/CooooolFrog/dragonfly/issues/1554))
- model version add training result ([#1558](https://github.com/CooooolFrog/dragonfly/issues/1558))
- storage calculate the count of records ([#1557](https://github.com/CooooolFrog/dragonfly/issues/1557))
- model and model version api removes auth ([#1556](https://github.com/CooooolFrog/dragonfly/issues/1556))
- add seed trace ([#1549](https://github.com/CooooolFrog/dragonfly/issues/1549))
- gc removes logrus ([#1548](https://github.com/CooooolFrog/dragonfly/issues/1548))
- add MultiReadCloser and storage add open func ([#1546](https://github.com/CooooolFrog/dragonfly/issues/1546))
- scheduler dynconfig returns more info ([#1545](https://github.com/CooooolFrog/dragonfly/issues/1545))
- scheduler and manager change graceful stop timeout ([#1540](https://github.com/CooooolFrog/dragonfly/issues/1540))
- schedulers create main peer record ([#1539](https://github.com/CooooolFrog/dragonfly/issues/1539))
- change update model api ([#1538](https://github.com/CooooolFrog/dragonfly/issues/1538))
- manager adds model and model version api ([#1530](https://github.com/CooooolFrog/dragonfly/issues/1530))
- when the request has a range header, object storage is no need to  to calculate md5 ([#1534](https://github.com/CooooolFrog/dragonfly/issues/1534))
- support grpc recursive download ([#1518](https://github.com/CooooolFrog/dragonfly/issues/1518))
- manager embed frontend assets ([#1523](https://github.com/CooooolFrog/dragonfly/issues/1523))
- can not return peer with the same host ([#1526](https://github.com/CooooolFrog/dragonfly/issues/1526))
- add daemon-socket-path ([#1521](https://github.com/CooooolFrog/dragonfly/issues/1521))
- store preheat result ([#1516](https://github.com/CooooolFrog/dragonfly/issues/1516))
- replace grpc package with https://github.com/dragonflyoss/api ([#1515](https://github.com/CooooolFrog/dragonfly/issues/1515))
- dfdaemon add Authorization and WWWAuthenticate headers ([#1513](https://github.com/CooooolFrog/dragonfly/issues/1513))
- auto switch to concurrent back source based on download speed ([#1494](https://github.com/CooooolFrog/dragonfly/issues/1494))
- enable dependabot ([#1501](https://github.com/CooooolFrog/dragonfly/issues/1501))
- scheduler adds filter range limit ([#1497](https://github.com/CooooolFrog/dragonfly/issues/1497))
- scheduler set workhome ([#1493](https://github.com/CooooolFrog/dragonfly/issues/1493))
- remove test print
- rename steal peers to candidate peers ([#1476](https://github.com/CooooolFrog/dragonfly/issues/1476))
- scheduler merge end of piece and piece from seed peer ([#1474](https://github.com/CooooolFrog/dragonfly/issues/1474))
- dfdaemon add default healthy config ([#1472](https://github.com/CooooolFrog/dragonfly/issues/1472))
- dag adds LenVertex and RangeVertex func ([#1470](https://github.com/CooooolFrog/dragonfly/issues/1470))
- generate dag mock
- add directed acyclic graph package ([#1468](https://github.com/CooooolFrog/dragonfly/issues/1468))
- proxy add defaultTag field ([#1462](https://github.com/CooooolFrog/dragonfly/issues/1462))
- manager support postgres ([#1459](https://github.com/CooooolFrog/dragonfly/issues/1459))
- use os.PathSeparator to generate object key
- scheduler add data dir ([#1453](https://github.com/CooooolFrog/dragonfly/issues/1453))
- dfdaemon is compatible with v2.0.2 ([#1452](https://github.com/CooooolFrog/dragonfly/issues/1452))
- add slices util package
- reload proxy option ([#1443](https://github.com/CooooolFrog/dragonfly/issues/1443))
- if peer back-to-source failed, return source metadata. ([#1444](https://github.com/CooooolFrog/dragonfly/issues/1444))
- report peer result with source error detail ([#1439](https://github.com/CooooolFrog/dragonfly/issues/1439))
- add dfstore command ([#1441](https://github.com/CooooolFrog/dragonfly/issues/1441))
- back source error detail ([#1437](https://github.com/CooooolFrog/dragonfly/issues/1437))
- change local cache ttl ([#1436](https://github.com/CooooolFrog/dragonfly/issues/1436))
- if service can not found fqdn, replace fqdn with hostname ([#1435](https://github.com/CooooolFrog/dragonfly/issues/1435))
- remove errors package ([#1434](https://github.com/CooooolFrog/dragonfly/issues/1434))
- concurrent multiple pieces back source ([#1426](https://github.com/CooooolFrog/dragonfly/issues/1426))
- dfstore closes writer ([#1424](https://github.com/CooooolFrog/dragonfly/issues/1424))
- use put object action ([#1422](https://github.com/CooooolFrog/dragonfly/issues/1422))
- GetObjectInput add range field ([#1421](https://github.com/CooooolFrog/dragonfly/issues/1421))
- rename client/clientutil to client/util ([#1420](https://github.com/CooooolFrog/dragonfly/issues/1420))
- rewrite interface{} to any ([#1419](https://github.com/CooooolFrog/dragonfly/issues/1419))
- update namely/protoc-all image version to 1.47_0 ([#1418](https://github.com/CooooolFrog/dragonfly/issues/1418))
- update golang to 1.18.3 ([#1417](https://github.com/CooooolFrog/dragonfly/issues/1417))
- remove github/pkg/errors package ([#1416](https://github.com/CooooolFrog/dragonfly/issues/1416))
- add dfstore client interface ([#1415](https://github.com/CooooolFrog/dragonfly/issues/1415))
- scheduler http status ([#1414](https://github.com/CooooolFrog/dragonfly/issues/1414))
- enable configuration of the tls parameter for the mysql connection. i.e. tls=preferred ([#1300](https://github.com/CooooolFrog/dragonfly/issues/1300))
- import object to seed peer with max replicas ([#1413](https://github.com/CooooolFrog/dragonfly/issues/1413))
- object storage add filter field ([#1412](https://github.com/CooooolFrog/dragonfly/issues/1412))
- dfdaemon add destroyObject rest api ([#1410](https://github.com/CooooolFrog/dragonfly/issues/1410))
- client add create object storage ([#1409](https://github.com/CooooolFrog/dragonfly/issues/1409))
- seed peer add object storage port ([#1408](https://github.com/CooooolFrog/dragonfly/issues/1408))
- rename digest func and add new digest func ([#1405](https://github.com/CooooolFrog/dragonfly/issues/1405))
- dfdaemon upload and object storage service add middlewares ([#1404](https://github.com/CooooolFrog/dragonfly/issues/1404))
- remove cdn ([#1401](https://github.com/CooooolFrog/dragonfly/issues/1401))
- redirect stdout and stderr to file ([#1399](https://github.com/CooooolFrog/dragonfly/issues/1399))
- dfdaemon add GetObject rest api ([#1398](https://github.com/CooooolFrog/dragonfly/issues/1398))
- add seed peer for list scheduler grpc interface ([#1393](https://github.com/CooooolFrog/dragonfly/issues/1393))
- dfdaemon add object storage rest api ([#1390](https://github.com/CooooolFrog/dragonfly/issues/1390))
- replace gin-gonic/gin with gorilla/mux ([#1389](https://github.com/CooooolFrog/dragonfly/issues/1389))
- add enable config to peer gauge ([#1382](https://github.com/CooooolFrog/dragonfly/issues/1382))
- dfdaemon add ns filter ([#1379](https://github.com/CooooolFrog/dragonfly/issues/1379))
- remove connection gc ([#1378](https://github.com/CooooolFrog/dragonfly/issues/1378))
- dynconfig add object storage ([#1369](https://github.com/CooooolFrog/dragonfly/issues/1369))
- manager add bucket interface ([#1368](https://github.com/CooooolFrog/dragonfly/issues/1368))
- add objectstorage pkg ([#1366](https://github.com/CooooolFrog/dragonfly/issues/1366))
- remove preheat tag validate with required ([#1363](https://github.com/CooooolFrog/dragonfly/issues/1363))
- e2e seed peer ([#1358](https://github.com/CooooolFrog/dragonfly/issues/1358))
- update console and helm-charts submodule ([#1355](https://github.com/CooooolFrog/dragonfly/issues/1355))
- use uid/gid as UserID and UserGroup if current user not found in passwd ([#1352](https://github.com/CooooolFrog/dragonfly/issues/1352))
- use 127.0.0.1 as IPv4 if there's no external IPv4 addr ([#1353](https://github.com/CooooolFrog/dragonfly/issues/1353))
- add security group id with scheduler cluster ([#1354](https://github.com/CooooolFrog/dragonfly/issues/1354))
- change pattern from cdn to seed peer and remove kustomize shell ([#1345](https://github.com/CooooolFrog/dragonfly/issues/1345))
- update casbin/gorm-adapter version and change e2e charts config
- update helm charts
- update dependencies
- add seed peer metrics ([#1342](https://github.com/CooooolFrog/dragonfly/issues/1342))
- grpc health probe support arm64 ([#1338](https://github.com/CooooolFrog/dragonfly/issues/1338))
- docker build with multi platforms ([#1337](https://github.com/CooooolFrog/dragonfly/issues/1337))
- add sync piece watchdog ([#1272](https://github.com/CooooolFrog/dragonfly/issues/1272))
- scheduler handles seed peer failed ([#1325](https://github.com/CooooolFrog/dragonfly/issues/1325))
- custom preheat tag parameters ([#1324](https://github.com/CooooolFrog/dragonfly/issues/1324))
- client add tls verify config ([#1323](https://github.com/CooooolFrog/dragonfly/issues/1323))
- scheduler register interface return task type ([#1318](https://github.com/CooooolFrog/dragonfly/issues/1318))
- get active peer count ([#1315](https://github.com/CooooolFrog/dragonfly/issues/1315))
- reduce dynconfig log ([#1312](https://github.com/CooooolFrog/dragonfly/issues/1312))
- back source when receive seed request ([#1309](https://github.com/CooooolFrog/dragonfly/issues/1309))
- update helm charts submodule ([#1308](https://github.com/CooooolFrog/dragonfly/issues/1308))
- add vsock network type support ([#1303](https://github.com/CooooolFrog/dragonfly/issues/1303))
- support response header ([#1292](https://github.com/CooooolFrog/dragonfly/issues/1292))
- add seed peer logic ([#1302](https://github.com/CooooolFrog/dragonfly/issues/1302))
- support more digest like sha256 ([#1304](https://github.com/CooooolFrog/dragonfly/issues/1304))
- implement pattern in client daemon ([#1231](https://github.com/CooooolFrog/dragonfly/issues/1231))
- scheduler add seed peer ([#1298](https://github.com/CooooolFrog/dragonfly/issues/1298))
- don't gc client rpc connection if ExpireTime is 0 ([#1296](https://github.com/CooooolFrog/dragonfly/issues/1296))
- update scheduler PeerResult validation ([#1294](https://github.com/CooooolFrog/dragonfly/issues/1294))
- manager add seed peer ([#1293](https://github.com/CooooolFrog/dragonfly/issues/1293))
- implement client seed mode ([#1247](https://github.com/CooooolFrog/dragonfly/issues/1247))
- scheduler peer result validation ([#1288](https://github.com/CooooolFrog/dragonfly/issues/1288))
- use a golang native file server to replace nginx ([#1258](https://github.com/CooooolFrog/dragonfly/issues/1258))
- support build arm&&arm64 dfget ([#1286](https://github.com/CooooolFrog/dragonfly/issues/1286))
- update filter parent ([#1279](https://github.com/CooooolFrog/dragonfly/issues/1279))
- in tree plugin ([#1276](https://github.com/CooooolFrog/dragonfly/issues/1276))
- move dfnet to pkg dir ([#1265](https://github.com/CooooolFrog/dragonfly/issues/1265))
- add dfcache rpm/deb packages and man pages and publish in goreleaser ([#1259](https://github.com/CooooolFrog/dragonfly/issues/1259))
- add AnnounceTask and StatTask metrics ([#1256](https://github.com/CooooolFrog/dragonfly/issues/1256))
- define and implement new dfdaemon APIs to make dragonfly2 work as a distributed cache ([#1227](https://github.com/CooooolFrog/dragonfly/issues/1227))
- redirect daemon stdout stderr to file ([#1244](https://github.com/CooooolFrog/dragonfly/issues/1244))
- registerTask returns to the task in time ([#1250](https://github.com/CooooolFrog/dragonfly/issues/1250))
- docker-compose write log to file ([#1236](https://github.com/CooooolFrog/dragonfly/issues/1236))
- update docker compose version ([#1235](https://github.com/CooooolFrog/dragonfly/issues/1235))
- update to v2.0.2 ([#1232](https://github.com/CooooolFrog/dragonfly/issues/1232))
- scheduler blocks steal peers ([#1224](https://github.com/CooooolFrog/dragonfly/issues/1224))
- update manager console ([#1222](https://github.com/CooooolFrog/dragonfly/issues/1222))
- manager validate with config ([#1218](https://github.com/CooooolFrog/dragonfly/issues/1218))
- remove kustomize template ([#1216](https://github.com/CooooolFrog/dragonfly/issues/1216))
- add back source fail metric in client ([#1214](https://github.com/CooooolFrog/dragonfly/issues/1214))
- cannot delete a cluster with existing instances ([#1213](https://github.com/CooooolFrog/dragonfly/issues/1213))
- add type to DownloadFailureCount ([#1212](https://github.com/CooooolFrog/dragonfly/issues/1212))
- if the number of failed peers in the task is greater than FailedPeerCountLimit, then scheduler notifies running peers of failure ([#1211](https://github.com/CooooolFrog/dragonfly/issues/1211))
- optimize get available task ([#1208](https://github.com/CooooolFrog/dragonfly/issues/1208))
- change scheduler and cdn listen ([#1205](https://github.com/CooooolFrog/dragonfly/issues/1205))
- scheduler add block peers set ([#1202](https://github.com/CooooolFrog/dragonfly/issues/1202))
- add grpc-health-probe to image ([#1196](https://github.com/CooooolFrog/dragonfly/issues/1196))
- add grpc health interface ([#1195](https://github.com/CooooolFrog/dragonfly/issues/1195))
- remove grpc error code validate ([#1191](https://github.com/CooooolFrog/dragonfly/issues/1191))
- generate grpc protos in namely/protoc-all image ([#1187](https://github.com/CooooolFrog/dragonfly/issues/1187))
- scheduler addresses log ([#1183](https://github.com/CooooolFrog/dragonfly/issues/1183))
- manage GetCDN interface return scheduler info ([#1184](https://github.com/CooooolFrog/dragonfly/issues/1184))
- dfdaemon match scheduler with case insensitive ([#1181](https://github.com/CooooolFrog/dragonfly/issues/1181))
- add RBAC to manager config interface ([#1179](https://github.com/CooooolFrog/dragonfly/issues/1179))
- dfdaemon get available scheduler addresses in the same cluster ([#1178](https://github.com/CooooolFrog/dragonfly/issues/1178))
- implement grpc client side sync pieces ([#1167](https://github.com/CooooolFrog/dragonfly/issues/1167))
- seacher return multiple scheduler clusters ([#1175](https://github.com/CooooolFrog/dragonfly/issues/1175))
- replace time.Now().Sub by time.Since ([#1173](https://github.com/CooooolFrog/dragonfly/issues/1173))
- change DefaultServerOptions to variable
- change default scheduler filter parent limit ([#1166](https://github.com/CooooolFrog/dragonfly/issues/1166))
- implement bidirectional fetch pieces ([#1165](https://github.com/CooooolFrog/dragonfly/issues/1165))
- scheduler add default biz tag ([#1164](https://github.com/CooooolFrog/dragonfly/issues/1164))
- optimize proxy performance ([#1137](https://github.com/CooooolFrog/dragonfly/issues/1137))
- host remove peer ([#1161](https://github.com/CooooolFrog/dragonfly/issues/1161))
- change reschdule config ([#1158](https://github.com/CooooolFrog/dragonfly/issues/1158))
- update git submodule ([#1153](https://github.com/CooooolFrog/dragonfly/issues/1153))
- scheduler metrics add default value of biz tag ([#1151](https://github.com/CooooolFrog/dragonfly/issues/1151))
- add user update interface and rename rest to service ([#1148](https://github.com/CooooolFrog/dragonfly/issues/1148))
- scheduler trace trigger cdn ([#1147](https://github.com/CooooolFrog/dragonfly/issues/1147))
- add scheduler traffic metrics ([#1143](https://github.com/CooooolFrog/dragonfly/issues/1143))
- update otel package version and fix otelgrpc goroutine leak ([#1141](https://github.com/CooooolFrog/dragonfly/issues/1141))
- add scheduler metrics ([#1139](https://github.com/CooooolFrog/dragonfly/issues/1139))
- scheduler remove inactive host ([#1135](https://github.com/CooooolFrog/dragonfly/issues/1135))
- task state for register ([#1132](https://github.com/CooooolFrog/dragonfly/issues/1132))
- change grpc client keepalive config ([#1125](https://github.com/CooooolFrog/dragonfly/issues/1125))
- scheduler change piece cost from nanosecond to millisecond ([#1119](https://github.com/CooooolFrog/dragonfly/issues/1119))
- support health probe in daemon ([#1120](https://github.com/CooooolFrog/dragonfly/issues/1120))
- when peer downloads finished, peer deletes parent ([#1116](https://github.com/CooooolFrog/dragonfly/issues/1116))
- change source client dialer config ([#1115](https://github.com/CooooolFrog/dragonfly/issues/1115))
- optimize scheduler log ([#1114](https://github.com/CooooolFrog/dragonfly/issues/1114))
- remove needless manager grpc proxy ([#1113](https://github.com/CooooolFrog/dragonfly/issues/1113))
- set grpc logger verbosity from env variable ([#1111](https://github.com/CooooolFrog/dragonfly/issues/1111))
- change back-to-source timeout ([#1112](https://github.com/CooooolFrog/dragonfly/issues/1112))
- optimize scheduler ([#1106](https://github.com/CooooolFrog/dragonfly/issues/1106))
- reuse partial completed task ([#1107](https://github.com/CooooolFrog/dragonfly/issues/1107))
- optimize depth limit func ([#1102](https://github.com/CooooolFrog/dragonfly/issues/1102))
- change client default load limit ([#1104](https://github.com/CooooolFrog/dragonfly/issues/1104))
- limit tree depth ([#1099](https://github.com/CooooolFrog/dragonfly/issues/1099))
- update load limit ([#1097](https://github.com/CooooolFrog/dragonfly/issues/1097))
- optimize peer range ([#1095](https://github.com/CooooolFrog/dragonfly/issues/1095))
- add cdn addresses log ([#1091](https://github.com/CooooolFrog/dragonfly/issues/1091))
- scheduler add limit count of filter parent func ([#1090](https://github.com/CooooolFrog/dragonfly/issues/1090))
- merge ranged request storage into parent ([#1078](https://github.com/CooooolFrog/dragonfly/issues/1078))
- add dynamic parallel count ([#1088](https://github.com/CooooolFrog/dragonfly/issues/1088))
- fix docker-compose ([#1087](https://github.com/CooooolFrog/dragonfly/issues/1087))
- add prefetch metric in client ([#1068](https://github.com/CooooolFrog/dragonfly/issues/1068))
- when scheduler blocks cdn, resource does not initialize cdn ([#1081](https://github.com/CooooolFrog/dragonfly/issues/1081))
- scheduler blocks cdn ([#1079](https://github.com/CooooolFrog/dragonfly/issues/1079))
- job trigger cdn by resource ([#1076](https://github.com/CooooolFrog/dragonfly/issues/1076))
- add client request log ([#1069](https://github.com/CooooolFrog/dragonfly/issues/1069))
- support change console log level ([#1055](https://github.com/CooooolFrog/dragonfly/issues/1055))
- manager support mysql ssl connection ([#1015](https://github.com/CooooolFrog/dragonfly/issues/1015))
- remove host and task when peer make tree ([#1042](https://github.com/CooooolFrog/dragonfly/issues/1042))
- cdn download tiny file ([#1040](https://github.com/CooooolFrog/dragonfly/issues/1040))
- If cdn only updates IP, set cdn peers state to PeerStateLeave ([#1038](https://github.com/CooooolFrog/dragonfly/issues/1038))
- generate grpc protoc ([#1027](https://github.com/CooooolFrog/dragonfly/issues/1027))
- manager config model add is_boot key ([#1025](https://github.com/CooooolFrog/dragonfly/issues/1025))
- scheduler download tiny file with range header ([#1024](https://github.com/CooooolFrog/dragonfly/issues/1024))
- change compatibility version to v2.0.2-rc.0 ([#1017](https://github.com/CooooolFrog/dragonfly/issues/1017))
- when cdn peer is failed, peer should be back-to-source ([#1005](https://github.com/CooooolFrog/dragonfly/issues/1005))
- add actions job timout ([#1008](https://github.com/CooooolFrog/dragonfly/issues/1008))
- set peer state to running when scope size is SizeScope_TINY ([#1004](https://github.com/CooooolFrog/dragonfly/issues/1004))
- update submodule charts ([#1002](https://github.com/CooooolFrog/dragonfly/issues/1002))
- task mutex replace sync kmutex ([#1000](https://github.com/CooooolFrog/dragonfly/issues/1000))
- stream send error code ([#986](https://github.com/CooooolFrog/dragonfly/issues/986))
- trace https proxy request ([#996](https://github.com/CooooolFrog/dragonfly/issues/996))
- add scheduler host gc ([#989](https://github.com/CooooolFrog/dragonfly/issues/989))
- update typo in local_storage.go ([#955](https://github.com/CooooolFrog/dragonfly/issues/955))
- update charts submodule version ([#985](https://github.com/CooooolFrog/dragonfly/issues/985))
- change task and peer ttl ([#984](https://github.com/CooooolFrog/dragonfly/issues/984))
- when write last piece, generate digest ([#982](https://github.com/CooooolFrog/dragonfly/issues/982))
- merge same tasks in daemon ([#977](https://github.com/CooooolFrog/dragonfly/issues/977))
- if cdn is deleted, clear cdn related information ([#967](https://github.com/CooooolFrog/dragonfly/issues/967))
- add default DiskGCThresholdPercent and ignore it when is 0 ([#971](https://github.com/CooooolFrog/dragonfly/issues/971))
- improve redirect to allow url rewrite ([#969](https://github.com/CooooolFrog/dragonfly/issues/969))
- Add useProxies to registryMirror allowing to mirror more anything ([#965](https://github.com/CooooolFrog/dragonfly/issues/965))
- change metrics port to 8000 ([#964](https://github.com/CooooolFrog/dragonfly/issues/964))
- add daemon metrics support ([#960](https://github.com/CooooolFrog/dragonfly/issues/960))
- support disk usage gc in client ([#953](https://github.com/CooooolFrog/dragonfly/issues/953))
- update source.Response and source client interface ([#945](https://github.com/CooooolFrog/dragonfly/issues/945))
- remove stat log from scheduler ([#946](https://github.com/CooooolFrog/dragonfly/issues/946))
- support recursive download in dfget ([#932](https://github.com/CooooolFrog/dragonfly/issues/932))
- add kmutex and krwmutex ([#934](https://github.com/CooooolFrog/dragonfly/issues/934))
- make idgen package public ([#931](https://github.com/CooooolFrog/dragonfly/issues/931))
- make dfpath public ([#929](https://github.com/CooooolFrog/dragonfly/issues/929))
- dfdaemon list scheduler cluster with multi idc ([#917](https://github.com/CooooolFrog/dragonfly/issues/917))
- update submodule ([#916](https://github.com/CooooolFrog/dragonfly/issues/916))
- update task access time ([#909](https://github.com/CooooolFrog/dragonfly/issues/909))
- optmize dfget package upgrade support ([#804](https://github.com/CooooolFrog/dragonfly/issues/804))
- support create container without docker-compose ([#915](https://github.com/CooooolFrog/dragonfly/issues/915))
- add data directory ([#910](https://github.com/CooooolFrog/dragonfly/issues/910))
- add data storage directory  ([#907](https://github.com/CooooolFrog/dragonfly/issues/907))
- dfdaemon update content length ([#895](https://github.com/CooooolFrog/dragonfly/issues/895))
- lint sh ([#876](https://github.com/CooooolFrog/dragonfly/issues/876))
- update helm charts ([#870](https://github.com/CooooolFrog/dragonfly/issues/870))
- update version to v2.0.1 ([#869](https://github.com/CooooolFrog/dragonfly/issues/869))
- add oauth timeout ([#867](https://github.com/CooooolFrog/dragonfly/issues/867))
- support customize transport in daemon ([#866](https://github.com/CooooolFrog/dragonfly/issues/866))
- console ([#865](https://github.com/CooooolFrog/dragonfly/issues/865))
- move dfnet to internal ([#862](https://github.com/CooooolFrog/dragonfly/issues/862))
- remove ifaceutils pkg ([#860](https://github.com/CooooolFrog/dragonfly/issues/860))
- move syncmap pkg([#859](https://github.com/CooooolFrog/dragonfly/issues/859))
- oauth interface auth ([#857](https://github.com/CooooolFrog/dragonfly/issues/857))
- add scopes validation ([#856](https://github.com/CooooolFrog/dragonfly/issues/856))
- log ([#852](https://github.com/CooooolFrog/dragonfly/issues/852))
- get scheduler list with advertise ip ([#848](https://github.com/CooooolFrog/dragonfly/issues/848))
- support mutli manager addrs ([#846](https://github.com/CooooolFrog/dragonfly/issues/846))
- searcher plugin change return params ([#844](https://github.com/CooooolFrog/dragonfly/issues/844))
- plugin log ([#843](https://github.com/CooooolFrog/dragonfly/issues/843))
- export searcher evaluate func ([#842](https://github.com/CooooolFrog/dragonfly/issues/842))
- add context for FindSchedulerCluster ([#841](https://github.com/CooooolFrog/dragonfly/issues/841))
- add application cdn clusters field ([#840](https://github.com/CooooolFrog/dragonfly/issues/840))
- update console submodule ([#838](https://github.com/CooooolFrog/dragonfly/issues/838))
- preheat compatible with harbor ([#837](https://github.com/CooooolFrog/dragonfly/issues/837))
- gin version ([#833](https://github.com/CooooolFrog/dragonfly/issues/833))
- update manager image ([#831](https://github.com/CooooolFrog/dragonfly/issues/831))
- update helm charts version ([#824](https://github.com/CooooolFrog/dragonfly/issues/824))
- add package reachable ([#822](https://github.com/CooooolFrog/dragonfly/issues/822))
- support list plugin ([#819](https://github.com/CooooolFrog/dragonfly/issues/819))
- scheduler and cdn report fqdn to manager ([#818](https://github.com/CooooolFrog/dragonfly/issues/818))
- dfdaemon get scheduler list dynamically from manager ([#812](https://github.com/CooooolFrog/dragonfly/issues/812))
- update image-spec version ([#808](https://github.com/CooooolFrog/dragonfly/issues/808))
- add security rule ([#806](https://github.com/CooooolFrog/dragonfly/issues/806))
- add idgen peer id ([#800](https://github.com/CooooolFrog/dragonfly/issues/800))
- optimize scheduler peer stat log ([#798](https://github.com/CooooolFrog/dragonfly/issues/798))
- replace sortedList with sortedUniqueList ([#793](https://github.com/CooooolFrog/dragonfly/issues/793))
- calculate piece metadata digest ([#787](https://github.com/CooooolFrog/dragonfly/issues/787))
- preheat skip certificate validation ([#786](https://github.com/CooooolFrog/dragonfly/issues/786))
- support traffic metrics by peer host ([#776](https://github.com/CooooolFrog/dragonfly/issues/776))
- support dump http content in client for debugging ([#770](https://github.com/CooooolFrog/dragonfly/issues/770))
- remove calculate total count service ([#772](https://github.com/CooooolFrog/dragonfly/issues/772))
- add user list interface ([#771](https://github.com/CooooolFrog/dragonfly/issues/771))
- clear hashcircler and maputils package ([#768](https://github.com/CooooolFrog/dragonfly/issues/768))
- add cdn task peers monitor log ([#764](https://github.com/CooooolFrog/dragonfly/issues/764))
- change config key name ([#759](https://github.com/CooooolFrog/dragonfly/issues/759))
- scheduler channel blocking ([#756](https://github.com/CooooolFrog/dragonfly/issues/756))
- add jobs api ([#751](https://github.com/CooooolFrog/dragonfly/issues/751))
- add config ([#746](https://github.com/CooooolFrog/dragonfly/issues/746))
- add preheat otel ([#741](https://github.com/CooooolFrog/dragonfly/issues/741))
- add job logger ([#740](https://github.com/CooooolFrog/dragonfly/issues/740))
- manager add grpc jaeger ([#738](https://github.com/CooooolFrog/dragonfly/issues/738))
- load limit ([#739](https://github.com/CooooolFrog/dragonfly/issues/739))
- preheat cluster ([#731](https://github.com/CooooolFrog/dragonfly/issues/731))
- nsswitch ([#737](https://github.com/CooooolFrog/dragonfly/issues/737))
- export e2e logs ([#732](https://github.com/CooooolFrog/dragonfly/issues/732))
- compatible with V1 preheat  ([#720](https://github.com/CooooolFrog/dragonfly/issues/720))
- add manager client list scheduler interface ([#694](https://github.com/CooooolFrog/dragonfly/issues/694))
- release fd ([#668](https://github.com/CooooolFrog/dragonfly/issues/668))
- add otel trace ([#650](https://github.com/CooooolFrog/dragonfly/issues/650))
- disable prepared statement ([#648](https://github.com/CooooolFrog/dragonfly/issues/648))
- update verison ([#640](https://github.com/CooooolFrog/dragonfly/issues/640))
- changelog ([#638](https://github.com/CooooolFrog/dragonfly/issues/638))
- update console submodule ([#637](https://github.com/CooooolFrog/dragonfly/issues/637))
- update submodule ([#632](https://github.com/CooooolFrog/dragonfly/issues/632))
- beautify scheduler & CDN log ([#618](https://github.com/CooooolFrog/dragonfly/issues/618))
- Print version information when the system starts up ([#620](https://github.com/CooooolFrog/dragonfly/issues/620))
- add piece download timeout ([#621](https://github.com/CooooolFrog/dragonfly/issues/621))
- notice client back source when rescheduled parent reach max times ([#611](https://github.com/CooooolFrog/dragonfly/issues/611))
- avoid report peer result fail due to context cancel & add backsource tracer ([#606](https://github.com/CooooolFrog/dragonfly/issues/606))
- optimize cdn check free space ([#603](https://github.com/CooooolFrog/dragonfly/issues/603))
- client back source ([#579](https://github.com/CooooolFrog/dragonfly/issues/579))
- enable manager user's features ([#598](https://github.com/CooooolFrog/dragonfly/issues/598))
- add sni proxy support ([#600](https://github.com/CooooolFrog/dragonfly/issues/600))
- compatibility e2e with matrix ([#599](https://github.com/CooooolFrog/dragonfly/issues/599))
- change scheduler cluster query params ([#596](https://github.com/CooooolFrog/dragonfly/issues/596))
- add oauth2 signin ([#591](https://github.com/CooooolFrog/dragonfly/issues/591))
- update scheduler cluster query params ([#587](https://github.com/CooooolFrog/dragonfly/issues/587))
- add time out when register ([#588](https://github.com/CooooolFrog/dragonfly/issues/588))
- skip verify when back to source ([#586](https://github.com/CooooolFrog/dragonfly/issues/586))
- update charts submodule ([#583](https://github.com/CooooolFrog/dragonfly/issues/583))
- support limit from dfget client ([#578](https://github.com/CooooolFrog/dragonfly/issues/578))
- add cdn cluster id for scheduler cluster ([#580](https://github.com/CooooolFrog/dragonfly/issues/580))
- start process ([#572](https://github.com/CooooolFrog/dragonfly/issues/572))
- gin log to file ([#574](https://github.com/CooooolFrog/dragonfly/issues/574))
- add manager cors middleware ([#573](https://github.com/CooooolFrog/dragonfly/issues/573))
- change rabc code struct ([#552](https://github.com/CooooolFrog/dragonfly/issues/552))
- empty scheduler job ([#565](https://github.com/CooooolFrog/dragonfly/issues/565))
- optimize manager startup process ([#562](https://github.com/CooooolFrog/dragonfly/issues/562))
- update git submodule ([#560](https://github.com/CooooolFrog/dragonfly/issues/560))
- optimize scheduler start server ([#558](https://github.com/CooooolFrog/dragonfly/issues/558))
- add console ([#559](https://github.com/CooooolFrog/dragonfly/issues/559))
- generate swagger api ([#557](https://github.com/CooooolFrog/dragonfly/issues/557))
- add console submodule ([#549](https://github.com/CooooolFrog/dragonfly/issues/549))
- optimize get permission name ([#548](https://github.com/CooooolFrog/dragonfly/issues/548))
- rename task to job ([#544](https://github.com/CooooolFrog/dragonfly/issues/544))
- Add distribute Schedule Tracer & Refactor scheduler ([#537](https://github.com/CooooolFrog/dragonfly/issues/537))
- add artifacthub badge ([#524](https://github.com/CooooolFrog/dragonfly/issues/524))
- update cdn host ([#530](https://github.com/CooooolFrog/dragonfly/issues/530))
- back source when no available peers or scheduler error ([#521](https://github.com/CooooolFrog/dragonfly/issues/521))
- add task manager ([#490](https://github.com/CooooolFrog/dragonfly/issues/490))
- rename manager grpc ([#510](https://github.com/CooooolFrog/dragonfly/issues/510))
- Add stress testing tool for daemon ([#506](https://github.com/CooooolFrog/dragonfly/issues/506))
- scheduler getevaluator lock ([#502](https://github.com/CooooolFrog/dragonfly/issues/502))
- rename search file to searcher ([#484](https://github.com/CooooolFrog/dragonfly/issues/484))
- Add schedule log ([#495](https://github.com/CooooolFrog/dragonfly/issues/495))
- Extract peer event processing function ([#489](https://github.com/CooooolFrog/dragonfly/issues/489))
- optimize scheduler dynconfig ([#480](https://github.com/CooooolFrog/dragonfly/issues/480))
- optimize jwt ([#476](https://github.com/CooooolFrog/dragonfly/issues/476))
- register service to manager ([#475](https://github.com/CooooolFrog/dragonfly/issues/475))
- add searcher to scheduler cluster ([#462](https://github.com/CooooolFrog/dragonfly/issues/462))
- CDN implementation supports HDFS type storage ([#420](https://github.com/CooooolFrog/dragonfly/issues/420))
- add is_default to scheduler_cluster table ([#458](https://github.com/CooooolFrog/dragonfly/issues/458))
- add host info for scheduler and cdn ([#457](https://github.com/CooooolFrog/dragonfly/issues/457))
- Install e2e script ([#451](https://github.com/CooooolFrog/dragonfly/issues/451))
- Manager user logic ([#419](https://github.com/CooooolFrog/dragonfly/issues/419))
- Add plugin support for resource ([#291](https://github.com/CooooolFrog/dragonfly/issues/291))
- changelog ([#326](https://github.com/CooooolFrog/dragonfly/issues/326))
- remove queue package ([#275](https://github.com/CooooolFrog/dragonfly/issues/275))
- add ci badge ([#265](https://github.com/CooooolFrog/dragonfly/issues/265))
- remove slidingwindow and assertutils package ([#263](https://github.com/CooooolFrog/dragonfly/issues/263))

### Feature
- prefetch ranged requests ([#1053](https://github.com/CooooolFrog/dragonfly/issues/1053))
- support e2e feature gates ([#1056](https://github.com/CooooolFrog/dragonfly/issues/1056))
- change log level in-flight ([#1023](https://github.com/CooooolFrog/dragonfly/issues/1023))
- update helm charts submodule ([#567](https://github.com/CooooolFrog/dragonfly/issues/567))
- Add manager charts with submodule ([#525](https://github.com/CooooolFrog/dragonfly/issues/525))
- support mysql 5.6 ([#520](https://github.com/CooooolFrog/dragonfly/issues/520))
- support customize base image ([#519](https://github.com/CooooolFrog/dragonfly/issues/519))
- add kustomize yaml for deploying ([#349](https://github.com/CooooolFrog/dragonfly/issues/349))
- support basic auth for proxy ([#250](https://github.com/CooooolFrog/dragonfly/issues/250))
- add disk quota gc for daemon ([#215](https://github.com/CooooolFrog/dragonfly/issues/215))

### Feature
- refresh proto file ([#615](https://github.com/CooooolFrog/dragonfly/issues/615))
- optimize manager project layout ([#540](https://github.com/CooooolFrog/dragonfly/issues/540))
- enable grpc tracing ([#531](https://github.com/CooooolFrog/dragonfly/issues/531))
- update multiple registries support docs ([#481](https://github.com/CooooolFrog/dragonfly/issues/481))
- add multiple registry mirrors support ([#479](https://github.com/CooooolFrog/dragonfly/issues/479))
- disable proxy when config is empty ([#455](https://github.com/CooooolFrog/dragonfly/issues/455))
- add pod labels in helm chart ([#447](https://github.com/CooooolFrog/dragonfly/issues/447))
- optimize failed reason not set ([#446](https://github.com/CooooolFrog/dragonfly/issues/446))
- report peer result when failed to register ([#433](https://github.com/CooooolFrog/dragonfly/issues/433))
- rename PeerHost to Daemon in client ([#438](https://github.com/CooooolFrog/dragonfly/issues/438))
- export peer.TaskManager for embedding dragonfly in custom binary ([#434](https://github.com/CooooolFrog/dragonfly/issues/434))
- optimize error message for proxy ([#428](https://github.com/CooooolFrog/dragonfly/issues/428))
- minimize daemon runtime capabilities ([#421](https://github.com/CooooolFrog/dragonfly/issues/421))
- add default filter in proxy for deployment and docs ([#417](https://github.com/CooooolFrog/dragonfly/issues/417))
- add jaeger for helm deployment ([#415](https://github.com/CooooolFrog/dragonfly/issues/415))
- update dfdaemon proxy port comment
- update cdn init container template ([#399](https://github.com/CooooolFrog/dragonfly/issues/399))
- update client config to Camel-Case format ([#393](https://github.com/CooooolFrog/dragonfly/issues/393))
- update helm charts deploy guide ([#386](https://github.com/CooooolFrog/dragonfly/issues/386))
- update helm charts ([#385](https://github.com/CooooolFrog/dragonfly/issues/385))
- support setns in client ([#378](https://github.com/CooooolFrog/dragonfly/issues/378))
- disable resolver server config ([#314](https://github.com/CooooolFrog/dragonfly/issues/314))
- update docs ([#307](https://github.com/CooooolFrog/dragonfly/issues/307))
- remove unsafe code in client/daemon/storage ([#258](https://github.com/CooooolFrog/dragonfly/issues/258))
- remove redundant configurations ([#216](https://github.com/CooooolFrog/dragonfly/issues/216))

### Ffix
- typo in Makefile ([#975](https://github.com/CooooolFrog/dragonfly/issues/975))

### Fix
- Interval in SyncProbesResponse ([#2466](https://github.com/CooooolFrog/dragonfly/issues/2466))
- e2e test dfget recursive ([#2458](https://github.com/CooooolFrog/dragonfly/issues/2458))
- announcer in scheduler ([#2451](https://github.com/CooooolFrog/dragonfly/issues/2451))
- delete host in network topology ([#2417](https://github.com/CooooolFrog/dragonfly/issues/2417))
- call MakeNamespaceKeyInScheduler function error ([#2383](https://github.com/CooooolFrog/dragonfly/issues/2383))
- package declaration error ([#2379](https://github.com/CooooolFrog/dragonfly/issues/2379))
- evaluate after filter ([#2363](https://github.com/CooooolFrog/dragonfly/issues/2363))
- when bufferSize is zero, storage can not write data to file ([#2366](https://github.com/CooooolFrog/dragonfly/issues/2366))
- SyncPieceViaHTTPS not work ([#2329](https://github.com/CooooolFrog/dragonfly/issues/2329))
- object downloads failed by dfstore when dfdaemon enabled concurrent ([#2328](https://github.com/CooooolFrog/dragonfly/issues/2328))
- redis validation in scheduler config ([#2287](https://github.com/CooooolFrog/dragonfly/issues/2287))
- local dynconfig panic in Notify ([#2266](https://github.com/CooooolFrog/dragonfly/issues/2266))
- client grpc dial non-block ([#2261](https://github.com/CooooolFrog/dragonfly/issues/2261))
- modify the traversal condition for Items ([#2250](https://github.com/CooooolFrog/dragonfly/issues/2250))
- ip and hostname params in FindSchedulerClusters ([#2249](https://github.com/CooooolFrog/dragonfly/issues/2249))
- traffic shaper record task not found ([#2226](https://github.com/CooooolFrog/dragonfly/issues/2226))
- fsm events failed when register task ([#2225](https://github.com/CooooolFrog/dragonfly/issues/2225))
- stat DownloadPeerCount and DownloadPieceCount ([#2180](https://github.com/CooooolFrog/dragonfly/issues/2180))
- manager metrics Subsystem ([#2175](https://github.com/CooooolFrog/dragonfly/issues/2175))
- remove unnecessary fmt.Sprintf calls ([#2159](https://github.com/CooooolFrog/dragonfly/issues/2159))
- validate daemon gcInterval config ([#2118](https://github.com/CooooolFrog/dragonfly/issues/2118))
- unregister task from scheduler in storage.deleteTask ([#2100](https://github.com/CooooolFrog/dragonfly/issues/2100))
- backsource first piece timeout ([#2083](https://github.com/CooooolFrog/dragonfly/issues/2083))
- peer GC clear all peers when peer's count large than PeerCountLimitForTask ([#2061](https://github.com/CooooolFrog/dragonfly/issues/2061))
- spelling mistakes ([#2027](https://github.com/CooooolFrog/dragonfly/issues/2027))
- dferror not convert ([#2001](https://github.com/CooooolFrog/dragonfly/issues/2001))
- dfstore typo ([#2000](https://github.com/CooooolFrog/dragonfly/issues/2000))
- manager typo ([#1995](https://github.com/CooooolFrog/dragonfly/issues/1995))
- daemon recognize Code_SchedForbidden ([#1994](https://github.com/CooooolFrog/dragonfly/issues/1994))
- count of total page in pagination ([#1993](https://github.com/CooooolFrog/dragonfly/issues/1993))
- manager grpc filename ([#1992](https://github.com/CooooolFrog/dragonfly/issues/1992))
- client bitMap extend capacity ([#1973](https://github.com/CooooolFrog/dragonfly/issues/1973))
- context of trigger seed peer ([#1971](https://github.com/CooooolFrog/dragonfly/issues/1971))
- config decode net.IP ([#1964](https://github.com/CooooolFrog/dragonfly/issues/1964))
- download context cancelled ([#1942](https://github.com/CooooolFrog/dragonfly/issues/1942))
- peer keepalive with manager ([#1940](https://github.com/CooooolFrog/dragonfly/issues/1940))
- panic caused by hashring not being built ([#1928](https://github.com/CooooolFrog/dragonfly/issues/1928))
- application not found ([#1924](https://github.com/CooooolFrog/dragonfly/issues/1924))
- remove advertiseIP config in e2e ([#1878](https://github.com/CooooolFrog/dragonfly/issues/1878))
- recursive download always return none error ([#1841](https://github.com/CooooolFrog/dragonfly/issues/1841))
- expire header timezone ([#1840](https://github.com/CooooolFrog/dragonfly/issues/1840))
- otel goroutine leak ([#1815](https://github.com/CooooolFrog/dragonfly/issues/1815))
- gorm-adaptor pkg version ([#1805](https://github.com/CooooolFrog/dragonfly/issues/1805))
- leave host ([#1803](https://github.com/CooooolFrog/dragonfly/issues/1803))
- daemon don't leaveHost when keepStorage=true ([#1790](https://github.com/CooooolFrog/dragonfly/issues/1790))
- did not call scheduler leave tasks in forceGC ([#1782](https://github.com/CooooolFrog/dragonfly/issues/1782))
- plugin builder ([#1775](https://github.com/CooooolFrog/dragonfly/issues/1775))
- add fallback registry mirror in gen-containerd-host.sh ([#1774](https://github.com/CooooolFrog/dragonfly/issues/1774))
- open end range in concurrent back source ([#1764](https://github.com/CooooolFrog/dragonfly/issues/1764))
- manager PeerGauge ([#1761](https://github.com/CooooolFrog/dragonfly/issues/1761))
- backsource temporary error judgement ([#1726](https://github.com/CooooolFrog/dragonfly/issues/1726))
- gorm can not update is_default field ([#1731](https://github.com/CooooolFrog/dragonfly/issues/1731))
- content length is zero when task succeed ([#1732](https://github.com/CooooolFrog/dragonfly/issues/1732))
- docker compose config ([#1713](https://github.com/CooooolFrog/dragonfly/issues/1713))
- hdfs not registered ([#1702](https://github.com/CooooolFrog/dragonfly/issues/1702))
- grpc download tidy file error ([#1697](https://github.com/CooooolFrog/dragonfly/issues/1697))
- manager redis config convert ([#1680](https://github.com/CooooolFrog/dragonfly/issues/1680))
- task CanBackToSource func ([#1663](https://github.com/CooooolFrog/dragonfly/issues/1663))
- manager embed assets ([#1642](https://github.com/CooooolFrog/dragonfly/issues/1642))
- dfstore flags invalid ([#1641](https://github.com/CooooolFrog/dragonfly/issues/1641))
- ci actions with docker ([#1613](https://github.com/CooooolFrog/dragonfly/issues/1613))
- dfdaemon can not shutdown ([#1580](https://github.com/CooooolFrog/dragonfly/issues/1580))
- scheduler can not exit gracefully due to machinery fatal log ([#1573](https://github.com/CooooolFrog/dragonfly/issues/1573))
- scheduler and manager tracing ([#1551](https://github.com/CooooolFrog/dragonfly/issues/1551))
- scheduler's MainParent func ([#1550](https://github.com/CooooolFrog/dragonfly/issues/1550))
- check same peer id for sync pieces ([#1525](https://github.com/CooooolFrog/dragonfly/issues/1525))
- auto switch to concurrent back source ([#1507](https://github.com/CooooolFrog/dragonfly/issues/1507))
- wait first peer packet fail ([#1500](https://github.com/CooooolFrog/dragonfly/issues/1500))
- one piece task sometimes backsource after succeed ([#1499](https://github.com/CooooolFrog/dragonfly/issues/1499))
- random vertices ([#1496](https://github.com/CooooolFrog/dragonfly/issues/1496))
- dfstore command-line tool name ([#1492](https://github.com/CooooolFrog/dragonfly/issues/1492))
- oss client judge directory bug ([#1488](https://github.com/CooooolFrog/dragonfly/issues/1488))
- dfdaemon unix socket ([#1489](https://github.com/CooooolFrog/dragonfly/issues/1489))
- init storage error ([#1486](https://github.com/CooooolFrog/dragonfly/issues/1486))
- back source error ([#1485](https://github.com/CooooolFrog/dragonfly/issues/1485))
- keepalive with ip
- upload_manager write header in time ([#1471](https://github.com/CooooolFrog/dragonfly/issues/1471))
- infinite loop in peer.Ancestors() ([#1469](https://github.com/CooooolFrog/dragonfly/issues/1469))
- upload_manager write header immediately when it is ready ([#1466](https://github.com/CooooolFrog/dragonfly/issues/1466))
- metrics reduces labels ([#1457](https://github.com/CooooolFrog/dragonfly/issues/1457))
- depth limit ([#1451](https://github.com/CooooolFrog/dragonfly/issues/1451))
- dfpath creates redundant directories ([#1446](https://github.com/CooooolFrog/dragonfly/issues/1446))
- release package name ([#1442](https://github.com/CooooolFrog/dragonfly/issues/1442))
- seed task metric panic ([#1427](https://github.com/CooooolFrog/dragonfly/issues/1427))
- pkg/strings comment typo
- downloadFromSource() doesn't validate response ([#1400](https://github.com/CooooolFrog/dragonfly/issues/1400))
- default repository does not exist and missing dependency containers ([#1395](https://github.com/CooooolFrog/dragonfly/issues/1395))
- validate rate limiter ([#1392](https://github.com/CooooolFrog/dragonfly/issues/1392))
- dfget ratelimit params ([#1391](https://github.com/CooooolFrog/dragonfly/issues/1391))
- count error & totalPage error ([#1373](https://github.com/CooooolFrog/dragonfly/issues/1373)) ([#1376](https://github.com/CooooolFrog/dragonfly/issues/1376))
- manager router middlewares order ([#1385](https://github.com/CooooolFrog/dragonfly/issues/1385))
- dfget build error ([#1381](https://github.com/CooooolFrog/dragonfly/issues/1381))
- preheat tack id ([#1375](https://github.com/CooooolFrog/dragonfly/issues/1375))
- register fail panic ([#1351](https://github.com/CooooolFrog/dragonfly/issues/1351))
- find partial completed overflow ([#1346](https://github.com/CooooolFrog/dragonfly/issues/1346))
- e2e charts config
- seed peer reuse value
- dfdaemon seed peer metrics namespace ([#1343](https://github.com/CooooolFrog/dragonfly/issues/1343))
- create_at timestamp ([#1341](https://github.com/CooooolFrog/dragonfly/issues/1341))
- reuse seed peer id is not exist ([#1335](https://github.com/CooooolFrog/dragonfly/issues/1335))
- reuse non-end range request ([#1333](https://github.com/CooooolFrog/dragonfly/issues/1333))
- http range header validation ([#1334](https://github.com/CooooolFrog/dragonfly/issues/1334))
- nfpms maintainer ([#1326](https://github.com/CooooolFrog/dragonfly/issues/1326))
- reuse seed panic ([#1319](https://github.com/CooooolFrog/dragonfly/issues/1319))
- seed peer did not send done seed result and no content length send ([#1316](https://github.com/CooooolFrog/dragonfly/issues/1316))
- remove container after generating protoc ([#1306](https://github.com/CooooolFrog/dragonfly/issues/1306))
- digest reader and unit tests ([#1305](https://github.com/CooooolFrog/dragonfly/issues/1305))
- scheduler typo ([#1297](https://github.com/CooooolFrog/dragonfly/issues/1297))
- keep accept header ([#1291](https://github.com/CooooolFrog/dragonfly/issues/1291))
- grpc total_piece_count comment ([#1289](https://github.com/CooooolFrog/dragonfly/issues/1289))
- run.sh threw error on mac ([#1285](https://github.com/CooooolFrog/dragonfly/issues/1285))
- docker compose run.sh ([#1282](https://github.com/CooooolFrog/dragonfly/issues/1282))
- legacy cdn peer ([#1283](https://github.com/CooooolFrog/dragonfly/issues/1283))
- filter parent condition ([#1277](https://github.com/CooooolFrog/dragonfly/issues/1277))
- dfget daemon console log invalid ([#1275](https://github.com/CooooolFrog/dragonfly/issues/1275))
- scheduler config validation ([#1274](https://github.com/CooooolFrog/dragonfly/issues/1274))
- run.sh threw error on mac ([#1273](https://github.com/CooooolFrog/dragonfly/issues/1273))
- tree infinite loop ([#1271](https://github.com/CooooolFrog/dragonfly/issues/1271))
- acquire empty dst pid ([#1268](https://github.com/CooooolFrog/dragonfly/issues/1268))
- skip unsupported kernel in systemd service ([#1261](https://github.com/CooooolFrog/dragonfly/issues/1261))
- client synchronizer report error lock and dial grpc timeout ([#1260](https://github.com/CooooolFrog/dragonfly/issues/1260))
- prevent traversal tree from infinite loop ([#1266](https://github.com/CooooolFrog/dragonfly/issues/1266))
- error message ([#1255](https://github.com/CooooolFrog/dragonfly/issues/1255))
- client sync piece panic ([#1246](https://github.com/CooooolFrog/dragonfly/issues/1246))
- client superfluous usage gc ([#1243](https://github.com/CooooolFrog/dragonfly/issues/1243))
- client sync send unsafe call ([#1240](https://github.com/CooooolFrog/dragonfly/issues/1240))
- client unexpected timeout ([#1239](https://github.com/CooooolFrog/dragonfly/issues/1239))
- goreleaser config
- make generate ([#1228](https://github.com/CooooolFrog/dragonfly/issues/1228))
- calculate FreeUploadLoad ([#1226](https://github.com/CooooolFrog/dragonfly/issues/1226))
- sync pieces hang ([#1221](https://github.com/CooooolFrog/dragonfly/issues/1221))
- client miss failed piece ([#1194](https://github.com/CooooolFrog/dragonfly/issues/1194))
- client break error ([#1190](https://github.com/CooooolFrog/dragonfly/issues/1190))
- rpc cdn sync piece tasks ([#1168](https://github.com/CooooolFrog/dragonfly/issues/1168))
- subscriber data race ([#1169](https://github.com/CooooolFrog/dragonfly/issues/1169))
- docker-compose run with mac throw error ([#1134](https://github.com/CooooolFrog/dragonfly/issues/1134))
- wrong md5 sign in cdn ([#1126](https://github.com/CooooolFrog/dragonfly/issues/1126))
- docker-compose preheat pending ([#1124](https://github.com/CooooolFrog/dragonfly/issues/1124))
- scheduler piece cost time ([#1118](https://github.com/CooooolFrog/dragonfly/issues/1118))
- when peer state is PeerStateSucceeded, return size scope is small ([#1103](https://github.com/CooooolFrog/dragonfly/issues/1103))
- delete peer's parent on PeerEventDownloadSucceeded event ([#1085](https://github.com/CooooolFrog/dragonfly/issues/1085))
- pull request template typo ([#1080](https://github.com/CooooolFrog/dragonfly/issues/1080))
- when cdn download failed, scheduler should set cdn peer state PeerStateFailed ([#1067](https://github.com/CooooolFrog/dragonfly/issues/1067))
- evaluate peer's parent ([#1064](https://github.com/CooooolFrog/dragonfly/issues/1064))
- scheduler download tiny file error ([#1052](https://github.com/CooooolFrog/dragonfly/issues/1052))
- docker actions typo ([#1041](https://github.com/CooooolFrog/dragonfly/issues/1041))
- cdn trigger peer error ([#1035](https://github.com/CooooolFrog/dragonfly/issues/1035))
- retrigger cdn panic ([#1034](https://github.com/CooooolFrog/dragonfly/issues/1034))
- calculate piece MD5 sign when last piece download ([#1006](https://github.com/CooooolFrog/dragonfly/issues/1006))
- register task with size scope ([#1003](https://github.com/CooooolFrog/dragonfly/issues/1003))
- when scheduler is not available, replace the scheduler client ([#999](https://github.com/CooooolFrog/dragonfly/issues/999))
- total pieces count not set cause digest invalid ([#992](https://github.com/CooooolFrog/dragonfly/issues/992))
- send piece result error not handled ([#987](https://github.com/CooooolFrog/dragonfly/issues/987))
- scheduler config typo ([#983](https://github.com/CooooolFrog/dragonfly/issues/983))
- schedulers send invalid direct piece ([#970](https://github.com/CooooolFrog/dragonfly/issues/970))
- use 'parent' as mainPeer in PeerPacket in removePeerFromCurrentTree() ([#957](https://github.com/CooooolFrog/dragonfly/issues/957))
- size scope empty ([#941](https://github.com/CooooolFrog/dragonfly/issues/941))
- not handle base.Code_SchedTaskStatusError in client ([#938](https://github.com/CooooolFrog/dragonfly/issues/938))
- infinitely get pieces when piece num is invalid ([#926](https://github.com/CooooolFrog/dragonfly/issues/926))
- plugin dir is empty ([#922](https://github.com/CooooolFrog/dragonfly/issues/922))
- peer gc ([#918](https://github.com/CooooolFrog/dragonfly/issues/918))
- go plugin test build error ([#912](https://github.com/CooooolFrog/dragonfly/issues/912))
- typo ([#911](https://github.com/CooooolFrog/dragonfly/issues/911))
- total pieces not set when back source ([#908](https://github.com/CooooolFrog/dragonfly/issues/908))
- mismatch digest peer task did not mark invalid ([#903](https://github.com/CooooolFrog/dragonfly/issues/903))
- dfget dfpath ([#901](https://github.com/CooooolFrog/dragonfly/issues/901))
- scheduler success event ([#891](https://github.com/CooooolFrog/dragonfly/issues/891))
- add cdn cluster to scheduler cluster ([#887](https://github.com/CooooolFrog/dragonfly/issues/887))
- small size task failed due to digest error ([#886](https://github.com/CooooolFrog/dragonfly/issues/886))
- searcher log ([#878](https://github.com/CooooolFrog/dragonfly/issues/878))
- error log ([#863](https://github.com/CooooolFrog/dragonfly/issues/863))
- file peer task back source digest not match ([#849](https://github.com/CooooolFrog/dragonfly/issues/849))
- manager typo and cdn peer id ([#809](https://github.com/CooooolFrog/dragonfly/issues/809))
- cdn AdvertiseIP not used ([#782](https://github.com/CooooolFrog/dragonfly/issues/782))
- add peer to task failed because InnerBucketMaxLength is small ([#765](https://github.com/CooooolFrog/dragonfly/issues/765))
- back source weight ([#762](https://github.com/CooooolFrog/dragonfly/issues/762))
- client load ([#753](https://github.com/CooooolFrog/dragonfly/issues/753))
- peer empty parent ([#724](https://github.com/CooooolFrog/dragonfly/issues/724))
- client panic ([#719](https://github.com/CooooolFrog/dragonfly/issues/719))
- client goroutine and fd leak ([#713](https://github.com/CooooolFrog/dragonfly/issues/713))
- skip check DisableAutoBackSource option when scheduler says back source ([#693](https://github.com/CooooolFrog/dragonfly/issues/693))
- go library cve ([#666](https://github.com/CooooolFrog/dragonfly/issues/666))
- return failed piece ([#619](https://github.com/CooooolFrog/dragonfly/issues/619))
- use string slice for header ([#601](https://github.com/CooooolFrog/dragonfly/issues/601))
- preheat-e2e timeout ([#602](https://github.com/CooooolFrog/dragonfly/issues/602))
- use getTask instead of taskStore.Get, for the error cause type ([#571](https://github.com/CooooolFrog/dragonfly/issues/571))
- adjust dfget download log ([#564](https://github.com/CooooolFrog/dragonfly/issues/564))
- wait available peer packet panic ([#561](https://github.com/CooooolFrog/dragonfly/issues/561))
- wrong content length in proxy
- cdn back source range size overflow ([#550](https://github.com/CooooolFrog/dragonfly/issues/550))
- scheduler concurrent dead lock ([#509](https://github.com/CooooolFrog/dragonfly/issues/509))
- scheduler pick candidate and associate child  encounter  dead lock ([#500](https://github.com/CooooolFrog/dragonfly/issues/500))
- generate proto file ([#483](https://github.com/CooooolFrog/dragonfly/issues/483))
- address typo ([#468](https://github.com/CooooolFrog/dragonfly/issues/468))
- dead lock when pt.failedPieceCh is full ([#466](https://github.com/CooooolFrog/dragonfly/issues/466))
- user table typo ([#453](https://github.com/CooooolFrog/dragonfly/issues/453))
- log specification ([#452](https://github.com/CooooolFrog/dragonfly/issues/452))
- wrong cache header ([#423](https://github.com/CooooolFrog/dragonfly/issues/423))
- close net namespace fd ([#418](https://github.com/CooooolFrog/dragonfly/issues/418))
- update static cdn config
- wrong daemon config and kubectl image tag ([#398](https://github.com/CooooolFrog/dragonfly/issues/398))
- update mapsturcture decode and remove unused config ([#396](https://github.com/CooooolFrog/dragonfly/issues/396))
- update DynconfigOptions typo ([#390](https://github.com/CooooolFrog/dragonfly/issues/390))
- gc test ([#370](https://github.com/CooooolFrog/dragonfly/issues/370))
- scheduler panic ([#356](https://github.com/CooooolFrog/dragonfly/issues/356))
- use seederName to replace the PeerID to generate the UUID ([#355](https://github.com/CooooolFrog/dragonfly/issues/355))
- check health too long when dfdaemon is unavailable ([#344](https://github.com/CooooolFrog/dragonfly/issues/344))
- when load config from cdn directory in dynconfig, skip sub directories ([#310](https://github.com/CooooolFrog/dragonfly/issues/310))
- Makefile and build.sh ([#309](https://github.com/CooooolFrog/dragonfly/issues/309))
- ci badge ([#281](https://github.com/CooooolFrog/dragonfly/issues/281))
- change peerPacketReady to buffer channel ([#256](https://github.com/CooooolFrog/dragonfly/issues/256))
- cdn gc dead lock ([#231](https://github.com/CooooolFrog/dragonfly/issues/231))
- cfgFile nil error ([#224](https://github.com/CooooolFrog/dragonfly/issues/224))
- change manager docs path ([#193](https://github.com/CooooolFrog/dragonfly/issues/193))
- **manager:** modify to config from scheduler_config in swagger yaml ([#317](https://github.com/CooooolFrog/dragonfly/issues/317))

### Fix
- [scheduler]  destPeer keepalive when downloaded by other peer ([#1865](https://github.com/CooooolFrog/dragonfly/issues/1865))
- source plugin not loaded ([#811](https://github.com/CooooolFrog/dragonfly/issues/811))
- proxy for stress testing tool ([#507](https://github.com/CooooolFrog/dragonfly/issues/507))
- add process level for scheduler peer task status ([#435](https://github.com/CooooolFrog/dragonfly/issues/435))
- infinite recursion in MkDirAll ([#358](https://github.com/CooooolFrog/dragonfly/issues/358))
- use atomic to avoid data race in client ([#254](https://github.com/CooooolFrog/dragonfly/issues/254))

### Refactor
- trainer server module ([#2486](https://github.com/CooooolFrog/dragonfly/issues/2486))
- network topology package ([#2412](https://github.com/CooooolFrog/dragonfly/issues/2412))
- probes package in network topology ([#2382](https://github.com/CooooolFrog/dragonfly/issues/2382))
- network topology package in scheduler ([#2380](https://github.com/CooooolFrog/dragonfly/issues/2380))
- improve the performance of the code ([#2162](https://github.com/CooooolFrog/dragonfly/issues/2162))
- optimize certifyCache Get function ([#2160](https://github.com/CooooolFrog/dragonfly/issues/2160))
- preheat job ([#2113](https://github.com/CooooolFrog/dragonfly/issues/2113))
- support reload scheduler addresses for local Dynconfig in client ([#2107](https://github.com/CooooolFrog/dragonfly/issues/2107))
- scheduling with v2 grpc ([#2104](https://github.com/CooooolFrog/dragonfly/issues/2104))
- package digest ([#2085](https://github.com/CooooolFrog/dragonfly/issues/2085))
- type of digest in task ([#2084](https://github.com/CooooolFrog/dragonfly/issues/2084))
- task.SizeScope with v2 grpc in scheduler ([#2082](https://github.com/CooooolFrog/dragonfly/issues/2082))
- task piece with v2 grpc ([#2080](https://github.com/CooooolFrog/dragonfly/issues/2080))
- resource task with v2 version of grpc ([#2078](https://github.com/CooooolFrog/dragonfly/issues/2078))
- parse http range ([#2071](https://github.com/CooooolFrog/dragonfly/issues/2071))
- peer resource with v2 version of the grpc ([#2039](https://github.com/CooooolFrog/dragonfly/issues/2039))
- announcer and dynconfig with v2 verison of the manager grpc ([#2037](https://github.com/CooooolFrog/dragonfly/issues/2037))
- resource host without scheduler v1 definition ([#2036](https://github.com/CooooolFrog/dragonfly/issues/2036))
- piece_dispatcher considering score of parent peer ([#1978](https://github.com/CooooolFrog/dragonfly/issues/1978))
- dynconfig without Unmarshal ([#1926](https://github.com/CooooolFrog/dragonfly/issues/1926))
- back-to-source configuration ([#1895](https://github.com/CooooolFrog/dragonfly/issues/1895))
- scheduler registers task ([#1766](https://github.com/CooooolFrog/dragonfly/issues/1766))
- obs of objectstorage pkg ([#1762](https://github.com/CooooolFrog/dragonfly/issues/1762))
- idgen pkg ([#1715](https://github.com/CooooolFrog/dragonfly/issues/1715))
- pkg basic ([#1712](https://github.com/CooooolFrog/dragonfly/issues/1712))
- manager and scheduler config ([#1701](https://github.com/CooooolFrog/dragonfly/issues/1701))
- listenIP and advertiseIP ([#1694](https://github.com/CooooolFrog/dragonfly/issues/1694))
- dfpath for certify cache dir ([#1618](https://github.com/CooooolFrog/dragonfly/issues/1618))
- dfnet package ([#1578](https://github.com/CooooolFrog/dragonfly/issues/1578))
- dfdaemon client and remove rpc connection pool ([#1576](https://github.com/CooooolFrog/dragonfly/issues/1576))
- set and dag with generics ([#1490](https://github.com/CooooolFrog/dragonfly/issues/1490))
- cache key for peer ([#1483](https://github.com/CooooolFrog/dragonfly/issues/1483))
- scheduler training configuration
- dag GetSourceVertices and GetSinkVertices func
- rewrite math max and min with generics ([#1447](https://github.com/CooooolFrog/dragonfly/issues/1447))
- scheduler announce task ([#1407](https://github.com/CooooolFrog/dragonfly/issues/1407))
- digest package ([#1403](https://github.com/CooooolFrog/dragonfly/issues/1403))
- pkg util ([#1402](https://github.com/CooooolFrog/dragonfly/issues/1402))
- scheduler grpc ([#1310](https://github.com/CooooolFrog/dragonfly/issues/1310))
- scheduler task SizeScope ([#1287](https://github.com/CooooolFrog/dragonfly/issues/1287))
- scheduler end and begin of piece ([#1189](https://github.com/CooooolFrog/dragonfly/issues/1189))
- manager grpc server ([#1047](https://github.com/CooooolFrog/dragonfly/issues/1047))
- scheduler grpc server ([#1046](https://github.com/CooooolFrog/dragonfly/issues/1046))
- docker workflows ([#1039](https://github.com/CooooolFrog/dragonfly/issues/1039))
- scheduler register task ([#924](https://github.com/CooooolFrog/dragonfly/issues/924))
- move from io/ioutil to io and os packages ([#906](https://github.com/CooooolFrog/dragonfly/issues/906))
- dfpath pkg ([#879](https://github.com/CooooolFrog/dragonfly/issues/879))
- scheduler evaluator ([#805](https://github.com/CooooolFrog/dragonfly/issues/805))
- scheduler supervisor ([#655](https://github.com/CooooolFrog/dragonfly/issues/655))
- rbac
- user interface
- manager server new instance ([#464](https://github.com/CooooolFrog/dragonfly/issues/464))
- update arch ([#319](https://github.com/CooooolFrog/dragonfly/issues/319))
- remove benchmark-rate and rename not-back-source ([#245](https://github.com/CooooolFrog/dragonfly/issues/245))
- support multi digest not only md5 ([#236](https://github.com/CooooolFrog/dragonfly/issues/236))
- simplify to make imports more format ([#230](https://github.com/CooooolFrog/dragonfly/issues/230))
- **manager:** modify mysql table schema, orm json tag. ([#283](https://github.com/CooooolFrog/dragonfly/issues/283))

### Test
- add client daemon network topology unit tests ([#2490](https://github.com/CooooolFrog/dragonfly/issues/2490))
- add TestAnnouncer_Serve test
- optimize announcer in scheduler ([#2448](https://github.com/CooooolFrog/dragonfly/issues/2448))
- announcer adds tests ([#2377](https://github.com/CooooolFrog/dragonfly/issues/2377))
- add storage unit tests of trainer ([#2437](https://github.com/CooooolFrog/dragonfly/issues/2437))
- add ParseProbedCountKeyInScheduler and Snapshot tests ([#2438](https://github.com/CooooolFrog/dragonfly/issues/2438))
- optimize network topology and probes unit tests ([#2425](https://github.com/CooooolFrog/dragonfly/issues/2425))
- add delete host unit tests ([#2424](https://github.com/CooooolFrog/dragonfly/issues/2424))
- add unit test for Bytes.Set ([#2422](https://github.com/CooooolFrog/dragonfly/issues/2422))
- add unit test for MakeProbedAtKeyInScheduler ([#2423](https://github.com/CooooolFrog/dragonfly/issues/2423))
- add probes and network topology unit tests ([#2414](https://github.com/CooooolFrog/dragonfly/issues/2414))
- optimize network topology and probes tests ([#2409](https://github.com/CooooolFrog/dragonfly/issues/2409))
- add unit test for ComputePieceCount ([#2401](https://github.com/CooooolFrog/dragonfly/issues/2401))
- add probes and network topology unit tests ([#2390](https://github.com/CooooolFrog/dragonfly/issues/2390))
- add unit tests in pkg/redis ([#2384](https://github.com/CooooolFrog/dragonfly/issues/2384))
- add slice packege tests ([#2386](https://github.com/CooooolFrog/dragonfly/issues/2386))
- add test case "new dfpath by dataDir" ([#2368](https://github.com/CooooolFrog/dragonfly/issues/2368))
- improve timeout in recursive download ([#2367](https://github.com/CooooolFrog/dragonfly/issues/2367))
- add new metrics test to service ([#2212](https://github.com/CooooolFrog/dragonfly/issues/2212))
- improve Test_parseByte ([#2173](https://github.com/CooooolFrog/dragonfly/issues/2173))
- add UT for byte String function ([#2170](https://github.com/CooooolFrog/dragonfly/issues/2170))
- improve TestMin ([#2168](https://github.com/CooooolFrog/dragonfly/issues/2168))
- add UT for MustParseRang ([#2158](https://github.com/CooooolFrog/dragonfly/issues/2158))
- improve TestFilterQuery ([#2157](https://github.com/CooooolFrog/dragonfly/issues/2157))
- add Validate test to scheduler config ([#2129](https://github.com/CooooolFrog/dragonfly/issues/2129))
- add Validate test to manager config ([#2128](https://github.com/CooooolFrog/dragonfly/issues/2128))
- refactor client validate ut ([#2126](https://github.com/CooooolFrog/dragonfly/issues/2126))
- add unit tests for DaemonConfig.Validate() ([#2119](https://github.com/CooooolFrog/dragonfly/issues/2119))
- remove random test in pieceDispatcherTest ([#2106](https://github.com/CooooolFrog/dragonfly/issues/2106))
- remove test main ([#1710](https://github.com/CooooolFrog/dragonfly/issues/1710))
- add test for daemon rpcserver ([#1704](https://github.com/CooooolFrog/dragonfly/issues/1704))
- find parent with ancestor ([#1482](https://github.com/CooooolFrog/dragonfly/issues/1482))
- update e2e charts config
- watchdog
- close dfget back-to-souce ([#1317](https://github.com/CooooolFrog/dragonfly/issues/1317))
- fix storage backups ([#1270](https://github.com/CooooolFrog/dragonfly/issues/1270))
- scheduler storage ([#1257](https://github.com/CooooolFrog/dragonfly/issues/1257))
- AnnounceTask and StatTask ([#1254](https://github.com/CooooolFrog/dragonfly/issues/1254))
- fix e2e preheat case ([#1170](https://github.com/CooooolFrog/dragonfly/issues/1170))
- cache expire interval ([#1160](https://github.com/CooooolFrog/dragonfly/issues/1160))
- add scheduler constructSuccessPeerPacket case ([#1154](https://github.com/CooooolFrog/dragonfly/issues/1154))
- scheduler service handlePieceFail ([#1146](https://github.com/CooooolFrog/dragonfly/issues/1146))
- FilterParentCount ([#1094](https://github.com/CooooolFrog/dragonfly/issues/1094))
- scheduler handle failed piece ([#1084](https://github.com/CooooolFrog/dragonfly/issues/1084))
- dump goroutine in e2e ([#980](https://github.com/CooooolFrog/dragonfly/issues/980))
- idgen peer id ([#913](https://github.com/CooooolFrog/dragonfly/issues/913))
- preheat image ([#794](https://github.com/CooooolFrog/dragonfly/issues/794))
- scheduler supervisor ([#742](https://github.com/CooooolFrog/dragonfly/issues/742))
- preheat e2e ([#627](https://github.com/CooooolFrog/dragonfly/issues/627))
- print merge commit ([#581](https://github.com/CooooolFrog/dragonfly/issues/581))
- compare image commit ([#538](https://github.com/CooooolFrog/dragonfly/issues/538))
- E2E download concurrency ([#467](https://github.com/CooooolFrog/dragonfly/issues/467))
- E2E test use kind's containerd ([#448](https://github.com/CooooolFrog/dragonfly/issues/448))
- manager config ([#392](https://github.com/CooooolFrog/dragonfly/issues/392))
- idgen add digest ([#243](https://github.com/CooooolFrog/dragonfly/issues/243))


<a name="v2.1.0-alpha.3"></a>
## [v2.1.0-alpha.3] - 2023-03-30
### Feat
- add logger.CoreLogger to searcher plugin ([#2232](https://github.com/CooooolFrog/dragonfly/issues/2232))


<a name="v2.1.0-alpha.2"></a>
## [v2.1.0-alpha.2] - 2023-03-30
### Chore
- update traffic shaper log ([#2227](https://github.com/CooooolFrog/dragonfly/issues/2227))

### Feat
- add log to searcher plugin ([#2231](https://github.com/CooooolFrog/dragonfly/issues/2231))

### Fix
- traffic shaper record task not found ([#2226](https://github.com/CooooolFrog/dragonfly/issues/2226))


<a name="v2.1.0-alpha.1"></a>
## [v2.1.0-alpha.1] - 2023-03-28
### Chore
- format ci action
- add Mohammed Farooq to MAINTAINERS ([#2211](https://github.com/CooooolFrog/dragonfly/issues/2211))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.224 to 1.44.229 ([#2221](https://github.com/CooooolFrog/dragonfly/issues/2221))
- **deps:** bump github.com/go-playground/validator/v10 from 10.11.2 to 10.12.0 ([#2220](https://github.com/CooooolFrog/dragonfly/issues/2220))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.37.0 to 0.40.0 ([#2219](https://github.com/CooooolFrog/dragonfly/issues/2219))
- **deps:** bump github.com/aliyun/aliyun-oss-go-sdk from 2.2.6+incompatible to 2.2.7+incompatible ([#2218](https://github.com/CooooolFrog/dragonfly/issues/2218))
- **deps:** bump gorm.io/driver/postgres from 1.4.8 to 1.5.0 ([#2217](https://github.com/CooooolFrog/dragonfly/issues/2217))

### Feat
- add probes struct ([#2190](https://github.com/CooooolFrog/dragonfly/issues/2190))
- add trainer config in scheduler ([#2214](https://github.com/CooooolFrog/dragonfly/issues/2214))
- add tfserving service to rpc package ([#2210](https://github.com/CooooolFrog/dragonfly/issues/2210))
- add trainer service to rpc package ([#2209](https://github.com/CooooolFrog/dragonfly/issues/2209))
- rename security client file name ([#2208](https://github.com/CooooolFrog/dragonfly/issues/2208))
- add CreateModel func to manager grpc client ([#2207](https://github.com/CooooolFrog/dragonfly/issues/2207))
- rename SecurityService to Security ([#2206](https://github.com/CooooolFrog/dragonfly/issues/2206))
- rename HostName to Hostname ([#2205](https://github.com/CooooolFrog/dragonfly/issues/2205))

### Fix
- fsm events failed when register task ([#2225](https://github.com/CooooolFrog/dragonfly/issues/2225))

### Test
- add new metrics test to service ([#2212](https://github.com/CooooolFrog/dragonfly/issues/2212))


<a name="v2.1.0-alpha.0"></a>
## [v2.1.0-alpha.0] - 2023-03-21

<a name="v2.0.10"></a>
## [v2.0.10] - 2023-12-11
### Chore
- migrate from k8s.gcr.io to registry.k8s.io ([#2186](https://github.com/CooooolFrog/dragonfly/issues/2186))
- add Mohammed Farooq to MAINTAINERS ([#2211](https://github.com/CooooolFrog/dragonfly/issues/2211))
- change codecov rules ([#2174](https://github.com/CooooolFrog/dragonfly/issues/2174))
- format ci action
- change gorm-adaptor version to v3.5.0 ([#2247](https://github.com/CooooolFrog/dragonfly/issues/2247))
- add features swagger config ([#2246](https://github.com/CooooolFrog/dragonfly/issues/2246))
- release v2.0.9 and generate changelog ([#2181](https://github.com/CooooolFrog/dragonfly/issues/2181))
- add build-man-page to makefile ([#2182](https://github.com/CooooolFrog/dragonfly/issues/2182))
- change the compatibility testing version of manager and scheduler to v2.0.9 ([#2184](https://github.com/CooooolFrog/dragonfly/issues/2184))
- change docker.yml timout
- update nydus-snapshotter helm-charts to v0.0.4 ([#2188](https://github.com/CooooolFrog/dragonfly/issues/2188))
- update traffic shaper log ([#2227](https://github.com/CooooolFrog/dragonfly/issues/2227))
- update dingtalk group qrcode ([#2262](https://github.com/CooooolFrog/dragonfly/issues/2262))
- **deps:** bump github.com/swaggo/gin-swagger from 1.5.3 to 1.6.0 ([#2256](https://github.com/CooooolFrog/dragonfly/issues/2256))
- **deps:** bump github.com/swaggo/swag from 1.8.9 to 1.8.10 ([#2197](https://github.com/CooooolFrog/dragonfly/issues/2197))
- **deps:** bump github.com/aliyun/aliyun-oss-go-sdk from 2.2.6+incompatible to 2.2.7+incompatible ([#2218](https://github.com/CooooolFrog/dragonfly/issues/2218))
- **deps:** bump github.com/casbin/casbin/v2 from 2.66.1 to 2.66.3 ([#2260](https://github.com/CooooolFrog/dragonfly/issues/2260))
- **deps:** bump github.com/go-playground/validator/v10 from 10.11.2 to 10.12.0 ([#2220](https://github.com/CooooolFrog/dragonfly/issues/2220))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.224 to 1.44.229 ([#2221](https://github.com/CooooolFrog/dragonfly/issues/2221))
- **deps:** bump actions/setup-go from 3 to 4 ([#2202](https://github.com/CooooolFrog/dragonfly/issues/2202))
- **deps:** bump google.golang.org/api from 0.109.0 to 0.114.0 ([#2201](https://github.com/CooooolFrog/dragonfly/issues/2201))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.209 to 1.44.224 ([#2200](https://github.com/CooooolFrog/dragonfly/issues/2200))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.13.0 to 3.13.1 ([#2199](https://github.com/CooooolFrog/dragonfly/issues/2199))
- **deps:** bump golang.org/x/oauth2 from 0.5.0 to 0.6.0 ([#2198](https://github.com/CooooolFrog/dragonfly/issues/2198))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.37.0 to 0.40.0 ([#2219](https://github.com/CooooolFrog/dragonfly/issues/2219))
- **deps:** bump google.golang.org/protobuf from 1.29.0 to 1.29.1 ([#2195](https://github.com/CooooolFrog/dragonfly/issues/2195))
- **deps:** bump github.com/casbin/gorm-adapter/v3 from 3.5.0 to 3.15.0 ([#2237](https://github.com/CooooolFrog/dragonfly/issues/2237))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.2 to 3.23.3 ([#2239](https://github.com/CooooolFrog/dragonfly/issues/2239))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.229 to 1.44.234 ([#2240](https://github.com/CooooolFrog/dragonfly/issues/2240))
- **deps:** bump github.com/casbin/casbin/v2 from 2.65.2 to 2.66.1 ([#2238](https://github.com/CooooolFrog/dragonfly/issues/2238))
- **deps:** bump github.com/gin-gonic/gin from 1.8.2 to 1.9.0 ([#2241](https://github.com/CooooolFrog/dragonfly/issues/2241))
- **deps:** bump gorm.io/driver/postgres from 1.4.8 to 1.5.0 ([#2217](https://github.com/CooooolFrog/dragonfly/issues/2217))
- **deps:** bump moul.io/zapgorm2 from 1.2.0 to 1.3.0 ([#2167](https://github.com/CooooolFrog/dragonfly/issues/2167))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.23.1 to 3.23.2 ([#2166](https://github.com/CooooolFrog/dragonfly/issues/2166))
- **deps:** bump github.com/casbin/casbin/v2 from 2.64.0 to 2.65.2 ([#2164](https://github.com/CooooolFrog/dragonfly/issues/2164))
- **deps:** bump golang.org/x/crypto from 0.6.0 to 0.7.0 ([#2163](https://github.com/CooooolFrog/dragonfly/issues/2163))

### Docs
- optimize Community description in README.md ([#2255](https://github.com/CooooolFrog/dragonfly/issues/2255))
- add Volcano Engine to ADOPTERS.md ([#2169](https://github.com/CooooolFrog/dragonfly/issues/2169))

### Feat
- add multiple aarch support ([#2308](https://github.com/CooooolFrog/dragonfly/issues/2308))
- local dynconfig notifies data in client
- local dynconfig notifies data in client ([#2264](https://github.com/CooooolFrog/dragonfly/issues/2264))
- update resource director ([#2243](https://github.com/CooooolFrog/dragonfly/issues/2243))
- add CreatedAt function ([#2244](https://github.com/CooooolFrog/dragonfly/issues/2244))
- add trainer configuration ([#2216](https://github.com/CooooolFrog/dragonfly/issues/2216))
- update d7y.io/api package and change cpu percent validation ([#2236](https://github.com/CooooolFrog/dragonfly/issues/2236))
- add authinfo injector ([#2149](https://github.com/CooooolFrog/dragonfly/issues/2149))
- when the cache is missing, change the error log to a warning log ([#2235](https://github.com/CooooolFrog/dragonfly/issues/2235))
-  if the scheduler feature is not in feature flags, then it will stop providing the featrue ([#2234](https://github.com/CooooolFrog/dragonfly/issues/2234))
- add train interval and trainer addresses ([#2223](https://github.com/CooooolFrog/dragonfly/issues/2223))
- add logger.CoreLogger to searcher plugin ([#2232](https://github.com/CooooolFrog/dragonfly/issues/2232))
- add log to searcher plugin ([#2231](https://github.com/CooooolFrog/dragonfly/issues/2231))
- add probes struct ([#2190](https://github.com/CooooolFrog/dragonfly/issues/2190))
- add trainer config in scheduler ([#2214](https://github.com/CooooolFrog/dragonfly/issues/2214))
- add tfserving service to rpc package ([#2210](https://github.com/CooooolFrog/dragonfly/issues/2210))
- add trainer service to rpc package ([#2209](https://github.com/CooooolFrog/dragonfly/issues/2209))
- rename security client file name ([#2208](https://github.com/CooooolFrog/dragonfly/issues/2208))
- add CreateModel func to manager grpc client ([#2207](https://github.com/CooooolFrog/dragonfly/issues/2207))
- rename SecurityService to Security ([#2206](https://github.com/CooooolFrog/dragonfly/issues/2206))
- rename HostName to Hostname ([#2205](https://github.com/CooooolFrog/dragonfly/issues/2205))
- remove model migration ([#2204](https://github.com/CooooolFrog/dragonfly/issues/2204))
- change default value of dynconfig cache ([#2203](https://github.com/CooooolFrog/dragonfly/issues/2203))
- add index uk_model to model table ([#2196](https://github.com/CooooolFrog/dragonfly/issues/2196))
- remove model api ([#2194](https://github.com/CooooolFrog/dragonfly/issues/2194))
- add inference model table in database ([#2192](https://github.com/CooooolFrog/dragonfly/issues/2192))
- rename manager/model to manager/models ([#2191](https://github.com/CooooolFrog/dragonfly/issues/2191))
- add advertisePort to manager ([#2189](https://github.com/CooooolFrog/dragonfly/issues/2189))
- add advertise port ([#2156](https://github.com/CooooolFrog/dragonfly/issues/2156))
- add error log to database in manager ([#2172](https://github.com/CooooolFrog/dragonfly/issues/2172))

### Fix
- goreleaser config in v2.0.10 ([#2935](https://github.com/CooooolFrog/dragonfly/issues/2935))
- upgraded high vulnerability package to latest in v2.0.9 release ([#2928](https://github.com/CooooolFrog/dragonfly/issues/2928))
- object downloads failed by dfstore when dfdaemon enabled concurrent ([#2328](https://github.com/CooooolFrog/dragonfly/issues/2328))
- redis validation in scheduler config ([#2287](https://github.com/CooooolFrog/dragonfly/issues/2287))
- local dynconfig panic in Notify ([#2266](https://github.com/CooooolFrog/dragonfly/issues/2266))
- client grpc dial non-block ([#2261](https://github.com/CooooolFrog/dragonfly/issues/2261))
- modify the traversal condition for Items ([#2250](https://github.com/CooooolFrog/dragonfly/issues/2250))
- ip and hostname params in FindSchedulerClusters ([#2249](https://github.com/CooooolFrog/dragonfly/issues/2249))
- traffic shaper record task not found ([#2226](https://github.com/CooooolFrog/dragonfly/issues/2226))
- fsm events failed when register task ([#2225](https://github.com/CooooolFrog/dragonfly/issues/2225))
- stat DownloadPeerCount and DownloadPieceCount ([#2180](https://github.com/CooooolFrog/dragonfly/issues/2180))
- manager metrics Subsystem ([#2175](https://github.com/CooooolFrog/dragonfly/issues/2175))

### Refactor
- improve the performance of the code ([#2162](https://github.com/CooooolFrog/dragonfly/issues/2162))

### Test
- add new metrics test to service ([#2212](https://github.com/CooooolFrog/dragonfly/issues/2212))
- improve Test_parseByte ([#2173](https://github.com/CooooolFrog/dragonfly/issues/2173))
- add UT for byte String function ([#2170](https://github.com/CooooolFrog/dragonfly/issues/2170))
- improve TestMin ([#2168](https://github.com/CooooolFrog/dragonfly/issues/2168))


<a name="v2.0.9-rc.2"></a>
## [v2.0.9-rc.2] - 2023-03-13
### Chore
- add Garen Wen to maintainers ([#2136](https://github.com/CooooolFrog/dragonfly/issues/2136))
- **deps:** bump gorm.io/gorm from 1.24.5 to 1.24.6 ([#2143](https://github.com/CooooolFrog/dragonfly/issues/2143))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.13.0 to 1.14.0 ([#2144](https://github.com/CooooolFrog/dragonfly/issues/2144))
- **deps:** bump gorm.io/driver/postgres from 1.4.6 to 1.4.8 ([#2142](https://github.com/CooooolFrog/dragonfly/issues/2142))
- **deps:** bump github.com/stretchr/testify from 1.8.1 to 1.8.2 ([#2140](https://github.com/CooooolFrog/dragonfly/issues/2140))
- **deps:** bump github.com/casbin/casbin/v2 from 2.61.1 to 2.64.0 ([#2123](https://github.com/CooooolFrog/dragonfly/issues/2123))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.189 to 1.44.209 ([#2122](https://github.com/CooooolFrog/dragonfly/issues/2122))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.37.0 to 0.39.0 ([#2120](https://github.com/CooooolFrog/dragonfly/issues/2120))

### Docs
- add OpenSSF badge to README.md ([#2138](https://github.com/CooooolFrog/dragonfly/issues/2138))
- add public cloud providers Adopters.md ([#2137](https://github.com/CooooolFrog/dragonfly/issues/2137))

### Feat
- add auth config to manager ([#2161](https://github.com/CooooolFrog/dragonfly/issues/2161))
- add metrics to service v2 ([#2153](https://github.com/CooooolFrog/dragonfly/issues/2153))
- add SearchSchedulerClusterCount metric to manager ([#2152](https://github.com/CooooolFrog/dragonfly/issues/2152))
- implement announce peer ([#2150](https://github.com/CooooolFrog/dragonfly/issues/2150))
- add handleRegisterSeedPeerRequest to service v2 in scheduler ([#2148](https://github.com/CooooolFrog/dragonfly/issues/2148))
- add handleRegisterSeedPeerRequest to AnnouncePeer in service v2 ([#2147](https://github.com/CooooolFrog/dragonfly/issues/2147))
- change ScheduleCandidateParentsForNormalPeer implement ([#2133](https://github.com/CooooolFrog/dragonfly/issues/2133))
- enhance daemon health check ([#2130](https://github.com/CooooolFrog/dragonfly/issues/2130))
- implement v2 version of scheduler service ([#2125](https://github.com/CooooolFrog/dragonfly/issues/2125))

### Fix
- remove unnecessary fmt.Sprintf calls ([#2159](https://github.com/CooooolFrog/dragonfly/issues/2159))

### Refactor
- optimize certifyCache Get function ([#2160](https://github.com/CooooolFrog/dragonfly/issues/2160))

### Test
- add UT for MustParseRang ([#2158](https://github.com/CooooolFrog/dragonfly/issues/2158))
- improve TestFilterQuery ([#2157](https://github.com/CooooolFrog/dragonfly/issues/2157))
- add Validate test to scheduler config ([#2129](https://github.com/CooooolFrog/dragonfly/issues/2129))
- add Validate test to manager config ([#2128](https://github.com/CooooolFrog/dragonfly/issues/2128))
- refactor client validate ut ([#2126](https://github.com/CooooolFrog/dragonfly/issues/2126))
- add unit tests for DaemonConfig.Validate() ([#2119](https://github.com/CooooolFrog/dragonfly/issues/2119))


<a name="v2.0.9-rc.1"></a>
## [v2.0.9-rc.1] - 2023-02-27
### Feat
- update golang version to 1.20.1 ([#2117](https://github.com/CooooolFrog/dragonfly/issues/2117))
- correct grpc error code and implement StatPeer and LeavePeer ([#2115](https://github.com/CooooolFrog/dragonfly/issues/2115))

### Fix
- validate daemon gcInterval config ([#2118](https://github.com/CooooolFrog/dragonfly/issues/2118))


<a name="v2.0.9-rc.0"></a>
## [v2.0.9-rc.0] - 2023-02-24
### Feat
- add SyncNetworkTopology and SyncProbes to scheduler client ([#2114](https://github.com/CooooolFrog/dragonfly/issues/2114))

### Refactor
- preheat job ([#2113](https://github.com/CooooolFrog/dragonfly/issues/2113))


<a name="v2.0.9-beta.4"></a>
## [v2.0.9-beta.4] - 2023-02-23
### Feat
- add CIDR affinity to searcher ([#2111](https://github.com/CooooolFrog/dragonfly/issues/2111))


<a name="v2.0.9-beta.3"></a>
## [v2.0.9-beta.3] - 2023-02-22
### Chore
- remove unused MarkInvalid in daemon ([#2101](https://github.com/CooooolFrog/dragonfly/issues/2101))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.4.0 to 1.5.0 ([#2097](https://github.com/CooooolFrog/dragonfly/issues/2097))
- **deps:** bump gorm.io/driver/mysql from 1.4.5 to 1.4.7 ([#2096](https://github.com/CooooolFrog/dragonfly/issues/2096))
- **deps:** bump golang.org/x/oauth2 from 0.4.0 to 0.5.0 ([#2094](https://github.com/CooooolFrog/dragonfly/issues/2094))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.12.0 to 1.13.0 ([#2093](https://github.com/CooooolFrog/dragonfly/issues/2093))

### Feat
- remove Scopes and SecurityGroup in seed peer cluster ([#2110](https://github.com/CooooolFrog/dragonfly/issues/2110))
- dynconfig resolves addresses with host ([#2109](https://github.com/CooooolFrog/dragonfly/issues/2109))
- enable oss client download object concurrently. ([#2105](https://github.com/CooooolFrog/dragonfly/issues/2105))
- support reload scheduler addresses for local Dynconfig ([#2091](https://github.com/CooooolFrog/dragonfly/issues/2091))
- oss client supports STS access (set security token in header) ([#2103](https://github.com/CooooolFrog/dragonfly/issues/2103))
- don't GC task if expire time is 0 ([#2102](https://github.com/CooooolFrog/dragonfly/issues/2102))
- avoid checking dir existence before MkdirAll ([#2090](https://github.com/CooooolFrog/dragonfly/issues/2090))
- add host ttl to scheduler ([#2089](https://github.com/CooooolFrog/dragonfly/issues/2089))
- rename scheduler package to scheduling ([#2087](https://github.com/CooooolFrog/dragonfly/issues/2087))
- use v2 version of host id and add Addrs func to seed peer ([#2086](https://github.com/CooooolFrog/dragonfly/issues/2086))

### Fix
- unregister task from scheduler in storage.deleteTask ([#2100](https://github.com/CooooolFrog/dragonfly/issues/2100))

### Refactor
- support reload scheduler addresses for local Dynconfig in client ([#2107](https://github.com/CooooolFrog/dragonfly/issues/2107))
- scheduling with v2 grpc ([#2104](https://github.com/CooooolFrog/dragonfly/issues/2104))
- package digest ([#2085](https://github.com/CooooolFrog/dragonfly/issues/2085))
- type of digest in task ([#2084](https://github.com/CooooolFrog/dragonfly/issues/2084))

### Test
- remove random test in pieceDispatcherTest ([#2106](https://github.com/CooooolFrog/dragonfly/issues/2106))


<a name="v2.0.9-beta.2"></a>
## [v2.0.9-beta.2] - 2023-02-15
### Chore
- **deps:** bump github.com/go-playground/validator/v10 from 10.11.1 to 10.11.2 ([#2077](https://github.com/CooooolFrog/dragonfly/issues/2077))
- **deps:** bump github.com/casbin/casbin/v2 from 2.60.0 to 2.61.1 ([#2075](https://github.com/CooooolFrog/dragonfly/issues/2075))
- **deps:** bump go.opentelemetry.io/otel from 1.12.0 to 1.13.0 ([#2074](https://github.com/CooooolFrog/dragonfly/issues/2074))
- **deps:** bump github.com/looplab/fsm from 1.0.0 to 1.0.1 ([#2073](https://github.com/CooooolFrog/dragonfly/issues/2073))

### Feat
- add networkTopology configuration to scheduler ([#2070](https://github.com/CooooolFrog/dragonfly/issues/2070))
- remove training configuration in scheduler ([#2081](https://github.com/CooooolFrog/dragonfly/issues/2081))
- change piece size to length ([#2079](https://github.com/CooooolFrog/dragonfly/issues/2079))

### Fix
- backsource first piece timeout ([#2083](https://github.com/CooooolFrog/dragonfly/issues/2083))

### Refactor
- task.SizeScope with v2 grpc in scheduler ([#2082](https://github.com/CooooolFrog/dragonfly/issues/2082))
- task piece with v2 grpc ([#2080](https://github.com/CooooolFrog/dragonfly/issues/2080))


<a name="v2.0.9-beta.1"></a>
## [v2.0.9-beta.1] - 2023-02-14
### Chore
- change e2e timeout ([#2062](https://github.com/CooooolFrog/dragonfly/issues/2062))
- add miHoYo to ADOPTERS.md ([#2054](https://github.com/CooooolFrog/dragonfly/issues/2054))
- ignore configs generate with docker compose ([#2034](https://github.com/CooooolFrog/dragonfly/issues/2034))
- change maintainers informations ([#2038](https://github.com/CooooolFrog/dragonfly/issues/2038))
- update issue templates ([#2041](https://github.com/CooooolFrog/dragonfly/issues/2041))
- **deps:** bump github.com/jarcoal/httpmock from 1.2.0 to 1.3.0 ([#2044](https://github.com/CooooolFrog/dragonfly/issues/2044))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.184 to 1.44.189 ([#2029](https://github.com/CooooolFrog/dragonfly/issues/2029))
- **deps:** bump gorm.io/gorm from 1.24.3 to 1.24.5 ([#2042](https://github.com/CooooolFrog/dragonfly/issues/2042))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.12 to 3.23.1 ([#2045](https://github.com/CooooolFrog/dragonfly/issues/2045))
- **deps:** bump docker/build-push-action from 3 to 4 ([#2047](https://github.com/CooooolFrog/dragonfly/issues/2047))
- **deps:** bump google.golang.org/grpc from 1.52.0 to 1.52.3 ([#2046](https://github.com/CooooolFrog/dragonfly/issues/2046))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.11.2 to 1.12.0 ([#2030](https://github.com/CooooolFrog/dragonfly/issues/2030))
- **deps:** bump google.golang.org/api from 0.107.0 to 0.109.0 ([#2043](https://github.com/CooooolFrog/dragonfly/issues/2043))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.7.0 to 2.7.1 ([#2028](https://github.com/CooooolFrog/dragonfly/issues/2028))
- **deps:** bump github.com/onsi/gomega from 1.25.0 to 1.26.0 ([#2024](https://github.com/CooooolFrog/dragonfly/issues/2024))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.180 to 1.44.184 ([#2022](https://github.com/CooooolFrog/dragonfly/issues/2022))
- **deps:** bump github.com/onsi/gomega from 1.24.2 to 1.25.0 ([#2021](https://github.com/CooooolFrog/dragonfly/issues/2021))
- **deps:** bump github.com/montanaflynn/stats from 0.6.6 to 0.7.0 ([#2020](https://github.com/CooooolFrog/dragonfly/issues/2020))
- **deps:** bump github.com/spf13/viper from 1.13.0 to 1.15.0 ([#2019](https://github.com/CooooolFrog/dragonfly/issues/2019))
- **deps:** bump gorm.io/gorm from 1.24.2 to 1.24.3 ([#2018](https://github.com/CooooolFrog/dragonfly/issues/2018))

### Docs
- change introduction in readem ([#2017](https://github.com/CooooolFrog/dragonfly/issues/2017))

### Feat
- set gorm log level ([#2063](https://github.com/CooooolFrog/dragonfly/issues/2063))
- change PeerCountLimitForTask to 1000 ([#2059](https://github.com/CooooolFrog/dragonfly/issues/2059))
- add v2 version of the idgen ([#2056](https://github.com/CooooolFrog/dragonfly/issues/2056))
- update task type from v1 to v2 ([#2053](https://github.com/CooooolFrog/dragonfly/issues/2053))
- add AnnouncePeers to task in resource ([#2051](https://github.com/CooooolFrog/dragonfly/issues/2051))
- add v2 version of dfdaemon client ([#2050](https://github.com/CooooolFrog/dragonfly/issues/2050))
- add DownloadTask to seed peer resource ([#2048](https://github.com/CooooolFrog/dragonfly/issues/2048))
- init AnnouncePeerStream of the peer ([#2040](https://github.com/CooooolFrog/dragonfly/issues/2040))
- update dingtalk qrcode ([#2016](https://github.com/CooooolFrog/dragonfly/issues/2016))

### Fix
- peer GC clear all peers when peer's count large than PeerCountLimitForTask ([#2061](https://github.com/CooooolFrog/dragonfly/issues/2061))
- spelling mistakes ([#2027](https://github.com/CooooolFrog/dragonfly/issues/2027))

### Refactor
- resource task with v2 version of grpc ([#2078](https://github.com/CooooolFrog/dragonfly/issues/2078))
- parse http range ([#2071](https://github.com/CooooolFrog/dragonfly/issues/2071))
- peer resource with v2 version of the grpc ([#2039](https://github.com/CooooolFrog/dragonfly/issues/2039))
- announcer and dynconfig with v2 verison of the manager grpc ([#2037](https://github.com/CooooolFrog/dragonfly/issues/2037))
- resource host without scheduler v1 definition ([#2036](https://github.com/CooooolFrog/dragonfly/issues/2036))


<a name="v2.0.9-beta.0"></a>
## [v2.0.9-beta.0] - 2023-01-19
### Chore
- fix workflows typo ([#2013](https://github.com/CooooolFrog/dragonfly/issues/2013))
- generate manager swagger ([#2009](https://github.com/CooooolFrog/dragonfly/issues/2009))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.175 to 1.44.180 ([#2005](https://github.com/CooooolFrog/dragonfly/issues/2005))
- **deps:** bump google.golang.org/api from 0.106.0 to 0.107.0 ([#2004](https://github.com/CooooolFrog/dragonfly/issues/2004))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.6.1 to 2.7.0 ([#2003](https://github.com/CooooolFrog/dragonfly/issues/2003))
- **deps:** bump gorm.io/driver/postgres from 1.4.5 to 1.4.6 ([#2002](https://github.com/CooooolFrog/dragonfly/issues/2002))

### Feat
- update helm charts ([#2015](https://github.com/CooooolFrog/dragonfly/issues/2015))
- add directed graph to pkg ([#2014](https://github.com/CooooolFrog/dragonfly/issues/2014))
- change peer's piece type in resource ([#2012](https://github.com/CooooolFrog/dragonfly/issues/2012))
- support source client option ([#2008](https://github.com/CooooolFrog/dragonfly/issues/2008))
- change ok to loaded in loading func ([#2010](https://github.com/CooooolFrog/dragonfly/issues/2010))
- remove NetTopology from scheduler and manager ([#2007](https://github.com/CooooolFrog/dragonfly/issues/2007))

### Fix
- dferror not convert ([#2001](https://github.com/CooooolFrog/dragonfly/issues/2001))
- dfstore typo ([#2000](https://github.com/CooooolFrog/dragonfly/issues/2000))

### Refactor
- piece_dispatcher considering score of parent peer ([#1978](https://github.com/CooooolFrog/dragonfly/issues/1978))


<a name="v2.0.9-alpha.10"></a>
## [v2.0.9-alpha.10] - 2023-01-13
### Chore
- update helm charts submodule ([#1997](https://github.com/CooooolFrog/dragonfly/issues/1997))

### Feat
- add v2 verison of the grpc to scheduler ([#1999](https://github.com/CooooolFrog/dragonfly/issues/1999))

### Fix
- manager typo ([#1995](https://github.com/CooooolFrog/dragonfly/issues/1995))
- daemon recognize Code_SchedForbidden ([#1994](https://github.com/CooooolFrog/dragonfly/issues/1994))
- count of total page in pagination ([#1993](https://github.com/CooooolFrog/dragonfly/issues/1993))
- manager grpc filename ([#1992](https://github.com/CooooolFrog/dragonfly/issues/1992))


<a name="v2.0.9-alpha.9"></a>
## [v2.0.9-alpha.9] - 2023-01-10
### Chore
- **deps:** bump github.com/schollz/progressbar/v3 from 3.12.2 to 3.13.0 ([#1989](https://github.com/CooooolFrog/dragonfly/issues/1989))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.171 to 1.44.175 ([#1988](https://github.com/CooooolFrog/dragonfly/issues/1988))
- **deps:** bump google.golang.org/api from 0.105.0 to 0.106.0 ([#1987](https://github.com/CooooolFrog/dragonfly/issues/1987))
- **deps:** bump golang.org/x/crypto from 0.4.0 to 0.5.0 ([#1986](https://github.com/CooooolFrog/dragonfly/issues/1986))
- **deps:** bump golang.org/x/time from 0.1.0 to 0.3.0 ([#1985](https://github.com/CooooolFrog/dragonfly/issues/1985))

### Docs
- fix manager swag error ([#1982](https://github.com/CooooolFrog/dragonfly/issues/1982))

### Feat
- add manager v2 api ([#1990](https://github.com/CooooolFrog/dragonfly/issues/1990))
- searcher can not found candidate scheduler clusters, return all scheduler clusters ([#1991](https://github.com/CooooolFrog/dragonfly/issues/1991))
- oras source client ([#1983](https://github.com/CooooolFrog/dragonfly/issues/1983))
- add fail_code in scheduler's DownloadFailureCount metric ([#1981](https://github.com/CooooolFrog/dragonfly/issues/1981))
- add udp ping to the net package ([#1979](https://github.com/CooooolFrog/dragonfly/issues/1979))

### Fix
- client bitMap extend capacity ([#1973](https://github.com/CooooolFrog/dragonfly/issues/1973))


<a name="v2.0.9-alpha.8"></a>
## [v2.0.9-alpha.8] - 2023-01-06
### Chore
- remove codecov patch feature ([#1977](https://github.com/CooooolFrog/dragonfly/issues/1977))
- update e2e timeout ([#1969](https://github.com/CooooolFrog/dragonfly/issues/1969))
- update charts version ([#1968](https://github.com/CooooolFrog/dragonfly/issues/1968))

### Feat
- add S3ForcePathStyle to object storage ([#1976](https://github.com/CooooolFrog/dragonfly/issues/1976))

### Fix
- context of trigger seed peer ([#1971](https://github.com/CooooolFrog/dragonfly/issues/1971))


<a name="v2.0.9-alpha.7"></a>
## [v2.0.9-alpha.7] - 2023-01-03
### Chore
- goreleaser set rlcp field ([#1967](https://github.com/CooooolFrog/dragonfly/issues/1967))
- update actions ([#1966](https://github.com/CooooolFrog/dragonfly/issues/1966))
- print e2e exec output ([#1963](https://github.com/CooooolFrog/dragonfly/issues/1963))
- change codecov coverage range ([#1965](https://github.com/CooooolFrog/dragonfly/issues/1965))
- create log dir ([#1947](https://github.com/CooooolFrog/dragonfly/issues/1947))
- change dingtalk image ([#1954](https://github.com/CooooolFrog/dragonfly/issues/1954))
- **deps:** bump moul.io/zapgorm2 from 1.1.3 to 1.2.0 ([#1961](https://github.com/CooooolFrog/dragonfly/issues/1961))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.167 to 1.44.171 ([#1958](https://github.com/CooooolFrog/dragonfly/issues/1958))
- **deps:** bump github.com/appleboy/gin-jwt/v2 from 2.9.0 to 2.9.1 ([#1949](https://github.com/CooooolFrog/dragonfly/issues/1949))
- **deps:** bump google.golang.org/api from 0.101.0 to 0.105.0 ([#1952](https://github.com/CooooolFrog/dragonfly/issues/1952))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.11 to 3.22.12 ([#1959](https://github.com/CooooolFrog/dragonfly/issues/1959))
- **deps:** bump github.com/gin-gonic/gin from 1.8.1 to 1.8.2 ([#1951](https://github.com/CooooolFrog/dragonfly/issues/1951))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.36.4 to 0.37.0 ([#1950](https://github.com/CooooolFrog/dragonfly/issues/1950))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.143 to 1.44.167 ([#1948](https://github.com/CooooolFrog/dragonfly/issues/1948))
- **deps:** bump gorm.io/driver/mysql from 1.4.4 to 1.4.5 ([#1962](https://github.com/CooooolFrog/dragonfly/issues/1962))

### Docs
- change dingtalk link

### Fix
- config decode net.IP ([#1964](https://github.com/CooooolFrog/dragonfly/issues/1964))


<a name="v2.0.9-alpha.6"></a>
## [v2.0.9-alpha.6] - 2022-12-22
### Chore
- optimize download log ([#1944](https://github.com/CooooolFrog/dragonfly/issues/1944))
- update helm charts version ([#1937](https://github.com/CooooolFrog/dragonfly/issues/1937))
- **deps:** bump k8s.io/component-base from 0.25.4 to 0.26.0 ([#1934](https://github.com/CooooolFrog/dragonfly/issues/1934))
- **deps:** bump goreleaser/goreleaser-action from 3 to 4 ([#1936](https://github.com/CooooolFrog/dragonfly/issues/1936))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.10 to 3.22.11 ([#1935](https://github.com/CooooolFrog/dragonfly/issues/1935))
- **deps:** bump github.com/swaggo/swag from 1.8.8 to 1.8.9 ([#1932](https://github.com/CooooolFrog/dragonfly/issues/1932))
- **deps:** bump github.com/onsi/gomega from 1.24.1 to 1.24.2 ([#1931](https://github.com/CooooolFrog/dragonfly/issues/1931))

### Feat
- corrupt data check ([#1946](https://github.com/CooooolFrog/dragonfly/issues/1946))
- create synchronizers concurrently ([#1941](https://github.com/CooooolFrog/dragonfly/issues/1941))
- register reflection on grpc server ([#1943](https://github.com/CooooolFrog/dragonfly/issues/1943))
- remove legacy peers support ([#1939](https://github.com/CooooolFrog/dragonfly/issues/1939))
- update fsm stable api ([#1938](https://github.com/CooooolFrog/dragonfly/issues/1938))
- add IPAddresses and DNSNames to sans of the cert ([#1930](https://github.com/CooooolFrog/dragonfly/issues/1930))
- change yaml field type from string to net.IP ([#1929](https://github.com/CooooolFrog/dragonfly/issues/1929))
- random pieces download ([#1918](https://github.com/CooooolFrog/dragonfly/issues/1918))

### Fix
- download context cancelled ([#1942](https://github.com/CooooolFrog/dragonfly/issues/1942))
- peer keepalive with manager ([#1940](https://github.com/CooooolFrog/dragonfly/issues/1940))
- panic caused by hashring not being built ([#1928](https://github.com/CooooolFrog/dragonfly/issues/1928))


<a name="v2.0.9-alpha.4"></a>
## [v2.0.9-alpha.4] - 2022-12-14
### Feat
- update version guage metrics ([#1927](https://github.com/CooooolFrog/dragonfly/issues/1927))


<a name="v2.0.9-alpha.3"></a>
## [v2.0.9-alpha.3] - 2022-12-14
### Chore
- add priority to dfget man page ([#1917](https://github.com/CooooolFrog/dragonfly/issues/1917))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.11.1 to 1.11.2 ([#1922](https://github.com/CooooolFrog/dragonfly/issues/1922))
- **deps:** bump github.com/casbin/casbin/v2 from 2.58.0 to 2.60.0 ([#1921](https://github.com/CooooolFrog/dragonfly/issues/1921))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.12.1 to 3.12.2 ([#1920](https://github.com/CooooolFrog/dragonfly/issues/1920))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.36.3 to 0.37.0 ([#1919](https://github.com/CooooolFrog/dragonfly/issues/1919))

### Feat
- remove callsystem and pattern ([#1925](https://github.com/CooooolFrog/dragonfly/issues/1925))
- client support 'priority' parameter ([#1911](https://github.com/CooooolFrog/dragonfly/issues/1911))

### Fix
- application not found ([#1924](https://github.com/CooooolFrog/dragonfly/issues/1924))

### Refactor
- dynconfig without Unmarshal ([#1926](https://github.com/CooooolFrog/dragonfly/issues/1926))


<a name="v2.0.9-alpha.2"></a>
## [v2.0.9-alpha.2] - 2022-12-12
### Feat
- handle application not found ([#1913](https://github.com/CooooolFrog/dragonfly/issues/1913))
- update priority api ([#1912](https://github.com/CooooolFrog/dragonfly/issues/1912))


<a name="v2.0.9-alpha.1"></a>
## [v2.0.9-alpha.1] - 2022-12-09
### Chore
- upload nydus e2e logs to artifact ([#1909](https://github.com/CooooolFrog/dragonfly/issues/1909))

### Feat
- support redis sentinal ([#1910](https://github.com/CooooolFrog/dragonfly/issues/1910))
- update submodule console ([#1908](https://github.com/CooooolFrog/dragonfly/issues/1908))


<a name="v2.0.9-alpha.0"></a>
## [v2.0.9-alpha.0] - 2022-12-09
### Chore
- add Baidu to ADOPTERS.md ([#1884](https://github.com/CooooolFrog/dragonfly/issues/1884))
- update api package verison ([#1893](https://github.com/CooooolFrog/dragonfly/issues/1893))
- optimize reregister ([#1888](https://github.com/CooooolFrog/dragonfly/issues/1888))
- releaser action disable cgo ([#1885](https://github.com/CooooolFrog/dragonfly/issues/1885))
- add e2e with nydus snapshotter ([#1860](https://github.com/CooooolFrog/dragonfly/issues/1860))
- release v2.0.8 ([#1877](https://github.com/CooooolFrog/dragonfly/issues/1877))
- **deps:** bump go.uber.org/zap from 1.23.0 to 1.24.0 ([#1900](https://github.com/CooooolFrog/dragonfly/issues/1900))
- **deps:** bump github.com/casbin/casbin/v2 from 2.56.0 to 2.58.0 ([#1899](https://github.com/CooooolFrog/dragonfly/issues/1899))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.11.0 to 3.12.1 ([#1898](https://github.com/CooooolFrog/dragonfly/issues/1898))
- **deps:** bump github.com/swaggo/swag from 1.8.7 to 1.8.8 ([#1897](https://github.com/CooooolFrog/dragonfly/issues/1897))
- **deps:** bump github.com/go-sql-driver/mysql from 1.6.0 to 1.7.0 ([#1896](https://github.com/CooooolFrog/dragonfly/issues/1896))
- **deps:** bump github.com/huaweicloud/huaweicloud-sdk-go-obs from 3.21.12+incompatible to 3.22.11+incompatible ([#1872](https://github.com/CooooolFrog/dragonfly/issues/1872))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.5.0 to 2.5.1 ([#1871](https://github.com/CooooolFrog/dragonfly/issues/1871))

### Feat
- storage collects upload piece count, peer cost and error details ([#1907](https://github.com/CooooolFrog/dragonfly/issues/1907))
- priority of the register parameter is higher than parameter of the application ([#1906](https://github.com/CooooolFrog/dragonfly/issues/1906))
- trigger task with priority ([#1904](https://github.com/CooooolFrog/dragonfly/issues/1904))
- dynconfig adds list application in scheduler ([#1903](https://github.com/CooooolFrog/dragonfly/issues/1903))
- rename url priority struct and remove PriorityLevel constants ([#1902](https://github.com/CooooolFrog/dragonfly/issues/1902))
- add priority to application in manager ([#1901](https://github.com/CooooolFrog/dragonfly/issues/1901))
- remove relation of application ([#1894](https://github.com/CooooolFrog/dragonfly/issues/1894))
- add backSourceCount validation ([#1892](https://github.com/CooooolFrog/dragonfly/issues/1892))
- add health check to service ([#1889](https://github.com/CooooolFrog/dragonfly/issues/1889))
- add pieceDownloadTimeout to docker compose template ([#1881](https://github.com/CooooolFrog/dragonfly/issues/1881))
- add the timeout of downloading piece to scheduler ([#1880](https://github.com/CooooolFrog/dragonfly/issues/1880))
- change log rotate size ([#1879](https://github.com/CooooolFrog/dragonfly/issues/1879))
- support reregister peer task in client ([#1876](https://github.com/CooooolFrog/dragonfly/issues/1876))
- if the scheduler cannot find the peer, then return Code_SchedReregister to dfdaemon ([#1875](https://github.com/CooooolFrog/dragonfly/issues/1875))

### Fix
- remove advertiseIP config in e2e ([#1878](https://github.com/CooooolFrog/dragonfly/issues/1878))

### Refactor
- back-to-source configuration ([#1895](https://github.com/CooooolFrog/dragonfly/issues/1895))


<a name="v2.0.8-rc.3"></a>
## [v2.0.8-rc.3] - 2022-11-25
### Feat
- change announcer validation ([#1869](https://github.com/CooooolFrog/dragonfly/issues/1869))
- add mysql read and write timeout ([#1868](https://github.com/CooooolFrog/dragonfly/issues/1868))


<a name="v2.0.8-rc.2"></a>
## [v2.0.8-rc.2] - 2022-11-25
### Chore
- add Kuaishou to ADOPTERS.md ([#1866](https://github.com/CooooolFrog/dragonfly/issues/1866))

### Feat
- store parent information ([#1867](https://github.com/CooooolFrog/dragonfly/issues/1867))
- remove MainParent from peer and add IsPieceBackToSource to piece

### Fix
- [scheduler]  destPeer keepalive when downloaded by other peer ([#1865](https://github.com/CooooolFrog/dragonfly/issues/1865))


<a name="v2.0.8-rc.1"></a>
## [v2.0.8-rc.1] - 2022-11-24
### Feat
- scheduler supports storage config ([#1864](https://github.com/CooooolFrog/dragonfly/issues/1864))
- store peer download information ([#1863](https://github.com/CooooolFrog/dragonfly/issues/1863))
- manager changes filterParentLimit value ([#1859](https://github.com/CooooolFrog/dragonfly/issues/1859))


<a name="v2.0.8-rc.0"></a>
## [v2.0.8-rc.0] - 2022-11-22
### Chore
- update dst peer log ([#1844](https://github.com/CooooolFrog/dragonfly/issues/1844))
- update e2e test ([#1839](https://github.com/CooooolFrog/dragonfly/issues/1839))
- remove unused code ([#1838](https://github.com/CooooolFrog/dragonfly/issues/1838))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.121 to 1.44.143 ([#1853](https://github.com/CooooolFrog/dragonfly/issues/1853))
- **deps:** bump github.com/prometheus/client_golang from 1.13.0 to 1.14.0 ([#1851](https://github.com/CooooolFrog/dragonfly/issues/1851))
- **deps:** bump github.com/aliyun/aliyun-oss-go-sdk from 2.2.5+incompatible to 2.2.6+incompatible ([#1849](https://github.com/CooooolFrog/dragonfly/issues/1849))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.3.3 to 1.4.0 ([#1848](https://github.com/CooooolFrog/dragonfly/issues/1848))
- **deps:** bump k8s.io/component-base from 0.25.3 to 0.25.4 ([#1847](https://github.com/CooooolFrog/dragonfly/issues/1847))
- **deps:** bump github.com/onsi/gomega from 1.23.0 to 1.24.1 ([#1832](https://github.com/CooooolFrog/dragonfly/issues/1832))
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.4.0 to 2.5.0 ([#1836](https://github.com/CooooolFrog/dragonfly/issues/1836))
- **deps:** bump github.com/mdlayher/vsock from 1.1.1 to 1.2.0 ([#1834](https://github.com/CooooolFrog/dragonfly/issues/1834))

### Feat
- optimize gc package ([#1855](https://github.com/CooooolFrog/dragonfly/issues/1855))
- add announcer to scheduler ([#1854](https://github.com/CooooolFrog/dragonfly/issues/1854))
- add announcer to dfdameon ([#1852](https://github.com/CooooolFrog/dragonfly/issues/1852))
- when dfdaemon disable object storage, dynconfig can't fetch manager ([#1845](https://github.com/CooooolFrog/dragonfly/issues/1845))
- optimize manager log ([#1846](https://github.com/CooooolFrog/dragonfly/issues/1846))
- scheduler adds announce host handler ([#1843](https://github.com/CooooolFrog/dragonfly/issues/1843))
- call all nodes in consistent hashing and reuse grpc connection ([#1842](https://github.com/CooooolFrog/dragonfly/issues/1842))
- update concurrent-map version ([#1837](https://github.com/CooooolFrog/dragonfly/issues/1837))

### Fix
- recursive download always return none error ([#1841](https://github.com/CooooolFrog/dragonfly/issues/1841))
- expire header timezone ([#1840](https://github.com/CooooolFrog/dragonfly/issues/1840))


<a name="v2.0.8-beta.3"></a>
## [v2.0.8-beta.3] - 2022-11-14
### Chore
- enable cache list metadata e2e ([#1829](https://github.com/CooooolFrog/dragonfly/issues/1829))

### Feat
- optimize scope size is error ([#1831](https://github.com/CooooolFrog/dragonfly/issues/1831))
- add timeout grpc and job ([#1830](https://github.com/CooooolFrog/dragonfly/issues/1830))


<a name="v2.0.8-beta.2"></a>
## [v2.0.8-beta.2] - 2022-11-14
### Feat
- optimize peer log ([#1828](https://github.com/CooooolFrog/dragonfly/issues/1828))
- optional save list metadata to p2p ([#1822](https://github.com/CooooolFrog/dragonfly/issues/1822))
- add s3 resource client and recursive e2e test ([#1826](https://github.com/CooooolFrog/dragonfly/issues/1826))


<a name="v2.0.8-beta.1"></a>
## [v2.0.8-beta.1] - 2022-11-11
### Chore
- daemon avoid alway open metadata files ([#1823](https://github.com/CooooolFrog/dragonfly/issues/1823))

### Feat
- optimize preheat log ([#1827](https://github.com/CooooolFrog/dragonfly/issues/1827))
- seed peer reuses traffic ([#1825](https://github.com/CooooolFrog/dragonfly/issues/1825))
- optimize preheat ([#1824](https://github.com/CooooolFrog/dragonfly/issues/1824))


<a name="v2.0.8-beta.0"></a>
## [v2.0.8-beta.0] - 2022-11-09
### Chore
- close out of use client grpc conn ([#1817](https://github.com/CooooolFrog/dragonfly/issues/1817))
- make SendMsg in doRecursiveDownload safe ([#1806](https://github.com/CooooolFrog/dragonfly/issues/1806))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.10.0 to 1.11.1 ([#1813](https://github.com/CooooolFrog/dragonfly/issues/1813))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.9 to 3.22.10 ([#1812](https://github.com/CooooolFrog/dragonfly/issues/1812))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.36.1 to 0.36.4 ([#1811](https://github.com/CooooolFrog/dragonfly/issues/1811))
- **deps:** bump github.com/gammazero/deque from 0.2.0 to 0.2.1 ([#1810](https://github.com/CooooolFrog/dragonfly/issues/1810))

### Feat
- returns an scheduling error if the peer state is not PeerStateRunning ([#1821](https://github.com/CooooolFrog/dragonfly/issues/1821))
- optimize peer gc ([#1819](https://github.com/CooooolFrog/dragonfly/issues/1819))
- peer.UpdateAt needs to be updated during download process ([#1818](https://github.com/CooooolFrog/dragonfly/issues/1818))
- statistical the traffic of reused peer ([#1816](https://github.com/CooooolFrog/dragonfly/issues/1816))
- add workHome and pluginDir to configuration ([#1807](https://github.com/CooooolFrog/dragonfly/issues/1807))
- add otel trace in log ([#1804](https://github.com/CooooolFrog/dragonfly/issues/1804))

### Fix
- otel goroutine leak ([#1815](https://github.com/CooooolFrog/dragonfly/issues/1815))


<a name="v2.0.8-alpha.3"></a>
## [v2.0.8-alpha.3] - 2022-11-02
### Fix
- gorm-adaptor pkg version ([#1805](https://github.com/CooooolFrog/dragonfly/issues/1805))


<a name="v2.0.8-alpha.2"></a>
## [v2.0.8-alpha.2] - 2022-11-02
### Chore
- add list log in rpc download ([#1802](https://github.com/CooooolFrog/dragonfly/issues/1802))


<a name="v2.0.8-alpha.1"></a>
## [v2.0.8-alpha.1] - 2022-11-01
### Chore
- **deps:** bump gorm.io/driver/mysql from 1.4.1 to 1.4.3 ([#1799](https://github.com/CooooolFrog/dragonfly/issues/1799))
- **deps:** bump google.golang.org/api from 0.97.0 to 0.101.0 ([#1800](https://github.com/CooooolFrog/dragonfly/issues/1800))
- **deps:** bump github.com/onsi/gomega from 1.22.1 to 1.23.0 ([#1798](https://github.com/CooooolFrog/dragonfly/issues/1798))
- **deps:** bump gorm.io/driver/postgres from 1.4.4 to 1.4.5 ([#1797](https://github.com/CooooolFrog/dragonfly/issues/1797))
- **deps:** bump github.com/spf13/cobra from 1.5.0 to 1.6.1 ([#1796](https://github.com/CooooolFrog/dragonfly/issues/1796))

### Feat
- add leave host logger ([#1801](https://github.com/CooooolFrog/dragonfly/issues/1801))
- scheduler's record adds ParentUploadCount and ParentUploadFailedCount ([#1795](https://github.com/CooooolFrog/dragonfly/issues/1795))

### Fix
- leave host ([#1803](https://github.com/CooooolFrog/dragonfly/issues/1803))


<a name="v2.0.8-alpha.0"></a>
## [v2.0.8-alpha.0] - 2022-10-28
### Chore
- add timestamp to stdout&stderr ([#1781](https://github.com/CooooolFrog/dragonfly/issues/1781))
- update grpc api proto verison ([#1779](https://github.com/CooooolFrog/dragonfly/issues/1779))
- add intel to ADOPTERS.md ([#1778](https://github.com/CooooolFrog/dragonfly/issues/1778))
- update helm-charts submodule
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.2.0 to 2.4.0 ([#1787](https://github.com/CooooolFrog/dragonfly/issues/1787))
- **deps:** bump github.com/casbin/gorm-adapter/v3 from 3.11.0 to 3.12.1 ([#1786](https://github.com/CooooolFrog/dragonfly/issues/1786))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.114 to 1.44.121 ([#1785](https://github.com/CooooolFrog/dragonfly/issues/1785))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.7 to 3.22.9 ([#1784](https://github.com/CooooolFrog/dragonfly/issues/1784))
- **deps:** bump go.opentelemetry.io/otel from 1.11.0 to 1.11.1 ([#1783](https://github.com/CooooolFrog/dragonfly/issues/1783))

### Feat
- support split running tasks ([#1794](https://github.com/CooooolFrog/dragonfly/issues/1794))
- add download header log ([#1793](https://github.com/CooooolFrog/dragonfly/issues/1793))
- grpc scheduler client dial options ([#1792](https://github.com/CooooolFrog/dragonfly/issues/1792))
- daemon call leaveHost when exit ([#1788](https://github.com/CooooolFrog/dragonfly/issues/1788))
- add calculateParentHostUploadSuccessScore to scheduler ([#1789](https://github.com/CooooolFrog/dragonfly/issues/1789))
- add LeaveHost handler ([#1780](https://github.com/CooooolFrog/dragonfly/issues/1780))

### Fix
- daemon don't leaveHost when keepStorage=true ([#1790](https://github.com/CooooolFrog/dragonfly/issues/1790))
- did not call scheduler leave tasks in forceGC ([#1782](https://github.com/CooooolFrog/dragonfly/issues/1782))


<a name="v2.0.7"></a>
## [v2.0.7] - 2022-10-19
### Chore
- release v2.0.7 ([#1776](https://github.com/CooooolFrog/dragonfly/issues/1776))

### Fix
- plugin builder ([#1775](https://github.com/CooooolFrog/dragonfly/issues/1775))


<a name="v2.0.7-rc.0"></a>
## [v2.0.7-rc.0] - 2022-10-18
### Chore
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.36.1 to 0.36.3 ([#1768](https://github.com/CooooolFrog/dragonfly/issues/1768))
- **deps:** bump github.com/swaggo/swag from 1.8.5 to 1.8.7 ([#1773](https://github.com/CooooolFrog/dragonfly/issues/1773))
- **deps:** bump k8s.io/component-base from 0.25.2 to 0.25.3 ([#1771](https://github.com/CooooolFrog/dragonfly/issues/1771))
- **deps:** bump github.com/swaggo/swag from 1.8.5 to 1.8.6 ([#1770](https://github.com/CooooolFrog/dragonfly/issues/1770))
- **deps:** bump github.com/casbin/casbin/v2 from 2.55.1 to 2.56.0 ([#1769](https://github.com/CooooolFrog/dragonfly/issues/1769))
- **deps:** bump go.opentelemetry.io/otel/trace from 1.10.0 to 1.11.0 ([#1767](https://github.com/CooooolFrog/dragonfly/issues/1767))

### Fix
- add fallback registry mirror in gen-containerd-host.sh ([#1774](https://github.com/CooooolFrog/dragonfly/issues/1774))


<a name="v2.0.7-beta.7"></a>
## [v2.0.7-beta.7] - 2022-10-17
### Chore
- check reuse file ([#1765](https://github.com/CooooolFrog/dragonfly/issues/1765))
- update golang version to 1.19 ([#1760](https://github.com/CooooolFrog/dragonfly/issues/1760))

### Docs
- add adopters.md ([#1759](https://github.com/CooooolFrog/dragonfly/issues/1759))

### Feat
- grpc_retry removes WithPerRetryTimeout ([#1763](https://github.com/CooooolFrog/dragonfly/issues/1763))
- obs object storage support ([#1758](https://github.com/CooooolFrog/dragonfly/issues/1758))

### Fix
- open end range in concurrent back source ([#1764](https://github.com/CooooolFrog/dragonfly/issues/1764))
- manager PeerGauge ([#1761](https://github.com/CooooolFrog/dragonfly/issues/1761))

### Refactor
- scheduler registers task ([#1766](https://github.com/CooooolFrog/dragonfly/issues/1766))
- obs of objectstorage pkg ([#1762](https://github.com/CooooolFrog/dragonfly/issues/1762))


<a name="v2.0.7-beta.6"></a>
## [v2.0.7-beta.6] - 2022-10-13
### Chore
- update console submodule ([#1755](https://github.com/CooooolFrog/dragonfly/issues/1755))
- update roundtrip log ([#1750](https://github.com/CooooolFrog/dragonfly/issues/1750))
- update console submodule ([#1748](https://github.com/CooooolFrog/dragonfly/issues/1748))
- change docker compose task ttl ([#1741](https://github.com/CooooolFrog/dragonfly/issues/1741))
- **deps:** bump github.com/casbin/gorm-adapter/v3 from 3.5.0 to 3.11.0 ([#1745](https://github.com/CooooolFrog/dragonfly/issues/1745))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.8.7 to 3.11.0 ([#1746](https://github.com/CooooolFrog/dragonfly/issues/1746))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.34.0 to 0.36.1 ([#1744](https://github.com/CooooolFrog/dragonfly/issues/1744))
- **deps:** bump gorm.io/driver/postgres from 1.3.10 to 1.4.4 ([#1743](https://github.com/CooooolFrog/dragonfly/issues/1743))
- **deps:** bump google.golang.org/grpc from 1.49.0 to 1.50.0 ([#1742](https://github.com/CooooolFrog/dragonfly/issues/1742))

### Feat
- available peer includes state is PeerStatePending ([#1756](https://github.com/CooooolFrog/dragonfly/issues/1756))
- peer will back-to-source when task switch state failed ([#1754](https://github.com/CooooolFrog/dragonfly/issues/1754))
- change FilterParentRangeLimit validation ([#1752](https://github.com/CooooolFrog/dragonfly/issues/1752))
- task state is TaskStateRunning can be registered ([#1751](https://github.com/CooooolFrog/dragonfly/issues/1751))
- gin logger rotation ([#1749](https://github.com/CooooolFrog/dragonfly/issues/1749))


<a name="v2.0.7-beta.5"></a>
## [v2.0.7-beta.5] - 2022-10-10
### Chore
- make lru cache safe ([#1737](https://github.com/CooooolFrog/dragonfly/issues/1737))

### Feat
- overwrite task url and url meta ([#1740](https://github.com/CooooolFrog/dragonfly/issues/1740))
- update source temporary error logic ([#1739](https://github.com/CooooolFrog/dragonfly/issues/1739))
- add seed peer back source traffic ([#1738](https://github.com/CooooolFrog/dragonfly/issues/1738))
- http request content log ([#1736](https://github.com/CooooolFrog/dragonfly/issues/1736))
- remove task and host gc ttl ([#1735](https://github.com/CooooolFrog/dragonfly/issues/1735))
- add http request log ([#1734](https://github.com/CooooolFrog/dragonfly/issues/1734))

### Fix
- backsource temporary error judgement ([#1726](https://github.com/CooooolFrog/dragonfly/issues/1726))


<a name="v2.0.7-beta.4"></a>
## [v2.0.7-beta.4] - 2022-10-09
### Chore
- change disk usage debug log format to decimal ([#1727](https://github.com/CooooolFrog/dragonfly/issues/1727))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.95 to 1.44.114 ([#1725](https://github.com/CooooolFrog/dragonfly/issues/1725))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.3.0 to 1.3.3 ([#1722](https://github.com/CooooolFrog/dragonfly/issues/1722))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.9.0 to 1.10.0 ([#1720](https://github.com/CooooolFrog/dragonfly/issues/1720))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.34.0 to 0.36.1 ([#1719](https://github.com/CooooolFrog/dragonfly/issues/1719))
- **deps:** bump github.com/appleboy/gin-jwt/v2 from 2.8.0 to 2.9.0 ([#1718](https://github.com/CooooolFrog/dragonfly/issues/1718))

### Feat
- add TaskStateLeave to task ([#1728](https://github.com/CooooolFrog/dragonfly/issues/1728))
- searcher calculates cluster type ([#1729](https://github.com/CooooolFrog/dragonfly/issues/1729))
- unregister failed task storage ([#1717](https://github.com/CooooolFrog/dragonfly/issues/1717))
- oss get metadata ([#1724](https://github.com/CooooolFrog/dragonfly/issues/1724))
- support concurrent recursive download ([#1714](https://github.com/CooooolFrog/dragonfly/issues/1714))

### Fix
- gorm can not update is_default field ([#1731](https://github.com/CooooolFrog/dragonfly/issues/1731))
- content length is zero when task succeed ([#1732](https://github.com/CooooolFrog/dragonfly/issues/1732))

### Refactor
- idgen pkg ([#1715](https://github.com/CooooolFrog/dragonfly/issues/1715))


<a name="v2.0.7-beta.3"></a>
## [v2.0.7-beta.3] - 2022-09-29
### Chore
- **deps:** bump github.com/onsi/ginkgo/v2 from 2.1.6 to 2.2.0 ([#1705](https://github.com/CooooolFrog/dragonfly/issues/1705))
- **deps:** bump google.golang.org/api from 0.94.0 to 0.97.0 ([#1709](https://github.com/CooooolFrog/dragonfly/issues/1709))
- **deps:** bump k8s.io/component-base from 0.25.0 to 0.25.2 ([#1708](https://github.com/CooooolFrog/dragonfly/issues/1708))
- **deps:** bump gorm.io/gorm from 1.23.9 to 1.23.10 ([#1707](https://github.com/CooooolFrog/dragonfly/issues/1707))
- **deps:** bump github.com/casbin/casbin/v2 from 2.55.0 to 2.55.1 ([#1706](https://github.com/CooooolFrog/dragonfly/issues/1706))

### Feat
- add traffic shaper for download tasks ([#1654](https://github.com/CooooolFrog/dragonfly/issues/1654))
- async create a record ([#1711](https://github.com/CooooolFrog/dragonfly/issues/1711))

### Fix
- docker compose config ([#1713](https://github.com/CooooolFrog/dragonfly/issues/1713))

### Refactor
- pkg basic ([#1712](https://github.com/CooooolFrog/dragonfly/issues/1712))

### Test
- remove test main ([#1710](https://github.com/CooooolFrog/dragonfly/issues/1710))
- add test for daemon rpcserver ([#1704](https://github.com/CooooolFrog/dragonfly/issues/1704))


<a name="v2.0.7-beta.2"></a>
## [v2.0.7-beta.2] - 2022-09-26
### Chore
- update api pkg ([#1700](https://github.com/CooooolFrog/dragonfly/issues/1700))

### Feat
- optimize storage log ([#1703](https://github.com/CooooolFrog/dragonfly/issues/1703))

### Fix
- hdfs not registered ([#1702](https://github.com/CooooolFrog/dragonfly/issues/1702))

### Refactor
- manager and scheduler config ([#1701](https://github.com/CooooolFrog/dragonfly/issues/1701))


<a name="v2.0.7-beta.1"></a>
## [v2.0.7-beta.1] - 2022-09-22
### Feat
- remove ipv4 and ipv6 log ([#1699](https://github.com/CooooolFrog/dragonfly/issues/1699))
- enable ipv6 in unit test ([#1698](https://github.com/CooooolFrog/dragonfly/issues/1698))
- ipv6 support ([#1685](https://github.com/CooooolFrog/dragonfly/issues/1685))
- update docker compose image ([#1696](https://github.com/CooooolFrog/dragonfly/issues/1696))
- manager add advertiseIP ([#1695](https://github.com/CooooolFrog/dragonfly/issues/1695))

### Fix
- grpc download tidy file error ([#1697](https://github.com/CooooolFrog/dragonfly/issues/1697))

### Refactor
- listenIP and advertiseIP ([#1694](https://github.com/CooooolFrog/dragonfly/issues/1694))


<a name="v2.0.7-beta.0"></a>
## [v2.0.7-beta.0] - 2022-09-20
### Chore
- update download rpc check ([#1684](https://github.com/CooooolFrog/dragonfly/issues/1684))
- **deps:** bump gorm.io/gorm from 1.23.8 to 1.23.9 ([#1691](https://github.com/CooooolFrog/dragonfly/issues/1691))
- **deps:** bump go.opentelemetry.io/otel/sdk from 1.9.0 to 1.10.0 ([#1692](https://github.com/CooooolFrog/dragonfly/issues/1692))
- **deps:** bump gorm.io/driver/postgres from 1.3.9 to 1.3.10 ([#1690](https://github.com/CooooolFrog/dragonfly/issues/1690))
- **deps:** bump github.com/go-playground/validator/v10 from 10.11.0 to 10.11.1 ([#1689](https://github.com/CooooolFrog/dragonfly/issues/1689))
- **deps:** bump d7y.io/api from 1.1.4 to 1.1.6 ([#1688](https://github.com/CooooolFrog/dragonfly/issues/1688))

### Feat
- empty file e2e ([#1687](https://github.com/CooooolFrog/dragonfly/issues/1687))
- support download empty file ([#1686](https://github.com/CooooolFrog/dragonfly/issues/1686))


<a name="v2.0.7-alpha.5"></a>
## [v2.0.7-alpha.5] - 2022-09-15
### Chore
- **deps:** bump github.com/spf13/viper from 1.12.0 to 1.13.0 ([#1676](https://github.com/CooooolFrog/dragonfly/issues/1676))
- **deps:** bump github.com/casbin/casbin/v2 from 2.53.2 to 2.55.0 ([#1679](https://github.com/CooooolFrog/dragonfly/issues/1679))
- **deps:** bump k8s.io/component-base from 0.23.3 to 0.25.0 ([#1674](https://github.com/CooooolFrog/dragonfly/issues/1674))
- **deps:** bump github.com/swaggo/gin-swagger from 1.5.2 to 1.5.3 ([#1673](https://github.com/CooooolFrog/dragonfly/issues/1673))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.91 to 1.44.95 ([#1672](https://github.com/CooooolFrog/dragonfly/issues/1672))

### Fix
- manager redis config convert ([#1680](https://github.com/CooooolFrog/dragonfly/issues/1680))


<a name="v2.0.7-alpha.4"></a>
## [v2.0.7-alpha.4] - 2022-09-09
### Feat
- stop grpc client ([#1671](https://github.com/CooooolFrog/dragonfly/issues/1671))
- change event DownloadFromBackToSource ([#1670](https://github.com/CooooolFrog/dragonfly/issues/1670))


<a name="v2.0.7-alpha.3"></a>
## [v2.0.7-alpha.3] - 2022-09-08
### Chore
- update helm-charts submodule version ([#1669](https://github.com/CooooolFrog/dragonfly/issues/1669))

### Feat
- dfget supports config file ([#1668](https://github.com/CooooolFrog/dragonfly/issues/1668))
- split concurrent back source e2e ([#1666](https://github.com/CooooolFrog/dragonfly/issues/1666))


<a name="v2.0.7-alpha.2"></a>
## [v2.0.7-alpha.2] - 2022-09-08
### Feat
- support redis cluster ([#1667](https://github.com/CooooolFrog/dragonfly/issues/1667))


<a name="v2.0.7-alpha.1"></a>
## [v2.0.7-alpha.1] - 2022-09-08
### Chore
- add disable seed peer action ([#1653](https://github.com/CooooolFrog/dragonfly/issues/1653))


<a name="v2.0.7-alpha.0"></a>
## [v2.0.7-alpha.0] - 2022-09-07
### Feat
- source changes ResponseHeaderTimeout and ExpectContinueTimeout ([#1662](https://github.com/CooooolFrog/dragonfly/issues/1662))
- change dfdaemon rate limit ([#1661](https://github.com/CooooolFrog/dragonfly/issues/1661))
- set created_at and updated_at to timestamp ([#1659](https://github.com/CooooolFrog/dragonfly/issues/1659))
- stat peer metrics with memory cache ([#1660](https://github.com/CooooolFrog/dragonfly/issues/1660))
- change storage strategy to simple ([#1658](https://github.com/CooooolFrog/dragonfly/issues/1658))
- add missing client version for ListSchedulers ([#1657](https://github.com/CooooolFrog/dragonfly/issues/1657))

### Fix
- task CanBackToSource func ([#1663](https://github.com/CooooolFrog/dragonfly/issues/1663))


<a name="v2.0.6"></a>
## [v2.0.6] - 2022-09-06
### Chore
- gitignore add .run
- update tls e2e cert ([#1626](https://github.com/CooooolFrog/dragonfly/issues/1626))
- release v2.0.6 version ([#1627](https://github.com/CooooolFrog/dragonfly/issues/1627))
- dependabot add github-actions ([#1629](https://github.com/CooooolFrog/dragonfly/issues/1629))
- **deps:** bump github.com/swaggo/swag from 1.8.4 to 1.8.5 ([#1636](https://github.com/CooooolFrog/dragonfly/issues/1636))
- **deps:** bump go.uber.org/zap from 1.21.0 to 1.23.0 ([#1635](https://github.com/CooooolFrog/dragonfly/issues/1635))
- **deps:** bump github.com/casbin/casbin/v2 from 2.52.2 to 2.53.2 ([#1644](https://github.com/CooooolFrog/dragonfly/issues/1644))
- **deps:** bump go.uber.org/atomic from 1.9.0 to 1.10.0 ([#1639](https://github.com/CooooolFrog/dragonfly/issues/1639))
- **deps:** bump google.golang.org/api from 0.92.0 to 0.94.0 ([#1638](https://github.com/CooooolFrog/dragonfly/issues/1638))
- **deps:** bump github.com/onsi/gomega from 1.20.0 to 1.20.2 ([#1637](https://github.com/CooooolFrog/dragonfly/issues/1637))
- **deps:** bump goreleaser/goreleaser-action from 2 to 3 ([#1650](https://github.com/CooooolFrog/dragonfly/issues/1650))
- **deps:** bump gorm.io/plugin/soft_delete from 1.1.0 to 1.2.0 ([#1643](https://github.com/CooooolFrog/dragonfly/issues/1643))
- **deps:** bump docker/setup-buildx-action from 1 to 2 ([#1634](https://github.com/CooooolFrog/dragonfly/issues/1634))
- **deps:** bump actions/setup-go from 2 to 3 ([#1633](https://github.com/CooooolFrog/dragonfly/issues/1633))
- **deps:** bump actions/checkout from 2 to 3 ([#1631](https://github.com/CooooolFrog/dragonfly/issues/1631))
- **deps:** bump codecov/codecov-action from 1 to 3 ([#1630](https://github.com/CooooolFrog/dragonfly/issues/1630))
- **deps:** bump actions/upload-artifact from 2 to 3 ([#1632](https://github.com/CooooolFrog/dragonfly/issues/1632))
- **deps:** bump github.com/aws/aws-sdk-go from 1.44.44 to 1.44.91 ([#1647](https://github.com/CooooolFrog/dragonfly/issues/1647))
- **deps:** bump docker/build-push-action from 2 to 3 ([#1648](https://github.com/CooooolFrog/dragonfly/issues/1648))
- **deps:** bump docker/login-action from 1 to 2 ([#1649](https://github.com/CooooolFrog/dragonfly/issues/1649))

### Feat
- add MaxConnectionIdle to grpc keepalive ([#1655](https://github.com/CooooolFrog/dragonfly/issues/1655))
- change consistent hashing element name ([#1652](https://github.com/CooooolFrog/dragonfly/issues/1652))

### Fix
- manager embed assets ([#1642](https://github.com/CooooolFrog/dragonfly/issues/1642))
- dfstore flags invalid ([#1641](https://github.com/CooooolFrog/dragonfly/issues/1641))


<a name="v2.0.6-beta.3"></a>
## [v2.0.6-beta.3] - 2022-09-01
### Chore
- workflows add tls e2e ([#1624](https://github.com/CooooolFrog/dragonfly/issues/1624))
- update debug info ([#1617](https://github.com/CooooolFrog/dragonfly/issues/1617))

### Feat
- add cert spec to security configuration ([#1621](https://github.com/CooooolFrog/dragonfly/issues/1621))
- support mutate all proxy requests ([#1623](https://github.com/CooooolFrog/dragonfly/issues/1623))


<a name="v2.0.6-beta.2"></a>
## [v2.0.6-beta.2] - 2022-08-31
### Feat
- check whether scheduler is in the same cluster ([#1620](https://github.com/CooooolFrog/dragonfly/issues/1620))
- manager add cert spec ([#1619](https://github.com/CooooolFrog/dragonfly/issues/1619))


<a name="v2.0.6-beta.1"></a>
## [v2.0.6-beta.1] - 2022-08-31
### Chore
- fix macos build ([#1609](https://github.com/CooooolFrog/dragonfly/issues/1609))
- add source error metrics ([#1560](https://github.com/CooooolFrog/dragonfly/issues/1560))
- update new manager ([#1597](https://github.com/CooooolFrog/dragonfly/issues/1597))
- **deps:** bump gorm.io/driver/postgres from 1.3.8 to 1.3.9 ([#1608](https://github.com/CooooolFrog/dragonfly/issues/1608))
- **deps:** bump github.com/aliyun/aliyun-oss-go-sdk from 2.2.4+incompatible to 2.2.5+incompatible ([#1607](https://github.com/CooooolFrog/dragonfly/issues/1607))
- **deps:** bump github.com/bits-and-blooms/bitset from 1.2.2 to 1.3.0 ([#1606](https://github.com/CooooolFrog/dragonfly/issues/1606))
- **deps:** bump github.com/gin-contrib/cors from 1.3.1 to 1.4.0 ([#1605](https://github.com/CooooolFrog/dragonfly/issues/1605))
- **deps:** bump github.com/swaggo/gin-swagger from 1.5.1 to 1.5.2 ([#1604](https://github.com/CooooolFrog/dragonfly/issues/1604))
- **deps:** bump github.com/casbin/casbin/v2 from 2.51.2 to 2.52.2 ([#1588](https://github.com/CooooolFrog/dragonfly/issues/1588))
- **deps:** bump github.com/swaggo/swag from 1.8.3 to 1.8.4 ([#1590](https://github.com/CooooolFrog/dragonfly/issues/1590))
- **deps:** bump k8s.io/apimachinery from 0.24.2 to 0.24.4 ([#1591](https://github.com/CooooolFrog/dragonfly/issues/1591))

### Feat
- add tls policy to scheduler grpc server ([#1616](https://github.com/CooooolFrog/dragonfly/issues/1616))
- set tls cert leaf ([#1615](https://github.com/CooooolFrog/dragonfly/issues/1615))
- resolver addr add ServerName ([#1614](https://github.com/CooooolFrog/dragonfly/issues/1614))
- refactor grpc credential ([#1612](https://github.com/CooooolFrog/dragonfly/issues/1612))
- add tls policy to manager grpc server ([#1611](https://github.com/CooooolFrog/dragonfly/issues/1611))
- add tls policy constants ([#1610](https://github.com/CooooolFrog/dragonfly/issues/1610))
- add grpc mux transport ([#1602](https://github.com/CooooolFrog/dragonfly/issues/1602))
- manager init cert for grpc server ([#1603](https://github.com/CooooolFrog/dragonfly/issues/1603))
- refactor peertask option ([#1600](https://github.com/CooooolFrog/dragonfly/issues/1600))
- add common serialize package ([#1601](https://github.com/CooooolFrog/dragonfly/issues/1601))
- add client grpc dial timeout ([#1599](https://github.com/CooooolFrog/dragonfly/issues/1599))
- support multiple certify cache ([#1598](https://github.com/CooooolFrog/dragonfly/issues/1598))
- PeerGauge adds version and commit labels ([#1596](https://github.com/CooooolFrog/dragonfly/issues/1596))
- daemon support auto issue certificate ([#1586](https://github.com/CooooolFrog/dragonfly/issues/1586))
- add default metrics address ([#1595](https://github.com/CooooolFrog/dragonfly/issues/1595))
- grpc dial adds context ([#1594](https://github.com/CooooolFrog/dragonfly/issues/1594))
- consistent hashing add picker log ([#1593](https://github.com/CooooolFrog/dragonfly/issues/1593))

### Fix
- ci actions with docker ([#1613](https://github.com/CooooolFrog/dragonfly/issues/1613))

### Refactor
- dfpath for certify cache dir ([#1618](https://github.com/CooooolFrog/dragonfly/issues/1618))


<a name="v2.0.6-beta.0"></a>
## [v2.0.6-beta.0] - 2022-08-19
### Feat
- remove golang +build tag ([#1585](https://github.com/CooooolFrog/dragonfly/issues/1585))
- manager add certificate config ([#1583](https://github.com/CooooolFrog/dragonfly/issues/1583))
- manager implements issue certificate grpc ([#1577](https://github.com/CooooolFrog/dragonfly/issues/1577))
- dfdaemon add convert interceptor ([#1582](https://github.com/CooooolFrog/dragonfly/issues/1582))


<a name="v2.0.6-alpha.3"></a>
## [v2.0.6-alpha.3] - 2022-08-18
### Feat
- dynconfig refresh and notify listeners ([#1579](https://github.com/CooooolFrog/dragonfly/issues/1579))

### Fix
- dfdaemon can not shutdown ([#1580](https://github.com/CooooolFrog/dragonfly/issues/1580))

### Refactor
- dfnet package ([#1578](https://github.com/CooooolFrog/dragonfly/issues/1578))
- dfdaemon client and remove rpc connection pool ([#1576](https://github.com/CooooolFrog/dragonfly/issues/1576))


<a name="v2.0.6-alpha.2"></a>
## [v2.0.6-alpha.2] - 2022-08-17
### Chore
- **deps:** bump gorm.io/driver/mysql from 1.3.4 to 1.3.6 ([#1567](https://github.com/CooooolFrog/dragonfly/issues/1567))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.33.0 to 0.34.0 ([#1566](https://github.com/CooooolFrog/dragonfly/issues/1566))
- **deps:** bump google.golang.org/api from 0.90.0 to 0.92.0 ([#1565](https://github.com/CooooolFrog/dragonfly/issues/1565))
- **deps:** bump github.com/prometheus/client_golang from 1.12.2 to 1.13.0 ([#1564](https://github.com/CooooolFrog/dragonfly/issues/1564))

### Feat
- add grpc client error interceptor ([#1575](https://github.com/CooooolFrog/dragonfly/issues/1575))
- grpc removes MaxConnectionIdle ([#1574](https://github.com/CooooolFrog/dragonfly/issues/1574))
- grpc add ratelimit ([#1572](https://github.com/CooooolFrog/dragonfly/issues/1572))
- refresh dynconfig addresses when grpc requests unavailable ([#1571](https://github.com/CooooolFrog/dragonfly/issues/1571))
- manager adds model and model version grpc api ([#1569](https://github.com/CooooolFrog/dragonfly/issues/1569))

### Fix
- scheduler can not exit gracefully due to machinery fatal log ([#1573](https://github.com/CooooolFrog/dragonfly/issues/1573))


<a name="v2.0.6-alpha.1"></a>
## [v2.0.6-alpha.1] - 2022-08-12
### Chore
- optimize source error log ([#1553](https://github.com/CooooolFrog/dragonfly/issues/1553))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin from 0.32.0 to 0.34.0 ([#1547](https://github.com/CooooolFrog/dragonfly/issues/1547))
- **deps:** bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0 ([#1544](https://github.com/CooooolFrog/dragonfly/issues/1544))
- **deps:** bump github.com/jarcoal/httpmock from 1.0.8 to 1.2.0 ([#1542](https://github.com/CooooolFrog/dragonfly/issues/1542))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.8.0 to 1.9.0 ([#1541](https://github.com/CooooolFrog/dragonfly/issues/1541))

### Feat
- dynconfig add refresh func ([#1563](https://github.com/CooooolFrog/dragonfly/issues/1563))
- manager client add context ([#1562](https://github.com/CooooolFrog/dragonfly/issues/1562))
- grpc add retry middleware ([#1561](https://github.com/CooooolFrog/dragonfly/issues/1561))
- grpc consistent hashing ([#1554](https://github.com/CooooolFrog/dragonfly/issues/1554))
- model version add training result ([#1558](https://github.com/CooooolFrog/dragonfly/issues/1558))
- storage calculate the count of records ([#1557](https://github.com/CooooolFrog/dragonfly/issues/1557))
- model and model version api removes auth ([#1556](https://github.com/CooooolFrog/dragonfly/issues/1556))
- add seed trace ([#1549](https://github.com/CooooolFrog/dragonfly/issues/1549))
- gc removes logrus ([#1548](https://github.com/CooooolFrog/dragonfly/issues/1548))
- add MultiReadCloser and storage add open func ([#1546](https://github.com/CooooolFrog/dragonfly/issues/1546))
- scheduler dynconfig returns more info ([#1545](https://github.com/CooooolFrog/dragonfly/issues/1545))
- scheduler and manager change graceful stop timeout ([#1540](https://github.com/CooooolFrog/dragonfly/issues/1540))
- schedulers create main peer record ([#1539](https://github.com/CooooolFrog/dragonfly/issues/1539))
- change update model api ([#1538](https://github.com/CooooolFrog/dragonfly/issues/1538))
- manager adds model and model version api ([#1530](https://github.com/CooooolFrog/dragonfly/issues/1530))
- when the request has a range header, object storage is no need to  to calculate md5 ([#1534](https://github.com/CooooolFrog/dragonfly/issues/1534))

### Fix
- scheduler and manager tracing ([#1551](https://github.com/CooooolFrog/dragonfly/issues/1551))
- scheduler's MainParent func ([#1550](https://github.com/CooooolFrog/dragonfly/issues/1550))


<a name="v2.0.6-alpha.0"></a>
## [v2.0.6-alpha.0] - 2022-08-04
### Chore
- dragonfly updates version to v2.0.5 ([#1498](https://github.com/CooooolFrog/dragonfly/issues/1498))
- **deps:** bump google.golang.org/grpc from 1.47.0 to 1.48.0 ([#1508](https://github.com/CooooolFrog/dragonfly/issues/1508))
- **deps:** bump github.com/casbin/casbin/v2 from 2.48.0 to 2.51.2 ([#1512](https://github.com/CooooolFrog/dragonfly/issues/1512))
- **deps:** bump github.com/shirou/gopsutil/v3 from 3.22.5 to 3.22.7 ([#1511](https://github.com/CooooolFrog/dragonfly/issues/1511))
- **deps:** bump google.golang.org/api from 0.86.0 to 0.90.0 ([#1510](https://github.com/CooooolFrog/dragonfly/issues/1510))
- **deps:** bump go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc from 0.32.0 to 0.33.0 ([#1509](https://github.com/CooooolFrog/dragonfly/issues/1509))
- **deps:** bump gorm.io/driver/postgres from 1.3.7 to 1.3.8 ([#1503](https://github.com/CooooolFrog/dragonfly/issues/1503))
- **deps:** bump go.opentelemetry.io/otel/exporters/jaeger from 1.7.0 to 1.8.0 ([#1506](https://github.com/CooooolFrog/dragonfly/issues/1506))
- **deps:** bump github.com/swaggo/gin-swagger from 1.5.0 to 1.5.1 ([#1505](https://github.com/CooooolFrog/dragonfly/issues/1505))
- **deps:** bump github.com/schollz/progressbar/v3 from 3.8.6 to 3.8.7 ([#1502](https://github.com/CooooolFrog/dragonfly/issues/1502))

### Docs
- add daemon-socket for daemon docs ([#1522](https://github.com/CooooolFrog/dragonfly/issues/1522))

### Feat
- support grpc recursive download ([#1518](https://github.com/CooooolFrog/dragonfly/issues/1518))
- manager embed frontend assets ([#1523](https://github.com/CooooolFrog/dragonfly/issues/1523))
- can not return peer with the same host ([#1526](https://github.com/CooooolFrog/dragonfly/issues/1526))
- add daemon-socket-path ([#1521](https://github.com/CooooolFrog/dragonfly/issues/1521))
- store preheat result ([#1516](https://github.com/CooooolFrog/dragonfly/issues/1516))
- replace grpc package with https://github.com/dragonflyoss/api ([#1515](https://github.com/CooooolFrog/dragonfly/issues/1515))
- dfdaemon add Authorization and WWWAuthenticate headers ([#1513](https://github.com/CooooolFrog/dragonfly/issues/1513))
- auto switch to concurrent back source based on download speed ([#1494](https://github.com/CooooolFrog/dragonfly/issues/1494))
- enable dependabot ([#1501](https://github.com/CooooolFrog/dragonfly/issues/1501))
- scheduler adds filter range limit ([#1497](https://github.com/CooooolFrog/dragonfly/issues/1497))

### Fix
- check same peer id for sync pieces ([#1525](https://github.com/CooooolFrog/dragonfly/issues/1525))
- auto switch to concurrent back source ([#1507](https://github.com/CooooolFrog/dragonfly/issues/1507))
- wait first peer packet fail ([#1500](https://github.com/CooooolFrog/dragonfly/issues/1500))
- one piece task sometimes backsource after succeed ([#1499](https://github.com/CooooolFrog/dragonfly/issues/1499))
- random vertices ([#1496](https://github.com/CooooolFrog/dragonfly/issues/1496))


<a name="v2.0.5-rc.0"></a>
## [v2.0.5-rc.0] - 2022-07-27
### Feat
- scheduler set workhome ([#1493](https://github.com/CooooolFrog/dragonfly/issues/1493))

### Fix
- dfstore command-line tool name ([#1492](https://github.com/CooooolFrog/dragonfly/issues/1492))


<a name="v2.0.5-beta.5"></a>
## [v2.0.5-beta.5] - 2022-07-26
### Fix
- oss client judge directory bug ([#1488](https://github.com/CooooolFrog/dragonfly/issues/1488))
- dfdaemon unix socket ([#1489](https://github.com/CooooolFrog/dragonfly/issues/1489))

### Refactor
- set and dag with generics ([#1490](https://github.com/CooooolFrog/dragonfly/issues/1490))


<a name="v2.0.5-beta.4"></a>
## [v2.0.5-beta.4] - 2022-07-25
### Feat
- remove test print

### Fix
- init storage error ([#1486](https://github.com/CooooolFrog/dragonfly/issues/1486))
- back source error ([#1485](https://github.com/CooooolFrog/dragonfly/issues/1485))
- keepalive with ip

### Refactor
- cache key for peer ([#1483](https://github.com/CooooolFrog/dragonfly/issues/1483))
- scheduler training configuration
- dag GetSourceVertices and GetSinkVertices func

### Test
- find parent with ancestor ([#1482](https://github.com/CooooolFrog/dragonfly/issues/1482))


<a name="v2.0.5-beta.3"></a>
## [v2.0.5-beta.3] - 2022-07-22
### Feat
- rename steal peers to candidate peers ([#1476](https://github.com/CooooolFrog/dragonfly/issues/1476))
- scheduler merge end of piece and piece from seed peer ([#1474](https://github.com/CooooolFrog/dragonfly/issues/1474))


<a name="v2.0.5-beta.2"></a>
## [v2.0.5-beta.2] - 2022-07-21
### Feat
- dfdaemon add default healthy config ([#1472](https://github.com/CooooolFrog/dragonfly/issues/1472))
- dag adds LenVertex and RangeVertex func ([#1470](https://github.com/CooooolFrog/dragonfly/issues/1470))

### Fix
- upload_manager write header in time ([#1471](https://github.com/CooooolFrog/dragonfly/issues/1471))
- infinite loop in peer.Ancestors() ([#1469](https://github.com/CooooolFrog/dragonfly/issues/1469))


<a name="v2.0.5-beta.1"></a>
## [v2.0.5-beta.1] - 2022-07-18
### Feat
- generate dag mock
- add directed acyclic graph package ([#1468](https://github.com/CooooolFrog/dragonfly/issues/1468))

### Fix
- upload_manager write header immediately when it is ready ([#1466](https://github.com/CooooolFrog/dragonfly/issues/1466))


<a name="v2.0.5-beta.0"></a>
## [v2.0.5-beta.0] - 2022-07-14
### Feat
- proxy add defaultTag field ([#1462](https://github.com/CooooolFrog/dragonfly/issues/1462))
- manager support postgres ([#1459](https://github.com/CooooolFrog/dragonfly/issues/1459))
- use os.PathSeparator to generate object key
- scheduler add data dir ([#1453](https://github.com/CooooolFrog/dragonfly/issues/1453))

### Fix
- metrics reduces labels ([#1457](https://github.com/CooooolFrog/dragonfly/issues/1457))


<a name="v2.0.5-alpha.3"></a>
## [v2.0.5-alpha.3] - 2022-07-12
### Chore
- check header length before update ([#1445](https://github.com/CooooolFrog/dragonfly/issues/1445))

### Feat
- dfdaemon is compatible with v2.0.2 ([#1452](https://github.com/CooooolFrog/dragonfly/issues/1452))
- add slices util package
- reload proxy option ([#1443](https://github.com/CooooolFrog/dragonfly/issues/1443))
- if peer back-to-source failed, return source metadata. ([#1444](https://github.com/CooooolFrog/dragonfly/issues/1444))
- report peer result with source error detail ([#1439](https://github.com/CooooolFrog/dragonfly/issues/1439))

### Fix
- depth limit ([#1451](https://github.com/CooooolFrog/dragonfly/issues/1451))
- dfpath creates redundant directories ([#1446](https://github.com/CooooolFrog/dragonfly/issues/1446))

### Refactor
- rewrite math max and min with generics ([#1447](https://github.com/CooooolFrog/dragonfly/issues/1447))


<a name="v2.0.5-alpha.2"></a>
## [v2.0.5-alpha.2] - 2022-07-07
### Chore
- update test/tools/no-content-length/main.go ([#1440](https://github.com/CooooolFrog/dragonfly/issues/1440))

### Fix
- release package name ([#1442](https://github.com/CooooolFrog/dragonfly/issues/1442))


<a name="v2.0.5-alpha.1"></a>
## [v2.0.5-alpha.1] - 2022-07-07
### Feat
- add dfstore command ([#1441](https://github.com/CooooolFrog/dragonfly/issues/1441))
- back source error detail ([#1437](https://github.com/CooooolFrog/dragonfly/issues/1437))
- change local cache ttl ([#1436](https://github.com/CooooolFrog/dragonfly/issues/1436))
- if service can not found fqdn, replace fqdn with hostname ([#1435](https://github.com/CooooolFrog/dragonfly/issues/1435))


<a name="v2.0.5-alpha.0"></a>
## [v2.0.5-alpha.0] - 2022-07-05
### Chore
- upgrade kind node version ([#1433](https://github.com/CooooolFrog/dragonfly/issues/1433))
- update docker compose ([#1431](https://github.com/CooooolFrog/dragonfly/issues/1431))
- exit when receive context done ([#1432](https://github.com/CooooolFrog/dragonfly/issues/1432))
- update codeql version ([#1428](https://github.com/CooooolFrog/dragonfly/issues/1428))

### Feat
- remove errors package ([#1434](https://github.com/CooooolFrog/dragonfly/issues/1434))
- concurrent multiple pieces back source ([#1426](https://github.com/CooooolFrog/dragonfly/issues/1426))

### Fix
- seed task metric panic ([#1427](https://github.com/CooooolFrog/dragonfly/issues/1427))


<a name="v2.0.4"></a>
## [v2.0.4] - 2022-07-01
### Chore
- release v2.0.4 ([#1425](https://github.com/CooooolFrog/dragonfly/issues/1425))


<a name="v2.0.4-rc.3"></a>
## [v2.0.4-rc.3] - 2022-06-30
### Feat
- dfstore closes writer ([#1424](https://github.com/CooooolFrog/dragonfly/issues/1424))


<a name="v2.0.4-rc.2"></a>
## [v2.0.4-rc.2] - 2022-06-30

<a name="v2.0.4-rc.1"></a>
## [v2.0.4-rc.1] - 2022-06-30
### Feat
- use put object action ([#1422](https://github.com/CooooolFrog/dragonfly/issues/1422))
- GetObjectInput add range field ([#1421](https://github.com/CooooolFrog/dragonfly/issues/1421))
- rename client/clientutil to client/util ([#1420](https://github.com/CooooolFrog/dragonfly/issues/1420))


<a name="v2.0.4-rc.0"></a>
## [v2.0.4-rc.0] - 2022-06-30
### Feat
- rewrite interface{} to any ([#1419](https://github.com/CooooolFrog/dragonfly/issues/1419))


<a name="v2.0.4-beta.2"></a>
## [v2.0.4-beta.2] - 2022-06-29
### Feat
- update namely/protoc-all image version to 1.47_0 ([#1418](https://github.com/CooooolFrog/dragonfly/issues/1418))
- update golang to 1.18.3 ([#1417](https://github.com/CooooolFrog/dragonfly/issues/1417))


<a name="v2.0.4-beta.1"></a>
## [v2.0.4-beta.1] - 2022-06-28
### Feat
- remove github/pkg/errors package ([#1416](https://github.com/CooooolFrog/dragonfly/issues/1416))
- add dfstore client interface ([#1415](https://github.com/CooooolFrog/dragonfly/issues/1415))
- scheduler http status ([#1414](https://github.com/CooooolFrog/dragonfly/issues/1414))
- enable configuration of the tls parameter for the mysql connection. i.e. tls=preferred ([#1300](https://github.com/CooooolFrog/dragonfly/issues/1300))


<a name="v2.0.4-beta.0"></a>
## [v2.0.4-beta.0] - 2022-06-27
### Chore
- update submodule version

### Feat
- import object to seed peer with max replicas ([#1413](https://github.com/CooooolFrog/dragonfly/issues/1413))
- object storage add filter field ([#1412](https://github.com/CooooolFrog/dragonfly/issues/1412))
- dfdaemon add destroyObject rest api ([#1410](https://github.com/CooooolFrog/dragonfly/issues/1410))
- client add create object storage ([#1409](https://github.com/CooooolFrog/dragonfly/issues/1409))
- seed peer add object storage port ([#1408](https://github.com/CooooolFrog/dragonfly/issues/1408))
- rename digest func and add new digest func ([#1405](https://github.com/CooooolFrog/dragonfly/issues/1405))
- dfdaemon upload and object storage service add middlewares ([#1404](https://github.com/CooooolFrog/dragonfly/issues/1404))

### Fix
- pkg/strings comment typo

### Refactor
- scheduler announce task ([#1407](https://github.com/CooooolFrog/dragonfly/issues/1407))
- digest package ([#1403](https://github.com/CooooolFrog/dragonfly/issues/1403))


<a name="v2.0.4-alpha.1"></a>
## [v2.0.4-alpha.1] - 2022-06-20
### Chore
- goreleaser remove cdn

### Feat
- remove cdn ([#1401](https://github.com/CooooolFrog/dragonfly/issues/1401))
- redirect stdout and stderr to file ([#1399](https://github.com/CooooolFrog/dragonfly/issues/1399))
- dfdaemon add GetObject rest api ([#1398](https://github.com/CooooolFrog/dragonfly/issues/1398))
- add seed peer for list scheduler grpc interface ([#1393](https://github.com/CooooolFrog/dragonfly/issues/1393))
- dfdaemon add object storage rest api ([#1390](https://github.com/CooooolFrog/dragonfly/issues/1390))
- replace gin-gonic/gin with gorilla/mux ([#1389](https://github.com/CooooolFrog/dragonfly/issues/1389))

### Fix
- downloadFromSource() doesn't validate response ([#1400](https://github.com/CooooolFrog/dragonfly/issues/1400))
- default repository does not exist and missing dependency containers ([#1395](https://github.com/CooooolFrog/dragonfly/issues/1395))
- validate rate limiter ([#1392](https://github.com/CooooolFrog/dragonfly/issues/1392))
- dfget ratelimit params ([#1391](https://github.com/CooooolFrog/dragonfly/issues/1391))
- count error & totalPage error ([#1373](https://github.com/CooooolFrog/dragonfly/issues/1373)) ([#1376](https://github.com/CooooolFrog/dragonfly/issues/1376))
- manager router middlewares order ([#1385](https://github.com/CooooolFrog/dragonfly/issues/1385))

### Refactor
- pkg util ([#1402](https://github.com/CooooolFrog/dragonfly/issues/1402))


<a name="v2.0.4-alpha.0"></a>
## [v2.0.4-alpha.0] - 2022-06-14
### Chore
- add check size workflows ([#1364](https://github.com/CooooolFrog/dragonfly/issues/1364))
- add hack/gen-containerd-hosts.sh ([#1361](https://github.com/CooooolFrog/dragonfly/issues/1361))
- release v2.0.3 ([#1360](https://github.com/CooooolFrog/dragonfly/issues/1360))
- update content range for partial content ([#1357](https://github.com/CooooolFrog/dragonfly/issues/1357))

### Docs
- update CHANGELOG
- update CHANGELOG
- update readme system features ([#1359](https://github.com/CooooolFrog/dragonfly/issues/1359))

### Feat
- add enable config to peer gauge ([#1382](https://github.com/CooooolFrog/dragonfly/issues/1382))
- dfdaemon add ns filter ([#1379](https://github.com/CooooolFrog/dragonfly/issues/1379))
- remove connection gc ([#1378](https://github.com/CooooolFrog/dragonfly/issues/1378))
- dynconfig add object storage ([#1369](https://github.com/CooooolFrog/dragonfly/issues/1369))
- manager add bucket interface ([#1368](https://github.com/CooooolFrog/dragonfly/issues/1368))
- add objectstorage pkg ([#1366](https://github.com/CooooolFrog/dragonfly/issues/1366))
- remove preheat tag validate with required ([#1363](https://github.com/CooooolFrog/dragonfly/issues/1363))
- e2e seed peer ([#1358](https://github.com/CooooolFrog/dragonfly/issues/1358))
- update console and helm-charts submodule ([#1355](https://github.com/CooooolFrog/dragonfly/issues/1355))
- use uid/gid as UserID and UserGroup if current user not found in passwd ([#1352](https://github.com/CooooolFrog/dragonfly/issues/1352))
- use 127.0.0.1 as IPv4 if there's no external IPv4 addr ([#1353](https://github.com/CooooolFrog/dragonfly/issues/1353))

### Fix
- dfget build error ([#1381](https://github.com/CooooolFrog/dragonfly/issues/1381))
- preheat tack id ([#1375](https://github.com/CooooolFrog/dragonfly/issues/1375))


<a name="v2.0.3-beta.8"></a>
## [v2.0.3-beta.8] - 2022-05-31
### Chore
- add check size action ([#1350](https://github.com/CooooolFrog/dragonfly/issues/1350))

### Docs
- readme typo
- readme add seed peer ([#1349](https://github.com/CooooolFrog/dragonfly/issues/1349))

### Feat
- add security group id with scheduler cluster ([#1354](https://github.com/CooooolFrog/dragonfly/issues/1354))
- change pattern from cdn to seed peer and remove kustomize shell ([#1345](https://github.com/CooooolFrog/dragonfly/issues/1345))

### Fix
- register fail panic ([#1351](https://github.com/CooooolFrog/dragonfly/issues/1351))
- find partial completed overflow ([#1346](https://github.com/CooooolFrog/dragonfly/issues/1346))


<a name="v2.0.3-beta.7"></a>
## [v2.0.3-beta.7] - 2022-05-31
### Feat
- update casbin/gorm-adapter version and change e2e charts config
- update helm charts
- update dependencies
- add seed peer metrics ([#1342](https://github.com/CooooolFrog/dragonfly/issues/1342))
- grpc health probe support arm64 ([#1338](https://github.com/CooooolFrog/dragonfly/issues/1338))
- docker build with multi platforms ([#1337](https://github.com/CooooolFrog/dragonfly/issues/1337))

### Fix
- e2e charts config
- seed peer reuse value
- dfdaemon seed peer metrics namespace ([#1343](https://github.com/CooooolFrog/dragonfly/issues/1343))
- create_at timestamp ([#1341](https://github.com/CooooolFrog/dragonfly/issues/1341))

### Test
- update e2e charts config


<a name="v2.0.3-beta.5"></a>
## [v2.0.3-beta.5] - 2022-05-25
### Chore
- check large files in pull request ([#1332](https://github.com/CooooolFrog/dragonfly/issues/1332))

### Fix
- reuse seed peer id is not exist ([#1335](https://github.com/CooooolFrog/dragonfly/issues/1335))
- reuse non-end range request ([#1333](https://github.com/CooooolFrog/dragonfly/issues/1333))
- http range header validation ([#1334](https://github.com/CooooolFrog/dragonfly/issues/1334))


<a name="v2.0.3-beta.4"></a>
## [v2.0.3-beta.4] - 2022-05-24
### Feat
- add sync piece watchdog ([#1272](https://github.com/CooooolFrog/dragonfly/issues/1272))
- scheduler handles seed peer failed ([#1325](https://github.com/CooooolFrog/dragonfly/issues/1325))
- custom preheat tag parameters ([#1324](https://github.com/CooooolFrog/dragonfly/issues/1324))
- client add tls verify config ([#1323](https://github.com/CooooolFrog/dragonfly/issues/1323))
- scheduler register interface return task type ([#1318](https://github.com/CooooolFrog/dragonfly/issues/1318))
- get active peer count ([#1315](https://github.com/CooooolFrog/dragonfly/issues/1315))
- reduce dynconfig log ([#1312](https://github.com/CooooolFrog/dragonfly/issues/1312))
- back source when receive seed request ([#1309](https://github.com/CooooolFrog/dragonfly/issues/1309))
- update helm charts submodule ([#1308](https://github.com/CooooolFrog/dragonfly/issues/1308))
- add vsock network type support ([#1303](https://github.com/CooooolFrog/dragonfly/issues/1303))
- support response header ([#1292](https://github.com/CooooolFrog/dragonfly/issues/1292))

### Fix
- nfpms maintainer ([#1326](https://github.com/CooooolFrog/dragonfly/issues/1326))
- reuse seed panic ([#1319](https://github.com/CooooolFrog/dragonfly/issues/1319))
- seed peer did not send done seed result and no content length send ([#1316](https://github.com/CooooolFrog/dragonfly/issues/1316))
- remove container after generating protoc ([#1306](https://github.com/CooooolFrog/dragonfly/issues/1306))

### Refactor
- scheduler grpc ([#1310](https://github.com/CooooolFrog/dragonfly/issues/1310))

### Test
- watchdog
- close dfget back-to-souce ([#1317](https://github.com/CooooolFrog/dragonfly/issues/1317))


<a name="v2.0.3-beta.2"></a>
## [v2.0.3-beta.2] - 2022-05-16
### Feat
- add seed peer logic ([#1302](https://github.com/CooooolFrog/dragonfly/issues/1302))
- support more digest like sha256 ([#1304](https://github.com/CooooolFrog/dragonfly/issues/1304))
- implement pattern in client daemon ([#1231](https://github.com/CooooolFrog/dragonfly/issues/1231))
- scheduler add seed peer ([#1298](https://github.com/CooooolFrog/dragonfly/issues/1298))
- don't gc client rpc connection if ExpireTime is 0 ([#1296](https://github.com/CooooolFrog/dragonfly/issues/1296))
- update scheduler PeerResult validation ([#1294](https://github.com/CooooolFrog/dragonfly/issues/1294))
- manager add seed peer ([#1293](https://github.com/CooooolFrog/dragonfly/issues/1293))

### Fix
- digest reader and unit tests ([#1305](https://github.com/CooooolFrog/dragonfly/issues/1305))
- scheduler typo ([#1297](https://github.com/CooooolFrog/dragonfly/issues/1297))


<a name="v2.0.3-beta.0"></a>
## [v2.0.3-beta.0] - 2022-05-06
### Chore
- add target peer id in sync piece trace ([#1278](https://github.com/CooooolFrog/dragonfly/issues/1278))

### Feat
- implement client seed mode ([#1247](https://github.com/CooooolFrog/dragonfly/issues/1247))
- scheduler peer result validation ([#1288](https://github.com/CooooolFrog/dragonfly/issues/1288))
- use a golang native file server to replace nginx ([#1258](https://github.com/CooooolFrog/dragonfly/issues/1258))
- support build arm&&arm64 dfget ([#1286](https://github.com/CooooolFrog/dragonfly/issues/1286))
- update filter parent ([#1279](https://github.com/CooooolFrog/dragonfly/issues/1279))

### Fix
- keep accept header ([#1291](https://github.com/CooooolFrog/dragonfly/issues/1291))
- grpc total_piece_count comment ([#1289](https://github.com/CooooolFrog/dragonfly/issues/1289))
- run.sh threw error on mac ([#1285](https://github.com/CooooolFrog/dragonfly/issues/1285))
- docker compose run.sh ([#1282](https://github.com/CooooolFrog/dragonfly/issues/1282))
- legacy cdn peer ([#1283](https://github.com/CooooolFrog/dragonfly/issues/1283))
- filter parent condition ([#1277](https://github.com/CooooolFrog/dragonfly/issues/1277))

### Refactor
- scheduler task SizeScope ([#1287](https://github.com/CooooolFrog/dragonfly/issues/1287))


<a name="v2.0.3-alpha.0"></a>
## [v2.0.3-alpha.0] - 2022-04-24
### Chore
- optimize create synchronizer logic ([#1269](https://github.com/CooooolFrog/dragonfly/issues/1269))
- add sync pieces trace and update sync pieces logic for done task ([#1263](https://github.com/CooooolFrog/dragonfly/issues/1263))
- add schedule cron with e2e testing ([#1262](https://github.com/CooooolFrog/dragonfly/issues/1262))
- optimize sync pieces ([#1253](https://github.com/CooooolFrog/dragonfly/issues/1253))
- update pull request template ([#1251](https://github.com/CooooolFrog/dragonfly/issues/1251))
- update compatibility version to v2.0.2
- update helm-charts commit
- generate change log
- update manager console commit ([#1219](https://github.com/CooooolFrog/dragonfly/issues/1219))
- print client stream task error log ([#1210](https://github.com/CooooolFrog/dragonfly/issues/1210))
- report client back source error ([#1209](https://github.com/CooooolFrog/dragonfly/issues/1209))

### Docs
- move document from /docs to d7y.io ([#1229](https://github.com/CooooolFrog/dragonfly/issues/1229))

### Feat
- in tree plugin ([#1276](https://github.com/CooooolFrog/dragonfly/issues/1276))
- move dfnet to pkg dir ([#1265](https://github.com/CooooolFrog/dragonfly/issues/1265))
- add dfcache rpm/deb packages and man pages and publish in goreleaser ([#1259](https://github.com/CooooolFrog/dragonfly/issues/1259))
- add AnnounceTask and StatTask metrics ([#1256](https://github.com/CooooolFrog/dragonfly/issues/1256))
- define and implement new dfdaemon APIs to make dragonfly2 work as a distributed cache ([#1227](https://github.com/CooooolFrog/dragonfly/issues/1227))
- redirect daemon stdout stderr to file ([#1244](https://github.com/CooooolFrog/dragonfly/issues/1244))
- registerTask returns to the task in time ([#1250](https://github.com/CooooolFrog/dragonfly/issues/1250))
- docker-compose write log to file ([#1236](https://github.com/CooooolFrog/dragonfly/issues/1236))
- update docker compose version ([#1235](https://github.com/CooooolFrog/dragonfly/issues/1235))
- update to v2.0.2 ([#1232](https://github.com/CooooolFrog/dragonfly/issues/1232))
- scheduler blocks steal peers ([#1224](https://github.com/CooooolFrog/dragonfly/issues/1224))
- update manager console ([#1222](https://github.com/CooooolFrog/dragonfly/issues/1222))
- manager validate with config ([#1218](https://github.com/CooooolFrog/dragonfly/issues/1218))
- remove kustomize template ([#1216](https://github.com/CooooolFrog/dragonfly/issues/1216))
- add back source fail metric in client ([#1214](https://github.com/CooooolFrog/dragonfly/issues/1214))
- cannot delete a cluster with existing instances ([#1213](https://github.com/CooooolFrog/dragonfly/issues/1213))
- add type to DownloadFailureCount ([#1212](https://github.com/CooooolFrog/dragonfly/issues/1212))
- if the number of failed peers in the task is greater than FailedPeerCountLimit, then scheduler notifies running peers of failure ([#1211](https://github.com/CooooolFrog/dragonfly/issues/1211))
- optimize get available task ([#1208](https://github.com/CooooolFrog/dragonfly/issues/1208))

### Fix
- dfget daemon console log invalid ([#1275](https://github.com/CooooolFrog/dragonfly/issues/1275))
- scheduler config validation ([#1274](https://github.com/CooooolFrog/dragonfly/issues/1274))
- run.sh threw error on mac ([#1273](https://github.com/CooooolFrog/dragonfly/issues/1273))
- tree infinite loop ([#1271](https://github.com/CooooolFrog/dragonfly/issues/1271))
- acquire empty dst pid ([#1268](https://github.com/CooooolFrog/dragonfly/issues/1268))
- skip unsupported kernel in systemd service ([#1261](https://github.com/CooooolFrog/dragonfly/issues/1261))
- client synchronizer report error lock and dial grpc timeout ([#1260](https://github.com/CooooolFrog/dragonfly/issues/1260))
- prevent traversal tree from infinite loop ([#1266](https://github.com/CooooolFrog/dragonfly/issues/1266))
- error message ([#1255](https://github.com/CooooolFrog/dragonfly/issues/1255))
- client sync piece panic ([#1246](https://github.com/CooooolFrog/dragonfly/issues/1246))
- client superfluous usage gc ([#1243](https://github.com/CooooolFrog/dragonfly/issues/1243))
- client sync send unsafe call ([#1240](https://github.com/CooooolFrog/dragonfly/issues/1240))
- client unexpected timeout ([#1239](https://github.com/CooooolFrog/dragonfly/issues/1239))
- goreleaser config
- make generate ([#1228](https://github.com/CooooolFrog/dragonfly/issues/1228))
- calculate FreeUploadLoad ([#1226](https://github.com/CooooolFrog/dragonfly/issues/1226))
- sync pieces hang ([#1221](https://github.com/CooooolFrog/dragonfly/issues/1221))

### Test
- fix storage backups ([#1270](https://github.com/CooooolFrog/dragonfly/issues/1270))
- scheduler storage ([#1257](https://github.com/CooooolFrog/dragonfly/issues/1257))
- AnnounceTask and StatTask ([#1254](https://github.com/CooooolFrog/dragonfly/issues/1254))


<a name="v2.0.2-rc.27"></a>
## [v2.0.2-rc.27] - 2022-03-29
### Chore
- update workflows compatibility version ([#1192](https://github.com/CooooolFrog/dragonfly/issues/1192))

### Docs
- add slack and google groups ([#1203](https://github.com/CooooolFrog/dragonfly/issues/1203))

### Feat
- change scheduler and cdn listen ([#1205](https://github.com/CooooolFrog/dragonfly/issues/1205))
- scheduler add block peers set ([#1202](https://github.com/CooooolFrog/dragonfly/issues/1202))
- add grpc-health-probe to image ([#1196](https://github.com/CooooolFrog/dragonfly/issues/1196))
- add grpc health interface ([#1195](https://github.com/CooooolFrog/dragonfly/issues/1195))

### Fix
- client miss failed piece ([#1194](https://github.com/CooooolFrog/dragonfly/issues/1194))

### Refactor
- scheduler end and begin of piece ([#1189](https://github.com/CooooolFrog/dragonfly/issues/1189))


<a name="v2.0.2-rc.26"></a>
## [v2.0.2-rc.26] - 2022-03-25
### Chore
- change golangci-lint min-complexity value ([#1188](https://github.com/CooooolFrog/dragonfly/issues/1188))
- optimize stream peer task ([#1186](https://github.com/CooooolFrog/dragonfly/issues/1186))
- always fallback to legacy get pieces ([#1180](https://github.com/CooooolFrog/dragonfly/issues/1180))

### Feat
- remove grpc error code validate ([#1191](https://github.com/CooooolFrog/dragonfly/issues/1191))
- generate grpc protos in namely/protoc-all image ([#1187](https://github.com/CooooolFrog/dragonfly/issues/1187))
- scheduler addresses log ([#1183](https://github.com/CooooolFrog/dragonfly/issues/1183))
- manage GetCDN interface return scheduler info ([#1184](https://github.com/CooooolFrog/dragonfly/issues/1184))
- dfdaemon match scheduler with case insensitive ([#1181](https://github.com/CooooolFrog/dragonfly/issues/1181))
- add RBAC to manager config interface ([#1179](https://github.com/CooooolFrog/dragonfly/issues/1179))
- dfdaemon get available scheduler addresses in the same cluster ([#1178](https://github.com/CooooolFrog/dragonfly/issues/1178))
- implement grpc client side sync pieces ([#1167](https://github.com/CooooolFrog/dragonfly/issues/1167))
- seacher return multiple scheduler clusters ([#1175](https://github.com/CooooolFrog/dragonfly/issues/1175))
- replace time.Now().Sub by time.Since ([#1173](https://github.com/CooooolFrog/dragonfly/issues/1173))
- change DefaultServerOptions to variable
- change default scheduler filter parent limit ([#1166](https://github.com/CooooolFrog/dragonfly/issues/1166))
- implement bidirectional fetch pieces ([#1165](https://github.com/CooooolFrog/dragonfly/issues/1165))

### Fix
- client break error ([#1190](https://github.com/CooooolFrog/dragonfly/issues/1190))
- rpc cdn sync piece tasks ([#1168](https://github.com/CooooolFrog/dragonfly/issues/1168))
- subscriber data race ([#1169](https://github.com/CooooolFrog/dragonfly/issues/1169))

### Test
- fix e2e preheat case ([#1170](https://github.com/CooooolFrog/dragonfly/issues/1170))


<a name="v2.0.2-rc.25"></a>
## [v2.0.2-rc.25] - 2022-03-16
### Feat
- scheduler add default biz tag ([#1164](https://github.com/CooooolFrog/dragonfly/issues/1164))
- optimize proxy performance ([#1137](https://github.com/CooooolFrog/dragonfly/issues/1137))


<a name="v2.0.2-rc.24"></a>
## [v2.0.2-rc.24] - 2022-03-15
### Chore
- update go mod ([#1156](https://github.com/CooooolFrog/dragonfly/issues/1156))
- add makefile note ([#1155](https://github.com/CooooolFrog/dragonfly/issues/1155))

### Feat
- host remove peer ([#1161](https://github.com/CooooolFrog/dragonfly/issues/1161))
- change reschdule config ([#1158](https://github.com/CooooolFrog/dragonfly/issues/1158))

### Test
- cache expire interval ([#1160](https://github.com/CooooolFrog/dragonfly/issues/1160))
- add scheduler constructSuccessPeerPacket case ([#1154](https://github.com/CooooolFrog/dragonfly/issues/1154))


<a name="v2.0.2-rc.23"></a>
## [v2.0.2-rc.23] - 2022-03-11
### Feat
- update git submodule ([#1153](https://github.com/CooooolFrog/dragonfly/issues/1153))
- scheduler metrics add default value of biz tag ([#1151](https://github.com/CooooolFrog/dragonfly/issues/1151))


<a name="v2.0.2-rc.22"></a>
## [v2.0.2-rc.22] - 2022-03-10
### Chore
- change scheduler config ([#1140](https://github.com/CooooolFrog/dragonfly/issues/1140))

### Feat
- add user update interface and rename rest to service ([#1148](https://github.com/CooooolFrog/dragonfly/issues/1148))
- scheduler trace trigger cdn ([#1147](https://github.com/CooooolFrog/dragonfly/issues/1147))
- add scheduler traffic metrics ([#1143](https://github.com/CooooolFrog/dragonfly/issues/1143))
- update otel package version and fix otelgrpc goroutine leak ([#1141](https://github.com/CooooolFrog/dragonfly/issues/1141))
- add scheduler metrics ([#1139](https://github.com/CooooolFrog/dragonfly/issues/1139))

### Test
- scheduler service handlePieceFail ([#1146](https://github.com/CooooolFrog/dragonfly/issues/1146))


<a name="v2.0.2-rc.21"></a>
## [v2.0.2-rc.21] - 2022-03-08
### Feat
- scheduler remove inactive host ([#1135](https://github.com/CooooolFrog/dragonfly/issues/1135))
- task state for register ([#1132](https://github.com/CooooolFrog/dragonfly/issues/1132))

### Fix
- docker-compose run with mac throw error ([#1134](https://github.com/CooooolFrog/dragonfly/issues/1134))


<a name="v2.0.2-rc.20"></a>
## [v2.0.2-rc.20] - 2022-03-04
### Fix
- wrong md5 sign in cdn ([#1126](https://github.com/CooooolFrog/dragonfly/issues/1126))


<a name="v2.0.2-rc.19"></a>
## [v2.0.2-rc.19] - 2022-03-04
### Chore
- fast back source when get pieces task failed ([#1123](https://github.com/CooooolFrog/dragonfly/issues/1123))

### Feat
- change grpc client keepalive config ([#1125](https://github.com/CooooolFrog/dragonfly/issues/1125))
- scheduler change piece cost from nanosecond to millisecond ([#1119](https://github.com/CooooolFrog/dragonfly/issues/1119))
- support health probe in daemon ([#1120](https://github.com/CooooolFrog/dragonfly/issues/1120))

### Fix
- docker-compose preheat pending ([#1124](https://github.com/CooooolFrog/dragonfly/issues/1124))


<a name="v2.0.2-rc.18"></a>
## [v2.0.2-rc.18] - 2022-03-03
### Chore
- optimize reuse logic ([#1110](https://github.com/CooooolFrog/dragonfly/issues/1110))

### Feat
- when peer downloads finished, peer deletes parent ([#1116](https://github.com/CooooolFrog/dragonfly/issues/1116))
- change source client dialer config ([#1115](https://github.com/CooooolFrog/dragonfly/issues/1115))
- optimize scheduler log ([#1114](https://github.com/CooooolFrog/dragonfly/issues/1114))
- remove needless manager grpc proxy ([#1113](https://github.com/CooooolFrog/dragonfly/issues/1113))
- set grpc logger verbosity from env variable ([#1111](https://github.com/CooooolFrog/dragonfly/issues/1111))
- change back-to-source timeout ([#1112](https://github.com/CooooolFrog/dragonfly/issues/1112))

### Fix
- scheduler piece cost time ([#1118](https://github.com/CooooolFrog/dragonfly/issues/1118))


<a name="v2.0.2-rc.17"></a>
## [v2.0.2-rc.17] - 2022-03-02
### Chore
- init url meta in rpc server ([#1098](https://github.com/CooooolFrog/dragonfly/issues/1098))

### Docs
- add plugin builder ([#1101](https://github.com/CooooolFrog/dragonfly/issues/1101))

### Feat
- optimize scheduler ([#1106](https://github.com/CooooolFrog/dragonfly/issues/1106))
- reuse partial completed task ([#1107](https://github.com/CooooolFrog/dragonfly/issues/1107))
- optimize depth limit func ([#1102](https://github.com/CooooolFrog/dragonfly/issues/1102))
- change client default load limit ([#1104](https://github.com/CooooolFrog/dragonfly/issues/1104))
- limit tree depth ([#1099](https://github.com/CooooolFrog/dragonfly/issues/1099))
- update load limit ([#1097](https://github.com/CooooolFrog/dragonfly/issues/1097))

### Fix
- when peer state is PeerStateSucceeded, return size scope is small ([#1103](https://github.com/CooooolFrog/dragonfly/issues/1103))


<a name="v2.0.2-rc.14"></a>
## [v2.0.2-rc.14] - 2022-02-25
### Feat
- optimize peer range ([#1095](https://github.com/CooooolFrog/dragonfly/issues/1095))

### Test
- FilterParentCount ([#1094](https://github.com/CooooolFrog/dragonfly/issues/1094))


<a name="v2.0.2-rc.13"></a>
## [v2.0.2-rc.13] - 2022-02-24

<a name="v2.0.2-rc.12"></a>
## [v2.0.2-rc.12] - 2022-02-24
### Feat
- add cdn addresses log ([#1091](https://github.com/CooooolFrog/dragonfly/issues/1091))
- scheduler add limit count of filter parent func ([#1090](https://github.com/CooooolFrog/dragonfly/issues/1090))


<a name="v2.0.2-rc.11"></a>
## [v2.0.2-rc.11] - 2022-02-23
### Feat
- merge ranged request storage into parent ([#1078](https://github.com/CooooolFrog/dragonfly/issues/1078))
- add dynamic parallel count ([#1088](https://github.com/CooooolFrog/dragonfly/issues/1088))
- fix docker-compose ([#1087](https://github.com/CooooolFrog/dragonfly/issues/1087))

### Fix
- delete peer's parent on PeerEventDownloadSucceeded event ([#1085](https://github.com/CooooolFrog/dragonfly/issues/1085))


<a name="v2.0.2-rc.10"></a>
## [v2.0.2-rc.10] - 2022-02-22
### Chore
- update gorelease ldflags ([#1086](https://github.com/CooooolFrog/dragonfly/issues/1086))

### Feat
- add prefetch metric in client ([#1068](https://github.com/CooooolFrog/dragonfly/issues/1068))

### Test
- scheduler handle failed piece ([#1084](https://github.com/CooooolFrog/dragonfly/issues/1084))


<a name="v2.0.2-rc.9"></a>
## [v2.0.2-rc.9] - 2022-02-17
### Feat
- when scheduler blocks cdn, resource does not initialize cdn ([#1081](https://github.com/CooooolFrog/dragonfly/issues/1081))

### Fix
- pull request template typo ([#1080](https://github.com/CooooolFrog/dragonfly/issues/1080))


<a name="v2.0.2-rc.8"></a>
## [v2.0.2-rc.8] - 2022-02-17
### Docs
- add metrics document ([#1075](https://github.com/CooooolFrog/dragonfly/issues/1075))
- add containerd private registry configuration ([#1074](https://github.com/CooooolFrog/dragonfly/issues/1074))
- add containerd private registry configuration ([#1073](https://github.com/CooooolFrog/dragonfly/issues/1073))
- add docs about preheat console ([#1072](https://github.com/CooooolFrog/dragonfly/issues/1072))

### Feat
- scheduler blocks cdn ([#1079](https://github.com/CooooolFrog/dragonfly/issues/1079))
- job trigger cdn by resource ([#1076](https://github.com/CooooolFrog/dragonfly/issues/1076))


<a name="v2.0.2-rc.7"></a>
## [v2.0.2-rc.7] - 2022-02-15
### Feat
- add client request log ([#1069](https://github.com/CooooolFrog/dragonfly/issues/1069))

### Fix
- when cdn download failed, scheduler should set cdn peer state PeerStateFailed ([#1067](https://github.com/CooooolFrog/dragonfly/issues/1067))


<a name="v2.0.2-rc.6"></a>
## [v2.0.2-rc.6] - 2022-02-14

<a name="v2.0.2-rc.5"></a>
## [v2.0.2-rc.5] - 2022-02-14
### Chore
- enable range feature gate in e2e ([#1059](https://github.com/CooooolFrog/dragonfly/issues/1059))
- add content length for fast stream peer task ([#1061](https://github.com/CooooolFrog/dragonfly/issues/1061))
- optimize https pass through ([#1054](https://github.com/CooooolFrog/dragonfly/issues/1054))

### Docs
- manager installation ([#1063](https://github.com/CooooolFrog/dragonfly/issues/1063))

### Feat
- support change console log level ([#1055](https://github.com/CooooolFrog/dragonfly/issues/1055))

### Feature
- prefetch ranged requests ([#1053](https://github.com/CooooolFrog/dragonfly/issues/1053))
- support e2e feature gates ([#1056](https://github.com/CooooolFrog/dragonfly/issues/1056))

### Fix
- evaluate peer's parent ([#1064](https://github.com/CooooolFrog/dragonfly/issues/1064))
- scheduler download tiny file error ([#1052](https://github.com/CooooolFrog/dragonfly/issues/1052))


<a name="v2.0.2-rc.4"></a>
## [v2.0.2-rc.4] - 2022-01-29
### Feat
- manager support mysql ssl connection ([#1015](https://github.com/CooooolFrog/dragonfly/issues/1015))
- remove host and task when peer make tree ([#1042](https://github.com/CooooolFrog/dragonfly/issues/1042))
- cdn download tiny file ([#1040](https://github.com/CooooolFrog/dragonfly/issues/1040))

### Refactor
- manager grpc server ([#1047](https://github.com/CooooolFrog/dragonfly/issues/1047))
- scheduler grpc server ([#1046](https://github.com/CooooolFrog/dragonfly/issues/1046))


<a name="v2.0.2-rc.3"></a>
## [v2.0.2-rc.3] - 2022-01-25
### Chore
- use buildx to build docker images in e2e ([#1018](https://github.com/CooooolFrog/dragonfly/issues/1018))
- add missing pod log volumes in e2e ([#1037](https://github.com/CooooolFrog/dragonfly/issues/1037))
- upgrade to ginkgo v2 ([#1036](https://github.com/CooooolFrog/dragonfly/issues/1036))
- add piece task metrics in daemon ([#1030](https://github.com/CooooolFrog/dragonfly/issues/1030))

### Feat
- If cdn only updates IP, set cdn peers state to PeerStateLeave ([#1038](https://github.com/CooooolFrog/dragonfly/issues/1038))

### Fix
- docker actions typo ([#1041](https://github.com/CooooolFrog/dragonfly/issues/1041))
- cdn trigger peer error ([#1035](https://github.com/CooooolFrog/dragonfly/issues/1035))
- retrigger cdn panic ([#1034](https://github.com/CooooolFrog/dragonfly/issues/1034))

### Refactor
- docker workflows ([#1039](https://github.com/CooooolFrog/dragonfly/issues/1039))


<a name="v2.0.2-rc.2"></a>
## [v2.0.2-rc.2] - 2022-01-21
### Chore
- update outdated log ([#1028](https://github.com/CooooolFrog/dragonfly/issues/1028))
- optimize metrics and trace in daemon ([#1022](https://github.com/CooooolFrog/dragonfly/issues/1022))

### Feat
- generate grpc protoc ([#1027](https://github.com/CooooolFrog/dragonfly/issues/1027))
- manager config model add is_boot key ([#1025](https://github.com/CooooolFrog/dragonfly/issues/1025))
- scheduler download tiny file with range header ([#1024](https://github.com/CooooolFrog/dragonfly/issues/1024))

### Feature
- change log level in-flight ([#1023](https://github.com/CooooolFrog/dragonfly/issues/1023))


<a name="v2.0.2-rc.1"></a>
## [v2.0.2-rc.1] - 2022-01-20
### Feat
- change compatibility version to v2.0.2-rc.0 ([#1017](https://github.com/CooooolFrog/dragonfly/issues/1017))


<a name="v2.0.2-rc.0"></a>
## [v2.0.2-rc.0] - 2022-01-20
### Chore
- register to scheduler after updated running tasks ([#1016](https://github.com/CooooolFrog/dragonfly/issues/1016))


<a name="v2.0.2-beta.6"></a>
## [v2.0.2-beta.6] - 2022-01-20
### Feat
- when cdn peer is failed, peer should be back-to-source ([#1005](https://github.com/CooooolFrog/dragonfly/issues/1005))


<a name="v2.0.2-beta.3"></a>
## [v2.0.2-beta.3] - 2022-01-20
### Chore
- optimize defer and test ([#1010](https://github.com/CooooolFrog/dragonfly/issues/1010))
- workflow add test timeout ([#1011](https://github.com/CooooolFrog/dragonfly/issues/1011))
- sync docker-compose scheduler config ([#1001](https://github.com/CooooolFrog/dragonfly/issues/1001))
- parameterize tests in peer task ([#994](https://github.com/CooooolFrog/dragonfly/issues/994))
- clarify daemon interface ([#991](https://github.com/CooooolFrog/dragonfly/issues/991))

### Feat
- add actions job timout ([#1008](https://github.com/CooooolFrog/dragonfly/issues/1008))
- set peer state to running when scope size is SizeScope_TINY ([#1004](https://github.com/CooooolFrog/dragonfly/issues/1004))
- update submodule charts ([#1002](https://github.com/CooooolFrog/dragonfly/issues/1002))
- task mutex replace sync kmutex ([#1000](https://github.com/CooooolFrog/dragonfly/issues/1000))
- stream send error code ([#986](https://github.com/CooooolFrog/dragonfly/issues/986))
- trace https proxy request ([#996](https://github.com/CooooolFrog/dragonfly/issues/996))
- add scheduler host gc ([#989](https://github.com/CooooolFrog/dragonfly/issues/989))
- update typo in local_storage.go ([#955](https://github.com/CooooolFrog/dragonfly/issues/955))
- update charts submodule version ([#985](https://github.com/CooooolFrog/dragonfly/issues/985))

### Fix
- calculate piece MD5 sign when last piece download ([#1006](https://github.com/CooooolFrog/dragonfly/issues/1006))
- register task with size scope ([#1003](https://github.com/CooooolFrog/dragonfly/issues/1003))
- when scheduler is not available, replace the scheduler client ([#999](https://github.com/CooooolFrog/dragonfly/issues/999))
- total pieces count not set cause digest invalid ([#992](https://github.com/CooooolFrog/dragonfly/issues/992))
- send piece result error not handled ([#987](https://github.com/CooooolFrog/dragonfly/issues/987))


<a name="v2.0.2-beta.1"></a>
## [v2.0.2-beta.1] - 2022-01-12
### Feat
- change task and peer ttl ([#984](https://github.com/CooooolFrog/dragonfly/issues/984))
- when write last piece, generate digest ([#982](https://github.com/CooooolFrog/dragonfly/issues/982))

### Ffix
- typo in Makefile ([#975](https://github.com/CooooolFrog/dragonfly/issues/975))

### Fix
- scheduler config typo ([#983](https://github.com/CooooolFrog/dragonfly/issues/983))


<a name="v2.0.2-beta.0"></a>
## [v2.0.2-beta.0] - 2022-01-12
### Chore
- change docker.pkg.github.com to ghcr.io ([#973](https://github.com/CooooolFrog/dragonfly/issues/973))

### Feat
- merge same tasks in daemon ([#977](https://github.com/CooooolFrog/dragonfly/issues/977))
- if cdn is deleted, clear cdn related information ([#967](https://github.com/CooooolFrog/dragonfly/issues/967))

### Test
- dump goroutine in e2e ([#980](https://github.com/CooooolFrog/dragonfly/issues/980))


<a name="v2.0.2-alpha.7"></a>
## [v2.0.2-alpha.7] - 2021-12-31
### Chore
- copy e2e proxy log to artifact ([#962](https://github.com/CooooolFrog/dragonfly/issues/962))
- add version metric ([#954](https://github.com/CooooolFrog/dragonfly/issues/954))
- optimize back source update digest logic ([#950](https://github.com/CooooolFrog/dragonfly/issues/950))

### Docs
- update plugin doc ([#951](https://github.com/CooooolFrog/dragonfly/issues/951))

### Feat
- add default DiskGCThresholdPercent and ignore it when is 0 ([#971](https://github.com/CooooolFrog/dragonfly/issues/971))
- improve redirect to allow url rewrite ([#969](https://github.com/CooooolFrog/dragonfly/issues/969))
- Add useProxies to registryMirror allowing to mirror more anything ([#965](https://github.com/CooooolFrog/dragonfly/issues/965))
- change metrics port to 8000 ([#964](https://github.com/CooooolFrog/dragonfly/issues/964))
- add daemon metrics support ([#960](https://github.com/CooooolFrog/dragonfly/issues/960))
- support disk usage gc in client ([#953](https://github.com/CooooolFrog/dragonfly/issues/953))
- update source.Response and source client interface ([#945](https://github.com/CooooolFrog/dragonfly/issues/945))
- remove stat log from scheduler ([#946](https://github.com/CooooolFrog/dragonfly/issues/946))
- support recursive download in dfget ([#932](https://github.com/CooooolFrog/dragonfly/issues/932))
- add kmutex and krwmutex ([#934](https://github.com/CooooolFrog/dragonfly/issues/934))

### Fix
- schedulers send invalid direct piece ([#970](https://github.com/CooooolFrog/dragonfly/issues/970))
- use 'parent' as mainPeer in PeerPacket in removePeerFromCurrentTree() ([#957](https://github.com/CooooolFrog/dragonfly/issues/957))
- size scope empty ([#941](https://github.com/CooooolFrog/dragonfly/issues/941))
- not handle base.Code_SchedTaskStatusError in client ([#938](https://github.com/CooooolFrog/dragonfly/issues/938))
- infinitely get pieces when piece num is invalid ([#926](https://github.com/CooooolFrog/dragonfly/issues/926))


<a name="v2.0.2-alpha.6"></a>
## [v2.0.2-alpha.6] - 2021-12-15
### Feat
- make idgen package public ([#931](https://github.com/CooooolFrog/dragonfly/issues/931))
- make dfpath public ([#929](https://github.com/CooooolFrog/dragonfly/issues/929))

### Refactor
- scheduler register task ([#924](https://github.com/CooooolFrog/dragonfly/issues/924))


<a name="v2.0.2-alpha.5"></a>
## [v2.0.2-alpha.5] - 2021-12-13
### Docs
- update plugin docs ([#921](https://github.com/CooooolFrog/dragonfly/issues/921))

### Fix
- plugin dir is empty ([#922](https://github.com/CooooolFrog/dragonfly/issues/922))


<a name="v2.0.2-alpha.4"></a>
## [v2.0.2-alpha.4] - 2021-12-13
### Feat
- dfdaemon list scheduler cluster with multi idc ([#917](https://github.com/CooooolFrog/dragonfly/issues/917))
- update submodule ([#916](https://github.com/CooooolFrog/dragonfly/issues/916))
- update task access time ([#909](https://github.com/CooooolFrog/dragonfly/issues/909))
- optmize dfget package upgrade support ([#804](https://github.com/CooooolFrog/dragonfly/issues/804))
- support create container without docker-compose ([#915](https://github.com/CooooolFrog/dragonfly/issues/915))

### Fix
- peer gc ([#918](https://github.com/CooooolFrog/dragonfly/issues/918))
- go plugin test build error ([#912](https://github.com/CooooolFrog/dragonfly/issues/912))
- typo ([#911](https://github.com/CooooolFrog/dragonfly/issues/911))

### Refactor
- move from io/ioutil to io and os packages ([#906](https://github.com/CooooolFrog/dragonfly/issues/906))

### Test
- idgen peer id ([#913](https://github.com/CooooolFrog/dragonfly/issues/913))


<a name="v2.0.2-alpha.3"></a>
## [v2.0.2-alpha.3] - 2021-12-09
### Feat
- add data directory ([#910](https://github.com/CooooolFrog/dragonfly/issues/910))

### Fix
- total pieces not set when back source ([#908](https://github.com/CooooolFrog/dragonfly/issues/908))


<a name="v2.0.2-alpha.2"></a>
## [v2.0.2-alpha.2] - 2021-12-09
### Chore
- support multi daemons e2e test ([#896](https://github.com/CooooolFrog/dragonfly/issues/896))

### Docs
- dir path ([#904](https://github.com/CooooolFrog/dragonfly/issues/904))

### Feat
- add data storage directory  ([#907](https://github.com/CooooolFrog/dragonfly/issues/907))

### Fix
- mismatch digest peer task did not mark invalid ([#903](https://github.com/CooooolFrog/dragonfly/issues/903))


<a name="v2.0.2-alpha.1"></a>
## [v2.0.2-alpha.1] - 2021-12-08
### Fix
- dfget dfpath ([#901](https://github.com/CooooolFrog/dragonfly/issues/901))


<a name="v2.0.2-alpha.0"></a>
## [v2.0.2-alpha.0] - 2021-12-08
### Chore
- update UnknownSourceFileLen ([#888](https://github.com/CooooolFrog/dragonfly/issues/888))
- update changelog
- upgrade to golang 1.17 and alpine 3.14 ([#861](https://github.com/CooooolFrog/dragonfly/issues/861))

### Docs
- add plugin guide ([#875](https://github.com/CooooolFrog/dragonfly/issues/875))
- keep alive ([#868](https://github.com/CooooolFrog/dragonfly/issues/868))
- **zh-CN:** refactor machine translation ([#783](https://github.com/CooooolFrog/dragonfly/issues/783))

### Feat
- dfdaemon update content length ([#895](https://github.com/CooooolFrog/dragonfly/issues/895))
- lint sh ([#876](https://github.com/CooooolFrog/dragonfly/issues/876))
- update helm charts ([#870](https://github.com/CooooolFrog/dragonfly/issues/870))
- update version to v2.0.1 ([#869](https://github.com/CooooolFrog/dragonfly/issues/869))
- add oauth timeout ([#867](https://github.com/CooooolFrog/dragonfly/issues/867))
- support customize transport in daemon ([#866](https://github.com/CooooolFrog/dragonfly/issues/866))
- console ([#865](https://github.com/CooooolFrog/dragonfly/issues/865))
- move dfnet to internal ([#862](https://github.com/CooooolFrog/dragonfly/issues/862))
- remove ifaceutils pkg ([#860](https://github.com/CooooolFrog/dragonfly/issues/860))
- move syncmap pkg([#859](https://github.com/CooooolFrog/dragonfly/issues/859))
- oauth interface auth ([#857](https://github.com/CooooolFrog/dragonfly/issues/857))

### Fix
- scheduler success event ([#891](https://github.com/CooooolFrog/dragonfly/issues/891))
- add cdn cluster to scheduler cluster ([#887](https://github.com/CooooolFrog/dragonfly/issues/887))
- small size task failed due to digest error ([#886](https://github.com/CooooolFrog/dragonfly/issues/886))
- searcher log ([#878](https://github.com/CooooolFrog/dragonfly/issues/878))
- error log ([#863](https://github.com/CooooolFrog/dragonfly/issues/863))

### Refactor
- dfpath pkg ([#879](https://github.com/CooooolFrog/dragonfly/issues/879))


<a name="v2.0.1-rc.7"></a>
## [v2.0.1-rc.7] - 2021-12-02
### Docs
- update quick-start.md format ([#850](https://github.com/CooooolFrog/dragonfly/issues/850))

### Feat
- add scopes validation ([#856](https://github.com/CooooolFrog/dragonfly/issues/856))
- log ([#852](https://github.com/CooooolFrog/dragonfly/issues/852))

### Fix
- file peer task back source digest not match ([#849](https://github.com/CooooolFrog/dragonfly/issues/849))


<a name="v2.0.1-rc.6"></a>
## [v2.0.1-rc.6] - 2021-12-01
### Feat
- get scheduler list with advertise ip ([#848](https://github.com/CooooolFrog/dragonfly/issues/848))


<a name="v2.0.1-rc.5"></a>
## [v2.0.1-rc.5] - 2021-12-01
### Feat
- support mutli manager addrs ([#846](https://github.com/CooooolFrog/dragonfly/issues/846))


<a name="v2.0.1-rc.4"></a>
## [v2.0.1-rc.4] - 2021-12-01
### Feat
- searcher plugin change return params ([#844](https://github.com/CooooolFrog/dragonfly/issues/844))


<a name="v2.0.1-rc.2"></a>
## [v2.0.1-rc.2] - 2021-12-01
### Feat
- plugin log ([#843](https://github.com/CooooolFrog/dragonfly/issues/843))


<a name="v2.0.1-rc.1"></a>
## [v2.0.1-rc.1] - 2021-11-30
### Feat
- export searcher evaluate func ([#842](https://github.com/CooooolFrog/dragonfly/issues/842))
- add context for FindSchedulerCluster ([#841](https://github.com/CooooolFrog/dragonfly/issues/841))
- add application cdn clusters field ([#840](https://github.com/CooooolFrog/dragonfly/issues/840))


<a name="v2.0.1-rc.0"></a>
## [v2.0.1-rc.0] - 2021-11-30
### Feat
- update console submodule ([#838](https://github.com/CooooolFrog/dragonfly/issues/838))


<a name="v2.0.1-beta.6"></a>
## [v2.0.1-beta.6] - 2021-11-29
### Chore
- unify binary directory ([#828](https://github.com/CooooolFrog/dragonfly/issues/828))

### Feat
- preheat compatible with harbor ([#837](https://github.com/CooooolFrog/dragonfly/issues/837))
- gin version ([#833](https://github.com/CooooolFrog/dragonfly/issues/833))
- update manager image ([#831](https://github.com/CooooolFrog/dragonfly/issues/831))
- update helm charts version ([#824](https://github.com/CooooolFrog/dragonfly/issues/824))


<a name="v2.0.1-beta.5"></a>
## [v2.0.1-beta.5] - 2021-11-24
### Docs
- metrics configuration ([#816](https://github.com/CooooolFrog/dragonfly/issues/816))

### Feat
- add package reachable ([#822](https://github.com/CooooolFrog/dragonfly/issues/822))
- support list plugin ([#819](https://github.com/CooooolFrog/dragonfly/issues/819))
- scheduler and cdn report fqdn to manager ([#818](https://github.com/CooooolFrog/dragonfly/issues/818))


<a name="v2.0.1-beta.4"></a>
## [v2.0.1-beta.4] - 2021-11-22
### Docs
- manager apis ([#814](https://github.com/CooooolFrog/dragonfly/issues/814))

### Feat
- dfdaemon get scheduler list dynamically from manager ([#812](https://github.com/CooooolFrog/dragonfly/issues/812))

### Fix
- source plugin not loaded ([#811](https://github.com/CooooolFrog/dragonfly/issues/811))


<a name="v2.0.1-beta.3"></a>
## [v2.0.1-beta.3] - 2021-11-19
### Feat
- update image-spec version ([#808](https://github.com/CooooolFrog/dragonfly/issues/808))
- add security rule ([#806](https://github.com/CooooolFrog/dragonfly/issues/806))
- add idgen peer id ([#800](https://github.com/CooooolFrog/dragonfly/issues/800))

### Fix
- manager typo and cdn peer id ([#809](https://github.com/CooooolFrog/dragonfly/issues/809))

### Refactor
- scheduler evaluator ([#805](https://github.com/CooooolFrog/dragonfly/issues/805))


<a name="v2.0.1-beta.2"></a>
## [v2.0.1-beta.2] - 2021-11-15
### Chore
- add lint errcheck  and fix errcheck([#766](https://github.com/CooooolFrog/dragonfly/issues/766))
- optimize client storage gc log ([#790](https://github.com/CooooolFrog/dragonfly/issues/790))

### Feat
- optimize scheduler peer stat log ([#798](https://github.com/CooooolFrog/dragonfly/issues/798))
- replace sortedList with sortedUniqueList ([#793](https://github.com/CooooolFrog/dragonfly/issues/793))

### Test
- preheat image ([#794](https://github.com/CooooolFrog/dragonfly/issues/794))


<a name="v2.0.1-beta.1"></a>
## [v2.0.1-beta.1] - 2021-11-10
### Feat
- calculate piece metadata digest ([#787](https://github.com/CooooolFrog/dragonfly/issues/787))


<a name="v2.0.1-alpha.10"></a>
## [v2.0.1-alpha.10] - 2021-11-09
### Chore
- add markdown lint ([#779](https://github.com/CooooolFrog/dragonfly/issues/779))
- update golang import lint ([#780](https://github.com/CooooolFrog/dragonfly/issues/780))

### Docs
- manager api ([#774](https://github.com/CooooolFrog/dragonfly/issues/774))
- **zh:** add zh docs ([#777](https://github.com/CooooolFrog/dragonfly/issues/777))

### Feat
- preheat skip certificate validation ([#786](https://github.com/CooooolFrog/dragonfly/issues/786))
- support traffic metrics by peer host ([#776](https://github.com/CooooolFrog/dragonfly/issues/776))

### Fix
- cdn AdvertiseIP not used ([#782](https://github.com/CooooolFrog/dragonfly/issues/782))

### Test
- scheduler supervisor ([#742](https://github.com/CooooolFrog/dragonfly/issues/742))


<a name="v2.0.1-alpha.8"></a>
## [v2.0.1-alpha.8] - 2021-10-29
### Chore
- optimize stream peer task ([#763](https://github.com/CooooolFrog/dragonfly/issues/763))

### Feat
- support dump http content in client for debugging ([#770](https://github.com/CooooolFrog/dragonfly/issues/770))
- remove calculate total count service ([#772](https://github.com/CooooolFrog/dragonfly/issues/772))
- add user list interface ([#771](https://github.com/CooooolFrog/dragonfly/issues/771))
- clear hashcircler and maputils package ([#768](https://github.com/CooooolFrog/dragonfly/issues/768))


<a name="v2.0.1-alpha.7"></a>
## [v2.0.1-alpha.7] - 2021-10-28
### Fix
- add peer to task failed because InnerBucketMaxLength is small ([#765](https://github.com/CooooolFrog/dragonfly/issues/765))


<a name="v2.0.1-alpha.6"></a>
## [v2.0.1-alpha.6] - 2021-10-28
### Chore
- check empty registry mirror ([#761](https://github.com/CooooolFrog/dragonfly/issues/761))

### Feat
- add cdn task peers monitor log ([#764](https://github.com/CooooolFrog/dragonfly/issues/764))
- change config key name ([#759](https://github.com/CooooolFrog/dragonfly/issues/759))

### Fix
- back source weight ([#762](https://github.com/CooooolFrog/dragonfly/issues/762))


<a name="v2.0.1-alpha.5"></a>
## [v2.0.1-alpha.5] - 2021-10-27
### Feat
- scheduler channel blocking ([#756](https://github.com/CooooolFrog/dragonfly/issues/756))


<a name="v2.0.1-alpha.4"></a>
## [v2.0.1-alpha.4] - 2021-10-26
### Chore
- optimize span context for report ([#747](https://github.com/CooooolFrog/dragonfly/issues/747))

### Docs
- add maxConcurrency comment ([#755](https://github.com/CooooolFrog/dragonfly/issues/755))
- add troubleshooting guide ([#752](https://github.com/CooooolFrog/dragonfly/issues/752))
- add load limit ([#745](https://github.com/CooooolFrog/dragonfly/issues/745))
- **en:** upgrade docs ([#673](https://github.com/CooooolFrog/dragonfly/issues/673))
- **runtime:** upgrade containerd runtime ([#748](https://github.com/CooooolFrog/dragonfly/issues/748))

### Feat
- add jobs api ([#751](https://github.com/CooooolFrog/dragonfly/issues/751))
- add config ([#746](https://github.com/CooooolFrog/dragonfly/issues/746))
- add preheat otel ([#741](https://github.com/CooooolFrog/dragonfly/issues/741))

### Fix
- client load ([#753](https://github.com/CooooolFrog/dragonfly/issues/753))


<a name="v2.0.1-alpha.3"></a>
## [v2.0.1-alpha.3] - 2021-10-20
### Feat
- add job logger ([#740](https://github.com/CooooolFrog/dragonfly/issues/740))


<a name="v2.0.1-alpha.2"></a>
## [v2.0.1-alpha.2] - 2021-10-20
### Feat
- manager add grpc jaeger ([#738](https://github.com/CooooolFrog/dragonfly/issues/738))
- load limit ([#739](https://github.com/CooooolFrog/dragonfly/issues/739))
- preheat cluster ([#731](https://github.com/CooooolFrog/dragonfly/issues/731))
- nsswitch ([#737](https://github.com/CooooolFrog/dragonfly/issues/737))
- export e2e logs ([#732](https://github.com/CooooolFrog/dragonfly/issues/732))


<a name="v2.0.1-alpha.1"></a>
## [v2.0.1-alpha.1] - 2021-10-13
### Chore
- repository name
- change docker registry name ([#725](https://github.com/CooooolFrog/dragonfly/issues/725))
- update config example ([#721](https://github.com/CooooolFrog/dragonfly/issues/721))
- release image to docker.pkg.github.com ([#703](https://github.com/CooooolFrog/dragonfly/issues/703))

### Docs
- update kubernetes docs ([#714](https://github.com/CooooolFrog/dragonfly/issues/714))
- add apis and preheat ([#712](https://github.com/CooooolFrog/dragonfly/issues/712))
- update kubernetes docs ([#705](https://github.com/CooooolFrog/dragonfly/issues/705))

### Feat
- compatible with V1 preheat  ([#720](https://github.com/CooooolFrog/dragonfly/issues/720))
- add grpc metric and refactor grpc server ([#686](https://github.com/CooooolFrog/dragonfly/issues/686))

### Fix
- peer empty parent ([#724](https://github.com/CooooolFrog/dragonfly/issues/724))
- client panic ([#719](https://github.com/CooooolFrog/dragonfly/issues/719))
- client goroutine and fd leak ([#713](https://github.com/CooooolFrog/dragonfly/issues/713))


<a name="v2.0.1-alpha.0"></a>
## [v2.0.1-alpha.0] - 2021-09-29
### Chore
- workflows ignore paths ([#697](https://github.com/CooooolFrog/dragonfly/issues/697))
- remove skip-duplicate-actions ([#690](https://github.com/CooooolFrog/dragonfly/issues/690))
- e2e workflows remove goproxy ([#677](https://github.com/CooooolFrog/dragonfly/issues/677))

### Docs
- scheduler config ([#698](https://github.com/CooooolFrog/dragonfly/issues/698))
- update kubernetes docs ([#696](https://github.com/CooooolFrog/dragonfly/issues/696))

### Feat
- add manager client list scheduler interface ([#694](https://github.com/CooooolFrog/dragonfly/issues/694))

### Fix
- skip check DisableAutoBackSource option when scheduler says back source ([#693](https://github.com/CooooolFrog/dragonfly/issues/693))

### Refactor
- scheduler supervisor ([#655](https://github.com/CooooolFrog/dragonfly/issues/655))


<a name="v2.0.1-a-rc2"></a>
## [v2.0.1-a-rc2] - 2021-09-23
### Chore
- export set log level ([#646](https://github.com/CooooolFrog/dragonfly/issues/646))
- enable calculate digest ([#656](https://github.com/CooooolFrog/dragonfly/issues/656))
- update build package config ([#653](https://github.com/CooooolFrog/dragonfly/issues/653))
- optimize advertise ip ([#652](https://github.com/CooooolFrog/dragonfly/issues/652))
- change zzy987 maintainers email ([#649](https://github.com/CooooolFrog/dragonfly/issues/649))
- update version ([#647](https://github.com/CooooolFrog/dragonfly/issues/647))
- add lucy-cl maintainer ([#645](https://github.com/CooooolFrog/dragonfly/issues/645))
- makefile typo

### Docs
- scheduler config ([#654](https://github.com/CooooolFrog/dragonfly/issues/654))

### Feat
- release fd ([#668](https://github.com/CooooolFrog/dragonfly/issues/668))
- add otel trace ([#650](https://github.com/CooooolFrog/dragonfly/issues/650))
- disable prepared statement ([#648](https://github.com/CooooolFrog/dragonfly/issues/648))

### Fix
- go library cve ([#666](https://github.com/CooooolFrog/dragonfly/issues/666))


<a name="v2.0.0"></a>
## [v2.0.0] - 2021-09-09
### Chore
- compatibility with v2.0.0 test ([#639](https://github.com/CooooolFrog/dragonfly/issues/639))
- skip e2e ([#631](https://github.com/CooooolFrog/dragonfly/issues/631))
- rename cdnsystem to cdn ([#626](https://github.com/CooooolFrog/dragonfly/issues/626))
- skip workflows ([#624](https://github.com/CooooolFrog/dragonfly/issues/624))
- update changelog ([#622](https://github.com/CooooolFrog/dragonfly/issues/622))
- update submodule version ([#608](https://github.com/CooooolFrog/dragonfly/issues/608))
- optimize app and tracer log ([#607](https://github.com/CooooolFrog/dragonfly/issues/607))

### Docs
- maintainers ([#636](https://github.com/CooooolFrog/dragonfly/issues/636))
- test guide link ([#635](https://github.com/CooooolFrog/dragonfly/issues/635))
- add manager preview ([#634](https://github.com/CooooolFrog/dragonfly/issues/634))
- install ([#628](https://github.com/CooooolFrog/dragonfly/issues/628))
- update document ([#625](https://github.com/CooooolFrog/dragonfly/issues/625))
- update docs/zh-CN/config/dfget.yaml ([#623](https://github.com/CooooolFrog/dragonfly/issues/623))
- Update documents ([#595](https://github.com/CooooolFrog/dragonfly/issues/595))
- update runtime guide in helm deploy ([#612](https://github.com/CooooolFrog/dragonfly/issues/612))

### Feat
- update verison ([#640](https://github.com/CooooolFrog/dragonfly/issues/640))
- changelog ([#638](https://github.com/CooooolFrog/dragonfly/issues/638))
- update console submodule ([#637](https://github.com/CooooolFrog/dragonfly/issues/637))
- update submodule ([#632](https://github.com/CooooolFrog/dragonfly/issues/632))
- beautify scheduler & CDN log ([#618](https://github.com/CooooolFrog/dragonfly/issues/618))
- Print version information when the system starts up ([#620](https://github.com/CooooolFrog/dragonfly/issues/620))
- add piece download timeout ([#621](https://github.com/CooooolFrog/dragonfly/issues/621))
- notice client back source when rescheduled parent reach max times ([#611](https://github.com/CooooolFrog/dragonfly/issues/611))
- avoid report peer result fail due to context cancel & add backsource tracer ([#606](https://github.com/CooooolFrog/dragonfly/issues/606))
- optimize cdn check free space ([#603](https://github.com/CooooolFrog/dragonfly/issues/603))

### Feature
- refresh proto file ([#615](https://github.com/CooooolFrog/dragonfly/issues/615))

### Fix
- return failed piece ([#619](https://github.com/CooooolFrog/dragonfly/issues/619))

### Test
- preheat e2e ([#627](https://github.com/CooooolFrog/dragonfly/issues/627))


<a name="v0.5.0"></a>
## [v0.5.0] - 2021-09-06
### Chore
- add compatibility test workflow ([#594](https://github.com/CooooolFrog/dragonfly/issues/594))

### Feat
- client back source ([#579](https://github.com/CooooolFrog/dragonfly/issues/579))
- enable manager user's features ([#598](https://github.com/CooooolFrog/dragonfly/issues/598))
- add sni proxy support ([#600](https://github.com/CooooolFrog/dragonfly/issues/600))
- compatibility e2e with matrix ([#599](https://github.com/CooooolFrog/dragonfly/issues/599))

### Fix
- use string slice for header ([#601](https://github.com/CooooolFrog/dragonfly/issues/601))
- preheat-e2e timeout ([#602](https://github.com/CooooolFrog/dragonfly/issues/602))


<a name="v0.4.0"></a>
## [v0.4.0] - 2021-09-02
### Chore
- add copyright ([#593](https://github.com/CooooolFrog/dragonfly/issues/593))

### Docs
- rbac swagger comment

### Feat
- change scheduler cluster query params ([#596](https://github.com/CooooolFrog/dragonfly/issues/596))
- add oauth2 signin ([#591](https://github.com/CooooolFrog/dragonfly/issues/591))
- update scheduler cluster query params ([#587](https://github.com/CooooolFrog/dragonfly/issues/587))
- add time out when register ([#588](https://github.com/CooooolFrog/dragonfly/issues/588))
- skip verify when back to source ([#586](https://github.com/CooooolFrog/dragonfly/issues/586))
- update charts submodule ([#583](https://github.com/CooooolFrog/dragonfly/issues/583))
- support limit from dfget client ([#578](https://github.com/CooooolFrog/dragonfly/issues/578))

### Refactor
- rbac
- user interface

### Test
- print merge commit ([#581](https://github.com/CooooolFrog/dragonfly/issues/581))


<a name="v0.3.0"></a>
## [v0.3.0] - 2021-08-25
### Feat
- add cdn cluster id for scheduler cluster ([#580](https://github.com/CooooolFrog/dragonfly/issues/580))
- start process ([#572](https://github.com/CooooolFrog/dragonfly/issues/572))
- gin log to file ([#574](https://github.com/CooooolFrog/dragonfly/issues/574))
- add manager cors middleware ([#573](https://github.com/CooooolFrog/dragonfly/issues/573))
- change rabc code struct ([#552](https://github.com/CooooolFrog/dragonfly/issues/552))

### Fix
- use getTask instead of taskStore.Get, for the error cause type ([#571](https://github.com/CooooolFrog/dragonfly/issues/571))


<a name="v0.2.0"></a>
## v0.2.0 - 2021-08-20
### Chore
- rename cdn server package to rpcserver ([#554](https://github.com/CooooolFrog/dragonfly/issues/554))
- optimize peer task report function ([#543](https://github.com/CooooolFrog/dragonfly/issues/543))
- optimize client rpc package name and other docs ([#541](https://github.com/CooooolFrog/dragonfly/issues/541))
- optimize grpc interceptor code ([#536](https://github.com/CooooolFrog/dragonfly/issues/536))
- optimize compute piece size function ([#528](https://github.com/CooooolFrog/dragonfly/issues/528))
- set GOPROXY with default value ([#463](https://github.com/CooooolFrog/dragonfly/issues/463))
- custom charts template namespace ([#416](https://github.com/CooooolFrog/dragonfly/issues/416))
- remove goreleaser go generate ([#409](https://github.com/CooooolFrog/dragonfly/issues/409))
- rename dfdaemon docker image ([#405](https://github.com/CooooolFrog/dragonfly/issues/405))
- remove macos ci ([#404](https://github.com/CooooolFrog/dragonfly/issues/404))
- add docs for dragonfly2.0 ([#234](https://github.com/CooooolFrog/dragonfly/issues/234))
- change bash to sh ([#383](https://github.com/CooooolFrog/dragonfly/issues/383))
- remove protoc.sh ([#341](https://github.com/CooooolFrog/dragonfly/issues/341))
- update CI timeout ([#328](https://github.com/CooooolFrog/dragonfly/issues/328))
- remove build script's git operation ([#321](https://github.com/CooooolFrog/dragonfly/issues/321))
- docker building workflow ([#323](https://github.com/CooooolFrog/dragonfly/issues/323))
- remove manager netcat-openbsd ([#298](https://github.com/CooooolFrog/dragonfly/issues/298))
- workflows remove main-rc branch ([#221](https://github.com/CooooolFrog/dragonfly/issues/221))
- change manager swagger docs path and add makefile swagger command ([#183](https://github.com/CooooolFrog/dragonfly/issues/183))
- add SECURITY.md ([#181](https://github.com/CooooolFrog/dragonfly/issues/181))
- change codeowners ([#179](https://github.com/CooooolFrog/dragonfly/issues/179))
- change codeowners to dragonfly2's maintainers and reviewers ([#169](https://github.com/CooooolFrog/dragonfly/issues/169))
- create custom issue template ([#168](https://github.com/CooooolFrog/dragonfly/issues/168))
- add pull request and issue templates ([#154](https://github.com/CooooolFrog/dragonfly/issues/154))

### Daemon
- add add additional peer id for some logs ([#205](https://github.com/CooooolFrog/dragonfly/issues/205))
- create output parent directory if not exists ([#188](https://github.com/CooooolFrog/dragonfly/issues/188))
- update default timeout and add context for downloading piece ([#190](https://github.com/CooooolFrog/dragonfly/issues/190))
- record failed code when unfinished and event for scheduler ([#176](https://github.com/CooooolFrog/dragonfly/issues/176))

### Docs
- Add dfget man page ([#388](https://github.com/CooooolFrog/dragonfly/issues/388))
- update v0.1.0-beta changelog ([#387](https://github.com/CooooolFrog/dragonfly/issues/387))
- add CHANGELOG.md
- add CODE_OF_CONDUCT.md ([#163](https://github.com/CooooolFrog/dragonfly/issues/163))

### Feat
- empty scheduler job ([#565](https://github.com/CooooolFrog/dragonfly/issues/565))
- optimize manager startup process ([#562](https://github.com/CooooolFrog/dragonfly/issues/562))
- update git submodule ([#560](https://github.com/CooooolFrog/dragonfly/issues/560))
- optimize scheduler start server ([#558](https://github.com/CooooolFrog/dragonfly/issues/558))
- add console ([#559](https://github.com/CooooolFrog/dragonfly/issues/559))
- generate swagger api ([#557](https://github.com/CooooolFrog/dragonfly/issues/557))
- add console submodule ([#549](https://github.com/CooooolFrog/dragonfly/issues/549))
- optimize get permission name ([#548](https://github.com/CooooolFrog/dragonfly/issues/548))
- rename task to job ([#544](https://github.com/CooooolFrog/dragonfly/issues/544))
- Add distribute Schedule Tracer & Refactor scheduler ([#537](https://github.com/CooooolFrog/dragonfly/issues/537))
- add artifacthub badge ([#524](https://github.com/CooooolFrog/dragonfly/issues/524))
- update cdn host ([#530](https://github.com/CooooolFrog/dragonfly/issues/530))
- back source when no available peers or scheduler error ([#521](https://github.com/CooooolFrog/dragonfly/issues/521))
- add task manager ([#490](https://github.com/CooooolFrog/dragonfly/issues/490))
- rename manager grpc ([#510](https://github.com/CooooolFrog/dragonfly/issues/510))
- Add stress testing tool for daemon ([#506](https://github.com/CooooolFrog/dragonfly/issues/506))
- scheduler getevaluator lock ([#502](https://github.com/CooooolFrog/dragonfly/issues/502))
- rename search file to searcher ([#484](https://github.com/CooooolFrog/dragonfly/issues/484))
- Add schedule log ([#495](https://github.com/CooooolFrog/dragonfly/issues/495))
- Extract peer event processing function ([#489](https://github.com/CooooolFrog/dragonfly/issues/489))
- optimize scheduler dynconfig ([#480](https://github.com/CooooolFrog/dragonfly/issues/480))
- optimize jwt ([#476](https://github.com/CooooolFrog/dragonfly/issues/476))
- register service to manager ([#475](https://github.com/CooooolFrog/dragonfly/issues/475))
- add searcher to scheduler cluster ([#462](https://github.com/CooooolFrog/dragonfly/issues/462))
- CDN implementation supports HDFS type storage ([#420](https://github.com/CooooolFrog/dragonfly/issues/420))
- add is_default to scheduler_cluster table ([#458](https://github.com/CooooolFrog/dragonfly/issues/458))
- add host info for scheduler and cdn ([#457](https://github.com/CooooolFrog/dragonfly/issues/457))
- Install e2e script ([#451](https://github.com/CooooolFrog/dragonfly/issues/451))
- Manager user logic ([#419](https://github.com/CooooolFrog/dragonfly/issues/419))
- Add plugin support for resource ([#291](https://github.com/CooooolFrog/dragonfly/issues/291))
- changelog ([#326](https://github.com/CooooolFrog/dragonfly/issues/326))
- remove queue package ([#275](https://github.com/CooooolFrog/dragonfly/issues/275))
- add ci badge ([#265](https://github.com/CooooolFrog/dragonfly/issues/265))
- remove slidingwindow and assertutils package ([#263](https://github.com/CooooolFrog/dragonfly/issues/263))

### Feature
- update helm charts submodule ([#567](https://github.com/CooooolFrog/dragonfly/issues/567))
- Add manager charts with submodule ([#525](https://github.com/CooooolFrog/dragonfly/issues/525))
- support mysql 5.6 ([#520](https://github.com/CooooolFrog/dragonfly/issues/520))
- support customize base image ([#519](https://github.com/CooooolFrog/dragonfly/issues/519))
- add kustomize yaml for deploying ([#349](https://github.com/CooooolFrog/dragonfly/issues/349))
- support basic auth for proxy ([#250](https://github.com/CooooolFrog/dragonfly/issues/250))
- add disk quota gc for daemon ([#215](https://github.com/CooooolFrog/dragonfly/issues/215))

### Feature
- optimize manager project layout ([#540](https://github.com/CooooolFrog/dragonfly/issues/540))
- enable grpc tracing ([#531](https://github.com/CooooolFrog/dragonfly/issues/531))
- remove proto redundant field ([#508](https://github.com/CooooolFrog/dragonfly/issues/508))
- update multiple registries support docs ([#481](https://github.com/CooooolFrog/dragonfly/issues/481))
- add multiple registry mirrors support ([#479](https://github.com/CooooolFrog/dragonfly/issues/479))
- disable proxy when config is empty ([#455](https://github.com/CooooolFrog/dragonfly/issues/455))
- add pod labels in helm chart ([#447](https://github.com/CooooolFrog/dragonfly/issues/447))
- optimize failed reason not set ([#446](https://github.com/CooooolFrog/dragonfly/issues/446))
- report peer result when failed to register ([#433](https://github.com/CooooolFrog/dragonfly/issues/433))
- rename PeerHost to Daemon in client ([#438](https://github.com/CooooolFrog/dragonfly/issues/438))
- move internal/rpc to pkg/rpc ([#436](https://github.com/CooooolFrog/dragonfly/issues/436))
- export peer.TaskManager for embedding dragonfly in custom binary ([#434](https://github.com/CooooolFrog/dragonfly/issues/434))
- optimize error message for proxy ([#428](https://github.com/CooooolFrog/dragonfly/issues/428))
- minimize daemon runtime capabilities ([#421](https://github.com/CooooolFrog/dragonfly/issues/421))
- add default filter in proxy for deployment and docs ([#417](https://github.com/CooooolFrog/dragonfly/issues/417))
- add jaeger for helm deployment ([#415](https://github.com/CooooolFrog/dragonfly/issues/415))
- update dfdaemon proxy port comment
- update cdn init container template ([#399](https://github.com/CooooolFrog/dragonfly/issues/399))
- update client config to Camel-Case format ([#393](https://github.com/CooooolFrog/dragonfly/issues/393))
- update helm charts deploy guide ([#386](https://github.com/CooooolFrog/dragonfly/issues/386))
- update helm charts ([#385](https://github.com/CooooolFrog/dragonfly/issues/385))
- support setns in client ([#378](https://github.com/CooooolFrog/dragonfly/issues/378))
- disable resolver server config ([#314](https://github.com/CooooolFrog/dragonfly/issues/314))
- update docs ([#307](https://github.com/CooooolFrog/dragonfly/issues/307))
- remove unsafe code in client/daemon/storage ([#258](https://github.com/CooooolFrog/dragonfly/issues/258))
- remove redundant configurations ([#216](https://github.com/CooooolFrog/dragonfly/issues/216))

### Fix
- wrong daemon config and kubectl image tag ([#398](https://github.com/CooooolFrog/dragonfly/issues/398))
- scheduler concurrent dead lock ([#509](https://github.com/CooooolFrog/dragonfly/issues/509))
- adjust dfget download log ([#564](https://github.com/CooooolFrog/dragonfly/issues/564))
- cdn back source range size overflow ([#550](https://github.com/CooooolFrog/dragonfly/issues/550))
- update mapsturcture decode and remove unused config ([#396](https://github.com/CooooolFrog/dragonfly/issues/396))
- scheduler pick candidate and associate child  encounter  dead lock ([#500](https://github.com/CooooolFrog/dragonfly/issues/500))
- generate proto file ([#483](https://github.com/CooooolFrog/dragonfly/issues/483))
- address typo ([#468](https://github.com/CooooolFrog/dragonfly/issues/468))
- dead lock when pt.failedPieceCh is full ([#466](https://github.com/CooooolFrog/dragonfly/issues/466))
- user table typo ([#453](https://github.com/CooooolFrog/dragonfly/issues/453))
- log specification ([#452](https://github.com/CooooolFrog/dragonfly/issues/452))
- update DynconfigOptions typo ([#390](https://github.com/CooooolFrog/dragonfly/issues/390))
- close net namespace fd ([#418](https://github.com/CooooolFrog/dragonfly/issues/418))
- update static cdn config
- wrong content length in proxy
- wait available peer packet panic ([#561](https://github.com/CooooolFrog/dragonfly/issues/561))
- wrong cache header ([#423](https://github.com/CooooolFrog/dragonfly/issues/423))
- gc test ([#370](https://github.com/CooooolFrog/dragonfly/issues/370))
- scheduler panic ([#356](https://github.com/CooooolFrog/dragonfly/issues/356))
- use seederName to replace the PeerID to generate the UUID ([#355](https://github.com/CooooolFrog/dragonfly/issues/355))
- check health too long when dfdaemon is unavailable ([#344](https://github.com/CooooolFrog/dragonfly/issues/344))
- change manager docs path ([#193](https://github.com/CooooolFrog/dragonfly/issues/193))
- when load config from cdn directory in dynconfig, skip sub directories ([#310](https://github.com/CooooolFrog/dragonfly/issues/310))
- Makefile and build.sh ([#309](https://github.com/CooooolFrog/dragonfly/issues/309))
- ci badge ([#281](https://github.com/CooooolFrog/dragonfly/issues/281))
- change peerPacketReady to buffer channel ([#256](https://github.com/CooooolFrog/dragonfly/issues/256))
- cdn gc dead lock ([#231](https://github.com/CooooolFrog/dragonfly/issues/231))
- cfgFile nil error ([#224](https://github.com/CooooolFrog/dragonfly/issues/224))
- **manager:** modify to config from scheduler_config in swagger yaml ([#317](https://github.com/CooooolFrog/dragonfly/issues/317))

### Fix
- proxy for stress testing tool ([#507](https://github.com/CooooolFrog/dragonfly/issues/507))
- add process level for scheduler peer task status ([#435](https://github.com/CooooolFrog/dragonfly/issues/435))
- infinite recursion in MkDirAll ([#358](https://github.com/CooooolFrog/dragonfly/issues/358))
- use atomic to avoid data race in client ([#254](https://github.com/CooooolFrog/dragonfly/issues/254))

### Refactor
- manager server new instance ([#464](https://github.com/CooooolFrog/dragonfly/issues/464))
- update arch ([#319](https://github.com/CooooolFrog/dragonfly/issues/319))
- remove benchmark-rate and rename not-back-source ([#245](https://github.com/CooooolFrog/dragonfly/issues/245))
- support multi digest not only md5 ([#236](https://github.com/CooooolFrog/dragonfly/issues/236))
- simplify to make imports more format ([#230](https://github.com/CooooolFrog/dragonfly/issues/230))
- **manager:** modify mysql table schema, orm json tag. ([#283](https://github.com/CooooolFrog/dragonfly/issues/283))

### Test
- compare image commit ([#538](https://github.com/CooooolFrog/dragonfly/issues/538))
- E2E download concurrency ([#467](https://github.com/CooooolFrog/dragonfly/issues/467))
- E2E test use kind's containerd ([#448](https://github.com/CooooolFrog/dragonfly/issues/448))
- manager config ([#392](https://github.com/CooooolFrog/dragonfly/issues/392))
- idgen add digest ([#243](https://github.com/CooooolFrog/dragonfly/issues/243))


[Unreleased]: https://github.com/CooooolFrog/dragonfly/compare/v2.3.1-beta.0...HEAD
[v2.3.1-beta.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.3.0...v2.3.1-beta.0
[v2.3.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.4-rc.2...v2.3.0
[v2.2.4-rc.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.4-rc.1...v2.2.4-rc.2
[v2.2.4-rc.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.4-rc.0...v2.2.4-rc.1
[v2.2.4-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.4-beta.1...v2.2.4-rc.0
[v2.2.4-beta.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.4-beta.0...v2.2.4-beta.1
[v2.2.4-beta.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.3...v2.2.4-beta.0
[v2.2.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.3-rc.4...v2.2.3
[v2.2.3-rc.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.3-rc.3...v2.2.3-rc.4
[v2.2.3-rc.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.3-rc.2...v2.2.3-rc.3
[v2.2.3-rc.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.3-rc.1...v2.2.3-rc.2
[v2.2.3-rc.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.3-rc.0...v2.2.3-rc.1
[v2.2.3-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.2...v2.2.3-rc.0
[v2.2.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.2-rc.0...v2.2.2
[v2.2.2-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.1...v2.2.2-rc.0
[v2.2.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.1-rc.3...v2.2.1
[v2.2.1-rc.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.1-rc.1...v2.2.1-rc.3
[v2.2.1-rc.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.1-rc.0...v2.2.1-rc.1
[v2.2.1-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.2.0...v2.2.1-rc.0
[v2.2.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.67...v2.2.0
[v2.1.67]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.66...v2.1.67
[v2.1.66]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.65...v2.1.66
[v2.1.65]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.64...v2.1.65
[v2.1.64]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.63...v2.1.64
[v2.1.63]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.62...v2.1.63
[v2.1.62]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.61...v2.1.62
[v2.1.61]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.60...v2.1.61
[v2.1.60]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.59...v2.1.60
[v2.1.59]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.58...v2.1.59
[v2.1.58]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.57...v2.1.58
[v2.1.57]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.56...v2.1.57
[v2.1.56]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.55...v2.1.56
[v2.1.55]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.54...v2.1.55
[v2.1.54]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.53...v2.1.54
[v2.1.53]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.52...v2.1.53
[v2.1.52]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.51...v2.1.52
[v2.1.51]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.50...v2.1.51
[v2.1.50]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.49...v2.1.50
[v2.1.49]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.48...v2.1.49
[v2.1.48]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.47...v2.1.48
[v2.1.47]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.46...v2.1.47
[v2.1.46]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.45...v2.1.46
[v2.1.45]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.44...v2.1.45
[v2.1.44]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.43...v2.1.44
[v2.1.43]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.42...v2.1.43
[v2.1.42]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.41...v2.1.42
[v2.1.41]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.40...v2.1.41
[v2.1.40]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.39...v2.1.40
[v2.1.39]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.38...v2.1.39
[v2.1.38]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.37...v2.1.38
[v2.1.37]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.36...v2.1.37
[v2.1.36]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.35...v2.1.36
[v2.1.35]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.34...v2.1.35
[v2.1.34]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.33...v2.1.34
[v2.1.33]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.32...v2.1.33
[v2.1.32]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.31...v2.1.32
[v2.1.31]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.30...v2.1.31
[v2.1.30]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.29...v2.1.30
[v2.1.29]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.28...v2.1.29
[v2.1.28]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.27...v2.1.28
[v2.1.27]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.26...v2.1.27
[v2.1.26]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.25...v2.1.26
[v2.1.25]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.24...v2.1.25
[v2.1.24]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.23...v2.1.24
[v2.1.23]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.22...v2.1.23
[v2.1.22]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.21...v2.1.22
[v2.1.21]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.20...v2.1.21
[v2.1.20]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.19...v2.1.20
[v2.1.19]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.18...v2.1.19
[v2.1.18]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.17...v2.1.18
[v2.1.17]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.16...v2.1.17
[v2.1.16]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.15...v2.1.16
[v2.1.15]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.14...v2.1.15
[v2.1.14]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.13...v2.1.14
[v2.1.13]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.12...v2.1.13
[v2.1.12]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.11...v2.1.12
[v2.1.11]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.10...v2.1.11
[v2.1.10]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.9...v2.1.10
[v2.1.9]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.8...v2.1.9
[v2.1.8]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.7...v2.1.8
[v2.1.7]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.6...v2.1.7
[v2.1.6]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.5...v2.1.6
[v2.1.5]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.4...v2.1.5
[v2.1.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.3...v2.1.4
[v2.1.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.2...v2.1.3
[v2.1.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.1...v2.1.2
[v2.1.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.0...v2.1.1
[v2.1.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.0-rc.0...v2.1.0
[v2.1.0-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.0-beta.4...v2.1.0-rc.0
[v2.1.0-beta.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.0-beta.3...v2.1.0-beta.4
[v2.1.0-beta.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.0-beta.2...v2.1.0-beta.3
[v2.1.0-beta.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.0-beta.1...v2.1.0-beta.2
[v2.1.0-beta.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.0-alpha.3...v2.1.0-beta.1
[v2.1.0-alpha.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.0-alpha.2...v2.1.0-alpha.3
[v2.1.0-alpha.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.0-alpha.1...v2.1.0-alpha.2
[v2.1.0-alpha.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.1.0-alpha.0...v2.1.0-alpha.1
[v2.1.0-alpha.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.10...v2.1.0-alpha.0
[v2.0.10]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-rc.2...v2.0.10
[v2.0.9-rc.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-rc.1...v2.0.9-rc.2
[v2.0.9-rc.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-rc.0...v2.0.9-rc.1
[v2.0.9-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-beta.4...v2.0.9-rc.0
[v2.0.9-beta.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-beta.3...v2.0.9-beta.4
[v2.0.9-beta.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-beta.2...v2.0.9-beta.3
[v2.0.9-beta.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-beta.1...v2.0.9-beta.2
[v2.0.9-beta.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-beta.0...v2.0.9-beta.1
[v2.0.9-beta.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-alpha.10...v2.0.9-beta.0
[v2.0.9-alpha.10]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-alpha.9...v2.0.9-alpha.10
[v2.0.9-alpha.9]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-alpha.8...v2.0.9-alpha.9
[v2.0.9-alpha.8]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-alpha.7...v2.0.9-alpha.8
[v2.0.9-alpha.7]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-alpha.6...v2.0.9-alpha.7
[v2.0.9-alpha.6]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-alpha.4...v2.0.9-alpha.6
[v2.0.9-alpha.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-alpha.3...v2.0.9-alpha.4
[v2.0.9-alpha.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-alpha.2...v2.0.9-alpha.3
[v2.0.9-alpha.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-alpha.1...v2.0.9-alpha.2
[v2.0.9-alpha.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.9-alpha.0...v2.0.9-alpha.1
[v2.0.9-alpha.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-rc.3...v2.0.9-alpha.0
[v2.0.8-rc.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-rc.2...v2.0.8-rc.3
[v2.0.8-rc.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-rc.1...v2.0.8-rc.2
[v2.0.8-rc.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-rc.0...v2.0.8-rc.1
[v2.0.8-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-beta.3...v2.0.8-rc.0
[v2.0.8-beta.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-beta.2...v2.0.8-beta.3
[v2.0.8-beta.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-beta.1...v2.0.8-beta.2
[v2.0.8-beta.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-beta.0...v2.0.8-beta.1
[v2.0.8-beta.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-alpha.3...v2.0.8-beta.0
[v2.0.8-alpha.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-alpha.2...v2.0.8-alpha.3
[v2.0.8-alpha.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-alpha.1...v2.0.8-alpha.2
[v2.0.8-alpha.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.8-alpha.0...v2.0.8-alpha.1
[v2.0.8-alpha.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7...v2.0.8-alpha.0
[v2.0.7]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-rc.0...v2.0.7
[v2.0.7-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-beta.7...v2.0.7-rc.0
[v2.0.7-beta.7]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-beta.6...v2.0.7-beta.7
[v2.0.7-beta.6]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-beta.5...v2.0.7-beta.6
[v2.0.7-beta.5]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-beta.4...v2.0.7-beta.5
[v2.0.7-beta.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-beta.3...v2.0.7-beta.4
[v2.0.7-beta.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-beta.2...v2.0.7-beta.3
[v2.0.7-beta.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-beta.1...v2.0.7-beta.2
[v2.0.7-beta.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-beta.0...v2.0.7-beta.1
[v2.0.7-beta.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-alpha.5...v2.0.7-beta.0
[v2.0.7-alpha.5]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-alpha.4...v2.0.7-alpha.5
[v2.0.7-alpha.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-alpha.3...v2.0.7-alpha.4
[v2.0.7-alpha.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-alpha.2...v2.0.7-alpha.3
[v2.0.7-alpha.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-alpha.1...v2.0.7-alpha.2
[v2.0.7-alpha.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.7-alpha.0...v2.0.7-alpha.1
[v2.0.7-alpha.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.6...v2.0.7-alpha.0
[v2.0.6]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.6-beta.3...v2.0.6
[v2.0.6-beta.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.6-beta.2...v2.0.6-beta.3
[v2.0.6-beta.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.6-beta.1...v2.0.6-beta.2
[v2.0.6-beta.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.6-beta.0...v2.0.6-beta.1
[v2.0.6-beta.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.6-alpha.3...v2.0.6-beta.0
[v2.0.6-alpha.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.6-alpha.2...v2.0.6-alpha.3
[v2.0.6-alpha.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.6-alpha.1...v2.0.6-alpha.2
[v2.0.6-alpha.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.6-alpha.0...v2.0.6-alpha.1
[v2.0.6-alpha.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.5-rc.0...v2.0.6-alpha.0
[v2.0.5-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.5-beta.5...v2.0.5-rc.0
[v2.0.5-beta.5]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.5-beta.4...v2.0.5-beta.5
[v2.0.5-beta.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.5-beta.3...v2.0.5-beta.4
[v2.0.5-beta.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.5-beta.2...v2.0.5-beta.3
[v2.0.5-beta.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.5-beta.1...v2.0.5-beta.2
[v2.0.5-beta.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.5-beta.0...v2.0.5-beta.1
[v2.0.5-beta.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.5-alpha.3...v2.0.5-beta.0
[v2.0.5-alpha.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.5-alpha.2...v2.0.5-alpha.3
[v2.0.5-alpha.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.5-alpha.1...v2.0.5-alpha.2
[v2.0.5-alpha.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.5-alpha.0...v2.0.5-alpha.1
[v2.0.5-alpha.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.4...v2.0.5-alpha.0
[v2.0.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.4-rc.3...v2.0.4
[v2.0.4-rc.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.4-rc.2...v2.0.4-rc.3
[v2.0.4-rc.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.4-rc.1...v2.0.4-rc.2
[v2.0.4-rc.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.4-rc.0...v2.0.4-rc.1
[v2.0.4-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.4-beta.2...v2.0.4-rc.0
[v2.0.4-beta.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.4-beta.1...v2.0.4-beta.2
[v2.0.4-beta.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.4-beta.0...v2.0.4-beta.1
[v2.0.4-beta.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.4-alpha.1...v2.0.4-beta.0
[v2.0.4-alpha.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.4-alpha.0...v2.0.4-alpha.1
[v2.0.4-alpha.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.3-beta.8...v2.0.4-alpha.0
[v2.0.3-beta.8]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.3-beta.7...v2.0.3-beta.8
[v2.0.3-beta.7]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.3-beta.5...v2.0.3-beta.7
[v2.0.3-beta.5]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.3-beta.4...v2.0.3-beta.5
[v2.0.3-beta.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.3-beta.2...v2.0.3-beta.4
[v2.0.3-beta.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.3-beta.0...v2.0.3-beta.2
[v2.0.3-beta.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.3-alpha.0...v2.0.3-beta.0
[v2.0.3-alpha.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.27...v2.0.3-alpha.0
[v2.0.2-rc.27]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.26...v2.0.2-rc.27
[v2.0.2-rc.26]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.25...v2.0.2-rc.26
[v2.0.2-rc.25]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.24...v2.0.2-rc.25
[v2.0.2-rc.24]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.23...v2.0.2-rc.24
[v2.0.2-rc.23]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.22...v2.0.2-rc.23
[v2.0.2-rc.22]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.21...v2.0.2-rc.22
[v2.0.2-rc.21]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.20...v2.0.2-rc.21
[v2.0.2-rc.20]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.19...v2.0.2-rc.20
[v2.0.2-rc.19]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.18...v2.0.2-rc.19
[v2.0.2-rc.18]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.17...v2.0.2-rc.18
[v2.0.2-rc.17]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.14...v2.0.2-rc.17
[v2.0.2-rc.14]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.13...v2.0.2-rc.14
[v2.0.2-rc.13]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.12...v2.0.2-rc.13
[v2.0.2-rc.12]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.11...v2.0.2-rc.12
[v2.0.2-rc.11]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.10...v2.0.2-rc.11
[v2.0.2-rc.10]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.9...v2.0.2-rc.10
[v2.0.2-rc.9]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.8...v2.0.2-rc.9
[v2.0.2-rc.8]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.7...v2.0.2-rc.8
[v2.0.2-rc.7]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.6...v2.0.2-rc.7
[v2.0.2-rc.6]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.5...v2.0.2-rc.6
[v2.0.2-rc.5]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.4...v2.0.2-rc.5
[v2.0.2-rc.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.3...v2.0.2-rc.4
[v2.0.2-rc.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.2...v2.0.2-rc.3
[v2.0.2-rc.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.1...v2.0.2-rc.2
[v2.0.2-rc.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-rc.0...v2.0.2-rc.1
[v2.0.2-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-beta.6...v2.0.2-rc.0
[v2.0.2-beta.6]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-beta.3...v2.0.2-beta.6
[v2.0.2-beta.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-beta.1...v2.0.2-beta.3
[v2.0.2-beta.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-beta.0...v2.0.2-beta.1
[v2.0.2-beta.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-alpha.7...v2.0.2-beta.0
[v2.0.2-alpha.7]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-alpha.6...v2.0.2-alpha.7
[v2.0.2-alpha.6]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-alpha.5...v2.0.2-alpha.6
[v2.0.2-alpha.5]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-alpha.4...v2.0.2-alpha.5
[v2.0.2-alpha.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-alpha.3...v2.0.2-alpha.4
[v2.0.2-alpha.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-alpha.2...v2.0.2-alpha.3
[v2.0.2-alpha.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-alpha.1...v2.0.2-alpha.2
[v2.0.2-alpha.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.2-alpha.0...v2.0.2-alpha.1
[v2.0.2-alpha.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-rc.7...v2.0.2-alpha.0
[v2.0.1-rc.7]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-rc.6...v2.0.1-rc.7
[v2.0.1-rc.6]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-rc.5...v2.0.1-rc.6
[v2.0.1-rc.5]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-rc.4...v2.0.1-rc.5
[v2.0.1-rc.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-rc.2...v2.0.1-rc.4
[v2.0.1-rc.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-rc.1...v2.0.1-rc.2
[v2.0.1-rc.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-rc.0...v2.0.1-rc.1
[v2.0.1-rc.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-beta.6...v2.0.1-rc.0
[v2.0.1-beta.6]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-beta.5...v2.0.1-beta.6
[v2.0.1-beta.5]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-beta.4...v2.0.1-beta.5
[v2.0.1-beta.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-beta.3...v2.0.1-beta.4
[v2.0.1-beta.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-beta.2...v2.0.1-beta.3
[v2.0.1-beta.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-beta.1...v2.0.1-beta.2
[v2.0.1-beta.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-alpha.10...v2.0.1-beta.1
[v2.0.1-alpha.10]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-alpha.8...v2.0.1-alpha.10
[v2.0.1-alpha.8]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-alpha.7...v2.0.1-alpha.8
[v2.0.1-alpha.7]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-alpha.6...v2.0.1-alpha.7
[v2.0.1-alpha.6]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-alpha.5...v2.0.1-alpha.6
[v2.0.1-alpha.5]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-alpha.4...v2.0.1-alpha.5
[v2.0.1-alpha.4]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-alpha.3...v2.0.1-alpha.4
[v2.0.1-alpha.3]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-alpha.2...v2.0.1-alpha.3
[v2.0.1-alpha.2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-alpha.1...v2.0.1-alpha.2
[v2.0.1-alpha.1]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-alpha.0...v2.0.1-alpha.1
[v2.0.1-alpha.0]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.1-a-rc2...v2.0.1-alpha.0
[v2.0.1-a-rc2]: https://github.com/CooooolFrog/dragonfly/compare/v2.0.0...v2.0.1-a-rc2
[v2.0.0]: https://github.com/CooooolFrog/dragonfly/compare/v0.5.0...v2.0.0
[v0.5.0]: https://github.com/CooooolFrog/dragonfly/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/CooooolFrog/dragonfly/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/CooooolFrog/dragonfly/compare/v0.2.0...v0.3.0
