#include <errno.h>
#include <error.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/epoll.h>
#include <unistd.h>	/* close() */
#include <string.h>
// #include <sha.h>
#include <openssl/sha.h>
#include <openssl/pem.h>

#define SIZE 20
// #define needed_key_field "Sec-WebSocket-Key:"

const char UPGRADE_HEADER[] = "HTTP/1.1 101 Switching Protocols\nUpgrade: websocket\nConnection: Upgrade";
const char HANDHACKE_STRING[] = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11";
char hand_responce[] = "HTTP/1.1 101 Switching Protocols\nUpgrade: websocket\nConnection: Upgrade\nSec-WebSocket-Accept: ";
// hand_responce[sizeof(hand_responce)] = '\0';
char sub_protocol[] = "Sec-WebSocket-Protocol: chat";
    
int* get_header(char buf[1024], char *needed_key_field, char *hval);
char* concat(const char *s1, const char *s2);
char *make_sec_key(char *buff);


char *base64encode2 (const void *b64_encode_this, int encode_this_many_bytes){
    BIO *b64_bio, *mem_bio;      //Declares two OpenSSL BIOs: a base64 filter and a memory BIO.
    BUF_MEM *mem_bio_mem_ptr;    //Pointer to a "memory BIO" structure holding our base64 data.
    b64_bio = BIO_new(BIO_f_base64());                      //Initialize our base64 filter BIO.
    mem_bio = BIO_new(BIO_s_mem());                           //Initialize our memory sink BIO.
    BIO_push(b64_bio, mem_bio);            //Link the BIOs by creating a filter-sink BIO chain.
    BIO_set_flags(b64_bio, BIO_FLAGS_BASE64_NO_NL);  //No newlines every 64 characters or less.
    BIO_write(b64_bio, b64_encode_this, encode_this_many_bytes); //Records base64 encoded data.
    BIO_flush(b64_bio);   //Flush data.  Necessary for b64 encoding, because of pad characters.
    BIO_get_mem_ptr(mem_bio, &mem_bio_mem_ptr);  //Store address of mem_bio's memory structure.
    BIO_set_close(mem_bio, BIO_NOCLOSE);   //Permit access to mem_ptr after BIOs are destroyed.
    BIO_free_all(b64_bio);  //Destroys all BIOs in chain, starting with b64 (i.e. the 1st one).
    BUF_MEM_grow(mem_bio_mem_ptr, (*mem_bio_mem_ptr).length + 1);   //Makes space for end null.
    (*mem_bio_mem_ptr).data[(*mem_bio_mem_ptr).length] = '\0';  //Adds null-terminator to tail.
    return (*mem_bio_mem_ptr).data; //Returns base-64 encoded data. (See: "buf_mem_st" struct).
}


static char encoding_table[] = {'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
                                'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
                                'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
                                'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f',
                                'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
                                'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
                                'w', 'x', 'y', 'z', '0', '1', '2', '3',
                                '4', '5', '6', '7', '8', '9', '+', '/'};
static char *decoding_table = NULL;
static int mod_table[] = {0, 2, 1};


char *base64_encode(const unsigned char *data,
                    size_t input_length,
                    size_t *output_length) {

    *output_length = 4 * ((input_length + 2) / 3);

    char *encoded_data = malloc(*output_length);
    if (encoded_data == NULL) return NULL;

    for (int i = 0, j = 0; i < input_length;) {

        uint32_t octet_a = i < input_length ? (unsigned char)data[i++] : 0;
        uint32_t octet_b = i < input_length ? (unsigned char)data[i++] : 0;
        uint32_t octet_c = i < input_length ? (unsigned char)data[i++] : 0;

        uint32_t triple = (octet_a << 0x10) + (octet_b << 0x08) + octet_c;

        encoded_data[j++] = encoding_table[(triple >> 3 * 6) & 0x3F];
        encoded_data[j++] = encoding_table[(triple >> 2 * 6) & 0x3F];
        encoded_data[j++] = encoding_table[(triple >> 1 * 6) & 0x3F];
        encoded_data[j++] = encoding_table[(triple >> 0 * 6) & 0x3F];
    }

    for (int i = 0; i < mod_table[input_length % 3]; i++)
        encoded_data[*output_length - 1 - i] = '=';

    return encoded_data;
}


unsigned char *base64_decode(const char *data,
                             size_t input_length,
                             size_t *output_length) {

    if (decoding_table == NULL) build_decoding_table();

    if (input_length % 4 != 0) return NULL;

    *output_length = input_length / 4 * 3;
    if (data[input_length - 1] == '=') (*output_length)--;
    if (data[input_length - 2] == '=') (*output_length)--;

    unsigned char *decoded_data = malloc(*output_length);
    if (decoded_data == NULL) return NULL;

    for (int i = 0, j = 0; i < input_length;) {

        uint32_t sextet_a = data[i] == '=' ? 0 & i++ : decoding_table[data[i++]];
        uint32_t sextet_b = data[i] == '=' ? 0 & i++ : decoding_table[data[i++]];
        uint32_t sextet_c = data[i] == '=' ? 0 & i++ : decoding_table[data[i++]];
        uint32_t sextet_d = data[i] == '=' ? 0 & i++ : decoding_table[data[i++]];

        uint32_t triple = (sextet_a << 3 * 6)
        + (sextet_b << 2 * 6)
        + (sextet_c << 1 * 6)
        + (sextet_d << 0 * 6);

        if (j < *output_length) decoded_data[j++] = (triple >> 2 * 8) & 0xFF;
        if (j < *output_length) decoded_data[j++] = (triple >> 1 * 8) & 0xFF;
        if (j < *output_length) decoded_data[j++] = (triple >> 0 * 8) & 0xFF;
    }

    return decoded_data;
}


void build_decoding_table() {

    decoding_table = malloc(256);

    for (int i = 0; i < 64; i++)
        decoding_table[(unsigned char) encoding_table[i]] = i;
}


void base64_cleanup() {
    free(decoding_table);
}

int main(void){

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
        int count = epoll_wait(pollfd, &ev, 2, 1000);
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
            printf("FROM LISTENERS SOCKET: \n%s\n", buff); 
            recv(fd, buff, sizeof(buff), 0);
            printf("Buff: %s\n", buff); 

            // buff[sizeof(buff)] = '\0';
            // char *clients_key ;
            // if (clients_key = get_header(buff, "sss") != NULL){
            //     printf("responce: %s\n", handshake(buff));
            // }; 

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
            int bytes = recv(fd, buff, sizeof(buff), 0);
            buff[sizeof(buff)] = 0;
            printf("FROM CLIENTS SOCKET: fd:%d; %d bytes\n%s\n", fd, bytes, buff); 

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
            char *sec_key = make_sec_key(&buff);

            printf("*sec_key: %s \n", sec_key);
            char *resp = concat(hand_responce, sec_key);
            int resp_len = strlen(resp);
            printf("response: %d bytes;\n %s\n", resp_len, resp);
            printf("Trying to send (fd: %d)\n", fd);
            sended = send(fd, resp, resp_len, 0);
            printf("sended: %d ", sended);

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

// void clear_row(char *row, int row_len){
//     int i=0;
//     for (i; i < row_len; i++)
//         row[i] = '0';
//     row[i+1] = '\0';
// }
int *get_header(char buf[1024], char *header_name, char *hval) {

    // printf("GET_HEADER()\n");
    // printf("header_name: %s\n", header_name);

    int HEADER_NAME_LEN = strlen(header_name);
    int hval_len = 0;
    int skip_row_flag = 0;
    // printf("HEADER_NAME_LEN]: %d\n", HEADER_NAME_LEN);

    for (int i=0, j=0; buf[j] != '\0'; i++, j++){
        if (i == HEADER_NAME_LEN){
            i=-1;
        }
        // printf("header_name[%d]: %c; buf[%d]: %c; skip_row_flag: %d\n", i, header_name[i], j, buf[j], skip_row_flag);
        if (!skip_row_flag){
            if (header_name[i] != buf[j]){
                // printf("header_name[i] != buf[j])\n");
                i = -1;
                skip_row_flag = 1;
            } else {
                // printf("header_name[%d] == buf[%d]: %c\n", i, j, header_name[i]);
                if (header_name[i] == buf[j] && i == HEADER_NAME_LEN-1){
                    printf("\n===\n");
                    j++;
                    while (buf[++j] != '\n')
                    {   
                        hval[hval_len++] = buf[j];
                    }
                    hval[hval_len] = '\0';
                    break;
                }
            }
            continue;
        } else {
            if (buf[j] == '\n'){
                skip_row_flag = 0;
                i = -1;
            }
            continue;
        }   
    }
    return hval_len;
}

char *make_sec_key(char *buff){
    // static int resp_len = 0;
    char *base64_encoded; 
    char *clients_key[25];
    int key_len = get_header(buff, "Sec-WebSocket-Key:", &clients_key);
    char *ccated = concat(clients_key, HANDHACKE_STRING);

    printf("clients_key: %s\n", clients_key);
    printf("CONCATED: %s\n", ccated);
    unsigned char hash[SHA_DIGEST_LENGTH] = {0,};
    SHA1(hash, strlen(ccated), ccated);
    printf("SHA-1 hash: %x\n", hash); 
    // unsigned char key_out[64] = {0,};
    // base64_encode(hash, key_out, sizeof(hash));
    free(ccated);
    
    int bytes_to_encode = strlen(hash);
    int sizz = 1024;
    int *out_len = &sizz;
    base64_encoded = base64_encode(hash, bytes_to_encode, out_len); 
    printf("Encoded (base-64): %s\n", base64_encoded); 
    return base64_encoded;
    // for (int i = 0; i < )
    // *--hand_responce = 
    // free(base64_encoded);
    // printf("hand_responce: \n%s\n", hand_responce); 

    // *resp = concat(hand_responce, base64_encoded);
    // printf("RESP IN: %s\n", hash); 

    // resp_len = strlen(resp);
    // return resp_len;
}