#!/bin/bash -e

cd ..
VERSION=$(semver get release)
MESSAGE=$(git log -1 --pretty=%B)

echo "creating version $VERSION"

git tag -a $VERSION -m "$MESSAGE"
git push origin $VERSION