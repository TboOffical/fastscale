package pipes

import "orchestrator"

func main() {
	node := orchestrator.PipesNode{
		Online:      true,
		UID:         "1234",
		Ip:          nil,
		Capacity:    100,
		Version:     "1.0",
		LastCheckin: nil,
	}
}