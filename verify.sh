#!/bin/bash
FILE="100MB.bin"
var=`ls /tmp/*${FILE}.part | wc -l`
cat `eval echo /tmp/{0..$(($var-1))}_${FILE}.part` > /tmp/temp.bin && rm -f /tmp/*${FILE}.part && md5sum /tmp/temp.bin
rm /tmp/temp.bin
md5sum /tmp/${FILE}
