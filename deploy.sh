#! /bin/bash

cp -R ./templates ./bin/

mkdir ./bin/videos

cd bin

nohup ./api &
nphup ./scheduler &
nohup ./streamserver &
nohup ./web &

echo "deploy finished"