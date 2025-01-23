package shell

import (
	"fmt"
	"sync"
)

type bgJob struct {
	id  int
	pid int
	cmd string
}

type jobStore struct {
	mu      sync.Mutex
	jobs    map[int]*bgJob
	nextID  int
}

func newJobStore() *jobStore {
	return &jobStore{jobs: make(map[int]*bgJob)}
}

func (js *jobStore) add(cmd string, pid int) int {
	js.mu.Lock()
	defer js.mu.Unlock()
	js.nextID++
	js.jobs[js.nextID] = &bgJob{id: js.nextID, pid: pid, cmd: cmd}
	return js.nextID
}

func (js *jobStore) remove(id int) {
	js.mu.Lock()
	defer js.mu.Unlock()
	delete(js.jobs, id)
}

func (js *jobStore) list() {
	js.mu.Lock()
	defer js.mu.Unlock()
	for _, j := range js.jobs {
		fmt.Printf("[%d] %d  %s\n", j.id, j.pid, j.cmd)
	}
}
