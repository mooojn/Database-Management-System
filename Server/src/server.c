#include "mongoose.h"
#include <stdio.h>
#include <string.h>

static void handler(struct mg_connection *c, int ev, void *ev_data) {
    if (ev == MG_EV_HTTP_MSG) {
        struct mg_http_message *hm = (struct mg_http_message *) ev_data;

        // Handle GET request (Just to confirm server is running)
        if (strncmp(hm->uri.buf, "/", hm->uri.len) == 0) {
            mg_http_reply(c, 200,
                          "Content-Type: application/json\r\n"
                          "Access-Control-Allow-Origin: *\r\n"
                          "Access-Control-Allow-Methods: GET, POST, OPTIONS\r\n"
                          "Access-Control-Allow-Headers: Content-Type\r\n",
                          "{\"message\": \"Hello from cxs backend!\"}");
        }
        // Handle POST request
        else if (strncmp(hm->uri.buf, "/save", hm->uri.len) == 0) {
            char body[1024] = {0};
            snprintf(body, sizeof(body), "%.*s", (int) hm->body.len, hm->body.buf);

            // Save data to file
            FILE *file = fopen("./server/data/server.txt", "a"); // Open in append mode
            if (file) {
                fprintf(file, "%s\n", body);
                fclose(file);
                mg_http_reply(c, 200,
                              "Content-Type: application/json\r\n"
                              "Access-Control-Allow-Origin: *\r\n"
                              "Access-Control-Allow-Methods: GET, POST, OPTIONS\r\n"
                              "Access-Control-Allow-Headers: Content-Type\r\n",
                              "{\"status\": \"success\", \"message\": \"Data saved\"}");
            } else {
                mg_http_reply(c, 500,
                              "Content-Type: application/json\r\n"
                              "Access-Control-Allow-Origin: *\r\n"
                              "Access-Control-Allow-Methods: GET, POST, OPTIONS\r\n"
                              "Access-Control-Allow-Headers: Content-Type\r\n",
                              "{\"status\": \"error\", \"message\": \"Failed to save data\"}");
            }
        }
    }
}

int main() {
    printf("Server is Running on 8000");
    struct mg_mgr mgr;
    mg_mgr_init(&mgr);
    mg_http_listen(&mgr, "http://localhost:8080", handler, &mgr);

    for (;;) mg_mgr_poll(&mgr, 1000);

    mg_mgr_free(&mgr);
    return 0;
}
