Requires a long timeout on your web server, otherwise requests will fail if they hit Newsblur's rate limit.

Like this for nginx:

proxy_connect_timeout       600;
proxy_send_timeout          600;
proxy_read_timeout          600;
send_timeout                600;
