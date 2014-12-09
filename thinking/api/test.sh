
curl localhost:7071/block -d'{"name":"A"}'
curl localhost:7071/block -d'{"name":"B"}'
curl localhost:7071/block -d'{"name":"C"}'


curl localhost:7071
echo "\n"
curl localhost:7071/group -d'{"ParentID":0, "ChildrenIDs":[], "MemberIDs":[1,2]}'
echo "\n"
curl localhost:7071
echo "\n"
