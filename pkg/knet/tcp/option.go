package tcp

type Options struct {
	Addr string
}

type Option func(*Options)

func defaultOptions() *Options {
	return &Options{
		Addr: "0.0.0.0:8000",
	}
}
