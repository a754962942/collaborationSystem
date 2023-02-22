package code_gen

import "testing"

func TestCenStruct(t *testing.T) {
	//GenStruct("project", "Project")
	GenProtoMessage("project", "ProjectMessage")
}
