# -*- mode: Python -*-
# This Tiltfile manages a Go project with live reload in Kubernetes

listen_port = 3000
metrics_port = 9090

# Use our wrapper to set the kube namespace
if os.getenv('TILT_WRAPPER') != '1':
    fail("could not run, ./hack/tilt/bin/tilt was not used to start tilt")

# Resource ids
holos_backend = 'Holos Server'
compile_id = 'Go Build'

# Default Registry.
# See: https://github.com/tilt-dev/tilt.build/blob/master/docs/choosing_clusters.md#manual-configuration
# Note, Tilt will append the image name to the registry uri path
# default_registry('{account}.dkr.ecr.{region}.amazonaws.com/holos-run/holos'.format(account=aws_account, region=aws_region))

# Set a name prefix specific to the user.  Multiple developers share the tilt-holos namespace.
developer = os.getenv('USER')
holos_server = 'holos'

# We always develop against the k3d-workload cluster
os.putenv('KUBECONFIG', os.path.abspath('./hack/tilt/kubeconfig'))

# Extensions are open-source, pre-packaged functions that extend Tilt
#
#   More info: https://github.com/tilt-dev/tilt-extensions
#   More info: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')
load('ext://k8s_attach', 'k8s_attach')
load('ext://git_resource', 'git_checkout')
load('ext://uibutton', 'cmd_button')

# Paths edited by the developer Tilt watches to trigger compilation.
# Generated files should be excluded to avoid an infinite build loop.
developer_paths = [
    './cmd',
    './internal/server',
    './internal/ent/schema',
    './frontend/package-lock.json',
    './frontend/src',
    './go.mod',
    './pkg',
    './service/holos',
]

# Builds the holos executable GOOS=linux
local_resource(compile_id, 'make linux', deps=developer_paths)

# Build Docker image
#   Tilt will automatically associate image builds with the resource(s)
#   that reference them (e.g. via Kubernetes or Docker Compose YAML).
#
#   More info: https://docs.tilt.dev/api.html#api.docker_build
#
docker_build_with_restart(
    'k3d-registry.holos.localhost:5100/holos',
    context='.',
    entrypoint=[
        '/app/bin/holos.linux',
        'server',
        '--log-format=text',
        '--oidc-issuer=https://login.holos.run',
        '--oidc-audience=275804490387516853@holos_quickstart', # auth proxy
        '--oidc-audience=270319630705329162@holos_platform', # holos cli
    ],
    dockerfile='./Dockerfile',
    only=['./bin'],
    # (Recommended) Updating a running container in-place
    # https://docs.tilt.dev/live_update_reference.html
    live_update=[
        # Sync files from host to container
        sync('./bin/', '/app/bin/'),
    ],
)

# Troubleshooting
def resource_name(id):
    print('resource: {}'.format(id))
    return id.name

workload_to_resource_function(resource_name)

# Customize a Kubernetes resource
#   By default, Kubernetes resource names are automatically assigned
#   based on objects in the YAML manifests, e.g. Deployment name.
#
#   Tilt strives for sane defaults, so calling k8s_resource is
#   optional, and you only need to pass the arguments you want to
#   override.
#
#   More info: https://docs.tilt.dev/api.html#api.k8s_resource
#
k8s_yaml(blob(str(read_file('./hack/tilt/k8s/dev-holos-app/deployment.yaml'))))

# Backend server process
k8s_resource(
    workload=holos_server,
    new_name=holos_backend,
    objects=[],
    resource_deps=[compile_id],
    links=[
        link('https://app.holos.localhost/ui/'.format(developer), "Holos Web UI")
    ],
)

# Database
print("✨ Tiltfile evaluated")
