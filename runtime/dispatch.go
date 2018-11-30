package runtime

import (
	"github.com/aosfather/myway/meta"
)

/**
  dispatch
*/

type DispatchManager struct {
	domainNode  map[string]*node               //特定域名下的node
	defaultNode *node                          //默认
	clusterMap  map[string]*meta.ServerCluster //集群列表
	apiMap      map[string]*meta.Api           //api列表
}

//根据域名和url获取对应的API
func (this *DispatchManager) GetApi(domain, url string) *meta.Api {

	return nil
}
