# WORKSPACE

# workspace file for running bazel in this repo

# --------------------------------------------------------------------------- #

# name the workspace
workspace(name = "wizter")

#################################################
### Initial Setup for Using Bazel and Gazelle ###
#################################################

# load rules for retriving dependencies with http
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# retrieve bazel build rules from the Internet
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "69de5c704a05ff37862f7e0f5534d4f479418afc21806c887db544a316f3cb6b",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.27.0/rules_go-v0.27.0.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.27.0/rules_go-v0.27.0.tar.gz",
    ],
)

# retrieve gazelle from the Internet
http_archive(
    name = "bazel_gazelle",
    sha256 = "62ca106be173579c0a167deb23358fdfe71ffa1e4cfdddf5582af26520f1c66f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
    ],
)



# load rules for working with Go
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
load("//:deps.bzl", "go_dependencies")

# download all go dependencies with gazelle and save them to a bild ./deps.blz

# the below gazelle: ... is a macro that tells gazelle to look for go dependencies in the ./deps.bzl file, then initialize them
# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

# overriding the version of go x tools due to a cycle dependency bug introduced upstream - https://github.com/bazelbuild/rules_go/issues/2479
# if `bazelisk run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies` is ran, then the below go_repository will
# get changed. This is here to override the version of X tools that gazelle uses. If allowed to use the latest version, a cycle dependency
# gets introduced, see above github issue. x tools, for now, needs to stay at:

    # sum = "h1:kRBLX7v7Af8W7Gdbbc908OJcdgtK8bOz9Uaj8/F1ACA=",
    # version = "v0.1.2",

# to update local repositories without overriding the below hard-coded
go_repository(
    name = "org_golang_x_tools",
    importpath = "golang.org/x/tools",
    sum = "h1:kRBLX7v7Af8W7Gdbbc908OJcdgtK8bOz9Uaj8/F1ACA=",
    version = "v0.1.2",
)

# init the dependency rules
go_rules_dependencies()

# set version of Go
go_register_toolchains(version = "1.16")

# init gazelle dependencies
gazelle_dependencies()

#########################################
### Rules for Using Protobuf and gRPC ###
#########################################

# get the rules for protobuff from the Internet
http_archive(
    name = "rules_proto",
    sha256 = "602e7161d9195e50246177e7c55b2f39950a9cf7366f74ed5f22fd45750cd208",
    strip_prefix = "rules_proto-97d8af4dc474595af3900dd85cb3a29ad28cc313",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
        "https://github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
    ],
)

# load the protobuf rules
load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

# init proto dependencies
rules_proto_dependencies()

# init the proto tool changes
rules_proto_toolchains()

#####################################
### Rules for Working with Docker ###
#####################################

# retrieve the Docker rules from the Internet
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "59d5b42ac315e7eadffa944e86e90c2990110a1c8075f1cd145f487e999d22b3",
    strip_prefix = "rules_docker-0.17.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.17.0/rules_docker-v0.17.0.tar.gz"],
)

# load the Docker rules
load("@io_bazel_rules_docker//repositories:repositories.bzl", container_repositories = "repositories",)

# init the repo to save Docker images
container_repositories()

# load rules for using Docker repos and dependencies
load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

# init the dependencies
container_deps()

# load the rules for pulling base images from Docker Hub
load("@io_bazel_rules_docker//container:container.bzl", "container_pull",)

# pull down a base image to use later
container_pull(
    name = "alpine_linux_amd64",
    registry = "index.docker.io",
    repository = "library/alpine",
    tag = "3.8",
)


########################################
### Rules for Build Go Docker Images ###
########################################

# load rules for making a Go-based .tar file to be loaded into Docker
load("@io_bazel_rules_docker//go:image.bzl", _go_image_repos = "repositories",)

# init the Docker repo so Bazel knows where to load Go images into
_go_image_repos()

#EOF