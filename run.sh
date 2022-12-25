#!/bin/bash

cd interpreter; go build -o interpreter; cd ..
cd compiler; go build -o compiler; cd ..

mkdir out

./compiler/compiler programs/hanoi.bf > out/hanoi.ll; clang -O3 out/hanoi.ll -o out/hanoi; ./out/hanoi
./compiler/compiler programs/mb.bf > out/mb.ll; clang -O3 out/mb.ll -o out/mb; ./out/mb
./compiler/compiler programs/hello-world.bf > out/hello-world.ll; clang -O3 out/hello-world.ll -o out/hello-world; ./out/hello-world