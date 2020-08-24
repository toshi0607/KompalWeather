package http

//https://cloud.google.com/run/docs/triggering/using-scheduler Cloud Schedulerから呼び出すエンドポイント
//Cloud Scheduler で使用しているサービスをデプロイするときは、未承認の呼び出しを許可しないでください。

type server struct {
}

func (s server) Run() error {
	//ctx := context.Background()

	// init secret manager client
	// init config
	// init real clients
	// init controller
	// signal handling
	// これをCloud Schedulerでたたくパスのハンドラーに渡してからRun
	//c := controller.New()
	//err := c.Run(ctx)
	return nil
}
