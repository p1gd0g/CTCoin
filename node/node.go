package node

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/p1gd0g/CTCoin/account"
	"github.com/p1gd0g/CTCoin/crypto"
	"github.com/p1gd0g/CTCoin/transaction"

	"github.com/Nik-U/pbc"
)

// Node includes node and wallet.
type Node struct {
	Account account.Account

	Pool      [99]transaction.Transaction
	PoolIndex int

	PeerList  [10]string
	PeerIndex int

	UsedImage      [99]string
	UsedImageIndex int
}

// NewTx is used to send new tx.
func (node *Node) NewTx() {

	var A, A1, A2 string
	var AT, AT1, AT2 string

	fmt.Println("\033[1;31;40m", "Please input the A:", "\033[0m")
	fmt.Scan(&A1, &A2)
	A = A1 + " " + A2

	fmt.Println("\033[1;31;40m", "Please input the AT:", "\033[0m")
	fmt.Scan(&AT1, &AT2)
	AT = AT1 + " " + AT2

	tx := transaction.NewTx(node.Account.Wallet[node.Account.WalletIndex-1], node.Pool[node.PoolIndex-1], node.Pool[node.PoolIndex-2], A, AT)
	node.Account.WalletIndex--

	node.Pool[node.PeerIndex] = tx
	node.PoolIndex++
	node.SendTransaction(tx)
}

// StartConn starts the connection.
func (node *Node) StartConn() {
	fmt.Println("\033[1;31;40m", "Please input your port number:", "\033[0m")

	fmt.Scan(&node.PeerList[node.PeerIndex])
	node.PeerIndex++

	l, _ := net.Listen("tcp", ":"+node.PeerList[0])
	go func() {
		for {
			conn, _ := l.Accept()
			node.handleConnection(conn)
		}
	}()

}

func (node *Node) handleConnection(conn net.Conn) {
	txJSON := make([]byte, 11111)
	lengh, _ := conn.Read(txJSON)

	var tx transaction.Transaction
	json.Unmarshal(txJSON[0:lengh], &tx)

	fmt.Println("\033[1;31;40m", "Received a new tx:", tx.P, "\033[0m")

	if node.usedI(tx.Sig.I) {
		fmt.Println("\033[1;31;40m", "Double spending!", "\033[0m")
	} else {
		if tx.Base == false {

			if transaction.Verify(tx) {
				fmt.Println("\033[1;31;40m", "Verified.", "\033[0m")
				if node.Account.Check(tx) == false {
					node.Pool[node.PoolIndex] = tx
					node.PoolIndex++
				}
			} else {
				fmt.Println("\033[1;31;40m", "Wrong transaction.", "\033[0m")
			}

		} else {
			node.Pool[node.PoolIndex] = tx
			node.PoolIndex++
		}
	}
}

// AddPeer adds the new peer port number.
func (node *Node) AddPeer() {
	fmt.Println("\033[1;31;40m", "Input the port number:", "\033[0m")

	fmt.Scan(&node.PeerList[node.PeerIndex])
	node.PeerIndex++
}

// SendTransaction sends the transaction to other nodes.
func (node *Node) SendTransaction(tx transaction.Transaction) {
	txJSON, _ := json.Marshal(tx)

	for i := 1; i < node.PeerIndex; i++ {
		conn, _ := net.Dial("tcp", ":"+node.PeerList[i])
		conn.Write(txJSON)
	}
}

// ShowPeerIndex is used to show peer index.
func (node *Node) ShowPeerIndex() {
	fmt.Println("\033[1;31;40m", node.PeerIndex, "\033[0m")

}

// ShowPoolIndex is uesd to show pool index.
func (node *Node) ShowPoolIndex() {
	fmt.Println("\033[1;31;40m", node.PoolIndex, "\033[0m")

}

// HandleCommand is uesd to handle command.
func (node *Node) HandleCommand(command string) {
	switch command {
	case "exit":
		os.Exit(0)
	case "addPeer":
		node.AddPeer()
	case "showAccount":
		node.Account.ShowAccount()
	case "mine":
		node.SendTransaction(node.Account.Mine())
	case "newTx":
		node.NewTx()
	case "showPeerIndex":
		node.ShowPeerIndex()
	case "showPoolIndex":
		node.ShowPoolIndex()
	case "showWalletIndex":
		node.Account.ShowWalletIndex()
	case "trace":
		node.trace()
	case "showTx":
		node.ShowTx()

	default:
		fmt.Println("\033[1;31;40m", "What?", "\033[0m")
	}
}

func (node *Node) usedI(I string) bool {
	for i := 0; i < node.UsedImageIndex; i++ {
		if I == node.UsedImage[i] {
			return true
		}
	}
	return false
}

func (node *Node) trace() {
	fmt.Println("\033[1;31;40m", "Please input t:", "\033[0m")

	var t string
	fmt.Scan(&t)

	params, _ := pbc.NewParamsFromString(crypto.Params)
	pairing := params.NewPairing()

	G, _ := pairing.NewG1().SetString(crypto.G, 10)
	pointt, _ := pairing.NewZr().SetString(t, 10)
	G.ThenPowZn(pointt)

	if G.String() == crypto.T {
		fmt.Println("\033[1;31;40m", "Hello, superman!", "\033[0m")

		var P1, P2, P string
		fmt.Println("\033[1;31;40m", "Input the tx's P:", "\033[0m")
		fmt.Scan(&P1, &P2)
		P = P1 + " " + P2

		for i := 0; i < node.PoolIndex; i++ {
			if node.Pool[i].P == P {
				node.Pool[i].Trace(t)
				break
			}
		}

	}
}

// ShowTx is uesd to show a transaction.
func (node *Node) ShowTx() {
	fmt.Println("\033[1;31;40m", "Please input P of tx:", "\033[0m")

	var P, P1, P2 string
	fmt.Scan(&P1, &P2)
	P = P1 + " " + P2

	if node.Account.ShowTx(P) == false {
		for i := 0; i < node.PoolIndex; i++ {
			if node.Pool[i].P == P {
				JSON, _ := json.MarshalIndent(node.Pool[i], "", "    ")
				fmt.Println(string(JSON))
				break
			}
		}

	}
}
