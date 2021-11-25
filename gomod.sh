#!/bin/bash
# 临时远程获取package，绕开公司网络
export GOPROXY=https://proxy.golang.org
go mod tidy
# 9865bfef540385e3d764d766071015bafd6ee7f90501ebd110611dae09cf