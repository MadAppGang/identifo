package s3

// TODO: a great idea is to implement file change listening with SNS HTTP service
// it will user push scheme to send all config file updates
// https://aws.amazon.com/sns/features/
// the only question is to think how to handle horizontal scaling
// to notify all instances
