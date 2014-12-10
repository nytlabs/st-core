
curl localhost:7071/block -d'{"name":"A"}'
curl localhost:7071/block -d'{"name":"B"}'
curl localhost:7071/block -d'{"name":"C"}'


printf "\n"
printf "\n"
curl localhost:7071
printf "\n"
curl localhost:7071/group -d'{"ParentID":0, "ChildIDs":[1,2]}'
printf "\n"
printf "\n"
curl localhost:7071
