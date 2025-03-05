# kvm-device-plugin

[![GitHub License](https://img.shields.io/github/license/anza-labs/kvm-device-plugin)][license]
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md)
[![GitHub issues](https://img.shields.io/github/issues/anza-labs/kvm-device-plugin)](https://github.com/anza-labs/kvm-device-plugin/issues)
[![GitHub release](https://img.shields.io/github/release/anza-labs/kvm-device-plugin)](https://GitHub.com/anza-labs/kvm-device-plugin/releases/)
[![Go Report Card](https://goreportcard.com/badge/github.com/anza-labs/kvm-device-plugin)](https://goreportcard.com/report/github.com/anza-labs/kvm-device-plugin)

`kvm-device-plugin` is a Kubernetes Device Plugin that manages access to `/dev/kvm` (Kernel-based Virtual Machine) devices. It allows workloads running in Kubernetes to request KVM access via the Device Plugin interface, ensuring proper communication with the kubelet.

- [kvm-device-plugin](#kvm-device-plugin)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)
  - [How It Works](#how-it-works)
  - [Compatibility](#compatibility)
  - [License](#license)
  - [Attributions](#attributions)

## Features

- Provides access to `/dev/kvm` for containers running in Kubernetes.
- Implements the Kubernetes Device Plugin API to manage KVM allocation.
- Ensures that only workloads explicitly requesting KVM access receive it.

## Installation

To deploy the `kvm-device-plugin`, apply the provided manifests:

```sh
LATEST="$(curl -s 'https://api.github.com/repos/anza-labs/kink/releases/latest' | jq -r '.tag_name')"
kubectl apply -k "https://github.com/anza-labs/kink/?ref=${LATEST}"
```

## Usage

To request access to `/dev/kvm` in a pod, specify the device resource in the `resources` section:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: kvm-checker
spec:
  restartPolicy: Never
  containers:
    - name: kvm-checker
      image: busybox
      command: ["sh", "-c", "[ -e /dev/kvm ]"]
      resources:
        requests:
          device.anza-labs.com/kvm: '1' # Request KVM device
        limits:
          device.anza-labs.com/kvm: '1' # Limit KVM device
```

## How It Works

1. The `kvm-device-plugin` registers with the kubelet and advertises available KVM devices.
2. When a pod requests the `device.anza-labs.com/kvm` resource, the device plugin assigns a `/dev/kvm` device to the container.
3. The container is granted access to `/dev/kvm` for virtualization tasks.

## Compatibility

- Kubernetes 1.20+
- Nodes must have KVM enabled (check with `lsmod | grep kvm`)

## License

`kvm-device-plugin` is licensed under the [Apache-2.0][license].

## Attributions

This codebase is inspired by:
- [github.com/squat/generic-device-plugin](https://github.com/squat/generic-device-plugin)
- [github.com/kubevirt/kubernetes-device-plugins](https://github.com/kubevirt/kubernetes-device-plugins)

<!-- Resources -->

[license]: https://github.com/anza-labs/kvm-device-plugin/blob/main/LICENSE
