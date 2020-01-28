package loc

import (
	"github.com/dhconnelly/rtreego"
)

//Build Builds new R-tree to store closest neighbors
func Build(dimensions, minBranches, maxBranches int, userLocs ...rtreego.Spatial) *rtreego.Rtree {
	locationTree := rtreego.NewTree(dimensions, minBranches, maxBranches, userLocs...)
	return locationTree
}

//Save saves R-tree in persistent storage
func Save(locationTree *rtreego.Rtree) *rtreego.Rtree {
	return locationTree
}
