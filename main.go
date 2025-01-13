package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    
    "lem-in/pkg/parser"
    "lem-in/pkg/pathfinder"
    "lem-in/pkg/simulator"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("ERROR: invalid data format, usage: go run . [filename]")
        os.Exit(1)
    }
    
    // Print input file contents first
    printFile(os.Args[1])
    fmt.Println()
    
    // Parse input file
    p := parser.NewParser()
    colony, err := p.Parse(os.Args[1])
    if err != nil {
        fmt.Printf("ERROR: invalid data format, %v\n", err)
        os.Exit(1)
    }
    
    // Find optimal paths
    pf := pathfinder.NewPathFinder(colony)
    paths := pf.FindPaths()
    if len(paths) == 0 {
        fmt.Println("ERROR: invalid data format, no valid paths found")
        os.Exit(1)
    }
    
    // Optimize paths based on number of ants
    optimalPaths := pf.OptimizePaths(paths)
    
    // Create simulator and run simulation
    sim := simulator.NewSimulator(colony, optimalPaths)
    sim.AssignPaths()
    
    // Run simulation until complete
    complete := false
    for !complete {
        moves, done := sim.SimulateMove()
        if len(moves) > 0 {
            fmt.Println(sim.FormatMoves(moves))
        }
        complete = done
    }
}

// printFile prints the contents of the input file
func printFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    reader := bufio.NewReader(file)
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                fmt.Print(line)
                break
            }
            return err
        }
        fmt.Print(line)
    }
    return nil
}