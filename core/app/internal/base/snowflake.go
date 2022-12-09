package base

import (
	"github.com/ichaly/go-api/core/app/pkg/util"
	"github.com/sony/sonyflake"
	"strings"
	"time"
)

var sf *sonyflake.Sonyflake

func init() {
	st := sonyflake.Settings{
		// machineID是个回调函数
		MachineID: getMachineID,
		StartTime: time.Date(2022, 12, 9, 0, 0, 0, 0, time.Local),
	}
	sf = sonyflake.NewSonyflake(st)
}

// 模拟获取本机的机器ID
func getMachineID() (mID uint16, err error) {
	result := strings.Join(util.GetMAC(), ",")
	mID = uint16(util.HashCode(result) % 10)
	return
}

func GenerateID() (uint64, error) {
	return sf.NextID()
}
