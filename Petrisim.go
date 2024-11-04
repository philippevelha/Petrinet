package main

// this program reads a petri network saved in json format
import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Place struct {
	ID     string `json:"id"`
	Tokens int    `json:"tokens"`
}

type Transition struct {
	ID        string `json:"id"`
	Delay     int    `json:"delay"` // delay in milliseconds
	LastFired time.Time
	Inputs    []string `json:"-"`
	Outputs   []string `json:"-"`
}

type Arc struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type PetriNet struct {
	Places      map[string]*Place
	Transitions map[string]*Transition
}

// Load the Petri net from a JSON file
func LoadPetriNet(filename string) (*PetriNet, error) {
	// Read the JSON file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var jsonData struct {
		Places      []Place      `json:"places"`
		Transitions []Transition `json:"transitions"`
		Arcs        []Arc        `json:"arcs"`
	}
	if err := json.NewDecoder(file).Decode(&jsonData); err != nil {
		return nil, err
	}

	// Initialize the Petri net structure
	petriNet := &PetriNet{
		Places:      make(map[string]*Place),
		Transitions: make(map[string]*Transition),
	}

	// Load places
	for _, place := range jsonData.Places {
		p := place // avoid referencing loop variable
		petriNet.Places[p.ID] = &p
	}

	// Load transitions
	for _, transition := range jsonData.Transitions {
		t := transition // avoid referencing loop variable
		t.LastFired = time.Now()
		petriNet.Transitions[t.ID] = &t
	}

	// Load arcs
	for _, arc := range jsonData.Arcs {
		if place, ok := petriNet.Places[arc.From]; ok {
			transition := petriNet.Transitions[arc.To]
			transition.Inputs = append(transition.Inputs, place.ID)
		} else if place, ok := petriNet.Places[arc.To]; ok {
			transition := petriNet.Transitions[arc.From]
			transition.Outputs = append(transition.Outputs, place.ID)
		}
	}

	return petriNet, nil
}

// Check if a transition is enabled
func (pn *PetriNet) IsTransitionEnabled(transition *Transition) bool {
	for _, placeID := range transition.Inputs {
		if pn.Places[placeID].Tokens == 0 {
			return false
		}
	}
	return true
}

// Fire a transition if it is enabled and satisfies timing constraints
func (pn *PetriNet) FireTransition(transition *Transition) bool {
	now := time.Now()
	if pn.IsTransitionEnabled(transition) && int(now.Sub(transition.LastFired).Milliseconds()) >= transition.Delay {
		fmt.Printf("Firing transition %s\n", transition.ID)
		for _, placeID := range transition.Inputs {
			pn.Places[placeID].Tokens--
		}
		for _, placeID := range transition.Outputs {
			pn.Places[placeID].Tokens++
		}
		transition.LastFired = now
		return true
	}
	return false
}

// Run the simulation and report timing violations
func (pn *PetriNet) Simulate() {
	for {
		anyFired := false
		for _, transition := range pn.Transitions {
			if pn.FireTransition(transition) {
				anyFired = true
			}
		}

		// Print current state of the places
		fmt.Println("Current State of Places:")
		for _, place := range pn.Places {
			fmt.Printf("Place %s: %d tokens\n", place.ID, place.Tokens)
		}

		if !anyFired {
			break // stop if no transitions can fire
		}
		//time.Sleep(500 * time.Millisecond) // slow down the simulation
	}
	fmt.Println("Simulation completed.")
	time.Sleep(10000 * time.Millisecond)
}

func main() {
	filename := "petri_net.json"
	petriNet, err := LoadPetriNet(filename)
	if err != nil {
		fmt.Println("Error loading Petri net:", err)
		return
	}

	fmt.Println("Starting Petri Net Simulation")
	petriNet.Simulate()
}
