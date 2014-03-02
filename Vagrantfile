# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant::Config.run do |config|
  config.vm.box = "base-0.4.0"
  config.vm.box_url = "http://paasta-boxes.s3.amazonaws.com/base-0.4.0-amd64-20131218-virtualbox.box"
  config.vm.customize ["modifyvm", :id, "--cpus", 2 ]
  config.vm.provision :shell,
    :inline => <<SCRIPT
PACKAGE=go1.2.linux-amd64.tar.gz
if ! [ -d /usr/local/go ]; then                                          
  cd /usr/local
  wget -nv http://go.googlecode.com/files/$PACKAGE
  tar xzf $PACKAGE
  echo 'export PATH="/usr/local/go/bin:$PATH"' > /etc/profile.d/golang.sh
fi 
SCRIPT
end

