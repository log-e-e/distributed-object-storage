# 生成哈希值
content=$1
# shellcheck disable=SC2006
echo "Start generate `${content}` hash code..."
echo -n "$content" | openssl dgst -sha1 -binary | base64