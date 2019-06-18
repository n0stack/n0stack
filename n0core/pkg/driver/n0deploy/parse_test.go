package n0deploy

import "testing"

func TestParse(t *testing.T) {
	parser := NewLocalParser()

	n0dep, err := parser.Parse(`
		RUN git clone ... \
		 && apt update && apt install -y ...
		COPY ./src /dst

		DEPLOY
		COPY ./src /dst
	`)
	if err != nil {
		t.Errorf("Parse() got error: err=%s", err.Error())
	}

	if n0dep.Bootstrap[0].String() != "RUN git clone ... 		 && apt update && apt install -y ..." {
		t.Errorf("Parse() return wrong bootstrap[0]:\n  have=%s\n  want=%s", n0dep.Bootstrap[0].String(), "RUN git clone ... 		 && apt update && apt install -y ...")
	}
	if n0dep.Bootstrap[1].String() != "COPY ./src /dst" {
		t.Errorf("Parse() return wrong bootstrap[1]:\n have=%s\n  want=%s", n0dep.Bootstrap[1].String(), "COPY ./src /dst")
	}
	if n0dep.Deploy[0].String() != "COPY ./src /dst" {
		t.Errorf("Parse() return wrong deploy[0]:\n have=%s\n  want=%s", n0dep.Bootstrap[1].String(), "COPY ./src /dst")
	}
}
