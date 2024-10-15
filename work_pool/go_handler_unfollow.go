package work_pool

type unFollowHandler struct{}

func initUnFollowHandler() {
	l := unFollowHandler{}
	_pool.addCmdHandler(1, l)
}

func (f unFollowHandler) Handle(d *GoWorkerData) {
	if d == nil {
		return
	}

	defer putGoWorkerData(d)

	_, ok := d.Data["follow_uid"].(uint64)
	if !ok {
		return
	}
	_ = uint64(d.ID)

}
