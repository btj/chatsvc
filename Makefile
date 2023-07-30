build:
	docker build -t bartjacobs/chatsvc chatsvc

push:
	docker push bartjacobs/chatsvc

restart:
	gcloud compute instances reset chatsvc

ssh:
	gcloud compute ssh chatsvc