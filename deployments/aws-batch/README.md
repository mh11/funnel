## Default AWS Batch Funnel Resources

The JSON documents in this directory show examples of the resources that would be created by running `funnel aws batch create-all-resources --region us-west-2`.

### Create IAM Roles

*AWSBatchServiceRole*

```
aws iam create-role --role-name AWSBatchServiceRole --assume-role-policy-document file://assume-role-policies/AWSBatchServiceRole.json
```

*ecsInstancerole*

```
aws iam create-role --role-name ecsInstanceRole --assume-role-policy-document file://assume-role-policies/ecsInstanceRole.json
```

*FunnelEcsTaskRole*

```
aws iam create-role --role-name FunnelEcsTaskRole --assume-role-policy-document file://assume-role-policies/FunnelEcsTaskRole.json
aws iam attach-role-policy --role-name FunnelEcsTaskRole --policy-name FunnelS3 --policy-document file://policies/FunnelS3.json
aws iam attach-role-policy --role-name FunnelEcsTaskRole --policy-name FunnelDynamodb --policy-document file://policies/FunnelDynamodb.json
```

## Create AWS Resources
