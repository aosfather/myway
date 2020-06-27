package jobs

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

//链表简单的实现
type LinkedNode struct {
	value      interface{}
	prev, next *LinkedNode
}

func newList() *LinkedList {
	return new(LinkedList)
}

type LinkedList struct {
	head, tail *LinkedNode
	size       int
}

func (l *LinkedList) push(v interface{}) *LinkedNode {
	n := &LinkedNode{value: v}
	if l.head == nil {
		l.head, l.tail = n, n
		l.size++
		return n
	}

	n.prev = l.tail
	n.next = nil

	l.tail.next = n
	l.tail = n
	l.size++
	return n
}

func (l *LinkedList) remove(n *LinkedNode) {
	if n == nil {
		return
	}

	prev, next := n.prev, n.next
	if prev == nil {
		l.head = next
	} else {
		prev.next = next
	}

	if next == nil {
		l.tail = prev
	} else {
		next.prev = prev
	}
	n = nil // 主动释放内存
	l.size--
}

func (l *LinkedList) String() (s string) {
	s = fmt.Sprintf("[%d]: ", l.size)
	for cur := l.head; cur != nil; cur = cur.next {
		s += fmt.Sprintf("%v <-> ", cur.value)
	}
	s += "<nil>"

	return s
}

/**
  时间轮定时处理
*/
type Slot struct {
	id    int
	tasks *LinkedList
}

func newSlot(id int) *Slot {
	return &Slot{id: id, tasks: newList()}
}

//触发处理函数
type TriggerHandle func() interface{}

//触发器
type Trigger struct {
	id       int64            // 在 slot 中的索引位置
	slotIdx  int              // 所属 slot
	interval time.Duration    // 任务执行间隔
	cycles   int64            // 延迟指定圈后执行
	do       TriggerHandle    // 执行任务
	resCh    chan interface{} // 传递任务执行结果
	repeat   int64            // 任务重复执行次数
}

func (t *Trigger) String() string {
	return fmt.Sprintf("[slot]:%d [interval]:%.fs [repeat]:%d [cycle]:%dth [idx]:%d ",
		t.slotIdx, t.interval.Seconds(), t.repeat, t.cycles, t.id)
}

// 计算 timeout 应在第几圈被执行
func cycle(interval time.Duration, cycleCost int64) (n int64) {
	n = 1 + int64(interval)/cycleCost
	return
}

type TimeWheel struct {
	ticker    *time.Ticker
	tickGap   time.Duration         // 每次 tick 时长
	slotNum   int                   // slot 数量
	curSlot   int                   // 当前 slot 序号
	slots     []*Slot               // 槽数组
	taskMap   map[int64]*LinkedNode // taskId -> taskPtr
	incrId    int64                 // 自增 id
	taskCh    chan *Trigger         // task 缓冲 channel
	lock      sync.RWMutex          // 数据读写锁
	cycleCost int64                 // 周期耗时
}

// 生成 slotNum 个以 tickGap 为时间间隔的时间轮
func NewTimeWheel(tickGap time.Duration, slotNum int) *TimeWheel {
	tw := &TimeWheel{
		ticker:  time.NewTicker(tickGap),
		tickGap: tickGap,
		slotNum: slotNum,
		slots:   make([]*Slot, 0, slotNum),
		taskMap: make(map[int64]*LinkedNode),
		taskCh:  make(chan *Trigger, 100),
		lock:    sync.RWMutex{},
	}
	tw.cycleCost = int64(tw.tickGap * time.Duration(tw.slotNum))
	for i := 0; i < slotNum; i++ {
		tw.slots = append(tw.slots, newSlot(i))
	}

	go tw.turn()

	return tw
}

func (tw *TimeWheel) NewTrigger(interval time.Duration, repeat int64, do TriggerHandle) *Trigger {
	return &Trigger{
		interval: interval,
		cycles:   cycle(interval, tw.cycleCost),
		repeat:   repeat,
		do:       do,
		resCh:    make(chan interface{}, 1),
	}
}

func (tw *TimeWheel) Append(interval time.Duration, repeat int64, do TriggerHandle) {
	trigger := tw.NewTrigger(interval, repeat, do)
	tw.AddTrigger(trigger)
}

func (tw *TimeWheel) AddTrigger(t *Trigger) {
	if t != nil {
		t.slotIdx = tw.convSlotIdx(t.interval)
		t.id = tw.slot2Task(t.slotIdx)
		slot := tw.slots[t.slotIdx]
		if slot != nil {
			node := slot.tasks.push(t)
			tw.taskMap[t.id] = node
		}
	}
}

// 执行延时任务
func (tw *TimeWheel) After(timeout time.Duration, do TriggerHandle) (int64, chan interface{}) {
	if timeout < 0 {
		return -1, nil
	}

	t := tw.NewTrigger(timeout, 1, do)
	tw.locate(t, t.interval, false)
	tw.taskCh <- t
	return t.id, t.resCh
}

// 取消任务
func (tw *TimeWheel) Cancel(tid int64) bool {
	tw.lock.Lock()
	defer tw.lock.Unlock()

	node, ok := tw.taskMap[tid]
	if !ok {
		return false // 任务已执行完毕或不存在
	}

	t := node.value.(*Trigger)
	t.resCh <- nil
	close(t.resCh) // 避免资源泄漏

	slot := tw.slots[t.slotIdx]
	slot.tasks.remove(node)
	delete(tw.taskMap, tid)
	return true
}

// 接收 task 并定时运行 slot 中的任务
func (tw *TimeWheel) turn() {
	idx := 0
	for {
		select {
		case <-tw.ticker.C:
			idx %= tw.slotNum
			tw.lock.Lock()
			tw.curSlot = idx // 锁粒度要细，不要重叠
			tw.lock.Unlock()
			tw.handleSlotTasks(idx)
			idx++
		case t := <-tw.taskCh:
			tw.lock.Lock()
			// fmt.Println(t)
			slot := tw.slots[t.slotIdx]
			tw.taskMap[t.id] = slot.tasks.push(t)
			tw.lock.Unlock()
		}
	}
}

// 计算 task 所在 slot 的编号
func (tw *TimeWheel) locate(t *Trigger, gap time.Duration, restart bool) {
	tw.lock.Lock()
	defer tw.lock.Unlock()
	if restart {
		t.slotIdx = tw.convSlotIdx(gap)
	} else {
		t.slotIdx = tw.curSlot + tw.convSlotIdx(gap)
	}
	t.id = tw.slot2Task(t.slotIdx)
}

// 执行指定 slot 中的所有任务
func (tw *TimeWheel) handleSlotTasks(idx int) {
	var expNodes []*LinkedNode

	tw.lock.RLock()
	slot := tw.slots[idx]
	for node := slot.tasks.head; node != nil; node = node.next {
		task := node.value.(*Trigger)
		task.cycles--
		if task.cycles > 0 {
			continue
		}
		// 重复任务恢复 cycle
		if task.repeat > 0 {
			task.cycles = cycle(task.interval, tw.cycleCost)
			task.repeat--
		}

		// 不重复任务或重复任务最后一次执行都将移除
		if task.repeat == 0 {
			expNodes = append(expNodes, node)
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("task exec paic: %v", err) // 出错暂只记录
				}
			}()

			var res interface{}
			if task.do != nil {
				res = task.do()
			}
			task.resCh <- res
			if task.repeat == 0 {
				close(task.resCh)
			}
		}()
	}
	tw.lock.RUnlock()

	//删除过期的触发器
	tw.lock.Lock()
	for _, n := range expNodes {
		slot.tasks.remove(n)                      // 剔除过期任务
		delete(tw.taskMap, n.value.(*Trigger).id) //
	}
	tw.lock.Unlock()
}

// 在指定 slot 中无重复生成新 task id
func (tw *TimeWheel) slot2Task(slotIdx int) int64 {
	return int64(slotIdx)<<32 + atomic.AddInt64(&tw.incrId, 1) // 保证去重优先
}

// 反向获取 task 所在的 slot
func (tw *TimeWheel) task2Slot(taskIdx int64) int {
	return int(taskIdx >> 32)
}

// 将指定间隔计算到指定的 slot 中
func (tw *TimeWheel) convSlotIdx(gap time.Duration) int {
	timeGap := gap % time.Duration(tw.cycleCost)
	slotGap := int(timeGap / tw.tickGap)
	return int(slotGap % tw.slotNum)
}

func (tw *TimeWheel) String() (s string) {
	for _, slot := range tw.slots {
		if slot.tasks.size > 0 {
			s += fmt.Sprintf("[%v]\t", slot.tasks)
		}
	}
	return
}
