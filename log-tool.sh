#!/bin/bash
# log-tool.sh: tmux + gum 選單版 Docker Compose log viewer
# 1. 左側 gum 選單選 service，右側 tmux pane 跑 log
# 2. 可直接複製 log，支援 Ctrl+C 快速退出
# 3. 自動偵測 docker-compose.local.yml 裡所有服務
trap 'rm -rf ./tmp' EXIT
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
SERVICES=($(docker compose -f docker-compose.local.yml config --services))
if [ ${#SERVICES[@]} -eq 0 ]; then
  echo -e "\033[1;31m[ERROR]\033[0m Can't find any service in docker-compose.local.yml"
  exit 1
fi

SESSION="logtool"
tmux kill-session -t "$SESSION" 2>/dev/null || true

mkdir -p ./tmp
echo "${SERVICES[0]}" > ./tmp/logtool_selected_service



# 左窗: gum 選單，選擇即寫入 TMP_SELECT
left_cmd='
TMP_SELECT="./tmp/logtool_selected_service"
while true; do \
  LAST_SELECTED=$(cat "$TMP_SELECT")
  SEL=$(gum choose --cursor="👉" --selected="${LAST_SELECTED}" '${SERVICES[@]}'); \
  if [ -n "$SEL" ]; then echo "$SEL" > "$TMP_SELECT"; fi; \
  sleep 0.2; \
done'

# 右窗: 監聽 TMP_SELECT，動態切換 log
right_cmd='LAST=""; while true; do \
  SEL="$(cat ./tmp/logtool_selected_service 2>/dev/null)"; \
  if [ "$SEL" != "$LAST" ] && [ -n "$SEL" ]; then \
    pkill -f "docker compose -f docker-compose.local.yml logs -f" 2>/dev/null || true; \
    clear; echo -e "\033[1;36m>>> $SEL log (Ctrl+C Exit)\033[0m"; \
    docker compose -f docker-compose.local.yml logs -f "$SEL" \
    | tee >(sed -r "s/\x1B\[[0-9;]*[mGKH]//g" >> "./tmp/${SEL}.log") & \
    LAST="$SEL"; \
  fi; \
  sleep 0.2; \
done'

# 啟動 tmux session

tmux new-session -d -s $SESSION bash -c "$left_cmd"
tmux split-window -v -p 70 -t $SESSION bash -c "$right_cmd"
tmux select-pane -t $SESSION:0.0
# 設置 Ctrl+C 快捷鍵來退出
tmux bind-key -n C-c kill-session
tmux attach-session -t $SESSION