package jobs

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

func TestIsSystemTask(t *testing.T) {
	t.Log(TT_FORK_JOIN.GetName())
	t.Log(TT_FORK_JOIN.IsSystemTask())
	t.Log(TT_USER_DEFINED.GetName())
	t.Log(TT_USER_DEFINED.IsSystemTask())
}

func TestStageType_UnmarshalYAML(t *testing.T) {
	files, e := ioutil.ReadFile("../example_process.yaml")
	if e != nil {
		t.Log(e.Error())
	}
	var job Job
	yaml.Unmarshal(files, &job)

	t.Log(job)
}
