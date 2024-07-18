from flask import Flask, jsonify, send_from_directory
from werkzeug.exceptions import NotFound
from dotenv import load_dotenv
import os
import socket
import winrm
from wakeonlan import send_magic_packet

load_dotenv()
WINRM_HOST = os.getenv("WINRM_HOST")
WINRM_PORT = 5986
PSSHUTDOWN_PATH = os.getenv("PSSHUTDOWN_PATH")
WOL_MAC = os.getenv("WOL_MAC")

app = Flask(__name__)

# todo: replace with nginx rule, only send api to wsgi
@app.route("/", defaults={"path": ""})
@app.route("/<path:path>")
def serve_static(path: str):
	try:
		return send_from_directory(app.static_folder, path)
	except NotFound as _:
		# if path.endswith("/"):
		return send_from_directory(app.static_folder, path + "/index.html")
		# raise e

@app.route("/api/ping")
def ping():
	s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
	s.settimeout(2)
	try:
		s.connect((WINRM_HOST, WINRM_PORT))
		s.shutdown(socket.SHUT_RDWR)
		return jsonify(True)
	except Exception as _:
		return jsonify(False)
	finally:
		s.close()

# https+ntlm
# session = winrm.Session(f"https://{WINRM_HOST}:{WINRM_PORT}/wsman",
# 	auth=(os.getenv("WINRM_USER"), os.getenv("WINRM_PASSWORD")),
# 	transport="ntlm",
# 	server_cert_validation="ignore"\
# )

# https+ssl (no user credentials sent)
session = winrm.Session(f"https://{WINRM_HOST}:{WINRM_PORT}/wsman",
	# auth=(os.getenv("WINRM_USER"), os.getenv("WINRM_PASSWORD")),
	auth=(None, None),
	transport="ssl",
	server_cert_validation="ignore",
	cert_pem="./user.pem", cert_key_pem="./key.pem"
)

@app.route("/api/boot-time")
def boot_time():
	result = session.run_ps("(Get-CimInstance -ClassName Win32_OperatingSystem).LastBootUpTime")
	if result.status_code != 0: return jsonify({"error": result.std_err.decode()})
	
	boot_time = result.std_out.decode().strip()
	return jsonify({"boot_time": boot_time})

@app.route("/api/sleep")
def sleep():
	result = session.run_cmd(f"{PSSHUTDOWN_PATH}  -d -t 5 -c")
	if result.status_code != 0: return jsonify({"error": result.std_err.decode()})
	return jsonify({"message": result.std_out.decode().strip().split("\n")[-1]})

@app.route("/api/wake")
def wake():
	send_magic_packet(WOL_MAC)

if __name__ == "__main__":
	app.run(host="0.0.0.0", port=5000, debug=True)

# # result = session.run_cmd("whoami", [])

# ps_script = """
# Start-Sleep -Seconds 5
# Add-Type -AssemblyName System.Windows.Forms
# [System.Windows.Forms.Application]::SetSuspendState('Suspend', $false, $false)
# """
# result = session.run_ps(ps_script)
# print(result.status_code)
# print(result.std_out.decode())