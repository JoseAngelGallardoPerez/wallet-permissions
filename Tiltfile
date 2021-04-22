print("Wallet Permissions")

load("ext://restart_process", "docker_build_with_restart")

cfg = read_yaml(
    "tilt.yaml",
    default = read_yaml("tilt.yaml.sample"),
)

local_resource(
    "permissions-build-binary",
    "make fast_build",
    deps = ["./cmd", "./internal"],
)
local_resource(
    "permissions-generate-protpbuf",
    "make gen-protobuf",
    deps = ["./rpc/permissions/permissions.proto"],
)

docker_build(
    "velmie/wallet-permissions-db-migration",
    ".",
    dockerfile = "Dockerfile.migrations",
    only = "migrations",
)
k8s_resource(
    "wallet-permissions-db-migration",
    trigger_mode = TRIGGER_MODE_MANUAL,
    resource_deps = ["wallet-permissions-db-init"],
)

wallet_permissions_options = dict(
    entrypoint = "/app/service_permissions",
    dockerfile = "Dockerfile.prebuild",
    port_forwards = [],
    helm_set = [],
)

if cfg["debug"]:
    wallet_permissions_options["entrypoint"] = "$GOPATH/bin/dlv --continue --listen :%s --accept-multiclient --api-version=2 --headless=true exec /app/service_permissions" % cfg["debug_port"]
    wallet_permissions_options["dockerfile"] = "Dockerfile.debug"
    wallet_permissions_options["port_forwards"] = cfg["debug_port"]
    wallet_permissions_options["helm_set"] = ["containerLivenessProbe.enabled=false", "containerPorts[0].containerPort=%s" % cfg["debug_port"]]

docker_build_with_restart(
    "velmie/wallet-permissions",
    ".",
    dockerfile = wallet_permissions_options["dockerfile"],
    entrypoint = wallet_permissions_options["entrypoint"],
    only = [
        "./build",
        "zoneinfo.zip",
    ],
    live_update = [
        sync("./build", "/app/"),
    ],
)
k8s_resource(
    "wallet-permissions",
    resource_deps = ["wallet-permissions-db-migration"],
    port_forwards = wallet_permissions_options["port_forwards"],
)

yaml = helm(
    "./helm/wallet-permissions",
    # The release name, equivalent to helm --name
    name = "wallet-permissions",
    # The values file to substitute into the chart.
    values = ["./helm/values-dev.yaml"],
    set = wallet_permissions_options["helm_set"],
)

k8s_yaml(yaml)
