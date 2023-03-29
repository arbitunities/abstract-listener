.PHONY: bindings abigen

run:
	go run cmd/listener/main.go

bindings:
	curl https://api.etherscan.io/api\?module\=contract\&action\=getabi\&address\=0x0576a174D229E3cFA37253523E645A78A0C91B57 | jq -r '.result' > ./bindings/entrypoint.abi

abigen: bindings
	abigen --abi ./bindings/entrypoint.abi --pkg erc4337 --out pkg/erc4337/entrypoint.go
	
