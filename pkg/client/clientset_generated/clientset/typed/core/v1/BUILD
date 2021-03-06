package(default_visibility = ["//visibility:public"])

licenses(["notice"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "componentstatus.go",
        "configmap.go",
        "core_client.go",
        "doc.go",
        "endpoints.go",
        "event.go",
        "event_expansion.go",
        "generated_expansion.go",
        "limitrange.go",
        "namespace.go",
        "namespace_expansion.go",
        "node.go",
        "node_expansion.go",
        "persistentvolume.go",
        "persistentvolumeclaim.go",
        "pod.go",
        "pod_expansion.go",
        "podtemplate.go",
        "replicationcontroller.go",
        "resourcequota.go",
        "secret.go",
        "service.go",
        "service_expansion.go",
        "serviceaccount.go",
    ],
    tags = ["automanaged"],
    deps = [
        "//pkg/api:go_default_library",
        "//pkg/api/v1:go_default_library",
        "//pkg/apimachinery/registered:go_default_library",
        "//pkg/apis/meta/v1:go_default_library",
        "//pkg/apis/policy/v1beta1:go_default_library",
        "//pkg/client/restclient:go_default_library",
        "//pkg/fields:go_default_library",
        "//pkg/runtime:go_default_library",
        "//pkg/runtime/schema:go_default_library",
        "//pkg/runtime/serializer:go_default_library",
        "//pkg/util/net:go_default_library",
        "//pkg/watch:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//pkg/client/clientset_generated/clientset/typed/core/v1/fake:all-srcs",
    ],
    tags = ["automanaged"],
)
