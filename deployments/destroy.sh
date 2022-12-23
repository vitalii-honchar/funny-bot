#!/usr/bin/env sh

instance_name='funny-bot'

instance_ids=$(aws ec2 describe-instances \
  --filters "Name=tag:Name,Values=$instance_name" "Name=instance-state-name,Values=running" \
  --query "Reservations[*].Instances[*].{InstanceID:InstanceId}" \
  --output text \
  --no-cli-pager | xargs)
aws ec2 terminate-instances --instance-ids "$instance_ids" --no-cli-pager