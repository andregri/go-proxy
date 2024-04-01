test-http:
	curl -L --proxy http://localhost:8080 http://google.com

test-https:
	curl -Lv --proxy http://localhost:8080 https://example.org

health:
	curl http://localhost:8080/health