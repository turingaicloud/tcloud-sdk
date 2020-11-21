#!/bin/bash
strA=`ifconfig`
nicA="enp1s0f1"
nicB="enp129s0f1"
if [[ $strA =~ $nicA ]]
then
	export GLOO_SOCKET_IFNAME=$nicA
else
	export GLOO_SOCKET_IFNAME=$nicB
fi
echo $GLOO_SOCKET_IFNAME
