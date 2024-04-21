

PS D:\Nuwas Cisco Laptop\Aggregator\Aggregator\release-k8s\Cdc\kafka_sample_client_service> Docker build -t nponnamb/kafka_sample_client:0.0.1 .
PS D:\Nuwas Cisco Laptop\Aggregator\Aggregator\release-k8s\Cdc\kafka_sample_client_service> docker push nponnamb/kafka_sample_client:0.0.1

[root@k8smaster debezium]# kubectl apply -f deployment.yaml


[root@k8smaster deploy]# curl -H "Accept:application/json" 192.168.0.203:8083
{"version":"3.6.1","commit":"5e3c2b738d253ff5","kafka_cluster_id":"6PMpHYL9QkeyXRj9Nrp4KA"}
[root@k8smaster deploy]#


[root@k8smaster deploy]# curl -X POST -H "Content-Type: application/json" -d '{"content":"Kafka007 Hello, Shalna How are you!!!!!"}' http://10.97.119.63:8080/message
{"status":"Send Message"}
[root@k8smaster deploy]#



 kubectl -n cdc logs pod/kafka-client-k8s-579b4fcc6d-kjtmw -f