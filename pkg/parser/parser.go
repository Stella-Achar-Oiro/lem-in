package parser

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    
    "lem-in/pkg/colony"
)

// Parser handles the parsing of input files
type Parser struct {
    colony *colony.Colony
}

// NewParser creates a new parser instance
func NewParser() *Parser {
    return &Parser{}
}

// Parse reads and parses the input file
func (p *Parser) Parse(filename string) (*colony.Colony, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("error opening file: %v", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    
    // Parse number of ants
    if !scanner.Scan() {
        return nil, fmt.Errorf("empty file")
    }
    
    numAnts, err := strconv.Atoi(scanner.Text())
    if err != nil || numAnts <= 0 {
        return nil, fmt.Errorf("invalid number of ants: %s", scanner.Text())
    }
    
    p.colony = colony.NewColony(numAnts)
    
    // Parse rooms and tunnels
    var parsingRooms = true
    for scanner.Scan() {
        line := scanner.Text()
        
        // Skip comments that don't indicate start/end
        if strings.HasPrefix(line, "#") {
            if line == "##start" || line == "##end" {
                if !scanner.Scan() {
                    return nil, fmt.Errorf("missing room definition after %s", line)
                }
                roomLine := scanner.Text()
                if err := p.parseRoom(roomLine, line == "##start"); err != nil {
                    return nil, err
                }
            }
            continue
        }
        
        // Empty line
        if line == "" {
            continue
        }
        
        // Check if we're switching from rooms to tunnels
        if strings.Contains(line, "-") {
            parsingRooms = false
        }
        
        if parsingRooms {
            if err := p.parseRoom(line, false); err != nil {
                return nil, err
            }
        } else {
            if err := p.parseTunnel(line); err != nil {
                return nil, err
            }
        }
    }
    
    if err := p.colony.Validate(); err != nil {
        return nil, err
    }
    
    if err := p.colony.Initialize(); err != nil {
        return nil, err
    }
    
    return p.colony, nil
}

// parseRoom parses a room definition line
func (p *Parser) parseRoom(line string, isStart bool) error {
    parts := strings.Fields(line)
    if len(parts) != 3 {
        return fmt.Errorf("invalid room format: %s", line)
    }
    
    name := parts[0]
    x, err1 := strconv.Atoi(parts[1])
    y, err2 := strconv.Atoi(parts[2])
    
    if err1 != nil || err2 != nil {
        return fmt.Errorf("invalid coordinates: %s", line)
    }
    
    if err := p.colony.AddRoom(name, x, y); err != nil {
        return err
    }
    
    if isStart {
        return p.colony.SetStart(name)
    }
    return nil
}

// parseTunnel parses a tunnel definition line
func (p *Parser) parseTunnel(line string) error {
    parts := strings.Split(line, "-")
    if len(parts) != 2 {
        return fmt.Errorf("invalid tunnel format: %s", line)
    }
    
    return p.colony.AddTunnel(parts[0], parts[1])
}