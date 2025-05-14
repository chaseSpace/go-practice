import logging
import sys
from logging.config import dictConfig

LOGGING_CONFIG = {
    "version": 1,
    "disable_existing_loggers": False,
    "formatters": {
        "default": {
            "()": "uvicorn.logging.DefaultFormatter",
            "fmt": "%(asctime)s.%(msecs)03d %(threadName)s %(levelprefix)s  %(message)s",
            "use_colors": True,
            "datefmt": "%Y-%m-%d %H:%M:%S"
        },
        "access": {
            "()": "uvicorn.logging.AccessFormatter",
            "fmt": "%(asctime)s.%(msecs)03d %(threadName)s %(levelprefix)s %(client_addr)s - '%(request_line)s' %(status_code)s",
            "use_colors": True,
            "datefmt": "%Y-%m-%d %H:%M:%S"
        }
    },
    "handlers": {
        "default": {
            "formatter": "default",
            "class": "logging.StreamHandler",
            "stream": sys.stdout,
        },
        # default 日志输出到本地文件，定义了日志路径、备份数量
        "default_file": {
            "formatter": "default",
            "class": "logging.handlers.TimedRotatingFileHandler",
            "filename": ".logs/default.log",
            "when": "midnight",
            "encoding": "utf-8",
            "backupCount": 5
        },
        "access": {
            "formatter": "access",
            "class": "logging.StreamHandler",
            "stream": sys.stdout,
        },
        # access 日志输出到本地文件，定义了日志路径、备份数量
        "access_file": {
            "formatter": "access",
            "class": "logging.handlers.TimedRotatingFileHandler",
            "filename": ".logs/access.log",
            "when": "midnight",
            "encoding": "utf-8",
            "backupCount": 5
        },
        'console': {
            'class': 'logging.StreamHandler',
            'level': 'INFO',
            'formatter': 'default',
            "stream": sys.stdout,
        },
        'server_file': {
            'class': 'logging.handlers.TimedRotatingFileHandler',
            'level': 'INFO',
            'formatter': 'default',
            'filename': '.logs/our_app.log',
            "when": "midnight",
            "encoding": "utf-8",
            'backupCount': 5,
        },
    },
    "loggers": {
        "uvicorn": {"handlers": ["default", "default_file"], "level": "INFO", "propagate": False},
        "uvicorn.error": {"level": "INFO"},
        "uvicorn.access": {"handlers": ["access", "access_file"], "level": "INFO", "propagate": False},

        "our_app": {"level": "INFO", 'handlers': ['console', 'server_file']},
    }
}

logger = logging.getLogger("our_app")


def configure_logging():
    """配置应用程序日志"""
    dictConfig(LOGGING_CONFIG)
