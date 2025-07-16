syntax = "v1"

info (
	title: "{{.serviceName}} http service"// TODO: add title
	desc: "{{.serviceName}} http service description"// TODO: add description
	author: "{{.gitUser}}"
	email: "{{.gitEmail}}"
)

type request {
	// TODO: add members here and delete this comment
}

type response {
	// TODO: add members here and delete this comment
}
@server (
	desc:    "{{.serviceName}} service"
	version: "v1"
	group:   ""
)
service {{.serviceName}} {
	@handler GetUser // TODO: set handler name and delete this comment
	get /users/id/:userId(request) returns(response)

	@handler CreateUser // TODO: set handler name and delete this comment
	post /users/create(request)
}
