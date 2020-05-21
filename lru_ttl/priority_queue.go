package lruttl

type priorityQueue []*lruEntry

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].ExpiryTime.Before(pq[j].ExpiryTime)
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].queueIdx = i
	pq[j].queueIdx = j
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	item.queueIdx = -1
	return item
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*lruEntry)
	item.queueIdx = len(*pq)
	*pq = append(*pq, item)
}
