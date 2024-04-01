test-http:
	curl -L --proxy https://localhost:8443 --proxy-cacert certs/server.pem http://google.com

test-https:
	curl -Lv --proxy https://localhost:8443 --proxy-cacert certs/server.pem https://google.com

health:
	curl --cacert certs/server.pem --key certs/server.key https://localhost:8443/health