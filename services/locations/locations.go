package locations

import (
	"github.com/dhconnelly/rtreego"
)

// UserLocations user locations
type UserLocations struct {
	locations *rtreego.Rtree
}

//Build Builds new R-tree to store closest neighbors
func Build(dimensions, minBranches, maxBranches int, userLocs ...rtreego.Spatial) *UserLocations {
	return &UserLocations{
		locations: rtreego.NewTree(dimensions, minBranches, maxBranches, userLocs...),
	}
}

//Save saves R-tree in some form of persistent storage
func (locationTree *UserLocations) Save() *UserLocations {
	return locationTree
}

//Rebuild rebuilds the R-tree and saves to some form of persistent storage if there exists some inconsistencies
func (locationTree *UserLocations) Rebuild() *UserLocations {
	return locationTree
}

//AddLocation adds to the the R-tree and saves addition in some form of persistent storage
func (locationTree *UserLocations) AddLocation() *UserLocations {
	return locationTree
}

//DeleteLocation deletes from the the R-tree and saves deletion in some form of persistent storage
func (locationTree *UserLocations) DeleteLocation() *UserLocations {
	return locationTree
}

func (locationTree *UserLocations) GetNNearestNeighbors() {

}

func conversion(addressLiteral string) {
	// TODO: Implement conversion from actual address to lat/long (potentially) for standard location storage
}
