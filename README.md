# upstreams
一种通用的地址连接池方案：

serverpool:管理上游服务地址  
resolver:提供服务发现的能力，如集成：nacos    
loadbalancer:提供对上游服务地址进行负载均衡的能力，如：wrr     
connpool:管理地址对应的连接，每个地址下都可以是个独立的连接池，也可以复用grpc、net http 本身的连接池能力   

