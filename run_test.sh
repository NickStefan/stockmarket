services=`cat services`

for service in $services; do
  cd $service && go test && cd ..
done
