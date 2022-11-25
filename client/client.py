import datetime
import os
import re
import socket
import threading
from unicodedata import category
from matplotlib.pyplot import close
from numpy import argsort

CLIENT_NAME = "hawk_client"	# use <ROLL_NUMBER_client>
CURRENT_IP = '192.168.0.126'
hospital_dict = { 'thor':'192.168.0.212', 'hulk':'192.168.0.180', 'gamora':'192.168.0.135', 'hawk':'192.168.0.126'}
#hospital_id = { 'Hospital1':'u0w7i9g1bz', 'Hospital2':'u0nmn3jhfk', 'Hospital3':'u0isumctro'}
os.environ["APIKEY"] = "u0i0a0kgqg-Wq4D05CxQ5dsvS8hRwYAWqivXTHCMlL+scbaAMQz8Nw="
os.environ["SUBMITTER"] = "u0w7i9g1bz"
os.environ["USER_ID"] = CLIENT_NAME


PORT = 6969
BUF_SIZE = 1024
DISCONNECT_MSG = '!EXIT'



# ask what this particular client needs
#enter item name
item = input('Enter What do you need?')



os.system("./rrr_bheem")
# i pass  category of item, item name
# above should create a json or a text file with relevant info dumped into it



#hospital_ips = {1:'192.168.0.101', 2:'192.168.0.102', 3:'192.168.0.103' }




#read file and create a list/dict of unique hospital which have the required item
my_file = open("result.txt", "r")
data = my_file.read()
data_into_list = data.split(",")
my_file.close()

# above we have a list of tuple like entries of form item:hospital
# from that we filter the unique tuple entries with the required item
# [id:item:owner, ............]
# hospital_list key is owner and val is id
hospital_list = {}
for i,entries in enumerate(data_into_list):
    temp = entries.split(":")
    if temp[1] == item:
        if temp[2] not in hospital_list.keys():
            hospital_list[temp[2]] = temp[0] 


#cheking parsing is correct
geeky_file = open('geekyfile.txt', 'wt')
geeky_file.write(str(hospital_list))
geeky_file.close()

# now with the list of hospitals calculate distance from source client
# remove option of self distance calculation (hard-code)

# for ping delays
def get_avg_rtt(ip, N):
    ping_output = os.popen(f"ping {ip} -c {N}").read()

    # ping_output should have N+5 lines?
    output_lines = ping_output.split('\n')
    print(len(output_lines))
    print(output_lines)
    val_line = output_lines[N+4]
    print(val_line)
    val_line = val_line[23:]
    print(val_line)
    val_line = val_line[:-3]
    print(val_line)
    rtt_vals = val_line.split('/')
    
    return rtt_vals[1]

# Avg ping delay over N attempts
N = 10
ping_delays = []  # get ping delays for all hospitals in list
hospital_ids = list(hospital_list.keys()) #get the list of ownwers
for id in hospital_ids:
    ping_delays.append(get_avg_rtt(hospital_dict[id], N))

#delayidx_sorted = argsort(ping_delays)
best = argsort(ping_delays)[0]  #just indiex amoung keys
SERVER_IP = hospital_dict[hospital_ids[best]]
ADDR = (SERVER_IP, PORT)

while True:
    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    # name of current client (hard-code)
    
    print(f'[CONNECTING] To {SERVER_IP}:{PORT}')
    client.connect(ADDR)
    access_recv = client.recv(BUF_SIZE).decode()
    print(access_recv)
    if 'GRANTED' in access_recv:
        break
    
    client.close()

print(f"[JOINED SUCESSFULLY], ")


#msg format is "item_id"

def send_allbytes(sock, data, flags=0):
    nbytes = sock.send(data, flags)
    if nbytes > 0:
        return send_allbytes(sock, data[nbytes:], flags)
    else:
        return None


def send_msg(msg):
    #curr_time = datetime.datetime.now().strftime('%d-%m-%Y %H:%M:%S')
    #msg_to_send = f"[{curr_time}] {name}: {msg}"
    send_allbytes(client, msg.encode())


def receive_msg():
    while True:
        global flag
        global FILE_NAME
        global FILE_SIZE
        message_recv = client.recv(BUF_SIZE).decode()
        if "ASSET TRANSFER COMPLETE" in message_recv:
            print(message_recv)
            client.send(DISCONNECT_MSG.encode())
            client.close()

        print(message_recv)



# thread = threading.Thread(target=receive_msg, daemon=True)
# thread.start()

#send first msg as the request for some item

message = "Transfer " + hospital_list[hospital_ids[best]]
send_msg(message)

while True:
    global flag
    global FILE_NAME
    global FILE_SIZE
    message_recv = client.recv(BUF_SIZE).decode()
    print(message_recv)
    if "ASSET TRANSFER COMPLETE" in message_recv:
        #print(message_recv)
        client.send(DISCONNECT_MSG.encode())
        client.close()
        break
    
