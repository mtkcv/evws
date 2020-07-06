package main

import (
	"syscall"
)

const (
	//PollIn POLLIN 1
	PollIn = syscall.EPOLLIN
	//PollHup POLLHUP 16
	PollHup = syscall.EPOLLHUP
	//PollOut POLLOUT 4
	PollOut = syscall.EPOLLOUT
	//PollErr POLLERR 8
	PollErr = syscall.EPOLLERR
)

//Poller poller
type Poller struct {
	fd int
}

//NewPoller poller constructor
func NewPoller() *Poller {
	ep := &Poller{}
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		return nil
	}
	ep.fd = fd

	return ep
}

//Add a fd
func (ep *Poller) Add(fd int) error {
	err := syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{Events: syscall.EPOLLIN, Fd: int32(fd)})
	if err != nil {
		return err
	}
	return nil
}

//Delete a fd
func (ep *Poller) Delete(fd int) error {
	err := syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	return nil
}

//Wait fd events
func (ep *Poller) Wait(callback func(fd, ev int)) error {
	events := make([]syscall.EpollEvent, 128)
	n, err := syscall.EpollWait(ep.fd, events, 128)
	if err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		callback(int(events[i].Fd), int(events[i].Events))
	}
	return nil
}

//Close poll
func (ep *Poller) Close() error {
	return syscall.Close(ep.fd)
}
