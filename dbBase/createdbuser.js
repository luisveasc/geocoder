use geocoderdb;
db.createUser({ user: 'geocoderuser', pwd: 'Q3Gb4wcjw9U83p!',
		roles: [ { role: "readWrite", db: "geocoderdb" } ]
	});

db.getCollection('SysVariables').insertMany([
	{
		"_id": "roles",
		"data": ["admin","user"]
	  }
])
exit;
