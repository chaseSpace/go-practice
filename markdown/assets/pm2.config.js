// 在开发、测试环境可以使用pm2来编排服务（替代k8s）
// note：但失去了容器化的优势（只能管理直接运行在裸机上的服务，所以无法多副本运行指定服务）
// 安装pm2: https://developer.aliyun.com/article/906699

// 1. 编辑此文件来添加服务
// 2. make deploy SVC=user ENV=prod （一键部署新服务）
// 3. make up SVC=user ENV=prod （一键更新服务）

// @其他命令
// pm2 start pm2.config.js --time --only "gateway" --env beta
// pm2 reload <服务名> -i x 重启时调整服务实例数。与restart区别：单实例数时没区别。多实例时，reload会逐个重启每个进程（至少1个可用），restart会同时重启所有进程
// pm2 kill [进程名|id] 杀死所有pm2服务
// pm2 del 进程名
// pm2 startup 开机启动pm2
// pm2 ls 查看进程列表
// pm2 scale <服务名> <数量> 调整实例数
// pm2 monit 进入监控ui
// pm2 ecosystem 创建pm2配置文件模板
// pm2 save 保存配置到本地，方便在机器重启时恢复

// @调试常用
// pm2 env <pm-id> 查看指定服务id的环境变量，调试使用
// pm2 attach <pm-id> 进入指定服务id的stdin/stdout

// @其他选项
// --max-memory-restart 100M

// @日志相关
// pm2 logs admin --lines 20
// pm2 logs admin --err 只查看err日志
// pm2 logs admin --out 只查看out（正常）日志
// pm2 flush [进程名|id] 删除（进程）全部日志
// rm -rf ~/.pm2/logs/gateway-out__*.log  仅删除gateway服务的out日志(不含正在写入的文件)
// rm -rf ~/.pm2/logs/gateway-error*.log  仅删除gateway服务的error日志(不含正在写入的文件)

// @滚动日志
// pm2 install pm2-logrotate
// pm2 conf pm2-logratate  查看配置
// pm2 set pm2-logrotate:max_size 50M
// pm2 set pm2-logrotate:retain 7
// pm2 set pm2-logrotate:compress false
// pm2 set pm2-logrotate:dateFormat YYYY-MM-DD_HH-mm-ss
// pm2 set pm2-logrotate:workerInterval 30
// pm2 set pm2-logrotate:rotateInterval 0 0 * * *   定时强制分割日志
// pm2 set pm2-logrotate:rotateModule true   分割 pm2 本身的日志文件


// 下面是一个微服务配置文件示例
module.exports = {
    apps: [{
        name: "gateway",
        script: ".bin/gateway",
        instances: 1, // 若服务监听固定端口，单机环境下实例数只能是1
        log_date_format: "YYYY-MM-DD HH:mm:ss", // 设置这个后，输出的所有日志行会加上时间前缀（重要！）
        max_memory_restart: "100M",
        env_beta: { // 格式固定：env_$ENV
            MICRO_SVC_ENV: "beta",
            MICRO_SVC_NO_PRINT_CFG: false,
            MICRO_SVC_LOG_LEVEL: "debug"
        },
        env_prod: {
            MICRO_SVC_ENV: "prod",
            MICRO_SVC_NO_PRINT_CFG: false,
            MICRO_SVC_LOG_LEVEL: "info"
        },
    },
        {
            name: "user",
            script: ".bin/user",
            instances: 1,
            log_date_format: "YYYY-MM-DD HH:mm:ss",
            env_beta: {
                MICRO_SVC_ENV: "beta",
                MICRO_SVC_NO_PRINT_CFG: false,
                MICRO_SVC_LOG_LEVEL: "debug"
            },
            env_prod: {
                MICRO_SVC_ENV: "prod",
                MICRO_SVC_NO_PRINT_CFG: false,
                MICRO_SVC_LOG_LEVEL: "info"
            },
        },
    ]
}
