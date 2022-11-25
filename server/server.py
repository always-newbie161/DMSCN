import socket
import threading
import os

SERVER_IP = '192.168.0.126'
SERVER_NAME = "hawk_server"	# use <ROLL_NUMBER_server>
clients_dict = {'192.168.0.212':'thor', '192.168.0.180':'hulk', '192.168.0.135':'gamora', '192.168.0.126':'hawk'}  # will store name and socket object
#hospital_id = { 'Hospital1':'u0w7i9g1bz', 'Hospital2':'u0nmn3jhfk', 'Hospital3':'u0isumctro'}
os.environ["APIKEY"] = "u0i0a0kgqg-Wq4D05CxQ5dsvS8hRwYAWqivXTHCMlL+scbaAMQz8Nw="
os.environ["SUBMITTER"] = "u0w7i9g1bz"
os.environ["USER_ID"] = SERVER_NAME

PORT = 6969
ADDR = (SERVER_IP, PORT)
DISCONNECT_MSG = '!EXIT'
BUF_SIZE = 1024

# making a server socket for devices in the same network to connect
server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
server.bind(ADDR)

def send_allbytes(sock, data, flags=0):
    nbytes = sock.send(data, flags)
    if nbytes > 0:
        return send_allbytes(sock, data[nbytes:], flags)
    else:
        return None


def send_client(name, msg):
    for clients in clients_dict:
        if clients == name:
            send_allbytes(clients_dict[clients], msg.encode())


def handle_client(name, conn, addr):
    print(f"[NEW CONNECTION] {addr} connected, Name: {name}")
    print(f"[ACTIVE CONNECTIONS] {threading.activeCount() - 1}") ##
    #send_all(name, f"[NEW CONNECTION] Name: {name}")
    
    # ASSUME: For now only for equipment... (message format is "CATEGORY ITEM")
    msg = conn.recv(BUF_SIZE).decode()
    asset_id = msg.split(' ')[1]
    
    # Transfer asset: Run rrr with required arguments
    # cases required for other categories? 
    os.system("./rrr_ram "+ asset_id+" "+name)
    # Inform client that asset has been transferred
    msg = f"ASSET TRANSFER COMPLETE: {msg} from {SERVER_IP} to {addr[0]}"
    send_allbytes(conn, msg.encode())

    msg2 = conn.recv(BUF_SIZE).decode()
    if msg2 == DISCONNECT_MSG:
        print(f"[USER DISCONNECTED] {addr} disconnected, Name: {name}")
        #send_all(name, f"[USER DISCONNECTED] Name: {name}")

    conn.close()


def start():
    server.listen()
    print(f"Server is listening on {SERVER_IP}")
    while True:
        # we wait on below line for a new connection
        conn, addr = server.accept()
        # when a new connection occur, we store and socket object's corresponding name
        # but first check if name already exists and if it does then ask client alternate name

        # ASSUMING: Every connection attempt is a legitimate request, so asset will get transferred
        accepted_msg = 'ACCESS GRANTED..'
        name = clients_dict[addr[0]]
        send_allbytes(conn, accepted_msg.encode())
        thread = threading.Thread(target=handle_client, args=(name, conn, addr), daemon=True)
        thread.start()


print("[STARTING] server is starting...")
start()
