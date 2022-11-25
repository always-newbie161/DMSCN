# Medi-Net
This is a DMSCN (Decentralized Medical Supply Chain Network) application we have developed to narrow the gap between demand of medical resources in hospitals. We have three broad categories in medical resources - organs, medicines and equipments.

# Submissions
This submission contains 3 subfolders:
* `DMSCN-cli`: This folder consists of the folders needed for DMSCN client. It contains an executable named `kgf` that provides an interface for user to interact with the deployed public blockchain network.
* `server`: Every hospital in MediNet has its own server which receives asset transfer requests from other hospitals. This folder contains a file named `server.py` and its dependencies that listens to client requests.
* `client`: Every hospital in MediNet has a client which can send asset transfer requests to other hospitals. This folder contains a file named `client.py` and its dependencies that take input as medical resource needed and sends requests for the same to the nearest hospital.
* `APIKEY.txt`: This file contains has the credentials to connect to MediNet, i.e., the blockchain network that we have deployed publicly using Kaleido.

**NOTE**: Please mail us at one of these accounts before you start the run. 
We would wake the environment for you to access the blockchain network. As the blockchain is deployed on cloud platforms, it costs us money to host it for 24 hours, we shall specifically host it whenever needed.

1. cs19btech11024@iith.ac.in
2. cs19btech11011@iith.ac.in
3. es19btech11007@iith.ac.in
4. ma19btech11005@iith.ac.in

# Instructions to run DMSCN-cli
1. Read the API Key in APIKEY.txt and use that to run this command.
	`` export APIKEY=<APIKEY>``
2. To build the executable run
	`` go build -o kgf ``
3. To run the executable run
	``./kgf``

The execution begins with taking input of Hospital you belong to as input. Then, the client would be given an option to choose between the various options provided by our blockchain network.

*NOTE*: Set USER\_ID=<ROLL\_NUMBER> in `main.go` file

# Instructions to run server and client
The following equipments would be needed to demonstrate the city-wide simulation using server and client workflow in our blockchain network.
1. A Router (repurposed as a switch)
2. 3 Laptops (preferably Ubuntu)
3. 3 LAN cables

First step is to connect all the laptops to the router using LAN cables and disconnect the network from Internet. Then, each laptop is connected to Internet using Mobile Hotspot (to ensure the devices are in no way connected to the IITH Wifi). Now we can confirm if pings are possible between laptops in the local network.

The `server.py` should be executed on each laptop. The `client.py` should be executed specifically on the laptop that is in need of medical resources. Once the transactions are done, you can verify it using the `kgf` executable.

*NOTE*: Set CLIENT\_NAME=<ROLL\_NUMBER\_client> in `client.py` file.
*NOTE* Set SERVER\_NAME=<ROLL\_NUMBER\_server> in `server.py` file.

# Workflow
(Assuming assets are already present in the network)
1. Client can request for one of 3 types of assets:
	* O - Organ 
	* E - Equipment
	* M - Medicine
2. getAllAssets() is invoked and information about all assets present in the network is retrieved from the ledger on the network. This is filtered and the list of hospitals which have the relevant assets is obtained.
3. The closest hospital from this list is determined so that the asset can be delivered quickly. The client sends an asset transfer request to the closest hospital.
4. The selected hospital's server responds to the request received from client by transfering assets to that client through blockchain.
