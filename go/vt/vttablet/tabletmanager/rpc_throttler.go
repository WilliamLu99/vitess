package tabletmanager

import (
	"context"
	"net/http"

	tabletmanagerdatapb "vitess.io/vitess/go/vt/proto/tabletmanagerdata"
	"vitess.io/vitess/go/vt/vttablet/tabletserver/throttle"
)

func (tm *TabletManager) ThrottlerCheck(ctx context.Context, request *tabletmanagerdatapb.ThrottlerCheckRequest) (*tabletmanagerdatapb.ThrottlerCheckResponse, error) {
	appName := request.App
	if appName == "" {
		appName = throttle.DefaultAppName
	}
	remoteAddr := "TODO"
	flags := &throttle.CheckFlags{
		// TODO: implement these flags
		// LowPriority:           (r.URL.Query().Get("p") == "low"),
		// SkipRequestHeartbeats: (r.URL.Query().Get("s") == "true"),
		LowPriority:           false,
		SkipRequestHeartbeats: false,

		// There's even more flags, but they're not used currently in the HTTP handler
		// ReadCheck             bool
		// OverrideThreshold     float64
		// LowPriority           bool
		// OKIfNotExists         bool
		// SkipRequestHeartbeats bool
	}
	checkResult := tm.ThrottlerService.CheckByType(ctx, appName, remoteAddr, flags, throttle.ThrottleCheckPrimaryWrite)
	if checkResult.StatusCode == http.StatusNotFound && flags.OKIfNotExists {
		checkResult.StatusCode = http.StatusOK // 200
	}

	// type CheckResult struct {
	// 	StatusCode int     `json:"StatusCode"`
	// 	Value      float64 `json:"Value"`
	// 	Threshold  float64 `json:"Threshold"`
	// 	Error      error   `json:"-"`
	// 	Message    string  `json:"Message"`
	// }

	return &tabletmanagerdatapb.ThrottlerCheckResponse{
		Message:   checkResult.Message,
		Value:     float32(checkResult.Value),
		Threshold: float32(checkResult.Threshold),
		// TODO: add StatusCode and Error
		// StatusCode: int32(checkResult.StatusCode),
		// Error: 	 checkResult.Error,
	}, nil
}
