# 部署说明

## Docker Compose 部署

1. 将 `docker-compose.yml` 和 `config.yaml` 放在同一目录下
2. 运行服务：
```bash
docker compose up -d
```

默认每分钟执行一次（*/1 * * * *），如果需要修改执行频率，可以编辑 `docker-compose.yml` 中的 crontab 表达式：
```yaml
echo "*/1 * * * * ./cubox-archiver from-file -f /app/config.yaml" > /crontab
```

常用的 cron 表达式：
- 每分钟执行：`*/1 * * * *`
- 每小时执行：`0 * * * *`
- 每天凌晨执行：`0 0 * * *`
- 每30分钟执行：`*/30 * * * *`

停止服务：
```bash
docker compose down
```

查看日志：
```bash
docker compose logs -f
```

## Kubernetes 部署

### 设置命名空间
```bash
# 设置命名空间变量
export NAMESPACE=default  # 可以修改为你想要的命名空间

# 创建命名空间（如果不存在）
kubectl create namespace $NAMESPACE
```

### 创建 ConfigMap
从 `config.yaml` 创建 ConfigMap：
```bash
kubectl create configmap cubox-archiver-config --from-file=config.yaml -n $NAMESPACE
```

### 部署 CronJob
```bash
# 替换命名空间并部署
cat cronjob.yaml | sed "s/namespace: default/namespace: $NAMESPACE/" | kubectl apply -f -
```

### 检查部署状态
```bash
# 查看 ConfigMap
kubectl get configmap cubox-archiver-config -n $NAMESPACE

# 查看 CronJob
kubectl get cronjob cubox-archiver -n $NAMESPACE
``` 
