package peerdisco

// Close 停止 gossip 循环、刷盘、关闭审计与传输。
func (c *Cluster) Close() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return ErrClosed
	}
	c.closed = true
	select {
	case <-c.stopCh:
	default:
		close(c.stopCh)
	}
	c.mu.Unlock()

	<-c.loopDone

	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.persistLocked(); err != nil {
		return err
	}
	if c.audit != nil {
		_ = c.audit.Close()
		c.audit = nil
	}
	if c.tr != nil {
		_ = c.tr.Close()
	}
	c.table.Clear()
	c.meta.Clear()
	c.susp.Clear()
	return nil
}

// StartLoop 启动后台探测循环（测试可跳过）。
func (c *Cluster) StartLoop() {
	c.mu.Lock()
	c.loopDone = make(chan struct{})
	stop := c.stopCh
	c.mu.Unlock()
	go func() {
		defer close(c.loopDone)
		for {
			select {
			case <-stop:
				return
			case <-c.sched.Tick():
				c.mu.Lock()
				if !c.closed {
					c.metrics.IncTicks()
				}
				c.mu.Unlock()
			}
		}
	}()
}
