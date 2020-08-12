# go-sync-gsheet
1. Create GCP project.
2. Get GCP OAuth 2.0 token, save to credentials.json.
3. Enable GCP Google Sheets API.
4. Create new spreadsheet, and setting spreadsheet id to config.yaml
5. Add AWS account for programmatic access, attach "AmazonEC2ReadOnlyAccess".
6. Add AWS auth info to config.yaml
7. Run docker, frist time, authorization required, then create token.json
8. Now syncing

example config.yaml
```yaml
Setting:
   SyncTimeInterval: 1m
   // https://docs.google.com/spreadsheets/d/{ID}/edit#gid=0
   SpreadsheetId: XXXXXXXXXXXXXXXXXXXXX
AWS:
  Auth:
    - Account: {AWS Account}
      AccessKey: {AWS AccessKey}
      SecretKey: {AWS SecretKey}
      Region: {AWS Region}
      Project: "AWS-Dev"

```

example docker run 
```
docker run -d --name sync-gsheet -v $(pwd)/config:/app kyos0109/go-sync-gsheet
```
