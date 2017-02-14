FROM flynn/busybox

COPY k8s-setup-network-environment /k8s-setup-network-environment

ENTRYPOINT ["/k8s-setup-network-environment"]
