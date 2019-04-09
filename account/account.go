package account

import (
	"encoding/json"
	"fmt"

	"github.com/Nik-U/pbc"
	"github.com/p1gd0g/CTCoin/crypto"
	"github.com/p1gd0g/CTCoin/transaction"
)

// Account includes keys.
type Account struct {
	a  string
	A  string
	AT string

	Wallet      [99]transaction.Transaction
	WalletIndex int
}

// New creats a new account.
func (account *Account) New() {

	params, _ := pbc.NewParamsFromString(crypto.Params)
	pairing := params.NewPairing()

	A := pairing.NewG1()
	AT := pairing.NewG1()
	a := pairing.NewZr()

	a.Rand()

	G := pairing.NewG1()
	T := pairing.NewG1()
	G.SetString(crypto.G, 10)
	T.SetString(crypto.T, 10)

	A.PowZn(G, a)
	AT.PowZn(T, a)

	account.a = a.String()
	account.A = A.String()
	account.AT = AT.String()
}

// ShowAccount shows a, A, AT.
func (account *Account) ShowAccount() {
	fmt.Println("a:", account.a)
	fmt.Println("A:", account.A)
	fmt.Println("AT:", account.AT)
}

// Mine is used to mine a coinbase transaction.
// Notice: this is not real mining work.
func (account *Account) Mine() transaction.Transaction {

	params, _ := pbc.NewParamsFromString(crypto.Params)
	pairing := params.NewPairing()

	var tx transaction.Transaction
	RA := pairing.NewG1()
	RT := pairing.NewG1()
	V := pairing.NewG1()
	P := pairing.NewG1()

	r := pairing.NewZr()
	r.Rand()

	T := pairing.NewG1()
	T.SetString(crypto.T, 10)

	A := pairing.NewG1()
	A.SetString(account.A, 10)

	tx.Base = true
	RA.PowZn(A, r)
	RT.PowZn(T, r)
	P.SetString(crypto.GetP(r.String(), account.AT, account.A), 10)
	V.SetString(crypto.GetV(r.String(), account.AT), 10)

	tx.RA = RA.String()
	tx.RT = RT.String()
	tx.P = P.String()
	tx.V = V.String()
	tx.Setp(crypto.Getp(account.a, RT.String()))

	account.Wallet[account.WalletIndex] = tx
	account.WalletIndex++

	fmt.Println("\033[1;31;40m", "We mined a new tx:", tx.P, "\033[0m")

	return tx
}

// Check is used to check if the tx belongs to me.
func (account *Account) Check(tx transaction.Transaction) bool {

	P := crypto.GetP(account.a, tx.RT, account.A)

	if P == tx.P {
		fmt.Println("\033[1;31;40m", "It belongs to us!", "\033[0m")

		tx.Setp(crypto.Getp(account.a, tx.RT))

		account.Wallet[account.WalletIndex] = tx
		account.WalletIndex++
		return true
	}
	return false

}

// ShowWalletIndex is uesd to show wallet index.
func (account *Account) ShowWalletIndex() {
	fmt.Println("\033[1;31;40m", account.WalletIndex, "\033[0m")

}

// ShowTx is used to show a transaction belongs to us.
func (account *Account) ShowTx(P string) bool {
	for i := 0; i < account.WalletIndex; i++ {
		if P == account.Wallet[i].P {
			JSON, _ := json.MarshalIndent(account.Wallet[i], "", "    ")

			fmt.Println(string(JSON))
			return true

		}
	}
	return false
}
