# kubelet requires swap off
swapoff -a
# keep swap off after reboot
sed -i '/swap/ s/^\(.*\)$/#\1/g' /etc/fstab

# Add br_netfilter module to get rid of "reply from unexpected source" in 12.
# Deploying the DNS Cluster Add-on. Not sure about details.
#
# https://github.com/kubernetes/kubernetes/issues/21613#issuecomment-343190401
modprobe br_netfilter
echo 'br_netfilter' >> /etc/modules
