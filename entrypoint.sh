#!/bin/bash
set -e

export EC2_AVAIL_ZONE=`curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone`
export EC2_REGION="`echo \"$EC2_AVAIL_ZONE\" | sed -e 's:\([0-9][0-9]*\)[a-z]*\$:\\1:'`"

# Inject SSM secrets
if [[ -n $SSM_SECRETS ]]
then
  for param in $SSM_SECRETS
    do
        secret=$(aws --region $EC2_REGION ssm get-parameters --names "$param" --query 'Parameters[*].[Value]' --output text --with-decryption)
        param_name=$(echo "$param" | cut -d "." -f2)
        export $param_name="$secret"
    done
fi

# Run application
exec "$@"