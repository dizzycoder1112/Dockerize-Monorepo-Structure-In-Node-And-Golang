# Kubernetes Helm Chart

Helm chart for deploying monorepo applications to Kubernetes with Infisical secrets management.

## Directory Structure

```
kubernetes/
├── Chart.yaml                  # Helm chart metadata
├── values.yaml                 # Default values template
├── templates/                  # Kubernetes manifest templates
│   ├── deployment.yaml        # Deployment configuration
│   ├── service.yaml           # Service configuration
│   ├── ingress.yaml           # Ingress rules
│   ├── hpa.yaml               # Horizontal Pod Autoscaler
│   └── secret.yaml            # Secret configuration
├── secrets/                    # Environment-specific secrets
│   ├── secrets.{app}.dev.yaml
│   ├── secrets.{app}.stage.yaml
│   └── secrets.{app}.prod.yaml
└── infisical/                  # Infisical integration
    ├── values.yaml
    ├── infisical-auth-secret.yaml
    └── infisical-secret-sync.yaml
```

## Prerequisites

- Kubernetes cluster (v1.19+)
- Helm 3.x installed
- kubectl configured
- Infisical account (for secrets management)

## Quick Start

### 1. Prepare Values File

Copy and customize the values file for your application:

```bash
cp values.yaml values.your-app.yaml
```

Edit `values.your-app.yaml`:

```yaml
app: your-app-name
replicaCount: 1

containers:
  port: 3000
  image:
    repository: registry.digitalocean.com/your-org/your-project
    pullPolicy: Always
    tag: '.dev'  # .dev / .stage / .prod
  env:
    - name: NODE_ENV
      value: 'development'
    - name: PORT
      value: '3000'
  secretsFrom: infisical-your-app-secrets

infisical:
  envSlug: dev
  projectId: 'your-infisical-project-id'

ingress:
  enabled: true
  host: your-app.your-domain.io
```

### 2. Deploy to Kubernetes

```bash
# Deploy or upgrade
helm upgrade --install <release-name> ./kubernetes \
  -f ./kubernetes/values.your-app.yaml \
  -n <your-namespace> \
  --create-namespace

# Example
helm upgrade --install my-app ./kubernetes \
  -f ./kubernetes/values.your-app.yaml \
  -n production
```

### 3. Verify Deployment

```bash
# Check deployment status
kubectl get pods -n <your-namespace>

# Check service
kubectl get svc -n <your-namespace>

# Check ingress
kubectl get ingress -n <your-namespace>

# View logs
kubectl logs -f deployment/<release-name> -n <your-namespace>
```

## Configuration

### Core Values

| Parameter | Description | Default |
|-----------|-------------|---------|
| `app` | Application name | `app-name` |
| `replicaCount` | Number of pod replicas | `1` |
| `containers.port` | Container port | `3000` |
| `containers.image.repository` | Docker image repository | - |
| `containers.image.tag` | Image tag suffix (.dev/.stage/.prod) | `.dev` |
| `containers.image.pullPolicy` | Image pull policy | `Always` |
| `containers.env` | Environment variables | `[]` |
| `containers.secretsFrom` | Infisical secret name | - |

### Resources

```yaml
resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 500m
    memory: 512Mi
```

### Service

```yaml
service:
  enabled: true
  protocol: TCP
  type: ClusterIP
```

### Ingress

```yaml
ingress:
  enabled: true
  host: your.domain.io
```

### Horizontal Pod Autoscaler (HPA)

```yaml
hpa:
  enabled: false
  minReplicas: 1
  maxReplicas: 5
  targetCPUUtilizationPercentage: 90
  targetMemoryUtilizationPercentage: 100
```

### Node Selector

```yaml
nodeSelector:
  doks.digitalocean.com/node-pool: "pool-258gs845a"
```

## Infisical Integration

This chart supports Infisical for secrets management.

### Setup Infisical

1. Create Infisical project and get project ID
2. Create auth secret:

```bash
kubectl apply -f ./kubernetes/infisical/infisical-auth-secret.yaml -n <namespace>
```

3. Configure values:

```yaml
infisical:
  envSlug: dev  # dev / staging / prod
  projectId: 'your-project-id'
```

4. Deploy secret sync:

```bash
kubectl apply -f ./kubernetes/infisical/infisical-secret-sync.yaml -n <namespace>
```

## Common Commands

### Deploy Different Environments

```bash
# Development
helm upgrade --install my-app-dev ./kubernetes \
  -f ./kubernetes/values.my-app.yaml \
  --set containers.image.tag=.dev \
  --set infisical.envSlug=dev \
  -n development

# Staging
helm upgrade --install my-app-stage ./kubernetes \
  -f ./kubernetes/values.my-app.yaml \
  --set containers.image.tag=.stage \
  --set infisical.envSlug=staging \
  -n staging

# Production
helm upgrade --install my-app-prod ./kubernetes \
  -f ./kubernetes/values.my-app.yaml \
  --set containers.image.tag=.prod \
  --set infisical.envSlug=prod \
  -n production
```

### Update Image Version

```bash
# Rollout new version
helm upgrade my-app ./kubernetes \
  -f ./kubernetes/values.my-app.yaml \
  --set containers.image.tag=.prod \
  -n production

# Force pod restart
kubectl rollout restart deployment/<release-name> -n <namespace>
```

### Rollback

```bash
# List releases
helm list -n <namespace>

# Rollback to previous version
helm rollback <release-name> -n <namespace>

# Rollback to specific revision
helm rollback <release-name> <revision> -n <namespace>
```

### Uninstall

```bash
# Uninstall release
helm uninstall <release-name> -n <namespace>

# Delete namespace
kubectl delete namespace <namespace>
```

## Troubleshooting

### Check Helm Release

```bash
# Get release status
helm status <release-name> -n <namespace>

# Get release values
helm get values <release-name> -n <namespace>

# Get release manifest
helm get manifest <release-name> -n <namespace>
```

### Debug Pods

```bash
# Describe pod
kubectl describe pod <pod-name> -n <namespace>

# Get pod logs
kubectl logs <pod-name> -n <namespace>

# Get previous container logs (if crashed)
kubectl logs <pod-name> -n <namespace> --previous

# Execute into pod
kubectl exec -it <pod-name> -n <namespace> -- /bin/sh
```

### Check Infisical Secrets

```bash
# Check secret exists
kubectl get secret <secret-name> -n <namespace>

# View secret (base64 encoded)
kubectl get secret <secret-name> -n <namespace> -o yaml

# Decode secret
kubectl get secret <secret-name> -n <namespace> -o jsonpath='{.data}' | jq
```

## Applications in This Monorepo

Current applications with Kubernetes deployments:

- **management-server** - Management server API
- **notify-worker** - Notification worker service
- **promotion-worker** - Promotion processing worker
- **share-profit-handler** - Profit sharing handler
- **weedza-app** - Main application

Each app has environment-specific secrets in `kubernetes/secrets/`.

## Best Practices

1. **Always use specific tags** - Avoid `:latest` in production
2. **Set resource limits** - Prevent resource exhaustion
3. **Use HPA for high-traffic apps** - Enable auto-scaling
4. **Review security** - Use Infisical for sensitive data
5. **Test in staging first** - Always validate before production
6. **Monitor deployments** - Use `kubectl get events` to watch

## References

- [Helm Documentation](https://helm.sh/docs/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Infisical Documentation](https://infisical.com/docs/)
