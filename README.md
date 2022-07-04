## aws-ec2-assign-elastic-ip (go)

This is a clone of discobean/aws-ec2-assign-elastic-ip inspired from skymill/aws-ec2-assign-elastic-ip
except that:
1. It is written in go
2. Allows to select from a Pool of EIPs by using tag key/values

#### Usage (instanceid/region from metadata):
```$xslt
./aws-ec2-assign-elastic-ip-darwin-amd64 
    --tag-name Application 
    --tag-value minecraft 
```

#### Usage (when specifying instanceid/region):
```$xslt
./aws-ec2-assign-elastic-ip-darwin-amd64 
    --tag-name Application 
    --tag-value minecraft 
    --region ap-southeast-2 
    --instanceid i-0f0e97a20a05ce74b
```

#### Building
```$xslt
$ ./build.sh
$ ls -1 bin/
aws-ec2-assign-elastic-ip-darwin-amd64
aws-ec2-assign-elastic-ip-darwin-arm64
aws-ec2-assign-elastic-ip-linux-amd64
aws-ec2-assign-elastic-ip-linux-arm64
```

#### Instance permissions required
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "AllowDescribeAssociate",
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeAddresses",
                "ec2:AssociateAddress"
            ],
            "Resource": "*"
        }
    ]
}
```