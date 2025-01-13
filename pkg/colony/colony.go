package colony

import "fmt"

// Room represents a single room in the ant colony
type Room struct {
	Name        string
	X, Y        int
	IsStart     bool
	IsEnd       bool
	Connections map[string]*Room
	CurrentAnt  *Ant // nil if no ant is present
}

// Ant represents a single ant in the colony
type Ant struct {
	ID            int
	CurrentRoom   *Room
	Path          []*Room
	HasReachedEnd bool
}

// Colony represents the entire ant colony structure
type Colony struct {
	Rooms     map[string]*Room
	StartRoom *Room
	EndRoom   *Room
	Ants      []*Ant
	NumAnts   int
}

// NewColony creates a new colony instance
func NewColony(numAnts int) *Colony {
	return &Colony{
		Rooms:   make(map[string]*Room),
		Ants:    make([]*Ant, numAnts),
		NumAnts: numAnts,
	}
}

// AddRoom adds a new room to the colony
func (c *Colony) AddRoom(name string, x, y int) error {
	if _, exists := c.Rooms[name]; exists {
		return fmt.Errorf("room %s already exists", name)
	}

	if name[0] == 'L' || name[0] == '#' {
		return fmt.Errorf("invalid room name: %s", name)
	}

	c.Rooms[name] = &Room{
		Name:        name,
		X:           x,
		Y:           y,
		Connections: make(map[string]*Room),
	}
	return nil
}

// SetStart sets the start room
func (c *Colony) SetStart(name string) error {
	room, exists := c.Rooms[name]
	if !exists {
		return fmt.Errorf("room %s does not exist", name)
	}
	room.IsStart = true
	c.StartRoom = room
	return nil
}

// SetEnd sets the end room
func (c *Colony) SetEnd(name string) error {
	room, exists := c.Rooms[name]
	if !exists {
		return fmt.Errorf("room %s does not exist", name)
	}
	room.IsEnd = true
	c.EndRoom = room
	return nil
}

// AddTunnel connects two rooms with a tunnel
func (c *Colony) AddTunnel(room1, room2 string) error {
	r1, exists1 := c.Rooms[room1]
	r2, exists2 := c.Rooms[room2]

	if !exists1 || !exists2 {
		return fmt.Errorf("one or both rooms do not exist")
	}

	if room1 == room2 {
		return fmt.Errorf("cannot connect room to itself")
	}

	if _, exists := r1.Connections[room2]; exists {
		return fmt.Errorf("tunnel already exists")
	}

	r1.Connections[room2] = r2
	r2.Connections[room1] = r1
	return nil
}

// Initialize creates all ants and places them in the start room
func (c *Colony) Initialize() error {
	if c.StartRoom == nil || c.EndRoom == nil {
		return fmt.Errorf("start or end room not set")
	}

	for i := 0; i < c.NumAnts; i++ {
		c.Ants[i] = &Ant{
			ID:          i + 1,
			CurrentRoom: c.StartRoom,
		}
	}
	return nil
}

// Validate checks if the colony structure is valid
func (c *Colony) Validate() error {
	if c.StartRoom == nil {
		return fmt.Errorf("no start room defined")
	}
	if c.EndRoom == nil {
		return fmt.Errorf("no end room defined")
	}
	if c.NumAnts <= 0 {
		return fmt.Errorf("invalid number of ants")
	}

	// Check if end is reachable from start using BFS
	visited := make(map[string]bool)
	queue := []*Room{c.StartRoom}
	visited[c.StartRoom.Name] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == c.EndRoom {
			return nil // Path found
		}

		for _, room := range current.Connections {
			if !visited[room.Name] {
				visited[room.Name] = true
				queue = append(queue, room)
			}
		}
	}

	return fmt.Errorf("no path exists from start to end")
}
