import os
import sys
from socket import socket, AF_INET, SOCK_STREAM
import logging
from datetime import datetime
import time

logger = logging.getLogger(__name__)
HOST = ''
PORT = 50009

soc = socket(AF_INET, SOCK_STREAM)
soc.bind((HOST, PORT))
soc.listen(5)



def find_logs(path):
    ...


processes = []


def reap_children():
    # убрать завершившиеся дочерние процессы,
    while processes:
        # иначе может переполниться системная таблица
        pid, stat = os.waitpid(0, os.WNOHANG) # не блокировать сервер, если
        if not pid: break
        # дочерний процесс не завершился
        processes.remove(pid)


def handle_client(connection):
    print("handling")
    while True:
        # чтение, запись в сокет клиента
        data = connection.recv(1024)
        # до получения признака eof, когда
        if not data: break
        # сокет будет закрыт клиентом
        print(data)
        connection.send(b'Echo=>' + data)
        connection.close()
        sys.exit()


def serve():
    while True:
        connection, addr = soc.accept()
        logger.info(f"logserver {datetime.now()} Поключился {addr}")
        reap_children()
        pid = os.fork()
        if pid == 0:
            handle_client(connection)
        else:
            processes.append(pid)
            print(processes, "<= processes")
            # sys.exit()


# if __name__ == "main":
serve()
