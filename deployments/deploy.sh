#!/usr/bin/env bash

instance_name='funny-bot'

destroy_previous_instance() {
  instance_ids=$(aws ec2 describe-instances \
    --filters "Name=tag:Name,Values=$instance_name" "Name=instance-state-name,Values=running" \
    --query "Reservations[*].Instances[*].{InstanceID:InstanceId}" \
    --output text \
    --no-cli-pager | xargs)
  aws ec2 terminate-instances --instance-ids "$instance_ids" --no-cli-pager > /dev/null
}

deploy_instance() {
  ami_id='ami-06143d8d167593285'
  instance_type='t3.nano'
  security_group_id='sg-0f67f523df045ab24'
  subnet_id='subnet-0319941ea41dabaf0'
  key_pair='test-key-pair'
  instance_id=$(aws ec2 run-instances \
    --image-id $ami_id \
    --count 1 \
    --instance-market-options "MarketType=spot,SpotOptions={SpotInstanceType=one-time}" \
    --instance-type $instance_type \
    --security-group-ids $security_group_id \
    --key-name $key_pair \
    --subnet-id $subnet_id \
    --output text \
    --query 'Instances[0].InstanceId' \
    --tag-specifications "ResourceType=instance,Tags=[{Key='Name',Value='$instance_name'}]" \
    --no-cli-pager)
  echo "$instance_id"
}

get_instance_public_ip() {
  ip_address=$(aws ec2 describe-instances \
    --instance-ids "$1" \
    --query "Reservations[0].Instances[0].PublicIpAddress" \
    --output text \
    --no-cli-pager)

  echo "$ip_address"
}

start_docker_image() {
  ssh -o StrictHostKeyChecking=no -i <(echo "$AWS_KEY") ubuntu@"$1" "docker run -e TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN -d weaxme/funny-bot:$VERSION"
}

destroy_previous_instance
instance_id=$(deploy_instance)
echo "Waiting for start instance $instance_id"
sleep 2m
echo "Initializing instance $instance_id"

ip_address=$(get_instance_public_ip "$instance_id")

start_docker_image "$ip_address"

aws ec2 describe-instances \
  --instance-ids "$instance_id" \
  --query "Reservations[*].Instances[*].{InstanceID:InstanceId,Name:Tags[?Key=='Name']|[0].Value,Status:State.Name,IP:PublicIpAddress}" \
  --output table --no-cli-pager
