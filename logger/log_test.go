package logger

import "testing"

func TestLog(t *testing.T) {
	SetLevel(LevelTrace) //设置日志级别
	Trace("Trace测试")
	Debug("Debug测试")
	Info("Info测试")
	Warn("Warn测试")
	Error("Error测试")
	Critical("Critical测试")

	Tracef("你好%s,年龄:%d", "菜的恒", 12)
	Debugf("你好%s,年龄:%d", "菜的恒", 12)
	Infof("你好%s,年龄:%d", "菜的恒", 12)
	Warnf("你好%s,年龄:%d", "菜的恒", 12)
	Errorf("你好%s,年龄:%d", "菜的恒", 12)
	Criticalf("你好%s,年龄:%d", "菜的恒", 12)
}
