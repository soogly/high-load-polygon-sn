#include <errno.h>
#include <error.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/epoll.h>
#include <unistd.h>	/* close() */
#include <string.h>
// #include <openssl/sha.h>

#define SIZE 20
#define needed_key_field "Sec-WebSocket-Key:"

char* get_key(char buf[1024]);
char* concat(const char *s1, const char *s2);

int main(void){

    const char UPGRADE_HEADER[] = "HTTP/1.1 101 Switching Protocols\nUpgrade: websocket\nConnection: Upgrade";
    const char HANDHACKE_STRING[] = "dGhlIHNhbXBsZSBub25jZQ==258EAFA5-E914-47DA-95CA-C5AB0DC85B11";
    char *clients_key;
    int listener;
    unsigned int port = 5555;
    struct sockaddr_in addr;
    int pollfd;
    char buff[1024];
    struct epoll_event ev = {
        .events = EPOLLIN,
    };

    listener = socket(AF_INET, SOCK_STREAM, 0);
    if (listener == -1)
        error(1, errno, "socket()");
    int optval = 1;
    setsockopt(listener, SOL_SOCKET, SO_REUSEADDR, &optval, sizeof(optval));
    // setsockopt(listener, SOL_SOCKET, SO_REUSEPORT, &optval, sizeof(optval));

    addr.sin_family = AF_INET;
    addr.sin_port = htons(port);
    addr.sin_addr.s_addr = htonl(INADDR_ANY);

    if (bind(listener, (struct sockaddr*)&addr, sizeof(addr))) {
        error(0, errno, "bind()");
        close(listener);
        return 1;
    }
    if (listen(listener, 5)) {
        error(0, errno, "listen()");
        close(listener);
        return 1;
    }
    pollfd = epoll_create1(0);
    if (pollfd == -1){
        error(0, errno, "epoll_create1()");
        close(listener);
        return 1;
    }
    ev.data.fd = listener;
    if (epoll_ctl(pollfd, EPOLL_CTL_ADD, listener, &ev)){
        error(0, errno, "epoll_ctl(ADD)");
        close(pollfd);
        close(listener);
        return 1;
    }

    while (1){
        int fd;
        printf("before epoll_wait (fd: %d)\n", fd);
        int count = epoll_wait(pollfd, &ev, 2, 1000);
        printf("after epoll_wait (fd: %d)\n", fd);
        if (count == -1){
            error(0, errno, "epoll_wait()");
            close(pollfd);
            close(listener);
            return 1;
        } else if (count == 0)
            continue;

        fd = ev.data.fd;
        if (fd == listener){    // новый клиент
            printf("Connected fd (fd: %d)\n", fd);
            
            // handshake
            printf("FROM LISTENERS SOCKET:\n%s\n", buff); 
            recv(fd, buff, sizeof(buff), 0);
            buff[sizeof(buff)]=0;
            clients_key = get_key(buff);
            char* ccated = concat(clients_key, HANDHACKE_STRING);
            printf("clients_key: %s\n", clients_key);
            printf("CONACTED: %s\n", ccated);
            // size_t length = strlen(ccated);
            // unsigned char hash[SHA_DIGEST_LENGTH];
            // SHA1(data, length, hash);
            free(ccated);

            printf("%s\n", buff); 
            printf("clients_key: %s\n", clients_key); 

            fd = accept(listener, NULL, NULL);
            // if (fd == -1) continue;
            printf("Connected client (fd: %d)\n", fd);
            ev.data.fd = fd;
            if (epoll_ctl(pollfd, EPOLL_CTL_ADD, fd, &ev)){
                error(0, errno, "epoll_ctl(ADD)");
                close(pollfd);
                close(listener);
                return 1;
            }
            int sended;
            sended = send(fd, UPGRADE_HEADER, sizeof(UPGRADE_HEADER), 0);
            if (sended == -1) {
                error(0, errno, "send()");
                close(pollfd);
                close(listener);
                return 1;
            }

        } else {     // клиент изменил данные
            int sended;
            printf("bef recv:\n%s\n", buff); 
            int bytes = recv(fd, buff, sizeof(buff), 0);
            // recv(fd, buff, sizeof(buff), 0);
            printf("aft recv:\n%s\n", buff); 
            buff[sizeof(buff)]=0;
            printf("FROM CLIENTS SOCKET:\n%s\n", buff); 

            if (bytes == -1){
                error(0, errno, "recv()");
                close(pollfd);
                close(listener);
                return 1;
            } else if (bytes == 0) {
                // continue;
                if (epoll_ctl(pollfd, EPOLL_CTL_DEL, fd, NULL)){
                    error(0, errno, "epoll_ctl(DEL)");
                    close(pollfd);
                    close(listener);
                    return 1;
                }
                close(fd);
                printf("Disconnected client (fd: %d)\n", fd);
                continue;
            }
            printf("Trying to send (fd: %d)\n", fd);

            sended = send(fd, UPGRADE_HEADER, sizeof(UPGRADE_HEADER), 0);
            if (sended == -1) {
                error(0, errno, "send()");
                close(pollfd);
                close(listener);
                return 1;
            }
        }

    }   
    close(pollfd);
    close(listener);
    return 0;
}

char* concat(const char *s1, const char *s2)
{
    const size_t len1 = strlen(s1);
    const size_t len2 = strlen(s2);
    char *result = malloc(len1 + len2 + 1); // +1 for the null-terminator
    // in real code you would check for errors in malloc here
    memcpy(result, s1, len1);
    memcpy(result + len1, s2, len2 + 1); // +1 to copy the null-terminator
    return result;
}

char* get_key(char buf[1024]) {
   char *token = strtok(buf, "\n");
   char *tmp;
   // loop through the string to extract all other tokens
   while( token != NULL ) {

      tmp = strtok(token, " ");
      if (!strcmp(tmp, needed_key_field)){
         tmp = strtok(NULL, " ");
         printf( "res => %s\n", tmp);
         break;
      }
      token = strtok(NULL, " ");
   }
   return tmp;
}