package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var GenAddr = map[string]string{
	// T1
	"g1us8428u2a5satrlxzagqqa5m6vmuze025anjlj": "Jae",
	"g1manfred47kzduec920z88wfr64ylksmdcedlf5": "Manfred",
	"g1e6gxg5tvc55mwsn7t7dymmlasratv7mkv0rap2": "Milos",
	"g1qhskthp2uycmg4zsdc9squ2jds7yv3t0qyrlnp": "Petar",
	"g18amm3fc00t43dcxsys6udug0czyvqt9e7p23rd": "Marc",
	"g19p3yzr3cuhzqa02j0ce6kzvyjqfzwemw3vam0x": "Guilhem",
	"g1mx4pum9976th863jgry4sdjzfwu03qan5w2v9j": "Ray",
	"g127l4gkhk0emwsx5tmxe96sp86c05h8vg5tufzq": "Maxwell",
	"g1m0rgan0rla00ygmdmp55f5m0unvsvknluyg2a4": "howl999", // Morgan
	"g1ker4vvggvsyatexxn3hkthp2hu80pkhrwmuczr": "Sergio",
	"g1aeddlftlfk27ret5rf750d7w5dume3kcsm8r8m": "aeddiold4242", // Antoine
	"g16tfrrul20g4jzt3z303raqw8vs8s2pqqh5clwu": "Ilker",
	"g1hy6zry03hg5d8le9s2w4fxme6236hkgd928dun": "jeronimo000", // Jeronimo
	"g15ruzptpql4dpuyzej0wkt5rq6r26kw4nxu9fwd": "Denis",
	"g1lckl8j2g3jyyuq6fx7pke3uz4kemht7lw4fg5l": "Danny",
	"g1778y2yphxs2wpuaflsy5y9qwcd4gttn4g5yjx5": "Michelle",
	"g1mq7g0jszdmn4qdpc9tq94w0gyex37su892n80m": "Alan",
	"g197q5e9v00vuz256ly7fq7v3ekaun5cr7wmjgfh": "Salvo",
	"g1mpkp5lm8lwpm0pym4388836d009zfe4maxlqsq": "Alexis",
	"g125em6arxsnj49vx35f0n0z34putv5ty3376fg5": "leon000",
	"g1whzkakk4hzjkvy60d5pwfk484xu67ar2cl62h2": "Kirk",
	"g1sw5xklxjjuv0yvuxy5f5s3l3mnj0nqq626a9wr": "Albert",

	// T2
	"g1jazghxvvgz3egnr2fc8uf72z4g0l03596y9ls7": "kouteki070", // Nemanja
	"g1dfr24yhk5ztwtqn2a36m8f6ud8cx5hww4dkjfl": "Antonio",
	"g12vx7dn3dqq89mz550zwunvg4qw6epq73d9csay": "onbloc", // Dongwon
	"g1r04aw56fgvzy859fachr8hzzhqkulkaemltr76": "Blake",
	"g17n4y745s08awwq4e0a38lagsgtntna0749tnxe": "Jinwoo",
	"g1ckae7tc5sez8ul3ssne75sk4muwgttp6ks2ky9": "ByeongJun",

	// T3
	"g14u5eaheavy0ux4dmpykg2gvxpvqvexm9cyg58a": "Norman",
	"g1qynsu9dwj9lq0m5fkje7jh6qy3md80ztqnshhm": "Rémi",
	"g17ernafy6ctpcz6uepfsq2js8x2vz0wladh5yc3": "Dragos",
}

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
			fmt.Printf("Voter: %-20s - %s Proposals voted: %d\n", entry.Name, addr, entry.Count)
		} else {
			fmt.Printf("Voter: %-20s - %s Proposals voted: %d\n", entry.Name, "no username", entry.Count)
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
			fmt.Printf("Missing address: %s\n", voter)
		}
	}

	fmt.Println("\nAddresses that didn't vote:")
	fmt.Println("---------------------------")
	for addr, name := range GenAddr {
		if Voters[name] == 0 && Voters[addr] == 0 {
			fmt.Printf("Non-voter: %s (%s)\n", name, addr)
		}
	}
}

func main() {
	fmt.Println("Gno Governance Proposal Voters Analysis")
	fmt.Println("======================================")

	getProposals(6)
	displayProposal()
}
