build:
	docker build -t bartjacobs/chatsvc chatsvc

run:
	docker run -e CHATSVC_TLS_CERT -e CHATSVC_TLS_PRIVATE_KEY -p 8443:8443 bartjacobs/chatsvc

push:
	docker push bartjacobs/chatsvc

restart:
	gcloud compute instances reset chatsvc

ssh:
	gcloud compute ssh chatsvc