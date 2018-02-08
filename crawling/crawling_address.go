package main

// Example
import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hunterhug/GoSpider/query"
	"github.com/hunterhug/GoSpider/spider"
	"io"
	"os"
	//"runtime"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	//some urls
	bad_url = []string{"blockchain.info", "webbtc.com", "bitcoinchain.com/block_explorer",
		"bitcoinwhoswho.com", "blockinfo", "hybrid-analysis", "paste4btc", "blkdb.cn", "learnmeabitcoin",
		"allinvain-addresses", "walletexplorer", "transactions-addresses", "tenthousandblocks", "www.coinprism.info",
		"block.bitbank.com", "block.okcoin.cn", "www.blocktrail.com", "live.blockcypher.com/btc/tx",
		"live.blockcypher.com", "/tx/", "/block/", "api.blockcypher.com"}
	d_time = time.Second / 1000
)

//get item's url from html
func get_item_url(html []byte, ret *map[string]int) error {
	doc, err := query.QueryBytes(html)
	if err != nil {
		spider.Log().Errorf("create doc err")
		return err
	}
	doc.Find(".b_algo").Each(func(i int, item *goquery.Selection) {
		title := item.Find("h2 a").Text()
		href, _ := item.Find("h2 a").Attr("href")
		spider.Log().Debugf("%d %s %s", i, title, href)
		if href != "" {
			(*ret)[href] = 0
		}
	})

	return nil
}

//get page's url from html
func get_page_url(html []byte) (map[string]int, error) {
	ret := make(map[string]int)
	doc, err := query.QueryBytes(html)
	if err != nil {
		spider.Log().Errorf("create doc err")
		return ret, err
	}
	doc.Find(".b_pag li").Each(func(i int, pg *goquery.Selection) {
		//var ntPage string[]
		href, _ := pg.Find("a").Attr("href")
		href = strings.Replace(href, " ", "", -1)
		href = strings.Replace(href, "/n", "", -1)

		if href != "" {
			spider.Log().Debugf("%d %s", i, href)
			ret[href] = 0
		}
	})

	return ret, nil
}

//get one url's html
func set_url_and_get_html(url string, sp *spider.Spider) ([]byte, error) {

	//set url
	sp.SetUrl(url).SetUa(spider.RandomUa()).SetMethod(spider.GET)

	//fetch
	html, err := sp.Go()
	if err != nil {
		spider.Log().Errorf(err.Error())
	}
	return html, err
}

//get website's content without spider
func get_html_from_url(url string) ([]byte, error) {
	sp, err := spider.New(nil)
	if err != nil {
		spider.Log().Errorf("new spider error %s", err.Error())
		return nil, err
	}

	sp.SetUrl(url).SetUa(spider.RandomUa()).SetMethod(spider.GET)
	// 3. Fetch
	html, err := sp.Go()
	if err != nil {
		spider.Log().Errorf("fetch err %s", err.Error())
	}
	return html, err
}

func get_url_from_bing(address string) (map[string]int, error) {

	ret := make(map[string]int)

	// 1. New a spider
	sp, _ := spider.New(nil)
	// 2. Set a URL
	//url := "https://www.bing.com/search?q=1M72Sfpbz1BPpXFHz9m3CdqATR44Jvaydd"
	main_url := "https://global.bing.com/search?q="
	//https://www.bing.com/search?q=1111111111111111111114oLvT2&qs=n&form=QBLH&sp=-1&sc=1-27&sk=&cvid=5AE66ADEFA6545C3A77D3D937B32B35B
	end_url := "&qs=HS&sc=8-0&cvid=7B8A327FCD4B455594B448CBE35EB5DA&FORM=QBLH&sp=1"
	url := main_url + address + end_url
	// 3. fetch the first page
	html, err := set_url_and_get_html(url, sp)
	if err != nil {
		spider.Log().Errorf(err.Error())
		return ret, err
	}
	// 4. get the page's item url
	err = get_item_url(html, &ret)
	if err != nil {
		spider.Log().Errorf(err.Error())
		return ret, err
	}
	// 5.get the page's index url
	pgs, _ := get_page_url(html)

	// 6.get the other pages info
	spider.Log().Debugf("get %d url and %d index from first page", len(ret), len(pgs))
	for nxturl, _ := range pgs {
		url = "https://www.bing.com" + nxturl
		html, err = set_url_and_get_html(url, sp)
		if err != nil {
			spider.Log().Errorf(err.Error())
		}
		err = get_item_url(html, &ret)
		if err != nil {
			spider.Log().Errorf(err.Error())
		}
		spider.Log().Debugf("after processed the url: %s\n we have %d urls\n", url, len(ret))
		time.Sleep(time.Duration(rand.Intn(100)) * d_time)
	}
	return ret, nil
}

func url_is_ok(address string, url string) bool {

	if strings.Contains(url, address) {
		return false
	}
	for _, s := range bad_url {
		if strings.Contains(strings.ToUpper(url), strings.ToUpper(s)) {
			return false
		}
	}
	return true
}

func url_review(address string, urls *map[string]int) error {
	t_urls := (*urls)

	http_num := 0
	https_num := 0

	for url, _ := range t_urls {
		spider.Log().Debugf("dealing with url %s", url)
		//1. check url	if contain the address
		if url_is_ok(address, url) == false {
			spider.Log().Debugf("the url contains the address or is a bad one")
			continue
		}

		//2. check the page's content if contain the address
		// html, err := get_html_from_url(url)
		// if err != nil {
		// 	spider.Log().Errorf("cann't fetch this url: %s. Err: %s", url, err.Error())
		// 	continue
		// }
		// if strings.Contains(string(html), address) == false {
		// 	spider.Log().Noticef("html don't contain the address")
		// 	continue
		// }

		if strings.Contains(url, "https") {
			t_urls[url] = 1
			https_num += 1
		} else {
			t_urls[url] = 2
			http_num += 1
		}
	}
	t_urls["http"] = http_num
	t_urls["https"] = https_num
	return nil
}

type ThreadResult struct {
	ads  string
	urls map[string]int
}

func worker(id int, wg *sync.WaitGroup, jobs <-chan string, results chan<- ThreadResult) {
	for ads := range jobs {
		spider.Log().Debugf(ads)
		urls, err := get_url_from_bing(ads)
		if err != nil {
			spider.Log().Errorf(err.Error())
			continue
		}
		url_review(ads, &urls)
		spider.Log().Debug("addreess:%s, all the urls num: %d, http: %d, https: %d", ads, len(urls)-2, urls["http"], urls["https"])
		if urls["http"] != 0 || urls["https"] != 0 {
			spider.Log().Debugf("get nothing about this address")
			results <- ThreadResult{ads, urls}
		}
		time.Sleep(time.Duration(rand.Intn(10)) * d_time)
		wg.Done()
	}
}

func deal_with_file(InFile string, OutFile string) {

	spider.Log().Debugf("begin deal with")
	f, err := os.Open(InFile)
	if err != nil {
		spider.Log().Errorf(err.Error())
		return
	}
	defer f.Close()
	var ads []string
	buff := bufio.NewReader(f)
	for {
		line, err := buff.ReadString('\n')
		if io.EOF == err {
			break
		}

		if err != nil {
			spider.Log().Errorf(err.Error())
			continue
		}

		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, " ", "", -1)
		ads = append(ads, line)
	}
	//spider.Log().Debugf("read file successed")
	//fmt.Println(ads)
	var wg sync.WaitGroup
	jobs := make(chan string, 1000)
	results := make(chan ThreadResult, 1000)
	for w := 1; w <= 20; w++ {
		go worker(w, &wg, jobs, results)
	}

	//fmt.Println("goggo")
	go func() {
		count := 0
		for _, address := range ads {
			//spider.Log().Noticef(address)
			jobs <- address
			count += 1
			wg.Add(1)
			if count%1000 == 0 {
				spider.Log().Noticef("time %s totoal: %d finished: %d  precent: %f",
					time.Now().String(), len(ads), count, float64(count)/float64(len(ads)))
			}
		}
		close(jobs)
	}()

	//not use this ?
	//wg.Wait()
	fmt.Println("ggg")
	ff, err := os.Create(OutFile)
	if err != nil {
		spider.Log().Errorf(err.Error())
		return
	}
	defer ff.Close()
	cnt := 0
	go func() {
		for rslt := range results {
			//save this to file, it may need mult-process
			spider.Log().Noticef("get one: %s http: %d: https:%d", rslt.ads, rslt.urls["http"], rslt.urls["https"])
			_, err = ff.WriteString(rslt.ads)
			if err != nil {
				spider.Log().Errorf("write %s err: %s", ads, err.Error())
				continue
			}
			//del the key: http and https
			delete(rslt.urls, "http")
			delete(rslt.urls, "https")

			for url, v := range rslt.urls {
				//don't record the bad url
				if v == 0 {
					continue
				}
				_, err = ff.WriteString("\t" + url)
				if err != nil {
					spider.Log().Errorf("write %s err: %s", url, err.Error())
					continue
				}
			}
			_, err = ff.WriteString("\n")
			if err != nil {
				spider.Log().Errorf("write %s err: %s", "enter", err.Error())
				continue
			}
			cnt += 1
			if cnt%100 == 0 {
				f.Sync()
			}
		}
	}()

	wg.Wait()
	close(results)

}

func deal_with_collected_address_single(file string) map[string]map[string]int {
	spider.SetLogLevel("Debug")

	ret := make(map[string]map[string]int)
	//read the file
	f, err := os.Open(file)
	if err != nil {
		spider.Log().Errorf(err.Error())
	}
	defer f.Close()
	buff := bufio.NewReader(f)

	//deal with each url
	for {
		line, err := buff.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, " ", "", -1)
		spider.Log().Debugf(line)
		//get the url form bing
		urls, err := get_url_from_bing(line)
		if err != nil {
			spider.Log().Errorf(err.Error())
			continue
		}
		// review the urls
		url_review(line, &urls)
		spider.Log().Noticef("addreess:%s, all the urls num: %d, http: %d, https: %d", line, len(urls)-2, urls["http"], urls["https"])
		if urls["http"] == 0 && urls["https"] == 0 {
			spider.Log().Debugf("get nothing about this address")
			continue
		}
		ret[line] = urls
	}
	return ret
}

func main() {
	spider.SetLogLevel("Info")
	//
	startCrawling()
	//test()
}

func startCrawling() {
	//this is the path to read the address
	InFile := "/home/bitcoin/addresses/blk00973.txt"
	//this is the path to save the result
	OutFile := "/home/bitcoin/urls/blk00973.txt"
	deal_with_file(InFile, OutFile)
}

func test2() {
	InFile := "/Users/anranli/script/bitcoin-blk-file-reader/blk00000.txt"
	OutFile := "/Users/anranli/script/bitcoin-blk-file-reader/blk00000_.dat"
	//deal_with_collected_address_single(file)
	deal_with_file(InFile, OutFile)
}

func test() {
	spider.SetLogLevel("Info")
	address := "1111111111111111111114oLvT2"
	ret, err := get_url_from_bing(address)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("we get all the urls:")
	fmt.Println(len(ret), ret)

	url_review(address, &ret)

	fmt.Println("we have reviewed urls:")
	fmt.Println(len(ret), ret)

	https_num := 0
	http_num := 0
	for url, v := range ret {
		if v == 1 {
			https_num += 1
			fmt.Printf("https: %s\n", url)
		} else if v == 2 {
			http_num += 1
			fmt.Printf("http: %s\n", url)
		}
	}

	fmt.Printf("https: %d\t http: %d\n", https_num, http_num)
	fmt.Printf("https: %d\t http:%d\n", ret["https"], ret["http"])
}
