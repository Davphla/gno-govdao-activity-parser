package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

///// Parameters /////

// taken from `genesis_addr.go`
var GenAddr = GenAddrTest7

// Number of proposal
var TotalProposal = 7

//////////////////////

var Voters = map[string]int{}

type VoterEntry struct {
	Name  string
	Count int
}

func getProposal(id int) {
	filename := fmt.Sprintf("data/proposal_%d", id)

	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading %s: %v", filename, err)
		return
	}

	// Parse JSON array
	var voters []string
	err = json.Unmarshal(data, &voters)
	if err != nil {
		log.Printf("Error parsing JSON from %s: %v", filename, err)
		return
	}

	for _, voter := range voters {
		if strings.HasPrefix(voter, "g1") {
			// address
			name := GenAddr[voter]
			if name == "" {
				name = voter
			}
			Voters[name] += 1
		} else {
			// username
			Voters[voter] += 1
		}
	}
}

func getProposals(maxID int) {
	for i := range maxID {
		getProposal(i)
	}
}

// Not ideal implementation, but it's working
func displayVoters(voterList []VoterEntry) {
	for _, entry := range voterList {
		// Find the address for this voter
		var addr string
		for genAddr, genName := range GenAddr {
			if genName == entry.Name {
				addr = genAddr
				break
			}
		}

		if addr != "" {
			fmt.Printf("%-20s (%s) - Proposals voted: %d\n", entry.Name, addr, entry.Count)
		} else {
			fmt.Printf("%-20s (%s) - Proposals voted: %d\n", entry.Name, "no username", entry.Count)
		}
	}

}

func displayProposal() {
	fmt.Println("\nVoting Statistics:")
	fmt.Println("------------------")

	var voterList []VoterEntry
	for voter, nb := range Voters {
		voterList = append(voterList, VoterEntry{Name: voter, Count: nb})
	}

	// Sort by vote count (descending), then by name (ascending)
	sort.Slice(voterList, func(i, j int) bool {
		if voterList[i].Count != voterList[j].Count {
			return voterList[i].Count > voterList[j].Count
		}
		return voterList[i].Name < voterList[j].Name
	})

	displayVoters(voterList)

	fmt.Println("\nMissing addresses (not in GenAddr map):")
	fmt.Println("---------------------------------------")
	for voter := range Voters {
		if strings.HasPrefix(voter, "g1") {
			// This is an address that wasn't found in GenAddr
			fmt.Printf("%s\n", voter)
		}
	}

	fmt.Println("\nAddresses that didn't vote:")
	fmt.Println("---------------------------")
	for addr, name := range GenAddr {
		if Voters[name] == 0 && Voters[addr] == 0 {
			fmt.Printf("%s (%s)\n", name, addr)
		}
	}
}

func main() {
	fmt.Println("Gno Governance Proposal Voters Analysis")
	fmt.Println("======================================")

	getProposals(TotalProposal)
	displayProposal()
}
