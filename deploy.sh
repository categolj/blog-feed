#!/bin/bash
kapp deploy -a blog-feed -c -f <(kbld -f k8s/blog-feed.yml)
