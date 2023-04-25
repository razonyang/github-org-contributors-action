#!/bin/sh -l
echo $1
echo $2
/org-contributors -org=$1 -output=$2
