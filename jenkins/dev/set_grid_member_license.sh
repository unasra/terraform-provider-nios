#!/usr/bin/env bash
set -euo pipefail

# set_grid_member_license.sh
# Automates temp license activation and CA cert check disabling on a NIOS Grid Member.
# Activates licenses 11 and 12 only.
#
# Usage:
#   ./set_grid_member_license.sh <host_url>
#
# Example:
#   ./set_grid_member_license.sh 192.0.2.2  # Grid Member host/IP
#
# Requires: expect

if [ $# -ne 1 ]; then
    echo "Usage: $0 <host_url>" >&2
    echo "Example: $0 192.0.2.2  # Grid Member host/IP" >&2
    exit 1
fi

if ! command -v expect &> /dev/null; then
    echo "ERROR: 'expect' command not found. Please install expect." >&2
    exit 1
fi

HOST_URL="$1"
SSH_USER="admin"
LICENSE_NUMBERS=(11 12)

echo "Connecting to ${SSH_USER}@${HOST_URL} ..."
echo ""
echo "License numbers to activate: ${LICENSE_NUMBERS[*]}"
echo ""
echo "──────────────────────────────────────"

# Use expect to handle interactive prompts properly
expect -c "
set timeout 60
log_user 1

spawn ssh \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -o LogLevel=ERROR \
    \"${SSH_USER}@${HOST_URL}\"

expect {
    timeout { send_user \"\\nERROR: Timeout waiting for prompt\\n\"; exit 1 }
    eof { send_user \"\\nERROR: Connection closed unexpectedly\\n\"; exit 1 }
    -re {[>#] }
}

# Activate each license
foreach lic_num {${LICENSE_NUMBERS[*]}} {
    send_user \"\\n→ Activating license \$lic_num\\n\"
    send \"set temp_license\\r\"
    
    expect {
        timeout { send_user \"\\nERROR: Timeout waiting for license selection prompt\\n\"; exit 1 }
        \"Select license\" {
            send \"\$lic_num\\r\"
        }
    }
    
    # Handle any (y or n): prompts that may appear
    while {1} {
        expect {
            timeout { break }
            -re {\\(y or n\\):} {
                send \"y\\r\"
            }
            -re {[>#] } {
                break
            }
        }
    }
}

# Disable strict CA cert check
send_user \"\\n→ Disabling strict CA cert check\\n\"
send \"set disable_strict_ca_cert_check\\r\"

expect {
    timeout { }
    -re {\\(y or n\\):} {
        send \"y\\r\"
    }
}

expect -re {[>#] }

send \"exit\\r\"
expect eof
"

EXIT_CODE=$?
echo "──────────────────────────────────────"
echo ""

if [ $EXIT_CODE -eq 0 ]; then
    echo "✅ License activation and CA cert check disable completed successfully for ${HOST_URL}"
else
    echo "❌ Command execution failed with exit code ${EXIT_CODE}" >&2
    exit $EXIT_CODE
fi
