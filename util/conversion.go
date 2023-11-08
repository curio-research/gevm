package util

import (
	// "reflect"
	logger "github.com/daweth/gevm/logger"
	vm "github.com/daweth/gevm/vm"
	gvm "github.com/ethereum/go-ethereum/core/vm"
	glogger "github.com/ethereum/go-ethereum/eth/tracers/logger"
	gparams "github.com/ethereum/go-ethereum/params"
)

func ConvertGConfigToConfig(a gvm.Config, lc glogger.Config) vm.Config {
	chainConfig := gparams.TestChainConfig
	logConfig := logger.Config{
		EnableMemory:     lc.EnableMemory,
		DisableStack:     lc.DisableStack,
		DisableStorage:   lc.DisableStack,
		EnableReturnData: lc.EnableReturnData,
		Debug:            lc.Debug,
		Limit:            lc.Limit,
		Overrides:        chainConfig,
	}
	logger := logger.NewStructLogger(&logConfig)

	return vm.Config{
		Tracer:                  logger,
		NoBaseFee:               a.NoBaseFee,
		EnablePreimageRecording: a.EnablePreimageRecording,
		ExtraEips:               a.ExtraEips,
	}

}
