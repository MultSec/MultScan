from flask import Flask
import logging
import flask.cli
import sys

# Disable logger to use our own
flask.cli.show_server_banner = lambda *args: None

log = logging.getLogger('werkzeug')
log.disabled = True

class Logger:
    @staticmethod
    def info(message):
        print(f"[\033[34m*\033[0m] {message}")

    @staticmethod
    def success(message):
        print(f"[\033[32m+\033[0m] {message}")

    @staticmethod
    def debug(message):
        print(f"[\033[33m^\033[0m] {message}")

    @staticmethod
    def error(message):
        print(f"[\033[31m!\033[0m] {message}")

    @staticmethod
    def section(message):
        print(f"\t[\033[93m-\033[0m] {message}")

    @staticmethod
    def subsection(message):
        print(f"\t\t[\033[95m>\033[0m] {message}")

# Init logger
Log = Logger()

app = Flask(__name__, static_folder='static')

from app import views