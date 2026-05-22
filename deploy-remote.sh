#!/bin/bash
set -e

echo "==> 1. 部署配置文件"
sudo mkdir -p /opt/go-demo/configs
sudo mv /tmp/config.yaml /opt/go-demo/configs/config.yaml
sudo chown -R godemo:godemo /opt/go-demo/configs

echo "==> 2. 替换二进制"
sudo systemctl stop go-demo
sudo mv /tmp/go-demo.new /opt/go-demo/bin/go-demo
sudo chown godemo:godemo /opt/go-demo/bin/go-demo
sudo chmod +x /opt/go-demo/bin/go-demo

echo "==> 3. 启动服务"
sudo systemctl start go-demo
sleep 2

echo "==> 4. 服务状态"
sudo systemctl is-active go-demo

echo "==> 服务日志(最近 10 行)"
sudo journalctl -u go-demo -n 10 --no-pager

echo
echo "==> 5. 老接口验证"
curl -s http://127.0.0.1:8080/api/claim/1
echo
echo

echo "==> 6. 新接口验证(读 jar_metrics)"
curl -s http://127.0.0.1:8080/api/v1/jars/1/metrics/latest
echo
echo

echo "==> 7. 全局窖藏环境"
curl -s http://127.0.0.1:8080/api/v1/home/cellar-env
echo

echo "==> 完成"
