# Copyright (c) 2021 Andres More

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {}

locals {
  name = "api-${random_string.suffix.result}"
}

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  enable_dns_hostnames = true

  tags = {
    Name = local.name
  }
}

resource "aws_eip" "nat" {
  count      = 2
  vpc        = true
  depends_on = [aws_internet_gateway.main]

  tags = {
    Name = "${local.name}-${count.index}"
  }
}

resource "aws_nat_gateway" "main" {
  count         = 2
  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id
  depends_on    = [aws_internet_gateway.main]
  tags = {
    Name = "${local.name}-${count.index}"
  }
}

resource "aws_subnet" "public" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(aws_vpc.main.cidr_block, 8, count.index)
  availability_zone = element(data.aws_availability_zones.available.names, count.index)

  map_public_ip_on_launch = true

  tags = {
    Name = "${local.name}-public-${count.index}"
  }
}

resource "aws_subnet" "private" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(aws_vpc.main.cidr_block, 8, count.index + length(aws_subnet.public))
  availability_zone = element(data.aws_availability_zones.available.names, count.index)

  tags = {
    Name = "${local.name}-public-${count.index}"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = local.name
  }
}

resource "aws_route_table" "private" {
  count  = 2
  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main[count.index].id
  }

  tags = {
    Name = "${local.name}-private-${count.index}"
  }

}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "${local.name}-public"
  }
}

resource "aws_route_table_association" "public" {
  count          = 2
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  count          = 2
  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}

data "aws_ami" "latest" {
  most_recent = true

  owners = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-arm64-gp2"]
  }
}

output "ami" {
  value = data.aws_ami.latest.description
}

resource "aws_iam_role" "agent" {
  name = local.name
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "cw" {
  role       = aws_iam_role.agent.name
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"
}

resource "aws_iam_role_policy_attachment" "ssm" {
  role       = aws_iam_role.agent.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "agent" {
  role = aws_iam_role.agent.name
}

resource "aws_launch_template" "main" {
  name          = local.name
  instance_type = "t4g.nano"
  image_id      = data.aws_ami.latest.id
  user_data     = filebase64("${path.module}/user_data.sh")

  vpc_security_group_ids = [aws_security_group.asg.id]
  iam_instance_profile {
    name = aws_iam_instance_profile.agent.name
  }

}

resource "aws_autoscaling_group" "main" {
  name                = local.name
  vpc_zone_identifier = [aws_subnet.private[0].id, aws_subnet.private[1].id]
  target_group_arns   = [aws_lb_target_group.main.arn]

  max_size         = 8
  desired_capacity = 1
  min_size         = 1

  launch_template {
    id      = aws_launch_template.main.id
    version = "$Latest"
  }

  instance_refresh {
    strategy = "Rolling"
  }

  tag {
    key                 = "Name"
    value               = local.name
    propagate_at_launch = true
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group" "lb" {
  name   = "${local.name}-lb"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "asg" {
  name   = "${local.name}-asg"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port       = 8080
    to_port         = 8080
    protocol        = "TCP"
    security_groups = [aws_security_group.lb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb" "main" {
  name            = local.name
  security_groups = [aws_security_group.lb.id]
  subnets         = [aws_subnet.public[0].id, aws_subnet.public[1].id]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_lb_listener" "main" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.main.arn
  }
}

resource "aws_lb_target_group" "main" {
  vpc_id   = aws_vpc.main.id
  port     = 8080
  protocol = "HTTP"
  health_check {
    path = "/v1/health"
  }
}

output "dns" {
  value = aws_lb.main.dns_name
}

resource "aws_autoscaling_policy" "up" {
  name                   = "${local.name}-up"
  scaling_adjustment     = 1
  adjustment_type        = "ChangeInCapacity"
  autoscaling_group_name = aws_autoscaling_group.main.name
}

resource "aws_autoscaling_policy" "down" {
  name                   = "${local.name}-down"
  scaling_adjustment     = -1
  adjustment_type        = "ChangeInCapacity"
  autoscaling_group_name = aws_autoscaling_group.main.name
}

resource "aws_cloudwatch_metric_alarm" "high" {
  alarm_name          = "${local.name}-high"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EC2"
  threshold           = "75"
  statistic           = "Average"
  alarm_actions       = [aws_autoscaling_policy.up.arn]
  dimensions = {
    AutoScalingGroupName = aws_autoscaling_group.main.name
  }
  evaluation_periods = 3
  period             = 600
}

resource "aws_cloudwatch_metric_alarm" "low" {
  alarm_name          = "${local.name}-low"
  comparison_operator = "LessThanOrEqualToThreshold"
  metric_name         = "CPUUtilization"
  namespace           = "AWS/EC2"
  threshold           = "25"
  statistic           = "Average"
  alarm_actions       = [aws_autoscaling_policy.down.arn]
  dimensions = {
    AutoScalingGroupName = aws_autoscaling_group.main.name
  }
  evaluation_periods = 3
  period             = 600
}
