{
  FunctionName: '{{ must_env `SERVICE_NAME` }}',
  Environment: {
    Variables: {
      SLACK_APP_TOKEN: 'dummy',
      SLACK_BOT_TOKEN: 'dummy',
      BEDROCK_MODEL_ID: '{{ must_env `BEDROCK_MODEL_ID` }}',
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
