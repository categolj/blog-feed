#!/bin/bash
set -ex
pack build ghcr.io/categolj/blog-feed --builder gcr.io/paketo-buildpacks/builder:tiny --publish
