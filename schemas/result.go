package schemas

import (
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/golang/glog"
)

var visNode map[string]*cgraph.Node

type ResponseCallBack func(*Result)

type Result struct {
	path   []*Grammar
	output []string
}

func NewResult() *Result {
	return &Result{
		path: []*Grammar{},
	}
}

func (r *Result) AddNode(n *Grammar) *Result {
	r.path = append(r.path, n)
	return r
}

func (r *Result) AddOutput(s string) *Result {
	r.output = append(r.output, s)
	return r
}

func (r *Result) GetPath() []*Grammar {
	return r.path
}

func (r *Result) GetOutput() []string {
	return r.output
}

func (r *Result) Visualize(filename string) {
	gh := graphviz.New()
	graph, _ := gh.Graph()

	var prev *cgraph.Node = nil
	var outputIdx int = 0
	for _, n := range r.GetPath() {
		current, err := graph.CreateNode(n.GetContent())
		if err != nil {
			glog.Fatalf("something unexpected when noding %s: %v", n.GetID(), err)
		}
		if n.gtype == GrammarTerminal {
			// current.Set(r.GetOutput()[outputIdx])
			outputNode, err := graph.CreateNode(r.GetOutput()[outputIdx])
			if err != nil {
				glog.Fatalf("something unexpected when outputing %s: %v", n.GetID(), err)
			}
			_, err = graph.CreateEdge("", current, outputNode)
			if err != nil {
				glog.Fatalf("something unexpected when connecting output %s: %v", n.GetID(), err)
			}
			outputIdx++
		}
		if prev != nil {
			_, err := graph.CreateEdge("", prev, current)
			if err != nil {
				glog.Fatalf("something unexpected when labeling %s: %v", n.GetID(), err)

			}
			//e.SetLabel(r.GetPath()[idx-1].GetOperator().GetText())
		}
		prev = current
	}

	err := gh.RenderFilename(graph, graphviz.PNG, filename)
	if err != nil {
		glog.Fatalf("something unexpected happened when saving %s: %v", filename, err)
	}
}
