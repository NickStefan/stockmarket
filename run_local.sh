services=`cat services`

for service in $services; do
  cd $service && go install && cd ..
done

for service in $services; do
  $GOPATH/bin/$service &
done
