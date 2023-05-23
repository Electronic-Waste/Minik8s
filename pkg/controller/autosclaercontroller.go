package controller

type DeploymentController struct {
	//Client
	//util

	// work queue
	queue   *queue.Queue
	nameMap map[interface{}]interface{}
	//channel chan struct{}
	//message *redis.Message
}