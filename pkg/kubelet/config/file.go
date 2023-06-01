package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"minik8s.io/pkg/apis/core"
	"minik8s.io/pkg/constant"
	"minik8s.io/pkg/kubelet/types"
	"minik8s.io/pkg/podmanager"
	"os"
	"time"

	"minik8s.io/pkg/util/listwatch"
	apiurl "minik8s.io/pkg/apiserver/util/url"
	"github.com/go-redis/redis/v8"
	"encoding/json"
)

// design : maintain a cache as the total pod in the dir
// for the purpose that we can determine which pod have been deleted

type podEventType int

const (
	retryPeriod = 1 * time.Second
)

const (
	podAdd podEventType = iota
	podModify
	podDelete

	eventBufferLen = 10
)

type watchEvent struct {
	fileName  string
	eventType podEventType
}

type sourceFile struct {
	// using the channel in mux
	update chan interface{}

	// file channel
	watch chan *watchEvent

	// listening path
	path string
}

// map from file name to pod name
type FileCache struct {
	PodMap 	map[string]string
	PodMap2 map[string]string
}

func NewSourceFile(ch chan interface{}) {
	cfg := newSourceFile(ch, constant.SysPodDir)
	// init the cache
	fileCache := FileCache{
		PodMap: make(map[string]string),	//file map
		PodMap2: make(map[string]string),	//pod map
	}
	go register(&fileCache)
	cfg.run(&fileCache)
}

func register (fileCache *FileCache){
	go listwatch.Watch(apiurl.PodStatusApplyURL, 
		func(msg *redis.Message){
			var Param core.ScheduleParam
			json.Unmarshal([]byte(msg.Payload), &Param)
			fmt.Println("filecache apply pod: ",Param.RunPod.Name)
			fileCache.PodMap2[Param.RunPod.Name] = Param.RunPod.Name
		})
	//need no update handler
	//go listwatch.Watch(apiurl.PodStatusUpdateURL,)
	go listwatch.Watch(apiurl.PodStatusDelURL, 
		func(msg *redis.Message){
			var podname string
			json.Unmarshal([]byte(msg.Payload), &podname)
			fmt.Println("filecache delete pod: ",podname)
			delete(fileCache.PodMap2, podname)
		})
}

func newSourceFile(ch chan interface{}, path string) *sourceFile {
	return &sourceFile{
		update: ch,
		watch:  make(chan *watchEvent, eventBufferLen),
		path:   path,
	}
}

func (cfg *sourceFile) run(fileCache *FileCache) {
	go func () {
		// start to receive message from sourceFile
		select {
		case e := <-cfg.watch:
			{
				if e.eventType == podAdd {
					// this code is only for file adding
					podUpdate := types.PodUpdate{}
					podUpdate.Op = types.ADD
					podUpdate.Source = types.FileSource
					pod, err := core.ParsePod(e.fileName)
					if err != nil {
						if err.Error() == "error file type" {
							fmt.Println("is not yaml file")
						} else {
							fmt.Println(err)
						}
					}
					fileCache.PodMap[e.fileName] = pod.Name
					podUpdate.Pods = append(podUpdate.Pods, pod)

					cfg.update <- podUpdate
				} else if e.eventType == podDelete {
					pod := core.Pod{}
					pod.Name = fileCache.PodMap[e.fileName]
					podUpdate := types.PodUpdate{}
					podUpdate.Op = types.DELETE
					podUpdate.Pods = append(podUpdate.Pods, &pod)
					podUpdate.Source = types.FileSource
					delete(fileCache.PodMap, e.fileName)
					cfg.update <- podUpdate
				} else if e.eventType == podModify {
					fmt.Println("don't support modify")
				}
			}
		}
	}()
	//polling to check container status
	
	go func () {
		timeout := time.Second * 30
		for {
			fmt.Println("check pod status")
			for _,podname := range fileCache.PodMap2{
				if !podmanager.IsPodRunning(podname) {
					fmt.Println("pod not running",podname)
					podmanager.DelSimpleContainer(podname)
				} else if !podmanager.IsCrashContainer(podname) {
					fmt.Println("contain crash",podname)
					podmanager.DelSimpleContainer(podname)
				}
			}
			time.Sleep(timeout)
		}
	}()
	
	cfg.startWatch(fileCache)
}

func (s *sourceFile) doWatch() error {
	_, err := os.Stat(s.path)
	if err != nil {
		fmt.Println("error in path of source")
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("unable to create inotify: %v", err)
	}
	defer w.Close()

	err = w.Add(s.path)
	if err != nil {
		return fmt.Errorf("unable to create inotify for path %q: %v", s.path, err)
	}

	for {
		select {
		case event := <-w.Events:
			if err = s.produceWatchEvent(&event); err != nil {
				return fmt.Errorf("error while processing inotify event (%+v): %v", event, err)
			}
		case err = <-w.Errors:
			return fmt.Errorf("error while watching %q: %v", s.path, err)
		}
	}
}

func (s *sourceFile) produceWatchEvent(e *fsnotify.Event) error {
	var eventType podEventType
	switch {
	case (e.Op & fsnotify.Create) > 0:
		eventType = podAdd
	case (e.Op & fsnotify.Write) > 0:
		eventType = podModify
	case (e.Op & fsnotify.Chmod) > 0:
		eventType = podModify
	case (e.Op & fsnotify.Remove) > 0:
		eventType = podDelete
	case (e.Op & fsnotify.Rename) > 0:
		eventType = podDelete
	default:
		// Ignore rest events
		return nil
	}

	s.watch <- &watchEvent{e.Name, eventType}
	return nil
}

func ListAllConfig(path string) ([]string, error) {
	s := []string{}
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}

	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := path + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

func (cfg *sourceFile) startWatch(fileCache *FileCache) {
	go func() {
		// start all pod in the dir first
		files, err := ListAllConfig(cfg.path)
		if err != nil {
			fmt.Println("error in List config")
			return
		}
		podSet := []*core.Pod{}
		for _, file := range files {
			// parse to Pod Object
			pod, err := core.ParsePod(file)
			if err != nil {
				if err.Error() == "error file type" {
					continue
				} else {
					fmt.Println(err)
					return
				}
			}
			podSet = append(podSet, pod)
			fileCache.PodMap[file] = pod.Name
		}
		podUpdate := types.PodUpdate{}
		podUpdate.Op = types.SET
		podUpdate.Source = types.FileSource
		podUpdate.Pods = podSet
		cfg.update <- podUpdate
		for {
			cfg.doWatch()
			time.Sleep(retryPeriod)
		}
	}()
}
