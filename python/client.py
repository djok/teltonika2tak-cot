import socket
import time

# Create a TCP/IP socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_address = ('127.0.0.1', 8887)
sock.connect(server_address)


IMEI = "000F333536333037303432343431303133"


# Test case (Codec 8)
packet_data = """00000000000000A7080400000113FC208DFF000F14F650209CCA80006F00D60400040004030101150316030001460000015D0000000113FC17610B000F14FFE0209CC580006E00C00500010004030101150316010001460000015E0000000113FC284945000F150F00209CD200009501080400000004030101150016030001460000015D0000000113FC267C5B000F150A50209CCCC0009300680400000004030101150016030001460000015B00040000BA48"""


# Test case (Codec 8Ex)
# packet_data = """000000000000005F8E010000015FBA40B620000F0DCDE420959D30008A000006000000000006000100EF0000010011001E000100100000CBDF0002000B000000003544C875000E0000000029BFE4D100010100000412345678010000D153"""


# Test case (Codec 16)
# packet_data = """000000000000005F10020000016BDBC7833000000000000000000000000000000000000B05040200010000030002000B00270042563A00000000016BDBC7871800000000000000000000000000000000000B05040200010000030002000B00260042563A00000200005FB312132"""


# Test case when data split into two pieces
# packet_data = """00000000000000A7080400000113FC208DFF000F14F650209CCA80006F00D60400040004030101150316030001460000015D0000000113FC17610B000F14FFE0209CC580006E00C00500010004030101150316010001460000015E0000000113FC284945000F150F00209CD200009501080400000004030101150016030001460000015D0000000113FC267C5B000F150A50209CCCC0009300680400000004030101150016030001460000015B00040000BA4800000000000000A7080400000113FC208DFF000F14F650209CCA80006F00D60400040004030101150316030001460000015D0000000113FC17610B000F14FFE0209CC580006E00C00500010004030101150316010001460000015E0000000113FC284945000F150F00209CD200009501080400000004030101150016030001460000015D0000000113FC267C5B000F150A50209CCCC0009300"""
#
# packet_data2 = """680400000004030101150016030001460000015B00040000B458"""


# Test case when packet split into three pieces
# packet_data = """00000000000000A7080400000113FC208DFF000F14F650209CCA80006F00D60400040004030101150316030001460000015D0000000113FC17610B000F14FFE0209CC580006E00C00500010004030101150"""
#
# packet_data2 = """316010001460000015E0000000113FC284945000F150F00209CD200009501080400000004030101150016030001460000015D0000000113FC267C5B000"""
#
# packet_data3 = """F150A50209CCCC0009300680400000004030101150016030001460000015B00040000BA48"""

# Send the message

sock.sendall(IMEI.encode('utf-8'))
data = sock.recv(1024)
response = data.decode('utf-8')

print("Response from server: ", response)

if response == '01':
    sock.sendall(packet_data.encode('utf-8'))
    # time.sleep(5)
    # sock.sendall(packet_data2.encode('utf-8'))
    # time.sleep(5)
    # sock.sendall(packet_data3.encode('utf-8'))
elif response == '00':
    print("IMEI corrupted!")
else:
    print('Bad connection!')
sock.close()