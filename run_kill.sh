services=`cat services`

for service in $services; do
  pkill -x $service
done
