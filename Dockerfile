FROM flynn/busybox

COPY setup-network-environment /setup-network-environment

ENTRYPOINT ["/setup-network-environment"]
