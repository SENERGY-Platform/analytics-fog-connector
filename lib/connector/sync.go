package connector

import (
	"context"
	"time"
)

func (connector *Connector) RequestSyncs(doneCtx context.Context, syncIntervalInSeconds time.Duration) error {
	ticker := time.NewTicker(syncIntervalInSeconds * time.Second)

	for {
		select {
		case <-doneCtx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			err := connector.RequestOperatorSync()
			if err != nil {
				return err 
			}
			err = connector.RequestUpstreamSync()
			return err
		}
	}

}