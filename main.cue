package main

import (
	"dagger.io/dagger"

	"universe.dagger.io/go"
)

dagger.#Plan & {

	client: filesystem: ".": read: exclude: [
		"bin",
	]
	client: filesystem: "./bin": write: contents: actions.build.output
	actions: {
		_source: client.filesystem["."].read.contents

		build: go.#Build & {
			source:  _source
			package: "."
			os:      client.platform.os
			arch:    client.platform.arch


			env: {
				CGO_ENABLED: "0"
			}
		}

    }
}