#!/bin/bash

# Define variables
DISK="/dev/nvme0n1"
PARTITION="${DISK}p1"
MOUNT_POINT="/mnt/data"

# Check if the disk is already partitioned
if [ -b "$PARTITION" ]; then
    echo "Partition $PARTITION already exists. Exiting script."
    exit 0
else
    # Create a new partition
    echo -e "n\np\n1\n\n\nw" | sudo fdisk $DISK
fi

# Create EXT4 filesystem on the partition
sudo mkfs.ext4 $PARTITION

# Create mount point if it doesn't exist
if [ ! -d "$MOUNT_POINT" ]; then
    sudo mkdir -p $MOUNT_POINT
fi

# Mount the filesystem
sudo mount $PARTITION $MOUNT_POINT

# Get the UUID of the partition
UUID=$(sudo blkid -s UUID -o value $PARTITION)

# Backup /etc/fstab
sudo cp /etc/fstab /etc/fstab.bak

# Add entry to /etc/fstab if it doesn't already exist
if ! grep -q "$UUID" /etc/fstab; then
    echo "UUID=$UUID  $MOUNT_POINT  ext4  defaults  0  2" | sudo tee -a /etc/fstab
fi

# Verify the fstab entry by mounting all filesystems
sudo mount -a

echo "EXT4 filesystem created and mounted on $MOUNT_POINT. Entry added to /etc/fstab."