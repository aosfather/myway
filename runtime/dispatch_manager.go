package runtime

import "github.com/aosfather/myway/meta"

//新增集群
func (this *DispatchManager) AddApplication(cluster *meta.ApplicationMapper) {
	if cluster != nil {
		c := this.clusterMap[cluster.App]
		if c == nil {
			this.clusterMap[cluster.App] = cluster
			//批量注册api
			apis := cluster.GetMappers()
			for _, api := range apis {
				this.AddApi("", "", api)
			}

		}
	}

}

//获取集群
func (this *DispatchManager) GetApplication(name string) *meta.ApplicationMapper {
	if name != "" {
		return this.clusterMap[name]
	}

	return nil
}

//删除集群
func (this *DispatchManager) DelApplication(name string) {
	delete(this.clusterMap, name)
}

//新增服务器
func (this *DispatchManager) AddServer(appName string, server *meta.Server) {
	if server != nil && appName != "" {
		c := this.clusterMap[appName]
		if c != nil {
			c.Cluster.AddServer(server)
		}
	}

}

//删除服务器
func (this *DispatchManager) DelServer(clusterName string, id int64) {
	if clusterName != "" && id != 0 {
		c := this.clusterMap[clusterName]
		if c != nil {
			var target []*meta.Server
			for _, v := range c.Cluster.Servers {
				if v.ID == id {
					continue
				}
				target = append(target, v)
			}

			c.Cluster.Servers = target

		}
	}
}

//新增api
func (this *DispatchManager) AddApi(domain, clusterName string, api *meta.ApiMapper) {
	if api == nil {
		return
	}

	if api.GetCluster() == nil && clusterName != "" {
		c := this.GetApplication(clusterName)
		c.AddMapper(api)

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
func (this *DispatchManager) DelApi(api *meta.ApiMapper) {
	if api == nil {
		return
	}

	//删除api的定义
	//路由映射,不用删除，因为在路由中只是存放了一个唯一的key
	delete(this.apiMap, api.Key())

}
