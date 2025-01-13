package pathfinder

import (
    "container/heap"
    "lem-in/pkg/colony"
)

// PathFinder handles finding optimal paths through the colony
type PathFinder struct {
    colony *colony.Colony
}

// Path represents a sequence of rooms from start to end
type Path struct {
    Rooms []*colony.Room
    Length int
}

// NewPathFinder creates a new pathfinder instance
func NewPathFinder(c *colony.Colony) *PathFinder {
    return &PathFinder{colony: c}
}

// item is a queue item for Dijkstra's algorithm
type item struct {
    room     *colony.Room
    priority int
    index    int
    path     []*colony.Room
}

// priorityQueue implements heap.Interface
type priorityQueue []*item

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq priorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}
func (pq *priorityQueue) Push(x interface{}) {
    n := len(*pq)
    item := x.(*item)
    item.index = n
    *pq = append(*pq, item)
}
func (pq *priorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    old[n-1] = nil
    item.index = -1
    *pq = old[0 : n-1]
    return item
}

// FindPaths finds multiple possible paths from start to end
func (pf *PathFinder) FindPaths() []Path {
    var paths []Path
    visited := make(map[string]bool)
    
    // Initialize priority queue with start room
    pq := make(priorityQueue, 0)
    heap.Init(&pq)
    heap.Push(&pq, &item{
        room:     pf.colony.StartRoom,
        priority: 0,
        path:     []*colony.Room{pf.colony.StartRoom},
    })
    
    for pq.Len() > 0 {
        current := heap.Pop(&pq).(*item)
        room := current.room
        
        // Skip if we've been here on a shorter path
        if visited[room.Name] {
            continue
        }
        
        // If we reached the end room, add the path to our collection
        if room == pf.colony.EndRoom {
            paths = append(paths, Path{
                Rooms:  current.path,
                Length: len(current.path) - 1,
            })
            // Don't mark end room as visited to allow finding multiple paths
            continue
        }
        
        visited[room.Name] = true
        
        // Add all unvisited neighbors to the queue
        for _, nextRoom := range room.Connections {
            if !visited[nextRoom.Name] {
                newPath := make([]*colony.Room, len(current.path))
                copy(newPath, current.path)
                heap.Push(&pq, &item{
                    room:     nextRoom,
                    priority: current.priority + 1,
                    path:     append(newPath, nextRoom),
                })
            }
        }
        
        visited[room.Name] = false // Allow room to be used in other paths
    }
    
    return paths
}

// OptimizePaths selects the best combination of paths for the number of ants
func (pf *PathFinder) OptimizePaths(paths []Path) []Path {
    if len(paths) == 0 {
        return nil
    }
    
    // Sort paths by length
    for i := 0; i < len(paths)-1; i++ {
        for j := i + 1; j < len(paths); j++ {
            if paths[i].Length > paths[j].Length {
                paths[i], paths[j] = paths[j], paths[i]
            }
        }
    }
    
    // Calculate optimal number of paths to use based on number of ants
    numAnts := pf.colony.NumAnts
    optimalPaths := make([]Path, 0)
    
    // Always include the shortest path
    optimalPaths = append(optimalPaths, paths[0])
    
    // Add additional paths if they help reduce the total number of turns
    for i := 1; i < len(paths); i++ {
        currentTurns := calculateTotalTurns(optimalPaths, numAnts)
        pathsCopy := make([]Path, len(optimalPaths))
        copy(pathsCopy, optimalPaths)
        pathsCopy = append(pathsCopy, paths[i])
        newTurns := calculateTotalTurns(pathsCopy, numAnts)
        
        if newTurns < currentTurns {
            optimalPaths = append(optimalPaths, paths[i])
        }
    }
    
    return optimalPaths
}

// calculateTotalTurns calculates the total number of turns needed to move all ants
func calculateTotalTurns(paths []Path, numAnts int) int {
    if len(paths) == 0 {
        return 0
    }
    
    // Calculate how many ants should use each path
    antsPerPath := make([]int, len(paths))
    remainingAnts := numAnts
    
    for remainingAnts > 0 {
        shortestTime := -1
        shortestPath := 0
        
        for i, path := range paths {
            timeWithExtraAnt := path.Length + antsPerPath[i]
            if shortestTime == -1 || timeWithExtraAnt < shortestTime {
                shortestTime = timeWithExtraAnt
                shortestPath = i
            }
        }
        
        antsPerPath[shortestPath]++
        remainingAnts--
    }
    
    // Calculate total turns needed
    maxTurns := 0
    for i, path := range paths {
        turns := path.Length + antsPerPath[i] - 1
        if turns > maxTurns {
            maxTurns = turns
        }
    }
    
    return maxTurns
}