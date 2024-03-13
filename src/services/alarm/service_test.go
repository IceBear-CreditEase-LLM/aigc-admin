package alarm

import (
	"context"
	"testing"
)

func TestService_Push(t *testing.T) {
	err := New("TraceId", Config{
		Host:        "http://localhost:8080",
		Namespace:   "aigc",
		ServiceName: "aigc-admin",
	}, nil).
		Push(context.Background(),
			"xxxx-发出预警",
			"5分钟内平均请求超过1分钟",
			"aigc-admin",
			LevelInfo,
			5,
		)
	if err != nil {
		t.Error(err)
		return
	}
}
