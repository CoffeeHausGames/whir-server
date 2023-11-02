package model

import (
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Manages the base structs and methods 

// Base represents a base owned by [User]s
type Base struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Owner       primitive.ObjectID `json:"owner,omitempty" bson:"owner,omitempty"`
	Buildings 	map[string][]Building `json:"buildings,omitempty" bson:"buildings,omitempty"`
	Grid        [][]*Building       `json:"grid,omitempty" bson:"grid,omitempty"`
}

func NewBase(user_id primitive.ObjectID) *Base {
	tower := Building{
		Name:     "tower",
		IsPlaced: true,
		PosX:     0,
		PosY:     0,
		PosZ:     0,
		Width:    10, // Set the width of the tower (adjust as needed)
		Height:   10, // Set the height of the tower (adjust as needed)
	}

	// Create the grid and initialize it with nil values
	gridWidth := 100 // Set the width of the grid (adjust as needed)
	gridHeight := 100 // Set the height of the grid (adjust as needed)

	grid := make([][]*Building, gridHeight)
	for i := range grid {
		grid[i] = make([]*Building, gridWidth)
	}

	// Calculate the middle of the grid
	middleX := gridWidth / 2
	middleZ := gridHeight / 2

	// Calculate the initial position of the tower to place it in the middle
	tower.PosX = float64(middleX) - tower.Width/2
	tower.PosZ = float64(middleZ) - tower.Height/2

	base := &Base{
		ID:        primitive.NewObjectID(),
		Owner:     user_id,
		Grid: 		 grid,
		Buildings: make(map[string][]Building),
	}

	// Add the tower to the building
	err := base.AddBuildingToBase(&tower)
	if err != nil {
		log.Panic(err)
	}

	return base
}

// Validate that the building is able to be placed
func (base *Base) ValidateBuildingPlacement(building *Building) error {

	// Check if base.Grid is nil or empty (no rows)
	if base.Grid == nil || len(base.Grid) == 0  {
		return fmt.Errorf("Base grid is not instantiated correctly")
	}

	gridSizeX := len(base.Grid[0])
	gridSizeZ := len(base.Grid)

	startX := int(building.PosX)
	startZ := int(building.PosZ)
	endX := startX + int(building.Width)
	endZ := startZ + int(building.Height)

	// Ensure the building is within the size of the grid
	if startX < 0 || startZ < 0 || endX > gridSizeX || endZ > gridSizeZ {
			return fmt.Errorf("Building Placement failed: building is out of grid bounds")
	}

	// Ensure there is no other building in the placement
	for i := startZ; i < endZ; i++ {
			for j := startX; j < endX; j++ {
					if base.Grid[i][j] != nil {
							return fmt.Errorf("Building Placement failed: building overlaps with existing building")
					}
			}
	}

	return nil
}

// Validate that the building is able to be placed
func (base *Base) ValidateBuildingRemoval(building *Building) error {
	gridSizeX := len(base.Grid[0])
	gridSizeZ := len(base.Grid)

	startX := int(building.PosX)
	startZ := int(building.PosZ)
	endX := startX + int(building.Width)
	endZ := startZ + int(building.Height)

	// Ensure the building is within the size of the grid
	if startX < 0 || startZ < 0 || endX > gridSizeX || endZ > gridSizeZ {
			return fmt.Errorf("Building Placement failed: building is out of grid bounds")
	}

	// Ensure this building is at the given location
	for i := startZ; i < endZ; i++ {
			for j := startX; j < endX; j++ {
					if base.Grid[i][j] == nil || !base.Grid[i][j].Equal(building) {
							return fmt.Errorf("Building removal failed: Not same building in location")
					}
			}
	}

	return nil
}

// Add a Building to the base
func (base *Base) AddBuildingToBase(building *Building) error {
	if err := base.ValidateBuildingPlacement(building); err != nil {
		return err
	}
	building.IsPlaced = true

	err := base.addToBuildings(building)
	if err != nil {
		building.IsPlaced = false
		return err
	}

	err = base.placeBuildingOnGrid(building)
	if err != nil {
		building.IsPlaced = false
		base.removeFromBuildings(building)
		return err
	}

	return nil
}

// Remove a Building from the base
func (base *Base) RemoveBuildingFromBase(building *Building) error {
	if err := base.ValidateBuildingRemoval(building); err != nil {
		return err
	}

	err, removedBuilding := base.removeFromBuildings(building)
	if err != nil {
		return err
	}

	err, removedBuilding  = base.removeBuildingFromGrid(building)
	if err != nil {
		base.addToBuildings(building)
		return err
	}
	removedBuilding.IsPlaced = false

	return nil
}

// Add a building to the grid of a base
func (base *Base) placeBuildingOnGrid(building *Building) error {
	startX := int(building.PosX)
	startZ := int(building.PosZ)
	endX := startX + int(building.Width)
	endZ := startZ + int(building.Height)

	for i := startZ; i < endZ; i++ {
			for j := startX; j < endX; j++ {
					base.Grid[i][j] = building
			}
	}

	return nil
}

// Remove a building from the grid of a base
func (base *Base) removeBuildingFromGrid(buildingToRemove *Building) (error, *Building) {
	startX := int(buildingToRemove.PosX)
	startZ := int(buildingToRemove.PosZ)
	endX := startX + int(buildingToRemove.Width)
	endZ := startZ + int(buildingToRemove.Height)

	emptyBuilding := new(Building)
	removedBuilding := new(Building)

	for i := startZ; i < endZ; i++ {
		for j := startX; j < endX; j++ {			
			// Check if the current cell contains the building to remove
				// Set the cell to an empty building
				fmt.Println("here")
				fmt.Println(base.Grid[i][j])
				removedBuilding = base.Grid[i][j]
				base.Grid[i][j] = emptyBuilding
		}
	}

	return nil, removedBuilding
}

// Add a new building to the Buildings map
func (base *Base) addToBuildings(newBuilding *Building) error {
	// Check if the map already has an entry for the building type
	buildingType := strings.ToLower(newBuilding.Name)
	
	if existingBuildings, ok := base.Buildings[buildingType]; ok {
			// If the building type exists, append the new building to the existing slice
			base.Buildings[buildingType] = append(existingBuildings, *newBuilding)
	} else {
			// If the building type doesn't exist, create a new entry in the map
			base.Buildings[buildingType] = []Building{*newBuilding}
	}
	return nil
}

// Remove a building from the Buildings map
func (base *Base) removeFromBuildings(buildingToRemove *Building) (error, *Building) {
	// Check if the map already has an entry for the building type
	buildingType := strings.ToLower(buildingToRemove.Name)
	if existingBuildings, ok := base.Buildings[buildingType]; ok {
		for i, existingBuilding := range existingBuildings {
			// Check if the existing building is the one to be removed
			if existingBuilding.Equal(buildingToRemove) {
				// Remove the building by slicing it out of the slice
				base.Buildings[buildingType] = append(existingBuildings[:i], existingBuildings[i+1:]...)
				return nil, &existingBuildings[i]
			}
		}
	}
	// Building type not found in the map, return an error
	return fmt.Errorf("Building type '%s' not found in the map", buildingType), nil
}
