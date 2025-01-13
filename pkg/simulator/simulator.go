package simulator

import (
    "fmt"
    "lem-in/pkg/colony"
    "lem-in/pkg/pathfinder"
)

// Move represents a single ant movement
type Move struct {
    AntID    int
    ToRoom   string
}

// Simulator handles the simulation of ant movements
type Simulator struct {
    colony    *colony.Colony
    paths     []pathfinder.Path
    antPaths  map[int][]string  // Maps ant IDs to their assigned paths
    positions map[int]int       // Maps ant IDs to their current position in their path
}

// NewSimulator creates a new simulator instance
func NewSimulator(c *colony.Colony, paths []pathfinder.Path) *Simulator {
    return &Simulator{
        colony:    c,
        paths:     paths,
        antPaths:  make(map[int][]string),
        positions: make(map[int]int),
    }
}

// AssignPaths assigns optimal paths to each ant
func (s *Simulator) AssignPaths() {
    numAnts := s.colony.NumAnts
    currentAnt := 1
    
    // Calculate delays for each path to optimize total moves
    delays := make([]int, len(s.paths))
    for i := range s.paths {
        if i > 0 {
            delays[i] = s.paths[i].Length - s.paths[0].Length
        }
    }
    
    // Assign paths to ants, accounting for optimal delays
    for currentAnt <= numAnts {
        shortestTime := -1
        bestPath := 0
        
        // Find the path that will get this ant to the end fastest
        for i, path := range s.paths {
            antsOnPath := 0
            for _, assignedPath := range s.antPaths {
                if len(assignedPath) == len(path.Rooms) {
                    antsOnPath++
                }
            }
            
            totalTime := path.Length + delays[i] + antsOnPath
            if shortestTime == -1 || totalTime < shortestTime {
                shortestTime = totalTime
                bestPath = i
            }
        }
        
        // Convert path to room names and assign to ant
        roomNames := make([]string, len(s.paths[bestPath].Rooms))
        for i, room := range s.paths[bestPath].Rooms {
            roomNames[i] = room.Name
        }
        s.antPaths[currentAnt] = roomNames
        s.positions[currentAnt] = 0
        
        currentAnt++
    }
}

// SimulateMove simulates one turn of ant movements
func (s *Simulator) SimulateMove() ([]Move, bool) {
    var moves []Move
    roomOccupancy := make(map[string]bool)
    
    // Mark start and end rooms as always available
    roomOccupancy[s.colony.StartRoom.Name] = false
    roomOccupancy[s.colony.EndRoom.Name] = false
    
    // Try to move each ant
    for antID := 1; antID <= s.colony.NumAnts; antID++ {
        path := s.antPaths[antID]
        pos := s.positions[antID]
        
        // Skip ants that have reached the end
        if pos >= len(path)-1 {
            continue
        }
        
        // Check if next room is available
        nextRoom := path[pos+1]
        if !roomOccupancy[nextRoom] {
            // Move ant to next room
            moves = append(moves, Move{
                AntID:  antID,
                ToRoom: nextRoom,
            })
            roomOccupancy[nextRoom] = true
            s.positions[antID]++
            
            // Mark current room as unoccupied
            if pos > 0 {
                roomOccupancy[path[pos]] = false
            }
        }
    }
    
    // Check if simulation is complete
    complete := true
    for antID := 1; antID <= s.colony.NumAnts; antID++ {
        if s.positions[antID] < len(s.antPaths[antID])-1 {
            complete = false
            break
        }
    }
    
    return moves, complete
}

// FormatMoves formats the moves for output
func (s *Simulator) FormatMoves(moves []Move) string {
    if len(moves) == 0 {
        return ""
    }
    
    result := ""
    for i, move := range moves {
        if i > 0 {
            result += " "
        }
        result += fmt.Sprintf("L%d-%s", move.AntID, move.ToRoom)
    }
    return result
}