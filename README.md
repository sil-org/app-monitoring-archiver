# app-monitoring-archiver
Gets app monitoring results and saves them to Google Sheets

## Description
This app gets the previous month's uptime values for each NodePing check that are associated with
a particular contact group.

It then adds/inserts them into a Google sheet.

 - The month headings go from cell B2 to the right.
 - Each month's results are in a column starting at row 3.
 - The month columns do not get overwritten, just added to (inserted in chronological order).
 - The NodePing check names go from A3 down.
 - Each row has the results for one NodePing check (beginning at column B, one column per run of this app).
 - New rows for NodePing checks are inserted in alphabetical order.  (If the existing checks are out of order,
   they will not be corrected.)


## Setup

### NodePing
 - Ensure that all the NodePing Checks you want included in the Google Sheet
   have a notification set to the group you will be polling for (e.g. "MyTeam Alerts").

### Google API
 - Set up a Google API project and authentication credentials using a
service account by following the instuctions at https://flaviocopes.com/google-api-authentication/
 - Give that service account edit permissions on your Google Sheet.

 Note: There is a 100 writes per 100 seconds rate limit on Google Sheets.

### Set environment variables

```sh
$ export NODEPING_TOKEN=EG123ABC
$ export GOOGLE_AUTH_CLIENT_EMAIL=example@myaccount-123.iam.gserviceaccount.com
$ export GOOGLE_AUTH_PRIVATE_KEY_ID=abc123
$ export GOOGLE_AUTH_PRIVATE_KEY=-----BEGIN PRIVATE KEY-----\nMIIE...\n...\nabc=\n-----END PRIVATE KEY-----\n
$ export GOOGLE_AUTH_TOKEN_URI=https://oauth2.googleapis.com/token
```
Note that these GOOGLE_AUTH_* variable values can be found in the json that Google provides when
creating a services account.  Also note that the Go code itself will convert "\n" to EOL in the
GOOGLE_AUTH_PRIVATE_KEY value.


## AWS CDK

To build and deploy:

* Build the Go binary:

```sh
CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -o bin/bootstrap cmd/lambda/main.go
```

* Deploy using CDK:

```sh
docker compose run --rm cdk cdk deploy
```

### Run from command line

```sh
$ go run main.go run --help
$ go run main.go run -g "MyTeams Alerts" -s EG123ABC
```

The SPREADSHEET_ID is the middle part of the url for the target Google Sheet when you just browse to it.

## Terraform / OIDC

Terraform Cloud authenticates to AWS using OIDC — no static access keys are needed.

### Trust policy

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/app.terraform.io"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "app.terraform.io:aud": "aws.workload.identity"
        },
        "StringLike": {
          "app.terraform.io:sub": "organization:${org}:project:*:workspace:app-monitoring-archiver:run_phase:*"
        }
      }
    }
  ]
}
```

### Permissions policy (inline, named `Deployment`)

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "IAMUserManagement",
      "Effect": "Allow",
      "Action": [
        "iam:CreateUser",
        "iam:DeleteUser",
        "iam:GetUser",
        "iam:TagUser",
        "iam:UntagUser",
        "iam:ListUserTags",
        "iam:CreateAccessKey",
        "iam:DeleteAccessKey",
        "iam:ListAccessKeys",
        "iam:AttachUserPolicy",
        "iam:DetachUserPolicy",
        "iam:ListAttachedUserPolicies",
        "iam:ListUserPolicies"
      ],
      "Resource": "arn:aws:iam::${AWS_ACCOUNT_ID}:user/app-monitoring-archiver-*"
    },
    {
      "Sid": "IAMPolicyManagement",
      "Effect": "Allow",
      "Action": [
        "iam:CreatePolicy",
        "iam:DeletePolicy",
        "iam:GetPolicy",
        "iam:GetPolicyVersion",
        "iam:ListPolicyVersions",
        "iam:ListPolicyTags",
        "iam:CreatePolicyVersion",
        "iam:DeletePolicyVersion",
        "iam:TagPolicy",
        "iam:UntagPolicy"
      ],
      "Resource": "arn:aws:iam::${AWS_ACCOUNT_ID}:policy/app-monitoring-archiver-*"
    }
  ]
}
```

### TFC workspace environment variables

| Key | Value |
|---|---|
| `TFC_AWS_PROVIDER_AUTH` | `true` |
| `TFC_AWS_RUN_ROLE_ARN` | `arn:aws:iam::${AWS_ACCOUNT_ID}:role/app-monitoring-archiver-hcp-terraform` |
