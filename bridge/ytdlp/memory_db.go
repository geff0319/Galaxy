package ytdlp

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/ge-fei-fan/gefflog"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/google/uuid"
)

// In-Memory Thread-Safe Key-Value Storage with optional persistence
type MemoryDB struct {
	table sync.Map
}

// Get a process pointer given its id
func (m *MemoryDB) Get(id string) (*Process, error) {
	entry, ok := m.table.Load(id)
	if !ok {
		return nil, errors.New("no process found for the given key")
	}

	return entry.(*Process), nil
}

// Store a pointer of a process and return its id
func (m *MemoryDB) Set(process *Process) string {
	id := uuid.NewString()

	m.table.Store(id, process)
	process.Id = id

	return id
}

// Removes a process progress, given the process id
func (m *MemoryDB) Delete(id string) {
	m.table.Delete(id)
}

func (m *MemoryDB) Keys() *[]string {
	var running []string

	m.table.Range(func(key, value any) bool {
		running = append(running, key.(string))
		return true
	})

	return &running
}

// Returns a slice of all currently stored processes progess
func (m *MemoryDB) All() *[]ProcessResponse {
	running := []ProcessResponse{}

	m.table.Range(func(key, value any) bool {
		running = append(running, ProcessResponse{
			Id:       key.(string),
			Url:      value.(*Process).Url,
			Info:     value.(*Process).Info,
			Progress: value.(*Process).Progress,
			Output:   value.(*Process).Output,
			Params:   value.(*Process).Params,
			//BiliMeta: value.(*Process).BiliMeta,
		})
		return true
	})
	sort.Slice(running, func(i, j int) bool {
		return running[i].Info.CreatedAt.After(running[j].Info.CreatedAt)
	})
	return &running
}

// Persist the database in a single file named "session.dat"
func (m *MemoryDB) Persist(basePath string) error {
	running := m.All()
	fmt.Println(running)
	sf := filepath.Join(basePath, "/data/session.dat")

	fd, err := os.Create(sf)
	if err != nil {
		return errors.Join(errors.New("failed to persist session"), err)
	}

	session := Session{Processes: *running}

	if err = gob.NewEncoder(fd).Encode(session); err != nil {
		return errors.Join(errors.New("failed to persist session"), err)
	}
	return nil
}

// Restore a persisted state
func (m *MemoryDB) Restore(basePath string, mq *MessageQueue) {
	fmt.Println("读取ytdlp下载内容")
	fd, err := os.Open(filepath.Join(basePath, "/data/session.dat"))
	if err != nil {
		return
	}
	fmt.Println(fd)
	var session Session

	if err = gob.NewDecoder(fd).Decode(&session); err != nil {
		gefflog.Err("Decode session error: " + err.Error())
		return
	}

	for _, proc := range session.Processes {
		restored := &Process{
			Id:       proc.Id,
			Url:      proc.Info.URL,
			Info:     proc.Info,
			Progress: proc.Progress,
			Output:   proc.Output,
			Params:   proc.Params,
			//BiliMeta: proc.BiliMeta,
			//Logger:   logger,
		}

		m.table.Store(proc.Id, restored)

		//if restored.Progress.Status != StatusCompleted {
		//	mq.Publish(restored)
		//}
		if restored.Progress.Status != StatusCompleted {
			restored.Progress.Status = StatusErrored
		}
	}
}

func (m *MemoryDB) IsProcessExist(process *Process) bool {
	isExist := false
	YdpConfig.Mdb.table.Range(func(key, value any) bool {
		if value.(*Process).Id == process.Id || value.(*Process).Info.Id != process.Info.Id {
			return true
		}
		isExist = true
		return false
	})
	return isExist
}
