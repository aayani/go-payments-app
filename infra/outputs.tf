// Payments URL output - exposes the ALB endpoint - lets operators access the local app
output "payments_url" {
  value = "http://localhost:4566/_alb/${aws_lb.payments.name}/"
}

// RDS endpoint output - exposes the database address - supports direct administration and debugging
output "rds_endpoint" {
  value = aws_db_instance.payments.endpoint
}

// Redis endpoint output - exposes the cache address and port - supports direct inspection and clients
output "redis_endpoint" {
  value = {
    address = aws_elasticache_cluster.payments.cache_nodes[0].address
    port    = aws_elasticache_cluster.payments.cache_nodes[0].port
  }
}
