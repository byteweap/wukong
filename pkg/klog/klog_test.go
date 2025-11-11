package klog_test

import (
	"testing"

	"github.com/byteweap/wukong/pkg/klog"
)

func TestDefaultLog_Debug(t *testing.T) {
	log := klog.New()
	log.Debug().Int("int", 1).Str("k", "INT").Msg("debug")
	log.Info().Int("int", 1).Str("k", "INT").Msg("info")
	log.Warn().Int("int", 1).Str("k", "INT").Msg("warn")
	log.Error().Int("int", 1).Str("k", "INT").Msg("error")
	log.Fatal().Int("int", 1).Str("k", "INT").Msg("fatal")
	log.Panic().Int("int", 1).Str("k", "INT").Msg("panic")
	log.Info().Int("int", 1).Str("k", "INT").Msg("info...")
}
