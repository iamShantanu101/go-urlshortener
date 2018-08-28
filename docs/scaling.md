## Scaling Guide:

This document outlines the available options for scaling the URL shortener application:

### Considerations:
1. For using same Persistent Volumes across all Pods, we need to use NFS like (GlusterFS, Ceph, etc) which supports `ReadWriteMany` option for Persistent Volumes. The suppport chart is documented [here](https://kubernetes.io/docs/concepts/storage/persistent-volumes/).
2. If, NFS seems to be overkill (since it adds maintenance overhead), we can use [emptyDir](https://kubernetes.io/docs/concepts/storage/volumes/#emptydir) instead. However, note that the lifecycle of an emptyDir is bound to its pods. The storage will be persistent across container crashes but will be removed if the pod get terminated/deleted. Example:
   ```yml   
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: url-shortener
     labels:
       app: url-shortener
   spec:
     selector:
       matchLabels:
         app: url-shortener
     template:
       metadata:
         labels:
           app: url-shortener
       spec:
         containers:
         - image: ishantanu16/gourlshortener:v1
           imagePullPolicy: Always
           name: gourlshortener
           ports:
           - containerPort: 8080
           volumeMounts:
           - name: bolt-volumeclaim
             mountPath: /boltdb-path
         volumes:
         - name: bolt-volumeclaim
           emptyDir: {}
   ```

### Scaling approach (considering shared storage):
1. Manually scale replicas of deployment. Example:
   ```sh
   kubectl scale deployment url-shortener --replicas=10
   ```
2. Enable autoscaling for cluster and autoscale the deployment based on metrics. For example, if we've to define an autoscaling policy based on CPU utilization, it would be something like:
   ```sh
   kubectl autoscale deployment url-shortener --min=10 --max=15 --cpu-percent=80
   ```
