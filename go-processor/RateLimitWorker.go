package rate_limit

type RateLimitWorker interface {
	// total number of work to be done
	WorkSize() int

	// get input for index
	GetInput(index int) (interface{}, error)

	StartIndex() int
	// do work for one input
	Work(interface{}) (interface{}, error)

	//process resp
	ProcessResp(resp interface{}) (interface{}, error)

	//Handle the output
	HandleOutput(output interface{})

	//post Handling the output
	Close()
}
