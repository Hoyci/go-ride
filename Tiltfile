# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

### K8s Config ###
k8s_yaml('./infra/development/k8s/app-config.yaml')
# k8s_yaml('./infra/development/k8s/secrets.yaml')
### End of K8s Config ###

### API Gateway ###
gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway/cmd'

local_resource(
  'api-gateway-compile',
  gateway_compile_cmd,
  deps=['./services/api-gateway', './shared'], labels="compiles")


docker_build_with_restart(
  'go-ride/api-gateway',
  '.',
  entrypoint=['/app/build/api-gateway'],
  dockerfile='./infra/development/docker/api-gateway.Dockerfile',
  only=[
    './build/api-gateway',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/api-gateway-deployment.yaml')
k8s_resource('api-gateway', port_forwards=8081, resource_deps=['api-gateway-compile'], labels="services")
### End of API Gateway ###

### User Service ###
user_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/user-service ./services/user-service/cmd'

local_resource(
  'user-service-compile',
  user_compile_cmd,
  deps=['./services/user-service', './shared'], labels="compiles")


docker_build_with_restart(
  'go-ride/user-service',
  '.',
  entrypoint=['/app/build/user-service'],
  dockerfile='./infra/development/docker/user-service.Dockerfile',
  only=[
    './build/user-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/user-service-deployment.yaml')
k8s_resource('user-service', resource_deps=['user-service-compile'], labels="services")
### End of user Service ###

### Trip Service ###
trip_compile_cmd  = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/trip-service ./services/trip-service/cmd'

local_resource(
  'trip-service-compile',
  trip_compile_cmd,
  deps=['./services/trip-service', './shared'], labels="compiles")


docker_build_with_restart(
  'go-ride/trip-service',
  '.',
  entrypoint=['/app/build/trip-service'],
  dockerfile='./infra/development/docker/trip-service.Dockerfile',
  only=[
    './build/trip-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/trip-service-deployment.yaml')
k8s_resource('trip-service', resource_deps=['trip-service-compile'], labels="services")
### End of Trip Service ###

### Web Service ###
web_compile_cmd = 'cd web && npm install --legacy-peer-deps && npm run build'

local_resource(
  'web-compile',
  web_compile_cmd,
  deps=['./web/src', './web/package.json', './web/vite.config.ts'], 
  labels="compiles"
)

docker_build_with_restart(
  'go-ride/web',
  '.',
  entrypoint=['nginx', '-g', 'daemon off;'],
  dockerfile='./infra/development/docker/web.Dockerfile',
  only=[
    './web/dist',
  ],
  live_update=[
    sync('./web/dist', '/usr/share/nginx/html'),
  ],
)

k8s_yaml('./infra/development/k8s/web-deployment.yaml')
k8s_resource('web-frontend', port_forwards='8080:80', resource_deps=['web-compile'], labels="services")
### End of Web Service ###