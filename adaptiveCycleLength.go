package ictl

type adaptiveCycleLength struct {
	cycleSentTotalSize int
	cycleSentCount     int
}

func (a *adaptiveCycleLength) first() bool {
	return a.cycleSentCount == 0
}

func (a *adaptiveCycleLength) sentKF(size int) {
	a.cycleSentTotalSize = size
	a.cycleSentCount = 1
}

func (a *adaptiveCycleLength) sentDF(size int) {
	a.cycleSentTotalSize += size
	a.cycleSentCount++
}

func (a *adaptiveCycleLength) shouldSendThisDF(size int) bool {
	return a.cycleSentTotalSize/a.cycleSentCount > size
}
