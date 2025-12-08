package types

type CantinaDetailsV1 struct {
	Name              string                     `firestore:"name"`
	Contact           string                     `firestore:"contact"`
	RecoveryAddresses []CantinaRecoveryAddressV1 `firestore:"recoveryAddresses"`
	Assets            []CantinaAssetsV1          `firestore:"assets"`
	BountyTerms       BountyTermsV1              `firestore:"bountyTerms"`
}

type CantinaRecoveryAddressV1 struct {
	Address string `firestore:"address"`
	Chain   string `firestore:"chain"`
}

type CantinaAssetsV1 struct {
	Name        string `firestore:"name"`
	Description string `firestore:"description"`
}
