#!/bin/bash
st_num=$1
ed_num=$2

if [[ $# -lt 2 ]]; then
	echo "plz enter start num and end num";
	exit;
fi

for((i=${st_num};i<=${ed_num};i++));
do

if [[ $i -lt 10 ]]; then
	name=blk0000$i;
elif [[ $i -lt 100 ]]; then
	name=blk000$i;
elif [[ $i -lt 1000 ]]; then
	name=blk00$i;
elif [[ $i -lt 10000 ]]; then
	name=blk0$i;
else
	name=blk$i;
fi;

echo "processing the file ${name}"

echo "CMD::python analyze.py /Users/anranli/Download/Bitcoin/blocks/${name}.dat | grep Address | sort | uniq | sed \"s/> Address: //g\" > ${name}.txt"
python analyze.py /Users/anranli/Download/Bitcoin/blocks/${name}.dat | grep Address | sort | uniq | sed "s/> Address: //g" > ${name}.txt

done