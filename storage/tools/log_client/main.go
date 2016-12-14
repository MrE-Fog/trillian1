package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/google/trillian"
	"github.com/google/trillian/storage/tools"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var treeIDFlag = flag.Int64("treeid", 3, "The tree id to use")
var startLeafFlag = flag.Int64("start_leaf", 0, "The first leaf index to fetch")
var numLeavesFlag = flag.Int64("num_leaves", 1, "The number of leaves to fetch")

func buildGetLeavesByIndexRequest(logID int64, startLeaf, numLeaves int64) *trillian.GetLeavesByIndexRequest {
	if startLeaf < 0 || numLeaves <= 0 {
		panic("Start leaf index and num_leaves must be >= 0")
	}

	var leafIndices []int64

	for l := int64(0); l < numLeaves; l++ {
		leafIndices = append(leafIndices, l+startLeaf)
	}

	return &trillian.GetLeavesByIndexRequest{LogId: logID, LeafIndex: leafIndices}
}

// TODO: Move this code out to a better place when we tidy up the initial test main stuff
// It's just a basic skeleton at the moment.
func main() {
	flag.Parse()

	port := tools.GetLogServerPort()

	// TODO: Other options apart from insecure connections
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure(), grpc.WithTimeout(time.Second*5))

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := trillian.NewTrillianLogClient(conn)

	req := buildGetLeavesByIndexRequest(*treeIDFlag, *startLeafFlag, *numLeavesFlag)
	getLeafByIndexResponse, err := client.GetLeavesByIndex(context.Background(), req)

	if err != nil {
		fmt.Printf("Got error in call: %v", err)
	} else {
		fmt.Printf("Got server response: %v", getLeafByIndexResponse)
	}
}
