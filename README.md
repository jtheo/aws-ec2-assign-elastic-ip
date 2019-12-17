## aws-ec2-assign-elastic-ip (go)

This is inspired from skymill/aws-ec2-assign-elastic-ip
except that:
1. It is written in go
2. Allows to select from a Pool of EIPs by using tag key/values

#### Usage (instanceid/region from metadata):
```$xslt
./aws-ec2-assign-elastic-ip-darwin-amd64 
    --eiptagkey Application 
    --eiptagvalue minecraft 
```

#### Usage (when specifying instanceid/region):
```$xslt
./aws-ec2-assign-elastic-ip-darwin-amd64 
    --eiptagkey Application 
    --eiptagvalue minecraft 
    --region ap-southeast-2 
    --instanceid i-0f0e97a20a05ce74b
```

#### Building
```$xslt
$ make build
$ ls -1 build/*
build/aws-ec2-assign-elastic-ip-darwin-amd64
build/aws-ec2-assign-elastic-ip-linux-amd64
build/aws-ec2-assign-elastic-ip-linux-arm
build/aws-ec2-assign-elastic-ip-windows-amd64
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