#!/bin/bash
# log-tool.sh: tmux + gum 選單版 Docker Compose log viewer
# 1. 左側 gum 選單選 service，右側 tmux pane 跑 log
# 2. 可直接複製 log，支援 Ctrl+C 快速退出
# 3. 自動偵測 docker-compose.dev.yml 裡所有服務

set -e

if ! command -v gum >/dev/null 2>&1; then
  echo -e "\033[1;31m[ERROR]\033[0m gum 未安裝，請先執行: brew install gum"
  exit 1
fi
if ! command -v tmux >/dev/null 2>&1; then
  echo -e "\033[1;31m[ERROR]\033[0m tmux 未安裝，請先執行: brew install tmux"
  exit 1
fi

# 取得所有 docker compose service 名稱
SERVICES=($(docker compose -f docker-compose.dev.yml config --services))
if [ ${#SERVICES[@]} -eq 0 ]; then
  echo -e "\033[1;31m[ERROR]\033[0m 找不到任何服務，請確認 docker-compose.dev.yml"
  exit 1
fi

LOG_DIR="logs"
TMP_SELECT=".service_select.tmp"
mkdir -p "$LOG_DIR"

echo "${SERVICES[0]}" > "$TMP_SELECT"

SESSION="logtool"

# 關閉舊 session
tmux kill-session -t "$SESSION" 2>/dev/null || true

# 左窗: gum 選單 (無限 loop, 選擇寫入 tmpfile)
left_cmd='
echo " Ctrl+C Exit the app"; \
while true; do \
  SEL=$(gum choose --cursor-prefix "👉" --selected.foreground 212 '${SERVICES[@]}'); \
  if [ -n "$SEL" ]; then echo "$SEL" > "$TMP_SELECT"; fi; \
  sleep 0.5; \
done'

# 右窗: 根據 tmpfile 選擇動態切換 log（監聽檔案變動）
right_cmd='LAST=""; while true; do \
  SEL="$(cat "$tmpfile" 2>/dev/null)"; \
  if [ "$SEL" != "$LAST" ] && [ -n "$SEL" ]; then \
    pkill -f "docker compose -f docker-compose.dev.yml logs -f" 2>/dev/null || true; \
    echo "\033[1;36m>>> $SEL log (Ctrl+C 停止, 左側可切換)\033[0m"; \
    docker compose -f docker-compose.dev.yml logs -f "$SEL" & \
    LAST="$SEL"; \
  fi; \
  sleep 0.5; \
done'

# 啟動 tmux session

tmux new-session -d -s "$SESSION" bash -c "$left_cmd"
tmux split-window -h -p 70 -t "$SESSION" bash -c "$right_cmd"
tmux select-pane -t 0
tmux select-layout -t "$SESSION" even-horizontal
tmux attach -t "$SESSION"
