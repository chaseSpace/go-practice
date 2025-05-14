# 秒数格式化
# e.g. 1 -> 1s; 20 -> 20s; 60 -> 1m; 61 -> 1m1s; 3661 -> 1h1m1s
def seconds_to_readable(seconds):
    minutes, seconds = divmod(seconds, 60)
    hours, minutes = divmod(minutes, 60)
    days, hours = divmod(hours, 24)

    parts = []
    if days:
        parts.append(f"{int(days)}d")
    if hours:
        parts.append(f"{int(hours)}h")
    if minutes:
        parts.append(f"{int(minutes)}m")
    if seconds:
        parts.append(f"{int(seconds)}s")

    return ''.join(parts) or "0s"
