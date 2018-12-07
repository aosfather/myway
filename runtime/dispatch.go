package runtime

import (
	"github.com/aosfather/myway/meta"
	"strings"
)

/**
  dispatch
*/

type DispatchManager struct {
	domainNode  map[string]*node               //特定域名下的node
	defaultNode *node                          //默认
	clusterMap  map[string]*meta.ServerCluster //集群列表
	apiMap      map[string]*meta.Api           //api列表
	env         meta.Env                       //环境列表
}

func (this *DispatchManager) Init() {
	this.domainNode = make(map[string]*node)
	this.clusterMap = make(map[string]*meta.ServerCluster)
	this.apiMap = make(map[string]*meta.Api)
	this.defaultNode = &node{}
}

//根据域名和url获取对应的API
func (this *DispatchManager) GetApi(domain, url string) *meta.Api {
	node := this.domainNode[domain]
	if node == nil {
		node = this.defaultNode
	}

	if node != nil {
		paramIndex := strings.Index(url, "?")
		realuri := url
		if paramIndex != -1 {
			realuri = strings.TrimSpace((url[:paramIndex]))
		}

		h, _, _ := node.getValue(realuri)
		if h != nil {
			key := h.(string)
			return this.apiMap[key]
		}

	}
	return nil

}
