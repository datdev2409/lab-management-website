- Build application from local using `make build`
- Run `make deploy` in local to copy the main file to /home/ec2-user/ directory in EC2
- Create environment config file `.env.production`
- We use systemd to run the go web as the background process. Copy the file `goweb.service` to /lib/systemd/system/ directory
- Run the following command to start the goweb.service

```
sudo systemctl daemon-reload
sudo systemctl start goweb
sudo systemctl status goweb
```

- View the service log using command `journalctl -u goweb`
