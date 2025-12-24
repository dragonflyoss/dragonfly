# Dragonfly

![alt][logo-linear]

[![GitHub release](https://img.shields.io/github/release/dragonflyoss/dragonfly.svg)](https://github.com/dragonflyoss/dragonfly/releases)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/dragonfly)](https://artifacthub.io/packages/search?repo=dragonfly)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/dragonflyoss/dragonfly/badge)](https://scorecard.dev/viewer/?uri=github.com/dragonflyoss/dragonfly)
[![CI](https://github.com/dragonflyoss/dragonfly/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/dragonflyoss/dragonfly/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dragonflyoss/dragonfly?style=flat-square)](https://goreportcard.com/report/github.com/dragonflyoss/dragonfly)
[![Discussions](https://img.shields.io/badge/discussions-on%20github-blue?style=flat-square)](https://github.com/dragonflyoss/dragonfly/discussions)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/10432/badge)](https://www.bestpractices.dev/projects/10432)
[![Twitter](https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Ftwitter.com%2Fdragonfly_oss)](https://twitter.com/dragonfly_oss)
[![LICENSE](https://img.shields.io/github/license/dragonflyoss/dragonfly.svg?style=flat-square)](https://github.com/dragonflyoss/dragonfly/blob/main/LICENSE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fdragonflyoss%2Fdragonfly.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fdragonflyoss%2Fdragonfly?ref=badge_shield&issueType=license)
[![LFX Health Score](https://img.shields.io/static/v1?label=Health%20Score&message=Healthy&color=A7F3D0&logo=linuxfoundation&logoColor=white&style=flat)](https://insights.linuxfoundation.org/project/d7y)
[![CNCF Status](https://img.shields.io/badge/CNCF%20Status-Graduated-informational)](https://www.cncf.io/projects/dragonfly/)

## Introduction

Delivers efficient, stable, and secure data distribution and acceleration powered by P2P technology,
with an optional content‑addressable filesystem that accelerates OCI container launch.
It aims to provide a best‑practice, standards‑based solution for cloud‑native architectures,
improving large‑scale delivery of files, container images, OCI artifacts, AI/ML models, caches,
logs, dependencies, etc.

## Components

### Dragonfly Injector

A Kubernetes Mutating Admission Webhook that automatically injects Dragonfly P2P capabilities into pods.
The injector simplifies pod configuration by automating the injection of:

- **P2P Proxy Environment Variables**: Routes application traffic through Dragonfly's P2P network
- **Dfdaemon Socket Mounts**: Enables CLI tools to communicate with dfdaemon
- **CLI Tools**: Injects dfget and other tools via init container for base images that don't include them

The injector supports both namespace-level (via labels) and pod-level (via annotations) injection policies.

**Quick Start:**
```bash
# Deploy the injector
kubectl apply -k deploy/injector/

# Enable injection for a namespace
kubectl label namespace my-namespace dragonfly.io/inject=true
```

For detailed documentation, see:
- [Injector README](cmd/injector/README.md)
- [Deployment Guide](deploy/injector/DEPLOYMENT.md)

## Documentation

You can find the full documentation on the [d7y.io][d7y.io].

## Community

Join the conversation and help the community grow. Here are the ways to get involved:

- **Slack Channel**: [#dragonfly](https://cloud-native.slack.com/messages/dragonfly/) on [CNCF Slack](https://slack.cncf.io/)
- **Github Discussions**: [Dragonfly Discussion Forum][discussion]
- **Developer Group**: <dragonfly-developers@googlegroups.com>
- **Mailing Lists**:
  - **Developers**: <dragonfly-developers@googlegroups.com>
  - **Maintainers**: <dragonfly-maintainers@googlegroups.com>
- **Twitter**: [@dragonfly_oss](https://twitter.com/dragonfly_oss)
- **DingTalk Group**: `22880028764`

You can also find community information in the [community repository](https://github.com/dragonflyoss/community).

## Roadmap

You can find the [roadmap](https://github.com/dragonflyoss/community/blob/master/ROADMAP.md)
in the [community repository](https://github.com/dragonflyoss/community).

## Security

### Security Audit

A third-party security audit of Dragonfly was performed by Trail of Bits, with the full report available at
[Dragonfly Comprehensive Report](docs/security/dragonfly-comprehensive-report-2023.pdf).

### Reporting security vulnerabilities

If you discover a vulnerability, please report it per our Security Policy at [Security Policy](https://github.com/dragonflyoss/community/blob/master/SECURITY.md),
and security insights are detailed in [SECURITY-INSIGHTS.yml](SECURITY-INSIGHTS.yml).

## Software bill of materials

We publish SBOMs with all of our releases. You can find them in Github release assets.

[arch]: docs/images/arch.png
[logo-linear]: docs/images/logo/dragonfly-linear.svg
[discussion]: https://github.com/dragonflyoss/dragonfly/discussions
[contributing]: https://github.com/dragonflyoss/community/blob/master/CONTRIBUTING.md
[codeconduct]: https://github.com/dragonflyoss/community/blob/master/CODE_OF_CONDUCT.md
[d7y.io]: https://d7y.io/

## Contributing

You should check out our
[CONTRIBUTING][contributing] and develop the project together.

## Code of Conduct

Please refer to our [Code of Conduct][codeconduct] which applies to all Dragonfly community members.
