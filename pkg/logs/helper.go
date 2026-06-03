package logs

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Helper struct {
	h *log.Helper
}

func NewHelper(logger log.Logger) *Helper {
	return &Helper{
		h: newHelper(logger),
	}
}

func (lh *Helper) Debug(ctx context.Context, a ...any) {
	lh.h.WithContext(ctx).Debug(a...)
}
func (lh *Helper) Debugf(ctx context.Context, format string, a ...any) {
	lh.h.WithContext(ctx).Debugf(format, a...)
}
func (lh *Helper) Debugw(ctx context.Context, keyvals ...any) {
	lh.h.WithContext(ctx).Debugw(keyvals...)
}
func (lh *Helper) Info(ctx context.Context, a ...any) {
	lh.h.WithContext(ctx).Info(a...)
}
func (lh *Helper) Infof(ctx context.Context, format string, a ...any) {
	lh.h.WithContext(ctx).Infof(format, a...)
}
func (lh *Helper) Infow(ctx context.Context, keyvals ...any) {
	lh.h.WithContext(ctx).Infow(keyvals...)
}
func (lh *Helper) Warn(ctx context.Context, a ...any) {
	lh.h.WithContext(ctx).Warn(a...)
}
func (lh *Helper) Warnf(ctx context.Context, format string, a ...any) {
	lh.h.WithContext(ctx).Warnf(format, a...)
}
func (lh *Helper) Warnw(ctx context.Context, keyvals ...any) {
	lh.h.WithContext(ctx).Warnw(keyvals...)
}
func (lh *Helper) Error(ctx context.Context, a ...any) {
	lh.h.WithContext(ctx).Error(a...)
}
func (lh *Helper) Errorf(ctx context.Context, format string, a ...any) {
	lh.h.WithContext(ctx).Errorf(format, a...)
}
func (lh *Helper) Errorw(ctx context.Context, keyvals ...any) {
	lh.h.WithContext(ctx).Errorw(keyvals...)
}
func (lh *Helper) Fatal(ctx context.Context, a ...any) {
	lh.h.WithContext(ctx).Fatal(a...)
}
func (lh *Helper) Fatalf(ctx context.Context, format string, a ...any) {
	lh.h.WithContext(ctx).Fatalf(format, a...)
}
func (lh *Helper) Fatalw(ctx context.Context, keyvals ...any) {
	lh.h.WithContext(ctx).Fatalw(keyvals...)
}
