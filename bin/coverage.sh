#!/usr/bin/env bash
shopt -s dotglob
outdir=$(mktemp -d)

for pkg in '' ; do
  go test \
    -covermode=atomic \
    -coverprofile="$outdir"/"$pkg".out \
    ./"$pkg"
done

cat - - <<<'mode: atomic' <(tail -n +2 -q "$outdir"/*.out) > coverage.out
rm -rf "$outdir"
