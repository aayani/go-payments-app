// PostgreSQL database - stores payment accounts and transactions - provides durable application state
resource "aws_db_instance" "payments" {
  identifier        = "payments"
  engine            = "postgres"
  engine_version    = "18"
  instance_class    = "db.t3.micro"
  allocated_storage = 20
  db_name           = "payments"
  username          = "payments"
  password          = "payments"

  lifecycle {
    prevent_destroy = true
    ignore_changes  = [max_allocated_storage]
  }
}

// Redis cache - provides storage for initialized transactions - supports the transaction initiation flow
resource "aws_elasticache_cluster" "payments" {
  cluster_id                 = "payments-redis"
  engine                     = "redis"
  node_type                  = "cache.t3.micro"
  num_cache_nodes            = 1
  port                       = 6379
  auto_minor_version_upgrade = false
}

// CloudWatch log group - stores ECS container output for seven days - centralizes logs with bounded retention
resource "aws_cloudwatch_log_group" "payments" {
  name              = "/ecs/payments"
  retention_in_days = 7
}

// ECS execution role - grants ECS permission to publish container logs - supplies the identity used by awslogs
resource "aws_iam_role" "payments_task" {
  name = "payments-task"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Service = "ecs-tasks.amazonaws.com"
      }
      Action = "sts:AssumeRole"
    }]
  })
}

// ECS log policy - allows log streams and events in the payments group - limits writes to this app's logs
resource "aws_iam_role_policy" "payments_task" {
  name = "payments-logs"
  role = aws_iam_role.payments_task.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["logs:CreateLogStream", "logs:PutLogEvents"]
      Resource = aws_cloudwatch_log_group.payments.arn
    }]
  })
}

// ECS cluster - groups the payments service and its tasks - provides the ECS deployment boundary
resource "aws_ecs_cluster" "payments" {
  name = "payments"
}

// ALB security group - allows public HTTP access and outbound target traffic - exposes only the listener publicly
resource "aws_security_group" "alb" {
  name        = "payments-alb"
  description = "Allow HTTP traffic to the payments ALB"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

// App security group - allows the ALB to reach the app on port 8080 - prevents direct public app access
resource "aws_security_group" "app" {
  name        = "payments-app"
  description = "Allow ALB traffic to the payments app"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

// Application load balancer - exposes the payments service over HTTP - provides a stable public entry point
resource "aws_lb" "payments" {
  name               = "payments"
  internal           = false
  load_balancer_type = "application"
  subnets            = data.aws_subnets.default.ids
  security_groups    = [aws_security_group.alb.id]
}

// ALB target group - routes requests to healthy app tasks on port 8080 - keeps unhealthy tasks out of rotation
resource "aws_lb_target_group" "payments" {
  name        = "payments"
  port        = 8080
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = data.aws_vpc.default.id

  health_check {
    path    = "/ping"
    matcher = "200"
  }
}

// ALB listener - forwards HTTP requests to the payments target group - connects incoming traffic to ECS
resource "aws_lb_listener" "payments" {
  load_balancer_arn = aws_lb.payments.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.payments.arn
  }
}

// ECS task definition - configures the payments container, networking, and logs - defines what Fargate runs
resource "aws_ecs_task_definition" "payments" {
  family                   = "payments"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.payments_task.arn

  container_definitions = jsonencode([{
    name      = "payments"
    image     = var.app_image
    essential = true

    portMappings = [{
      containerPort = 8080
      protocol      = "tcp"
    }]

    environment = [
      { name = "PORT", value = "8080" },
      {
        name  = "DATABASE_URL"
        value = "postgres://payments:payments@host.docker.internal:${aws_db_instance.payments.port}/payments?sslmode=disable"
      },
      {
        name  = "REDIS_URL"
        value = "redis://${aws_elasticache_cluster.payments.cache_nodes[0].address}:${aws_elasticache_cluster.payments.cache_nodes[0].port}/0"
      }
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.payments.name
        "awslogs-region"        = "us-east-1"
        "awslogs-stream-prefix" = "ecs"
      }
    }
  }])
}

// ECS service - keeps one payments task running behind the ALB - minimizes local container use
resource "aws_ecs_service" "payments" {
  name                 = "payments"
  cluster              = aws_ecs_cluster.payments.id
  task_definition      = aws_ecs_task_definition.payments.arn
  launch_type          = "FARGATE"
  desired_count        = 1
  force_new_deployment = true

  triggers = {
    redeployment = plantimestamp()
  }

  network_configuration {
    subnets          = data.aws_subnets.default.ids
    security_groups  = [aws_security_group.app.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.payments.arn
    container_name   = "payments"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener.payments, aws_iam_role_policy.payments_task]
}
