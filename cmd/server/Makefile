
run:
	go generate

test:
	@go test -mod vendor .

secret:
	echo "$CLOUDFUNC_KEY" > "./secret/ncrypt-key.pem"

flood:
	echo "GET https://us-central1-ncrypt.cloudfunctions.net/Stream" \
       | vegeta attack -rate=5 -duration=120s