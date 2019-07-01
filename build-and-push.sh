#!/bin/bash
set -e
cf local stage blog-feed
cf local export blog-feed -r making/blog-feed
docker push making/blog-feed
