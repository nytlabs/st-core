curl "localhost:7071/groups" -d '{"group":0}'
curl "localhost:7071/blocks" -d '{"type":"+","group":1}'
curl "localhost:7071/blocks" -d '{"type":"delay", "group":1}'
curl "localhost:7071/blocks/3/routes/1" -X PUT -d '{"type":"const","value":"1s"}'
curl "localhost:7071/connections" -d '{"source":{"id":2, "Route":0}, "target":{"id":3, "Route":0}}'
curl "localhost:7071/groups/1/export"



