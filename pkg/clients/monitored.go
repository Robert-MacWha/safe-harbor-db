package clients

// Monitored represents a protocol that is being monitored
type Monitored struct {
	Addresses []AccountM
	MockData  bool
}

// AccountM represents information about the EOA or contract or child contracts to include.
type AccountM struct {
	Name     string
	Address  string
	Chains   []int
	Children []AccountM
}
