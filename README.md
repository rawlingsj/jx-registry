# jx registry

[![Documentation](https://godoc.org/github.com/jenkins-x-plugins/jx-registry?status.svg)](https://pkg.go.dev/mod/github.com/jenkins-x-plugins/jx-registry)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenkins-x-plugins/jx-registry)](https://goreportcard.com/report/github.com/jenkins-x-plugins/jx-registry)
[![Releases](https://img.shields.io/github/release-pre/jenkins-x/helmboot.svg)](https://github.com/jenkins-x-plugins/jx-registry/releases)
[![LICENSE](https://img.shields.io/github/license/jenkins-x/helmboot.svg)](https://github.com/jenkins-x-plugins/jx-registry/blob/master/LICENSE)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://slack.k8s.io/)

`jx-registry` is a simple command line tool for working with container registries.

The main use case is initially to support lazy creation of AWS ECR registries on demand. Most other registries allow a registry to be created and used for different images.


foo 3

## Getting Started

Download the [jx-registry binary](https://github.com/jenkins-x-plugins/jx-registry/releases) for your operating system and add it to your `$PATH`.

## Enabling Cache images

If you wish to also create a cache image in addition to the ECR image for your repository enable the `CACHE_SUFFIX` environment variable.

e.g. in your local `.lighthouse/jenkins-x/release.yaml` file you could do something like:

```yaml
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  creationTimestamp: null
  name: release
spec:
  pipelineSpec:
    tasks:
    - name: from-build-pack
      resources: {}
      taskSpec:
        metadata: {}
        stepTemplate:
          image: uses:jenkins-x/jx3-pipeline-catalog/tasks/javascript/release.yaml@versionStream
        steps:
        - image: uses:jenkins-x/jx3-pipeline-catalog/tasks/git-clone/git-clone.yaml@versionStream
          name: ""
          resources: {}
        - name: next-version
          resources: {}
        - name: jx-variables
          resources: {}
        - name: build-npm-install
          resources: {}
        - name: build-npm-test
          resources: {}
        - name: check-registry
          env:
          - name: CACHE_SUFFIX
            value: "/cache"
          resources: {}
        - name: build-container-build
          resources: {}
        - name: promote-changelog
          resources: {}
        - name: promote-helm-release
          resources: {}
        - name: promote-jx-promote
          resources: {}
```

## Commands

See the [jx-registry command reference](https://github.com/jenkins-x-plugins/jx-registry/blob/master/docs/cmd/jx-registry.md#jx-registry)

