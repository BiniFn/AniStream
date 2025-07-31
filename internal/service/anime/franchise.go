package anime

import "github.com/coeeter/aniways/internal/infra/client/shikimori"

func deriveWatchOrder(fr *shikimori.FranchiseResponse, malID int) []int {
	forward := make(map[int]int, len(fr.Links))
	backward := make(map[int]int, len(fr.Links))
	for _, link := range fr.Links {
		if link.Relation == "sequel" {
			forward[link.SourceID] = link.TargetID
			backward[link.TargetID] = link.SourceID
		}
	}

	first := malID
	for {
		prev, ok := backward[first]
		if !ok {
			break
		}
		first = prev
	}

	var watchOrder []int
	for cur := first; ; cur = forward[cur] {
		watchOrder = append(watchOrder, cur)
		_, ok := forward[cur]
		if !ok {
			break
		}
	}

	return watchOrder
}

func deriveRelated(fr *shikimori.FranchiseResponse, malID int, watchOrder []int) []int {
	inWatchOrder := make(map[int]bool, len(watchOrder))
	for _, id := range watchOrder {
		inWatchOrder[id] = true
	}

	var related []int
	for _, node := range fr.Nodes {
		if node.ID != malID && !inWatchOrder[node.ID] {
			related = append(related, node.ID)
		}
	}

	return related
}

func deriveFullFranchise(fr *shikimori.FranchiseResponse) []int {
	ids := make([]int, 0, len(fr.Nodes))
	for _, node := range fr.Nodes {
		ids = append(ids, node.ID)
	}
	return ids
}
