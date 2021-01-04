#!/bin/sh
mkdir build
cd build
mkdir onbio
cd onbio
mkdir bin
mkdir conf
mkdir tools
cp  ../../bin/*  bin/
cp  ../../conf.json conf/
cp  ../../tools/*  tools

cd ..

tar -zcvf onbio.tgz onbio

