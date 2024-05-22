package strategies

import (
	"time"

	log "github.com/sirupsen/logrus"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type BaseStrategy struct {
	kConfig *rest.Config
	kClient *kubernetes.Clientset

	log              *log.Entry
	tolerateAllNodes bool
	timeout          time.Duration
}

func NewBaseStrategy(config *rest.Config, client *kubernetes.Clientset, tolerateAllNodes bool, timeout *time.Duration) BaseStrategy {
	var t time.Duration
	if timeout == nil {
		t = 60 * time.Second
	} else {
		t = *timeout
	}
	return BaseStrategy{
		kConfig:          config,
		kClient:          client,
		tolerateAllNodes: tolerateAllNodes,
		timeout:          t,
		log:              log.WithField("component", "strategy"),
	}
}

type Strategy interface {
	CompatibleWithContext(MigrationContext) error
	Description() string
	Identifier() string
	Do(sourcePVC *v1.PersistentVolumeClaim, destTemplate *v1.PersistentVolumeClaim, WaitForTempDestPVCBind bool) error
}

type MigrationContext struct {
	PVCControllers []interface{}
	SourcePVC      v1.PersistentVolumeClaim
}

func StrategyInstances(b BaseStrategy) []Strategy {
	s := []Strategy{
		NewCopyTwiceNameStrategy(b),
		NewExportStrategy(b),
		NewImportStrategy(b),
	}
	return s
}
