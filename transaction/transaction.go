package transaction

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/Nik-U/pbc"
	"github.com/p1gd0g/CTCoin/crypto"
)

// Transaction is the transaction struction.
type Transaction struct {
	Base   bool
	Inputs [3]Input
	P      string
	V      string
	RT     string
	RA     string
	Sig    Sig

	p string
}

// Sig is the signature struction.
type Sig struct {
	I string
	C [3]string
	R [3]string
}

// Input includes P and V.
type Input struct {
	P string
	V string
}

// NewTx is used to construct the new tx being sent.
func NewTx(
	tx1 Transaction, tx2 Transaction, tx3 Transaction, A string, AT string) (
	tx Transaction) {

	tx.Base = false

	params, _ := pbc.NewParamsFromString(crypto.Params)
	pairing := params.NewPairing()

	pointA := pairing.NewG1()
	pointA.SetString(A, 10)

	pointr := pairing.NewZr()
	pointr.Rand()

	pointG := pairing.NewG1()
	pointG.SetString(crypto.G, 10)

	pointT := pairing.NewG1()
	pointT.SetString(crypto.T, 10)

	tx.P = crypto.GetP(pointr.String(), AT, A)
	tx.V = crypto.GetV(pointr.String(), AT)

	pointRA := pairing.NewG1()
	pointRT := pairing.NewG1()

	pointRA.PowZn(pointA, pointr)
	pointRT.PowZn(pointT, pointr)

	tx.RA = pointRA.String()
	tx.RT = pointRT.String()

	pointp := pairing.NewZr()
	pointp.SetString(tx1.p, 10)

	tx.Inputs[0].P = tx1.P
	tx.Inputs[1].P = tx2.P
	tx.Inputs[2].P = tx3.P

	tx.Inputs[0].V = tx1.V
	tx.Inputs[1].V = tx2.V
	tx.Inputs[2].V = tx3.V

	V1 := pairing.NewG1()
	V1.SetString(tx.Inputs[0].V, 10)
	V2 := pairing.NewG1()
	V2.SetString(tx.Inputs[1].V, 10)
	V3 := pairing.NewG1()
	V3.SetString(tx.Inputs[2].V, 10)

	pointV := pairing.NewG1()
	pointV.SetString(tx.V, 10)

	pointI := pairing.NewG1()
	pointI.PowZn(V1, pointp)

	tx.Sig.I = pointI.String()

	pointm := pairing.NewZr()

	m, _ := json.Marshal(tx)
	sum := sha256.Sum256(m)
	pointm.SetFromHash(sum[:])

	mLR := pairing.NewG1()

	L1 := pairing.NewG1()
	L2 := pairing.NewG1()
	L3 := pairing.NewG1()

	R1 := pairing.NewG1()
	R2 := pairing.NewG1()
	R3 := pairing.NewG1()

	// P1, _ := pairing.NewG1().SetString(tx1.P, 10)
	P2, _ := pairing.NewG1().SetString(tx2.P, 10)
	P3, _ := pairing.NewG1().SetString(tx3.P, 10)

	q := pairing.NewZr()
	w := pairing.NewZr()

	c1 := pairing.NewZr()
	c2 := pairing.NewZr()
	c3 := pairing.NewZr()

	q.Rand()
	w.Rand()
	L3.PowZn(pointG, q)
	L3.ThenAdd(pairing.NewG1().PowZn(P3, w))
	R3.PowZn(V3, q)

	R3.ThenAdd(pairing.NewG1().PowZn(pointI, w))
	c3.SetString(w.String(), 10)
	tx.Sig.C[2] = w.String()
	tx.Sig.R[2] = q.String()

	q.Rand()
	w.Rand()
	L2.PowZn(pointG, q)
	L2.ThenAdd(pairing.NewG1().PowZn(P2, w))
	R2.PowZn(V2, q)
	R2.ThenAdd(pairing.NewG1().PowZn(pointI, w))
	c2.SetString(w.String(), 10)
	tx.Sig.C[1] = w.String()
	tx.Sig.R[1] = q.String()

	q.Rand()
	w.Rand()
	L1.PowZn(pointG, q)
	R1.PowZn(V1, q)

	mLR.ThenAdd(L1).ThenAdd(L2).ThenAdd(L3).ThenAdd(R1).ThenAdd(R2).ThenAdd(R3)
	mLR.ThenPowZn(pointm)
	c := pairing.NewZr()
	c.SetString(crypto.Hash(mLR).String(), 10)
	c1.SetString(c.String(), 10)
	c1.ThenSub(c3).ThenSub(c2)

	tx.Sig.C[0] = c1.String()
	tx.Sig.R[0] = q.ThenSub(c1.ThenMulZn(pointp)).String()

	return
}

// Setp is used to set p.
func (tx *Transaction) Setp(p string) {
	tx.p = p
}

// Verify is used to verify the new tx.
func Verify(tx Transaction) bool {

	params, _ := pbc.NewParamsFromString(crypto.Params)
	pairing := params.NewPairing()

	c1 := pairing.NewZr()
	c1.SetString(tx.Sig.C[0], 10)
	c2 := pairing.NewZr()
	c2.SetString(tx.Sig.C[1], 10)
	c3 := pairing.NewZr()
	c3.SetString(tx.Sig.C[2], 10)

	r1 := pairing.NewZr()
	r1.SetString(tx.Sig.R[0], 10)
	r2 := pairing.NewZr()
	r2.SetString(tx.Sig.R[1], 10)
	r3 := pairing.NewZr()
	r3.SetString(tx.Sig.R[2], 10)

	L1 := pairing.NewG1()
	L2 := pairing.NewG1()
	L3 := pairing.NewG1()

	R1 := pairing.NewG1()
	R2 := pairing.NewG1()
	R3 := pairing.NewG1()

	P1 := pairing.NewG1()
	P1.SetString(tx.Inputs[0].P, 10)
	P1.ThenPowZn(c1)
	P2 := pairing.NewG1()
	P2.SetString(tx.Inputs[1].P, 10)
	P2.ThenPowZn(c2)
	P3 := pairing.NewG1()
	P3.SetString(tx.Inputs[2].P, 10)
	P3.ThenPowZn(c3)

	V1 := pairing.NewG1()
	V1.SetString(tx.Inputs[0].V, 10)
	V1.ThenPowZn(r1)
	V2 := pairing.NewG1()
	V2.SetString(tx.Inputs[1].V, 10)
	V2.ThenPowZn(r2)
	V3 := pairing.NewG1()
	V3.SetString(tx.Inputs[2].V, 10)
	V3.ThenPowZn(r3)

	c := pairing.NewZr()
	c.ThenAdd(c1).ThenAdd(c2).ThenAdd(c3)

	L1.SetString(crypto.G, 10)
	L1.ThenPowZn(r1).ThenAdd(P1)
	L2.SetString(crypto.G, 10)
	L2.ThenPowZn(r2).ThenAdd(P2)
	L3.SetString(crypto.G, 10)
	L3.ThenPowZn(r3).ThenAdd(P3)

	R1.SetString(tx.Sig.I, 10)
	R1.ThenPowZn(c1).ThenAdd(V1)
	R2.SetString(tx.Sig.I, 10)
	R2.ThenPowZn(c2).ThenAdd(V2)
	R3.SetString(tx.Sig.I, 10)
	R3.ThenPowZn(c3).ThenAdd(V3)

	for i := 0; i < 3; i++ {
		tx.Sig.C[i] = ""
		tx.Sig.R[i] = ""
	}

	m, _ := json.Marshal(tx)
	sum := sha256.Sum256(m)
	pointm := pairing.NewZr()
	pointm.SetFromHash(sum[:])

	mLR := pairing.NewG1()
	mLR.ThenAdd(L1).ThenAdd(L2).ThenAdd(L3).ThenAdd(R1).ThenAdd(R2).ThenAdd(R3)
	mLR.ThenPowZn(pointm)

	if crypto.Hash(mLR).String() == c.String() {
		return true
	}
	return false

}

// Trace is used to trace a transaction.
func (tx *Transaction) Trace(t string) {
	var A, AT string
	A, AT = crypto.GetReceiver(tx.P, t, tx.RA)

	fmt.Println("\033[1;31;40m", "A:", A, "\033[0m")
	fmt.Println("\033[1;31;40m", "AT:", AT, "\033[0m")

	x := crypto.Getx(tx.Sig.I)

	for i := 0; i < 3; i++ {
		if x == crypto.Gete(t, tx.Inputs[i].P) {
			fmt.Println("\033[1;31;40m", "Real transaction:",
				tx.Inputs[i].P, "\033[0m")
			break
		}
	}

}
