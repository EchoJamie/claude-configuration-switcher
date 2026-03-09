# CCS (Claude Configuration Switcher)

CCS 是一个用于管理 Claude API 配置的 CLI 工具，支持多配置切换。

## 安装

```bash
make build
```

## 配置文件

- **全局配置**: `~/.ccs/config.json` - 存储所有可用的 Claude 配置
- **目标文件**: 命令执行目录的 `.claude/settings.local.json` 文件中的 env 配置项

### 配置格式

全局配置存储格式：
```json
{
  "profiles": {
    "default": {
      "name": "default",
      "ANTHROPIC_API_KEY": "sk-XXXXXXXXXXXXXXXXXXXXXX",
      "ANTHROPIC_BASE_URL": "https://api.anthropic.com",
      "ANTHROPIC_MODEL": "opus",
      "ANTHROPIC_DEFAULT_OPUS_MODEL": "opus-4-5-20250514",
      "ANTHROPIC_DEFAULT_SONNET_MODEL": "sonnet-4-20250514",
      "ANTHROPIC_DEFAULT_HAIKU_MODEL": "haiku-3-20250514",
      "ANTHROPIC_SMALL_FAST_MODEL": "haiku-3-20250514"
    }
  }
}
```

## 支持的命令

```bash
ccs init                                    # 配置文件初始化
ccs add                                     # 交互式添加新配置
ccs list                                    # 列出所有可用配置
ccs use <name>                              # 切换到指定配置
ccs show                                    # 显示当前使用的配置名称
ccs clear                                   # 清理当前配置
ccs rename <old> <new>                      # 重命名配置
ccs rm <name>                               # 删除指定配置
ccs import [path]                           # 从 cc-switch 导入配置
ccs edit                                    # 使用默认编辑器打开配置文件
```

## 命令详细说明

### ccs init
- 初始化全局配置文件 `~/.ccs/config.json`
- 如果文件已存在，则跳过

### ccs add
交互式添加新配置，需要输入：
- 配置名称
- API Key（必填）
- Base URL（必填）
- 默认模型
- Opus/Sonnet/Haiku/快速模型（可选）

### ccs list
列出所有可用的配置名称

### ccs use \<name\>
切换到指定配置，将配置写入 `.claude/settings.local.json` 的 env 配置项

### ccs show
显示当前使用的配置名称

### ccs clear
清理当前配置，删除 `.claude/settings.local.json` 中的 env 字段

### ccs rename \<old\> \<new\>
重命名配置

### ccs rm \<name\>
删除指定配置（需要确认）

### ccs import [path]
从 cc-switch 数据库导入配置
- 默认路径: `~/.cc-switch/cc-switch.db`
- 支持指定自定义路径
- 多选导入，自动检测缺失字段

### ccs edit
使用默认编辑器打开全局配置文件

## 使用示例

```bash
# 初始化
ccs init

# 添加配置
ccs add

# 列出配置
ccs list

# 切换配置
ccs use default

# 查看当前配置
ccs show

# 清理配置
ccs clear
```