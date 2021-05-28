#!/bin/bash
strA=$(ifconfig)
nicA="rdma0"
nicB="rdma1"
if [[ $strA =~ $nicA ]]; then
	export GLOO_SOCKET_IFNAME=$nicA
else
	export GLOO_SOCKET_IFNAME=$nicB
fi
echo $GLOO_SOCKET_IFNAME
