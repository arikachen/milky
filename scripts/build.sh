#!/bin/bash

OUTPUT=_output
go build -o $OUTPUT/milky-ctrl ./cmd/ctrl &
go build -o $OUTPUT/milky-agent ./cmd/agent &
go build -o $OUTPUT/openshift-sdn ./cmd/cni &
wait
