// Copyright © 2015 Couchbase, Inc.
// Copyright © 2011-12 Qtrac Ltd.
//
// This program or package and any associated files are licensed under the
// Apache License, Version 2.0 (the "License"); you may not use these files
// except in compliance with the License. You can get a copy of the License
// at: http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mapserver

type mapServer chan commandData

type commandData struct {
	action  commandAction
	key     string
	value   interface{}
	result  chan<- interface{}
	data    chan<- map[string]interface{}
	updater UpdateFunc
}

type commandAction int

const (
	remove commandAction = iota
	end
	find
	insert
	length
	update
	snapshot
)

type findResult struct {
	value interface{}
	found bool
}

type MapServer interface {
	Insert(string, interface{})
	Delete(string)
	Find(string) (value interface{}, found bool)
	Len() int
	Update(string, UpdateFunc)
	Close() map[string]interface{}
	Snapshot() map[string]interface{}
}

type UpdateFunc func(interface{}, bool) interface{}

func NewMapserver() MapServer {
	sm := make(mapServer) // type mapServer chan commandData
	go sm.run()
	return sm
}

func (sm mapServer) run() {
	store := make(map[string]interface{})
	for command := range sm {
		switch command.action {
		case insert:
			store[command.key] = command.value
		case remove:
			delete(store, command.key)
		case find:
			value, found := store[command.key]
			command.result <- findResult{value, found}
		case length:
			command.result <- len(store)
		case update:
			value, found := store[command.key]
			store[command.key] = command.updater(value, found)
		case snapshot:
			mapCopy := make(map[string]interface{}, len(store))
			for k, v := range store {
				mapCopy[k] = v
			}
			command.data <- mapCopy
		case end:
			close(sm)
			command.data <- store
		}
	}
}

func (sm mapServer) Insert(key string, value interface{}) {
	sm <- commandData{action: insert, key: key, value: value}
}

func (sm mapServer) Delete(key string) {
	sm <- commandData{action: remove, key: key}
}

func (sm mapServer) Find(key string) (value interface{}, found bool) {
	reply := make(chan interface{})
	sm <- commandData{action: find, key: key, result: reply}
	result := (<-reply).(findResult)
	return result.value, result.found
}

func (sm mapServer) Len() int {
	reply := make(chan interface{})
	sm <- commandData{action: length, result: reply}
	return (<-reply).(int)
}

// If the updater calls a mapServer method we will get deadlock!
func (sm mapServer) Update(key string, updater UpdateFunc) {
	sm <- commandData{action: update, key: key, updater: updater}
}

// Close() may only be called once per safe map; all other methods can be
// called as often as desired from any number of goroutines
func (sm mapServer) Close() map[string]interface{} {
	reply := make(chan map[string]interface{})
	sm <- commandData{action: end, data: reply}
	return <-reply
}

// Snapshot returns a snapshot copy of the map
func (sm mapServer) Snapshot() map[string]interface{} {
	reply := make(chan map[string]interface{})
	sm <- commandData{action: snapshot, data: reply}
	return <-reply
}
