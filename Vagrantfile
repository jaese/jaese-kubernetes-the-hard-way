Vagrant.configure("2") do |config|
  config.vm.box = "bento/ubuntu-20.04"

  config.vm.provider :virtualbox do |v|
    v.memory = "512"
  end

  config.vm.define "controller-0" do |v|
    v.vm.hostname = "controller-0"
    v.vm.network :private_network, ip: "192.168.60.20"
  end

  # pod-cidr=10.200.0.0/24
  (0..1).each do |n|
    config.vm.define "worker-#{n}" do |v|
      v.vm.hostname = "worker-#{n}"
      v.vm.network :private_network, ip: "192.168.60.3#{n}"
    end
  end
end
