/*
 * Copyright (c) 2022 Avi Misra
 *
 * Use of this work is governed by a MIT License.
 * You may find a license copy in project root.
 */

package etherscan

import "fmt"

// GetLogs gets logs that match "topicsOrOps" filter emitted by the specified "address" between the "fromBlock" and "toBlock"
// Leave nil if needed. Atleast on topic is a must. "topicOrOps" is a sequence of topics of 32 bytes and operators either `and`
// or `or`.
func (c *Client) GetLogs(fromBlock, toBlock *int, address *string, topicsOrOps ...string) (logs []Log, err error) {
	if len(topicsOrOps)%2 != 1 {
		return nil, fmt.Errorf("atleast one topic is required and each topic must have an operator in between")
	}

	if len(topicsOrOps) > 7 {
		return nil, fmt.Errorf("cannot provide more than 4 topics")
	}

	param := M{
		"fromBlock": fromBlock,
		"toBlock":   toBlock,
		"address":   address,
	}

	// Not all the operators are supported (https://docs.etherscan.io/api-endpoints/logs).
	// Only operator in between topics are supported.
	for index, topicOrOp := range topicsOrOps {
		if index%2 == 0 {
			param[fmt.Sprintf("topic%v", index)] = topicOrOp
		} else {
			if topicOrOp != "and" && topicOrOp != "or" {
				return nil, fmt.Errorf("invalid operator")
			}

			param[fmt.Sprintf("topic%v_%v_opr", index-1, index)] = topicOrOp
		}
	}

	err = c.call("logs", "getLogs", param, &logs)
	return
}
