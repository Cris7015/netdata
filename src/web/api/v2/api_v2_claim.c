// SPDX-License-Identifier: GPL-3.0-or-later

#include "api_v2_calls.h"
#include "claim/claim.h"

#if defined(OS_WINDOWS)
#include <windows.h>
#include <sys/cygwin.h>
#endif

static char *netdata_random_session_id_filename = NULL;
static nd_uuid_t netdata_random_session_id = { 0 };

bool netdata_random_session_id_generate(void) {
    static char guid[UUID_STR_LEN] = "";

    uuid_generate_random(netdata_random_session_id);
    uuid_unparse_lower(netdata_random_session_id, guid);

    char filename[FILENAME_MAX + 1];
    snprintfz(filename, FILENAME_MAX, "%s/netdata_random_session_id", netdata_configured_varlib_dir);

    bool ret = true;

    (void)unlink(filename);

    // save it
    int fd = open(filename, O_WRONLY|O_CREAT|O_TRUNC|O_CLOEXEC, 640);
    if(fd == -1) {
        netdata_log_error("Cannot create random session id file '%s'.", filename);
        ret = false;
    }
    else {
        if (write(fd, guid, UUID_STR_LEN - 1) != UUID_STR_LEN - 1) {
            netdata_log_error("Cannot write the random session id file '%s'.", filename);
            ret = false;
        } else {
            ssize_t bytes = write(fd, "\n", 1);
            UNUSED(bytes);
        }
        close(fd);
    }

    if(ret && (!netdata_random_session_id_filename || strcmp(netdata_random_session_id_filename, filename) != 0)) {
        freez(netdata_random_session_id_filename);
        netdata_random_session_id_filename = strdupz(filename);
    }

    return ret;
}

static const char *netdata_random_session_id_get_filename(void) {
    if(!netdata_random_session_id_filename)
        netdata_random_session_id_generate();

    return netdata_random_session_id_filename;
}

static bool netdata_random_session_id_matches(const char *guid) {
    if(uuid_is_null(netdata_random_session_id))
        return false;

    nd_uuid_t uuid;

    if(uuid_parse(guid, uuid))
        return false;

    if(uuid_compare(netdata_random_session_id, uuid) == 0)
        return true;

    return false;
}

static bool check_claim_param(const char *s) {
    if(!s || !*s) return true;

    do {
        if(isalnum((uint8_t)*s) || *s == '.' || *s == ',' || *s == '-' || *s == ':' || *s == '/' || *s == '_')
            ;
        else
            return false;

    } while(*++s);

    return true;
}

int api_v2_claim(RRDHOST *host __maybe_unused, struct web_client *w, char *url) {
    char *key = NULL;
    char *token = NULL;
    char *rooms = NULL;
    char *base_url = NULL;

    while (url) {
        char *value = strsep_skip_consecutive_separators(&url, "&");
        if (!value || !*value) continue;

        char *name = strsep_skip_consecutive_separators(&value, "=");
        if (!name || !*name) continue;
        if (!value || !*value) continue;

        if(!strcmp(name, "key"))
            key = value;
        else if(!strcmp(name, "token"))
            token = value;
        else if(!strcmp(name, "rooms"))
            rooms = value;
        else if(!strcmp(name, "url"))
            base_url = value;
    }

    BUFFER *wb = w->response.data;
    buffer_flush(wb);
    buffer_json_initialize(wb, "\"", "\"", 0, true, BUFFER_JSON_OPTIONS_DEFAULT);

    time_t now_s = now_realtime_sec();
    CLOUD_STATUS status = buffer_json_cloud_status(wb, now_s);

    bool can_be_claimed = false;
    switch(status) {
        case CLOUD_STATUS_AVAILABLE:
        case CLOUD_STATUS_OFFLINE:
        case CLOUD_STATUS_INDIRECT:
            can_be_claimed = true;
            break;

        case CLOUD_STATUS_BANNED:
        case CLOUD_STATUS_ONLINE:
            can_be_claimed = false;
            break;
    }

    buffer_json_member_add_boolean(wb, "can_be_claimed", can_be_claimed);

    if(can_be_claimed && key) {
        if(!netdata_random_session_id_matches(key)) {
            buffer_reset(wb);
            buffer_strcat(wb, "invalid key");
            netdata_random_session_id_generate(); // generate a new key, to avoid an attack to find it
            return HTTP_RESP_FORBIDDEN;
        }

        if(!token || !base_url || !check_claim_param(token) || !check_claim_param(base_url) || (rooms && !check_claim_param(rooms))) {
            buffer_reset(wb);
            buffer_strcat(wb, "invalid parameters");
            netdata_random_session_id_generate(); // generate a new key, to avoid an attack to find it
            return HTTP_RESP_BAD_REQUEST;
        }

        netdata_random_session_id_generate(); // generate a new key, to avoid an attack to find it

        bool success = false;
        const char *msg;
        if(claim_agent(base_url, token, rooms, cloud_config_proxy_get(), cloud_config_insecure_get())) {
            msg = "ok";
            success = true;
            can_be_claimed = false;
            status = claim_reload_and_wait_online();
        }
        else
            msg = claim_agent_failure_reason_get();

        // our status may have changed
        // refresh the status in our output
        buffer_flush(wb);
        buffer_json_initialize(wb, "\"", "\"", 0, true, BUFFER_JSON_OPTIONS_DEFAULT);
        now_s = now_realtime_sec();
        buffer_json_cloud_status(wb, now_s);

        // and this is the status of the claiming command we run
        buffer_json_member_add_boolean(wb, "success", success);
        buffer_json_member_add_string_or_empty(wb, "message", msg);
    }

    if(can_be_claimed) {
        const char *filename = netdata_random_session_id_get_filename();
        CLEAN_BUFFER *buffer = buffer_create(0, NULL);

        const char *os_filename;
        const char *os_prefix;
        const char *os_quote;
        const char *os_message;

#if defined(OS_WINDOWS)
        char win_path[MAX_PATH];
        cygwin_conv_path(CCP_POSIX_TO_WIN_A, filename, win_path, sizeof(win_path));
        os_filename = win_path;
        os_prefix = "more";
        os_message = "We need to verify this Windows server is yours. So, open a Command Prompt on this server to run the command. It will give you a UUID. Copy and paste this UUID to this box:";
#else
        os_filename = filename;
        os_prefix = "sudo cat";
        os_message = "We need to verify this server is yours. SSH to this server and run this command. It will give you a UUID. Copy and paste this UUID to this box:";
#endif

        // add quotes only when the filename has a space
        if(strchr(os_filename, ' '))
            os_quote = "\"";
        else
            os_quote = "";

        buffer_sprintf(buffer, "%s %s%s%s", os_prefix, os_quote, os_filename, os_quote);
        buffer_json_member_add_string(wb, "key_filename", os_filename);
        buffer_json_member_add_string(wb, "cmd", buffer_tostring(buffer));
        buffer_json_member_add_string(wb, "help", os_message);
    }

    buffer_json_agents_v2(wb, NULL, now_s, false, false);
    buffer_json_finalize(wb);

    return HTTP_RESP_OK;
}