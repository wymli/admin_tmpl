# 安装just:
# 

install:
    # 1. install just: curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/bin
    # 2. install gvm:
    # 3. install nvm: 
    # 4: install air: go install github.com/air-verse/air@latest


check:
    #!/usr/bin/env bash
    [[ $(go version) =~ "go1.23.0" ]] || (echo "must go 1.23.0" && exit -1)
    [[ $(node -v) =~ "v22.13.1" ]] || (echo "must node 22.13.1" && exit -1)


dev $type="split":
    #!/usr/bin/env bash
    if [[ $type == "split" ]];then
        # 分别启动前后端的调试模式
        echo "start split dev..."
        cd fe
        source $NVM_DIR/nvm.sh
        nvm use
        NODE_ENV=x npm run dev 2>&1 & # 
        cd -
        sleep 1s && echo "fe dev..."
        air 2>&1 &
        sleep 1s && echo "be dev..."
        echo "fe & be dev..."
        wait
    elif [[ $type == "join" ]]; then
        # 编译前端后，只启动后端的调试模式，前端的dist目录被embed
        cd fe && npm run build && cd -
        air 2>&1 &
        sleep 1s && echo "be dev..."
        wait
    elif [[ $type == "be" ]];then
        # 只启动后端的调试模式，前端的dist目录被embed
        air 2>&1 &
        sleep 1s && echo "be dev..."
        wait
    fi
