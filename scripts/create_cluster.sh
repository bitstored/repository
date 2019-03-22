curl -X POST -u Administrator:password http://127.0.0.1:8091/pools/default/buckets \
-d name=bitstored -d ramQuotaMB=100 -d authType=none \
-d replicaNumber=0