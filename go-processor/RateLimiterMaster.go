package rate_limit

import (
	"bufio"
	"errors"
	"github.com/cheggaaa/pb/v3"
	"go.uber.org/ratelimit"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type RateLimitExecutor struct {
	Limit       int
	Concurrency int
	rateLimiter ratelimit.Limiter
	close       chan error
}

func (r *RateLimitExecutor) Execute(worker RateLimitWorker) {
	workSize := (worker).WorkSize()

	log.Printf("Running ....")
	log.Printf("RateLimit %d", r.Limit)
	log.Printf("Work Size = %d", workSize)
	log.Printf("Concurrency = %d", r.Concurrency)

	SetupCloseHandler(r)
	wg := &sync.WaitGroup{}
	executorWg := &sync.WaitGroup{}
	inputChannel := make(chan interface{})
	outputChannel := make(chan interface{}, r.Concurrency)
	r.rateLimiter = ratelimit.New(r.Limit)
	r.close = make(chan error, r.Concurrency)
	wg.Add(1)
	go r.ingestor(worker, inputChannel, wg)

	executorWg.Add(r.Concurrency)
	for i := 0; i < r.Concurrency; i++ {
		go r.executor(worker, inputChannel, outputChannel, executorWg)
	}
	wg.Add(1)
	go r.outputConsumer(worker, outputChannel, wg)

	go r.takeInputRateLimit()
	executorWg.Wait()
	log.Println("completed execution")
	close(outputChannel)
	wg.Wait()
	worker.Close()
}
func (r *RateLimitExecutor)takeInputRateLimit(){
	for{
		consoleReader := bufio.NewReader(os.Stdin)

		input, err := consoleReader.ReadString('\n')
		if err!=nil{
			if err==io.EOF{
				continue
			}
			log.Println(err)
			continue
		}
		input=strings.Trim(input,"\n")
		newRate,err:=strconv.Atoi(input)
		if err!=nil{
			log.Println(input,"invalid input for rate",err)
			continue
		}
		log.Println("Changing rate of consumption to",newRate)
		r.rateLimiter=ratelimit.New(newRate)
	}
}
func SetupCloseHandler(r *RateLimitExecutor) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		r.close <- errors.New("Ctrl+C Pressed")
		os.Exit(0)
	}()
}
func (r *RateLimitExecutor) ingestor(worker RateLimitWorker, inputChannel chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(inputChannel)
	workSize := worker.WorkSize()
	bar := pb.StartNew(workSize)
	bar.Add(worker.StartIndex())
	for i := worker.StartIndex(); i < workSize; i++ {
		r.rateLimiter.Take()
		input, err := worker.GetInput(i)
		if err != nil {
			r.close <- err
		}
		select {
		case inputChannel <- input:
			bar.Increment()
		case err := <-r.close:
			log.Println(err, i)
			return
		}
	}
	bar.Finish()
}
func (r *RateLimitExecutor) executor(worker RateLimitWorker, inputChannel <-chan interface{}, outputChannel chan<- interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	var output, resp interface{}
	var err error
	for input := range inputChannel {
		resp, err = worker.Work(input)
		if err != nil {
			r.close <- err
			continue
		}
		output, err = worker.ProcessResp(resp)
		if err != nil {
			r.close <- err
			continue
		}
		outputChannel <- output
	}

}
func (r *RateLimitExecutor) outputConsumer(worker RateLimitWorker, outputChannel <-chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for output := range outputChannel {
		worker.HandleOutput(output)
	}
}
