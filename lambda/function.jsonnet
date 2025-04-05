{
  FunctionName: '{{ must_env `SERVICE_NAME` }}',
  Environment: {
    Variables: {
    },
  },
  Handler: 'bootstarp.sh',
  MemorySize: 128,
  Role: '{{ must_env `ROLE_ARN` }}',
  Runtime: 'provided.al2023',
  Architectures: [
    'arm64',
  ],
  Tags: {
    Service: '{{ must_env `SERVICE_NAME` }}',
  },
  Timeout: 30,
}
