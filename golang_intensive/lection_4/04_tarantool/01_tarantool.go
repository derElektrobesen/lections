package main

import (
	"log"
	"time"

	tarantool "github.com/tarantool/go-tarantool"
	"github.com/tarantool/go-tarantool/queue"
)

var db *tarantool.Connection

func init() {
	opts := tarantool.Opts{
		Timeout:       time.Second,
		Reconnect:     time.Second,
		MaxReconnects: 5,
		User:          "rs_events_user",
		Pass:          "V8SkhukQ8lbP8QxWyHFvP7UKzo",
	}

	var err error
	db, err = tarantool.Connect("isaev1.e.smailru.net:6766", opts)

	if err != nil {
		log.Fatalf("Can't connect to tarantool: %s", err)
	}
}

func testQueue() {
	q := queue.New(db, "test_queue")
	err := q.Create(queue.Cfg{
		Temporary:   true,
		IfNotExists: true,
		Kind:        queue.FIFO_TTL,
		Opts: queue.Opts{
			Ttl:   100 * time.Second, // как долго таск может жить в очереди (таски старше Ttl удаляются)
			Ttr:   5 * time.Second,   // как долго таск может обрабатываться воркером (если дольше, другой воркер может его взять)
			Delay: 30 * time.Second,  // на сколько таск можно отложить в очередь
			Pri:   1,
		},
	})

	if err != nil {
		log.Fatalf("Can't create queue: %s", err)
	}

	defer q.Drop()

	task, err := q.PutWithOpts([]interface{}{time.Now().Unix(), "test task"}, queue.Opts{Delay: 5 * time.Second})
	if err != nil {
		log.Fatalf("Can't put task: %s", err)
	}

	log.Printf("Task added into queue: %+v", task.Data())
	log.Println("Waiting for task, attempt 1...")

	task, err = q.TakeTimeout(3 * time.Second)
	if err != nil {
		log.Fatalf("Can't take task with timeout: %s", err)
	}

	if task != nil {
		log.Fatalf("Unexpected task found: %+v", task.Data())
	}

	log.Println("Waiting for task, attempt 2...")

	task, err = q.TakeTimeout(100 * time.Second)
	if err != nil {
		log.Fatalf("Can't take task with timeout: %s", err)
	}

	if task == nil {
		log.Fatalf("Task wasn't found")
	}

	log.Printf("Task taken: %+v", task.Data())

	task.Ack()
}

func main() {
	space := db.Schema.Spaces["test_space"]
	if space == nil {
		log.Fatalf("Space not found")
	}

	resp, err := db.Insert(space, []interface{}{time.Now().Unix(), "test", "xxxx", 22})
	if err != nil {
		log.Fatalf("Can't store tuple in space: %s", err)
	}

	log.Printf("Response is %+v", resp.Tuples())

	resp, err = db.Call("test_function", []interface{}{"test"})
	if err != nil {
		log.Fatalf("Can't call test function: %s", err)
	}

	log.Printf("Response for function call is %+v", resp.Tuples())

	testQueue()
}
