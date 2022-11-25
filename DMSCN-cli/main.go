package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"						// the libraries are imported here
	"sync"
	"time"
	"encoding/json"

	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kaleido-io/kaleido-fabric-go/fabric"
	"github.com/kaleido-io/kaleido-fabric-go/kaleido"
)

type Asset struct {
	ID             string `json:"ID"`
	Type		   string `json:"asset_type"`
	Name		   string `json:"name"`
	Bloodgroup	   string `json:"bloodgroup"`		// the structure of Asset is defined here
	Stock          int    `json:"stock"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

func main() {				/* MAIN STARTS HERE */
	rand.Seed(time.Now().UTC().UnixNano())

	username := os.Getenv("USER_ID")
	if username == "" {
		username = "hawk"					// username needs to be configured here
	}

	ccname := os.Getenv("CCNAME")
	if ccname == "" {
		ccname = "asset_transfer"
	}

	var err error
	var channel *kaleido.Channel

	fmt.Println("==========================================================================================================================")
	
	fmt.Println("Using the Fabric SDK for transaction submission")
	network := kaleido.NewNetwork()
	network.Initialize()
	config, err := fabric.BuildConfig(network)			// the setting up of network and its related configs
	if err != nil {
		fmt.Printf("Failed to generate network configuration for the SDK: %s\n", err)
		os.Exit(1)
	}

	sdk1 := newSDK(config)
	defer sdk1.Close()

	wallet := kaleido.NewWallet(username, *network, sdk1)
	err = wallet.InitIdentity()
	if err != nil || wallet.Signer == nil {				// the wallet is defined here for the user
		fmt.Printf("Failed to initiate wallet: %v\n", err)
		os.Exit(1)
	}

	fabric.AddTlsConfig(config, wallet.Signer)

	sdk2 := newSDK(config)
	defer sdk2.Close()

	channel = kaleido.NewChannel("default-channel", sdk2)
	err = channel.Connect(wallet.Signer.Identifier())		// the connection to channel happens here
	if err != nil {
		fmt.Printf("Failed to connect to channel: %s\n", err)
		os.Exit(1)
	}

	initChaincode := os.Getenv("INIT_CC")
	if initChaincode == "" {
		initChaincode = "false"
	} else {
		initChaincode = strings.ToLower(initChaincode)
	}

	if initChaincode == "true" {
		var Args [][]byte
		err = channel.InitChaincode(ccname, "InitLedger", Args)
												// if the initialization is not done for the chaincode it will be done
		if err != nil {
			fmt.Printf("Failed to initialize chaincode: %s\n", err)
		}
	} else {
		var count int
		countStr := os.Getenv("TX_COUNT")
		if countStr != "" {
			count, err = strconv.Atoi(countStr)
			if err != nil {
				fmt.Printf("Failed to convert %s to integer", countStr)
				os.Exit(1)
			}
		} else {
			count = 1
		}
		if count > 50 {
			fmt.Println("Error: TX_COUNT cannot exceed 50")
			os.Exit(1)
		}
		var batches int
		batchStr := os.Getenv("BATCHES")
		if batchStr != "" {
			batches, err = strconv.Atoi(batchStr)
			if err != nil {
				fmt.Printf("Failed to convert %s to integer", batchStr)
				os.Exit(1)
			}
		} else {
			batches = 1
		}

		for i := 0; true; i++ {
			var wg sync.WaitGroup
			for j := 0; j < count; j++ {
				wg.Add(1)
				go func(idx int) {

					defer wg.Done()

					var option int
					ListAllFuncs()
					fmt.Printf("Kya Chahiye rey tereko?? ")
					fmt.Scanf("%d",&option)

					var fun string

					switch option {

					case 1:
						var tp,name,bg,q,owner,av string

						fmt.Printf("Type of Asset: (O/M/E) ")
						fmt.Scanf("%s",&tp)

						fmt.Printf("Name: ")
						fmt.Scanf("%s",&name)

						if tp == "O" {
							fmt.Printf("Blood Group: ")
							fmt.Scanf("%s",&bg)

							q = "1"

						} else
						if tp == "M" {
							fmt.Printf("Quantity: ")
							fmt.Scanf("%s",&q)
						} else
						if tp == "E" {
							fmt.Printf("Quantity: ")
							fmt.Scanf("%s",&q)
						}

						owner = username
						
						fmt.Printf("Appraised Value: ")
						fmt.Scanf("%s",&av)

						fun = "CreateAsset"
						Args := [][]byte{[]byte(tp), []byte(name), []byte(bg), []byte(q), []byte(owner), []byte(av)}

						fmt.Printf("=> Batch %d: Send transaction %d of %d (Option: %d)\n", i+1, idx+1, count)
						err = channel.ExecChaincode(ccname, fun, Args)
						
						if err != nil {
							fmt.Printf("=> Batch %d: Failed to send transaction %d (Option: %d). %s\n", i+1, idx+1, err)
						}

					case 2:
						var id string

						fmt.Printf("Asset ID: ")
						fmt.Scanf("%s",&id)

						fun = "DeleteAsset"
						Args := [][]byte{[]byte(id)}

						fmt.Printf("=> Batch %d: Send transaction %d of %d (Option: %d)\n", i+1, idx+1, count)
						err = channel.ExecChaincode(ccname, fun, Args)
						
						if err != nil {
							fmt.Printf("=> Batch %d: Failed to send transaction %d (Option: %d). %s\n", i+1, idx+1, err)
						}

					case 3:
						var id,newOwner string

						fmt.Printf("Asset ID: ")
						fmt.Scanf("%s",&id)

						fmt.Printf("New Owner: ")
						fmt.Scanf("%s",&newOwner)

						fun = "TransferAsset"
						Args := [][]byte{[]byte(id), []byte(newOwner)}

						fmt.Printf("=> Batch %d: Send transaction %d of %d (Option: %d)\n", i+1, idx+1, count)
						err = channel.ExecChaincode(ccname, fun, Args)
						
						if err != nil {
							fmt.Printf("=> Batch %d: Failed to send transaction %d (Option: %d). %s\n", i+1, idx+1, err)
						}

					case 4:
						var id string

						fmt.Printf("Asset ID: ")
						fmt.Scanf("%s",&id)

						fun = "ReadAsset"
						Args := [][]byte{[]byte(id)}

						fmt.Printf("=> Batch %d: Send transaction %d of %d (Option: %d)\n", i+1, idx+1, count)
						res := channel.ReadChaincode(ccname, fun, Args)
						
						data := fmt.Sprintf("%+v\n", res.Responses[0])
						data = strings.Split(data,">")[0]
						data = strings.Split(data,"payload:")[1]
						data = strings.Replace(data, "\\", "", -1)
						data = data[1:len(data)-2]

						var asset Asset
						err := json.Unmarshal( []byte(data), &asset)
						
						if err != nil {
							panic(err)
						}

						printJSON(asset)

					case 5:
						var id, attr, newVal string

						fmt.Printf("Asset ID: ")
						fmt.Scanf("%s",&id)

						fmt.Printf("Attribute to modify: ")
						fmt.Scanf("%s",&attr)

						fmt.Printf("New Value: ")
						fmt.Scanf("%s",&newVal)

						fun = "UpdateAsset"
						Args := [][]byte{[]byte(id), []byte(attr), []byte(newVal)}

						fmt.Printf("=> Batch %d: Send transaction %d of %d (Option: %d)\n", i+1, idx+1, count)
						err = channel.ExecChaincode(ccname, fun, Args)
						
						if err != nil {
							fmt.Printf("=> Batch %d: Failed to send transaction %d (Option: %d). %s\n", i+1, idx+1, err)
						}

					case 6:
						fun = "GetAllAssets"
						Args := [][]byte{[]byte("No Args")}

						fmt.Printf("=> Batch %d: Send transaction %d of %d (Option: %d)\n", i+1, idx+1, count)
						res := channel.ReadChaincode(ccname, fun, Args)
						
						data := fmt.Sprintf("%+v\n", res.Responses[0])
						data = strings.Split(data,">")[0]
						data = strings.Split(data,"payload:")[1]

						dataArr := strings.Split(data, "},")
						assets := strtoArray(dataArr)

						fmt.Println("[")
						for i:=0; i<len(assets); i++ {
							printJSON(*assets[i])
							println("")
						}
						fmt.Println("]")
					}
				}(j)
			}
			wg.Wait()

			fmt.Printf("\nCompleted batch %d of %d\n\n", i+1, batches)

			var isdone string
			fmt.Printf("Done? ")
			fmt.Scanf("%s", &isdone)
	
			if isdone == "yes" {
				break
			}
		}

		fmt.Printf("\nAll Done!\n")
	}
}


func newSDK(config map[string]interface{}) *fabsdk.FabricSDK {
	configProvider, err := fabric.NewConfigProvider(config)
	if err != nil {
		fmt.Printf("Failed to create config provider from config map: %s\n", err)
		os.Exit(1)
	}

	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		fmt.Printf("Failed to instantiate an SDK: %s\n", err)
		os.Exit(1)
	}
	return sdk
}

func ListAllFuncs() {
	
	var functions[7]string

	functions[0] = "1. Create an Asset"
	functions[1] = "2. Delete an Asset"
	functions[2] = "3. Transfer an Asset"
	functions[3] = "4. Read an Asset"
	functions[4] = "5. Update an Asset"
	functions[5] = "6. Get All Assets"

	for i:=0; i<len(functions); i++ {
		fmt.Println(functions[i])
	}
}

func printJSON(asset Asset) {
	fmt.Println("Asset{")
	fmt.Println("ID:", asset.ID)
	fmt.Println("Type:", asset.Type)
	fmt.Println("Name:", asset.Name)
	fmt.Println("Blood Group:", asset.Bloodgroup)
	fmt.Println("Stock:", asset.Stock)
	fmt.Println("Owner:", asset.Owner)
	fmt.Println("Appraisedvalue:", asset.AppraisedValue)
	fmt.Println("}")
}

func strtoArray(data []string) []*Asset{
	
	data[0] = data[0][2:]

	for i:=0 ; i<len(data)-1; i++ {
		data[i] = data[i] + "}"
	}
	data[len(data)-1] = data[len(data)-1][:len(data[len(data)-1])-3]

	var assets []*Asset
	for i:=0; i<len(data); i++ {

		var asset Asset
		d := data[i]

		d = strings.Replace(d, "\\", "", -1)
		err := json.Unmarshal([]byte(d), &asset)

		if err != nil {
			panic(err)
		}

		assets = append(assets, &asset)
	}

	return assets
}
