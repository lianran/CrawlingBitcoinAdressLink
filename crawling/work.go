package crawling

func worker(id int, wg *sync.WaitGroup, jobs <-chan int, results chan<- ThreadResult) {
	for j := range jobs {
		//fmt.Println("worker", id, "processing job", j)
		//time.Sleep(time.Second)
		stt := time.Now()
		cell, bonus, _ := bonusFunc3(j)
		results <- ThreadResult{cell, bonus, time.Now().Sub(stt).Seconds()}
		wg.Done()
	}
}

//计算idxs列表里每个probe的bonus
func calcBonus1(idxs []int) {
	var wg sync.WaitGroup
	jobs := make(chan int, 200)
	results := make(chan ThreadResult, 300)
	for w := 1; w <= 4; w++ {
		go worker(w, &wg, jobs, results)
	}

	go func() {
		for _, j := range idxs {
			jobs <- j
			wg.Add(1)
		}
		close(jobs)
	}()

	wg.Wait()
	//fmt.Println("worker finished")
	total_time := 0.0
	for a := 1; a <= len(idxs); a++ {
		result := <-results
		M_bonus[result.cell] = result.bonus
		total_time += result.ts
	}
	//fmt.Printf("execute over, average ts : %f\n", total_time/float64(len(idxs)))
}
