#!/bin/sh

swagger generate server \
	--exclude-main \
	--flag-strategy=pflag \
	-f swagger.yaml
