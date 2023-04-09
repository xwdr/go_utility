package utils

// OptionFunc 范式方法定义
type OptionFunc func() error

// OptionFuncExec 执行多个函数方法
func OptionFuncExec(opts ...OptionFunc) error {
	for _, opt := range opts {
		if err := opt(); err != nil {
			return err
		}
	}
	return nil
}

// AsyncOptionFuncExec 异步执行多个函数方法提升执行时间
func AsyncOptionFuncExec(opts ...OptionFunc) error {
	ch := make(chan error, len(opts))
	defer close(ch)
	var reErr error
	for _, opt := range opts {
		//异步执行调用链函数
		go func(opt OptionFunc) {
			// 执行调用，并监听上下文执行状态
			ch <- opt()
		}(opt)
	}

	// 检查函数执行返回错误
	for i := len(opts); i > 0; i-- {
		// 监听错误
		err := <-ch
		if err != nil && reErr == nil {
			reErr = err
		}
	}
	return reErr
}
