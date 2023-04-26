package params

import (
	"time"
)

const BlockCapacity = 1048576

const TimeToWaitNextBlock = 10 * time.Second

const InitialBlockTimestamp = "2023-02-25 00:00:00"
const InitialBlockPrevHash = "0000000000000000000000000000000000000000000000000000000000000000"

const ConnectionBreaksValidPeriod = 30 * 24 * time.Hour
const TrafficValidPeriod = 10 * 24 * time.Hour
