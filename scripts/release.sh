#!/usr/bin/env bash
set -o errexit; set -o nounset; set -o pipefail

function next_version {
	local next_version_base=$(date --utc "+%Y.%-m")
	local num_released_this_month=$(git tag -l | grep -E "^${next_version_base}" | wc -l)

	echo "${next_version_base}.${num_released_this_month}"
}

function main {
	local last_release=${OVERRIDE_LAST_VERSION:-$(git tag -l | tail -n 1)}
	local new_release=$(next_version)

	echo "Releasing version ${new_release}"

	EMAIL=$(git config user.email) gbp dch --release --new-version=${new_release}

	git commit -am "Release version ${new_release}"
	git tag "${new_release}"
}

main
