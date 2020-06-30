package helpers

type MockReply []map[string]interface{}

type QueryMock struct {
	Query string
	Reply MockReply
}

type ExpectedResponse struct {
	StatusCode int
	BodyPart   string
}

type TestCase struct{
	TestName  string
	Request   interface{}
	QueryMock *QueryMock
	Expected  ExpectedResponse
}
