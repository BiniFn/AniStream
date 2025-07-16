package shikimori

type FranchiseResponse struct {
	Links     []Link `json:"links"`
	Nodes     []Node `json:"nodes"`
	CurrentID int    `json:"current_id"`
}

type Link struct {
	ID       int    `json:"id"`
	SourceID int    `json:"source_id"`
	TargetID int    `json:"target_id"`
	Relation string `json:"relation"`
}

type Node struct {
	ID     int `json:"id"`
	Weight int `json:"weight"`
}
