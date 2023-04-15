package tasks

import (
	chain_task "bvpn-prototype/internal/chain/task"
	peer_task "bvpn-prototype/internal/peer/task"
)

func Init() {
	peer_task.Run()
	chain_task.Run()
}
