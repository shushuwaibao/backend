package model

type Node struct {
	NodeID  int    `gorm:"primaryKey" json:"nodeId"`
	CPU     int    `gorm:"type:int" json:"cpu"`
	GPU     int    `gorm:"type:int" json:"gpu"`
	Memory  int    `gorm:"type:int" json:"memory"`
	Storage int    `gorm:"type:int" json:"storage"`
	Label   string `gorm:"size:255" json:"label"`
}

type NodeRemainingResources struct {
	NodeID        int `json:"nodeId"`
	RemainingCPU  int `json:"remainingCPU"`
	RemainingGPU  int `json:"remainingGPU"`
	RemainingMem  int `json:"remainingMem"`
	RemainingStor int `json:"remainingStor"`
}

func GetRemainingResources(nodeid int) ([]NodeRemainingResources, error) {
	return nil, nil
}
