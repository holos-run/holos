awk 'BEGIN {RS="---"} NR==1 {print > "service.yaml"} NR==2 {print > "deployment.yaml"}' $1
mv service.yaml $2/rendered-service.yaml
mv deployment.yaml $2/rendered-deployment.yaml
