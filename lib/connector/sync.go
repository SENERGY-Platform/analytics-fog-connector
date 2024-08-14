package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/SENERGY-Platform/analytics-fog-connector/lib/logging"
)

func (connector *Connector) RequestSyncs(doneCtx context.Context, syncIntervalInSeconds time.Duration) {
	ticker := time.NewTicker(syncIntervalInSeconds * time.Second)

	for {
		select {
		case <-doneCtx.Done():
			ticker.Stop()
		case <-ticker.C:
			err := connector.RequestOperatorSync()
			if err != nil {
				logging.Logger.Error(fmt.Sprintf("Cant request operator sync %s", err.Error())) 
			}
			err = connector.RequestUpstreamSync()
			if err != nil {
				logging.Logger.Error(fmt.Sprintf("Cant request upstream sync %s", err.Error())) 
			}
		}
	}

}