#sudo certbot --nginx -d whittileaks.com -d curso.whittileaks.com
upstream buffalo_app {
    server 127.0.0.1:3000;
}

#server {
#	listen 443 ssl;
#	server_name curso.whittileaks.com;
#	ssl_certificate /etc/nginx/ssl/cert.pem;
#	ssl_certificate_key /etc/nginx/ssl/key.pem;
#	location / {
#			proxy_redirect   off;
#			proxy_set_header Host              $http_host;
#			proxy_set_header X-Real-IP         $remote_addr;
#			proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
#			proxy_set_header X-Forwarded-Proto $scheme;
#
#			proxy_pass       http://buffalo_app;
#			# First attempt to serve request as file, then
#			# as directory, then fall back to displaying a 404.
#			# try_files $uri $uri/ =404;
#		}
#	# return 301 http://$server_name$request_uri;
#}

server {
	listen 80 default_server;
	#listen [::]:80 default_server;
	server_name curso.whittileaks.com;

	# Hide NGINX version (security best practice)
	server_tokens off;

	

	location / {
		proxy_redirect   off;
		proxy_set_header Host              $http_host;
		proxy_set_header X-Real-IP         $remote_addr;
		proxy_set_header X-Forwarded-For   $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Proto $scheme;

		proxy_pass       http://buffalo_app;
		# First attempt to serve request as file, then
		# as directory, then fall back to displaying a 404.
		# try_files $uri $uri/ =404;
	}

	# deny access to .htaccess files, if Apache's document root
	# concurs with nginx's one
	#
	#location ~ /\.ht {
	#	deny all;
	#}
}


# Virtual Host configuration for example.com
#
# You can move that to a different file under sites-available/ and symlink that
# to sites-enabled/ to enable it.
#
#server {
#	listen 80;
#	listen [::]:80;
#
#	server_name example.com;
#
#	root /var/www/example.com;
#	index index.html;
#
#	location / {
#		try_files $uri $uri/ =404;
#	}
#}

	# SSL configuration
	#
	# listen 443 ssl default_server;
	# listen [::]:443 ssl default_server;
	#
	# Note: You should disable gzip for SSL traffic.
	# See: https://bugs.debian.org/773332
	#
	# Read up on ssl_ciphers to ensure a secure configuration.
	# See: https://bugs.debian.org/765782
	#
	# Self signed certs generated by the ssl-cert package
	# Don't use them in a production server!
	#
	# include snippets/snakeoil.conf;

	Set Up NGINX
    certbot can automatically configure NGINX for SSL/TLS. It looks for and modifies the server block in your NGINX configuration that contains a server_name directive with the domain name you’re requesting a certificate for. In our example, the domain is www.example.com.

    Assuming you’re starting with a fresh NGINX install, use a text editor to create a file in the /etc/nginx/conf.d directory named domain‑name.conf (so in our example, www.example.com.conf).

    Specify your domain name (and variants, if any) with the server_name directive:

    server {
        listen 80 default_server;
        listen [::]:80 default_server;
        root /var/www/html;
        server_name example.com www.example.com;
    }
    Save the file, then run this command to verify the syntax of your configuration and restart NGINX:

    $ nginx -t && nginx -s reload
    3. Obtain the SSL/TLS Certificate
    The NGINX plug‑in for certbot takes care of reconfiguring NGINX and reloading its configuration whenever necessary.

    Run the following command to generate certificates with the NGINX plug‑in:

    $ sudo certbot --nginx -d example.com -d www.example.com
    Respond to prompts from certbot to configure your HTTPS settings, which involves entering your email address and agreeing to the Let’s Encrypt terms of service.

    When certificate generation completes, NGINX reloads with the new settings. certbot generates a message indicating that certificate generation was successful and specifying the location of the certificate on your server.