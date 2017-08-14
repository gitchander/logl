package logl

import "sync"

type fakeLocker struct{}

func (fakeLocker) Lock()   {}
func (fakeLocker) Unlock() {}

var _ sync.Locker = fakeLocker{}

func getLocker(notSafe bool) sync.Locker {
	if notSafe {
		return fakeLocker{}
	} else {
		return new(sync.Mutex)
	}
}
