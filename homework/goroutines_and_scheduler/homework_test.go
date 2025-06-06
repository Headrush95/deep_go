package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

type item struct {
	value    Task
	priority int
	index    int
}

type pqMaxHeap []*item

func newMaxHeapPriorityQueue(size int) *pqMaxHeap {
	pq := make(pqMaxHeap, 0, size)
	heap.Init(&pq)
	return &pq
}

func (pq *pqMaxHeap) Len() int { return len(*pq) }

func (pq *pqMaxHeap) Less(i, j int) bool {
	if (*pq)[i] == nil || (*pq)[j] == nil {
		return i > j
	}
	return (*pq)[i].priority > (*pq)[j].priority
}

func (pq *pqMaxHeap) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
	if (*pq)[i] != nil {
		(*pq)[i].index = i
	}
	if (*pq)[j] != nil {
		(*pq)[j].index = j
	}
}

func (pq *pqMaxHeap) Push(x any) {
	n := pq.Len()
	item, successfully := x.(*item)
	if !successfully {
		return
	}
	item.index = n
	*pq = append(*pq, item)
}

func (pq *pqMaxHeap) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // избегаем утечки памяти
	item.index = -1 // для безопасности
	*pq = old[0 : n-1]
	return item
}

func (pq *pqMaxHeap) Update(item *item, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}

type Scheduler struct {
	queue      *pqMaxHeap
	uniqueTask map[int]*item
}

const (
	magicSize = 1_024 // ну не могу я с 0 инициализировать) И в сигнатуру конструктора как будто неправильно это вносить
)

type Task struct {
	Identifier int
	Priority   int
}

func NewScheduler() Scheduler {
	sc := Scheduler{
		queue:      newMaxHeapPriorityQueue(magicSize),
		uniqueTask: make(map[int]*item, magicSize),
	}
	return sc
}

func (s *Scheduler) AddTask(task Task) {
	if _, exist := s.uniqueTask[task.Identifier]; exist {
		// можно, конечно, тут обновлять приоритет задачи, но для таких целей в API есть ChangeTaskPriority
		return
	}
	taskItem := &item{
		value:    task,
		priority: task.Priority,
	}
	heap.Push(s.queue, taskItem)
	s.uniqueTask[task.Identifier] = taskItem
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	taskItem, exist := s.uniqueTask[taskID]
	if !exist {
		return
	}
	s.queue.Update(taskItem, newPriority)
}

func (s *Scheduler) GetTask() Task {
	if len(*s.queue) == 0 {
		return Task{}
	}
	taskItem := heap.Pop(s.queue)
	return taskItem.(*item).value
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	// а почему тут мы обновляем приоритет у задачи, но при этом ожидаем старое значение. Как-то не логично...
	// но какое ТЗ имеем, с таким и работаем)
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
