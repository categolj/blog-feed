#!/bin/bash
set -ex
pack build making/blog-feed --builder gcr.io/paketo-buildpacks/builder:tiny --publish
