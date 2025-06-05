// configs/simple_pipeline.json
{
"name": "SimplePipeline",
"stages": [
{
"name": "validation",
"executionMode": "sequential",
"components": [
{
"type": "ValidateID",
"name": "idValidator",
"params": {
"minID": 1
}
},
{
"type": "ValidateName",
"name": "nameValidator"
}
]
},
{
"name": "dataLoading",
"executionMode": "parallel",
"components": [
{
"type": "LoadUserData",
"name": "userDataLoader",
"provides": ["user"]
},
{
"type": "LoadUserPermissions",
"name": "permissionLoader",
"requires": ["user"],
"dependsOn": ["userDataLoader"]
}
]
},
{
"name": "businessLogic",
"executionMode": "sequential",
"components": [
{
"type": "ProcessUser",
"name": "userProcessor",
"params": {
"debug": true
}
}
]
}
]
}