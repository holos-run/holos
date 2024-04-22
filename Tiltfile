# -*- mode: Python -*-
# This Tiltfile manages a Go project with live leload in Kubernetes

listen_port = 3000
metrics_port = 9090

# Use our wrapper to set the kube namespace
if os.getenv('TILT_WRAPPER') != '1':
    fail("could not run, ./hack/tilt/bin/tilt was not used to start tilt")

# AWS Account to work in
aws_account = '271053619184'
aws_region = 'us-east-2'

# Resource ids
holos_backend = 'Holos Backend'
pg_admin = 'pgAdmin'
pg_cluster = 'PostgresCluster'
pg_svc = 'Database Pod'
compile_id = 'Go Build'
auth_id = 'Auth Policy'
lint_id = 'Run Linters'
tests_id = 'Run Tests'

# PostgresCluster resource name in k8s
pg_cluster_name = 'holos'
# Database name inside the PostgresCluster
pg_database_name = 'holos'
# PGAdmin name
pg_admin_name = 'pgadmin'

# Default Registry.
# See: https://github.com/tilt-dev/tilt.build/blob/master/docs/choosing_clusters.md#manual-configuration
# Note, Tilt will append the image name to the registry uri path
default_registry('{account}.dkr.ecr.{region}.amazonaws.com/holos-run/holos-server'.format(account=aws_account, region=aws_region))

# Set a name prefix specific to the user.  Multiple developers share the tilt-holos namespace.
developer = os.getenv('USER')
holos_server = 'holos'
# See ./hack/tilt/bin/tilt
namespace = os.getenv('NAMESPACE')
# We always develop against the k1 cluster.
os.putenv('KUBECONFIG', os.path.abspath('./hack/tilt/kubeconfig'))
# The context defined in ./hack/tilt/kubeconfig
allow_k8s_contexts('sso@k1')
allow_k8s_contexts('sso@k2')
allow_k8s_contexts('sso@k3')
allow_k8s_contexts('sso@k4')
allow_k8s_contexts('sso@k5')
# PG db connection for localhost -> k8s port-forward
os.putenv('PGHOST', 'localhost')
os.putenv('PGPORT', '15432')
# We always develop in the dev aws account.
os.putenv('AWS_CONFIG_FILE', os.path.abspath('./hack/tilt/aws.config'))
os.putenv('AWS_ACCOUNT', aws_account)
os.putenv('AWS_DEFAULT_REGION', aws_region)
os.putenv('AWS_PROFILE', 'dev-holos')
os.putenv('AWS_SDK_LOAD_CONFIG', '1')
# Authenticate to AWS ECR when tilt up is run by the developer
local_resource('AWS Credentials', './hack/tilt/aws-login.sh', auto_init=True)

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

# Builds the holos-server executable
local_resource(compile_id, 'make build', deps=developer_paths)

# Build Docker image
#   Tilt will automatically associate image builds with the resource(s)
#   that reference them (e.g. via Kubernetes or Docker Compose YAML).
#
#   More info: https://docs.tilt.dev/api.html#api.docker_build
#
docker_build_with_restart(
    'holos',
    context='.',
    entrypoint=[
        '/app/bin/holos',
        'server',
        '--listen-port={}'.format(listen_port),
        '--oidc-issuer=https://login.ois.run',
        '--oidc-audience=262096764402729854@holos_platform',
        '--log-level=debug',
        '--metrics-port={}'.format(metrics_port),
    ],
    dockerfile='./hack/tilt/Dockerfile',
    only=['./bin'],
    # (Recommended) Updating a running container in-place
    # https://docs.tilt.dev/live_update_reference.html
    live_update=[
        # Sync files from host to container
        sync('./bin', '/app/bin'),
        # Wait for aws-login https://github.com/tilt-dev/tilt/issues/3048
        sync('./tilt/aws-login.last', '/dev/null'),
        # Execute commands in the container when paths change
        # run('/app/hack/codegen.sh', trigger=['./app/api'])
    ],
)


# Run local commands
#   Local commands can be helpful for one-time tasks like installing
#   project prerequisites. They can also manage long-lived processes
#   for non-containerized services or dependencies.
#
#   More info: https://docs.tilt.dev/local_resource.html
#
# local_resource('install-helm',
#                cmd='which helm > /dev/null || brew install helm',
#                # `cmd_bat`, when present, is used instead of `cmd` on Windows.
#                cmd_bat=[
#                    'powershell.exe',
#                    '-Noninteractive',
#                    '-Command',
#                    '& {if (!(Get-Command helm -ErrorAction SilentlyContinue)) {scoop install helm}}'
#                ]
# )

# Teach tilt about our custom resources (Note, this may be intended for workloads)
# k8s_kind('authorizationpolicy')
# k8s_kind('requestauthentication')
# k8s_kind('virtualservice')
k8s_kind('pgadmin')


# Troubleshooting
def resource_name(id):
    print('resource: {}'.format(id))
    return id.name


workload_to_resource_function(resource_name)

# Apply Kubernetes manifests
#   Tilt will build & push any necessary images, re-deploying your
#   resources as they change.
#
#   More info: https://docs.tilt.dev/api.html#api.k8s_yaml
#

def holos_yaml():
    """Return a k8s Deployment personalized for the developer."""
    k8s_yaml_template = str(read_file('./hack/tilt/k8s.yaml'))
    return k8s_yaml_template.format(
        name=holos_server,
        developer=developer,
        namespace=namespace,
        listen_port=listen_port,
        metrics_port=metrics_port,
        tz=os.getenv('TZ'),
    )

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
k8s_yaml(blob(holos_yaml()))

# Backend server process
k8s_resource(
    workload=holos_server,
    new_name=holos_backend,
    objects=[
        '{}:serviceaccount'.format(holos_server),
        '{}:servicemonitor'.format(holos_server),
    ],
    resource_deps=[compile_id],
    links=[
        link('https://{}.app.dev.k2.holos.run/ui/'.format(developer), "Holos Web UI")
    ],
)


# AuthorizationPolicy - Beyond Corp functionality
k8s_resource(
    new_name=auth_id,
    objects=[
      '{}:virtualservice'.format(holos_server),
    ],
)

# Database
# Note: Tilt confuses the backup pods with the database server pods, so this code is careful to tease the pods
# apart so logs are streamed correctly.
# See: https://github.com/tilt-dev/tilt.specs/blob/master/resource_assembly.md

# pgAdmin Web UI
k8s_resource(
    workload=pg_admin_name,
    new_name=pg_admin,
    port_forwards=[
        port_forward(15050, 5050, pg_admin),
    ],
)

# Disabled because these don't group resources nicely
# k8s_kind('postgrescluster')

# Postgres database in-cluster
k8s_resource(
    new_name=pg_cluster,
    objects=['holos:postgrescluster'],
)

# Needed to select the database by label
# https://docs.tilt.dev/api.html#api.k8s_custom_deploy
k8s_custom_deploy(
    pg_svc,
    apply_cmd=['./hack/tilt/k8s-get-db-sts', pg_cluster_name],
    delete_cmd=['echo', 'Skipping delete.  Object managed by custom resource.'],
    deps=[],
)
k8s_resource(
    pg_svc,
    port_forwards=[
        port_forward(15432, 5432, 'psql'),
    ],
    resource_deps=[pg_cluster]
)


# Run tests
local_resource(
    tests_id,
    'make test',
    allow_parallel=True,
    auto_init=False,
    deps=developer_paths,
)

# Run linter
local_resource(
    lint_id,
    'make lint',
    allow_parallel=True,
    auto_init=False,
    deps=developer_paths,
)

# UI Buttons for helpful things.
# Icons: https://fonts.google.com/icons
os.putenv("GH_FORCE_TTY", "80%")
cmd_button(
    '{}:go-test-failfast'.format(tests_id),
    argv=['./hack/tilt/go-test-failfast'],
    resource=tests_id,
    icon_name='quiz',
    text='Fail Fast',
)
cmd_button(
    '{}:issues'.format(holos_server),
    argv=['./hack/tilt/gh-issues'],
    resource=holos_backend,
    icon_name='folder_data',
    text='Issues',
)
cmd_button(
    '{}:gh-issue-view'.format(holos_server),
    argv=['./hack/tilt/gh-issue-view'],
    resource=holos_backend,
    icon_name='task',
    text='View Issue',
)
cmd_button(
    '{}:get-pgdb-creds'.format(holos_server),
    argv=['./hack/tilt/get-pgdb-creds', pg_cluster_name, pg_database_name],
    resource=pg_svc,
    icon_name='lock_open_right',
    text='DB Creds',
)
cmd_button(
    '{}:get-pgdb-creds'.format(pg_admin_name),
    argv=['./hack/tilt/get-pgdb-creds', pg_cluster_name, pg_database_name],
    resource=pg_admin,
    icon_name='lock_open_right',
    text='DB Creds',
)
cmd_button(
    '{}:get-pgadmin-creds'.format(pg_admin_name),
    argv=['./hack/tilt/get-pgadmin-creds', pg_admin_name],
    resource=pg_admin,
    icon_name='lock_open_right',
    text='pgAdmin Login',
)

print("✨ Tiltfile evaluated")
