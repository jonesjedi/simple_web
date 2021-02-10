#!/bin/sh
mkdir build
cd build
mkdir onbio
cd onbio
mkdir bin
mkdir conf
mkdir tools
mkdir logs
cp  ../../bin/*  bin/
cp  ../../conf/* conf/
cp  ../../tools/*  tools

cd ..

tar -zcvf onbio.tgz onbio

