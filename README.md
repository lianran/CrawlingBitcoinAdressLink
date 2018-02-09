### How to crawling bitcoin address with its web links

#### Download the bitcoin block files

Maybe you can download the bitcoin block files with the bitcoin client, or you can download the blcok files through the [bt](https://getbitcoinblockchain.com/), find the files from the internet?

	

So you will get the block files named blkxxxxx.dat, the xxxxx is the number from 00001 to 01000 mabye.


#### Extract address from the block files

1. Go to the **bitcoin-blk-file-reader** folder and open the **get_the_addresses.sh**, you should change the path to your path where u store the block files *in line 27* and where u want to save the address *in line 28*:

	```
	echo "CMD::python analyze.py /path/where/u/store/the/block/files/${name}.dat | grep Address | sort | uniq | sed \"s/> Address: //g\" > ${name}.txt"
	python analyze.py /path/where/u/want/to/save/address/${name}.dat | grep Address | sort | uniq | sed "s/> Address: //g" > ${name}.txt
	```

2. open a terminal to run the script **get_the_addresses.sh**

	```
	./get_the_addresses.sh
	```

After u finish this, u can get new file named blockxxxxx.txt whic store the addresses extracted from the blockxxxxx.dat.

#### Crawling the addresses's links

1. Go to the **crawling** folder and then open the crawling file **crawling_address.go**, you should set the path of address file *at the line 372* and the path for save the crawling result *at the line 374*

2. Run the crawling file:

	```
	go run crawling_address.go
	```
	
