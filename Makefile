build:
	docker build -t bartjacobs/chatsvc chatsvc

run:
	docker run -p 8080:8080 bartjacobs/chatsvc

push:
	docker push bartjacobs/chatsvc

restart:
	gcloud compute instances reset chatsvc

ssh:
	gcloud compute ssh chatsvc