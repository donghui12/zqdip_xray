#!/bin/bash

# Define color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to display usage information
usage() {
  echo "Usage: $0 --server <server_file> --file <file_to_copy>"
  exit 1
}

# Function to print a message in green (success)
log_success() {
  echo -e "${GREEN}$1${NC}"
}

# Function to print a message in red (error)
log_error() {
  echo -e "${RED}$1${NC}"
}

# Function to print a message in yellow (info)
log_info() {
  echo -e "${YELLOW}$1${NC}"
}

# Function to validate required arguments
validate_args() {
  if [ -z "$server_file" ]; then
    log_error "Error: --server argument is required."
    usage
  fi

  if [ -z "$file_to_copy" ]; then
    log_error "Error: --file argument is required."
    usage
  fi

  if [ ! -f "$server_file" ]; then
    log_error "Error: File '$server_file' not found."
    exit 1
  fi

  if [ ! -f "$file_to_copy" ]; then
    log_error "Error: File '$file_to_copy' not found."
    exit 1
  fi
}

# Function to copy file to the server
copy_file_to_server() {
  local ip=$1
  local user=$2
  local password=$3
  local file_to_copy=$4

  log_info "[$ip] Copying file to server..."
  sshpass -p "$password" scp -o StrictHostKeyChecking=no "$file_to_copy" "$user@$ip:/tmp/"
  if [ $? -eq 0 ]; then
    log_success "[$ip] File copied successfully."
  else
    log_error "[$ip] Failed to copy file."
    return 1
  fi
}

# Function to run commands on the server
run_remote_commands() {
  local ip=$1
  local user=$2
  local password=$3
  local file_to_copy=$4

  log_info "[$ip] Connecting to server..."

  sshpass -p "$password" ssh -o StrictHostKeyChecking=no "$user@$ip" << EOF
    echo -e "${YELLOW}[$ip] Unzipping file...${NC}"
    cd /tmp/
    unzip -qo $(basename "$file_to_copy")   # Unzip quietly (no output)
    if [ \$? -eq 0 ]; then
      echo -e "${GREEN}[$ip] Unzipped successfully!${NC}"
    else
      echo -e "${RED}[$ip] Unzip failed!${NC}"
      exit 1
    fi

    echo -e "${YELLOW}[$ip] Preparing to run prepare.sh...${NC}"
    cd $(basename "$file_to_copy" .zip)    # Enter the unzipped directory
    chmod +x prepare.sh                        # Make prepare.sh executable
    echo -e "${YELLOW}[$ip] Executing prepare.sh...${NC}"
    # ./prepare.sh                               # Run the prepare.sh script
EOF

  if [ $? -eq 0 ]; then
    log_success "[$ip] Successfully executed prepare.sh."
    log_success "[$ip] Install successfully!!!"

    # Fetching quick_link.txt from zqpid_xray and appending to host.txt
    log_info "[$ip] Fetching quick_link.txt from zqpid_xray..."
    sshpass -p "$password" scp -o StrictHostKeyChecking=no "$user@$ip:/tmp/$(basename "$file_to_copy" .zip)/zqpid_xray/quick_link.txt" ./quick_link_$ip.txt
    if [ $? -eq 0 ]; then
      log_success "[$ip] Successfully fetched quick_link.txt."
      cat quick_link_$ip.txt >> host.txt
      log_success "[$ip] Appended quick_link.txt to host.txt."
      rm quick_link_$ip.txt  # Clean up
    else
      log_error "[$ip] Failed to fetch quick_link.txt."
    fi
  else
    log_error "[$ip] Failed to execute prepare.sh."
    log_error "[$ip] Install Failed!!!"
  fi
}

# Main function to process each server
process_servers() {
  while read -r line; do
    ip=$(echo $line | awk '{print $1}')
    user=$(echo $line | awk '{print $2}')
    password=$(echo $line | awk '{print $3}')

    copy_file_to_server "$ip" "$user" "$password" "$file_to_copy"
    if [ $? -eq 0 ]; then
      run_remote_commands "$ip" "$user" "$password" "$file_to_copy"
    fi

  done < "$server_file"
}

# Parse command-line arguments
while [[ "$#" -gt 0 ]]; do
  case $1 in
    --server) server_file="$2"; shift ;;
    --file) file_to_copy="$2"; shift ;;
    *) log_error "Unknown parameter: $1"; usage ;;
  esac
  shift
done

# Validate arguments
validate_args

# Process servers from the list
process_servers