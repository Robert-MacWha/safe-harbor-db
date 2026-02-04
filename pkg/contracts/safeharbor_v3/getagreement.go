package safeharbor_v3

import (
	"SHDB/pkg/contracts/adoptiondetails"
	"context"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetAgreement(
	txhash common.Hash,
	eClient *ethclient.Client,
) (agreementAddress *common.Address, agreement *adoptiondetails.AgreementDetailsV1, err error) {
	contract, err := NewSafeHarborV3Filterer(common.HexToAddress("0x0"), eClient)
	if err != nil {
		return nil, nil, fmt.Errorf("contracts.NewContracts: %w", err)
	}

	receipt, err := eClient.TransactionReceipt(context.Background(), txhash)
	if err != nil {
		return nil, nil, fmt.Errorf("rpc.TransactionReceipt: %w", err)
	}

	for _, log := range receipt.Logs {
		if log == nil {
			slog.Warn("nil log in receipt")
			continue
		}

		adoption, err := contract.ParseSafeHarborAdoption(*log)
		if err != nil {
			continue
		}

		if adoption == nil {
			continue
		}

		agreementAddress = &adoption.AgreementAddress
		agreementContract, err := adoptiondetails.NewAdoptiondetails(adoption.AgreementAddress, eClient)
		if err != nil {
			return nil, nil, fmt.Errorf("contracts.NewAdoptiondetails: %w", err)
		}

		details, err := agreementContract.GetDetails(nil)
		if err != nil {
			return nil, nil, fmt.Errorf("adoptiondetails.GetDetails: %w", err)
		}

		agreement = &details
		return agreementAddress, agreement, nil
	}

	return nil, nil, fmt.Errorf("no safe harbor adoption logs found in transaction receipt")
}

// GetAgreementAddress parses the SafeHarborAdoption event in the tx receipt and returns the new details contract address.
func GetAgreementAddress(
	txhash common.Hash,
	eClient *ethclient.Client,
) (*common.Address, error) {
	contract, err := NewSafeHarborV3Filterer(common.HexToAddress("0x0"), eClient)
	if err != nil {
		return nil, fmt.Errorf("contracts.NewContracts: %w", err)
	}

	receipt, err := eClient.TransactionReceipt(context.Background(), txhash)
	if err != nil {
		return nil, fmt.Errorf("rpc.TransactionReceipt: %w", err)
	}

	for _, log := range receipt.Logs {
		slog.Info("log found", "address", log.Address.Hex(), "txhash", log.TxHash.Hex())
		if log == nil {
			slog.Warn("nil log in receipt")
			continue
		}

		adoption, err := contract.ParseSafeHarborAdoption(*log)
		if err != nil || adoption == nil {
			continue
		}

		addr := adoption.AgreementAddress
		return &addr, nil
	}

	return nil, fmt.Errorf("no safe harbor adoption logs found in transaction receipt")
}
