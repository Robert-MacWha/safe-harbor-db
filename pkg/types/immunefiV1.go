package types

type ImmunefiDetailsV1 struct {
	Name        string            `firestore:"name"`
	Contact     string            `firestore:"contact"`
	Chains      []ImmunefiChainV1 `firestore:"chains"`
	BountyTerms BountyTermsV1     `firestore:"bountyTerms"`
}

type ImmunefiChainV1 struct {
	ID       int                 `firestore:"id"`
	Accounts []ImmunefiAccountV1 `firestore:"accounts"`
}

type ImmunefiAccountV1 struct {
	Name    string `firestore:"name"`
	Address string `firestore:"address"`
}
