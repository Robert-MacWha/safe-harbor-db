package types

type AgreementVersion struct {
	Version SafeHarborVersion `firestore:"version"`
}
