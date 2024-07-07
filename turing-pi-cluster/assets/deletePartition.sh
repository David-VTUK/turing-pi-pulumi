#!/bin/bash

# Define the file path
AGENT_UNINSTALL="/usr/local/bin/k3s-agent-uninstall.sh"
SERVER_UNINSTALL="/usr/local/bin/k3s-uninstall.sh"


# Check if the agent is still running, if so, delete
if [ -f "$AGENT_UNINSTALL" ]; then
    # If the file exists, execute it
    "$AGENT_UNINSTALL"
fi

# Check if the server is still running, if so, delete
if [ -f "$SERVER_UNINSTALL" ]; then
    # If the file exists, execute it
    "$SERVER_UNINSTALL"
fi

sleep 30

# Define variables
DISK="/dev/nvme0n1"
PARTITION="${DISK}p1"
MOUNT_POINT="/mnt/data"

# Check if the partition exists
if [ -b "$PARTITION" ]; then
    echo "Partition $PARTITION exists. Proceeding to delete the partition."

    # Unmount the partition if it is mounted
    if mountpoint -q $MOUNT_POINT; then
        sudo umount $MOUNT_POINT
        echo "Unmounted $MOUNT_POINT."
    fi

    # Delete the partition using fdisk
    echo -e "d\nw" | sudo fdisk $DISK

    # Inform the user
    echo "Partition $PARTITION has been deleted."
else
    echo "Partition $PARTITION does not exist. Exiting script."
    exit 0
fi

# Remove the mount point directory if it exists
if [ -d "$MOUNT_POINT" ]; then
    sudo rmdir $MOUNT_POINT
    echo "Removed mount point directory $MOUNT_POINT."
fi

# Backup /etc/fstab
sudo cp /etc/fstab /etc/fstab.bak

# Remove the entry from /etc/fstab if it exists
UUID=$(sudo blkid -s UUID -o value "$PARTITION")
if [ -n "$UUID" ]; then
    sudo sed -i.bak "/UUID=$UUID/d" /etc/fstab
    echo "Removed $PARTITION entry from /etc/fstab."
fi

echo "Removed Partition successfully."