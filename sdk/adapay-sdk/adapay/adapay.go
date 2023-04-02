package adapay

import (
	adapayCore2 "go-practice/sdk/adapay-sdk/adapay-core"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Adapay struct {
	MultiMerchSysConfigs map[string]*adapayCore2.MerchSysConfig

	DefaultMerchSysConfig *adapayCore2.MerchSysConfig
}

func InitDefaultMerchSysConfig(filePath string) (*Adapay, error) {

	config, err := adapayCore2.ReadMerchConfig(filePath)
	if err != nil {
		return nil, err
	}

	ada := &Adapay{
		DefaultMerchSysConfig: config,
	}

	return ada, nil
}

func InitMultiMerchSysConfigs(fileDir string) (*Adapay, error) {

	dirs, _ := ioutil.ReadDir(fileDir)

	configs := map[string]*adapayCore2.MerchSysConfig{}

	for _, f := range dirs {

		ext := filepath.Ext(f.Name())
		if ext == ".json" {
			config, err := adapayCore2.ReadMerchConfig(fileDir + f.Name())
			if err != nil {
				continue
			}

			key := strings.Replace(f.Name(), ".json", "", -1)
			configs[key] = config
		}
	}

	ada := &Adapay{
		MultiMerchSysConfigs: configs,
	}

	return ada, nil
}

func (a *Adapay) HandleConfig(multiMerchConfigId ...string) *adapayCore2.MerchSysConfig {
	if multiMerchConfigId == nil {
		return a.DefaultMerchSysConfig
	} else {
		return a.MultiMerchSysConfigs[multiMerchConfigId[0]]
	}
}

func (a *Adapay) Payment() *Payment {
	return &Payment{Adapay: a}
}

func (a *Adapay) PaymentConfirm() *PaymentConfirm {
	return &PaymentConfirm{Adapay: a}
}

func (a *Adapay) PaymentReverse() *PaymentReverse {
	return &PaymentReverse{Adapay: a}
}

func (a *Adapay) SettleAccount() *SettleAccount {
	return &SettleAccount{Adapay: a}
}

func (a *Adapay) AdapayTools() *AdapayTools {
	return &AdapayTools{Adapay: a}
}

func (a *Adapay) Drawcash() *Drawcash {
	return &Drawcash{Adapay: a}
}

func (a *Adapay) CorpMember() *CorpMember {
	return &CorpMember{Adapay: a}
}

func (a *Adapay) Member() *Member {
	return &Member{Adapay: a}
}

func (a *Adapay) Refund() *Refund {
	return &Refund{Adapay: a}
}

func (a *Adapay) Wallet() *Wallet {
	return &Wallet{Adapay: a}
}

func (a *Adapay) Account() *Account {
	return &Account{Adapay: a}
}

func (a *Adapay) Checkout() *Checkout {
	return &Checkout{Adapay: a}
}

func (a *Adapay) FastPay() *FastPay {
	return &FastPay{Adapay: a}
}

func (a *Adapay) FreezeAccount() *FreezeAccount {
	return &FreezeAccount{Adapay: a}
}

func (a *Adapay) UnFreezeAccount() *UnFreezeAccount {
	return &UnFreezeAccount{Adapay: a}
}

func (a *Adapay) Transfer() *Transfer {
	return &Transfer{Adapay: a}
}

func (a *Adapay) Version() string {
	return "1.3.1"
}
