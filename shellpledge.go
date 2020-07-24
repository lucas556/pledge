package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)
type tefunc func() (error)
func main() {
	// shellpledge --interval=10 --count=5 --batch=60 --total=200s
	interval := flag.Int("interval", 10, "value of interval")
	count := flag.Int("count", 5, "value of count")
	batch := flag.Int("batch", 60, "value of batch")
	total := flag.Int("total", 100, "value of total")
	flag.Parse()

	openLogfile, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	log.SetOutput(openLogfile)
	//log.New(openLogfile,"logger:",log.Ldate|log.Ltime|log.Lshortfile)
	log.Println(fmt.Sprintf("interval=%d,count=%d,batch=%d,total=%d",*interval,*count,*batch,*total))
	//log.Println(fmt.Sprintf("interval=%d,count=%d,batch=%d,total=%d",*interval,*count,*batch,*total))
	schedle(*interval,*count,*batch,*total,func()(error){
		cmd := exec.Command("lotus-miner sectors pledge")
		buf, err :=cmd.Output()
		log.Println(fmt.Sprintf("任务执行结果返回：%s",buf))
		return err
	})

}
func schedle(interval,count,batch,total int,cc tefunc) {
	totalCount:=1
	for {
		for j := 0; j < count; j++ {
			var err = cc()
			if err!=nil{
				log.Println("任务执行出错:",err)
			}
			log.Println(fmt.Sprintf("任务执行第%d次",totalCount))
			totalCount++
			if (totalCount-total > 0){
				return
			}
			time.Sleep(time.Second * 60 * time.Duration(interval))
		}
		time.Sleep(time.Second * 60 * time.Duration(batch))
	}
}