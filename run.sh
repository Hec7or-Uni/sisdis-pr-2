#!/usr/bin/env bash

CONNFILE="ms/users.txt"
NUM_LINES=$(wc -l < "$CONNFILE")
NUM_EACH=$(($NUM_LINES / 2))

# Create files for each user
for i in $(seq 1 $NUM_LINES); do
  touch cmd/files/fileP${i}.txt
done

# Initialize the processes
for i in $(seq 1 $NUM_EACH); do
  go run cmd/escritor/main.go $i $CONNFILE cmd/files/fileP${i}.txt &
done

for i in $(seq $(($NUM_EACH + 1)) $NUM_LINES); do
  go run cmd/lector/main.go $i $CONNFILE cmd/files/fileP${i}.txt &
done
