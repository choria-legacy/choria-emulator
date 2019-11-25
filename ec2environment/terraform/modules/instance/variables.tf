variable "subnet_id" {
  description = "The subnet to place the instance in"
}

variable "type" {
  description = "Instance type to use"
  default     = "t2.large"
}

variable "security_group_ids" {
  description = "Security groups the instance should get applied"
  type        = list(string)
}

variable "user_data" {
  description = "User data to use"
}

variable "tags" {
  description = "Tags to apply to the resource"
  type        = map(string)
}

variable "instance_count" {
  description = "Number of instances to create"
  default     = 1
}

# product code for centos 7 on the marketplace from https://wiki.centos.org/Cloud/AWS
#
# set -e
# for region in $(aws ec2 describe-regions|jq -r '.Regions | map( .RegionName)|.[]')
# do
#   ami=$(aws ec2 describe-images --owners 'aws-marketplace' --filters 'Name=product-code,Values=aw0evgkw8e5c1q413zgy5pjce' --query 'sort_by(Images, &CreationDate)[-1].[ImageId]' --output 'text' --region $region)
#   echo "${region} = \"${ami}\""
# done
variable "amis" {
  description = "Base AMI to launch the instances with"
  default = {
    ap-northeast-1 = "ami-045f38c93733dd48d"
    ap-northeast-2 = "ami-06cf2a72dadf92410"
    ap-south-1     = "ami-02e60be79e78fef21"
    ap-southeast-1 = "ami-0b4dd9d65556cac22"
    ap-southeast-2 = "ami-08bd00d7713a39e7d"
    ca-central-1   = "ami-033e6106180a626d0"
    eu-central-1   = "ami-04cf43aca3e6f3de3"
    eu-north-1     = "ami-5ee66f20"
    eu-west-1      = "ami-0ff760d16d9497662"
    eu-west-2      = "ami-0eab3a90fc693af19"
    eu-west-3      = "ami-0e1ab783dc9489f34"
    sa-east-1      = "ami-0b8d86d4bf91850af"
    us-east-1      = "ami-02eac2c0129f6376b"
    us-east-2      = "ami-0f2b4fc905b0bd1f1"
    us-west-1      = "ami-074e2d6769f445be5"
    us-west-2      = "ami-01ed306a12b7d1c96"
  }
}
