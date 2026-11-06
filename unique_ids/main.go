package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
	"unique-ids/internal/snowflake"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	gen := snowflake.Generator{}
	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "generate_ok"
		nodeId, err := strconv.Atoi(n.ID()[1:])
		if err != nil {
			return err
		}

		id, err := gen.NewID(time.Now(), int64(nodeId))
		if err != nil {
			return err
		}
		body["id"] = id.String()
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
