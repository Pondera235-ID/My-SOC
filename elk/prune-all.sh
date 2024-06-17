#!/bin/bash

echo "Pruning all Docker data..."

# Remove all stopped containers
echo "Removing all stopped containers..."
sudo docker container prune -f

# Remove all unused images
echo "Removing all unused images..."
#sudo docker image prune -a -f

# Remove all unused networks
echo "Removing all unused networks..."
#sudo docker network prune -f

# Remove all unused volumes
echo "Removing all unused volumes..."
sudo docker volume prune -f

# Optional: Remove all dangling build cache
echo "Removing all dangling build cache..."
#sudo docker builder prune -a -f

echo "Docker data pruned successfully."