# Tiltfile
load("ext://helm_remote", "helm_remote")
load("ext://git_resource", "git_checkout")
load("ext://secret", "secret_from_dict")

# Load CRDs
k8s_yaml('local/promcrds.yaml')

# Load namespaces
k8s_yaml('local/namespace.yaml')

#Load secrets
k8s_yaml('local/secrets.yaml')

git_checkout(
    repository_url="git@github.com:saltpay/architecture-schema-registry-server#tilt/open-kafka",
    checkout_dir=".tilt_sources/kafka-stack/",
    unsafe_mode=True,
)

load_dynamic(".tilt_sources/kafka-stack/Tiltfile")

helm_remote(
    "postgresql",
    repo_name="bitnami",
    repo_url="https://charts.bitnami.com/bitnami",
    set=[
        "global.postgresql.auth.username=postgres",
        "global.postgresql.auth.password=postgres",
        "global.postgresql.auth.postgresPassword=postgres",
        "global.postgresql.auth.database=postgres"
    ]
)

# Docker build
docker_build("transaction-api-operational-mario", "src-mario-api/")
docker_build("transaction-api-operational-luigi", "src-luigi-transformer/")
# docker_build("transaction-api-operational-kamek", "src-kamek-cronjob/")
# docker_build("kamek-integration-tests", "src-kamek-cronjob", dockerfile="src-kamek-cronjob/integration.Dockerfile")

# k8s_yaml('./src-kamek-cronjob/integration.yaml' )

k8s_yaml(
    helm(
        "charts/trapi",
        name="transaction-api-operational",
        values=["charts/trapi/values.yaml"],
        set=[
            "luigi-transformer.image.repository=transaction-api-operational-luigi",
            "luigi-transformer.image.tag=latest",
            "luigi-transformer.serviceMonitor.enabled=false",
            "luigi-transformer.serviceAccount.create=true",
            "mario-api.enabled=true",
            "mario-api.image.repository=transaction-api-operational-mario",
            "mario-api.image.tag=latest",
            "mario-api.postgres.mockedData=true",
            "mario-api.service.port=8081",
            "mario-api.env.log_level=DEBUG"
        ],
    )
)

k8s_yaml(
    secret_from_dict(
        "transaction-api-operational-msk-eventstreaming",
        inputs={
            "endpoint": "kafka-0.kafka-headless.default.svc.cluster.local:9092",
            "username": "",
            "password": "",
        },
    )
)

k8s_yaml(
    secret_from_dict(
        "transaction-api-operational-db-transactions",
        inputs={
           "endpoint": "postgresql",
           "port": "5432",
           "username": "postgres",
           "password": "postgres",
           "database": "postgres",
           "url" : "postgresql://postgres:postgres@postgresql:5432/postgres",
           "read_url" : "postgresql://postgres:postgres@postgresql:5432/postgres"
        },
    )
)

k8s_resource(
    "postgresql",
    labels=["hell"],
    port_forwards=[
        port_forward(5432, 5432, "postgresql"),
    ],
)

k8s_resource(
    "transaction-api-operational-luigi",
    resource_deps=["kafka"],
    labels=["luigi"],
)

k8s_resource(
    "transaction-api-operational-mario",
    resource_deps=["kafka"],
    labels=["mario"],
    port_forwards=['8081']
)
