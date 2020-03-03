## Creating AWS Batch Resources

### With a CloudFormation template

```
aws --region us-west-2 cloudformation create-stack --stack-name funnel-stack --template-body file://cloud-formation/aws-batch-resources.yaml --capabilities CAPABILITY_IAM
```

### On the command line with JSON documents

#### IAM Roles

*FunnelBatchServiceRole*

```
aws iam create-role --role-name FunnelBatchServiceRole --assume-role-policy-document file://json-resources/assume-role-policies/AWSBatchServiceRole.json
```

*FunnelEcsInstancerole*

```
aws iam create-role --role-name FunnelEcsInstanceRole --assume-role-policy-document file://json-resources/assume-role-policies/ecsInstanceRole.json
```

*FunnelEcsTaskRole*

```
aws iam create-role --role-name FunnelEcsTaskRole --assume-role-policy-document file://json-resources/assume-role-policies/FunnelEcsTaskRole.json
aws iam attach-role-policy --role-name FunnelEcsTaskRole --policy-name FunnelS3 --policy-document file://json-resources/policies/FunnelS3.json
aws iam attach-role-policy --role-name FunnelEcsTaskRole --policy-name FunnelDynamodb --policy-document file://json-resources/policies/FunnelDynamodb.json
```

#### Batch Resources

*funnel-compute-environment*

```
aws batch create-compute-environment --cli-input-json file://json-resources/ComputeEnvironment.json
```

*funnel-job-queue*

```
aws batch create-job-queue --cli-input-json file://json-resources/JobQueue.json
```

*funnel-job-definition*

```
aws batch register-job-definition --cli-input-json file://json-resources/JobDefinition.json
```
