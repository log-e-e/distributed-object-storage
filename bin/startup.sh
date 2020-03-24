systemFolderPath=/home/${USER}/storage-system
# make objects directory as the storage folder
mkdir -p "${systemFolderPath}"/storage/objects

# list the ip address if could get
chmod +x ./get_ip_addr.sh
./get_ip_addr.sh

# startup a background task to monitor the server
# if server stopped then clear the all of the files created by the server
nohup ./port_monitor.sh "${systemFolderPath}" &
if [ $? -eq 0 ]; then
  echo "SUCCESS: Process monitor startup success"
else
  echo "FAILED: Process monitor startup failed"
fi

# write STORAGE_ROOT & LISTEN_ADDRESS to file
port=12345
envFilePath=../config/server.env

if [ -e $envFilePath ]; then
  rm $envFilePath
fi
echo "LISTEN_ADDRESS=:$port" >> $envFilePath
echo "STORAGE_ROOT=${systemFolderPath}/storage" >> $envFilePath

# startup server
cd ..
echo "--------------------------- Server Start ---------------------------"
echo "Server Port: $port"
go run server.go
