Vagrant.configure("2") do |config|
  config.vm.box = "bento/ubuntu-20.04"

  config.vm.provider :virtualbox do |v|
    v.memory = "512"
  end

  config.vm.define "controller-0" do |v|
    v.vm.hostname = "controller-0"
    v.vm.network :private_network, ip: "192.168.60.20"
  end

  (0..2).each do |n|
    config.vm.define "worker-#{n}" do |v|
      v.vm.hostname = "worker-#{n}"
      v.vm.network :private_network, ip: "192.168.60.3#{n}"
      v.vm.provision "shell", path: "setup-worker.sh"

      # Provision pod network routes.
      (0..2).each do |route_to|
        if route_to != n
          v.vm.provision "shell",
            run: "always",
            inline: "ip route replace 10.200.#{route_to}.0/24 via 192.168.60.3#{route_to}"
        end
      end
    end
  end
end
