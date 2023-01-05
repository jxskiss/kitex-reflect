package kitexreflect

import (
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
)

type Option struct {
	applyDescProv func(impl *ProviderImpl)
}

const defaultDebounceInterval = 10 * time.Second

func defaultErrorHandler(err error, msg string) {
	klog.Errorf("KitexReflect: %s: %v", msg, err)
}

// WithDebounceInterval sets max interval to debounce requests
// sent to backend service. The default is 10 seconds.
func WithDebounceInterval(d time.Duration) Option {
	return Option{
		applyDescProv: func(impl *ProviderImpl) {
			impl.debounceInterval = d
		},
	}
}

// WithErrorHandler optionally specifies an error handler function
// for errors happened when updating service descriptor from
// backend service.
// By default, it logs error using the [klog] package.
func WithErrorHandler(f func(err error, msg string)) Option {
	return Option{
		applyDescProv: func(impl *ProviderImpl) {
			impl.errorHandler = f
		},
	}
}

// InitAsync tells the description provider to initialize service
// descriptor asynchronously, don't wait for the first descriptor
// being ready.
func InitAsync() Option {
	return Option{
		applyDescProv: func(impl *ProviderImpl) {
			impl.initAsync = true
		},
	}
}
