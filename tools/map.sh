#!/bin/sh

# Copyright 2011 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This code parses mapping.old and finds a correspondance from the old
# urls (e.g. #42) to the corresponding path (e.g. /concurrency/3).

function findURL {
	title="$1"
	file=$(grep -l "* $title\$" *.article)
	if [[ -z $file ]]
	then
		echo "undefined"
		return 1
	fi
	titles=$(grep "^* " $file | awk '{print NR, $0}')
	page=$(echo "$titles" | grep "* $title\$" | awk '{print $1}')
	if [[ $(echo "$page" | wc -l) -gt "1" ]]
	then
		echo "multiple matches found for $title; find 'CHOOSE BETWEEN' in the output" 1>&2
		page="CHOOSE BETWEEN $page"
	fi

	page=$(echo $page)
	lesson=$(echo "$file" | rev | cut -c 9- | rev)
	echo "'/$lesson/$page'"
	return 0
}

mapping=`cat mapping.old`

pushd ../content
echo "$mapping" | while read page; do
	num=$(echo "$page" | awk '{print $1}')
	title=$(echo "$page" | sed "s/[0-9]* //")
	url=$(findURL "$title")
	echo "    '#$num': $url, // $title"
done
popd > /dev/null
