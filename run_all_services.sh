#!/bin/bash

# run_all_services.sh (tmux+gum 左選單右 log 動態切換)
# 1. 啟動 user, admin, creator 三個 service（air 背景，log 分開）
# 2. tmux 左窗 gum 選單，右窗動態 tail log
# 3. 支援 -e 參數指定 ENV，預設 local
# 4. 離開自動清理

set -e

if ! command -v gum >/dev/null 2>&1; then
  echo -e "\033[1;31m[ERROR]\033[0m gum 未安裝，請先執行: brew install gum"
  exit 1
fi
if ! command -v tmux >/dev/null 2>&1; then
  echo -e "\033[1;31m[ERROR]\033[0m tmux 未安裝，請先執行: brew install tmux"
  exit 1
fi

while getopts "e:" opt; do
  case $opt in
    e) ENV="$OPTARG" ;;
    *) ENV="local" ;;
  esac
done
ENV=${ENV:-local}

LOG_DIR="logs"
SERVICES=("user" "admin" "creator")
PIS=()
TMP_SELECT=".service_select.tmp"

mkdir -p "$LOG_DIR"

# 若 local 環境，自動啟動資料庫 container
if [[ "$ENV" == "local" ]]; then
  echo "啟動本地資料庫（docker-compose.db.yaml）..."
  docker-compose -f docker-compose.db.yaml up -d
fi

# 啟動所有服務（air 背景執行，log 分開）
echo "啟動三個服務..."
for SVC in "${SERVICES[@]}"; do
  SVC_DIR="services/$SVC"
  LOG_FILE="$LOG_DIR/$SVC.$ENV.log"
  if [ ! -d "$SVC_DIR" ]; then
    echo "$SVC_DIR 不存在，略過。"
    continue
  fi
  (
    cd "$SVC_DIR"
    APP_ENV="$ENV" air > "../../$LOG_FILE" 2>&1 &
    echo $! > "../../$LOG_DIR/$SVC.$ENV.pid"
  )
  PIDS+=("$(cat "$LOG_DIR/$SVC.$ENV.pid")")
done
sleep 2

echo "user" > "$TMP_SELECT"

# 右窗格 tail log，根據暫存檔動態切換
right_pane_cmd='\
  LOG_DIR="logs"; ENV="'$ENV'"; TMP_SELECT=".service_select.tmp"; \
  LAST=""; TAIL_PID=""; \
  while true; do \
    SVC=$(cat "$TMP_SELECT" 2>/dev/null); \
    LOG_FILE="$LOG_DIR/$SVC.$ENV.log"; \
    if [ "$SVC" != "$LAST" ]; then \
      if [ -n "$TAIL_PID" ]; then kill $TAIL_PID 2>/dev/null; fi; \
      clear; \
      echo -e "\033[1;36m@magnetar/$SVC-api#dev\033[0m"; \
      echo -e "\033[0;33mtail -f $LOG_FILE\033[0m"; \
      echo "--------------------------------------"; \
      if [ -f "$LOG_FILE" ]; then \
        tail -n 20 -f "$LOG_FILE" & \
        TAIL_PID=$!; \
      else \
        echo -e "\033[1;31m找不到 log 檔案: $LOG_FILE\033[0m"; \
        TAIL_PID=""; \
      fi; \
      LAST="$SVC"; \
    fi; \
    sleep 0.5; \
  done'

# 左窗格 gum 選單，選擇後寫入暫存檔
left_pane_cmd='\
  SERVICES=(user admin creator); TMP_SELECT=".service_select.tmp"; \
  echo -e "\033[1;33m按 Ctrl+C 可快速退出並停止所有服務\033[0m"; \
  while true; do \
    SVC=$(gum choose --cursor-prefix "👉 " --selected.foreground 212 "${SERVICES[@]}"); \
    if [ -n "$SVC" ]; then echo "$SVC" > "$TMP_SELECT"; fi; \
    sleep 0.2; \
  done'

SESSION="night_debug"
tmux kill-session -t $SESSION 2>/dev/null || true
# 直接讓左側 pane 只跑 gum 選單，右側只跑 log watcher
# 這樣不會有 motd、shell prompt 等雜訊

# 設置快速退出的函數
quick_exit() {
  tmux kill-session -t $SESSION 2>/dev/null || true
  cleanup
  exit 0
}

tmux new-session -d -s $SESSION "bash -c '$left_pane_cmd'"
tmux split-window -h -p 70 -t $SESSION "bash -c '$right_pane_cmd'"
tmux select-pane -t $SESSION:0.0

# 設置 Ctrl+C 快捷鍵來退出
tmux bind-key -n C-c kill-session

echo "\033[1;32m請在 tmux 視窗操作，左側選單切換服務，右側即時 log。\033[0m"
echo "\033[1;33m按 Ctrl+C 可快速退出並停止所有服務\033[0m"
tmux attach-session -t $SESSION

# 離開時清理
cleanup() {
  rm -f "$TMP_SELECT"
  # 直接結束所有服務，不再詢問
  for SVC in "${SERVICES[@]}"; do
    PID_FILE="$LOG_DIR/$SVC.$ENV.pid"
    if [ -f "$PID_FILE" ]; then
      PID=$(cat "$PID_FILE")
      if kill -0 $PID 2>/dev/null; then
        kill $PID 2>/dev/null || true
        sleep 1
        if kill -0 $PID 2>/dev/null; then
          kill -9 $PID 2>/dev/null || true
        fi
      fi
      rm "$PID_FILE"
    fi
  done
  echo "三個 service 已結束。"
  clear
  echo "Bye!"
}

trap cleanup EXIT