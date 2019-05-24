package runtime

import "github.com/aosfather/myway/meta"

//新增集群
func (this *DispatchManager) AddCluster(cluster *meta.ServerCluster) {
	if cluster != nil {
		c := this.clusterMap[cluster.ID]
		if c == nil {
			this.clusterMap[cluster.ID] = cluster
		}
	}

}

//获取集群
func (this *DispatchManager) GetCluster(name string) *meta.ServerCluster {
	if name != "" {
		return this.clusterMap[name]
	}

	return nil
}

//删除集群
func (this *DispatchManager) DelCluster(name string) {
	delete(this.clusterMap, name)
}

//新增服务器
func (this *DispatchManager) AddServer(clusterName string, server *meta.Server) {
	if server != nil && clusterName != "" {
		c := this.clusterMap[clusterName]
		if c != nil {
			c.Servers = append(c.Servers, server)
		}
	}

}

//删除服务器
func (this *DispatchManager) DelServer(clusterName string, id int64) {
	if clusterName != "" && id != 0 {
		c := this.clusterMap[clusterName]
		if c != nil {
			var target []*meta.Server
			for _, v := range c.Servers {
				if v.ID == id {
					continue
				}
				target = append(target, v)
			}

			c.Servers = target

		}
	}
}

//新增api
func (this *DispatchManager) AddApi(domain, clusterName string, api *meta.Api) {
	if api == nil {
		return
	}

	if api.Cluster == nil {
		api.Cluster = this.GetCluster(clusterName)
	}

	this.apiMap[api.Key()] = api
	var apiNode *node
	if domain == "" {
		apiNode = this.defaultNode
	} else {
		if domain == this.env.Domain {
			apiNode = this.defaultNode
		}

		//处理不同的域名的映射
		if apiNode == nil {
			apiNode = this.domainNode[domain]
			if apiNode == nil {
				apiNode = &node{}
				this.domainNode[domain] = apiNode
			}
		}

	}

	if api != nil {
		apiNode.addRoute(api.Url, api.Key())
	}
}

//删除Api
func (this *DispatchManager) DelApi(api *meta.Api) {
	if api == nil {
		return
	}

	//删除api的定义
	//路由映射,不用删除，因为在路由中只是存放了一个唯一的key
	delete(this.apiMap, api.Key())

}
