package crypto

// t = 29604782977066162081595096072912354690597154740

import (
	"crypto/sha256"

	"github.com/Nik-U/pbc"
)

var (
	// Params is the params of curve.
	Params = "type a q 8780710799663312522437781984754049815806883199414208211028653399266475630880222957078625179422662221423155858769582317459277713367317481324925129998224791 h 12016012264891146079388821366740534204802954401251311822919615131047207289359704531102844802183906537786776	r 730750818665451621361119245571504901405976559617 exp2 159	exp1 107 sign1 1 sign0 1"
	// T is the point which only CT knows the t.
	T = "[5727231827648595756897531848864122869918498496389756357381956054802360556819220753878539820754263673737876517302591142859995575085760327485328632960142531, 1045876788512743813205685182282517734295812023819168609016928302617792470636319077962426882438035762869150434937085117731863275756782969659017147641202555]"
	// G is the base point.
	G = "[3284601785124645061415289361210903980486649982340224943676791530933218310418638211100170192559965291154150008722329421133655442419580766313306988057699959, 5465062710369093639733293185865191620639276093999042698514377738820441053578196551708707048356302767994301600169627093794992328413765563974934120780024253]"
)

// GetP is used to get P.
func GetP(r string, AT string, A string) string {
	params, _ := pbc.NewParamsFromString(Params)
	pairing := params.NewPairing()

	pointG := pairing.NewG1()
	pointG.SetString(G, 10)

	pointAT := pairing.NewG1()
	pointAT.SetString(AT, 10)

	pointr := pairing.NewZr()
	pointr.SetString(r, 10)

	pointA := pairing.NewG1()
	pointA.SetString(A, 10)

	pointP := pairing.NewG1()
	pointP.PowZn(pointAT, pointr)

	hash := Hash(pointP)
	pointP.PowZn(pointG, hash)
	pointP.ThenAdd(pointA)

	return pointP.String()
}

// GetV is used to get V.
func GetV(r string, AT string) string {
	params, _ := pbc.NewParamsFromString(Params)
	pairing := params.NewPairing()

	pointG := pairing.NewG1()
	pointG.SetString(G, 10)

	pointT := pairing.NewG1()
	pointT.SetString(T, 10)

	pointAT := pairing.NewG1()
	pointAT.SetString(AT, 10)

	pointr := pairing.NewZr()
	pointr.SetString(r, 10)

	pointV := pairing.NewG1()
	pointV.PowZn(pointAT, pointr)

	hash := Hash(pointV)
	pointV.PowZn(pointT, hash)
	pointV.ThenAdd(pointAT)

	return pointV.String()
}

// Getp is used to get p.
func Getp(a string, RT string) string {
	params, _ := pbc.NewParamsFromString(Params)
	pairing := params.NewPairing()

	pointRT := pairing.NewG1()
	pointRT.SetString(RT, 10)

	pointa := pairing.NewZr()
	pointa.SetString(a, 10)

	pointRT.ThenPowZn(pointa)

	pointp := pairing.NewZr()
	pointp.SetString(Hash(pointRT).String(), 10)

	pointp.ThenAdd(pointa)

	return pointp.String()
}

// Hash is used to convert a G1 to Zr.
func Hash(x *pbc.Element) (temp *pbc.Element) {
	params, _ := pbc.NewParamsFromString(Params)
	pairing := params.NewPairing()
	temp = pairing.NewZr()

	sum := sha256.Sum256(x.Bytes())
	temp.SetFromHash(sum[:])
	return
}

// GetReceiver is used to get A and AT of transaction.
func GetReceiver(P string, t string, RA string) (string, string) {
	params, _ := pbc.NewParamsFromString(Params)
	pairing := params.NewPairing()

	pointP, _ := pairing.NewG1().SetString(P, 10)
	pointRA, _ := pairing.NewG1().SetString(RA, 10)
	pointt, _ := pairing.NewZr().SetString(t, 10)
	pointG, _ := pairing.NewG1().SetString(G, 10)

	pointRA.ThenPowZn(pointt)

	hash := Hash(pointRA)
	pointG.ThenPowZn(hash)
	pointP.ThenSub(pointG)

	pointPT := pairing.NewG1().PowZn(pointP, pointt)

	return pointP.String(), pointPT.String()

}

// Getx is used to get x = e(I,G).
func Getx(I string) string {
	params, _ := pbc.NewParamsFromString(Params)
	pairing := params.NewPairing()

	x := pairing.NewGT()
	pointG, _ := pairing.NewG1().SetString(G, 10)
	pointI, _ := pairing.NewG1().SetString(I, 10)

	x.Pair(pointG, pointI)
	return x.String()
}

// Gete is used to get e(P,tP).
func Gete(t string, P string) string {

	params, _ := pbc.NewParamsFromString(Params)
	pairing := params.NewPairing()

	pointP, _ := pairing.NewG1().SetString(P, 10)
	pointt, _ := pairing.NewZr().SetString(t, 10)
	pointPT := pairing.NewG1().PowZn(pointP, pointt)

	return pairing.NewGT().Pair(pointP, pointPT).String()
}
