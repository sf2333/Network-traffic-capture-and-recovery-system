package capture

import "sync"

type Queue struct {
	list []string
	Size int
	sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		list: make([]string, 0),
		Size: 0,
	}
}

func (q Queue) List() []string {
	tmp := make([]string,len(q.list))
	copy(tmp,q.list)
	return tmp
}

func (q *Queue) Front() string {
	if q.Size <= 0 {
		return ""
	}

	num := q.list[0]
	return num
}

func (q *Queue) Push(data string) {
	q.list = append(q.list, data)
	q.Size += 1
}

func (q *Queue) Pop() bool {
	if q.Size <= 0 {
		return false
	}

	q.Size--
	q.list = q.list[1:]
	return true
}

//对指定的值进行重置
//如果该值已存在在队列中，则先删除，然后加到队列尾部
func (q *Queue) ResetValue(value string) {
	q.Lock()
	for i, v := range q.list {
		if v == value {
			q.list = removeSlice(q.list, i)
			break
		}
	}

	q.list = append(q.list, value)
	q.Unlock()
}

func (q *Queue) RemoveValue(value string) {
	q.Lock()
	for i, v := range q.list {
		if v == value {
			q.Size--
			q.list = removeSlice(q.list, i)
			break
		}
	}
	q.Unlock()
}

func removeSlice(list []string, index int) []string {
	return append(list[:index], list[index+1:]...)
}
