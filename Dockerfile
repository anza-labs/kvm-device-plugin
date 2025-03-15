# Easy crosscomple toolkit
FROM ghcr.io/grpc-ecosystem/grpc-health-probe:v0.4.37 AS probe

# Build the plugin binary
FROM --platform=$BUILDPLATFORM docker.io/library/rust:1.85.0 AS builder

### TODO

# Use distroless as minimal base image to package the plugin binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# hadolint ignore=DL3007
FROM gcr.io/distroless/static:latest
WORKDIR /
COPY --from=builder /workspace/kvm-device-plugin .
COPY --from=probe /ko-app/grpc-health-probe /grpc_health_probe

ENTRYPOINT ["/kvm-device-plugin"]
